package notice

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"sync"
)

var audioMu sync.Mutex

type AudioConfig struct {
	Rate   int //语速
	Volume int //音量
}

type audio struct {
	cfg *AudioConfig
}

func (this *audio) Publish(msg *Message) error {
	audioMu.Lock()
	defer audioMu.Unlock()
	if err := ole.CoInitialize(0); err != nil {
		return err
	}
	unknown, err := oleutil.CreateObject("SAPI.SpVoice")
	if err != nil {
		return err
	}
	voice, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(voice, "Rate", this.cfg.Rate)
	if err != nil {
		return err
	}
	_, err = oleutil.PutProperty(voice, "Volume", this.cfg.Volume)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(voice, "Speak", msg.Content)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(voice, "WaitUntilDone", 0)
	if err != nil {
		return err
	}
	voice.Release()
	ole.CoUninitialize()
	return nil
}
