package share

import (
	"testing"
	"time"
)

const (
	name = "test"
	size = 20
)

func TestNewMemoryWrite(t *testing.T) {
	m, err := NewMemory(name, size)
	if err != nil {
		t.Error(err)
	}
	for {
		time.Sleep(time.Second * 3)
		n, err := m.WriteAt([]byte(time.Now().String()), 0)
		if err != nil {
			t.Error(err)
		}
		t.Log(n)
	}
}

func TestNewMemoryRead(t *testing.T) {
	m, err := NewMemory(name, size)
	if err != nil {
		t.Error(err)
	}
	for {
		time.Sleep(time.Second)
		p := make([]byte, 20)
		_, err := m.ReadAt(p, 0)
		if err != nil {
			t.Error(err)
		}
		t.Log(string(p))
	}
}
