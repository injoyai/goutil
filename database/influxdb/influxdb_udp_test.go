package influx

import (
	"testing"
)

func TestNewUDPClient(t *testing.T) {
	c := NewUDPClient(&UDPOption{Addr: "localhost:8089"})
	if err := c.Err(); err != nil {
		t.Error(err)
		return
	}
	c.Ping()
	{
		result, err := c.Exec("select * from test")
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(result)
	}
}
