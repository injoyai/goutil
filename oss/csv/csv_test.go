package csv

import (
	"testing"

	"github.com/injoyai/goutil/oss"
)

func TestExport(t *testing.T) {
	buf, err := Export([][]string{{"a", "b"}, {"c", "d"}})
	if err != nil {
		t.Fatal(err)
	}
	filename := "./test.csv"
	err = oss.New(filename, buf)
	if err != nil {
		t.Fatal(err)
	}

	// 读取文件验证内容
	content, err := Import(filename)
	if err != nil {
		t.Fatal(err)
	}

	buf, err = Export(content)
	if err != nil {
		t.Fatal(err)
	}
	err = oss.New(filename, buf)
	if err != nil {
		t.Fatal(err)
	}

}
