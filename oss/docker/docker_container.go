package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
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
func (c Client) ContainerCreate(req *ContainerCreateReq) (container.CreateResponse, error) {
	config := &container.Config{
		Image:           req.Image,
		Cmd:             req.Cmd,
		Entrypoint:      req.Entrypoint,
		Env:             req.Env,
		Labels:          req.Labels,
		ExposedPorts:    req.GetPortSet(),
		OpenStdin:       req.OpenStdin,
		Tty:             req.TTY,
		NetworkDisabled: req.NetworkDisabled,
		Volumes:         req.GetVolumesMap(),
	}
	hostConf := &container.HostConfig{
		Binds:           req.GetVolumes(),
		LogConfig:       req.LogConfig.Get(),
		PortBindings:    req.GetPortMap(),
		RestartPolicy:   req.RestartPolicy.Get(),
		AutoRemove:      req.AutoRemove,
		Privileged:      req.Privileged,
		PublishAllPorts: req.PublishAllPorts,
		Resources: container.Resources{
			Memory:    int64(req.MemoryShares * 1024 * 1024),
			NanoCPUs:  int64(req.NanoCPUs * 1e9),
			CPUShares: int64(req.CPUShares),
		},
	}

	//网络配置
	networkConf := &network.NetworkingConfig{}
	switch req.Network {
	case "host", "none", "bridge":
		hostConf.NetworkMode = container.NetworkMode(req.Network)
		networkConf.EndpointsConfig = map[string]*network.EndpointSettings{req.Network: {}}
	case "":
	default:
		//自定义网络
		networkConf.EndpointsConfig = map[string]*network.EndpointSettings{req.Network: {}}
	}

	//创建容器
	resp, err := c.Client.ContainerCreate(c.ctx, config, hostConf, networkConf, &v1.Platform{}, req.Name)
	if err != nil {
		c.Client.ContainerRemove(c.ctx, resp.ID, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
		return container.CreateResponse{}, err
	}

	return resp, nil
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
	return c.Client.ContainerRemove(c.ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true, RemoveLinks: true, Force: true})
}

func (c Client) ContainerPause(id string) error {
	return c.Client.ContainerPause(c.ctx, id)
}

func (c Client) ContainerUnpause(id string) error {
	return c.Client.ContainerUnpause(c.ctx, id)
}

func (c Client) ContainerInspect(id string) (types.ContainerJSON, error) {
	return c.Client.ContainerInspect(c.ctx, id)
}
