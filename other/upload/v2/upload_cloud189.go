package upload

import (
	"github.com/tickstep/cloudpan189-api/cloudpan"
	"io"
	"path/filepath"
)

var _ Uploader = (*Cloud189)(nil)

func DialCloud189(username, password string) (*Cloud189, error) {
	appToken, err := cloudpan.AppLogin(username, password)
	if err != nil {
		return nil, err
	}
	webToken := &cloudpan.WebLoginToken{}
	webTokenStr := cloudpan.RefreshCookieToken(appToken.SessionKey)
	if webTokenStr != "" {
		webToken.CookieLoginUser = webTokenStr
	}
	panClient := cloudpan.NewPanClient(*webToken, *appToken)
	return &Cloud189{PanClient: panClient}, nil
}

type Cloud189 struct {
	*cloudpan.PanClient
}

func (this *Cloud189) Upload(filename string, reader io.Reader) (URL, error) {
	//this.PanClient.AppUploadFileCommitOverwrite()
	//TODO implement me
	panic("implement me")
}

func (this *Cloud189) List(join ...string) ([]*Info, error) {
	info, err := this.PanClient.FileInfoByPath(filepath.Join(join...))
	if err != nil {
		return nil, err
	}
	result, err := this.PanClient.FileList(&cloudpan.FileListParam{FileId: info.FileId})
	if err != nil {
		return nil, err
	}
	ls := []*Info(nil)
	for _, v := range result.Data {
		ls = append(ls, &Info{
			Name: v.FileName,
			Size: v.FileSize,
			Dir:  v.IsFolder,
			Time: 0, //v.CreateTime,
		})
	}
	return ls, nil
}
