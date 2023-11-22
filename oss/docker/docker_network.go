package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"strings"
)

type NetworkSearch struct {
	PageSize
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

// NetworkList 获取网络列表,模拟分页
func (c Client) NetworkList(req *NetworkSearch) (result []types.NetworkResource, co int64, err error) {
	list, err := c.Client.NetworkList(c.ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, 0, err
	}
	co = int64(len(list))

	//数据筛选
	for _, v := range list {
		if len(req.Name) > 0 && !strings.Contains(v.Name, req.Name) {
			continue
		}
		if len(req.Driver) > 0 && v.Driver != req.Driver {
			continue
		}
		result = append(result, v)
	}

	//分页
	start, end := req.PageSize.Limit(len(result))
	result = result[start:end]

	return
}

type NetworkCreate struct {
	Name string `json:"name"`
	types.NetworkCreate
}

// NetworkCreate 创建网络
func (c Client) NetworkCreate(req *NetworkCreate) error {
	_, err := c.Client.NetworkCreate(c.ctx, req.Name, req.NetworkCreate)
	return err
}

// NetworkExist 判断网络是否存在
func (c Client) NetworkExist(name string) (bool, error) {
	options := types.NetworkListOptions{}
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	networks, err := c.Client.NetworkList(c.ctx, options)
	return len(networks) > 0, err
}

// NetworkDelete 删除网络
func (c Client) NetworkDelete(id string) error {
	return c.NetworkRemove(c.ctx, id)
}
