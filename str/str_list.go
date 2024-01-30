package str

import (
	"sort"
	"strings"
)

type List []string

func (this List) Len() int {
	return len(this)
}

func (this List) TotalLen() int {
	length := 0
	for _, v := range this {
		length += len(v)
	}
	return length
}

func (this List) Cap() int {
	return cap(this)
}

func (this List) Equal(ls []string) bool {
	if (this == nil) != (ls == nil) {
		return false
	}
	if this.Len() != len(ls) {
		return false
	}
	for i := 0; i < this.Len(); i++ {
		if this[i] != ls[i] {
			return false
		}
	}
	return true
}

func (this List) Upper() {
	for i, v := range this {
		this[i] = strings.ToUpper(v)
	}
}

func (this List) Lower() {
	for i, v := range this {
		this[i] = strings.ToLower(v)
	}
}

func (this List) Join(sep string) string {
	return strings.Join(this, sep)
}

func (this List) Joinln() string {
	return this.Join("\n")
}

func (this List) Split(sep string) List {
	ls := List{}
	for _, v := range this {
		ls = append(ls, strings.Split(v, sep)...)
	}
	return ls
}

func (this List) Sort() List {
	sort.Strings(this)
	return this
}

// Swap 实现排序接口,2个元素都存在则交换元素
func (this List) Swap(i, j int) {
	i = this.getIdx(i)
	j = this.getIdx(j)
	if i >= 0 && j >= 0 {
		this[i], this[j] = this[j], this[i]
	}
}

// Reverse 倒序元素
func (this List) Reverse() List {
	for i := 0; i < this.Len()/2; i++ {
		this.Swap(i, this.Len()-1-i)
	}
	return this
}

// Replace 替换自定元素
func (this List) Replace(idx int, v string) bool {
	if idx = this.getIdx(idx); idx >= 0 {
		this[idx] = v
		return true
	}
	return false
}

// Cut 裁剪元素,底层指针是同一个
func (this List) Cut(start, end int) List {
	if this.Len() == 0 {
		return this[:0]
	}
	start2 := this.getIdx(start)
	end2 := this.getIdx(end)
	if end > 0 && end2 < 0 {
		end2 = this.Len()
	}
	if start2 < 0 {
		start2 = 0
	}
	if end2 <= start2 {
		return this[:0]
	}
	return this[start2:end2]
}

func (this List) Copy() List {
	return append(List{}, this...)
}

func (this List) Exist(idx int) bool {
	return this.getIdx(idx) >= 0
}

func (this List) Get(idx int) (string, bool) {
	if idx = this.getIdx(idx); idx >= 0 {
		return this[idx], true
	}
	return "", false
}

// MustGet 获取元素,不存在返回nil
func (this List) MustGet(idx int, def ...string) string {
	if idx = this.getIdx(idx); idx >= 0 {
		return this[idx]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (this List) GetFirst(def ...string) string {
	return this.MustGet(0, def...)
}

func (this List) GetLast(def ...string) string {
	return this.MustGet(-1, def...)
}

// getIdx 处理下标,支持负数-1表示最后1个,同python
func (this List) getIdx(idx int) int {
	length := this.Len()
	if idx < length && idx >= 0 {
		return idx
	}
	if idx < 0 && -idx <= length {
		return length + idx
	}
	return -1
}
