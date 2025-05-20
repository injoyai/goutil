package lua

import (
	"encoding/json"
	"github.com/yuin/gopher-lua"
	"reflect"
)

func (this *Client) Value(i any) lua.LValue {
	switch value := (i).(type) {
	case error:
		return lua.LString(value.Error())
	case int:
		return lua.LNumber(value)
	case int8:
		return lua.LNumber(value)
	case int16:
		return lua.LNumber(value)
	case int32:
		return lua.LNumber(value)
	case int64:
		return lua.LNumber(value)
	case uint:
		return lua.LNumber(value)
	case uint8:
		return lua.LNumber(value)
	case uint16:
		return lua.LNumber(value)
	case uint32:
		return lua.LNumber(value)
	case uint64:
		return lua.LNumber(value)
	case float32:
		return lua.LNumber(value)
	case float64:
		return lua.LNumber(value)
	case string:
		return lua.LString(value)
	case bool:
		return lua.LBool(value)
	default:
		r := reflect.ValueOf(i)
		switch r.Kind() {
		case reflect.Ptr:
			return this.Value(r.Elem().Interface())
		case reflect.Map:
			bs, _ := json.Marshal(r.Interface())
			m := make(map[string]any, r.Len())
			_ = json.Unmarshal(bs, &m)
			t := this.client.NewTable()
			for i, v := range m {
				t.RawSet(lua.LString(i), this.Value(v))
			}
			return t
		case reflect.Func:
			if value, ok := i.(lua.LGFunction); ok {
				return this.client.NewFunction(value)
			}
		}
		return lua.LNil
	}
}
