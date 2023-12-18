package in

import (
	"encoding/json"
)

type ExitModel struct {
	Mark   string            //退出标识,需要独一无二(如果不止用in包)
	Type   string            //退出类型,暂时无用
	Code   int               //响应状态码
	Header map[string]string //响应请求头
	Value  []byte            //响应内容,body
}

// NewExit 新建退出实例
func NewExit(code int, bytes []byte) *ExitModel {
	return &ExitModel{
		Mark:   DefaultOption.ExitMark,
		Type:   DefaultOption.ExitOk,
		Code:   code,
		Header: map[string]string{},
		Value:  bytes,
	}
}

// NewExitJson 新建json退出实例
func NewExitJson(code int, data interface{}) *ExitModel {
	val, err := json.Marshal(data)
	if err != nil {
		return NewExit(500, []byte(err.Error())).SetHeaderJson()
	}
	return NewExit(code, val).SetHeaderJson()
}

// SetHeader 设置响应头部
func (this *ExitModel) SetHeader(i, v string) *ExitModel {
	this.Header[i] = v
	return this
}

// SetHeaderJson 设置响应头部json格式
func (this *ExitModel) SetHeaderJson() *ExitModel {
	return this.SetHeader("Content-Type", "application/json;charset=utf-8")
}

// ClearHeader 清除所有自定义响应头部
func (this *ExitModel) ClearHeader() *ExitModel {
	this.Header = map[string]string{}
	return this
}

// Exit 退出panic
func (this *ExitModel) Exit() {
	bs, _ := json.Marshal(this)
	panic(string(bs))
}
