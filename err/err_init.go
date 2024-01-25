package err

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type Entity struct {
	items []Item
}

func New(msg string) *Item {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return &Item{
		funcName: funcName,
		file:     file,
		line:     line,
		Msg:      msg,
	}
}

type Item struct {
	funcName string
	file     string
	line     int
	Msg      string
}

func (this *Item) String() string {
	return this.Error()
}

func (this *Item) Error() string {
	return fmt.Sprintf("\n\t>>> %s:%d: %s", filepath.Base(this.file), this.line, this.Msg)
}
