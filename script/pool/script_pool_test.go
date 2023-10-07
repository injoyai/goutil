package script_pool

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	p := New()
	for i := 0; i < 10000; i++ {
		go func(i int) {
			if _, err := p.Exec(fmt.Sprintf("print(%d)", i)); err != nil {
				t.Log(err)
			}
		}(i)
	}

	<-time.After(time.Second * 10)
	t.Log(p.count)
}
