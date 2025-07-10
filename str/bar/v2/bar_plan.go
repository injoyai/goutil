package bar

import (
	"fmt"
	"github.com/fatih/color"
)

type plan struct {
	tag    string       //tag 例 [导出] [>>>   ]
	prefix string       //前缀 例 [
	suffix string       //后缀 例 ]
	style  byte         //进度条风格 例 >
	color  *color.Color //整体颜色
	width  int          //宽度

	current int64 //当前
	total   int64 //总数
}

func (this *plan) SetTag(tag string) {
	this.tag = tag
}

func (this *plan) SetPrefix(prefix string) {
	this.prefix = prefix
}

func (this *plan) SetSuffix(suffix string) {
	this.suffix = suffix
}

func (this *plan) SetStyle(style byte) {
	this.style = style
}

func (this *plan) SetWidth(width int) {
	this.width = width
}

func (this *plan) SetColor(a color.Attribute) {
	this.color = color.New(a)
}

func (this *plan) String() string {
	rate := float64(this.current) / float64(this.total)
	nowWidth := ""
	for i := 0; i < int(float64(this.width)*rate); i++ {
		nowWidth += string(this.style)
	}

	barStr := fmt.Sprintf(fmt.Sprintf("%s%%-%ds%s", this.prefix, this.width, this.suffix), nowWidth)
	if this.color != nil {
		barStr = this.color.Sprint(barStr)
	}
	if len(this.tag) > 0 {
		barStr = fmt.Sprintf("[%s] %s", this.tag, barStr)
	}
	return barStr
}
