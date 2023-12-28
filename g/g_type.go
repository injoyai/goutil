package g

import (
	"fmt"
	"github.com/injoyai/base/bytes"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/script/js"
	json "github.com/json-iterator/go"
	"sync"
)

type (
	Var   = conv.Var
	DMap  = conv.Map
	Any   = interface{}
	List  []interface{}
	Bytes = bytes.Entity
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
				ScriptPool = js.NewPool()
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

//========================================Map========================================

type M = Map

type Map map[string]interface{}

// Struct json Marshal
func (this Map) Struct(ptr interface{}) error {
	bs, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, ptr)
}

// Json map转json
func (this Map) String() string {
	return string(this.Bytes())
}

// Bytes map转字节
func (this Map) Bytes() bytes.Entity {
	bs, _ := json.Marshal(this)
	return bs
}

// GetVar 实现conv.Extend接口
func (this Map) GetVar(key string) *conv.Var {
	return conv.New(this[key])
}

// Merge 合并多个map
func (this Map) Merge(m ...Map) Map {
	for _, v := range m {
		for key, val := range v {
			this[key] = val
		}
	}
	return this
}

func (this Map) Conv() conv.Extend {
	return conv.NewExtend(this)
}

//========================================Debugger========================================

type Debugger bool

func (this *Debugger) Debug(b ...bool) {
	*this = Debugger(len(b) == 0 || b[0])
}

//========================================Interface========================================

type Stringer interface{ String() string }

type GoStringer interface{ GoString() string }
