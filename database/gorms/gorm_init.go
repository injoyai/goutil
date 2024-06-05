package gorms

import (
	"errors"
	"github.com/jinzhu/gorm"
)

func New(Type, dsn string, options ...func(db *gorm.DB)) (*Engine, error) {
	return nil, errors.New("未实现")
}

type Engine struct{}

func WithMaxOpenConns(n int) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		db.DB().SetMaxOpenConns(n)
	}
}
