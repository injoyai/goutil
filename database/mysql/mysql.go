package mysql

import "github.com/injoyai/goutil/database/xorms"

func NewXorm(dsn string, options ...xorms.Option) (*xorms.Engine, error) {
	return xorms.NewMysql(dsn, options...)
}
