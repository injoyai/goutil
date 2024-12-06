package upload

import (
	"context"
	"errors"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"io"
	"time"
)

var _ Uploader = (*Qiniu)(nil)

func NewQiniu(cfg *QiniuConfig) (*Qiniu, error) {
	return &Qiniu{cfg}, nil
}

type Qiniu struct {
	*QiniuConfig
}

type QiniuConfig struct {
	Key    string //key
	Secret string //secret
	Domain string //前缀
	Space  string //空间
}

// GetPrivateURL 七牛云，获取私有下载URL,有效期1小时
func (this *Qiniu) GetPrivateURL(key string) string {
	mac := qbox.NewMac(this.Key, this.Secret)
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	return storage.MakePrivateURL(mac, this.Domain, key, deadline)
}

// GetPrivateToken 七牛云，获取私有上传Token
func (this *Qiniu) GetPrivateToken() string {
	putPolicy := storage.PutPolicy{
		Scope: this.Space,
	}
	mac := qbox.NewMac(this.Key, this.Secret)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

//===============================================

func (this *Qiniu) Upload(name string, reader io.Reader) (URL, error) {
	mac := qbox.NewMac(this.Key, this.Secret)
	putPolicy := storage.PutPolicy{Scope: this.Space}
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong, // 空间对应的机房
		UseHTTPS:      false,                // 是否使用https域名
		UseCdnDomains: false,                // 上传是否使用CDN上传加速
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.Put(context.Background(), &ret, upToken, name, reader, -1, new(storage.PutExtra))
	return HttpUrl(this.Domain + ret.Key), err
}

func (this *Qiniu) List(join ...string) ([]*Info, error) {
	return nil, errors.New("暂不支持")
}
