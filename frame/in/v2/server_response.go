package in

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"os"
)

//=================================Return=================================//

func Return(code int, data any) { DefaultClient.Text(code, data) }

func Return200(data any) { DefaultClient.Text(http.StatusOK, data) }

func Text(code int, data any) { DefaultClient.Text(code, data) }

func Text200(data any) { DefaultClient.Text(http.StatusOK, data) }

func Yaml(code int, data any) { DefaultClient.Yaml(code, data) }

func Yaml200(data any) { DefaultClient.Yaml(http.StatusOK, data) }

func Toml(code int, data any) { DefaultClient.Toml(code, data) }

func Toml200(data any) { DefaultClient.Toml(http.StatusOK, data) }

func Xml(code int, data any) { DefaultClient.Xml(code, data) }

func Xml200(data any) { DefaultClient.Xml(http.StatusOK, data) }

func Proto(code int, data proto.Message) { DefaultClient.Proto(code, data) }

func Proto200(data proto.Message) { DefaultClient.Proto(http.StatusOK, data) }

func Msgpack(code int, data proto.Message) { DefaultClient.Msgpack(code, data) }

func Msgpack200(data proto.Message) { DefaultClient.Msgpack(http.StatusOK, data) }

func Html(code int, data any) { DefaultClient.Html(code, data) }

func Html200(data any) { DefaultClient.Html(http.StatusOK, data) }

// Redirect301 永久重定向,GET和HEAD自动请求,其他需要用户确认
func Redirect301(addr string) { DefaultClient.Redirect(http.StatusMovedPermanently, addr) }

// Redirect302 临时重定向,GET和HEAD自动请求,其他需要用户确认
func Redirect302(addr string) { DefaultClient.Redirect(http.StatusFound, addr) }

// Redirect307 Temporary Redirect临时重定向,不改变请求方式
func Redirect307(addr string) { DefaultClient.Redirect(http.StatusTemporaryRedirect, addr) }

// Redirect308 Permanent Redirect永久重定向,不改变请求方式
func Redirect308(addr string) { DefaultClient.Redirect(http.StatusPermanentRedirect, addr) }

//=================================File=================================//

// FilePath 返回本地文件
func FilePath(name, path string) {
	bs, err := os.ReadFile(path)
	CheckErr(err)
	DefaultClient.File(name, bs)
}

// FileLocal 返回本地文件
func FileLocal(name, path string) {
	FilePath(name, path)
}

func FileReader(name string, reader io.Reader) {
	//适用于小文件
	bs, err := io.ReadAll(reader)
	CheckErr(err)
	DefaultClient.File(name, bs)
}

// FileBytes 返回文件
func FileBytes(name string, bs []byte) {
	DefaultClient.File(name, bs)
}

//=================================Other=================================//

func CopyReader(w http.ResponseWriter, filename string, r io.Reader) {
	DefaultClient.CopyFile(w, filename, r)
}

func Copy(w http.ResponseWriter, r io.Reader) {
	DefaultClient.Copy(w, r)
}

func Proxy(w http.ResponseWriter, r *http.Request, uri string) {
	DefaultClient.Proxy(w, r, uri)
}

//=================================Json=================================//

// Json 返回json
func Json(httpCode int, data any) { DefaultClient.Json(httpCode, data) }

func Json200(data any) { Json(http.StatusOK, data) }

func Json400(data any) { Json(http.StatusBadRequest, data) }

func Json401() { Json(http.StatusUnauthorized, "验证失败") }

func Json403() { Json(http.StatusForbidden, "没有权限") }

func Json404() { Json(http.StatusNotFound, "接口不存在") }

func Json415(data any) { Json(http.StatusUnsupportedMediaType, data) }

func Json500(data any) { Json(http.StatusInternalServerError, data) }

//=================================Injoy=================================//

func Succ(data any, count ...int64) { DefaultClient.Succ(data, count...) }

func Fail(data any) { DefaultClient.Fail(data) }

// Err 退出,并校验错误
func Err(data any, succData ...any) {
	if data == nil {
		Succ(conv.Default[any](nil, succData...))
	} else {
		Fail(data)
	}
}

// Errf 退出格式化错误信息
func Errf(format string, args ...any) {
	Err(fmt.Sprintf(format, args...))
}

// CheckErr 检测错误(遇到错误结束)
func CheckErr(err error, failMsg ...string) {
	if err != nil {
		Err(conv.Default[string](err.Error(), failMsg...))
	}
}
