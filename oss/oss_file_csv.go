package oss

import (
	"encoding/csv"
	"github.com/injoyai/conv"
	"os"
)

type CSVFile struct {
	*os.File
	*csv.Writer
}

func (this *CSVFile) Write(v ...any) (err error) {
	switch len(v) {
	case 0:
	case 1:
		err = this.Writer.Write(conv.Strings(v[0]))
	default:
		err = this.Writer.Write(conv.Strings(v))
	}
	return
}

func (this *CSVFile) WriteFlush(v ...any) error {
	defer this.Flush()
	return this.Write(v...)
}
