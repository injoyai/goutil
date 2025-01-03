package upload

import (
	"context"
	"github.com/injoyai/goutil/oss"
	"github.com/jlaffaye/ftp"
	"io"
	"path/filepath"
	"strings"
	"time"
)

var _ Uploader = (*FTP)(nil)

func DialFTP(address, username, password string) (*FTP, error) {

	c, err := ftp.Dial(address,
		ftp.DialWithContext(context.Background()),
		ftp.DialWithTimeout(time.Second*10),
	)
	if err != nil {
		return nil, err
	}

	if err := c.Login(username, password); err != nil {
		return nil, err
	}

	return &FTP{
		c:        c,
		address:  address,
		username: username,
		password: password,
	}, nil
}

type FTP struct {
	c        *ftp.ServerConn
	address  string
	username string
	password string
}

func (this *FTP) Upload(filename string, reader io.Reader) (URL, error) {
	err := this.c.Stor(filename, reader)
	return FTPUrl{
		address:  this.address,
		username: this.username,
		password: this.password,
		Filename: filename,
	}, err
}

func (this *FTP) List(join ...string) ([]*Info, error) {
	dir := filepath.Join(join...)
	dir, _ = strings.CutPrefix(dir, "/")
	dir, _ = strings.CutPrefix(dir, "\\")
	ls, err := this.c.List(dir)
	if err != nil {
		return nil, err
	}
	result := []*Info(nil)
	for _, v := range ls {
		result = append(result, &Info{
			Name: v.Name,
			Size: int64(v.Size),
			Dir:  v.Type == ftp.EntryTypeFolder,
			Time: v.Time.Unix(),
		})
	}
	return result, nil
}

type FTPUrl struct {
	address  string
	username string
	password string
	Filename string
}

func (this FTPUrl) String() string {
	return this.Filename
}

func (this FTPUrl) Download(filename string) error {
	c, err := ftp.Dial(this.address,
		ftp.DialWithContext(context.Background()),
		ftp.DialWithTimeout(time.Second*10),
	)
	if err != nil {
		return err
	}
	if err := c.Login(this.username, this.password); err != nil {
		return err
	}
	resp, err := c.Retr(this.Filename)
	if err != nil {
		return err
	}
	defer resp.Close()
	return oss.New(filename, resp)
}
