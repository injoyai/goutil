package mssql

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/injoyai/goutil/database/xorms"
)

func NewXorm(op *xorms.Option) (*xorms.Engine, error) {
	return xorms.New(&xorms.Config{
		Type:        "mssql",
		DSN:         op.DSN,
		FieldSync:   op.FieldSync,
		TablePrefix: op.TablePrefix,
	})
}

func NewXormWithDSN(dsn string) (*xorms.Engine, error) {
	return NewXorm(&xorms.Option{
		DSN:         dsn,
		FieldSync:   true,
		TablePrefix: "",
	})
}
