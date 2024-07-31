package in

import (
	"net/http"
)

var DefaultClient = New(WithDefault())

func Recover(h http.Handler) http.Handler {
	return DefaultClient.Recover(h)
}

func SetSuccFailCode(succ, fail interface{}) *Client {
	return DefaultClient.SetSuccFailCode(succ, fail)
}

func SetSuccFail(f func(c *Client, succ bool, data interface{}, count ...int64)) *Client {
	return DefaultClient.SetSuccFail(f)
}

func MiddleRecover(err interface{}, w http.ResponseWriter) {
	DefaultClient.MiddleRecover(err, w)
}
