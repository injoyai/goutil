package sms

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"strings"
)

type aliyun struct {
	option *Option
	*dysmsapi.Client
	err error
}

func (this *Option) aliyun() *aliyun {
	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(this.SecretID, this.SecretKey)
	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", config, credential)
	return &aliyun{
		option: this,
		Client: client,
		err:    err,
	}
}

func (this *aliyun) Send(msg *Message) error {
	if this.err != nil {
		return this.err
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = this.option.SignName
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
