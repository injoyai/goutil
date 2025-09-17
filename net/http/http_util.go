package http

import (
	"io"
)

type ListenWrite struct {
	io.Writer
	p       *Plan
	OnWrite func(*Plan)
}

func (this *ListenWrite) Write(p []byte) (n int, err error) {
	if this.p == nil {
		this.p = &Plan{
			Index:   0,
			Total:   0,
			Current: 0,
			Bytes:   nil,
		}
	}
	this.p.Index++
	this.p.Current += int64(len(p))
	this.p.Bytes = p
	if this.OnWrite != nil {
		this.OnWrite(this.p)
	}
	if this.p.Err != nil {
		return 0, this.p.Err
	}
	return this.Writer.Write(this.p.Bytes)
}

type Plan struct {
	Index   int64
	Total   int64
	Current int64
	Bytes   []byte
	Err     error
}

func (this *Plan) SetTotal(total int64) *Plan {
	this.Total = total
	return this
}
