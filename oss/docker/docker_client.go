package docker

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/goutil/database/sqlite"
	"github.com/injoyai/goutil/database/xorms"
	"xorm.io/xorm"
)

type Client struct {
	*client.Client                 //
	ctx            context.Context //
	configPath     string          //docker 配置目录
	storeCache     *maps.Safe      //仓库数据缓存
	DB             *xorm.Engine    //数据库,仓库需要数据库连接
}

func NewClient() (Client, error) {
	return NewClientContext(context.Background())
}

func NewClientContext(ctx context.Context) (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return Client{}, err
	}
	e, err := sqlite.NewXorm("./data/docker/store.db", xorms.WithTablePrefix("docker_"))
	if err != nil {
		return Client{}, err
	}
	return Client{
		ctx:        ctx,
		Client:     cli,
		configPath: "/etc/docker/daemon.json",
		storeCache: maps.NewSafe(),
		DB:         e.Engine,
	}, nil
}
