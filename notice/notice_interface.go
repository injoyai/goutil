package notice

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/io"
)

type Interface interface {

	// Publish 发布通知消息
	Publish(message *Message) error
}

type Message struct {
	Target  string `json:"target"`  //目标
	Param   g.Map  `json:"param"`   //其它参数
	Content string `json:"content"` //内容
	Tag     g.Map  `json:"tag"`     //标签,可以记录操作人等信息
}

func NewAudio(cfg *AudioConfig) (Interface, error) {
	return &audio{cfg}, nil
}

func NewPhoneAliyun(cfg *PhoneAliyunConfig) (Interface, error) {
	regionId := conv.SelectString(len(cfg.RegionID) == 0, "cn-hangzhou", cfg.RegionID)
	client, err := dyvmsapi.NewClientWithAccessKey(regionId, cfg.SecretID, cfg.SecretKey)
	return &phoneAliyun{
		cfg:    cfg,
		Client: client,
	}, err
}

func NewPhoneTencent(cfg *PhoneTencentConfig) (Interface, error) {
	return nil, nil
}

func NewIO(dial io.DialFunc, options ...io.OptionClient) (Interface, error) {
	return &IO{Client: io.Redial(dial, options...)}, nil
}
