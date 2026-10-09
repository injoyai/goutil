//go:build (darwin || linux) && !cgo
// +build darwin linux
// +build !cgo

package share

import (
	"errors"
)

// shmi 非CGO环境兜底实现, 共享内存需要CGO支持
type shmi struct {
	name string
	size int32
}

func create(name string, size int32) (*shmi, error) {
	return nil, errors.New("共享内存需要启用CGO")
}

func open(name string, size int32) (*shmi, error) {
	return nil, errors.New("共享内存需要启用CGO")
}

func (o *shmi) close() error {
	return errors.New("共享内存需要启用CGO")
}

func (o *shmi) readAt(p []byte, off int64) (n int, err error) {
	return 0, errors.New("共享内存需要启用CGO")
}

func (o *shmi) writeAt(p []byte, off int64) (n int, err error) {
	return 0, errors.New("共享内存需要启用CGO")
}
