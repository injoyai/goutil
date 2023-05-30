package sms

type Interface interface {
	Send(message *Message) error
}

type Message struct {
	Phone      []string `json:"phone"`      //手机号
	Param      string   `json:"param"`      //参数,阿里云是json,腾讯云是列表,隔开
	TemplateID string   `json:"templateID"` //模板id
}
