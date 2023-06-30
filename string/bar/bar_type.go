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

	// SetPrefix 设置前缀,默认格式生效
	SetPrefix(prefix string) Interface

	// SetSuffix 设置后缀,默认格式生效
	SetSuffix(suffix string) Interface

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
	Bar    Element
	Rate   Element
	Size   Element
	Used   Element
	Remain Element
}

type Formatter func(e *Entity) string
