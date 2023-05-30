package sms

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"strings"
)

func NewAliyun(cfg *AliyunConfig) (Interface, error) {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(cfg.SecretID, cfg.SecretKey)
	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", config, credential)
	return &aliyun{
		cfg:    cfg,
		Client: client,
	}, err
}

type AliyunConfig struct {
	SecretID  string //
	SecretKey string //
	SignName  string //签名
}

type aliyun struct {
	cfg *AliyunConfig
	*dysmsapi.Client
}

func (this *aliyun) Send(msg *Message) error {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = this.cfg.SignName
	request.TemplateCode = msg.TemplateID
	request.PhoneNumbers = strings.Join(msg.Phone, ",")
	request.TemplateParam = msg.Param
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
