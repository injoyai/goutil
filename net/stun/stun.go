package stun

import (
	"github.com/pion/stun"
)

var DefaultServer = "stun.l.google.com:19302"

// GetNetAddr 获取网络地址,可自定义stun服务器
func GetNetAddr(serverAddr ...string) (addr stun.XORMappedAddress, err error) {

	serAddr := DefaultServer
	if len(serverAddr) > 0 {
		serAddr = serverAddr[0]
	}

	// 创建一个 STUN 客户端连接到公共 STUN 服务器
	conn, err := stun.Dial("udp", serAddr)
	if err != nil {
		return stun.XORMappedAddress{}, err
	}
	defer conn.Close()

	// 创建一个 STUN 请求
	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	// 发送请求并接收响应
	if err := conn.Do(message, func(res stun.Event) {
		if res.Error != nil {
			err = res.Error
			return
		}
		// 解析 STUN 响应中的外网地址
		err = addr.GetFrom(res.Message)
	}); err != nil {
		return addr, err
	}

	return
}
