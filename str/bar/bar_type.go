package bar

import (
	"github.com/fatih/color"
	"io"
)

type Interface interface {

	/*
		SetFormatter 自定义格式
		func(e *Entity) string {
			return fmt.Sprintf("\r%s%s %s %s %s %s",
				this.prefix,
				e.Bar,
				e.Rate,
				e.Size,
				e.Remain,
				this.suffix,
			)
	*/
	SetFormatter(f Formatter) Interface

	// SetWidth 设置宽度
	SetWidth(width int) Interface

	// SetTotal 设置总数据大小
	SetTotal(total int64) Interface

	// SetStyle 设置进度条风格
	SetStyle(style byte) Interface

	// SetColor 设置颜色
	SetColor(color color.Attribute) Interface

	// Add 添加数据
	Add(n int64) Interface

	// Done 结束
	Done() Interface

	// Run 运行
	Run() <-chan struct{}

	//衍生功能,

	// Copy 复制数据加入进度条
	Copy(w io.Writer, r io.Reader) error

	// CopyN 复制数据加入进度条
	CopyN(w io.Writer, r io.Reader, num int64) error

	// DownloadHTTP 下载http
	DownloadHTTP(url, filename string) error
}

type Element interface {
	String() string
}

type element func() string

func (this element) String() string { return this() }

type Entity struct {
	Bar      Element
	Rate     Element
	Size     Element
	SizeUnit Element
	Speed    Element
	Used     Element
	Remain   Element
}

type Formatter func(e *Entity) string

func ToB(b int64) (float64, string) {
	var mapB = map[int]string{
		0:   "B",
		10:  "KB",
		20:  "MB",
		30:  "GB",
		40:  "TB",
		50:  "PB",
		60:  "EB",
		70:  "ZB",
		80:  "YB",
		90:  "BB",
		100: "NB",
		110: "DB",
		120: "CB",
		130: "XB",
	}

	for n := 0; n <= 130; n += 10 {
		if b < 1<<(n+10) {
			if n == 0 {
				return float64(b), mapB[n]
			}
			return float64(b) / float64(int64(1)<<n), mapB[n]
		}
	}
	return float64(b), mapB[0]
}
