package upload

import (
	"bytes"
	"testing"
)

func TestDialFTP(t *testing.T) {
	c, err := DialFTP("192.168.192.2:21", "download", "download")
	if err != nil {
		t.Error(err)
		return
	}
	data := bytes.NewBufferString("Hello World")
	u, err := c.Upload("/test/test.txt", data)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(u)

	t.Log(u.Download("./est.txt"))

	ls, err := c.List("/test/")
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range ls {
		t.Log(v)
	}

}
