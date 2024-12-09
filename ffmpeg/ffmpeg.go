package ffmpeg

import (
	"fmt"
	"github.com/injoyai/goutil/oss/shell"
	"log"
	"strings"
)

var DefaultDebug = false

// Changer 格式转换,例mp4转avi ffmpeg -i xxx.mp4 xxx.avi
func Changer(input, output string) error {
	return New().Input(input).Output(output).Cover().Run()
}

// Merge 合并视频,通过文件
// 例 ffmpeg -f concat -i tslist.txt -codec copy out_ts.ts
// 例 ffmpeg -i "concat:1.ts|2.ts|3.ts" -codec copy out_ts.ts
// 文件内容如下 file '1.ts'\nfile '2.ts'\nfile '3.ts'...  或者concat:1.ts|2.ts|3.ts...
func Merge(filename string, output string) error {
	if strings.HasPrefix(filename, "concat:") {
		return New().
			Input(filename).
			Output(output).
			Cover().
			Codec(Copy).
			Run()
	}
	return New().
		Input(filename, Concat).
		Output(output).
		Cover().
		Codec(Copy).
		Run()
}

// Capture 截取一张图片
func Capture(input string, at string, output string) error {
	return New().
		Input(input).
		StartTime(at).
		Output(output).
		Cover().
		Run()
}

// Gif 制作gif
func Gif(input string, at string, sec, fps, w, h int, output string) error {
	return New().
		Input(input).
		StartTime(at).
		FPS(fps).
		Duration(sec).
		Size(w, h).
		Output(output).
		Cover().
		Run()
}

// ToAudio 转换为音频
func ToAudio(input string, at string, sec int, name string) error {
	return New().
		Input(input).
		At(at).
		Duration(sec).
		Output(name).
		Cover().
		Run()
}

type Option func(*Ffmpeg)

func New(op ...Option) *Ffmpeg {
	f := &Ffmpeg{args: make(map[string]interface{})}
	for _, v := range op {
		v(f)
	}
	return f
}

type Ffmpeg struct {
	shell  func(string) error     //shell
	bin    string                 //ffmpeg命令,可空,默认ffmpeg
	input  []string               //输入资源
	output string                 //输出资源
	args   map[string]interface{} //参数
	debug  *bool                  //打印命令
}

// Debug 调试,打印命令
func (this *Ffmpeg) Debug(b ...bool) *Ffmpeg {
	*this.debug = len(b) == 0 || b[0]
	return this
}

// Run 开始执行
func (this *Ffmpeg) Run() error {
	if len(this.input) == 0 {
		return fmt.Errorf("input is empty")
	}
	if len(this.output) == 0 {
		return fmt.Errorf("output is empty")
	}
	if this.shell == nil {
		this.shell = func(s string) error {
			result, err := shell.Exec(s)
			if err != nil {
				if err.Error() == "exit status 1" {
					ls := strings.Split(result.String(), "\n")
					if len(ls) > 1 {
						return fmt.Errorf(ls[len(ls)-2])
					}
					return fmt.Errorf(result.String())
				}
			}
			return err
		}
	}
	if this.bin == "" {
		this.bin = "ffmpeg"
	}
	input := ""
	for _, v := range this.input {
		input += fmt.Sprintf(" -i %s", v)
	}
	args := ""
	for k, v := range this.args {
		if v == nil {
			args += fmt.Sprintf(" %s", k)
			continue
		}
		args += fmt.Sprintf(" %s %v", k, v)
	}
	s := fmt.Sprintf("%s%s%s %s", this.bin, input, args, this.output)
	if this.debug == nil || *this.debug {
		log.Println(s)
	}
	return this.shell(s)
}

// Shell 设置shell
func (this *Ffmpeg) Shell(f func(string) error) *Ffmpeg {
	this.shell = f
	return this
}

// Arg 设置自定义参数
func (this *Ffmpeg) Arg(key string, val ...interface{}) *Ffmpeg {
	if !strings.HasPrefix(key, "-") {
		key = "-" + key
	}
	this.args[key] = func() interface{} {
		if len(val) == 0 {
			return nil
		}
		return val[0]
	}()
	return this
}

// Args 设置自定义参数
func (this *Ffmpeg) Args(args map[string]interface{}) *Ffmpeg {
	for k, v := range args {
		this.Arg(k, v)
	}
	return this
}

// Input 设置输入,例 input.mp4
func (this *Ffmpeg) Input(input string, frame ...string) *Ffmpeg {
	if len(frame) > 0 {
		this.input = append(this.input, fmt.Sprintf("-f %s %s", frame[0], input))
	} else {
		this.input = append(this.input, input)
	}
	return this
}

// Output 设置输出
func (this *Ffmpeg) Output(output string, frame ...string) *Ffmpeg {
	if len(frame) > 0 {
		this.output = fmt.Sprintf("-f %s %s", frame[0], output)
	} else {
		this.output = output
	}
	return this
}

// OutputCover 输出并覆盖原文件
func (this *Ffmpeg) OutputCover(output string, frame ...string) *Ffmpeg {
	return this.Output(output, frame...).Cover()
}

// Cover 覆盖输出文件
func (this *Ffmpeg) Cover() *Ffmpeg {
	this.Arg("-y")
	return this
}

func (this *Ffmpeg) Codec(codec string) *Ffmpeg {
	this.args["-codec"] = codec
	return this
}

// VideoCodec 设置输出的视频编码
func (this *Ffmpeg) VideoCodec(codec string) *Ffmpeg {
	this.args["-codec:v"] = codec
	return this
}

// AudioCodec 设置输出的音频编码
func (this *Ffmpeg) AudioCodec(codec string) *Ffmpeg {
	this.args["-codec:a"] = codec
	return this
}

// Frames 设置视频总帧数,或图片数量
func (this *Ffmpeg) Frames(n int) *Ffmpeg {
	this.args["-vframes"] = n
	return this
}

// At 设置起始时间
func (this *Ffmpeg) At(at string) *Ffmpeg {
	return this.StartTime(at)
}

// StartTime 设置起始时间
func (this *Ffmpeg) StartTime(s string) *Ffmpeg {
	if len(s) > 0 {
		this.args["-ss"] = s
	}
	return this
}

// Duration 设置时长,秒
func (this *Ffmpeg) Duration(sec int) *Ffmpeg {
	if sec > 0 {
		this.args["-t"] = sec
	}
	return this
}

// Format 强制设置输入输出格式,否则根据名字后缀等判断
func (this *Ffmpeg) Format(format string) *Ffmpeg {
	this.args["-f"] = format
	return this
}

// FPS 设置帧率,例ffmpeg -i xxx.mp4 -r 10 xxx.mp4		# 注意，不能加 codec copy 否之转换无效
func (this *Ffmpeg) FPS(fps int) *Ffmpeg {
	this.args["-r"] = fps
	return this
}

// Size 设置视频长宽,像素
func (this *Ffmpeg) Size(w, h int) *Ffmpeg {
	this.args["-s"] = fmt.Sprintf("%dx%d", w, h)
	return this
}

// Aspect 设置长宽比,例16:9
func (this *Ffmpeg) Aspect(w, h float32) *Ffmpeg {
	this.args["-aspect"] = fmt.Sprintf("%f:%f", w, h)
	return this
}

// DeleteVideo 删除视频
func (this *Ffmpeg) DeleteVideo() *Ffmpeg {
	this.args["-vn"] = nil
	return this
}

// DeleteAudio 删除音频
func (this *Ffmpeg) DeleteAudio() *Ffmpeg {
	this.args["-an"] = nil
	return this
}

// VideoBitrate 设置视频位率,例 ffmpeg -i xxx.mp4 -b:v 400k xxx.mp4
func (this *Ffmpeg) VideoBitrate(s string) *Ffmpeg {
	this.args["-b:v"] = s
	return this
}

// AudioBitrate 设置音频位率,常见128k、192k、256k、320k,例 ffmpeg -i xxx.mp4 -b:a 192k xxx.mp4
func (this *Ffmpeg) AudioBitrate(s string) *Ffmpeg {
	this.args["-b:a"] = s
	return this
}

// AudioSampleRate 设置音频采样率,常见有44100Hz,48000Hz,例 ffmpeg -i xxx.mp4 -ar 44100 xxx.mp4
func (this *Ffmpeg) AudioSampleRate(rate int) *Ffmpeg {
	this.args["-ar"] = rate
	return this
}

// AudioChannels 设置音频通道数量,例ffmpeg -i xxx.mp3 -ac 2 s16.wav
func (this *Ffmpeg) AudioChannels(channels int) *Ffmpeg {
	this.args["-ac"] = channels
	return this
}
