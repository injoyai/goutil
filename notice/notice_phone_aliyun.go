package notice

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/injoyai/conv"
)

type PhoneAliyunConfig struct {
	SecretID  string `json:"secretID"`  //
	SecretKey string `json:"secretKey"` //
	SignName  string `json:"signName"`  //签名
	RegionID  string `json:"regionId"`  //地域,如cn-hangzhou
}

type phoneAliyun struct {
	cfg *PhoneAliyunConfig
	*dyvmsapi.Client
}

func (this *phoneAliyun) Publish(msg *Message) error {
	if len(msg.Target) == 0 {
		return nil
	}
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.AcceptFormat = "json"
	request.CalledNumber = msg.Target
	request.TtsCode = conv.String(msg.Param["TemplateID"])
	request.TtsParam = conv.String(msg.Param["Param"])
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
