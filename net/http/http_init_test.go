package http

import (
	"testing"
)

func TestGetBytes(t *testing.T) {
	//DefaultClient //.Debug() //.SetTimeout(time.Millisecond)
	bs, err := GetBytes("http://www.baidu.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bs))
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
