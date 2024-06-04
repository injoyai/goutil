package oss

import (
	"bytes"
	"testing"
	"time"
)

func TestSizeString(t *testing.T) {
	t.Log(SizeString(986))                       //986.0B
	t.Log(SizeString(1024))                      //1.0KB
	t.Log(SizeString(2024, 0))                   //2KB
	t.Log(SizeString(20240))                     //19.8KB
	t.Log(SizeString(314 * 1024 * 1024))         //314.0MB
	t.Log(SizeSpeed(314*1024*1024, time.Second)) //314.0MB/s
	t.Log(SizeSpeed(314*5, time.Second))         //1.5KB/s
}

func TestNew(t *testing.T) {
	New("./test.txt", bytes.NewBuffer([]byte("123456")))
	New("./test2.txt", "123456")
}

func TestSizeUnit(t *testing.T) {
	t.Log(SizeUnit(314 * 1024 * 1024))

	//t.Log(SizeUnit(math.MaxUint64))
}
