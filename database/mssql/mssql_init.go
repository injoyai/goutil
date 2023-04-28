package mssql

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/injoyai/goutil/database/xorms"
)

func NewXorm(op *xorms.Option) *xorms.Engine {
	return xorms.New(&xorms.Config{
		Type:        "mssql",
		DSN:         op.DSN,
		FieldSync:   op.FieldSync,
		TablePrefix: op.TablePrefix,
	})
}

func NewXormWithDSN(dsn string) *xorms.Engine {
	return NewXorm(&xorms.Option{
		DSN:         dsn,
		FieldSync:   true,
		TablePrefix: "",
	})
}
