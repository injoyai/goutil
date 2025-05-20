package in

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/golang/protobuf/proto"
	"github.com/injoyai/conv"
	"github.com/pelletier/go-toml/v2"
	"github.com/ugorji/go/codec"
	"gopkg.in/yaml.v3"
)

type Marshal func(v any) ([]byte, error)

type IMarshal interface {
	Bytes() ([]byte, error)
	ContentType() []string
}

type TEXT struct {
	Data any
}

func (this *TEXT) Bytes() ([]byte, error) {
	return []byte(conv.String(this.Data)), nil
}

func (this *TEXT) ContentType() []string {
	return []string{"text/plain; charset=utf-8"}
}

type JSON struct {
	Data any
}

func (this *JSON) Bytes() ([]byte, error) {
	return json.Marshal(this.Data)
}

func (this *JSON) ContentType() []string {
	return []string{"application/json; charset=utf-8"}
}

type YAML struct {
	Data any
}

func (this *YAML) Bytes() ([]byte, error) {
	return yaml.Marshal(this.Data)
}

func (this *YAML) ContentType() []string {
	return []string{"application/yaml; charset=utf-8"}
}

type TOML struct {
	Data any
}

func (this *TOML) Bytes() ([]byte, error) {
	return toml.Marshal(this.Data)
}

func (this *TOML) ContentType() []string {
	return []string{"application/toml; charset=utf-8"}
}

type XML struct {
	Data any
}

func (this *XML) Bytes() ([]byte, error) {
	return xml.Marshal(this.Data)
}

func (this *XML) ContentType() []string {
	return []string{"application/xml; charset=utf-8"}
}

type PROTO struct {
	Data proto.Message
}

func (this *PROTO) Bytes() ([]byte, error) {
	return proto.Marshal(this.Data)
}

func (this *PROTO) ContentType() []string {
	return []string{"application/x-protobuf"}
}

type MSGPACK struct {
	Data any
}

func (this *MSGPACK) Bytes() ([]byte, error) {
	var mh codec.MsgpackHandle
	w := bytes.NewBuffer(nil)
	err := codec.NewEncoder(w, &mh).Encode(this.Data)
	return w.Bytes(), err
}

func (this *MSGPACK) ContentType() []string {
	return []string{"application/msgpack; charset=utf-8"}
}

type HTML struct {
	Data any
}

func (this *HTML) Bytes() ([]byte, error) {
	return []byte(conv.String(this.Data)), nil
}

func (this *HTML) ContentType() []string {
	return []string{"text/html; charset=utf-8"}
}

type FILE struct {
	Data any
}

func (this *FILE) Bytes() ([]byte, error) {
	return []byte(conv.String(this.Data)), nil
}

func (this *FILE) ContentType() []string {
	return []string{"application/octet-stream"}
}
