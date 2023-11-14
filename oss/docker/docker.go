package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func NewClient() (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return Client{}, err
	}

	return Client{
		cli: cli,
	}, nil
}

// GetContainerList 获取容器列表
func (c Client) GetContainerList() ([]types.Container, error) {
	var options types.ContainerListOptions
	containers, err := c.cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (c Client) GetContainerListByName(names []string) ([]types.Container, error) {
	var (
		options types.ContainerListOptions
		res     []types.Container
	)
	options.All = true
	if len(names) > 0 {
		var array []filters.KeyValuePair
		for _, n := range names {
			array = append(array, filters.Arg("name", n))
		}
		options.Filters = filters.NewArgs(array...)
	}
	containers, err := c.cli.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		if container.Names[0] == "/"+names[0] {
			res = append(res, container)
		}
	}
	return res, nil
}

// CreateNetwork 创建网络
func (c Client) CreateNetwork(name string) error {
	_, err := c.cli.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Driver: "bridge",
	})
	return err
}

// NetworkExist 判断网络是否存在
func (c Client) NetworkExist(name string) bool {
	var options types.NetworkListOptions
	options.Filters = filters.NewArgs(filters.Arg("name", name))
	networks, err := c.cli.NetworkList(context.Background(), options)
	if err != nil {
		return false
	}
	return len(networks) > 0
}

// PullImage 拉取镜像
func (c Client) PullImage(imageName string, force bool) error {
	if !force {
		exist, err := c.ImageExist(imageName)
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
	}
	if _, err := c.cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{}); err != nil {
		return err
	}
	return nil
}

// DeleteImage 删除镜像
func (c Client) DeleteImage(imageID string) error {
	if _, err := c.cli.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: true}); err != nil {
		return err
	}
	return nil
}

// GetImageIDByName 根据名称获取镜像
func (c Client) GetImageIDByName(imageName string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("reference", imageName)
	list, err := c.cli.ImageList(context.Background(), types.ImageListOptions{
		Filters: filter,
	})
	if err != nil {
		return "", err
	}
	if len(list) > 0 {
		return list[0].ID, nil
	}
	return "", nil
}

// ImageExist 校验镜像是否存在
func (c Client) ImageExist(imageName string) (bool, error) {
	filter := filters.NewArgs()
	filter.Add("reference", imageName)
	list, err := c.cli.ImageList(context.Background(), types.ImageListOptions{
		Filters: filter,
	})
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}
