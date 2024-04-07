package g

import (
	"fmt"
	"github.com/injoyai/base/bytes"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/script/js"
	"github.com/injoyai/goutil/str"
	json "github.com/json-iterator/go"
	"sort"
	"sync"
)

type (
	Var     = conv.Var
	DMap    = conv.Map
	Any     = interface{}
	List    []interface{}
	Bytes   = bytes.Entity
	Strings = str.List
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
			ScriptPool = js.NewPool(20)
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

func (this Map) GetInt(key string) int {
	return conv.Int(this[key])
}

func (this Map) GetInt64(key string) int64 {
	return conv.Int64(this[key])
}

func (this Map) GetString(key string) string {
	return conv.String(this[key])
}

func (this Map) GetBool(key string) bool {
	return conv.Bool(this[key])
}

func (this Map) GetFloat(key string) float64 {
	return conv.Float64(this[key])
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

type Maps []Map

func (this Maps) Len() int {
	return len(this)
}

func (this Maps) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this Maps) Sort(fn func(i, j Map) bool) {
	sort.Sort(&_sortMaps{
		Maps: this,
		less: fn,
	})
}

type _sortMaps struct {
	Maps
	less func(i, j Map) bool
}

func (this *_sortMaps) Less(i, j int) bool {
	return this.less(this.Maps[i], this.Maps[j])
}

//========================================Key========================================

type Key string

func (this *Key) SetKey(key string) { *this = Key(key) }

func (this *Key) Set(key string) { *this = Key(key) }

func (this *Key) GetKey() string { return string(*this) }

func (this *Key) Get() string { return string(*this) }

//========================================Debugger========================================

type Debugger bool

func (this *Debugger) Debug(b ...bool) {
	*this = Debugger(len(b) == 0 || b[0])
}
