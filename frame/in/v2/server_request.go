package in

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/net/ghttp"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/other/upload"
	"io"
	"net/http"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
)

// Read 解析json数据
func Read(r interface{}, v interface{}) {
	if err := json.Unmarshal(GetBodyBytes(r), v); err != nil {
		Json415(err)
	}
}

// ReadJson 解析json数据
func ReadJson(r interface{}, ptr interface{}) {
	Read(r, ptr)
}

// GetBody 获取body的流
func GetBody(r interface{}) io.ReadCloser {
	switch v := r.(type) {
	case *http.Request:
		return v.Body
	case *ghttp.Request:
		return io.NopCloser(bytes.NewReader(v.GetBody()))
	case *gin.Context:
		read, err := v.Request.GetBody()
		CheckErr(err)
		return read
	default:
		Err(ErrInvalidRequest)
	}
	return nil
}

// GetBodyBytes 获取body字节
func GetBodyBytes(r interface{}) []byte {
	switch v := r.(type) {
	case *http.Request:
		bs, err := io.ReadAll(v.Body)
		CheckErr(err)
		v.Body.Close()
		v.Body = io.NopCloser(bytes.NewReader(bs))
		return bs
	case *ghttp.Request:
		return v.GetBody()
	case *gin.Context:
		bs, err := io.ReadAll(v.Request.Body)
		CheckErr(err)
		v.Request.Body.Close()
		v.Request.Body = io.NopCloser(bytes.NewReader(bs))
		return bs
	default:
		Err(ErrInvalidRequest)
	}
	return []byte{}
}

// GetBodyString 获取body字符
func GetBodyString(r interface{}) string {
	return string(GetBodyBytes(r))
}

// GetRequest 获取请求对象
func GetRequest(r interface{}) *http.Request {
	switch v := r.(type) {
	case *http.Request:
		return v
	case *ghttp.Request:
		return v.Request
	case *gin.Context:
		return v.Request
	default:
		Err(ErrInvalidRequest)
	}
	return nil
}

// GetHeader 请求头
func GetHeader(r interface{}) http.Header {
	return GetRequest(r).Header
}

// Header 请求头
func Header(r interface{}) http.Header {
	return GetRequest(r).Header
}

// GetFile 获取上传的文件流
func GetFile(r interface{}, name string) io.ReadCloser {
	switch v := r.(type) {
	case *http.Request:
		f, _, err := v.FormFile(name)
		CheckErr(err)
		return f
	case *ghttp.Request:
		f, err := v.GetUploadFile(name).Open()
		CheckErr(err)
		return f
	case *gin.Context:
		fh, err := v.FormFile(name)
		CheckErr(err)
		f, err := fh.Open()
		CheckErr(err)
		return f
	default:
		Err(ErrInvalidRequest)
	}
	return nil
}

// GetFileBytes 获取上传的文件字节
func GetFileBytes(r interface{}, name string) []byte {
	f := GetFile(r, name)
	defer f.Close()
	bs, err := io.ReadAll(f)
	CheckErr(err)
	return bs
}

// SaveFile 保存上传的文件
func SaveFile(r interface{}, name, filename string) string {
	return UploadFile(r, name, filename, upload.DefaultLocal)
}

// UploadFile 获取上传的文件,上传到...
// name: 上传文件的字段名
// filename: 保存的文件名
// upload: 保存对象接口
func UploadFile(r interface{}, name, filename string, upload upload.Interface) string {
	f := GetFile(r, name)
	defer f.Close()
	s, err := upload.Save(filename, f)
	CheckErr(err)
	return s
}

func GetBodyMap(r interface{}) *conv.Map {
	return conv.NewMap(GetBodyBytes(r))
}

// Get 获取参数
func Get(r interface{}, key string, def ...interface{}) *conv.Var {
	switch v := r.(type) {
	case *http.Request:
		if val, ok := v.URL.Query()[key]; ok {
			return conv.New(val[0])
		}
		if x := GetBodyMap(v).GetVar(key); !x.IsNil() {
			return x
		}
		bs := GetBodyBytes(v)
		CheckErr(v.ParseForm())
		if val, ok := v.Form[key]; ok {
			v.Body = io.NopCloser(bytes.NewReader(bs))
			return conv.New(val)
		}
		if val, ok := v.Header[key]; ok {
			return conv.New(val)
		}
	case *ghttp.Request:
		return conv.New(v.Get(key, def...))
	case *gin.Context:
		if val, ok := v.GetQuery(key); ok {
			return conv.New(val)
		}
		if x := GetBodyMap(v).GetVar(key); !x.IsNil() {
			return x
		}
		bs := GetBodyBytes(v)
		if val, ok := v.GetPostForm(key); ok {
			v.Request.Body = io.NopCloser(bytes.NewReader(bs))
			return conv.New(val)
		}
		if val, ok := v.Request.Header[key]; ok {
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

func GetString(r interface{}, key string, def ...interface{}) string {
	return Get(r, key, def...).String()
}

func GetBool(r interface{}, key string, def ...interface{}) bool {
	return Get(r, key, def...).Bool()
}

func GetInt(r interface{}, key string, def ...interface{}) int {
	return Get(r, key, def...).Int()
}

func GetInt64(r interface{}, key string, def ...interface{}) int64 {
	return Get(r, key, def...).Int64()
}

func GetFloat(r interface{}, key string, def ...interface{}) float64 {
	return Get(r, key, def...).Float64()
}

func GetPageNum(r interface{}) int {
	return DefaultClient.GetPageNum(r)
}

func GetPageSize(r interface{}) int {
	return DefaultClient.GetPageSize(r)
}
