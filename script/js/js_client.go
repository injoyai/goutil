package js

import (
	"fmt"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/script"
	"github.com/robertkrimen/otto"
	"sync"
)

var (
	Nil = otto.NullValue()
	_   = script.Client(new(Client))
)

func New(option ...func(c script.Client)) *Client {
	vm := otto.New()
	cli := &Client{
		Otto: vm,
	}
	cli.Set("nil", otto.NullValue())
	cli.SetFunc("print", func(args *script.Args) (interface{}, error) {
		fmt.Println(args.Interfaces()...)
		return nil, nil
	})
	cli.SetFunc("println", func(args *script.Args) (interface{}, error) {
		fmt.Println(args.Interfaces()...)
		return nil, nil
	})
	cli.SetFunc("printf", func(args *script.Args) (interface{}, error) {
		a := args.Interfaces()
		if len(a) > 0 {
			fmt.Printf(conv.String(a[0]), a[1:]...)
		} else {
			fmt.Printf("")
		}
		return nil, nil
	})
	for _, v := range option {
		v(cli)
	}
	return cli
}

type Client struct {
	*otto.Otto
	mu  sync.Mutex
	tag *maps.Safe
}

func (this *Client) Tag() *maps.Safe {
	if this.tag == nil {
		this.tag = maps.NewSafe()
	}
	return this.tag
}

func (this *Client) Exec(text string, option ...func(client script.Client)) (result interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	for _, v := range option {
		v(this)
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	value, err := this.Otto.Run(text)
	if err != nil {
		return nil, err
	}
	return value.Export()
}

func (this *Client) GetVar(key string) *conv.Var {
	val, _ := this.Otto.Get(key)
	value, _ := val.Export()
	return conv.New(value)
}

func (this *Client) Set(key string, value interface{}) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	switch fn := value.(type) {
	case script.Func:
		value = this.toFunc(fn)
	case func(*script.Args) (interface{}, error):
		value = this.toFunc(fn)
	case func(*script.Args) interface{}:
		value = this.toFunc(func(args *script.Args) (interface{}, error) {
			return fn(args), nil
		})
	case func(*script.Args) error:
		value = this.toFunc(func(args *script.Args) (interface{}, error) {
			return nil, fn(args)
		})
	case func(*script.Args):
		value = this.toFunc(func(args *script.Args) (interface{}, error) {
			fn(args)
			return nil, nil
		})
	case func():
		value = this.toFunc(func(args *script.Args) (interface{}, error) {
			fn()
			return nil, nil
		})
	}
	return this.Otto.Set(key, value)
}

func (this *Client) SetFunc(key string, fn script.Func) error {
	return this.Set(key, this.toFunc(fn))
}

func (this *Client) Close() error {
	return nil
}

func (this *Client) toFunc(fn script.Func) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) (val otto.Value) {
		defer func() {
			if err := recover(); err != nil {
				panic(call.Otto.MakeCustomError("", fmt.Sprint(err)))
			}
		}()
		args := []*conv.Var(nil)
		for _, v := range call.ArgumentList {
			if v.IsFunction() {
				args = append(args, conv.New(func(i ...interface{}) (otto.Value, error) {
					return v.Call(Nil)
				}))
				continue
			}
			val, err := v.Export()
			if err != nil {
				panic(err)
			}
			args = append(args, conv.New(val))
		}
		arg := &script.Args{
			This: this,
			Args: args,
		}

		value, err := fn(arg)
		if err != nil {
			panic(err)
		}
		result, err := this.ToValue(value)
		if err != nil {
			panic(err)
		}
		return result
	}
}
