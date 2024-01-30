package str

import (
	"testing"
)

func TestInterface(t *testing.T) {
	x := DecodeJson(`["a","b","c"]`)
	t.Logf("%T: %v", x, x)
	x = DecodeJson(`1`)
	t.Logf("%T: %v", x, x)
	x = DecodeJson(`"aaa"`)
	t.Logf("%T: %v", x, x)
	x = DecodeJson(`{"a":1}`)
	t.Logf("%T: %v", x, x)
}
