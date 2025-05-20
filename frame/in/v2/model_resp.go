package in

import (
	"github.com/injoyai/conv"
)

func NewRespMap(code any, data any, count ...int64) map[string]any {
	m := map[string]any{
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
