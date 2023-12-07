package in

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
)

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
