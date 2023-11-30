package gpio

import (
	"github.com/injoyai/goutil/oss/linux/bash"
)

func Open(pinNum int) error {
	_, err := bash.Execf("echo %d > %s/export", pinNum, gpioDir)
	if err == nil {
		bash.Execf("chmod 777 %s/gpio%d/value", gpioDir, pinNum)
	}
	return err
}

func Close(pinNum int) error {
	_, err := bash.Execf("echo %d > %s/unexport", pinNum, gpioDir)
	return err
}

func GetModel(pinNum int) (Model, error) {
	result, err := bash.Execf("cat %s/gpio%d/direction", gpioDir, pinNum)
	return Model(result), err
}

func SetModel(Type Model) error {
	_, err := bash.Execf("echo %d > %s/direction", Type, gpioDir)
	return err
}

func GetValue(pinNum int) (bool, error) {
	result, err := bash.Execf("cat %s/gpio%d/value", gpioDir, pinNum)
	return result == "1", err
}

func SetValue(pinNum int, value bool) error {
	_, err := bash.Execf("echo %d > %s/gpio%d/value", value, gpioDir, pinNum)
	return err
}
