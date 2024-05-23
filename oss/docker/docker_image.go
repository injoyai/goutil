package docker

import (
	"encoding/base64"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/injoyai/goutil/oss"
	json "github.com/json-iterator/go"
	"io"
	"os"
	"strings"
)

type ImageSearch struct {
	PageSize
	ID  string `json:"id"`
	Tag string `json:"tag"`
}

type ImageInfo struct {
	ID       string   `json:"id"`       //
	Tags     []string `json:"tags"`     //tag
	Size     int64    `json:"size"`     //字节数
	SizeText string   `json:"sizeText"` //可视的字节
	CreateAt int64    `json:"createAt"` //创建时间
}

func newImageInfo(info types.ImageSummary) *ImageInfo {
	return &ImageInfo{
		ID:       info.ID,
		Tags:     info.RepoTags,
		Size:     info.Size,
		SizeText: oss.SizeString(info.Size),
		CreateAt: info.Created,
	}
}

func (c Client) ImageList(req *ImageSearch) (result []*ImageInfo, co int64, err error) {
	list, err := c.Client.ImageList(c.ctx, types.ImageListOptions{})
	if err != nil {
		return nil, 0, err
	}
	co = int64(len(list))

	//数据筛选
	for _, v := range list {
		if len(v.ID) > 0 && !strings.Contains(v.ID, req.ID) {
			continue
		}
		if len(req.Tag) > 0 && !func() bool {
			for _, tag := range v.RepoTags {
				if strings.Contains(tag, req.Tag) {
					return true
				}
			}
			return false
		}() {
			continue
		}
		result = append(result, newImageInfo(v))
	}

	//分页
	start, end := req.PageSize.Limit(len(result))
	result = result[start:end]

	return result, co, nil
}

type ImagePullReq struct {
	Name     string `json:"name"`     //镜像名称
	Domain   string `json:"domain"`   //域名,默认空 docker.io
	Username string `json:"username"` //账号,默认空
	Password string `json:"password"` //密码,默认空
}

// ImagePull 拉取镜像
func (c Client) ImagePull(req *ImagePullReq) (io.ReadCloser, error) {
	refStr := req.Name
	options := types.ImagePullOptions{}
	authConfig := registry.AuthConfig{
		Username: req.Username,
		Password: req.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	options.RegistryAuth = authStr
	if req.Domain != "" {
		refStr = req.Domain + "/" + req.Name
	}
	return c.Client.ImagePull(c.ctx, refStr, options)
}

type ImagePushReq struct {
	Name     string `json:"name"`     //镜像名称
	Domain   string `json:"domain"`   //域名,默认空 docker.io
	Username string `json:"username"` //账号,默认空
	Password string `json:"password"` //密码,默认空
}

// ImagePush 推送镜像到仓库
func (c Client) ImagePush(req *ImagePushReq) (io.ReadCloser, error) {
	options := types.ImagePushOptions{}
	authConfig := registry.AuthConfig{
		Username: req.Username,
		Password: req.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	options.RegistryAuth = authStr
	return c.Client.ImagePush(c.ctx, req.Name, options)
}

// ImageImport 导入镜像
func (c Client) ImageImport() error {
	return nil
}

type ImageExportReq struct {
	Tag      string `json:"tag"`
	Filename string `json:"filename"` //文件名称 例 /home/test.tar
}

// ImageExport 导出镜像,生成tar
func (c Client) ImageExport(req *ImageExportReq) error {
	out, err := c.ImageSave(c.ctx, []string{req.Tag})
	if err != nil {
		return err
	}
	defer out.Close()
	file, err := os.OpenFile(req.Filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, out)
	return err
}

type ImageBuildReq struct {
	Dockerfile string   `json:"dockerfile"`
	Tag        []string `json:"tag"`
}

// ImageBuild 构建镜像,根据dockerfile
func (c Client) ImageBuild(req *ImageBuildReq) error {
	opts := types.ImageBuildOptions{
		Dockerfile: req.Dockerfile,
		Tags:       req.Tag,
		Remove:     true,
		Labels:     map[string]string{},
	}
	_, err := c.Client.ImageBuild(c.ctx, strings.NewReader(req.Dockerfile), opts)
	return err
}

type ImageTagReq struct {
	SourceID   string `json:"id"` //镜像id
	TargetName string `json:"targetName"`
}

func (c Client) ImageTag(req *ImageTagReq) error {
	err := c.Client.ImageTag(c.ctx, req.SourceID, req.TargetName)
	return err
}

// ImageClear 清理镜像
// 清理未使用的镜像,清理无效的镜像(无标签,大小为0)
func (c Client) ImageClear() error {
	return nil
}

// ImageDelete 删除镜像,根据镜像id
func (c Client) ImageDelete(id string) error {
	_, err := c.ImageRemove(c.ctx, id, types.ImageRemoveOptions{Force: true})
	return err
}

// ImageExist 校验镜像是否存在
func (c Client) ImageExist(id string) (bool, error) {
	_, co, err := c.ImageList(&ImageSearch{ID: id})
	return co > 0, err
}
