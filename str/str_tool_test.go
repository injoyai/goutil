package str

import (
	"strings"
	"testing"
)

func TestMustSplitN(t *testing.T) {
	testStr, testSep := "a,b,c", ","
	t.Log(strings.SplitN(testStr, testSep, 10))
	t.Log(MustSplitN(testStr, testSep, 3))
	t.Log(MustSplitN(testStr, testSep, 2))
	t.Log(len(MustSplitN(testStr, testSep, 2)))
	t.Log(MustSplitN(testStr, testSep, 1))
	t.Log(MustSplitN(testStr, testSep, 0))
	t.Log(MustSplitN(testStr, testSep, -1))
}

func TestBytes(t *testing.T) {
	s := "xxx"
	t.Logf("%p", &s)
	t.Logf("%p", &s)
	bs := Bytes(s)
	t.Logf("%p", bs)
	t.Log(bs)
	t.Log(string(bs))
	bs[0] = 48

	t.Log(s)
}
