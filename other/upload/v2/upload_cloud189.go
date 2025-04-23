package upload

import (
	"bytes"
	"github.com/injoyai/base/bytes/crypt/md5"
	"github.com/tickstep/cloudpan189-api/cloudpan"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var _ Uploader = (*Cloud189)(nil)

func NewCloud189(username, password string) (*Cloud189, error) {
	appToken, err := cloudpan.AppLogin(username, password)
	if err != nil {
		return nil, err
	}

	//webTokenStr := cloudpan.RefreshCookieToken(appToken.SessionKey)
	//webToken.CookieLoginUser = webTokenStr

	return &Cloud189{
		Token:  *appToken,
		Client: cloudpan.NewPanClient(cloudpan.WebLoginToken{}, *appToken),
	}, nil
}

type Cloud189 struct {
	Token  cloudpan.AppLoginToken //token
	Client *cloudpan.PanClient    //客户端
	Limit  int64                  //分片大小,待实现
}

func (this *Cloud189) Upload(filename string, r io.Reader) (URL, error) {
	dir, name := filepath.Split(filename)
	dirInfo, err := this.Client.AppFileInfoByPath(0, dir)
	if err != nil {
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

	//上传数据
	err = this.Client.AppUploadFileData(createRes.FileUploadUrl, createRes.UploadFileId, createRes.XRequestId, &cloudpan.AppFileUploadRange{}, func(method, url string, headers map[string]string) (resp *http.Response, err error) {
		req, err := http.NewRequest(method, url, bytes.NewReader(bs))
		if err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return http.DefaultClient.Do(req)
	})
	if err != nil {
		return nil, err
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

func (this *Cloud189) login() {

}

type Url string

func (this Url) String() string {
	return string(this)
}

func (this Url) Download(filename string) error {
	return nil
}
