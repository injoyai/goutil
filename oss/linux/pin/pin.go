package pin

import (
	"errors"
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/oss/shell"
	"regexp"
	"strings"
)

type Model string

const (
	Input   Model = "in"  //输入
	Output  Model = "out" //输出
	High          = 1     //高电平
	Low           = 0     //低电平
	gpioDir       = "/sys/class/gpio/"
)

type Pin interface {
	GetNumber() int

	// GetModel 获取当前引脚模式 输入/输出
	GetModel() (Model, error)

	// SetModel 设置引脚模式 输入/输出
	SetModel(Model) error

	// GetValue 获取pin当前状态
	GetValue() (bool, error)

	// SetValue 设置高低电平
	SetValue(int) error

	// SetHigh 设置高电平
	SetHigh() error

	// SetLow 设置低电平
	SetLow() error

	// Close 关闭引脚
	Close() error
}

func List() ([]Pin, error) {
	list := []Pin(nil)
	err := oss.RangeFileInfo(gpioDir, func(info *oss.FileInfo) (bool, error) {
		if info.IsDir() && regexp.MustCompile(`(gpio)\d+`).MatchString(info.Name()) {
			after, _ := strings.CutPrefix(info.Name(), "gpio")
			list = append(list, &pin{number: conv.Int(after)})
		}
		return true, nil
	})
	return list, err
}

func Get(number int) (Pin, error) {
	if oss.Exists(fmt.Sprintf("%s/gpio%d", gpioDir, number)) {
		return &pin{number: number}, nil
	}
	return nil, errors.New("pin未开启")
}

// Open 打开引脚
func Open(number int) (Pin, error) {
	p := &pin{number: number}
	return p, p.Open()
}

// Close 关闭引脚
func Close(number int) error {
	_, err := shell.SH.Execf("echo %d > %s/unexport", number, gpioDir)
	return err
}

type pin struct {
	number int
}

func (this *pin) Open() error {
	//判断pin是否已经打开
	if oss.Exists(fmt.Sprintf("%s/gpio%d", gpioDir, this.number)) {
		return nil
	}
	_, err := shell.SH.Execf("echo %d > %s/export", this.number, gpioDir)
	if err == nil {
		_, err = shell.SH.Execf("chmod 777 %s/gpio%d/value", gpioDir, this.number)
	}
	return err
}

func (this *pin) GetNumber() int {
	return this.number
}

func (this *pin) GetModel() (Model, error) {
	result, err := shell.SH.Execf("cat %s/gpio%d/direction", gpioDir, this.number)
	if err != nil {
		return "", err
	}
	mode, _ := strings.CutSuffix(result.String(), "\n")
	return Model(mode), nil
}

func (this *pin) SetModel(model Model) error {
	_, err := shell.SH.Execf("echo %d > %s/direction", model, gpioDir)
	return err
}

func (this *pin) GetValue() (bool, error) {
	result, err := shell.SH.Execf("cat %s/gpio%d/value", gpioDir, this.number)
	if err != nil {
		return false, err
	}
	return result.String() == "1\n", nil
}

func (this *pin) SetValue(b int) error {
	_, err := shell.SH.Execf("echo %d > %s/gpio%d/value", b, gpioDir, this.number)
	return err
}

func (this *pin) SetHigh() error {
	return this.SetValue(High)
}

func (this *pin) SetLow() error {
	return this.SetValue(Low)
}

func (this *pin) Close() error {
	return Close(this.number)
}
