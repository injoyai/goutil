package oss

import (
	"testing"
)

func TestSizeString(t *testing.T) {
	t.Log(SizeString(986))
	t.Log(SizeString(1024))
	t.Log(SizeString(2024, 1))
	t.Log(SizeString(20240))
	t.Log(SizeString(314 * 1024 * 1024))
}
