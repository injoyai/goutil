package main

import (
	"context"
	"fmt"
	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
	"github.com/DrmagicE/gmqtt/server"
	"github.com/injoyai/base/oss"
	"github.com/injoyai/base/oss/shell"
	"github.com/injoyai/goutil/string/bar"
	"github.com/injoyai/io/dial"
	"github.com/injoyai/logs"
	"github.com/spf13/cobra"
	"github.com/tebeka/selenium"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

//====================SeleniumServer====================//

func handlerSeleniumServer(cmd *cobra.Command, args []string, flags *Flags) {
	port := flags.GetInt("port")
	selenium.SetDebug(flags.GetBool("debug"))
	ser, err := selenium.NewChromeDriverService(flags.GetString("chromedriver"), port)
	if err != nil {
		logs.Err(err)
		return
	}
	defer ser.Stop()
	log.Printf("[%d] 开启驱动成功\n", port)
	select {}
}

//====================TCPServer====================//

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

//====================MQTTServer====================//

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

//====================EdgeServer====================//

func handlerEdgeServer(cmd *cobra.Command, args []string, flags *Flags) {
	userDir, _ := oss.UserHome()
	userDir = filepath.Join(userDir, "AppData/Local/injoy")
	os.MkdirAll(userDir, 0666)
	filename := filepath.Join(userDir, "edge.exe")
	if !oss.Exists(filename) || flags.GetBool("download") {
		for logs.PrintErr(bar.Download("http://192.168.10.102:8888/gateway/edge/-/raw/v1.0.12(%E5%90%88%E5%B9%B6%E5%88%86%E6%94%AF%E7%89%88%E6%9C%AC)/bin/windows/edge.exe", filename)) {
			<-time.After(time.Second)
		}
	}
	shell.Start(filename)
}
