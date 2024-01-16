package g

import "testing"

func TestRandString(t *testing.T) {
	t.Log(RandString(10))
	t.Log(RandString(10))
	t.Log(RandString(10))
	t.Log(RandString(10, "01"))
}

func TestRandFloat(t *testing.T) {
	t.Log(RandInt(0, 10))
	t.Log(RandFloat(0, 10))
}
