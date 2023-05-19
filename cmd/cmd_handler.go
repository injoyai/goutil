package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/DrmagicE/gmqtt"
	_ "github.com/DrmagicE/gmqtt/persistence"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/server"
	_ "github.com/DrmagicE/gmqtt/topicalias/fifo"
	"github.com/fatih/color"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/injoyai/base/oss"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/goutil/cmd/crud"
	"github.com/injoyai/goutil/string/bar"
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial"
	"github.com/injoyai/io/dial/proxy"
	"github.com/injoyai/logs"
	"github.com/spf13/cobra"
	"github.com/tebeka/selenium"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func handleVersion(cmd *cobra.Command, args []string, flags *Flags) {
	fmt.Println("版本: v1.0.0")
	fmt.Println("系统: ", runtime.GOOS)
	dir, _ := os.Getwd()
	fmt.Println("路径: ", dir)
	fmt.Println("GO版本: ", runtime.Version())
	fmt.Println("GOROOT: ", runtime.GOROOT())
	_, file, _, _ := runtime.Caller(1)
	fmt.Println("GOROOT: ", file)
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
	case "all":

		//todo 安装全部

	case "upx":

		logs.PrintErr(oss.New("./upx.exe", upx))

	case "rsrc":

		logs.PrintErr(oss.New("./rsrc.exe", rsrc))

	case "chromedriver":

		logs.Debug("未实现")

	case "downloader":

		url := "https://github.com/injoyai/downloader/releases/latest/download/downloader.exe"
		resp, err := http.Get(url)
		if err != nil {
			logs.Err(err)
			return
		}
		defer resp.Body.Close()
		buff := bufio.NewReader(resp.Body)

		f, err := os.Create("./downloader.exe")
		if err != nil {
			logs.Err(err)
			return
		}
		defer f.Close()
		b := bar.New()
		b.SetTotalSize(conv.Float64(resp.Header.Get("Content-Length")))
		if ca := flags.GetInt("color"); ca > 0 {
			b.SetColor(color.Attribute(ca))
		}
		go b.Wait()

		for {
			buf := make([]byte, 1<<20)
			n, err := buff.Read(buf)
			if err != nil {
				if err == io.EOF {
					return
				}
				logs.Err(err)
				return
			}
			b.Add(float64(n))
			if _, err := f.Write(buf[:n]); err != nil {
				logs.Err(err)
				return
			}
		}

	case "swag":

		logs.PrintErr(oss.New("./swag.exe", swag))

	case "influxdb":

		oss.New("./influxdb.temp", influxdb)
		defer oss.Remove("./influxdb.temp")
		logs.PrintErr(DecodeZIP("./influxdb.temp", "./"))

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

//====================MQTTServer====================

func handlerTCPServer(cmd *cobra.Command, args []string, flags *Flags) {
	port := flags.GetInt("port", 10086)
	s, err := dial.NewTCPServer(port)
	if err != nil {
		log.Printf("[错误] %s", err.Error())
		return
	}
	s.Debug(flags.GetBool("debug"))
	s.Run()
}

func handlerMQTTServer(cmd *cobra.Command, args []string, flags *Flags) {

	port := flags.GetInt("port", 1883)
	debug := flags.GetBool("debug")

	fmt.Printf("ERROR:%v", func() error {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return err
		}
		srv := server.New(server.WithTCPListener(ln))
		if err := srv.Init(server.WithHook(server.Hooks{
			OnConnected: func(ctx context.Context, client server.Client) {
				if debug {
					log.Printf("新的客户端连接:%s", client.ClientOptions().ClientID)
				}
				srv.SubscriptionService().Subscribe(client.ClientOptions().ClientID, &gmqtt.Subscription{
					TopicFilter: client.ClientOptions().ClientID,
					QoS:         packets.Qos0,
				})
			},
			OnMsgArrived: func(ctx context.Context, client server.Client, req *server.MsgArrivedRequest) error {
				if debug {
					log.Printf("发布主题:%s,消息内容:%s", req.Message.Topic, string(req.Message.Payload))
				}
				return nil
			},
		})); err != nil {
			return err
		}
		log.Printf("[信息][:%d] 开启MQTT服务成功...\n", port)
		if err := srv.Run(); err != nil {
			return err
		}
		return nil
	}())
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
	logs.PrintErr(func(msg string) (err error) {
		if err := ole.CoInitialize(0); err != nil {
			return err
		}
		defer ole.CoUninitialize()
		unknown, err := oleutil.CreateObject("SAPI.SpVoice")
		if err != nil {
			return err
		}
		voice, err := unknown.QueryInterface(ole.IID_IDispatch)
		if err != nil {
			return err
		}
		defer voice.Release()
		_, err = oleutil.PutProperty(voice, "Rate", flags.GetInt("rate"))
		if err != nil {
			return err
		}
		_, err = oleutil.PutProperty(voice, "Volume", flags.GetInt("volume", 100))
		if err != nil {
			return err
		}
		_, err = oleutil.CallMethod(voice, "Speak", msg)
		if err != nil {
			return err
		}
		_, err = oleutil.CallMethod(voice, "WaitUntilDone", 0)
		if err != nil {
			return err
		}
		return nil
	}(fmt.Sprint(conv.Interfaces(args)...)))
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

	c := proxy.NewPortForwardingClient(serverAddr, sn, func(ctx context.Context, c *io.Client, e *proxy.Entity) {
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

func handlerSeleniumServer(cmd *cobra.Command, args []string, flags *Flags) {
	port := flags.GetInt("port")
	selenium.SetDebug(flags.GetBool("debug"))
	ser, err := selenium.NewChromeDriverService(flags.GetString("chromedriver"), port)
	if err != nil {
		logs.Err(err)
		return
	}
	defer ser.Stop()
	logs.Debugf("[%d] 开启驱动成功", port)
	select {}
}

func handlerDial(cmd *cobra.Command, args []string, flags *Flags) {
	switch true {
	case len(args) < 1:
		logs.Err("无效连接类型(tcp,serial...)")
	case len(args) < 2:
		logs.Err("无效连接地址")
	default:
		r := bufio.NewReader(os.Stdin)
		op := func(ctx context.Context, c *io.Client) {
			c.Debug()
			if !flags.GetBool("redial") {
				c.SetRedialWithNil()
			}
			go func(ctx context.Context) {
				for {
					select {
					case <-ctx.Done():
					default:
						bs, _, err := r.ReadLine()
						logs.PrintErr(err)
						msg := string(bs)
						if len(msg) > 2 && msg[0] == '0' && (msg[1] == 'x' || msg[1] == 'X') {
							_, err := c.WriteHEX(msg[2:])
							logs.PrintErr(err)
						} else {
							_, err := c.WriteASCII(msg)
							logs.PrintErr(err)
						}
					}
				}
			}(ctx)
		}
		switch args[0] {
		case "serial":
			c := dial.RedialSerial(&dial.SerialConfig{
				Address:  args[1],
				BaudRate: flags.GetInt("baudRate"),
				DataBits: flags.GetInt("dataBits"),
				StopBits: flags.GetInt("stopBits"),
				Parity:   flags.GetString("parity"),
				Timeout:  0,
			}, op)
			defer c.Close()
			oss.ListenExit(func() { c.CloseAll() })
		case "websocket", "ws":
			dial.RedialWebsocket(args[1], nil, op)
		default:
			dial.RedialTCP(args[1], op)
		}
		select {}
	}
}
