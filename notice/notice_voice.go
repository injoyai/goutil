package notice

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"sync"
)

var (
	DefaultVoice = NewVoice(&VoiceConfig{
		Rate:   0,
		Volume: 100,
	})
)

func NewVoice(cfg *VoiceConfig) *voice {
	if cfg == nil {
		return DefaultVoice
	}
	return &voice{cfg}
}

var voiceMu sync.Mutex

type VoiceConfig struct {
	Rate   int `json:"rate"`   //语速
	Volume int `json:"volume"` //音量
}

type voice struct {
	cfg *VoiceConfig
}

func (this *voice) Publish(msg *Message) error {
	return this.Speak(msg.Content)
}

func (this *voice) Speak(content string) error {
	voiceMu.Lock()
	defer voiceMu.Unlock()
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
	_, err = oleutil.CallMethod(voice, "Speak", content)
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

func (this *voice) Save(filename, msg string) error {
	voiceMu.Lock()
	defer voiceMu.Unlock()
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
	saveFile, err := oleutil.CreateObject("SAPI.SpFileStream")
	if err != nil {
		return err
	}
	ff, err := saveFile.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(ff, "Open", filename, 3, true)
	if err != nil {
		return err
	}
	_, err = oleutil.PutPropertyRef(voice, "AudioOutputStream", ff)
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
	_, err = oleutil.CallMethod(voice, "Speak", msg)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(voice, "WaitUntilDone", 0)
	if err != nil {
		return err
	}
	_, err = oleutil.CallMethod(ff, "Close")
	if err != nil {
		return err
	}
	ff.Release()
	voice.Release()
	ole.CoUninitialize()
	return nil
}
