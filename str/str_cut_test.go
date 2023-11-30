package str

import "testing"

func TestCutLeast(t *testing.T) {
	t.Log(len(CutLeast("123", 3)))
	t.Log(len(CutLeast("123", 4)))
	t.Log(len(CutLeast("123", -1)))
	t.Log(len(CutLeast("123", 0)))
	t.Log(len(CutLeast("123", 1)))
}
