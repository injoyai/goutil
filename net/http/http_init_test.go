package http

import (
	"net/http"
	"testing"
)

func TestGetBytes(t *testing.T) {
	resp := NewRequest(http.MethodGet, "https://gitee.com/injoyai/file/releases/download/v0.0.1/in.exe", nil).SetHeader("Connection", "keep-alive").Debug().Get()
	if resp.Err() != nil {
		t.Error(resp.Err())
		return
	}
	t.Log(len(resp.GetBodyBytes()))
	t.Log(resp.GetBodyString())
}

func TestUrl(t *testing.T) {
	var result string
	resp := Url("http://www.baidu.com").SetQuery("1", "2").
		AddHeader("1", "2").Retry(3).Debug().Bind(&result).Get()
	if resp.Err() != nil {
		t.Error(resp.Err())
	}
	t.Log(result)
}

func TestGetToWriter(t *testing.T) {
	if err := DefaultClient.SetProxy("http://127.0.0.1:1081"); err != nil {
		t.Error(err)
	}
	err := GetToFile("https://github.com/injoyai/downloader/releases/latest/download/downloader.exe", "./test.exe")
	if err != nil {
		t.Error(err)
	}
}
