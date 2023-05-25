package influx

import (
	"crypto/tls"
	client "github.com/influxdata/influxdb1-client/v2"
	"net/http"
	"net/url"
	"time"
)

type HTTPOption struct {

	// Database 数据库
	Database string

	// Precision 精度s
	Precision string

	// Addr should be of the form "http://host:port"
	// or "http://[ipv6-host%zone]:port".
	Addr string

	// Username is the influxdb username, optional.
	Username string

	// Password is the influxdb password, optional.
	Password string

	// UserAgent is the http User Agent, defaults to "InfluxDBClient".
	UserAgent string

	// Timeout for influxdb writes, defaults to no timeout.
	Timeout time.Duration

	// InsecureSkipVerify gets passed to the http client, if true, it will
	// skip https certificate verification. Defaults to false.
	InsecureSkipVerify bool

	// TLSConfig allows the user to set their own TLS config for the HTTP
	// Client. If set, this option overrides InsecureSkipVerify.
	TLSConfig *tls.Config

	// Proxy configures the Proxy function on the HTTP client.
	Proxy func(req *http.Request) (*url.URL, error)
}

func (this *HTTPOption) new() {
	if len(this.Database) == 0 {
		this.Database = "_default"
	}
	if len(this.Precision) == 0 {
		this.Precision = "ns"
	}
	if len(this.Addr) == 0 {
		this.Addr = "http://127.0.0.1:8086"
	}
	if len(this.Username) == 0 {
		this.Username = "admin"
	}
}

// NewHTTPClient 新建客户端
// show databases //查库
// show measurements //查表
func NewHTTPClient(op *HTTPOption) *Client {
	op.new()
	c := &Client{option: &option{
		Database:  op.Database,
		Precision: op.Precision,
	}}
	c.client, c.err = client.NewHTTPClient(client.HTTPConfig{
		Addr:               op.Addr,
		Username:           op.Username,
		Password:           op.Password,
		UserAgent:          op.UserAgent,
		Timeout:            op.Timeout,
		InsecureSkipVerify: op.InsecureSkipVerify,
		TLSConfig:          op.TLSConfig,
		Proxy:              op.Proxy,
	})
	if c.err != nil {
		return c
	}
	_, c.err = c.Exec("CREATE DATABASE " + op.Database)
	c.newClient = func() *Client { return NewHTTPClient(op) }
	return c
}
