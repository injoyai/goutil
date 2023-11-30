package http

import "net/http"

const (
	HeaderKeySpend         = "Injoy-Spend"
	HeaderKeyTry           = "Injoy-Try"
	HeaderKeyAuthorization = "Authorization"
	HeaderKeyUserAgent     = "User-Agent"
	HeaderKeyAccept        = "Accept"
	HeaderKeyContentType   = "Content-Type"
	HeaderKeyConnection    = "Connection"
	HeaderKeyReferer       = "Referer"
)

var (
	HeaderBase = http.Header{
		HeaderKeyUserAgent:   {UserAgentDefault},
		HeaderKeyContentType: {"application/json;charset=utf-8"}, //发送的数据格式
		HeaderKeyAccept:      {"application/json"},               //希望接收的数据格式
		HeaderKeyConnection:  {"close"},                          //短连接
	}

	HeaderCORS = http.Header{
		"AllowOrigin":      []string{"*"},
		"AllowMethods":     []string{"GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"},
		"AllowCredentials": []string{"true"},
		"AllowHeaders":     []string{"Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With"},
		"MaxAge":           []string{"3628800"},
	}

	CORS = HeaderCORS
)

type Header = http.Header
