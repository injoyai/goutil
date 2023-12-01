package upload

import (
	"bytes"
	"testing"
)

func TestLocal_Save(t *testing.T) {
	i := NewLocal(false)
	s, err := i.Save("./a/b/c.txt", bytes.NewReader([]byte("hello")))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(s)
}
