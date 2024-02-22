package win

import (
	"fmt"
	"github.com/injoyai/goutil/g"
	"golang.org/x/sys/windows/registry"
	"strings"
)

// APPPath 从注册表获取软件路径 golang.org/x/sys v0.0.0-20220412211240-33da011f77ad
func APPPath(appName string) (result []string, err error) {

	fn := func(root registry.Key) error {
		path := "Software\\Microsoft\\Windows\\CurrentVersion\\App Paths\\"
		res, err := GetRegister(root, path)
		if err != nil {
			return err
		}
		for _, v := range res.Dir {
			if strings.Contains(v, appName) {
				appKey, err := registry.OpenKey(root, path+"\\"+v, registry.READ)
				if err != nil {
					return err
				}
				filename, _, err := appKey.GetStringValue("")
				appKey.Close()
				if err != nil {
					return err
				}
				result = append(result, filename)
			}
		}
		return nil
	}

	err = g.StopWithErr(
		func() error { return fn(registry.LOCAL_MACHINE) },
		func() error { return fn(registry.CURRENT_USER) },
	)

	return
}

const (
	REGISTER_ROOT             = REGISTER_CLASSES_ROOT
	REGISTER_USER             = REGISTER_CURRENT_USER
	REGISTER_LOCAL            = REGISTER_LOCAL_MACHINE
	REGISTER_CLASSES_ROOT     = registry.CLASSES_ROOT
	REGISTER_CURRENT_USER     = registry.CURRENT_USER
	REGISTER_LOCAL_MACHINE    = registry.LOCAL_MACHINE
	REGISTER_USERS            = registry.USERS
	REGISTER_CURRENT_CONFIG   = registry.CURRENT_CONFIG
	REGISTER_PERFORMANCE_DATA = registry.PERFORMANCE_DATA
)

// RegisterURLProtocol 注册URL协议
// @protocol 协议前缀,例如 download,浏览器输入 download://xxx 时会关联filename对应的程序
// @filename 协议对应的执行文件
func RegisterURLProtocol(root registry.Key, protocol string, filename string) error {
	if err := SetRegister(
		root,
		protocol,
		Value{Name: "URL Protocol", Value: ``},
	); err != nil {
		return err
	}
	if err := SetRegister(
		root,
		protocol+"/DefaultIcon",
	); err != nil {
		return err
	}
	return SetRegister(
		root,
		protocol+"/shell/open/command",
		Value{Name: "", Value: fmt.Sprintf(`"%s" "%%1"`, filename)},
	)
}

// GetRegister 获取注册表信息
func GetRegister(root registry.Key, path string) (*Register, error) {
	path = strings.ReplaceAll(path, "/", "\\")
	key, err := registry.OpenKey(root, path, registry.READ)
	if err != nil {
		return nil, err
	}
	defer key.Close()
	kList, err := key.ReadValueNames(0)
	if err != nil {
		return nil, err
	}
	dirs, err := key.ReadSubKeyNames(0)
	if err != nil {
		return nil, err
	}
	m := []Value(nil)
	for _, k := range kList {
		v, t, err := key.GetStringValue(k)
		if err != nil {
			return nil, err
		}
		m = append(m, Value{
			Name:  k,
			Type:  t,
			Value: v,
		})
	}
	return &Register{
		Dir:   dirs,
		Value: m,
	}, nil
}

// SetRegister 注册到注册表,文件,可选键值对
func SetRegister(root registry.Key, path string, values ...Value) error {
	path = strings.ReplaceAll(path, "/", "\\")
	key, _, err := registry.CreateKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	for _, v := range values {
		err = key.SetStringValue(v.Name, v.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// DelRegister 删除注册表的文件或者键值对
func DelRegister(root registry.Key, path string, names ...string) error {
	path = strings.ReplaceAll(path, "/", "\\")
	if len(names) == 0 {
		return registry.DeleteKey(root, path)
	}
	key, err := registry.OpenKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	for _, name := range names {
		if err := key.DeleteValue(name); err != nil {
			return err
		}
	}
	return nil
}

type Register struct {
	Dir   []string `json:"dir"`   //下级目录
	Value []Value  `json:"value"` //值
}

func (this *Register) String() string {
	s := "\n-dir:\n"
	for _, v := range this.Dir {
		s += fmt.Sprintf("\t%s\n", v)
	}
	s += "-value:\n"
	for _, v := range this.Value {
		if len(v.Name) == 0 {
			v.Name = "(默认)"
		}
		s += fmt.Sprintf("\t%s: %s\n", v.Name, v.Value)
	}
	return s
}

type Value struct {
	Name  string `json:"name"`
	Type  uint32 `json:"type"`
	Value string `json:"value"`
}
