package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/injoyai/goutil/database/xorms"
)

func NewXorm(dsn string, options ...xorms.Option) (*xorms.Engine, error) {
	return xorms.NewMysql(dsn, options...)
}
