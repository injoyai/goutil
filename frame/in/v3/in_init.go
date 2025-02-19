package in

import (
	"bytes"
	"fmt"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"os"
	"reflect"
	"time"
)

var DefaultClient = New(WithDefault())

// SetCacheByHandler 尝试从缓存中获取数据,如果不存在则通过函数获取,执行函数时,其他相同的key会等待此次结果
func SetCacheByHandler(key interface{}, handler func() interface{}, expiration time.Duration) interface{} {
	value, err := DefaultClient.GetOrSetByHandler(key, func() (interface{}, error) { return handler(), nil }, expiration)
	CheckErr(err)
	return value
}

// DelCache 删除缓存数据
func DelCache(key ...interface{}) {
	for _, v := range key {
		DefaultClient.Del(v)
	}
}

// SetCache 设置缓存,覆盖缓存
func SetCache(key interface{}, value interface{}, expiration time.Duration) {
	DefaultClient.Set(key, value, expiration)
}

// HTTPHandler 对http.HandlerFunc使用中间件,http.ListenAndServe(":80",in.HTTPHandler(h))
func HTTPHandler(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return DefaultClient.Recover(http.HandlerFunc(h))
}

// Recover 对http.Handler使用中间件
func Recover(h http.Handler) http.Handler {
	return DefaultClient.Recover(h)
}

// RecoverFunc 对http.Handler使用中间件
func RecoverFunc(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return DefaultClient.Recover(http.HandlerFunc(f))
}

// MiddleRecover 捕捉panic,或自定义panic,并输出到http.ResponseWriter
func MiddleRecover(e interface{}, w http.ResponseWriter) {
	DefaultClient.MiddleRecover(e, w)
}

// SetHandlerWithCode 设置常用响应的状态码
func SetHandlerWithCode(succ, fail, unauthorized, forbidden interface{}) *Client {
	return DefaultClient.SetHandlerWithCode(succ, fail, unauthorized, forbidden)
}

func GetPageNum(r *http.Request) int {
	return DefaultClient.GetPageNum(r)
}

func GetPageSize(r *http.Request) int {
	return DefaultClient.GetPageSize(r)
}

//=================================Return=================================//

func Return(code int, data interface{}) { DefaultClient.Text(code, data) }

func Return200(data interface{}) { Return(http.StatusOK, data) }

func Text(code int, data interface{}) { DefaultClient.Text(code, data) }

func Text200(data interface{}) { Return(http.StatusOK, data) }

func Html(code int, data interface{}) { DefaultClient.Html(code, data) }

func Html200(data interface{}) { DefaultClient.Html(http.StatusOK, data) }

func Reader(code int, r io.ReadCloser) { DefaultClient.Reader(code, r) }

func Reader200(r io.ReadCloser) { DefaultClient.Reader(http.StatusOK, r) }

// Redirect301 永久重定向,GET和HEAD自动请求,其他需要用户确认
func Redirect301(addr string) { DefaultClient.Redirect(http.StatusMovedPermanently, addr) }

// Redirect302 临时重定向,GET和HEAD自动请求,其他需要用户确认
func Redirect302(addr string) { DefaultClient.Redirect(http.StatusFound, addr) }

// Redirect307 Temporary Redirect临时重定向,不改变请求方式
func Redirect307(addr string) { DefaultClient.Redirect(http.StatusTemporaryRedirect, addr) }

// Redirect308 Permanent Redirect永久重定向,不改变请求方式
func Redirect308(addr string) { DefaultClient.Redirect(http.StatusPermanentRedirect, addr) }

//=================================File=================================//

// FileLocal 返回本地文件
func FileLocal(name, filename string) {
	f, err := os.Open(filename)
	CheckErr(err)
	i, err := f.Stat()
	CheckErr(err)
	DefaultClient.File(name, i.Size(), f)
}

// FileReader 返回文件流
func FileReader(name string, r io.ReadCloser) {
	DefaultClient.File(name, -1, r)
}

// FileBytes 返回文件,字节
func FileBytes(name string, bs []byte) {
	DefaultClient.File(name, int64(len(bs)), io.NopCloser(bytes.NewReader(bs)))
}

//=================================Other=================================//

func Proxy(w http.ResponseWriter, r *http.Request, uri string) {
	DefaultClient.Proxy(w, r, uri)
}

//=================================Json=================================//

// Json 返回json
func Json(httpCode int, data interface{}) { DefaultClient.Json(httpCode, data) }

func Json200(data interface{}) { Json(http.StatusOK, data) }

func Json400(data interface{}) { Json(http.StatusBadRequest, data) }

func Json401() { Json(http.StatusUnauthorized, "验证失败") }

func Json403() { Json(http.StatusForbidden, "没有权限") }

func Json404() { Json(http.StatusNotFound, "接口不存在") }

func Json415(data interface{}) { Json(http.StatusUnsupportedMediaType, data) }

func Json500(data interface{}) { Json(http.StatusInternalServerError, data) }

//=================================Injoy=================================//

func Succ(data interface{}, count ...int64) { DefaultClient.Succ(data, count...) }

func Fail(data interface{}) { DefaultClient.Fail(data) }

func Forbidden() { DefaultClient.Forbidden() }

func Unauthorized() { DefaultClient.Unauthorized() }

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

// Range 取值范围
func Range(value interface{}, min, max int, message string) {
	Min(value, min, message)
	Max(value, max, message)
}

// Min 介于当前go mod版本
func Min(value interface{}, min int, message string) {
	fail := true
	switch val := value.(type) {
	case nil:

	case int:
		fail = val < min

	case int8:
		fail = int(val) < min

	case int16:
		fail = int(val) < min

	case int32:
		fail = int(val) < min

	case int64:
		fail = int(val) < min

	case uint:
		fail = int(val) < min

	case uint8:
		fail = int(val) < min

	case uint16:
		fail = int(val) < min

	case uint32:
		fail = int(val) < min

	case uint64:
		fail = int(val) < min

	case float32:
		fail = val < float32(min)

	case float64:
		fail = val < float64(min)

	case string:
		fail = len(val) < min

	case []byte:
		fail = len(val) < min

	default:

		r := reflect.ValueOf(value)
		switch r.Kind() {
		case reflect.Pointer:
			Min(r.Elem().Interface(), min, message)
			return

		case reflect.Slice, reflect.Array:
			fail = r.Len() < min

		case reflect.Map:
			fail = r.Len() < min

		case reflect.Chan:
			fail = r.Len() < min

		default:

		}

	}

	if fail {
		Err(message)
	}
}

// Max 介于当前go mod版本
func Max(value interface{}, max int, message string) {
	fail := true
	switch val := value.(type) {
	case nil:

	case int:
		fail = val > max

	case int8:
		fail = int(val) > max

	case int16:
		fail = int(val) > max

	case int32:
		fail = int(val) > max

	case int64:
		fail = int(val) > max

	case uint:
		fail = int(val) > max

	case uint8:
		fail = int(val) > max

	case uint16:
		fail = int(val) > max

	case uint32:
		fail = int(val) > max

	case uint64:
		fail = int(val) > max

	case float32:
		fail = val > float32(max)

	case float64:
		fail = val > float64(max)

	case string:
		fail = len(val) > max

	case []byte:
		fail = len(val) > max

	default:

		r := reflect.ValueOf(value)
		switch r.Kind() {
		case reflect.Pointer:
			Max(r.Elem().Interface(), max, message)
			return

		case reflect.Slice, reflect.Array:
			fail = r.Len() > max

		case reflect.Map:
			fail = r.Len() > max

		case reflect.Chan:
			fail = r.Len() > max

		default:

		}

	}

	if fail {
		Err(message)
	}
}
