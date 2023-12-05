package in

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
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

func (this *Resp) Default() (int, g.Map) {
	return http.StatusOK, this.Map()
}

func (this *Resp) Map() g.Map {
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
	return m
}

func NewRespMap(code interface{}, data interface{}, count ...int64) g.Map {
	m := map[string]interface{}{
		"code": code,
		"data": data,
	}
	if len(count) > 0 {
		m["count"] = count[0]
	}
	switch val := data.(type) {
	case error:
		errMsg := conv.String(val)
		m["msg"] = errMsg
		m["data"] = errMsg
	}
	return m
}
