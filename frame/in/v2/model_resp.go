package in

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/logs"
)

func newResp(code interface{}, data interface{}, msg string, count ...int64) *Resp {
	return &Resp{
		Code:  code,
		Data:  data,
		Msg:   msg,
		Count: count,
	}
}

type Resp struct {
	Code  interface{} `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Count []int64     `json:"count"`
}

func (this *Resp) Bytes() []byte {
	m := map[string]interface{}{
		"code": this.Code,
		"data": this.Data,
		"msg":  this.Msg,
	}
	logs.Debug(m)
	if len(this.Count) > 0 {
		m["count"] = this.Count
	}
	switch val := this.Data.(type) {
	case error:
		errMsg := conv.String(val)
		m["msg"] = errMsg
		m["data"] = errMsg
	}
	return conv.Bytes(m)
}
