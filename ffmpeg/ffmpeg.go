package ffmpeg

import "fmt"

func New() *Ffmpeg {
	return &Ffmpeg{}
}

type Ffmpeg struct {
	Input  string                 //输入资源
	Output string                 //输出资源
	Args   map[string]interface{} //参数
}

func (this *Ffmpeg) SetInput(input string) *Ffmpeg {
	this.Input = input
	return this
}

func (this *Ffmpeg) SetOutput(output string) *Ffmpeg {
	this.Output = output
	return this
}

// SetVideoCodec 设置输出的视频编码
func (this *Ffmpeg) SetVideoCodec(codec string) *Ffmpeg {
	this.Args["-codec:v"] = codec
	return this
}

// SetAudioCodec 设置输出的音频编码
func (this *Ffmpeg) SetAudioCodec(codec string) *Ffmpeg {
	this.Args["-codec:a"] = codec
	return this
}

// SetFrame 设置视频总帧数
func (this *Ffmpeg) SetFrame(n int) *Ffmpeg {
	this.Args["-vframes"] = n
	return this
}

// SetFPS 设置帧率
func (this *Ffmpeg) SetFPS(fps int) *Ffmpeg {
	this.Args["-r"] = fps
	return this
}

// SetWidthHeight 设置视频长宽,像素
func (this *Ffmpeg) SetWidthHeight(w, h int) *Ffmpeg {
	this.Args["-s"] = fmt.Sprintf("%dx%d", w, h)
	return this
}

// SetAspect 设置长宽比,例16:9
func (this *Ffmpeg) SetAspect(w, h float32) *Ffmpeg {
	this.Args["-aspect"] = fmt.Sprintf("%f:%f", w, h)
	return this
}

// IgnoredVideo 不处理视频
func (this *Ffmpeg) IgnoredVideo() *Ffmpeg {
	this.Args["-vn"] = nil
	return this
}

// IgnoredAudio 不处理音频
func (this *Ffmpeg) IgnoredAudio() *Ffmpeg {
	this.Args["-an"] = nil
	return this
}
