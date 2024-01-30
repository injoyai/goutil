package ip

import (
	"testing"
)

func TestParseV4(t *testing.T) {
	t.Log(ParseV4("127.0.0.1"))
	t.Log(ParseV4("192.168.10.15:10086"))
}
