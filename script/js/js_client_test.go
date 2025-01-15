package js

import (
	"errors"
	"fmt"
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

func TestObj(t *testing.T) {
	x := New()
	err := x.Set("obj", &Obj{Name: "test", Age: 18})
	if err != nil {
		t.Error(err)
		return
	}

	x.Exec(`print(obj.Name,obj.Age)`)
	_, err = x.Exec(`obj.Print()`)
	if err != nil {
		t.Error(err)
		return
	}
}

type Obj struct {
	Name string
	Age  int
}

func (this *Obj) Print() {
	fmt.Println(this.Name, this.Age)
}
