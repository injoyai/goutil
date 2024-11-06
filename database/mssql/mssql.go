package mssql

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/injoyai/goutil/database/xorms"
)

func NewXorm(dsn string, options ...xorms.Option) (*xorms.Engine, error) {
	return xorms.NewMssql(dsn, options...)
}
