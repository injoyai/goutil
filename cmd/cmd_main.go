package main

import (
	"github.com/injoyai/logs"
	"github.com/spf13/cobra"
)

func init() {

}

func main() {

	logs.DefaultErr.SetWriter(logs.Stdout, logs.Trunk)
	logs.SetShowColor(false)

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
			Short:   "查看版本",
			Example: "in version",
			Run:     handleVersion,
		},

		&Command{
			Use:     "upgrade",
			Short:   "自我升级",
			Example: "in upgrade",
			Run:     handlerUpgrade,
		},

		&Command{
			Flag: []*Flag{
				{Name: "g", Short: "g"},
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
			Flag: []*Flag{
				{Name: "color", Short: "c", Memo: "日志颜色"},
			},
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
			Example: "in heap localhost:6060",
			Run:     handlerPprof,
		},

		&Command{
			Use:     "profile",
			Short:   "profile",
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
			Run:   handlerTCPServer,
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
			Flag: []*Flag{
				{Name: "rate", Short: "r", Memo: "语速"},
				{Name: "volume", Short: "v", DefValue: "100", Memo: "音量"},
			},
			Use:     "now",
			Short:   "当前时间",
			Example: "in now",
			Run:     handlerNow,
		},

		&Command{
			Flag: []*Flag{
				{Name: "rate", Short: "r", DefValue: "", Memo: "语速"},
				{Name: "volume", Short: "v", DefValue: "100", Memo: "音量"},
			},
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

		&Command{
			Flag: []*Flag{
				{Name: "port", Short: "p", Memo: "监听端口", DefValue: "20165"},
				{Name: "chromedriver", Short: "c", Memo: "驱动路径", DefValue: "./chromedriver.exe"},
				{Name: "debug", Short: "d", Memo: "打印日志", DefValue: "true"},
			},
			Use:     "seleniumServer",
			Short:   "自动化服务",
			Example: "in seleniumServer",
			Run:     handlerSeleniumServer,
		},

		&Command{
			Flag: []*Flag{
				{Name: "redial", Short: "r", Memo: "自动重连", DefValue: "true"},
				{Name: "debug", Short: "d", Memo: "打印日志", DefValue: "true"},
				{Name: "timeout", Short: "t", Memo: "超时时间(ms)", DefValue: "500"},

				{Name: "baudRate", Memo: "波特率", DefValue: "9600"},
				{Name: "dataBits", Memo: "数据位", DefValue: "8"},
				{Name: "stopBits", Memo: "停止位", DefValue: "1"},
				{Name: "parity", Memo: "校验", DefValue: "N"},
			},
			Use:     "dial",
			Short:   "连接",
			Example: "in dial tcp 127.0.0.1:80 -r false",
			Run:     handlerDial,
		},

		&Command{
			Flag: []*Flag{
				{Name: "number", Short: "n", Memo: "扫描数量", DefValue: "-1"},
			},
			Use:     "scan",
			Short:   "扫描",
			Example: "in scan icmp",
			Run:     handlerScan,
		},

		&Command{
			Use: "demo",
			Child: []*Command{
				{
					Use:   "build",
					Short: "build.sh",
					Run:   handlerDemo("./build.sh", build),
				},
				{
					Use:   "dockerfile",
					Short: "dockerfile",
					Run:   handlerDemo("./Dockerfile", dockerfile),
				},
				{
					Use:   "service",
					Short: "service",
					Run:   handlerDemo("./service.service", service),
				},
				{
					Use:   "install_minio",
					Short: "install_minio",
					Run:   handlerDemo("./install_minio.sh", installMinio),
				},
				{
					Use:   "install_nodered",
					Short: "install_nodered",
					Run:   handlerDemo("./install_nodered.sh", installNodeRed),
				},
			},
		},
	)

	root.Execute()

}
