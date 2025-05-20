package lua

import (
	"github.com/injoyai/goutil/script"
	"testing"
)

func TestSpend(t *testing.T) {
	l := New()
	l.Set("add", Add)
	l.Set("a", 100)

	fun := `
b=10
c=add(a,b)
return c+1
`
	for i := 0; i < 10000; i++ {
		_, err := l.Exec(fun)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestNew(t *testing.T) {
	l := New()
	l.Set("add", Add)
	l.Set("a", 100)

	fun := `
b=10
c=add(a,b)
return c+1
`
	result, err := l.Exec(fun)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func Add(i *script.Args) any {
	return i.GetInt(1) + i.GetInt(2)
}
