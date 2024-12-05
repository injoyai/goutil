package cache

import "github.com/injoyai/conv"

var DefaultDir = "./data/cache/"

// NewFile 新建文件缓存
// 万次读写速度4.18秒
// 万次协程读写速度2.21秒
func NewFile(name string, groups ...string) *File {
	group := conv.GetDefaultString("var", groups...)
	return newFile(name, group)
}
