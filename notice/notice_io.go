package notice

import "github.com/injoyai/io"

func NewIO(dial io.DialFunc, options ...io.OptionClient) (Interface, error) {
	c, err := io.NewDial(dial, options...)
	if err != nil {
		return nil, err
	}
	return &IO{Client: c.Redial()}, nil
}

type IO struct {
	*io.Client
}

func (this *IO) Publish(msg *Message) error {
	_, err := this.WriteAny(msg)
	return err
}
