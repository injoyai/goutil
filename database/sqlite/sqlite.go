package sqlite

import "github.com/injoyai/goutil/database/xorms"

func NewXorm(filename string, options ...xorms.Option) (*xorms.Engine, error) {
	return xorms.NewSqlite(filename, options...)
}
