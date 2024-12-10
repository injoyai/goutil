package types

import (
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/script/js"
	"sync"
)

//========================================Type========================================

const (
	String Type = "string"
	Bool   Type = "bool"
	Int    Type = "int"
	Float  Type = "float"
	Array  Type = "array"
	Object Type = "object"
	Script Type = "script"
)

// Type 数据类型
type Type string

func (this Type) Int() int {
	switch this {
	case String:
		return 1
	case Bool:
		return 2
	case Int:
		return 3
	case Float:
		return 4
	case Object:
		return 5
	case Array:
		return 6
	case Script:
		return 7
	}
	return 0
}

func (this Type) Name() string {
	switch this {
	case Float:
		return "浮点"
	case Int:
		return "整数"
	case String:
		return "字符"
	case Bool:
		return "布尔"
	case Array:
		return "数组"
	case Object:
		return "对象"
	case Script:
		return "脚本"
	}
	return "未知"
}

func (this Type) Value(v interface{}) interface{} {
	switch this {
	case String:
		return conv.String(v)
	case Bool:
		return conv.Bool(v)
	case Int:
		return conv.Int(v)
	case Float:
		return conv.Float64(v)
	case Script:
		scriptPoolOnce.Do(func() {
			if ScriptPool == nil {
				ScriptPool = js.NewPool(20)
			}
		})
		val, _ := ScriptPool.Exec(conv.String(v))
		return val
	}
	return v
}

var (
	ScriptPool     *js.Pool
	scriptPoolOnce sync.Once
)

func (this Type) Check() error {
	switch this {
	case String, Bool, Int, Float, Array, Object, Script:
	default:
		return fmt.Errorf("未知数据类型:%s", this)
	}
	return nil
}

//========================================Debugger========================================

type Debugger bool

func (this *Debugger) Debug(b ...bool) {
	*this = Debugger(len(b) == 0 || b[0])
}

//========================================error========================================

// Err 错误,好处是能定义在const
type Err string

func (this Err) Error() string { return string(this) }
