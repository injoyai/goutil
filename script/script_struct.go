package script

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/cfg/v2"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/client/dial"
	"github.com/injoyai/logs"
	"net"
)

func NewGlobal() *Global {
	return &Global{
		Map: maps.NewSafe(),
	}
}

type Global struct {
	Map *maps.Safe
}

type Net struct{}

func (this *Net) Dial(network, address string) net.Conn {
	panic(789)
	c, err := net.Dial(network, address)
	if err != nil {
		panic(err)
	}
	return c
}

func (this *Net) DialTCP(address string) net.Conn {
	return this.Dial("tcp", address)
}

func (this *Net) DialUDP(address string) net.Conn {
	return this.Dial("udp", address)
}

type Ios struct{}

func (this *Ios) Dial(args *Args) (*client.Client, error) {
	switch args.GetString(1) {
	case "tcp":
		return dial.TCP(args.GetString(2))
	case "mqtt":

	}
	return nil, nil
}

type Logs struct{}

func (this *Logs) Debug(args ...interface{}) (int, error) {
	return logs.Debug(args...)
}

func (this *Logs) Debugf(format string, args ...interface{}) (int, error) {
	return logs.Debugf(format, args...)
}

func (this *Logs) Info(args ...interface{}) (int, error) {
	return logs.Info(args...)
}

func (this *Logs) Infof(format string, args ...interface{}) (int, error) {
	return logs.Infof(format, args...)
}

func (this *Logs) Err(args ...interface{}) (int, error) {
	return logs.Err(args...)
}

func (this *Logs) Errf(format string, args ...interface{}) (int, error) {
	return logs.Errf(format, args...)
}

func (this *Logs) Error(args ...interface{}) (int, error) {
	return logs.Error(args...)
}

func (this *Logs) Errorf(format string, args ...interface{}) (int, error) {
	return logs.Errorf(format, args...)
}

func NewHTTP() *HTTP {
	return &HTTP{
		DefaultClient: http.DefaultClient,
	}
}

type HTTP struct {
	DefaultClient *http.Client
}

func (this *HTTP) Url(url string) *http.Request {
	return http.Url(url)
}

func NewConv() *Conv {
	return &Conv{
		Json: codec.Json,
		Toml: codec.Toml,
		Yaml: codec.Yaml,
		Ini:  codec.Ini,
	}
}

type Conv struct {
	Json codec.Interface
	Toml codec.Interface
	Yaml codec.Interface
	Ini  codec.Interface
}

func (this *Conv) New(i interface{}) *conv.Var {
	return conv.New(i)
}

func (this *Conv) NewMap(i interface{}, codec ...codec.Interface) *conv.Map {
	return conv.NewMap(i, codec...)
}

func NewCfg() *Cfg {
	return &Cfg{
		Extend: cfg.Default,
	}
}

type Cfg struct {
	conv.Extend
}

func (this *Cfg) Init(i ...conv.IGetVar) {
	cfg.Init(i...)
}

func (this *Cfg) New(i ...conv.IGetVar) *cfg.Entity {
	return cfg.New(i...)
}

func (this *Cfg) WithAny(i interface{}) conv.IGetVar {
	return cfg.WithAny(i)
}

func (this *Cfg) WithFile(path string, codec ...codec.Interface) conv.IGetVar {
	return cfg.WithFile(path, codec...)
}

func (this *Cfg) WithEnv() conv.IGetVar {
	return cfg.WithEnv()
}
