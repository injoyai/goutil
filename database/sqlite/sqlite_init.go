package sqlite

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/injoyai/goutil/database/xorms"
	"os"
	"path/filepath"
)

func NewXorm(op *xorms.Option) (*xorms.Engine, error) {
	dir, _ := filepath.Split(op.DSN)
	_ = os.MkdirAll(dir, 0777)
	db, err := xorms.New(&xorms.Config{
		Type:        "sqlite",
		DSN:         op.DSN,
		FieldSync:   op.FieldSync,
		TablePrefix: op.TablePrefix,
	})
	if err != nil {
		return nil, err
	}
	//sqlite是文件数据库,只能打开一次(即一个连接)
	db.SetMaxOpenConns(1)
	return db, nil
}

func NewXormWithPath(path string) (*xorms.Engine, error) {
	return NewXorm(&xorms.Option{
		DSN:       path,
		FieldSync: true,
	})
}
