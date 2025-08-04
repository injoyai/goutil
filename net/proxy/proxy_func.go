package proxy

import (
	"github.com/elazarl/goproxy"
	"net/http"
	"regexp"
	"strings"
)

type Condition = goproxy.RespCondition

func DomainIs(domain ...string) Condition {
	return goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		for _, d := range domain {
			if strings.HasSuffix(req.Host, "."+d) {
				return true
			}
		}
		return false
	})
}

func HostLike(reg ...string) Condition {
	return goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		for _, r := range reg {
			if regexp.MustCompile(r).MatchString(req.Host) {
				return true
			}
		}
		return false
	})
}

func HostIs(host ...string) Condition {
	return goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		for _, h := range host {
			if req.Host == h {
				return true
			}
		}
		return false
	})
}

func PathIs(path ...string) Condition {
	return goproxy.ReqConditionFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		if req == nil {
			return false
		}
		for _, p := range path {
			if req.URL.Path == p {
				return true
			}
		}
		return false
	})
}

/*

1.2 常见的媒体格式
text/html：HTML格式
text/plain：纯文本格式
text/XML：XML格式
image/gif：gif图片格式
image/jpeg：jpg图片格式
image/png：png图片格式
以application开头的媒体格式类型：

application/xhtml+xml：XHTML格式
application/xml：XML数据格式
application/atom+xml：Atom XML聚合格式
application/json：JSON数据格式【常用】
application/pdf：pdf数据格式
application/msword：Word文档格式
application/octet-stream：二进制流格式（如常见的文件下载）
application/x-www-form-urlencoded：<form encType=" ">中默认的encType，form表单数据编码为key/value格式发送到服务器（表单默认的提交数据格式）【常用】另外一种常见的媒体格式是上传文件时使用的：【常用】
multipart/form-data： ['mʌlti:pɑ:t] 需要在表单进行文件上传时，就需要使用该格式



*/

func NewResponse(r *http.Request, status int, body string, contentType string) *http.Response {
	return goproxy.NewResponse(r, contentType, status, body)
}

func NewTextResponse(r *http.Request, text string) *http.Response {
	return NewResponse(r, http.StatusAccepted, text, "text/plain;charset=UTF-8")
}

func NewHtmlResponse(r *http.Request, text string) *http.Response {
	return NewResponse(r, http.StatusAccepted, text, "text/html;charset=UTF-8")
}

func NewJsonResponse(r *http.Request, text string) *http.Response {
	return NewResponse(r, http.StatusAccepted, text, "application/json;charset=UTF-8")
}

func NewPngResponse(r *http.Request, text []byte) *http.Response {
	return NewResponse(r, http.StatusAccepted, string(text), "image/png")
}

func NewJpgResponse(r *http.Request, text []byte) *http.Response {
	return NewResponse(r, http.StatusAccepted, string(text), "image/jpeg")
}

func NewGifResponse(r *http.Request, text []byte) *http.Response {
	return NewResponse(r, http.StatusAccepted, string(text), "image/gif")
}
