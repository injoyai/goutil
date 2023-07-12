package main

import (
	"fmt"
	_ "github.com/DrmagicE/gmqtt/persistence"
	_ "github.com/DrmagicE/gmqtt/topicalias/fifo"
	"github.com/injoyai/base/oss"
	"github.com/injoyai/base/oss/shell"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/goutil/cmd/crud"
	"github.com/injoyai/goutil/net/ip"
	oss2 "github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/other/notice/voice"
	"github.com/injoyai/goutil/string/bar"
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial/proxy"
	"github.com/injoyai/logs"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

func handleVersion(cmd *cobra.Command, args []string, flags *Flags) {
	fmt.Println(version)
	fmt.Println(details)
}

func handlerSwag(cmd *cobra.Command, args []string, flags *Flags) {
	param := []string{"init"}
	flags.Range(func(key string, val *Flag) bool {
		param = append(param, fmt.Sprintf(" -%s %s", val.Short, val.Value))
		return true
	})
	bs, _ := handlerShell("swag", append(param, args...)...)
	fmt.Println(bs)
}

func handleBuild(cmd *cobra.Command, args []string, flags *Flags) {
	os.Setenv("GOOS", "windows")
	os.Setenv("GOARCH", "amd64")
	os.Setenv("GO111MODULE", "on")
	list := append([]string{"go", "build"}, args...)
	result, _ := handlerShell(strings.Join(list, " "))
	fmt.Println(result)
}

func handlerInstall(cmd *cobra.Command, args []string, flags *Flags) {
	if len(args) == 0 {
		fmt.Println("请输入需要安装的应用")
		return
	}
	switch args[0] {

	case "in":

		url := "https://github.com/injoyai/goutil/raw/main/cmd/in.exe"
		logs.PrintErr(bar.Download(url, "./in.exe"))

	case "upgrade":

		logs.PrintErr(oss.New("./in_upgrade.exe", upgrade))

	case "upx":

		logs.PrintErr(oss.New("./upx.exe", upx))

	case "rsrc":

		logs.PrintErr(oss.New("./rsrc.exe", rsrc))

	case "chromedriver":

		if _, err := installChromedriver(oss2.UserDefaultDir(), flags.GetBool("download")); err != nil {
			log.Printf("[错误] %s", err.Error())
		}

	case "downloader":

		url := "https://github.com/injoyai/downloader/releases/latest/download/downloader.exe"
		logs.PrintErr(bar.Download(url, "./downloader.exe"))

	case "swag":

		logs.PrintErr(oss.New("./swag.exe", swag))

	case "hfs":

		logs.PrintErr(oss.New("./hfs.exe", hfs))

	case "influxdb":

		url := "https://dl.influxdata.com/influxdb/releases/influxdb2-2.7.1-windows-amd64.zip"
		logs.PrintErr(bar.Download(url, "./influxdb.zip"))
		logs.PrintErr(DecodeZIP("./influxdb.zip", "./"))
		os.Remove("./influxdb.zip")

	default:

		bs, _ := exec.Command("go", "install", args[0]).CombinedOutput()
		fmt.Println(string(bs))

	}
}

func handlerGo(cmd *cobra.Command, args []string, flags *Flags) {
	bs, _ := exec.Command("go", args...).CombinedOutput()
	fmt.Println(string(bs))
}

func handlerPprof(cmd *cobra.Command, args []string, flags *Flags) {
	if len(args) == 0 {
		fmt.Println("输入地址,例: http://localhost:6060 , localhost:6060")
		return
	}
	switch cmd.Use {
	case "profile":
		fmt.Println("正在读取数据,需要20秒...")
		handlerPprof2(args[0] + "/pprof/profile?seconds=20")
	case "heap":
		handlerPprof2(args[0] + "/pprof/heap")
	}
}

func handlerPprof2(url string, param ...string) {
	if !strings.Contains(url, "http://") {
		url = "http://" + url
	}
	param = append(param, url)
	param = append([]string{"tool", "pprof"}, param...)
	result, _ := handlerShell("go", param...)
	fmt.Println(result)
}

func handlerShell(name string, args ...string) (string, error) {
	bs, err := exec.Command(name, args...).CombinedOutput()
	return string(bs), err
}

func handlerCrud(cmd *cobra.Command, args []string, flags *Flags) {
	if len(args) == 0 {
		log.Printf("[错误] %s", "请输入模块名称 例: in curd test")
	}
	logs.PrintErr(crud.New(args[0]))
}

func handlerNow(cmd *cobra.Command, args []string, flags *Flags) {
	fmt.Println(time.Now().String())
}

func handlerSpeak(cmd *cobra.Command, args []string, flags *Flags) {
	msg := fmt.Sprint(conv.Interfaces(args)...)
	voice.Speak(msg)
}

func handlerProxy(cmd *cobra.Command, args []string, flags *Flags) {
	//服务端地址
	serverAddr := cfg.GetString("addr", flags.GetString("serverAddr"))
	if runtime.GOOS == "windows" && len(serverAddr) == 0 {
		fmt.Println("请输入服务地址(默认121.36.99.197:9000):")
		fmt.Scanln(&serverAddr)
		if len(serverAddr) == 0 {
			serverAddr = "121.36.99.197:9000"
		}
	}

	//客户端唯一标识
	sn := cfg.GetString("sn", flags.GetString("sn"))
	if runtime.GOOS == "windows" && len(sn) == 0 {
		fmt.Println("请输入SN(默认cmd):")
		fmt.Scanln(&sn)
		if len(sn) == 0 {
			sn = "cmd"
		}
	}

	//代理地址
	proxyAddr := flags.GetString("proxyAddr")
	if runtime.GOOS == "windows" && len(proxyAddr) == 0 {
		fmt.Println("请输入代理地址(默认代理全部):")
		fmt.Scanln(&proxyAddr)
	}

	c := proxy.NewPortForwardingClient(serverAddr, sn, func(c *io.Client, e *proxy.Entity) {
		c.SetPrintWithBase()
		c.Debug()
		if len(proxyAddr) > 0 {
			e.SetWriteFunc(func(msg *proxy.Message) (*proxy.Message, error) {
				msg.Addr = proxyAddr
				return msg, nil
			})
		}
	})
	c.Run()
	select {}
}

func handlerScan(cmd *cobra.Command, args []string, flags *Flags) {
	switch true {
	case len(args) == 0:
		log.Println("[错误]", "缺少扫描类型(icmp,serial...)")
	default:

		number := flags.GetInt("number")

		switch args[0] {
		case "icmp":

			gateIPv4 := []byte(net.ParseIP(ip.GetLocal())[12:15])
			wg := sync.WaitGroup{}
			for i := conv.Uint32(append(gateIPv4, 0)); i <= conv.Uint32(append(gateIPv4, 255)); i++ {
				ipv4 := net.IPv4(uint8(i>>24), uint8(i>>16), uint8(i>>8), uint8(i))
				wg.Add(1)
				go func(ipv4 net.IP) {
					defer wg.Done()
					used, err := ip.Ping(ipv4.String(), time.Second)
					if err == nil {
						fmt.Printf("%s: %s\n", ipv4, used.String())
					}
				}(ipv4)
			}
			wg.Wait()

		case "ssh":

			gateIPv4 := []byte(net.ParseIP(ip.GetLocal())[12:15])
			wg := sync.WaitGroup{}
			for i := conv.Uint32(append(gateIPv4, 0)); i <= conv.Uint32(append(gateIPv4, 255)); i++ {
				ipv4 := net.IPv4(uint8(i>>24), uint8(i>>16), uint8(i>>8), uint8(i))
				wg.Add(1)
				go func(ipv4 net.IP) {
					defer wg.Done()
					c, err := net.Dial("tcp", ipv4.String()+":22")
					if err == nil {
						c.Close()
						fmt.Printf("%s\n", ipv4)
					}
				}(ipv4)
			}
			wg.Wait()

		case "serial":

			list, err := serial.GetPortsList()
			if err != nil {
				logs.Err(err)
				return
			}
			fmt.Println(strings.Join(list, "\n"))

		case "edge":

			ipv4 := ip.GetLocal()
			startIP := append(net.ParseIP(ipv4)[:15], 0)
			endIP := append(net.ParseIP(ipv4)[:15], 255)
			ch, ctx := handlerScanEdge(startIP, endIP)
			for i := 1; ; i++ {
				select {
				case <-ctx.Done():
					return
				case data := <-ch:
					fmt.Printf("%v: %v\n", data.IP, data.SN)
					if flags.GetBool("open") {
						logs.PrintErr(shell.OpenBrowser(fmt.Sprintf("http://%s:10001", data.IP)))
					}
					if number > 0 && i >= number {
						return
					}
				}
			}

		}
	}
}

func handlerDemo(name string, bs []byte) func(cmd *cobra.Command, args []string, flags *Flags) {
	return func(cmd *cobra.Command, args []string, flags *Flags) {
		oss.New(name, bs)
		fmt.Println("success")
	}
}

func handlerOpen(cmd *cobra.Command, args []string, flags *Flags) {
	if len(args) == 0 {
		fmt.Printf("请输入参数,例(in open hosts)")
		return
	}
	switch strings.ToLower(args[0]) {
	case "hosts":
		shell.Start("C:\\Windows\\System32\\drivers\\etc\\hosts")
	case "injoy":
		shell.Start(oss2.UserDefaultDir())
	case "appdata":
		shell.Start(oss2.UserDataDir())
	case "startup":
		shell.Start(oss2.UserStartupDir())
	default:
		shell.Start(args[0])
	}
}
