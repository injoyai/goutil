package proxy

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/elazarl/goproxy"
	"github.com/injoyai/logs"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type (
	ReqHandler = func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response)

	RespHandler = func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response
)

type ReqAction struct {
	*goproxy.ReqProxyConds
	log *logs.Entity
}

func (this *ReqAction) Do(fs ...ReqHandler) *ReqAction {
	for _, f := range fs {
		this.ReqProxyConds.DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if req == nil || f == nil {
				return req, nil
			}
			return f(req, ctx)
		})
	}
	return this
}

func (this *ReqAction) SetHeader(k, v string) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Set(k, v)
		return req, nil
	})
}

func (this *ReqAction) SetHeaders(header http.Header) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header = header
		return req, nil
	})
}

func (this *ReqAction) DelCookie() *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Del("Cookie")
		return req, nil
	})
}

func (this *ReqAction) AddCookie(cookie *http.Cookie) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Add("Cookie", cookie.String())
		return req, nil
	})
}

func (this *ReqAction) SetBody(body []byte) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Body = io.NopCloser(bytes.NewReader(body))
		return req, nil
	})
}

func (this *ReqAction) DoNothing() *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, nil
	})
}

func (this *ReqAction) Response(resp *http.Response) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, resp
	})
}

func (this *ReqAction) ResponseHtml(body string) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, NewHtmlResponse(req, body)
	})
}

func (this *ReqAction) ResponseHtmlFile(filename string) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		bs, _ := os.ReadFile(filename)
		return req, NewHtmlResponse(req, string(bs))
	})
}

func (this *ReqAction) ResponsePng(body []byte) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, NewPngResponse(req, body)
	})
}

func (this *ReqAction) ResponseJpg(body []byte) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, NewJpgResponse(req, body)
	})
}

func (this *ReqAction) ResponseGif(body []byte) *ReqAction {
	return this.Do(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, NewGifResponse(req, body)
	})
}

/*






 */

type RespAction struct {
	*goproxy.ProxyConds
	log *logs.Entity
}

func (this *RespAction) Do(fs ...RespHandler) *RespAction {
	for _, f := range fs {
		this.ProxyConds.DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if resp == nil || f == nil {
				return resp
			}
			return f(resp, ctx)
		})
	}
	return this
}

func (this *RespAction) DoNothing() *RespAction {
	return this.Do(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		return resp
	})
}

func (this *RespAction) OnURL(f func(u *url.URL)) *RespAction {
	return this.Do(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if f == nil || resp.Request == nil || resp.Request.URL == nil {
			return resp
		}
		f(resp.Request.URL)
		return resp
	})
}

func (this *RespAction) OnQuery(f func(q url.Values)) *RespAction {
	return this.OnURL(func(u *url.URL) {
		if f != nil {
			f(u.Query())
		}
	})
}

func (this *RespAction) PrintHost() *RespAction {
	return this.Do(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		fmt.Println(resp.Request.Host)
		return resp
	})
}

func (this *RespAction) PrintRequest(body ...bool) *RespAction {
	return this.Do(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		bs, err := httputil.DumpRequest(resp.Request, len(body) > 0 && body[0])
		if err != nil {
			this.log.Printf("[错误] %v\n", err)
			return resp
		}
		fmt.Println("=============================================================")
		fmt.Println(string(bs))
		return resp
	})
}

func (this *RespAction) PrintResponse(body ...bool) *RespAction {
	return this.Do(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		bs, err := httputil.DumpResponse(resp, len(body) > 0 && body[0])
		if err != nil {
			this.log.Printf("[错误] %v\n", err)
			return resp
		}
		fmt.Println(string(bs))
		return resp
	})
}

func (this *RespAction) Print(body ...bool) *RespAction {
	return this.Do(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		reqBs, err := httputil.DumpRequest(resp.Request, len(body) > 0 && body[0])
		if err != nil {
			this.log.Printf("[错误] %v\n", err)
			return resp
		}

		respBs, err := httputil.DumpResponse(resp, len(body) > 0 && body[0])
		if err != nil {
			this.log.Printf("[错误] %v\n", err)
			return resp
		}
		fmt.Println("=============================================================")
		fmt.Println(string(reqBs))
		fmt.Println(string(respBs))
		return resp
	})
}

func (this *RespAction) ReplaceBody(old, new string) *RespAction {
	return this.Do(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp == nil || resp.Body == nil {
			return resp
		}

		// 读取原始响应体body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp
		}
		resp.Body.Close()

		// 替换内容
		modified := strings.ReplaceAll(string(bodyBytes), old, new)

		// 设置新的响应体body
		resp.Body = io.NopCloser(strings.NewReader(modified))
		resp.ContentLength = int64(len(modified))
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(modified)))

		return resp
	})
}

/*
Document 解析body
使用方法参考
https://blog.csdn.net/qq_38334677/article/details/129225231

	//查找标签: 		doc.Find("body,div,...") 多个用,隔开
	//查找ID: 		doc.Find("#id1")
	//查找class: 	doc.Find(".class1")
	//查找属性: 		doc.Find("div[lang]") doc.Find("div[lang=zh]") doc.Find("div[id][lang=zh]")
	//查找子节点: 	doc.Find("body>div")
	//过滤数据: 		doc.Find("div:contains(xxx)")
	//过滤节点: 		dom.Find("span:has(div)")
	doc.Find("body").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})

选择器					说明
Find(“div[lang]”)		筛选含有lang属性的div元素
Find(“div[lang=zh]”)	筛选lang属性为zh的div元素
Find(“div[lang!=zh]”)	筛选lang属性不等于zh的div元素
Find(“div[lang¦=zh]”)	筛选lang属性为zh或者zh-开头的div元素
Find(“div[lang*=zh]”)	筛选lang属性包含zh这个字符串的div元素
Find(“div[lang~=zh]”)	筛选lang属性包含zh这个单词的div元素，单词以空格分开的
Find(“div[lang$=zh]”)	筛选lang属性以zh结尾的div元素，区分大小写
Find(“div[lang^=zh]”)	筛选lang属性以zh开头的div元素，区分大小写
*/
func (this *RespAction) Document(f func(resp *http.Response, ctx *goproxy.ProxyCtx, doc *goquery.Document)) *RespAction {
	return this.Do(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if f == nil {
			return resp
		}
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			this.log.Printf("[错误] %v\n", err)
			return resp
		}
		resp.Body = io.NopCloser(bytes.NewReader(bs))
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bs))
		if err != nil {
			this.log.Printf("[错误] %v\n", err)
			return resp
		}
		f(resp, ctx, doc)
		return resp
	})
}
