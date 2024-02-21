package notice

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/injoyai/conv"
)

func NewAliyunSMS(cfg *AliyunConfig) (Interface, error) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(cfg.SecretID, cfg.SecretKey)
	regionId := conv.SelectString(len(cfg.RegionID) == 0, "cn-hangzhou", cfg.RegionID)
	client, err := dysmsapi.NewClientWithOptions(regionId, config, credential)
	return &aliyunSMS{
		cfg:    cfg,
		Client: client,
	}, err
}

type aliyunSMS struct {
	cfg *AliyunConfig
	*dysmsapi.Client
}

// Publish 发送自定义参数是json格式{"a":"1"}
func (this *aliyunSMS) Publish(msg *Message) error {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = this.cfg.SignName
	request.TemplateCode = conv.String(msg.Param["TemplateID"])
	request.PhoneNumbers = msg.Target
	request.TemplateParam = conv.String(msg.Param["Param"])
	response, err := this.Client.SendSms(request)
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

/*



 */

func NewAliyunPhone(cfg *AliyunConfig) (Interface, error) {
	regionId := conv.SelectString(len(cfg.RegionID) == 0, "cn-hangzhou", cfg.RegionID)
	client, err := dyvmsapi.NewClientWithAccessKey(regionId, cfg.SecretID, cfg.SecretKey)
	return &aliyunPhone{
		cfg:    cfg,
		Client: client,
	}, err
}

type aliyunPhone struct {
	cfg *AliyunConfig
	*dyvmsapi.Client
}

// Publish 发送自定义参数是json格式{"a":"1"}
func (this *aliyunPhone) Publish(msg *Message) error {
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

/*



 */

type AliyunConfig struct {
	SecretID  string `json:"secretID"`  //
	SecretKey string `json:"secretKey"` //
	SignName  string `json:"signName"`  //签名
	RegionID  string `json:"regionId"`  //地域,如cn-hangzhou
}
