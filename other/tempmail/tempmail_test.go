package tempmail

import (
	"testing"
)

func TestNew(t *testing.T) {
	c := New("xxx")
	ls, err := c.List(10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range ls {
		t.Log(v)
	}
	d, err := c.Details(2969874524)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(d)
}
