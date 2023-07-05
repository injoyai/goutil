package oss

import "testing"

func TestDir(t *testing.T) {
	t.Log(ExecName())
	t.Log(ExecDir())
	t.Log(FuncName())
	t.Log(FuncDir())
	t.Log(UserDir())
	t.Log(UserDataDir(DefaultName))
}

func TestUserStartupDir(t *testing.T) {
	t.Log(UserStartupDir())
}
