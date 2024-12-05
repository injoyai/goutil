package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"io/ioutil"
	"path/filepath"
)

func newFile(name, group string) *File {
	data := &File{
		name:  name,
		group: group,
	}
	bs, _ := ioutil.ReadFile(data.Filename())
	data.Map = conv.NewMap(bs)
	return data
}

type File struct {
	name  string
	group string
	*conv.Map
}

// Name 名字
func (this *File) Name() string {
	return this.name
}

// Clear 清空数据
func (this *File) Clear() *File {
	*this.Map = *conv.NewMap(nil)
	return this
}

// GetAndSetByExtend
// 根据conv.Extend获取数据,不存在则取File中的数据,并设置到File中
// 尝试从用户那里获取数据,存在则覆盖
func (this *File) GetAndSetByExtend(key string, extend conv.Extend) interface{} {
	old := this.GetInterface(key)
	val := extend.GetInterface(key, old)
	this.Set(key, val)
	return val
}

// Set 设置参数
func (this *File) Set(key string, val interface{}) *File {
	this.Map.Set(key, val)
	return this
}

// SetMap 批量设置参数
func (this *File) SetMap(m map[string]interface{}) *File {
	for k, v := range m {
		this.Set(k, v)
	}
	return this
}

// Del 删除参数
func (this *File) Del(key string) *File {
	this.Map.Del(key)
	return this
}

// Save 保存配置文件,存在则覆盖
func (this *File) Save() error {
	filename := this.Filename()
	return oss.New(filename, this.Map.String())
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
