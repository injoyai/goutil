package excel

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := (&Export{}).Set(func() (list [][]any) {
		list = append(list, []any{"序号", "名字"})
		list = append(list, []any{1, "哈哈"})
		return
	}())
	t.Log(s)
}

func TestNew1(t *testing.T) {
	s := make(Export)
	s.Set(func() (list [][]any) {
		list = append(list, []any{"序号", "名字"})
		list = append(list, []any{1, "哈哈"})
		return
	}())
	t.Log(s.Buffer())
}
