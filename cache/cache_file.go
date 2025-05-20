package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"os"
	"path/filepath"
)

func newFile(name, group string) *File {
	data := &File{
		name:  name,
		group: group,
		Safe:  maps.NewSafe(),
	}
	bs, _ := os.ReadFile(data.Filename())
	m := make(map[string]any)
	_ = json.Unmarshal(bs, &m)
	for i, v := range m {
		data.Set(i, v)
	}
	data.Extend = conv.NewExtend(data)
	return data
}

type File struct {
	name  string
	group string
	*maps.Safe
	conv.Extend
}

// Name 名字
func (this *File) Name() string {
	return this.name
}

// GetVar 实现接口
func (this *File) GetVar(key string) *conv.Var {
	return conv.New(this.MustGet(key))
}

// DMap *conv.Map
func (this *File) DMap() *conv.Map {
	return conv.NewMap(this.Safe.GMap())
}

// Clear 清空数据
func (this *File) Clear() *File {
	this.Safe = maps.NewSafe()
	return this
}

// GetAndSetByExtend
// 根据conv.Extend获取数据,不存在则取File中的数据,并设置到File中
// 尝试从用户那里获取数据,存在则覆盖
func (this *File) GetAndSetByExtend(key string, extend conv.Extend) any {
	old := this.GetInterface(key)
	val := extend.GetInterface(key, old)
	this.Set(key, val)
	return val
}

// Set 设置参数
func (this *File) Set(key string, val any) *File {
	this.Safe.Set(key, val)
	return this
}

// SetMap 批量设置参数
func (this *File) SetMap(m map[string]any) *File {
	for k, v := range m {
		this.Set(k, v)
	}
	return this
}

// Del 删除参数
func (this *File) Del(key string) *File {
	this.Safe.Del(key)
	return this
}

// Save 保存配置文件,存在则覆盖
func (this *File) Save() error {
	filename := this.Filename()
	bs, err := json.Marshal(this.Safe.GMap())
	if err != nil {
		return err
	}
	return oss.New(filename, bs)
}

// Cover 覆盖
func (this *File) Cover() error {
	return this.Save()
}

func (this *File) Filename() string {
	fileDir, filename := DefaultDir, this.name
	if dir, file := filepath.Split(this.name); len(dir) > 0 {
		fileDir, filename = dir, file
	}
	h := md5.New()
	h.Write([]byte(filename))
	filename = fmt.Sprintf("%s@%s", this.group, hex.EncodeToString(h.Sum(nil)))
	return filepath.Join(fileDir, filename)
}
