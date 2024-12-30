package win

import (
	"golang.org/x/sys/windows/registry"
)

// GetCurrentEnv 获取当前用户环境变量
func GetCurrentEnv(key string) (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	val, _, err := k.GetStringValue(key)
	return val, err
}

// SetCurrentEnv 设置当前用户环境变量
func SetCurrentEnv(key, value string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue(key, value)
}

// GetRootEnv 获取系统环境变量
func GetRootEnv(key string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	val, _, err := k.GetStringValue(key)
	return val, err
}

// SetRootEnv 设置系统环境变量
func SetRootEnv(key, value string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue(key, value)
}
