package script

import (
	"github.com/injoyai/base/maps"
	"github.com/injoyai/base/maps/wait"
	"github.com/injoyai/base/types"
	"github.com/injoyai/conv"
	"github.com/injoyai/conv/cfg"
	"github.com/injoyai/conv/codec"
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/mux"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/goutil/net/ip"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/oss/shell"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/client/dial"
	"github.com/injoyai/logs"
	"io"
	"net"
	gohttp "net/http"
	"runtime"
	"time"
)

func NewGlobal() *Global {
	return &Global{
		Map: maps.NewSafe(),
	}
}

type Global struct {
	Map *maps.Safe
}

type Net struct{}

func (this *Net) Dial(network, address string) net.Conn {
	switch network {
	case "icmp":
		return &Pinger{ip.NewPinger().SetHost([]string{address})}
	}
	c, err := net.Dial(network, address)
	if err != nil {
		panic(err)
	}
	return c
}

func (this *Net) DialTCP(address string) net.Conn {
	return this.Dial("tcp", address)
}

func (this *Net) DialUDP(address string) net.Conn {
	return this.Dial("udp", address)
}

func (this *Net) DialICMP(hosts ...string) *ip.Pinger {
	return ip.NewPinger().SetHost(hosts)
}

type Pinger struct {
	*ip.Pinger
}

func (this *Pinger) For(num int, interval float64) {
	for i := 0; i < num; i++ {
		<-time.After(time.Millisecond * time.Duration(interval*1000))
		this.Ping()
	}
}

type Ios struct{}

func (this *Ios) Dial(args *Args) *client.Client {
	switch args.GetString(1) {
	case "tcp":
		c, err := dial.TCP(args.GetString(2))
		if err != nil {
			panic(err)
		}
		return c
	case "mqtt":

	}
	return nil
}

func NewLogs() *Logs {
	return &Logs{
		TimeFormatter: logs.TimeFormatter,
		LevelAll:      logs.LevelAll,
		LevelTrace:    logs.LevelTrace,
		LevelDebug:    logs.LevelDebug,
		LevelWrite:    logs.LevelWrite,
		LevelRead:     logs.LevelRead,
		LevelInfo:     logs.LevelInfo,
		LevelWarn:     logs.LevelWarn,
		LevelError:    logs.LevelError,
		LevelNone:     logs.LevelNone,
	}
}

type Logs struct {
	TimeFormatter logs.IFormatter
	LevelAll      logs.Level
	LevelTrace    logs.Level
	LevelDebug    logs.Level
	LevelWrite    logs.Level
	LevelRead     logs.Level
	LevelInfo     logs.Level
	LevelWarn     logs.Level
	LevelError    logs.Level
	LevelNone     logs.Level
}

func (this *Logs) Debug(args ...any) (int, error) {
	return logs.Debug(args...)
}

func (this *Logs) Debugf(format string, args ...any) (int, error) {
	return logs.Debugf(format, args...)
}

func (this *Logs) Info(args ...any) (int, error) {
	return logs.Info(args...)
}

func (this *Logs) Infof(format string, args ...any) (int, error) {
	return logs.Infof(format, args...)
}

func (this *Logs) Err(args ...any) (int, error) {
	return logs.Err(args...)
}

func (this *Logs) Errf(format string, args ...any) (int, error) {
	return logs.Errf(format, args...)
}

func (this *Logs) Error(args ...any) (int, error) {
	return logs.Error(args...)
}

func (this *Logs) Errorf(format string, args ...any) (int, error) {
	return logs.Errorf(format, args...)
}

func NewHTTP() *HTTP {
	return &HTTP{
		DefaultClient: http.DefaultClient,
	}
}

type HTTP struct {
	DefaultClient *http.Client
}

func (this *HTTP) Url(url string) *http.Request {
	return http.Url(url)
}

func NewConv() *Conv {
	return &Conv{
		Json: codec.Json,
		Toml: codec.Toml,
		Yaml: codec.Yaml,
		Ini:  codec.Ini,
	}
}

type Conv struct {
	Json codec.Interface
	Toml codec.Interface
	Yaml codec.Interface
	Ini  codec.Interface
}

func (this *Conv) New(i any) *conv.Var {
	return conv.New(i)
}

func (this *Conv) NewMap(i any, codec ...codec.Interface) *conv.Map {
	return conv.NewMap(i, codec...)
}

func (this *Conv) Bytes(i any) types.Bytes {
	return conv.Bytes(i)
}

func (this *Conv) String(i any) string {
	return conv.String(i)
}

func (this *Conv) Int(i any) int {
	return conv.Int(i)
}

func (this *Conv) Float(i any) float64 {
	return conv.Float64(i)
}

func (this *Conv) Bool(i any) bool {
	return conv.Bool(i)
}

func (this *Conv) Interfaces(i any) []any {
	return conv.Interfaces(i)
}

func (this *Conv) DMap(i any) *conv.Map {
	return conv.DMap(i)
}

func (this *Conv) GMap(i any) map[string]any {
	return conv.GMap(i)
}

func (this *Conv) SMap(i any) map[string]string {
	return conv.SMap(i)
}

func (this *Conv) Duration(i any) time.Duration {
	return conv.Duration(i)
}

func NewCfg() *Cfg {
	return &Cfg{
		Extend: cfg.Default,
	}
}

type Cfg struct {
	conv.Extend
}

func (this *Cfg) Init(i ...conv.IGetVar) {
	cfg.Init(i...)
}

func (this *Cfg) New(i ...conv.IGetVar) *cfg.Entity {
	return cfg.New(i...)
}

func (this *Cfg) WithAny(i any) conv.IGetVar {
	return cfg.WithAny(i)
}

func (this *Cfg) WithFile(path string, codec ...codec.Interface) conv.IGetVar {
	return cfg.WithFile(path, codec...)
}

func (this *Cfg) WithEnv() conv.IGetVar {
	return cfg.WithEnv()
}

func NewOS() *OS {
	return &OS{
		CPUNum:  runtime.NumCPU(),
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		Version: runtime.Version(),
		Root:    runtime.GOROOT(),
		SH:      shell.SH,
		Bash:    shell.Bash,
		CMD:     shell.CMD,
	}
}

type OS struct {
	CPUNum  int
	OS      string
	Arch    string
	Version string
	Root    string
	SH      *shell.Shell
	Bash    *shell.Shell
	CMD     *shell.Shell
}

func (this *OS) New(filename string, i ...any) {
	panicErr(oss.New(filename, i...))
}

func (this *OS) NewNotExist(filename string, i ...any) {
	panicErr(oss.NewNotExist(filename, i...))
}

func (this *OS) Read(filename string) types.Bytes {
	bs, err := oss.Read(filename)
	panicErr(err)
	return bs
}

func (this *OS) Shell(args ...string) string {
	result, err := shell.Exec(args...)
	panicErr(err)
	return result.String()
}

func (this *OS) Start(filename string) {
	err := shell.Start(filename)
	panicErr(err)
}

func (this *OS) Stop(filename string) {
	err := shell.Stop(filename)
	panicErr(err)
}

type Bytes struct{}

func (this *Bytes) New(i any) types.Bytes {
	return conv.Bytes(i)
}

func (this *Bytes) Sum(bs []byte) byte {
	return types.Bytes(bs).Sum()
}

func (this *Bytes) Reverse(bs []byte) types.Bytes {
	return types.Bytes(bs).Reverse()
}

func (this *Bytes) Upper(bs []byte) types.Bytes {
	return types.Bytes(bs).Upper()
}

func (this *Bytes) Lower(bs []byte) types.Bytes {
	return types.Bytes(bs).Lower()
}

func (this *Bytes) Base64(bs []byte) string {
	return types.Bytes(bs).Base64()
}

func (this *Bytes) BIN(bs []byte) string {
	return types.Bytes(bs).BIN()
}

func (this *Bytes) Int(bs []byte) int {
	return int(types.Bytes(bs).Int64())
}

func (this *Bytes) Add(bs []byte, b byte) types.Bytes {
	return types.Bytes(bs).AddByte(b)
}

func (this *Bytes) Sub(bs []byte, b byte) types.Bytes {
	return types.Bytes(bs).SubByte(b)
}

type Mux struct{}

func (this *Mux) New() *mux.Server {
	return mux.New()
}

func (this *Mux) Json200(i any) {
	in.Json200(i)
}

func (this *Mux) Succ(i any) {
	in.Succ(i)
}

func NewIn() *In {
	return &In{
		DefaultClient: in.DefaultClient,
	}
}

type In struct {
	DefaultClient *in.Client
}

func (this *In) Exit() {
	in.DefaultClient.Exit()
}

func (this *In) Recover(h gohttp.Handler) gohttp.Handler {
	return in.Recover(h)
}

func (this *In) SetHandlerWithCode(succ, fail, unauthorized, forbidden any) *in.Client {
	return in.SetHandlerWithCode(succ, fail, unauthorized, forbidden)
}

func (this *In) Succ(i any) {
	in.Succ(i)
}

func (this *In) Fail(i any) {
	in.Fail(i)
}

func (this *In) Forbidden() {
	in.Forbidden()
}

func (this *In) Unauthorized() {
	in.Unauthorized()
}

func (this *In) Err(data any, succData ...any) {
	in.Err(data, succData...)
}

func (this *In) CheckErr(err error) {
	in.CheckErr(err)
}

func (this *In) FileLocal(name, filename string) {
	in.FileLocal(name, filename)
}

func (this *In) FileReader(name string, r io.ReadCloser) {
	in.FileReader(name, r)
}

func (this *In) FileBytes(name string, bs []byte) {
	in.FileBytes(name, bs)
}

func (this *In) Text(code int, data any) {
	in.Text(code, data)
}

func (this *In) Text200(data any) {
	in.Text200(data)
}

func (this *In) Json(code int, data any) {
	in.Json(code, data)
}

func (this *In) Json200(data any) {
	in.Json200(data)
}

func (this *In) Html(code int, data any) {
	in.Html(code, data)
}

func (this *In) Html200(data any) {
	in.Html200(data)
}

func (this *In) Reader(code int, r io.ReadCloser) {
	in.Reader(code, r)
}

func (this *In) Reader200(r io.ReadCloser) {
	in.Reader200(r)
}

type Maps struct{}

func (this *Maps) New() *maps.Safe {
	return maps.NewSafe()
}

func (this *Maps) NewSafe() *maps.Safe {
	return maps.NewSafe()
}

func (this *Maps) NewWait(f float64) *wait.Entity {
	return wait.New(floatToDuration(f))
}

func NewTime() *Time {
	return &Time{
		Day:         time.Hour * 24,
		Hour:        time.Hour,
		Minute:      time.Minute,
		Second:      time.Second,
		Millisecond: time.Millisecond,
		Microsecond: time.Microsecond,
		Nanosecond:  time.Nanosecond,
	}
}

type Time struct {
	Day         time.Duration
	Hour        time.Duration
	Minute      time.Duration
	Second      time.Duration
	Millisecond time.Duration
	Microsecond time.Duration
	Nanosecond  time.Duration
}

func (this *Time) Now() time.Time {
	return time.Now()
}

func (this *Time) Unix(sec, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func (this *Time) Date(year, month, day, hour, min, sec int) time.Time {
	return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
}

func (this *Time) Sleep(d time.Duration) {
	time.Sleep(d)
}

/*



 */

func floatToDuration(f float64) time.Duration {
	return time.Millisecond * time.Duration(f*1000)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
