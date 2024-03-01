package notice

const (
	TargetPop   = "pop"
	TargetPopup = "popup"
)

type Interface interface {

	// Publish 发布通知消息
	Publish(message *Message) error
}

type Message struct {
	Target  string                 `json:"target"`  //目标
	Title   string                 `json:"title"`   //标题
	Content string                 `json:"content"` //内容
	Param   map[string]interface{} `json:"param"`   //其它参数
	Tag     map[string]interface{} `json:"tag"`     //标签,可以记录操作人等信息
}
