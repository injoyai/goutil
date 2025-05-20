package js

import (
	"errors"
	"fmt"
	"github.com/injoyai/goutil/script"
	"testing"
)

func TestNew(t *testing.T) {
	n := New()
	n.SetFunc("panic", func(args *script.Args) (any, error) {
		panic("panic")
	})
	n.SetFunc("err", func(args *script.Args) (any, error) {
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
	x := New(script.WithObject, script.WithFunc)
	x.SetFunc("test", func(args *script.Args) (any, error) {
		return nil, errors.New("错误")
	})
	result, err := x.Exec(`
//test()
//panic(668)
x=logs.Debug(666)
logs.Debug(x)
logs.Debug(x[0],x[1])
http.Url("http://192.168.192.2:8181").Get().GetBodyString()
c=net.Dial("tcp",":10086")
logs.Debug(c)
sleep(0.1)
logs.Debug(c==nil)
c.Close()
global.Map.Set("key",123)
value=global.Map.GetString("key")
x=logs.Debug(value)
logs.Debug(conv.Bytes("666").HEX())
logs.Debug(x)
logs.Debug(x)
logs.Debug(net.Dial("icmp","192.168.192.2").Ping())
//net.Dial("icmp","192.168.192.2").For(3,1)
logs.Debug(bytes.Sum(conv.Bytes("666")))
logs.Debug(bytes.Sum([1,2,3]))
logs.Debug(bytes.Reverse([1,2,3]).HEX())

s=mux.New()
s.ALL("/",function(r){
logs.Info(r.GetBodyString())
mux.Succ(666)
})
go(function(){
sleep(5)
s.Close()
})
logs.Debug("run")
s.Run()
os.Shell("start www.baidu.com")
`)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}
