package cache

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"os"
	"path/filepath"
)

// NewFile 新建文件缓存
// 万次读写速度4.18秒
// 万次协程读写速度2.21秒
func NewFile(filename string) (*File, error) {
	data := &File{
		filename: filename,
	}
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data.Map = conv.NewMap(bs)
	return data, nil
}

type File struct {
	filename string
	*conv.Map
}

// Name 名字
func (this *File) Name() string {
	return filepath.Base(this.filename)
}

// Clear 清空数据
func (this *File) Clear() *File {
	*this.Map = *conv.NewMap(nil)
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
	this.Map.Set(key, val)
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
	this.Map.Del(key)
	return this
}

// Save 保存配置文件,存在则覆盖
func (this *File) Save() error {
	return oss.New(this.filename, this.Map.String())
}
