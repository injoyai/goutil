/*

	不同进程之间的内存共享

	copy from https://github.com/hidez8891/shm

*/

package share

import (
	"fmt"
	"io"
)

// Cover 覆盖更新,读取/写入会从下标0开始
type Cover struct {
	*Memory
}

// NewCover 新建共享内存,覆盖更新
func NewCover(name string, size int) (*Cover, error) {
	m, err := NewMemory(name, size)
	return &Cover{m}, err
}

// Write from same position
func (this *Cover) Write(p []byte) (int, error) {
	return this.Memory.WriteAt(p, 0)
}

// Read from same position
func (this *Cover) Read(p []byte) (int, error) {
	return this.Memory.ReadAt(p, 0)
}

// Close the shared memory
func (this *Cover) Close() error {
	return this.Memory.Close()
}

// Memory is shared memory struct
type Memory struct {
	m   *shmi
	pos int64
}

// NewMemory will create shared memory
func NewMemory(name string, size int) (*Memory, error) {
	m, err := create(name, int32(size))
	if err != nil {
		return nil, err
	}
	return &Memory{m, 0}, nil
}

// Close will close & discard shared memory
func (o *Memory) Close() (err error) {
	if o.m != nil {
		err = o.m.close()
		if err == nil {
			o.m = nil
		}
	}
	return err
}

// Read is read shared memory (current position)
func (o *Memory) Read(p []byte) (n int, err error) {
	n, err = o.ReadAt(p, o.pos)
	if err != nil {
		return 0, err
	}
	o.pos += int64(n)
	return n, nil
}

// ReadAt is read shared memory (offset)
func (o *Memory) ReadAt(p []byte, off int64) (n int, err error) {
	return o.m.readAt(p, off)
}

// Seek is move read/write position at shared memory
func (o *Memory) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		offset += int64(0)
	case io.SeekCurrent:
		offset += o.pos
	case io.SeekEnd:
		offset += int64(o.m.size)
	}
	if offset < 0 || offset >= int64(o.m.size) {
		return 0, fmt.Errorf("invalid offset")
	}
	o.pos = offset
	return offset, nil
}

// Write is write shared memory (current position)
func (o *Memory) Write(p []byte) (n int, err error) {
	n, err = o.WriteAt(p, o.pos)
	if err != nil {
		return 0, err
	}
	o.pos += int64(n)
	return n, nil
}

// WriteAt will write shared memory (offset)
func (o *Memory) WriteAt(p []byte, off int64) (n int, err error) {
	return o.m.writeAt(p, off)
}
