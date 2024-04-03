package cache

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
)

var DefaultDir = "./data/cache/"

// NewFile 新建文件缓存
// 万次读写速度4.18秒
// 万次协程读写速度2.21秒
func NewFile(name string, groups ...string) *File {
	group := conv.GetDefaultString("var", groups...)
	return newFile(name, group)
}

// NewCycle 新建内存缓存,存储记录,循环使用
// 百万次读写速度4.31秒
// 百万次携程读写速度3.76秒
func NewCycle(num int) *Cycle {
	return newCycle(num)
}

// NewMap 新建内存缓存
// 千万次写速度 8.3s,
// 千万次读速度 2.3s
// 百万次读写速度1.2秒
// 百万次携程读写速度1.25秒
func NewMap(m ...maps.Map) *maps.Safe {
	return maps.NewSafe(m...)
}

// NewFileLog 文件日志
// 单文件万次写入速度 0.08s/万次
// 单文件百万次写入速度 2.11s/百万次
// 多文件千万次写入速度 19.88/千万次
func NewFileLog(cfg *FileLogConfig) *FileLog {
	return newFileLog(cfg)
}
