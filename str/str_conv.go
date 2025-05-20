package str

import (
	"encoding/json"
)

func DecodeJson(s string) any {
	s = `{"data":` + s + `}`
	m := make(map[string]any)
	_ = json.Unmarshal([]byte(s), &m)
	return m["data"]
}
