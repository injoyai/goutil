package docker

import (
	"errors"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/oss/linux/systemctl"
	"github.com/injoyai/goutil/oss/shell"
	json "github.com/json-iterator/go"
	"os"
	"strings"
	"time"
)

var (
	defaultStore = &Store{
		ID:       0,
		Name:     "DockerHub",
		Username: "",
		Password: "",
		Domain:   "https://docker.io",
		Created:  0,
	}
)

type Store struct {
	ID         int64  `jso:"id"`                      //主键
	Name       string `jso:"name"`                    //名称
	Username   string `json:"username"`               //账号
	Password   string `json:"password"`               //密码
	Domain     string `json:"domain"`                 //域名,下载地址前缀
	Created    int64  `json:"created" xorm:"created"` //创建时间
	Status     string `json:"status"`                 //状态信息,success成功 fail失败
	StatusText string `json:"statusText"`             //状态详情,错误信息
}

type StoreSearch struct {
	PageSize
	Name string `json:"name"`
}

// StoreList 仓库列表
func (c Client) StoreList(req *StoreSearch) ([]*Store, int64, error) {
	data := []*Store{}
	session := c.DB.Desc("ID")
	if len(req.Name) > 0 {
		session.Where("Name like ?", "%"+req.Name+"%")
	}
	if req.PageSize.PageSize > 0 {
		session.Limit(req.PageSize.PageSize, req.PageSize.PageNum*req.PageSize.PageSize)
	}
	co, err := session.FindAndCount(&data)
	if err != nil {
		return nil, 0, err
	}
	//增加默认仓库hub.docker.com
	co += 1
	if len(data) < req.PageSize.PageSize || req.PageSize.PageSize <= 0 {
		data = append(data, defaultStore)
	}
	return data, co, nil
}

// StoreGet 获取仓库详情
func (c Client) StoreGet(id int64) (*Store, error) {
	data := new(Store)
	has, err := c.DB.Where("ID=?", id).Get(data)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("仓库不存在")
	}
	return data, nil
}

type StoreCreateReq struct {
	Name     string `jso:"name"`      //名称
	Username string `json:"username"` //账号
	Password string `json:"password"` //密码
	Domain   string `json:"domain"`   //域名,下载地址前缀
}

// check 校验参数
// 校验仓库的账号密码
func (this *StoreCreateReq) check() error {
	result, err := shell.Exec("docker", "login", "-u", this.Username, "-p", this.Password, this.Domain)
	if err != nil {
		return err
	}
	if strings.Contains(result.String(), "Login Succeeded") {
		return nil
	}
	return errors.New(result.String())
}

// StoreCreate 添加新仓库
func (c Client) StoreCreate(req *StoreCreateReq) error {
	//校验账号密码是否正确
	connErr := req.check()
	//读取配置文件,增加第三方仓库的地址
	if err := c.storeConfigDeal("create", "", req.Domain); err != nil {
		return err
	}
	data := &Store{
		Name:       req.Name,
		Username:   req.Username,
		Password:   req.Password,
		Domain:     req.Domain,
		Created:    time.Now().Unix(),
		Status:     conv.Select[string](connErr == nil, "success", "fail"),
		StatusText: conv.New(connErr).String(""),
	}
	//添加到数据库
	_, err := c.DB.Insert(data)
	if err == nil {
		c.storeCache.Set(data.ID, data)
	}
	return err
}

type StoreUpdateReq struct {
	ID int64 `json:"id"` //主键
	*StoreCreateReq
}

// StoreUpdate 更新仓库信息
func (c Client) StoreUpdate(req *StoreUpdateReq) error {
	//获取仓库信息
	data, err := c.StoreGet(req.ID)
	if err != nil {
		return err
	}
	//校验账号密码是否正确
	connErr := req.check()
	//读取配置文件,增加第三方仓库的地址
	if err := c.storeConfigDeal("update", data.Domain, req.Domain); err != nil {
		return err
	}
	{ //更新字段
		data.Name = req.Name
		data.Username = req.Username
		data.Password = req.Password
		data.Domain = req.Domain
		data.Status = conv.Select[string](connErr == nil, "success", "fail")
		data.StatusText = conv.New(connErr).String("")
	}
	//同步到数据库
	_, err = c.DB.Where("ID=?", data.ID).AllCols().Update(data)
	if err == nil {
		c.storeCache.Set(data.ID, data)
	}
	return err
}

// StoreDelete 删除仓库
func (c Client) StoreDelete(id int64) error {

	//获取仓库信息
	data, err := c.StoreGet(id)
	if err != nil {
		return err
	}

	//读取配置文件,增加第三方仓库的地址
	if err := c.storeConfigDeal("update", data.Domain, ""); err != nil {
		return err
	}
	_, err = c.DB.Where("ID=?", id).Delete(new(Store))
	if err == nil {
		c.storeCache.Del(id)
	}
	return err
}

// storeConfigDeal 更新配置文件信息
func (c Client) storeConfigDeal(Type string, oldDomain, newDomain string) error {
	file, err := os.ReadFile(c.configPath)
	if err != nil {
		return err
	}
	cfgMap := map[string]interface{}{}
	if err := json.Unmarshal(file, &cfgMap); err != nil {
		return err
	}
	registries := conv.Interfaces(cfgMap["insecure-registries"])
	switch Type {
	case "create":
		registries = append(registries, newDomain)
	case "update":
		registries = append(registries, newDomain)
		for i, v := range registries {
			if conv.String(v) == oldDomain {
				registries = append(registries[:i], registries[i+1:]...)
				break
			}
		}
	case "delete":
		for i, v := range registries {
			if conv.String(v) == oldDomain {
				registries = append(registries[:i], registries[i+1:]...)
				break
			}
		}
	}
	cfgMap["insecure-registries"] = registries
	if err := oss.New(c.configPath, cfgMap); err != nil {
		return err
	}
	return c.dockerRestart(time.Second, 20)
}

// dockerRestart 重启docker服务
func (c Client) dockerRestart(interval time.Duration, num int) error {
	err := systemctl.Restart("docker")
	if err != nil {
		return err
	}
	var active bool
	for i := 0; i < num; i++ {
		active, err = systemctl.IsActive("docker")
		if err != nil {
			return err
		}
		if active {
			break
		}
		time.Sleep(interval)
	}
	if !active {
		return errors.New("docker服务重启失败")
	}
	return nil
}
