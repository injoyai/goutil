package ec800e

import (
	"errors"
	"fmt"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/conv"
	"strings"
)

/*

打开
+QIOPEN=contextID(1~15),connectID(0~11),"TCP","远程地址",远程端口,本地端口,模式
contextID: 场景ID
模式:
	0是缓存模式(需要主动去读数据),
	1是直吐模式(收到的数据会发送至串口),
	2透传模式(只能一个,透传和串口的所有数据)

先返回收到命令
OK
再收到连接结果
+QIOPEN: connectID(0~11),result(0成功,错误信息)

发送数据
AT+QISEND=connectID,数据长度
等收到>后再发送数据
响应
SEND OK

被关闭
+QIURC: "closed",connectID

需要再进行关闭操作
AT+QICLOSE=connectID

除非设置了自动关闭
AT+QICFG="passiveclosed",1(1标识被关闭后,主动关闭客户端)

*/

type TCPOpenResult string

func (this TCPOpenResult) Err() error {
	if this == "0" {
		return nil
	}
	return this
}

func (this TCPOpenResult) Error() string {
	switch this {
	case "0":
		return "成功"
	case "550":
		return "未知错误"
	case "551":
		return "操作受阻"
	case "552":
		return "无效参数"
	case "553":
		return "内存不足"
	case "554":
		return "创建socket失败"
	case "555":
		return "操作不支持"
	case "556":
		return "socket绑定失败"
	case "557":
		return "socket监听失败"
	case "558":
		return "socket写入失败"
	case "559":
		return "socket读取失败"
	case "560":
		return "socket接受失败"
	case "561":
		return "打开PDP场景失败"
	case "562":
		return "关闭PDP场景失败"
	case "563":
		return "socket标识被占用"
	case "564":
		return "DNS忙碌"
	case "565":
		return "DNS解析失败"
	default:
		return "未知错误"
	}
}

// TCPDial 建立TCP连接 127.0.0.1:10086
func (this *AT) TCPDial(num int8, addr string) error {

	s := strings.ReplaceAll(addr, ":", `",`)
	//contextID 0~15
	//connectID 0~11
	//model 0是缓存模式(需要主动去读数据),1是直吐模式(收到的数据会发送至串口),2透传模式(只能一个,透传和串口的所有数据)
	//参数分别是contextID,connectID,remoteIP,remotePort,localPort,model(2是透传)
	command := fmt.Sprintf(`+QIOPEN=1,%d,"TCP","%s,0,1`, num, s)
	_, err := this.send2(command)
	if err != nil {
		return err
	}

	//等待连接成功回调
	result, err := this.wait.Wait("+QIOPEN: " + conv.String(num))
	if err != nil {
		return err
	}
	if err := TCPOpenResult(result.(string)).Err(); err != nil {
		return err
	}

	return nil

}

func (this *AT) TCPClose(connectID int8) error {
	command := "+QICLOSE=" + conv.String(connectID)
	_, err := this.send2(command)
	return err
}

func (this *AT) TCPWrite(connectID int8, p []byte) error {

	//拆分数据为小段
	size := 1440
	for i := 0; i < len(p); i += size {
		bs := p[i:]
		if len(bs) > size {
			bs = p[i : i+size]
		}

		command := fmt.Sprintf("+QISEND=%d,%d", connectID, len(bs))
		_, err := this.send2(command)
		if err != nil {
			return err
		}

		//防止数据串,发送数据的时候无标识的,固加锁
		this.mu.Lock()
		_, err = this.write(p)
		this.mu.Unlock()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *AT) newTCP(connectID int8) *TCP {
	return &TCP{
		at:        this,
		connectID: connectID,
		Closer: safe.NewCloserErr(errors.New("未初始化")).SetCloseFunc(func(err error) error {
			return this.TCPClose(connectID)
		}),
		ch: chans.NewIO(1),
	}
}

// TCP 直吐模式
type TCP struct {
	at        *AT
	connectID int8
	*safe.Closer
	ch *chans.IO
}

func (this *TCP) Dial(address string) error {
	this.at.MQTTClose(this.connectID)
	if err := this.at.TCPDial(this.connectID, address); err != nil {
		return err
	}
	//重置错误
	this.Closer.Reset()
	return nil
}

func (this *TCP) ReadMessage() ([]byte, error) {
	return this.ch.ReadMessage()
}

func (this *TCP) Read(p []byte) (int, error) {
	return this.ch.Read(p)
}

func (this *TCP) Write(p []byte) (int, error) {
	err := this.at.TCPWrite(this.connectID, p)
	return len(p), err
}
