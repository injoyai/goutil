package gpio

type Model string

const (
	In      Model = "in"
	Out     Model = "out"
	gpioDir       = "/sys/class/gpio/"

	High = true
	Low  = false
)

type PIN interface {
	// GetModel 获取当前引脚模式 输入/输出
	GetModel() (Model, error)
	SetModel(Model) error
	GetValue() (bool, error)
	SetValue(bool) error
	// SetHigh 设置高电平
	SetHigh() error
	// SetLow 设置低电平
	SetLow() error
	// Close 关闭引脚
	Close() error
}

func OpenPin(pinNum int, models ...Model) (PIN, error) {
	if err := Open(pinNum); err != nil {
		return nil, err
	}
	model := Out
	if len(models) > 0 {
		model = models[0]
	}
	if err := SetModel(model); err != nil {
		return nil, err
	}

	return &pin{number: pinNum}, nil
}

type pin struct {
	number int
}

func (this *pin) GetModel() (Model, error) {
	return GetModel(this.number)
}

func (this *pin) SetModel(model Model) error {
	return SetModel(model)
}

func (this *pin) GetValue() (bool, error) {
	return GetValue(this.number)
}

func (this *pin) SetValue(b bool) error {
	return SetValue(this.number, b)
}

func (this *pin) SetHigh() error {
	return SetValue(this.number, High)
}

func (this *pin) SetLow() error {
	return SetValue(this.number, Low)
}

func (this *pin) Close() error {
	return Close(this.number)
}
