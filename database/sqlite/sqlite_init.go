package sqlite

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/injoyai/goutil/database/xorms"
	"os"
	"path/filepath"
)

func NewXorm(op *xorms.Option) *xorms.Engine {
	dir, _ := filepath.Split(op.DSN)
	_ = os.MkdirAll(dir, 0777)
	db := xorms.New(&xorms.Config{
		Type:        "sqlite",
		DSN:         op.DSN,
		FieldSync:   op.FieldSync,
		TablePrefix: op.TablePrefix,
	})
	if db.Err() == nil {
		db.SetMaxOpenConns(1)
	}
	return db
}

func NewXormWithPath(path string) *xorms.Engine {
	return NewXorm(&xorms.Option{
		DSN:       path,
		FieldSync: true,
	})
}
