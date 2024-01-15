package voice

type Interface interface {
	Call(msg *Message) error
}

type Message struct {
	TemplateID string `json:"templateID"` //模板id
	Phone      string `json:"phone"`      //手机号
	Param      string `json:"param"`      //参数,阿里云是json,腾讯云是列表,隔开
}
