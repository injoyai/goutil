package in

import (
	"encoding/json"
	"github.com/injoyai/conv"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/net/ghttp"
)

func read(i any, v any) {
	bs := getBody(i)
	if err := json.Unmarshal(bs, v); err != nil {
		Json415(err)
	}
}

func getBody(r any) []byte {
	switch v := r.(type) {
	case *http.Request:
		bs, err := ioutil.ReadAll(v.Body)
		CheckErr(err)
		return bs
	case *ghttp.Request:
		return v.GetBody()
	case *gin.Context:
		read, err := v.Request.GetBody()
		CheckErr(err)
		bs, err := ioutil.ReadAll(read)
		CheckErr(err)
		return bs
	default:
		Err(ErrInvalidRequest)
	}
	return []byte{}
}

func getHeader(r any) http.Header {
	switch v := r.(type) {
	case *http.Request:
		return v.Header
	case *ghttp.Request:
		return v.Header
	case *gin.Context:
		return v.Request.Header
	default:
		Err(ErrInvalidRequest)
	}
	return http.Header{}
}

func getFile(r any, name string) (bs []byte) {
	switch v := r.(type) {
	case *http.Request:
		f, _, err := v.FormFile(name)
		CheckErr(err)
		bs, err = ioutil.ReadAll(f)
		CheckErr(err)
	case *ghttp.Request:
		f, err := v.GetUploadFile(name).Open()
		CheckErr(err)
		defer f.Close()
		bs, err = ioutil.ReadAll(f)
		CheckErr(err)
	case *gin.Context:
		fh, err := v.FormFile(name)
		CheckErr(err)
		f, err := fh.Open()
		CheckErr(err)
		bs, err = ioutil.ReadAll(f)
		CheckErr(err)
	default:
		Err(ErrInvalidRequest)
	}
	return
}

func get(r any, key string, def ...any) *conv.Var {
	switch v := r.(type) {
	case *http.Request:
		m := v.URL.Query()
		if val, ok := m[key]; ok {
			return conv.New(val[0])
		}
	case *ghttp.Request:
		return conv.New(v.Get(key, def...))
	case *gin.Context:
		if val, ok := v.GetQuery(key); ok {
			return conv.New(val)
		}
	default:
		Err(ErrInvalidRequest)
	}
	if len(def) > 0 {
		return conv.New(def[0])
	}
	return conv.New(nil)
}

func getPageSize(r any) (int, int) {
	return get(r, DefaultOption.Page, 1).Int() - 1,
		get(r, DefaultOption.Size, DefaultOption.SizeDef).Int()
}
