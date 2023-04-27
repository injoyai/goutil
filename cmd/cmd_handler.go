package main

import (
	"context"
	"fmt"
	"github.com/DrmagicE/gmqtt"
	_ "github.com/DrmagicE/gmqtt/persistence"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/server"
	_ "github.com/DrmagicE/gmqtt/topicalias/fifo"
	"github.com/injoyai/base/oss"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/goutil/cmd/crud"
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial"
	"github.com/injoyai/io/dial/proxy"
	"github.com/injoyai/logs"
	"github.com/spf13/cobra"
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
	fmt.Println("v1.0.0")
}

func handlerPing(cmd *cobra.Command, args []string, flags *Flags) {
	fmt.Println("pong")
}

type Swag struct {
	g string
	d string
	p string
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

	case "downloader":

		url := "https://github.com/injoyai/downloader/releases/latest/download/downloader.exe"
		resp, err := http.Get(url)
		if err != nil {
			logs.Err(err)
			return
		}
		defer resp.Body.Close()
		logs.PrintErr(oss.New("./downloader.exe", resp.Body))

	case "swag":

		logs.PrintErr(oss.New("./swag.exe", swag))

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
