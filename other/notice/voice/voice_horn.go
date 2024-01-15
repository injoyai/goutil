package voice

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/injoyai/goutil/g"
	"sync"
)

// Speak 生成音频并播放
func Speak(msg string) error {
	return newLocal(nil).Call(&Message{
		Param: msg,
	})
}

// Save 保存成文件 ./wav xxx
func Save(path, msg string) error {
	return newLocal(nil).Save(path, msg)
}

func newLocal(cfg *LocalConfig) *local {
	if cfg == nil {
		cfg = &LocalConfig{
			Rate:   0,
			Volume: 100,
		}
	}
	return &local{cfg: cfg}
}

var mu sync.Mutex

type LocalConfig struct {
	Rate   int //语速
	Volume int //音量
}

type local struct {
	cfg *LocalConfig
}

//func (this *local) Publish(msg *notice.Message) error {
//	mu.Lock()
//	defer mu.Unlock()
//	if err := ole.CoInitialize(0); err != nil {
//		return err
//	}
//	unknown, err := oleutil.CreateObject("SAPI.SpVoice")
//	if err != nil {
//		return err
//	}
//	voice, err := unknown.QueryInterface(ole.IID_IDispatch)
//	if err != nil {
//		return err
//	}
//	_, err = oleutil.PutProperty(voice, "Rate", this.cfg.Rate)
//	if err != nil {
//		return err
//	}
//	_, err = oleutil.PutProperty(voice, "Volume", this.cfg.Volume)
//	if err != nil {
//		return err
//	}
//	_, err = oleutil.CallMethod(voice, "Speak", msg.Content)
//	if err != nil {
//		return err
//	}
//	_, err = oleutil.CallMethod(voice, "WaitUntilDone", 0)
//	if err != nil {
//		return err
//	}
//	voice.Release()
//	ole.CoUninitialize()
//	return nil
//}

func (this *local) Call(msg *Message) error {
	mu.Lock()
	defer mu.Unlock()
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
	_, err = oleutil.CallMethod(voice, "Speak", msg.Param)
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

func (this *local) Save(path, msg string) (err error) {
	mu.Lock()
	defer mu.Unlock()
	defer g.Recover(&err)
	g.PanicErr(ole.CoInitialize(0))
	unknown, err := oleutil.CreateObject("SAPI.SpVoice")
	g.PanicErr(err)
	voice, err := unknown.QueryInterface(ole.IID_IDispatch)
	g.PanicErr(err)
	saveFile, err := oleutil.CreateObject("SAPI.SpFileStream")
	g.PanicErr(err)
	ff, err := saveFile.QueryInterface(ole.IID_IDispatch)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(ff, "Open", path, 3, true)
	g.PanicErr(err)
	_, err = oleutil.PutPropertyRef(voice, "AudioOutputStream", ff)
	g.PanicErr(err)
	_, err = oleutil.PutProperty(voice, "Rate", this.cfg.Rate)
	g.PanicErr(err)
	_, err = oleutil.PutProperty(voice, "Volume", this.cfg.Volume)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(voice, "Speak", msg)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(voice, "WaitUntilDone", 0)
	g.PanicErr(err)
	_, err = oleutil.CallMethod(ff, "Close")
	g.PanicErr(err)
	ff.Release()
	voice.Release()
	ole.CoUninitialize()
	return nil
}
