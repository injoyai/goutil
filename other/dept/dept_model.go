package dept

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"sort"
)

func NewManage() *Manage {
	m := &Manage{
		Safe: maps.NewSafe(),
	}
	return m
}

type Manage struct {
	*maps.Safe
}

// IsRoot 是否根部门
func (this *Manage) IsRoot(id any) bool {
	val, ok := this.Get(id)
	if ok {
		return conv.IsDefault(val.ParentID)
	}
	return false
}

// Exist 是否存在
func (this *Manage) Exist(id any) bool {
	_, ok := this.Get(id)
	return ok
}

// Get 获取部门
func (this *Manage) Get(id any) (*Dept, bool) {
	val, ok := this.Safe.Get(id)
	if ok {
		return val.(*Dept), true
	}
	return nil, false
}

func (this *Manage) Set(val ...*Dept) {
	for _, v := range val {
		this.Safe.Set(v.ID, v)
	}
}

func (this *Manage) GetFather(id any) (*Dept, bool) {
	e, ok := this.Get(id)
	if !ok {
		return nil, false
	}
	return this.Get(e.ParentID)
}

func (this *Manage) GetFatherAll(id any) []*Dept {
	result := []*Dept(nil)
	for {
		father, ok := this.GetFather(id)
		if !ok {
			break
		}
		result = append(result, father)
		id = father.ID
	}
	return result
}

func (this *Manage) GetChildrenID(id any, level int) []any {
	result := []any(nil)
	for _, v := range this.GetChildren(id, level) {
		result = append(result, v.ID)
	}
	return result
}

func (this *Manage) GetChildren(id any, level int) (list []*Dept) {
	if level == 0 || !this.Exist(id) {
		return
	}
	this.Range(func(key, value any) bool {
		if e := value.(*Dept); e.ParentID == id {
			list = append(list, e)
			list = append(list, this.GetChildren(e.ID, level-1)...)
		}
		return true
	})
	return
}

func (this *Manage) GetChildrenAll(id any) []*Dept {
	return this.GetChildren(id, -1)
}

// GetTree 获取整理成树状结构
func (this *Manage) GetTree(id any) (*Dept, bool) {
	e, ok := this.Get(id)
	if !ok {
		return nil, false
	}
	child := this.GetChildrenAll(id)
	m := map[any][]*Dept{}
	for _, v := range child {
		m[v.ParentID] = append(m[v.ParentID], v)
	}
	e.Children = m[e.ID]
	for i, v := range child {
		list := m[v.ID]
		sort.Sort(Depts(list))
		child[i].Children = list
	}
	return e, true
}

type Dept struct {
	ID       any
	ParentID any
	Name     string
	Children []*Dept
	Data     any
}

type Depts []*Dept

func (this Depts) Len() int {
	return len(this)
}

func (this Depts) Less(i, j int) bool {
	if conv.IsInt(this[i].ID) {
		return conv.Int64(this[i].ID) < conv.Int64(this[j].ID)
	}
	if conv.IsFloat(this[i].ID) {
		return conv.Float64(this[i].ID) < conv.Float64(this[j].ID)
	}
	return conv.String(this[i].ID) < conv.String(this[j].ID)
}

func (this Depts) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
