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

func TestHTTP(t *testing.T) {
	x := New(script.WithBaseFunc)
	result, err := x.Exec(`
x=logs.Debug(666)
logs.Debug(x)
logs.Debug(x[0],x[1])
http.Url("http://192.168.192.2:8181").Get().GetBodyString()
c=net.Dial("tcp",":10086")
logs.Debug(c)
sleep(5)
c.Close()
global.Map.Set("key",123)
value=global.Map.GetString("key")
logs.Debug(value)


`)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}
