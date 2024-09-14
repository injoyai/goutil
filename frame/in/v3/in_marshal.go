package in

import (
	"bytes"
	"encoding/json"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"strconv"
)

type Marshal func(v interface{}) ([]byte, error)

type IMarshal interface {
	io.ReadCloser
	Header() http.Header
}

type TEXT struct {
	Data   interface{}
	reader *bytes.Reader
}

func (this *TEXT) Read(p []byte) (int, error) {
	if this.reader == nil {
		this.reader = bytes.NewReader([]byte(conv.String(this.Data)))
	}
	return this.reader.Read(p)
}

func (this *TEXT) Header() http.Header {
	return http.Header{
		"Content-Type": []string{"text/plain; charset=utf-8"},
	}
}

func (this *TEXT) Close() error { return nil }

type JSON struct {
	Data   interface{}
	reader *bytes.Reader
}

func (this *JSON) Read(p []byte) (int, error) {
	if this.reader == nil {
		bs, err := json.Marshal(this.Data)
		if err != nil {
			return 0, err
		}
		this.reader = bytes.NewReader(bs)
	}
	return this.reader.Read(p)
}

func (this *JSON) Header() http.Header {
	return http.Header{
		"Content-Type": []string{"application/json; charset=utf-8"},
	}
}

func (this *JSON) Close() error { return nil }

type HTML struct {
	Data   interface{}
	reader *bytes.Reader
}

func (this *HTML) Read(p []byte) (int, error) {
	if this.reader == nil {
		this.reader = bytes.NewReader([]byte(conv.String(this.Data)))
	}
	return this.reader.Read(p)
}

func (this *HTML) Header() http.Header {
	return http.Header{
		"Content-Type": []string{"text/html; charset=utf-8"},
	}
}

func (this *HTML) Close() error { return nil }

type FILE struct {
	Name string
	Size int64
	io.ReadCloser
}

func (this *FILE) Read(p []byte) (int, error) {
	if this.ReadCloser == nil {
		return 0, io.EOF
	}
	return this.ReadCloser.Read(p)
}

func (this *FILE) Header() http.Header {
	return http.Header{
		"Content-Type":        []string{"application/octet-stream"},
		"Content-Disposition": []string{"attachment; filename=" + this.Name},
		"Content-Length":      []string{strconv.FormatInt(this.Size, 10)},
	}
}