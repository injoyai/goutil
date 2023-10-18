package oss

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/injoyai/conv"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	defaultPerm = 0777
)

// Read 读取文件
func Read(filename string) ([]byte, error) {
	return ReadBytes(filename)
}

// ReadReader 读取reader
func ReadReader(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

// ReadBytes 读取文件内容
func ReadBytes(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// ReadString 读取文件
func ReadString(filename string) (string, error) {
	bs, err := ReadBytes(filename)
	return string(bs), err
}

// ReadHEX 读取文件
func ReadHEX(filename string) (string, error) {
	bs, err := ReadBytes(filename)
	return hex.EncodeToString(bs), err
}

// ReadBase64 读取文件
func ReadBase64(filename string) (string, error) {
	bs, err := ReadBytes(filename)
	return base64.StdEncoding.EncodeToString(bs), err
}

// NewFile 新建文件
func NewFile(filename string) (*os.File, error) {
	return os.Create(filename)
}

// OpenFile 打开文件
func OpenFile(filename string) (*os.File, error) {
	return os.Open(filename)
}

// New 新建文件,会覆盖
func New(filename string, v ...interface{}) error {
	dir := filepath.Dir(filename)
	name := filepath.Base(filename)
	if err := os.MkdirAll(dir, defaultPerm); err != nil {
		return err
	}
	if len(name) == 0 {
		return nil
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if len(v) == 0 {
		return nil
	}
	data := []byte(nil)
	for _, k := range v {
		data = append(data, conv.Bytes(k)...)
	}
	_, err = f.Write(data)
	return err
}

// NewNotExist 如果不存在,则新建
func NewNotExist(filename string, v ...interface{}) error {
	if !Exists(filename) {
		return New(filename, v...)
	}
	return nil
}

// OpenFunc 打开文件,并执行函数
func OpenFunc(filename string, fn func(f *os.File) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return fn(f)
}

func WithOpen(filename string, fn func(f *os.File) error) error {
	return OpenFunc(filename, fn)
}

// WithCopyTo 打开文件,并写入到io.Writer
func WithCopyTo(filename string, writer io.Writer) error {
	return OpenFunc(filename, func(f *os.File) error {
		_, err := io.Copy(writer, f)
		return err
	})
}

// WithCopyFrom 打开文件,并从io.Reader写入
func WithCopyFrom(filename string, reader io.Reader) error {
	return OpenFunc(filename, func(f *os.File) error {
		_, err := io.Copy(f, reader)
		return err
	})
}
