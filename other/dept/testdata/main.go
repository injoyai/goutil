package main

import (
	. "github.com/injoyai/goutil/other/dept"
	"github.com/injoyai/logs"
)

func main() {
	m := NewManage()
	m.Set(&Dept{
		ID:       1,
		ParentID: 0,
		Name:     "部门1",
	}, &Dept{
		ID:       2,
		ParentID: 1,
		Name:     "部门11",
	}, &Dept{
		ID:       3,
		ParentID: 1,
		Name:     "部门12",
	}, &Dept{
		ID:       4,
		ParentID: 2,
		Name:     "部门21",
	}, &Dept{
		ID:       5,
		ParentID: 2,
		Name:     "部门22",
	}, &Dept{
		ID:       6,
		ParentID: 5,
		Name:     "部门51",
	}, &Dept{
		ID:       7,
		ParentID: 6,
		Name:     "部门61",
	})
	{
		logs.Debug("GetChildrenAll")
		for _, v := range m.GetChildrenAll(2) {
			logs.Debug(*v)
		}
	}
	{
		logs.Debug("GetTree")
		top, _ := m.GetTree(2)
		logs.Debug(top.Children)
		for _, k := range top.Children {
			logs.Debug(*k)
			logs.Debug(k.Children)
		}
	}
}
