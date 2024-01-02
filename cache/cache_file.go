package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"io/ioutil"
	"path/filepath"
)

type File struct {
	name string
	tag  string
	*maps.Safe
	conv.Extend
}

func (this *File) GetVar(key string) *conv.Var {
	return conv.New(this.MustGet(key))
}

func newFile(name, tag string) *File {
	data := &File{
		name: name,
		tag:  tag,
		Safe: maps.NewSafe(),
	}
	bs, _ := ioutil.ReadFile(data.getPath(tag, name))
	m := make(map[string]interface{})
	_ = json.Unmarshal(bs, &m)
	for i, v := range m {
		data.Set(i, v)
	}
	data.Extend = conv.NewExtend(data)
	return data
}

// Name 名字
func (this *File) Name() string {
	return this.name
}

// Clear 清空数据
func (this *File) Clear() *File {
	this.Safe = maps.NewSafe()
	return this
}

// GetAndSetByExtend 根据conv.Extend获取数据,不存在则取File中的数据,并设置到File中
func (this *File) GetAndSetByExtend(key string, extend conv.Extend) interface{} {
	old := this.GetInterface(key)
	val := extend.GetInterface(key, old)
	this.Set(key, val)
	return val
}

// Set 设置参数
func (this *File) Set(key string, val interface{}) *File {
	this.Safe.Set(key, val)
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
	this.Safe.Del(key)
	return this
}

func (this *File) Save() error {
	filename := this.getPath(this.tag, this.name)
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

func (this *File) getPath(tag, name string) string {
	fileDir, filename := DefaultDir, name
	if dir, file := filepath.Split(name); len(dir) > 0 {
		fileDir, filename = dir, file
	}
	h := md5.New()
	h.Write([]byte(filename))
	return fileDir + fmt.Sprintf("%s@%s", tag, hex.EncodeToString(h.Sum(nil)))
}
