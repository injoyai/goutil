package license

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/injoyai/base/crypt/aes"
	"github.com/injoyai/goutil/oss"
)

func WithFilename(name string) Option {
	return func(l *License) {
		l.Filename = name
	}
}

func WithEncode(fn func(bs []byte) ([]byte, error)) Option {
	return func(l *License) {
		l.Encode = fn
	}
}

func WithDecode(fn func(bs []byte) ([]byte, error)) Option {
	return func(l *License) {
		l.Decode = fn
	}
}

func New(appName string, op ...Option) *License {
	key := "License.App.Key."
	l := &License{
		AppName:  appName,
		Filename: oss.UserInjoyDir(appName, "License.txt"),
		Encode: func(bs []byte) ([]byte, error) {
			return aes.EncryptCBC(bs, []byte(key))
		},
		Decode: func(bs []byte) ([]byte, error) {
			return aes.DecryptCBC(bs, []byte(key))
		},
	}
	for _, v := range op {
		v(l)
	}
	return l
}

type Option func(*License)

type License struct {
	AppName  string
	Filename string
	Encode   func(bs []byte) ([]byte, error)
	Decode   func(bs []byte) ([]byte, error)
}

func (this *License) License(c Code) (string, error) {
	return this.encode(c)
}

func (this *License) Activate(code string) error {

	//解密code
	c := new(Code)
	err := this.decode(code, c)
	if err != nil {
		return err
	}

	now := time.Now()
	if now.Unix() > c.End {
		return errors.New("激活码已失效")
	}

	//读取唯一标识
	id, err := machineid.ProtectedID(this.AppName)
	if err != nil {
		return err
	}

	//生成激活信息
	info := &Info{
		MachineID:  id,
		ExpireTime: c.Expire,
		ActivateAt: now.Unix(),
	}

	//加密激活信息
	s, err := this.encode(info)
	if err != nil {
		return err
	}

	//存储激活信息
	return oss.New(this.Filename, s)
}

func (this *License) Valid() (*Info, bool, error) {
	//读取激活信息
	info, err := this.loadingInfo()
	if err != nil {
		return nil, false, err
	}
	//读取唯一标识
	id, err := machineid.ProtectedID(this.AppName)
	if err != nil {
		return nil, false, err
	}
	//比对机器码
	if info.MachineID != id {
		return info, false, nil
	}
	//判断有效期
	if info.Expired() {
		return info, false, nil
	}
	return info, true, nil
}

func (this *License) loadingInfo() (*Info, error) {
	//读取激活信息
	bs, err := oss.Read(this.Filename)
	if err != nil {
		return nil, err
	}
	i := new(Info)
	//解密激活信息
	err = this.decode(string(bs), i)
	return i, err
}

func (this *License) encode(a any) (string, error) {
	bs, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	bs, err = this.Encode(bs)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bs), nil
}

func (this *License) decode(s string, ptr any) error {
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	bs, err = this.Decode(bs)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, ptr)
}

type Code struct {
	End    int64
	Expire int64
}

type Info struct {
	MachineID  string
	ExpireTime int64
	ActivateAt int64
}

func (this *Info) Expired() bool {
	return time.Now().Unix() > this.ExpireTime
}
