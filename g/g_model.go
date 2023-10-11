package g

type Resp struct {
	Code int         `json:"code"` //状态
	Data interface{} `json:"data"` //数据
	Msg  string      `json:"msg"`  //消息
}
