package upload

import (
	"github.com/injoyai/goutil/oss/fss"
	"io"
	"os"
	"path/filepath"
)

var _ Uploader = (*Smb)(nil)

type SmbConfig = fss.SmbConfig

func NewSmb(cfg *SmbConfig) (*Smb, error) {

	root, err := fss.NewSmb(cfg)
	return &Smb{Smb: root}, err
}

type Smb struct {
	*fss.Smb
}

func (this *Smb) Upload(filename string, reader io.Reader) (string, error) {
	f, err := this.Share.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return filename, err
}

func (this *Smb) Download(filename, localFilename string) error {
	f, err := this.Share.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	f2, err := os.Create(localFilename)
	if err != nil {
		return err
	}
	defer f2.Close()
	_, err = io.Copy(f2, f)
	return err
}

func (this *Smb) Dir(join ...string) ([]*Info, error) {
	infos, err := this.Share.ReadDir(filepath.Join(join...))
	if err != nil {
		return nil, err
	}
	ls := make([]*Info, len(infos))
	for i := range infos {
		ls[i] = &Info{
			Name: infos[i].Name(),
			Size: infos[i].Size(),
			Dir:  infos[i].IsDir(),
			Time: infos[i].ModTime().Unix(),
		}
	}
	return ls, nil
}
