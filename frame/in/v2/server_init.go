package in

import (
	"net/http"
)

var DefaultClient = New(WithDefault())

func Recover(h http.Handler) http.Handler {
	return DefaultClient.Recover(h)
}

func SetSuccFailCode(succ, fail any) *Client {
	return DefaultClient.SetSuccFailCode(succ, fail)
}

func SetSuccFail(f func(c *Client, succ bool, data any, count ...int64)) *Client {
	return DefaultClient.SetSuccFail(f)
}

func MiddleRecover(err any, w http.ResponseWriter) {
	DefaultClient.MiddleRecover(err, w)
}
