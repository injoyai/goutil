package file_log

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	dir := "./output/log/"
	x := New(&Config{
		Dir:      dir,
		Layout:   "2006-01-02-15-04.log",
		SaveTime: time.Minute * 10,
	})
	oss.ListenExit(func() {
		os.RemoveAll(dir)
	})

	go func() {
		for {
			<-time.After(time.Second * 20)
			c, err := x.GetLogCurve(
				time.Now().Add(-time.Minute*2),
				time.Now(),
				time.Second*20,
				&DecodeFunc{
					DecodeFunc: func(bs []byte) (any, error) {
						return string(bs), nil
					},
					ReportFunc: func(node int64, list []any) (any, error) {
						return list, nil
					},
				},
			)
			if err != nil {
				t.Error(err)
				continue
			}
			t.Logf("\n%s\n", c)
		}
	}()

	for i := 0; ; i++ {
		<-time.After(time.Second)
		if _, err := x.WriteString(conv.String(i)); err != nil {
			t.Error(err)
		}
	}

}
