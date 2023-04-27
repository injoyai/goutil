package main

import (
	"github.com/spf13/cobra"
)

func main() {

	root := &cobra.Command{
		Use:   "in",
		Short: "Cli工具",
		Long:  "Cli工具",
	}

	addCommand := func(cmd ...ICommand) {
		for _, v := range cmd {
			root.AddCommand(v.command())
		}
	}

	addCommand(

		&Command{
			Use:     "version",
			Short:   "version",
			Long:    "查看版本",
			Example: "in version",
			Run:     handleVersion,
		},

		&Command{
			Flag: []*Flag{
				{Name: "g", Short: "g", DefValue: "", Memo: ""},
			},
			Use:     "swag",
			Short:   "swag",
			Long:    "生成swagger文档",
			Example: "in swag -g /cmd/main.go",
			Run:     handlerSwag,
		},

		&Command{
			Use:   "build",
			Short: "build",
			Long:  "编译go文件",
			Run:   handleBuild,
		},

		&Command{
			Use:     "install",
			Short:   "install",
			Long:    "安装应用",
			Example: "in install github.com/xxx/xxx",
			Run:     handlerInstall,
		},

		&Command{
			Use:     "go",
			Short:   "go",
			Long:    "go cmd",
			Example: "in go version",
			Run:     handlerGo,
		},

		&Command{
			Use:     "heap",
			Short:   "heap",
			Long:    "尝试",
			Example: "in heap localhost:6060",
			Run:     handlerPprof,
		},

		&Command{
			Use:     "profile",
			Short:   "profile",
			Long:    "尝试",
			Example: "in profile localhost:6060",
			Run:     handlerPprof,
		},

		&Command{
			Flag: []*Flag{
				{Name: "port", Short: "p", DefValue: "10086", Memo: "监听端口"},
				{Name: "debug", Short: "d", DefValue: "true", Memo: "打印日志"},
			},
			Use:   "tcpServer",
			Short: "tcp server",
			Run:   handlerPprof,
		},

		&Command{
			Flag: []*Flag{
				{Name: "port", Short: "p", DefValue: "1883", Memo: "监听端口"},
				{Name: "debug", Short: "d", DefValue: "true", Memo: "打印日志"},
			},
			Command: &cobra.Command{
				Use:     "mqttServer",
				Short:   "mqtt server",
				Example: "in mqttServer -p 1883",
			},
			Run: handlerMQTTServer,
		},

		&Command{
			Use:     "crud",
			Short:   "生成增删改查",
			Example: "in curd test",
			Run:     handlerCrud,
		},

		&Command{
			Use:     "now",
			Short:   "当前时间",
			Example: "in now",
			Run:     handlerNow,
		},

		&Command{
			Use:     "speak",
			Short:   "文字转语音",
			Example: "in speak 哈哈哈",
			Run:     handlerSpeak,
		},

		&Command{
			Flag: []*Flag{
				{Name: "serverAddr", Short: "s", Memo: "服务地址"},
				{Name: "sn", Short: "k", Memo: "客户端标识"},
				{Name: "proxyAddr", Short: "p", Memo: "代理地址"},
			},
			Use:     "proxy",
			Short:   "代理",
			Example: "in proxy",
			Run:     handlerProxy,
		},
	)

	root.Execute()

}
