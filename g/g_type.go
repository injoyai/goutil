package g

import (
	"github.com/injoyai/base/bytes"
	"github.com/injoyai/conv"
	json "github.com/json-iterator/go"
	"sort"
)

type (
	Var   = conv.Var
	DMap  = conv.Map
	Any   = interface{}
	List  []interface{}
	Bytes = bytes.Entity
	M     = Map
)

//========================================Map========================================

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

func (this Map) GetFloat32(key string) float32 {
	return conv.Float32(this[key])
}

func (this Map) GetFloat64(key string) float64 {
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

//========================================Maps========================================

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
