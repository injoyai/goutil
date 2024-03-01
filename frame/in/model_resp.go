package in

import (
	"github.com/injoyai/conv"
)

func NewRespMap(code interface{}, data interface{}, count ...int64) map[string]interface{} {
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
