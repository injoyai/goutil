package voice

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
)

func NewAliyun(cfg *AliyunConfig) (Interface, error) {
	client, err := dyvmsapi.NewClientWithAccessKey("cn-hangzhou", cfg.SecretID, cfg.SecretKey)
	return &aliyun{
		cfg:    cfg,
		Client: client,
	}, err
}

type AliyunConfig struct {
	SecretID  string `json:"secretID"`  //
	SecretKey string `json:"secretKey"` //
	SignName  string `json:"signName"`  //签名
}

type aliyun struct {
	cfg *AliyunConfig
	*dyvmsapi.Client
}

func (this *aliyun) Call(msg *Message) error {
	if len(msg.Phone) == 0 {
		return nil
	}
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.AcceptFormat = "json"
	request.CalledNumber = msg.Phone
	request.TtsCode = msg.TemplateID
	request.TtsParam = msg.Param
	response := dyvmsapi.CreateAddRtcAccountResponse()
	err := this.DoAction(request, response)
	if err != nil {
		return err
	}
	if !response.IsSuccess() {
		return errors.New(response.Message)
	}
	if response.Code != "OK" {
		return errors.New(response.Message)
	}
	return nil
}
