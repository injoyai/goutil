package http

import "net/http"

const (
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
		"Access-Control-Allow-Origin":      []string{"*"},
		"Access-Control-Allow-Methods":     []string{"GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"},
		"Access-Control-Allow-Credentials": []string{"true"},
		"Access-Control-Allow-Headers":     []string{"Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With"},
		"Access-Control-Allow-Max-Age":     []string{"3628800"},
	}

	CORS = HeaderCORS
)
