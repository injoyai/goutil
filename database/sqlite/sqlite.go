package sqlite

import (
	_ "github.com/glebarez/go-sqlite"
	"github.com/injoyai/goutil/database/xorms"
)

func NewXorm(filename string, options ...xorms.Option) (*xorms.Engine, error) {
	return xorms.NewSqlite(filename, options...)
}
