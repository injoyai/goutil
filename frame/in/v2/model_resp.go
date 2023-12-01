package in

import (
	"github.com/injoyai/conv"
	"net/http"
)

func newResp(httpCode int, data interface{}, count ...int64) *Resp {
	return &Resp{
		Code:  httpCode,
		Data:  data,
		Count: count,
	}
}

type Resp struct {
	Code  interface{}
	Data  interface{}
	Count []int64
}

func (this *Resp) Default() (int, []byte) {
	return http.StatusOK, this.Bytes()
}

func (this *Resp) Bytes() []byte {
	m := map[string]interface{}{
		"code": this.Code,
		"data": this.Data,
	}
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
