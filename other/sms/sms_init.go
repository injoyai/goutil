package sms

type Client interface {
	Send(message *Message) error
}

type Message struct {
	Phone      []string `json:"phone"`      //手机号
	Param      string   `json:"param"`      //参数,阿里云是json,腾讯云是列表,隔开
	TemplateID string   `json:"templateID"` //模板id
}

type Option struct {
	Model     int    //0是腾讯云,1是阿里云
	SecretID  string //
	SecretKey string //

	SignName string //签名
	AppID    string //腾讯云用
}

func New(op *Option) Client {
	switch op.Model {
	case 1: //阿里云
		return op.aliyun()
	default: //腾讯云
		return op.tencent()
	}
}
