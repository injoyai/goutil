package http

import (
	"net/http"
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

	HeaderSpend         = "Injoy-Spend"
	HeaderTry           = "Injoy-Try"
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
)

func Get(url string, bind ...interface{}) *Response {
	return DefaultClient.Get(url, bind...)
}

func GetBytes(url string) ([]byte, error) {
	return DefaultClient.GetBytes(url)
}

func Post(url string, body interface{}, bind ...interface{}) *Response {
	return DefaultClient.Post(url, body, bind...)
}

func Put(url string, body interface{}, bind ...interface{}) *Response {
	return DefaultClient.Put(url, body, bind...)
}

func Delete(url string, body interface{}, bind ...interface{}) *Response {
	return DefaultClient.Delete(url, body, bind...)
}

func Url(url string) *Request {
	return NewRequest("", url, nil)
}
