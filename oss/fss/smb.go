package fss

import (
	"github.com/hirochachacha/go-smb2"
	"github.com/injoyai/base/safe"
	"io"
	"io/fs"
	"net"
)

func NewSmb(cfg *SmbConfig) (*Smb, error) {

	conn, err := net.Dial("tcp", cfg.Host+":445")
	if err != nil {
		return nil, err
	}

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     cfg.Username,
			Password: cfg.Password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}

	fs, err := s.Mount(cfg.ShareName) // 共享名
	if err != nil {
		return nil, err
	}

	return &Smb{
		Share: fs,
		Closer: safe.NewCloser().SetCloseFunc(func(err error) error {
			fs.Umount()
			s.Logoff()
			return conn.Close()
		}),
	}, nil
}

type SmbConfig struct {
	Host      string
	Username  string
	Password  string
	ShareName string
}

type Smb struct {
	*smb2.Share
	io.Closer
}

func (this *Smb) Open(filename string) (fs.File, error) {
	return this.Share.Open(filename)
}

func (this *Smb) Create(filename string) (fs.File, error) {
	return this.Share.Create(filename)
}

func (this *Smb) Remove(filename string) error {
	return this.Share.RemoveAll(filename)
}

func (this *Smb) Stat(filename string) (fs.FileInfo, error) {
	return this.Share.Stat(filename)
}

func (this *Smb) Rename(oldFilename, newFilename string) error {
	return this.Share.Rename(oldFilename, newFilename)
}

func (this *Smb) ReadDir(dir string) ([]fs.FileInfo, error) {
	return this.Share.ReadDir(dir)
}

func (this *Smb) Mkdir(dir string, perm fs.FileMode) error {
	return this.Share.MkdirAll(dir, perm)
}
