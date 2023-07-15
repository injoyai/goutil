package regexps

import (
	"testing"
)

func TestIsTest(t *testing.T) {
	//t.Log(Is(`((http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?)`, "http://www.baidu.om2"))

	t.Log(Is(`.(png|jpg)$`, "钱测试.jpg"))
	t.Log(Is(``, "钱测试.jpg"))

	//t.Log(Find("(2)[0-9]{3}(年)[0-9]{1,2}(月)", "2021年3月治村指数"))
}

func TestIsPwd(t *testing.T) {

	s := "1234567aA*"

	t.Log(Is(`^(?=.*\d)(?=.*[a-zA-Z])(?=.*[~!@#$%^&*])[\da-zA-Z~!@#$%^&*]{8,}$`, s))
	t.Log(Is(`[0-9]{1,}`, s))
	t.Log(Is(`[a-z]{1,}`, s))
	t.Log(Is(`[A-Z]{1,}`, s))
	t.Log(Is(`[A-Z]{1,}`, s) && Is(`[0-9]{1,}`, s) && Is(`[a-z]{1,}`, s))
}

func TestIsPwd1(t *testing.T) {
	t.Log(IsPwd("000"))
	t.Log(IsPwd("000000"))
	t.Log(IsPwd("000123."))
	t.Log(IsPwd("000abc"))
	t.Log(IsPwd("000aaa."))
	t.Log(IsPwd("a000Abc0."))

}
