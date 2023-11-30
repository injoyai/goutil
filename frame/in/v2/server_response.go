package in

import (
	"fmt"
	"github.com/injoyai/conv"
	"io"
	"os"
)

var DefaultClient = New(WithDefault())

//=================================Return=================================//

func Return(code int, data interface{}) { DefaultClient.Exit(code, data) }

func Return200(data interface{}) { Return(200, data) }

func Return400(data interface{}) { Return(400, data) }

func Return401(data interface{}) { Return(401, data) }

func Return403(data interface{}) { Return(403, data) }

func Return415(data interface{}) { Return(415, data) }

func Return500(data interface{}) { Return(500, data) }

// Redirect 重定向
func Redirect(addr string) { Return(302, nil) }

//=================================File=================================//

// ReturnFilePath 返回本地文件
func ReturnFilePath(name, path string) {
	bs, err := os.ReadFile(path)
	CheckErr(err)
	DefaultClient.File(name, bs)
}

// ReturnFileLocal 返回本地文件
func ReturnFileLocal(name, path string) { ReturnFilePath(name, path) }

func ReturnFileReader(name string, reader io.Reader) {
	bs, err := io.ReadAll(reader)
	CheckErr(err)
	DefaultClient.File(name, bs)
}

// ReturnFileBytes 返回文件
func ReturnFileBytes(name string, bs []byte) { DefaultClient.File(name, bs) }

//=================================Json=================================//

// Json 返回json
func Json(httpCode int, data interface{}, count ...int64) {
	DefaultClient.Json(httpCode, data, count...)
}

func Json200(data interface{}, count ...int64) { Json(200, data, count...) }

func Json400(data interface{}) { Json(400, data) }

func Json401() { Json(401, "验证失败") }

func Json403() { Json(403, "没有权限") }

func Json415(data interface{}) { Json(415, data) }

func Json500(data interface{}) { Json(500, data) }

//=================================Injoy=================================//

func Succ(data interface{}, count ...int64) { DefaultClient.Succ(data, count...) }

func Fail(data interface{}) { DefaultClient.Fail(data) }

// Err 退出,并校验错误
func Err(data interface{}, succData ...interface{}) {
	if data == nil {
		Succ(conv.DefaultInterface(nil, succData...))
	} else {
		Fail(data)
	}
}

// Errf 退出格式化错误信息
func Errf(format string, args ...interface{}) {
	Err(fmt.Sprintf(format, args...))
}

// CheckErr 检测错误(遇到错误结束)
func CheckErr(err error, failMsg ...string) {
	if err != nil {
		Err(conv.DefaultString(err.Error(), failMsg...))
	}
}
