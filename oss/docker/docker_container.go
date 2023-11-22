package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"strings"
)

type ContainerSearch struct {
	PageSize
	Name      string `json:"name"`      //容器名称
	ImageName string `json:"imageName"` //镜像名称
	ImageID   string `json:"imageID"`   //镜像id
	IP        string `json:"ip"`        //ip
	Port      uint16 `json:"port"`      //端口
	Status    string `json:"status"`    //容器状态
}

// ContainerList 获取容器列表
func (c Client) ContainerList(req *ContainerSearch) (result []types.Container, co int64, err error) {
	list, err := c.Client.ContainerList(c.ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, 0, err
	}
	co = int64(len(list))

	//数据筛选
	for _, v := range list {
		//筛选容器名称
		if len(req.Name) > 0 && func() bool {
			for _, name := range v.Names {
				if strings.Contains(name, req.Name) {
					return true
				}
			}
			return false
		}() {
			continue
		}
		// 筛选ip
		if len(req.IP) > 0 && func() bool {
			for _, port := range v.Ports {
				if strings.Contains(port.IP, req.IP) {
					return true
				}
			}
			return false
		}() {
			continue
		}
		// 筛选映射端口
		if req.Port > 0 && func() bool {
			for _, port := range v.Ports {
				if port.PrivatePort == req.Port || port.PublicPort == req.Port {
					return true
				}
			}
			return false
		}() {
			continue
		}
		//筛选镜像名称
		if len(req.ImageName) > 0 && !strings.Contains(v.Image, req.ImageName) {
			continue
		}
		//筛选镜像id
		if len(req.ImageID) > 0 && !strings.Contains(v.ImageID, req.ImageID) {
			continue
		}
		result = append(result, v)
	}

	//分页
	start, end := req.PageSize.Limit(len(result))
	result = result[start:end]

	return
}

// ContainerCreate 容器创建
func (c Client) ContainerCreate() error {
	return nil
}

func (c Client) ContainerStart(id string) error {
	return c.Client.ContainerStart(c.ctx, id, types.ContainerStartOptions{})
}

func (c Client) ContainerStop(id string) error {
	return c.Client.ContainerStop(c.ctx, id, container.StopOptions{})
}

func (c Client) ContainerKill(id string) error {
	return c.Client.ContainerKill(c.ctx, id, "SIGKILL")
}

func (c Client) ContainerRemove(id string) error {
	return c.Client.ContainerRemove(c.ctx, id, types.ContainerRemoveOptions{})
}

func (c Client) ContainerPause(id string) error {
	return c.Client.ContainerPause(c.ctx, id)
}

func (c Client) ContainerUnpause(id string) error {
	return c.Client.ContainerUnpause(c.ctx, id)
}

func (c Client) ContainerInspect(id string) (*types.ContainerJSON, error) {
	return nil, nil
}
