package docker

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/injoyai/conv"
)

type ContainerCreateReq struct {
	Name            string            `json:"name"`            //容器名称
	Image           string            `json:"image"`           //镜像
	Network         string            `json:"network"`         //网络
	Volumes         []Volume          `json:"volumes"`         //存储卷
	PublishAllPorts bool              `json:"publishAllPorts"` //映射全部端口
	ExposedPorts    []Port            `json:"exposedPorts"`    //映射端口
	ForcePull       bool              `json:"forcePull"`       //强制拉取
	Entrypoint      []string          `json:"entrypoint"`
	Cmd             []string          `json:"cmd"`
	Labels          map[string]string `json:"labels"`
	Env             []string          `json:"env"`
	OpenStdin       bool              `json:"openStdin"`
	TTY             bool              `json:"tty"`
	NetworkDisabled bool              `json:"networkDisabled"` //禁用容器网络

	Privileged    bool          `json:"privileged"`    //授予完全权限
	AutoRemove    bool          `json:"autoRemove"`    //退出容器后自动删除
	LogConfig     LogConfig     `json:"logConfig"`     //日志
	RestartPolicy RestartPolicy `json:"restartPolicy"` //重启方式 no(不重启) always(始终重启) unless-stopped(除非用户手动停止) on-failure(最大重试次数)
	CPUShares     int64         `json:"cpuShares"`     //CPU分配 单位核 0表示不限制
	NanoCPUs      float64       `json:"nanoCPUs"`      //CPU权重 单位纳秒 0表示不分配
	MemoryShares  float64       `json:"memoryShares"`  //内存分配 ,单位MB 0表示不限制
}

// GetPortSet 端口映射 conf
func (this *ContainerCreateReq) GetPortSet() nat.PortSet {
	portSet := nat.PortSet{}
	for _, v := range this.ExposedPorts {
		protocol := conv.SelectString(len(v.Protocol) > 0, v.Protocol, "tcp")
		port := v.HostPort + "/" + protocol
		portSet[nat.Port(port)] = struct{}{}
	}
	return portSet
}

// GetPortMap 端口映射 hostConf
func (this *ContainerCreateReq) GetPortMap() nat.PortMap {
	portMap := nat.PortMap{}
	for _, v := range this.ExposedPorts {
		protocol := conv.SelectString(len(v.Protocol) > 0, v.Protocol, "tcp")
		port := v.HostPort + "/" + protocol
		portMap[nat.Port(port)] = []nat.PortBinding{
			{HostPort: v.ContainerPort},
		}
	}
	return portMap
}

func (this *ContainerCreateReq) GetVolumesMap() map[string]struct{} {
	m := make(map[string]struct{})
	for _, v := range this.Volumes {
		m[v.ContainerDir] = struct{}{}
	}
	return m
}

func (this *ContainerCreateReq) GetVolumes() []string {
	ls := make([]string, 0)
	for _, v := range this.Volumes {
		ls = append(ls, fmt.Sprintf("%s:%s:%s", v.SourceDir, v.ContainerDir, v.Mode))
	}
	return ls
}

type Volume struct {
	SourceDir    string `json:"sourceDir"`
	ContainerDir string `json:"containerDir"`
	Mode         string `json:"mode"`
}

type Port struct {
	HostIP        string `json:"hostIP"`
	HostPort      string `json:"hostPort"`
	ContainerPort string `json:"containerPort"`
	Protocol      string `json:"protocol"` //tcp udp
}

type RestartPolicy struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (this RestartPolicy) Get() container.RestartPolicy {
	return container.RestartPolicy{
		Name:              this.Name,
		MaximumRetryCount: this.Count,
	}
}

type LogConfig struct {
	Type   string            `json:"type"` //awslogs fluentd gcplogs gelf journald json-file local logentries splunk syslog
	Config map[string]string `json:"config"`
}

func (this LogConfig) Get() container.LogConfig {
	return container.LogConfig{
		Type:   this.Type,
		Config: this.Config,
	}
}
