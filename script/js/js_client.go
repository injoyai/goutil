package js

import (
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/script"
	"github.com/robertkrimen/otto"
	"sync"
)

var (
	Nil = otto.NullValue()
	_   = script.Interface(new(Client))
)

func New() *Client {
	vm := otto.New()
	cli := &Client{
		Otto: vm,
	}
	cli.Set("print", cli.toFunc(func(args *script.Args) interface{} {
		fmt.Println(args.Interfaces()...)
		return Nil
	}))
	cli.Exec("var console={\nlog:function(any){\nprint(any)\n}\n}")
	return cli
}

type Client struct {
	*otto.Otto
	mu sync.Mutex
}

func (this *Client) Exec(text string) (*conv.Var, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	value, err := this.Otto.Run(text)
	if err != nil {
		return conv.Nil(), err
	}
	val, _ := value.Export()
	return conv.New(val), nil
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
	case func(*script.Args) interface{}:
		value = this.toFunc(fn)
	case func():
		value = this.toFunc(func(args *script.Args) interface{} {
			fn()
			return nil
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
	return func(call otto.FunctionCall) otto.Value {
		it, _ := call.This.Export()
		args := []*conv.Var(nil)
		for _, v := range call.ArgumentList {
			val, _ := v.Export()
			args = append(args, conv.New(val))
		}
		arg := &script.Args{
			This:      conv.New(it),
			Args:      args,
			Interface: this,
		}
		result, _ := otto.ToValue(fn(arg))
		return result
	}
}
