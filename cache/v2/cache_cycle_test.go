package cache

import (
	"testing"
)

func TestNewCycle(t *testing.T) {
	x := NewCycle[int](20)
	x.Append(100)
	t.Log(x.List())
	x.Padding(101)
	t.Log(x.List())
	x.Append(102)
	t.Log(x.List())
}
