package docker

import (
	"errors"
	"github.com/docker/docker/api/types/volume"
	"strings"
)

type VolumeSearch struct {
	PageSize
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

// VolumeList 存储卷列表
func (c Client) VolumeList(req *VolumeSearch) (result []*volume.Volume, co int64, err error) {
	resp, err := c.Client.VolumeList(c.ctx, volume.ListOptions{})
	if err != nil {
		return nil, 0, err
	}
	co = int64(len(resp.Volumes))

	//数据筛选
	for _, v := range resp.Volumes {
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

type VolumeCreateReq struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	DriverOpts map[string]string `json:"driverOpts"`
	Labels     map[string]string `json:"labels"`
}

// VolumeCreate 存储卷新建
func (c Client) VolumeCreate(req *VolumeCreateReq) error {
	list, _, err := c.VolumeList(&VolumeSearch{Name: req.Name})
	if err != nil {
		return err
	}
	for _, v := range list {
		if v.Name == req.Name {
			return errors.New("存储卷已存在")
		}
	}
	_, err = c.Client.VolumeCreate(c.ctx, volume.CreateOptions{
		Driver:     req.Driver,
		DriverOpts: req.DriverOpts,
		Labels:     req.Labels,
		Name:       req.Name,
	})
	return err
}

// VolumeDelete 存储卷删除
func (c Client) VolumeDelete(id string) error {
	err := c.Client.VolumeRemove(c.ctx, id, true)
	if err != nil {
		if strings.Contains(err.Error(), "volume is in use") {
			return errors.New("储存卷正在使用")
		}
	}
	return err
}
