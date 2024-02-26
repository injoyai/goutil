package oss

import (
	"bytes"
	"math"
	"testing"
	"time"
)

func TestSizeString(t *testing.T) {
	t.Log(SizeString(986))
	t.Log(SizeString(1024))
	t.Log(SizeString(2024, 0))
	t.Log(SizeString(20240))
	t.Log(SizeString(314 * 1024 * 1024))
	t.Log(SizeSpeed(314*1024*1024, time.Second))
	t.Log(SizeSpeed(314*5, time.Second))
}

func TestNew(t *testing.T) {
	New("./test.txt", bytes.NewBuffer([]byte("123456")))
	New("./test2.txt", "123456")
}

func TestSizeUnit(t *testing.T) {
	t.Log(SizeUnit(314 * 1024 * 1024))

	t.Log(SizeUnit(math.MaxUint64))
}
