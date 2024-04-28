package js

import (
	"errors"
	"github.com/injoyai/goutil/script"
	"testing"
)

func TestNew(t *testing.T) {
	n := New()
	n.SetFunc("panic", func(args *script.Args) (interface{}, error) {
		panic("panic")
	})
	n.SetFunc("err", func(args *script.Args) (interface{}, error) {
		return nil, errors.New("错误")
	})
	t.Log(n.Exec(`
 throw new Error('The value is not a number.');
`))
	t.Log(n.Exec(`panic()`))
	t.Log(n.Exec(`err()`))
	t.Log(n.Exec(`100`))
	t.Log(n.Exec(`"test"`))

}
