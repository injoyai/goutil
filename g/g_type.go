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

type KV struct {
	K string      `json:"key"`
	V interface{} `json:"value"`
	L string      `json:"label,omitempty"`
}

type Msg struct {
	Type string      `json:"type"`           //请求类型,例如测试连接ping,写入数据write... 推荐请求和响应通过code区分
	Code int         `json:"code,omitempty"` //请求结果,推荐 请求:0(或null)  响应: 200成功,500失败... 同http好记一点
	UID  string      `json:"uid,omitempty"`  //消息的唯一ID,例如UUID
	Data interface{} `json:"data,omitempty"` //请求响应的数据
	Msg  string      `json:"msg,omitempty"`  //消息
}

func (this *Msg) IsRequest() bool {
	return this.Code == 0
}

func (this *Msg) IsResponse() bool {
	return this.Code != 0
}

func (this *Msg) Response(code int, data interface{}, msg string) *Msg {
	return &Msg{
		Type: this.Type,
		Code: code,
		UID:  this.UID,
		Data: data,
		Msg:  msg,
	}
}

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

//========================================Interface========================================

type Stringer interface{ String() string }

type GoStringer interface{ GoString() string }
