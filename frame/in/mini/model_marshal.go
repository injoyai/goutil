package in

import (
	"encoding/json"
	"github.com/injoyai/conv"
)

type Marshal func(v interface{}) ([]byte, error)

type IMarshal interface {
	Bytes() ([]byte, error)
	ContentType() []string
}

type TEXT struct {
	Data interface{}
}

func (this *TEXT) Bytes() ([]byte, error) {
	return []byte(conv.String(this.Data)), nil
}

func (this *TEXT) ContentType() []string {
	return []string{"text/plain; charset=utf-8"}
}

type JSON struct {
	Data interface{}
}

func (this *JSON) Bytes() ([]byte, error) {
	return json.Marshal(this.Data)
}

func (this *JSON) ContentType() []string {
	return []string{"application/json; charset=utf-8"}
}

type HTML struct {
	Data interface{}
}

func (this *HTML) Bytes() ([]byte, error) {
	return []byte(conv.String(this.Data)), nil
}

func (this *HTML) ContentType() []string {
	return []string{"text/html; charset=utf-8"}
}

type FILE struct {
	Data interface{}
}

func (this *FILE) Bytes() ([]byte, error) {
	return []byte(conv.String(this.Data)), nil
}

func (this *FILE) ContentType() []string {
	return []string{"application/octet-stream"}
}
