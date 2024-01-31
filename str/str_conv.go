package str

import (
	"encoding/json"
)

func DecodeJson(s string) interface{} {
	s = `{"data":` + s + `}`
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(s), &m)
	return m["data"]
}
