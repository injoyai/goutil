package upload

import (
	"bytes"
	"testing"
)

func TestLocal_Save(t *testing.T) {
	i := NewLocal("./data/upload/")
	i.Save("d.txt", bytes.NewReader([]byte("hello")))
	_, err := i.Save("./a/b/c.txt", bytes.NewReader([]byte("hello")))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(i.List())
}
