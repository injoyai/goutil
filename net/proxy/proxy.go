package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/fatih/color"
	"github.com/injoyai/logs"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type Option func(*Proxy)

func WithCA(certFile, keyFile string) Option {
	return func(p *Proxy) {
		err := p.SetCA(certFile, keyFile)
		logs.PrintErr(err)
	}
}

func WithCABytes(certFile, keyFile []byte) Option {
	return func(p *Proxy) {
		err := p.SetCABytes(certFile, keyFile)
		logs.PrintErr(err)
	}
}

func WithProxy(u string) Option {
	return func(p *Proxy) {
		err := p.SetProxy(u)
		logs.PrintErr(err)
	}
}

func WithProxyPac(u string, domains []string) Option {
	return func(p *Proxy) {
		err := p.SetProxyPac(u, domains)
		logs.PrintErr(err)
	}
}

func WithDebug(b ...bool) Option {
	return func(p *Proxy) {
		p.Debug(b...)
	}
}

func WithMitm() Option {
	return func(p *Proxy) {
		p.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
			return &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&p.ca)}, host
		}))
	}
}

func WithOptions(op ...Option) Option {
	return func(p *Proxy) {
		p.SetOptions(op...)
	}
}

func WithPort(port int) Option {
	return func(p *Proxy) {
		p.SetPort(port)
	}
}

func Default(op ...Option) *Proxy {
	return New(
		WithPort(DefaultPort),
		WithCABytes([]byte(DefaultCrt), []byte(DefaultKey)), //100年的证书
		WithMitm(),
		WithOptions(op...),
	)
}

func New(op ...Option) *Proxy {
	p := &Proxy{
		ProxyHttpServer: goproxy.NewProxyHttpServer(),
		log:             logs.New("").SetFormatter(logs.TimeFormatter).SetColor(color.FgGreen),
		ca:              goproxy.GoproxyCa,
		port:            DefaultPort,
	}
	for _, v := range op {
		v(p)
	}
	return p
}

type Proxy struct {
	*goproxy.ProxyHttpServer
	log   *logs.Entity
	ca    tls.Certificate
	port  int
	debug bool
}

func (this *Proxy) SetPort(port int) {
	this.port = port
}

func (this *Proxy) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", this.port),
		Handler: this.ProxyHttpServer,
	}
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	this.log.Printf("[信息] [%s] 代理开启成功...\n", srv.Addr)
	c := make(chan struct{})
	defer close(c)
	go func() {
		select {
		case <-ctx.Done():
		case <-c:
		}
		ln.Close()
	}()
	return srv.Serve(ln)
}

// SetOptions 设置选项
func (this *Proxy) SetOptions(op ...Option) {
	for _, v := range op {
		v(this)
	}
}

// SetCABytes 设置ca证书
func (this *Proxy) SetCABytes(crt, key []byte) error {
	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return err
	}
	this.ca = cert
	return nil
}

// SetCA 设置ca证书
func (this *Proxy) SetCA(certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	this.ca = cert
	return nil
}

// SetProxy 设置代理
func (this *Proxy) SetProxy(u string) error {
	if len(u) == 0 {
		this.ProxyHttpServer.Tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
		}
		return nil
	}

	proxyUrl, err := url.Parse(u)
	if err != nil {
		return err
	}
	t := &http.Transport{}
	switch proxyUrl.Scheme {
	case "socks5", "socks5h":
		dialer, err := proxy.FromURL(proxyUrl, this)
		if err != nil {
			return err
		}
		t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	default: //"http", "https"
		t.Proxy = http.ProxyURL(proxyUrl)
	}
	this.ProxyHttpServer.Tr = t
	return nil
}

func (this *Proxy) SetProxyPac(u string, domains []string) error {
	if len(u) == 0 {
		this.ProxyHttpServer.Tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
		}
		return nil
	}

	m := map[string]struct{}{}
	for _, v := range domains {
		m[v] = struct{}{}
	}

	f := func(host string) bool {
		ls := strings.Split(host, ".")
		if len(ls) >= 2 {
			_, ok := m[ls[len(ls)-2]+"."+ls[len(ls)-1]]
			return ok
		}
		return false
	}

	proxyUrl, err := url.Parse(u)
	if err != nil {
		return err
	}
	t := &http.Transport{}
	switch proxyUrl.Scheme {
	case "socks5", "socks5h":
		dialer, err := proxy.FromURL(proxyUrl, this)
		if err != nil {
			return err
		}
		t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			if f(addr) {
				return dialer.Dial(network, addr)
			}
			return this.Dial(network, addr)
		}
	default: //"http", "https"
		t.Proxy = func(req *http.Request) (*url.URL, error) {
			if f(req.Host) {
				return proxyUrl, nil
			}
			return nil, nil
		}
	}
	this.ProxyHttpServer.Tr = t
	return nil
}

// Debug 打印通讯数据
func (this *Proxy) Debug(b ...bool) {
	this.ProxyHttpServer.Verbose = len(b) == 0 || b[0]
	this.debug = len(b) == 0 || b[0]
}

// Dial 实现接口
func (this *Proxy) Dial(network, addr string) (c net.Conn, err error) {
	return this.ProxyHttpServer.ConnectDial(network, addr)
}

// OnRequest 请求事件,Do完之后生效
func (this *Proxy) OnRequest(c ...Condition) *ReqAction {
	cs := make([]goproxy.ReqCondition, len(c))
	for i, v := range c {
		switch x := v.(type) {
		case goproxy.ReqCondition:
			cs[i] = x
		}
	}
	return &ReqAction{
		ReqProxyConds: this.ProxyHttpServer.OnRequest(cs...),
		log:           this.log,
	}
}

// OnResponse 响应事件,Do完之后生效
func (this *Proxy) OnResponse(c ...Condition) *RespAction {
	return &RespAction{
		ProxyConds: this.ProxyHttpServer.OnResponse(c...),
		log:        this.log,
	}
}
