package xorms

import (
	"xorm.io/core"
	"xorm.io/xorm"
)

type Option struct {
	DSN         string //账号密码
	FieldSync   bool   //字段和表名同步
	TablePrefix string //表名前缀
}

type Config struct {
	Type        string //连接方式
	DSN         string //账号密码
	FieldSync   bool   //字段和表名同步
	TablePrefix string //表名前缀
}

var defaultConfig = &Config{
	Type:        "mysql",
	DSN:         "root:root@tcp(127.0.0.1:3306)/test",
	FieldSync:   true,
	TablePrefix: "",
}

func (this *Config) SetType(s string) *Config {
	this.Type = s
	return this
}

func (this *Config) SetDSN(s string) *Config {
	this.DSN = s
	return this
}

func (this *Config) SetFieldSync(b ...bool) *Config {
	this.FieldSync = !(len(b) > 0 && !b[0])
	return this
}

func (this *Config) SetTablePrefix(s string) *Config {
	this.TablePrefix = s
	return this
}

func (this *Config) Open() *Engine {
	db, err := xorm.NewEngine(this.Type, this.DSN)
	data := &Engine{
		Engine: db,
		cfg:    this,
		err:    err,
	}
	if db != nil {
		if this.FieldSync {
			db.SetMapper(core.SameMapper{}) //字段同步
		}
		db.SetTableMapper(core.NewPrefixMapper(core.SameMapper{}, this.TablePrefix))
	}
	return data
}
