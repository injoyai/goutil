package ip

import (
	"testing"
	"time"
)

func TestNewPinger(t *testing.T) {
	p := NewPinger()
	for {
		s, err := p.Ping()
		if err != nil {
			t.Log(err)
		}
		t.Log(s)
		<-time.After(time.Second)
	}

}
