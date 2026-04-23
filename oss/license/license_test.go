package license

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	x := New("test")
	code, err := x.License(Code{
		End:    time.Now().Add(time.Hour * 24).Unix(),
		Expire: time.Now().AddDate(1, 0, 0).Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(code), code)
	_c := new(Code)
	err = x.decode(code, _c)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(_c)

	err = x.Activate(code)
	if err != nil {
		t.Fatal(err)
	}
	info, err := x.loadingInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}
