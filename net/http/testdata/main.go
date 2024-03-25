package main

import (
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/logs"
)

func main() {

	HTTP := http.DefaultClient

	err := HTTP.Url("http://127.0.0.1:8080").
		SetQuery("1", "2").
		SetUrl("http://www.baidu.com:80/a").
		//SetUrl("http://127.0.0.1").
		Debug().Get().Err()

	logs.PrintErr(err)

}
