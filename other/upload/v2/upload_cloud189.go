package upload

import (
	"bytes"
	"encoding/json"
	"github.com/injoyai/base/crypt/md5"
	"github.com/injoyai/goutil/oss"
	"github.com/tickstep/cloudpan189-api/cloudpan"
	"github.com/tickstep/cloudpan189-api/cloudpan/apierror"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var _ Uploader = (*Cloud189)(nil)

type Cloud189ConfigCache interface {
	Set(*cloudpan.AppLoginToken) error
	Get() (*cloudpan.AppLoginToken, error)
}

type Cloud189Config struct {
	Username string
	Password string
	Size     int
	Retry    int
	Cache    Cloud189ConfigCache
}

func NewCloud189(cfg Cloud189Config) (*Cloud189, error) {
	c := &Cloud189{Config: cfg}
	return c, c.refreshToken()
}

type Cloud189 struct {
	Token  cloudpan.AppLoginToken //token
	Client *cloudpan.PanClient    //客户端
	Config Cloud189Config
}

func (this *Cloud189) refreshToken() error {
	if this.Config.Cache == nil {
		appToken, err := cloudpan.AppLogin(this.Config.Username, this.Config.Password)
		if err != nil {
			return err
		}
		this.Client = cloudpan.NewPanClient(cloudpan.WebLoginToken{}, *appToken)
		return nil
	}

	appToken, err := this.Config.Cache.Get()
	if err != nil || appToken.SskAccessTokenExpiresIn < time.Now().UnixMilli() {
		token, err := cloudpan.AppLogin(this.Config.Username, this.Config.Password)
		if err != nil {
			return err
		}
		appToken = token
	}
	if err != nil {
		return err
	}
	this.Client = cloudpan.NewPanClient(cloudpan.WebLoginToken{}, *appToken)
	this.Config.Cache.Set(appToken)
	return nil
}

func (this *Cloud189) Rename(filename, newname string) error {
	fileInfo, err := this.Client.AppFileInfoByPath(0, filename)
	if err != nil {
		return err
	}
	_, err = this.Client.AppRenameFile(fileInfo.FileId, newname)
	return err
}

func (this *Cloud189) Upload(filename string, r io.Reader) (URL, error) {
	dir, name := filepath.Split(filename)
	dirInfo, err := this.Client.AppFileInfoByPath(0, dir)
	if err != nil {
		if err.Code == apierror.ApiCodeFileNotFoundCode {
			_, err = this.Client.AppMkdirRecursive(0, "", dir, 0, strings.Split(dir, "/"))
			if err != nil {
				return nil, err
			}
			return this.Upload(filename, r)
		}
		return nil, err
	}

	bs, er := io.ReadAll(r)
	if er != nil {
		return nil, er
	}

	//创建上传请求
	createRes, err := this.Client.AppCreateUploadFile(&cloudpan.AppCreateUploadFileParam{
		ParentFolderId: dirInfo.FileId,
		FileName:       name,
		Size:           int64(len(bs)),
		Md5:            md5.Encrypt(string(bs)),
	})
	if err != nil {
		return nil, err
	}

	size := this.Config.Size
	if size <= 0 {
		size = len(bs)
	}

	//上传数据
	for offset := 0; offset < len(bs); offset += this.Config.Size {
		if offset+this.Config.Size > len(bs) {
			size = len(bs) - offset
		}
		for n := 0; n <= this.Config.Retry; n++ {
			err = this.Client.AppUploadFileData(createRes.FileUploadUrl, createRes.UploadFileId, createRes.XRequestId,
				&cloudpan.AppFileUploadRange{
					Offset: int64(offset),
					Len:    int64(size),
				}, func(method, url string, headers map[string]string) (resp *http.Response, err error) {
					req, err := http.NewRequest(method, url, bytes.NewReader(bs))
					if err != nil {
						return nil, err
					}
					for k, v := range headers {
						req.Header.Set(k, v)
					}
					return http.DefaultClient.Do(req)
				})
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}
	}

	//提交数据
	commitRes, err := this.Client.AppUploadFileCommitOverwrite(createRes.FileCommitUrl, createRes.UploadFileId, createRes.XRequestId, true)
	if err != nil {
		return nil, err
	}
	_ = commitRes

	return Url(filename), nil
}

func (this *Cloud189) List(join ...string) ([]*Info, error) {
	dir := filepath.Join(join...)
	dir = strings.ReplaceAll(dir, "\\", "/")
	if len(dir) == 0 || dir[0] != '/' {
		dir = "/" + dir
	}
	fi, err := this.Client.AppFileInfoByPath(0, dir)
	if err != nil {
		return nil, err
	}
	res, err := this.Client.AppGetAllFileList(&cloudpan.AppFileListParam{FileId: fi.FileId})
	if err != nil {
		return nil, err
	}
	ls := []*Info(nil)
	for _, v := range res.FileList {
		t, err := time.Parse(time.DateTime, v.CreateTime)
		if err != nil {
			return nil, err
		}
		ls = append(ls, &Info{
			Name: v.FileName,
			Size: v.FileSize,
			Dir:  v.IsFolder,
			Time: t.Unix(),
		})
	}
	return ls, nil
}

type Url string

func (this Url) String() string {
	return string(this)
}

func (this Url) Download(filename string) error {
	return nil
}

func NewCloud189FileCache(filename string) Cloud189ConfigCache {
	oss.NewNotExist(filename, "{}")
	return &cloud189ConfigCache{filename: filename}
}

type cloud189ConfigCache struct {
	filename string
}

func (this *cloud189ConfigCache) Get() (*cloudpan.AppLoginToken, error) {
	cfg := new(cloudpan.AppLoginToken)
	bs, err := os.ReadFile(this.filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, cfg)
	return cfg, nil
}

func (this *cloud189ConfigCache) Set(token *cloudpan.AppLoginToken) error {
	bs, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return os.WriteFile(this.filename, bs, 0644)
}
