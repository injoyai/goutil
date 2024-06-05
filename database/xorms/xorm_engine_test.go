package xorms

import (
	"testing"
)

func TestNewSqlite(t *testing.T) {
	e, err := NewSqlite("./test.db")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(e.Ping())
	type Person struct {
		Name string
	}
	t.Log(e.Sync2(new(Person)))
	//t.Log(e.Insert(Person{Name: "injoy"}))
	result := []*Person{}
	t.Log(e.Find(&result))
	for _, v := range result {
		t.Log(*v)
	}
}

func TestNewMysql(t *testing.T) {
	e, err := NewMysql("test:RxTnAppSbsRc4jpJ@tcp(192.168.10.23:3306)/test?charset=utf8mb4")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(e.Ping())
}
