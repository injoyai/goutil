package influx

import client "github.com/influxdata/influxdb1-client/v2"

type UDPOption struct {

	// Database 数据库
	Database string

	// Precision 精度s
	Precision string

	// Addr should be of the form "host:port"
	// or "[ipv6-host%zone]:port".
	Addr string

	// PayloadSize is the maximum size of a UDP client message, optional
	// Tune this based on your network. Defaults to UDPPayloadSize.
	PayloadSize int
}

func (this *UDPOption) new() {
	if len(this.Database) == 0 {
		this.Database = "_default"
	}
	if len(this.Precision) == 0 {
		this.Precision = "ns"
	}
	if len(this.Addr) == 0 {
		this.Addr = "localhost:8086"
	}
}

func NewUDPClient(op *UDPOption) *Client {
	op.new()
	c := &Client{option: &option{
		Database:  op.Database,
		Precision: op.Precision,
	}}
	c.client, c.err = client.NewUDPClient(client.UDPConfig{
		Addr:        op.Addr,
		PayloadSize: op.PayloadSize,
	})
	if c.err != nil {
		return c
	}
	_, c.err = c.Exec("CREATE DATABASE " + op.Database)
	return c
}
