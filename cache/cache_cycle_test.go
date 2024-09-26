package cache

import (
	"testing"
)

func TestCycle_Padding(t *testing.T) {
	x := NewCycle(20)
	x.Padding(101)
	t.Log(x.List())
}
