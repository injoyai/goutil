package oss

import "testing"

func TestCopy(t *testing.T) {
	err := Copy("./oss_dir.go", "./oss_dir_copy.go")
	t.Log(err)
}
