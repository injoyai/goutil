package oss

import (
	"testing"
	"time"
)

func TestSizeString(t *testing.T) {
	t.Log(SizeString(986))
	t.Log(SizeString(1024))
	t.Log(SizeString(2024, 0))
	t.Log(SizeString(20240))
	t.Log(SizeString(314 * 1024 * 1024))
	t.Log(SizeSpendString(314*1024*1024, time.Second))
	t.Log(SizeSpendString(314*5, time.Second))
}
