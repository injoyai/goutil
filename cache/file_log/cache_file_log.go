package file_log

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/injoyai/conv"
)

func New(cfg *Config) *FileLog {
	if len(cfg.Dir) == 0 {
		cfg.Dir = "./output/log/"
	}
	if len(cfg.Layout) == 0 {
		cfg.Layout = "2006-01-02-15.log"
	}
	if len(cfg.Split) == 0 {
		cfg.Split = "\n"
	}
	if cfg.FileSize == 0 {
		cfg.FileSize = 1024 * 1024 * 5
	}
	if cfg.SaveTime == 0 {
		cfg.SaveTime = time.Hour * 24 * 10
	}
	return &FileLog{
		Dir:      cfg.Dir,
		Layout:   cfg.Layout,
		Split:    cfg.Split,
		FileSize: cfg.FileSize,
		SaveTime: cfg.SaveTime,
	}
}

type Config struct {
	Dir      string
	Layout   string
	SaveTime time.Duration
	Split    string
	FileSize int
}

type FileLog struct {

	//日志保存时间,当写入时并且新建文件时触发
	SaveTime time.Duration

	//写入文件内容分隔符 默认 "\n"
	//选用数据中不会出现的字符,否则数据中包含该分隔符的话会丢失(后半部分)
	Split string

	//单位字节 1MB=1024*1024=1048576
	//一个日志文件的大小,约等于,取决于最后写入字节的大小
	//最大值是 设置的FileSize+最后写入字节-1
	FileSize int

	//文件保存目录 默认 "./output/log/"
	Dir string

	//文件命名规则 例如 "日志2006-01-02-15.log"
	Layout string

	currentCache    []*Data  //当前缓存
	currentFile     *os.File //当前文件
	currentFilename string   //当前文件名称
}

// WriteAny 写入任意数据,根据配置写入到不同的文件
func (this *FileLog) WriteAny(p any) (int, error) {
	return this.Write([]byte(conv.String(p)))
}

// WriteString 写入字符串
func (this *FileLog) WriteString(s string) (int, error) {
	return this.Write([]byte(s))
}

func (this *FileLog) Write(p []byte) (n int, err error) {

	//获取当前时间,以便生成预期文件名称
	now := time.Now()
	layout := now.Format(this.Layout)
	//预期文件名称,预期的目录
	expectFilename := filepath.Join(this.Dir, layout)
	expectDir := filepath.Dir(expectFilename)

	//判断当前文件是否和预期一直,初始化操作
	if this.currentFilename != expectFilename {

		//关闭老文件
		if this.currentFile != nil {
			_ = this.currentFile.Close()
		}

		//删除保留时间之外的文件,错误也问题不大,下次还会尝试删除
		fs, _ := os.ReadDir(expectDir)
		for _, v := range fs {
			if !v.IsDir() && this.SaveTime > 0 && v.Name() < now.Add(-this.SaveTime).Format(this.Layout) {
				os.Remove(this.Dir + v.Name())
			}
		}

		//新建日志文件
		os.MkdirAll(expectDir, 0666)
		file, err := os.OpenFile(expectFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return 0, err
		}
		this.currentCache = nil
		this.currentFile = file
		this.currentFilename = expectFilename

	}

	//写入数据到文件
	t := now.UnixNano()
	bs := append([]byte(strconv.FormatInt(t, 10)+":"), p...)
	n, err = this.currentFile.Write(append(bs, this.Split...))
	if err == nil {
		this.currentCache = append(this.currentCache, &Data{Time: t, Bytes: p})
	}
	return n, err
}

type Search struct {
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
	PageSize  int       `json:"pageSize"`
}

func (this *FileLog) GetLog(s *Search) ([]*Data, error) {

	if s.EndTime.IsZero() {
		s.EndTime = time.Now()
	}

	//生成文件的范围名称
	startFilename := s.StartTime.Format(filepath.Join(this.Dir, this.Layout))
	endFilename := s.EndTime.Format(filepath.Join(this.Dir, this.Layout))
	list := []*Data(nil)

	//读取目录
	fs, err := os.ReadDir(this.Dir)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		filename := filepath.Join(this.Dir, f.Name())
		if filename >= startFilename && filename <= endFilename {
			ls, err := this.readFile(filename)
			if err != nil {
				return nil, err
			}
			for _, bs := range ls {
				list = append(list, bs)
				if s.PageSize > 0 && len(list) >= s.PageSize {
					return list, nil
				}
			}
		}
	}
	return list, nil
}

func (this *FileLog) readFile(filename string) ([]*Data, error) {
	//判断读取是否是当前的文件,是则返回缓存
	if filename == this.currentFilename {
		return this.currentCache, nil
	}
	//打开文件,从文件中读取数据
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//分割字节
	list := []*Data(nil)
	for _, v := range bytes.Split(bs, []byte(this.Split)) {
		//拆分时间戳
		l2 := bytes.SplitN(v, []byte(":"), 2)
		if len(l2) == 2 && len(l2[1]) > 0 {
			list = append(list, &Data{
				Time:  conv.Int64(string(l2[0])),
				Bytes: l2[1],
			})
		}
	}
	return list, nil
}

// GetLogCurve 获取日志并生成曲线图,统计间隔至少1秒
func (this *FileLog) GetLogCurve(start, end time.Time, merge time.Duration, d Decoder) (*Curve, error) {
	//统计间隔
	interval := merge
	interval = conv.Select[time.Duration](interval < time.Second, time.Second, interval)
	//获取日志列表
	list, err := this.GetLog(&Search{
		StartTime: start.Add(-interval),
		EndTime:   end,
		PageSize:  0,
	})
	if err != nil {
		return nil, err
	}

	//填充没有数据的节点
	curve := &Curve{}
	mapData := map[int64]*[]any{}
	for t := start.UnixNano(); t <= end.UnixNano(); t += int64(interval) {
		ls := &[]any{}
		curve.Time = append(curve.Time, t)
		curve.list = append(curve.list, ls)
		curve.Value = append(curve.Value, nil)
		mapData[t] = ls
	}

	//合并数据
	for _, v := range list {

		a, err := d.Decode(v.Bytes)
		if err != nil {
			return nil, err
		}

		//把n-m的数据归类到m,而不是n
		node := v.Time - (v.Time-start.UnixNano())%int64(interval) + int64(interval)
		if _, ok := mapData[node]; ok {
			*mapData[node] = append(*mapData[node], a)
		}
	}

	for i, v := range curve.list {
		curve.Value[i], err = d.Report(curve.Time[i], *v)
		if err != nil {
			return nil, err
		}
	}

	return curve, nil
}

type Data struct {
	Time  int64  `json:"time"`
	Bytes []byte `json:"bytes"`
}

type Curve struct {
	Time  []int64  `json:"time"`
	Value []any    `json:"value"` //计算后的值
	list  []*[]any //未计算的原始值
}

func (this *Curve) String() string {
	s := ""
	for i, v := range this.Value {
		if i != 0 {
			s += ","
		}
		s += conv.String(v)
	}
	return s
}

type Decoder interface {
	// Decode 字节转对象,这个是接口速度的关键,不推荐使用json
	Decode([]byte) (any, error)
	// Report 整理对象,合并统计,曲线的按时间统计,平均值或者最大值等等
	Report(node int64, list []any) (any, error)
}

type DecodeFunc struct {
	DecodeFunc func([]byte) (any, error)
	ReportFunc func(node int64, list []any) (any, error)
}

func (this *DecodeFunc) Decode(bs []byte) (any, error) {
	return this.DecodeFunc(bs)
}

func (this *DecodeFunc) Report(node int64, list []any) (any, error) {
	return this.ReportFunc(node, list)
}
