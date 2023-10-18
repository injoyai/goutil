package http

import (
	"io/ioutil"
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
	request := NewRequest("", url, nil)
	if len(bind) > 0 {
		request.Bind(bind)
	}
	return request.Get()
}

func GetBytes(url string) ([]byte, error) {
	return GetBody(url)
}

func GetBody(uri string) ([]byte, error) {
	resp, err := DefaultClient.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func Post(url string, bind ...interface{}) *Response {
	request := NewRequest("", url, nil)
	if len(bind) > 0 {
		request.Bind(bind)
	}
	return request.Post()
}

func Put(url string, bind ...interface{}) *Response {
	request := NewRequest("", url, nil)
	if len(bind) > 0 {
		request.Bind(bind)
	}
	return request.Put()
}

func Delete(url string, bind ...interface{}) *Response {
	request := NewRequest("", url, nil)
	if len(bind) > 0 {
		request.Bind(bind)
	}
	return request.Delete()
}

func Url(url string) *Request {
	return NewRequest("", url, nil)
}
