package file_log

import (
	"github.com/injoyai/base/g"
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
		for range g.Interval(time.Second * 20) {
			c, err := x.GetLogCurve(
				time.Now().Add(-time.Minute*2),
				time.Now(),
				time.Second*20,
				&DecodeFunc{
					DecodeFunc: func(bs []byte) (interface{}, error) {
						return string(bs), nil
					},
					ReportFunc: func(node int64, list []interface{}) (interface{}, error) {
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

	for i := range g.Interval(time.Second) {
		if _, err := x.WriteString(conv.String(i)); err != nil {
			t.Error(err)
		}
	}

}
