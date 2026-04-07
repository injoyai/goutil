package http

import (
	"net/http"
)

type (
	ResponseWriter = http.ResponseWriter
	Writer         = http.ResponseWriter
	Header         = http.Header
	Cookie         = http.Cookie
)

const (
	MethodGet     = http.MethodGet
	MethodHead    = http.MethodHead
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodConnect = http.MethodConnect
	MethodOptions = http.MethodOptions
	MethodTrace   = http.MethodTrace
)
