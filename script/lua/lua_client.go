package lua

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/script"
	"github.com/yuin/gopher-lua"
	"sync"
)

/*

================ ========================= ================== =======================
 Type name        Go type                   Type() value       Constants
================ ========================= ================== =======================
 ``LNilType``      (constants)              ``LTNil``          ``LNil``
 ``LBool``         (constants)              ``LTBool``         ``LTrue``, ``LFalse``
 ``LNumber``        float64                 ``LTNumber``       ``-``
 ``LString``        string                  ``LTString``       ``-``
 ``LFunction``      struct pointer          ``LTFunction``     ``-``
 ``LUserData``      struct pointer          ``LTUserData``     ``-``
 ``LState``         struct pointer          ``LTThread``       ``-``
 ``LTable``         struct pointer          ``LTTable``        ``-``
 ``LChannel``       chan LValue             ``LTChannel``      ``-``
================ ========================= ================== =======================


数据类型	描述
nil	这个最简单，只有值nil属于该类，表示一个无效值（在条件表达式中相当于false）。
boolean	包含两个值：false和true。
number	表示双精度类型的实浮点数
string	字符串由一对双引号或单引号来表示
function	由 C 或 Lua 编写的函数
userdata	表示任意存储在变量中的C数据结构
thread	表示执行的独立线路，用于执行协同程序
table	Lua 中的表（table）其实是一个"关联数组"（associative arrays），数组的索引可以是数字、字符串或表类型。在 Lua 里，table 的创建是通过"构造表达式"来完成，最简单构造表达式是{}，用来创建一个空表。


*/

var (
	Nil = lua.LNil
	_   = script.Interface(new(Client))
)

type Option = lua.Options

// New 新建实例 万次执行速度0.16秒(简易函数)
// 能覆盖自带函数
// 可以协程执行,不能声明多个对象协程执行,数据会冲突
// 堆积执行数量过高会直接exit退出,捕捉不到
func New(op ...Option) *Client {
	options := []lua.Options(nil)
	for _, v := range op {
		options = append(options, v)
	}
	L := lua.NewState(options...)
	return &Client{
		client: L,
		mProto: make(map[string]*lua.FunctionProto),
	}
}

type Client struct {
	client *lua.LState
	mu     sync.RWMutex
	mProto map[string]*lua.FunctionProto
}

/*
Exec 执行文本
local test=1
test=test+1
return test

map类型
local Param = { ['id'] = n,['name'] = 'jyj' }

不加function运行更方便
*/
func (this *Client) Exec(text string) (*conv.Var, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if err := this.client.DoString(text); err != nil {
		return conv.Nil(), dealErr(err)
	}
	result := this.client.Get(-1)
	if result == nil || result == lua.LNil {
		return conv.Nil(), nil
	}
	this.client.Pop(1)
	return conv.New(result.String()), nil
}

// GetVar 获取全局变量
func (this *Client) GetVar(key string) *conv.Var {
	this.mu.RLock()
	defer this.mu.RUnlock()
	value := this.client.GetGlobal(key)
	if value == lua.LNil {
		return conv.Nil()
	}
	return conv.New(value.String())
}

// SetFunc 设置函数,相当于Set,为了书写方便
func (this *Client) SetFunc(key string, value script.Func) error {
	return this.Set(key, value)
}

// Set 设置全局变量 禁止并发设置
func (this *Client) Set(key string, value interface{}) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	switch fn := value.(type) {
	case script.Func:
		value = this.toFunc(fn)
	case func(*script.Args) interface{}:
		value = this.toFunc(fn)
	case func():
		value = this.toFunc(func(*script.Args) interface{} {
			fn()
			return nil
		})
	}
	this.client.SetGlobal(key, this.Value(value))
	return nil
}

// Close 关闭解释器
func (this *Client) Close() error {
	this.client.Close()
	return nil
}

func (this *Client) toFunc(fn script.Func) lua.LGFunction {
	return func(call *lua.LState) int {
		args := []*conv.Var(nil)
		for i := 1; i <= call.GetTop(); i++ {
			args = append(args, conv.New(call.Get(i).String()))
		}
		arg := &script.Args{
			This: this,
			Args: args,
		}
		call.Push(this.Value(fn(arg)))
		return 1
	}
}
