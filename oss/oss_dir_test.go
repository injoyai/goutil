package oss

import (
	"fmt"
	"testing"
)

func TestDir(t *testing.T) {
	t.Log(ExecName())
	t.Log(ExecDir("/a/b"))
	t.Log(FuncName())
	t.Log(FuncDir())
	t.Log(UserDir())
	t.Log(UserDataDir(DefaultName))
}

func TestUserStartupDir(t *testing.T) {
	t.Log(UserStartupDir())
}

func TestUserDesktopDir(t *testing.T) {
	t.Log(UserDesktopDir())
}

func TestUserDir(t *testing.T) {
	t.Log(UserDir())
}

func TestRangeFileInfo(t *testing.T) {
	t.Log(ReadFilenames("./", -1))
}

func TestReadTree(t *testing.T) {
	{
		d, err := ReadTree("./", 2)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("\n%s", d)
	}

	{
		d, err := ReadTree("./compress")
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("\n%s", d)
	}

	{
		d, err := ReadTree("./compress/gzip")
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("\n%s", d)
	}
}

func TestReadTreeFormat(t *testing.T) {
	d, err := ReadTree("./compress")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("\n%s", d.Format("- ", "- ", "> ", "> "))

}
