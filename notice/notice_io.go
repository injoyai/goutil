package notice

import (
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
)

func NewIO(dial ios.DialFunc, options ...client.Option) (Interface, error) {
	c, err := client.Dial(dial, options...)
	if err != nil {
		return nil, err
	}
	c.Redial()
	return &IO{Client: c}, nil
}

type IO struct {
	*client.Client
}

func (this *IO) Publish(msg *Message) error {
	return this.WriteAny(msg)
}
