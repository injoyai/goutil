package bar

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	x := New()
	x.SetPrefix("进度:")
	x.SetWidth(50)
	x.SetTotalSize(1100)
	//x.SetColor(color.Violet + 10)
	x.SetStyle(">")
	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(time.Millisecond * 100)
			x.Add(10)
		}
	}()
	go func() {
		time.Sleep(time.Second * 3)
		//x.Done()
	}()
	<-x.Run()
}

func TestDemo(t *testing.T) {
	Demo()
}
