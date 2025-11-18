package http

import (
	"testing"
)

func TestUrl(t *testing.T) {
	{
		_, err := Url("www.baidu.com").Debug().Do()
		if err != nil {
			t.Error(err)
			return
		}
	}
	{
		_, err := Url("http://www.baidu.com").Debug().Do()
		if err != nil {
			t.Error(err)
			return
		}
	}
	{
		_, err := Url("http://www.baidu.com").Do()
		if err != nil {
			t.Error(err)
			return
		}
	}

}
