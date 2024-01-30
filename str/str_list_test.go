package str

import (
	"reflect"
	"testing"
)

func TestCut(t *testing.T) {
	{
		ls := List{"a", "b", "c"}
		x := ls[1:] //底层指针是同一个
		x[0] = "d"
		t.Log(ls) //[a d c]
		t.Log(x)  //[d c]
		if ls.Join(",") != "a,d,c" {
			t.Fail()
		}
		if x.Join(",") != "d,c" {
			t.Fail()
		}
		t.Log(x.Get(3))            // "",false
		t.Log(x.MustGet(3, "def")) //def
		t.Log(x.GetFirst())        //d
		t.Log(x.GetLast())         //c
		t.Log(ls.Reverse())        //[c d a]
		t.Log(x.Reverse())         //[a d]
		t.Log(x.Copy().Equal(x))   //true
		t.Log(x.Equal(ls))         //false

	}
}

func TestList_Sort(t *testing.T) {
	ls := List{}
	ls = append(ls, "b")
	ls = append(ls, "a")
	ls = append(ls, "c")
	ls.Sort()
	t.Log(ls)
	if !reflect.DeepEqual(ls, List{"a", "b", "c"}) {
		t.Fail()
	}
}
