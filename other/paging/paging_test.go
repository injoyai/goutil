package paging

import (
	"testing"
)

func TestLimit(t *testing.T) {
	t.Log(Limit(List{1, 2, 3}, 2, 1))  //[2 3]
	t.Log(Limit(List{1, 2, 3}, 3, 1))  //[2 3]
	t.Log(Limit(List{1, 2, 3}, 4, 1))  //[2 3]
	t.Log(Limit(List{1, 2, 3}, 4, 2))  //[3]
	t.Log(Limit(List{1, 2, 3}, 4, 3))  //[]
	t.Log(Limit(List{1, 2, 3}, 4, 4))  //[]
	t.Log(Limit(List{1, 2, 3}, 2, -1)) //[1]
	t.Log(Limit(List{1, 2, 3}, 2))     //[1 2]
	t.Log(Limit(List{1, 2, 3}, 4))     //[1 2 3]
	t.Log(Limit(List{1, 2, 3}, -1))    //[]
}

type List []int

func (this List) Len() int {
	return len(this)
}

func (this List) Cut(i1 int, i2 int) interface{} {
	return this[i1:i2]
}
