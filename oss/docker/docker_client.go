package docker

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/goutil/database/sqlite"
	"github.com/injoyai/goutil/database/xorms"
	"github.com/injoyai/goutil/g"
	"xorm.io/xorm"
)

type Client struct {
	*client.Client                 //
	ctx            context.Context //
	configPath     string          //docker 配置目录
	storeCache     *maps.Safe      //仓库数据缓存
	DB             *xorm.Engine    //数据库,仓库需要数据库连接
}

func NewClient(ctx ...context.Context) (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return Client{
		ctx:        g.Ctx(ctx...),
		Client:     cli,
		configPath: "/etc/docker/daemon.json",
		storeCache: maps.NewSafe(),
		DB: sqlite.NewXorm(&xorms.Option{
			DSN:         "./data/docker/store.db",
			FieldSync:   true,
			TablePrefix: "docker_",
		}).Engine,
	}, err
}
