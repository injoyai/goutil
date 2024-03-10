package err

import (
	"github.com/injoyai/logs"
	"testing"
)

func TestNew(t *testing.T) {
	logs.SetWriter(logs.Stdout)
	er := newErr2()
	if er != nil {
		logs.Err(er)
	}
}

func newErr2() *Item {
	return newErr()
}

func newErr() *Item {
	return New("测试错误")
}
