package http

import (
	"github.com/injoyai/io"
	"net/http"
)

type (
	Plan           = io.Plan
	ResponseWriter = http.ResponseWriter
	Writer         = http.ResponseWriter
	Header         = http.Header
	Cookie         = http.Cookie
)
