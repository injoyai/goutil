package notice

import "github.com/injoyai/io"

type IO struct {
	*io.Client
}

func (this *IO) Publish(msg *Message) error {
	_, err := this.WriteString(msg.Content)
	return err
}
