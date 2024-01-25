package spider

import (
	"fmt"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/logs"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
)

const (
	ByID              = selenium.ByID
	ByXPATH           = selenium.ByXPATH
	ByLinkText        = selenium.ByLinkText
	ByPartialLinkText = selenium.ByPartialLinkText
	ByName            = selenium.ByName
	ByTagName         = selenium.ByTagName
	ByClassName       = selenium.ByClassName
	ByCSSSelector     = selenium.ByCSSSelector
)

type Browser string

const (
	Chrome  Browser = "chrome"
	Firefox Browser = "firefox"
	Opera   Browser = "opera"
)

// New
// 新建实例需要下载chromedriver
// 查看浏览器版本Chrome://version
// http://chromedriver.storage.googleapis.com/index.html
func New(browserPath, driverPath string, option ...Option) *Entity {
	e := &Entity{
		showWindow:    oss.IsWindows(),
		showImg:       true,
		browser:       Chrome,
		browserPath:   browserPath,
		driverPath:    driverPath,
		seleniumPort:  20165,
		seleniumDebug: true,
		userAgent:     http.UserAgentDefault,
		retry:         3,
	}
	for _, v := range option {
		v(e)
	}
	return e
}

type Option func(e *Entity)

type Entity struct {
	showWindow    bool    //显示窗口
	showImg       bool    //显示图片
	browser       Browser //浏览器
	browserPath   string  //浏览器目录
	driverPath    string  //chromedriver路径
	seleniumPort  int     //selenium端口
	seleniumDebug bool    //selenium调试模式
	userAgent     string  //User-Agent
	retry         uint    //重试次数
}

func (this *Entity) SetRetry(n uint) *Entity {
	this.retry = n
	return this
}

func (this *Entity) SetBrowser(b Browser) *Entity {
	this.browser = b
	return this
}

func (this *Entity) SetUserAgent(ua string) *Entity {
	this.userAgent = ua
	return this
}

func (this *Entity) SetUserAgentDefault() *Entity {
	return this.SetUserAgent(http.UserAgentDefault)
}

func (this *Entity) SetUserAgentRand() *Entity {
	idx := g.RandInt(0, len(http.UserAgentList)-1)
	return this.SetUserAgent(http.UserAgentList[idx])
}

// ShowWindow 显示窗口linux系统无效
func (this *Entity) ShowWindow(b ...bool) *Entity {
	this.showWindow = !(len(b) > 0 && !b[0])
	return this
}

// ShowImg 是否加载图片
func (this *Entity) ShowImg(b ...bool) *Entity {
	this.showImg = !(len(b) > 0 && !b[0])
	return this
}

// SetPort 设置端口
func (this *Entity) SetPort(port int) *Entity {
	this.seleniumPort = port
	return this
}

// Debug 是否打印日志
func (this *Entity) Debug(b ...bool) *Entity {
	this.seleniumDebug = len(b) == 0 || b[0]
	return this
}

// Run 执行,记得保留加载时间
func (this *Entity) Run(f func(w *WebDriver) error, option ...selenium.ServiceOption) error {

	selenium.SetDebug(this.seleniumDebug)
	serviceOption := []selenium.ServiceOption{
		selenium.Output(logs.DefaultErr), // Output debug information to STDERR.
	}
	serviceOption = append(serviceOption, option...)
	//新建seleniumServer
	service, err := selenium.NewChromeDriverService(
		this.driverPath,
		this.seleniumPort,
		serviceOption...,
	)
	if nil != err {
		return err
	}
	defer service.Stop()

	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{"browserName": Chrome}
	//设置浏览器参数
	caps.AddChrome(chrome.Capabilities{
		Path: this.browserPath,
		Prefs: map[string]interface{}{
			//是否禁止图片加载，加快渲染速度
			"profile.managed_default_content_settings.images": conv.SelectInt(this.showWindow && this.showImg, 1, 2),
		},
		Args: []string{
			"--user-agent=" + this.userAgent,
			conv.SelectString(!oss.IsWindows() || !this.showWindow, "--headless", ""),
		},
	})

	// 调起chrome浏览器
	web, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", this.seleniumPort))
	if err != nil {
		return err
	}
	defer web.Close()

	return g.Retry(func() error {
		return f(&WebDriver{web})
	}, int(this.retry))
}

/*



 */

type WebDriver struct {
	selenium.WebDriver
}

func (this *WebDriver) Wait(t time.Duration) *WebDriver {
	<-time.After(t)
	return this
}

func (this *WebDriver) WaitSec(n ...int) *WebDriver {
	return this.Wait(time.Duration(conv.GetDefaultInt(1, n...)) * time.Second)
}

func (this *WebDriver) WaitMin(n ...int) *WebDriver {
	return this.Wait(time.Duration(conv.GetDefaultInt(1, n...)) * time.Minute)
}

// Text 返回页面数据
func (this *WebDriver) Text() (string, error) {
	return this.PageSource()
}

// Open 打开网页
func (this *WebDriver) Open(url string) error {
	return this.Get(url)
}

// FindXPaths 查找所有XPath
func (this *WebDriver) FindXPaths(path string) ([]*Element, error) {
	es, err := this.FindElements(ByXPATH, path)
	if err != nil {
		return nil, err
	}
	list := []*Element(nil)
	for _, v := range es {
		list = append(list, &Element{v})
	}
	return list, nil
}

// FindXPath 查找所有XPath
func (this *WebDriver) FindXPath(path string) (*Element, error) {
	e, err := this.FindElement(ByXPATH, path)
	return &Element{e}, err
}

// FindSelects 查找所有Select
func (this *WebDriver) FindSelects(path string) ([]*Element, error) {
	es, err := this.FindElements(ByCSSSelector, path)
	if err != nil {
		return nil, err
	}
	list := []*Element(nil)
	for _, v := range es {
		list = append(list, &Element{v})
	}
	return list, nil
}

// FindSelect 查找所有Select
func (this *WebDriver) FindSelect(path string) (*Element, error) {
	e, err := this.FindElement(ByCSSSelector, path)
	return &Element{e}, err
}

type Element struct {
	selenium.WebElement
}

func (this *Element) Wait(t time.Duration) *Element {
	<-time.After(t)
	return this
}

func (this *Element) WaitSec(n ...int) *Element {
	return this.Wait(time.Duration(conv.GetDefaultInt(1, n...)) * time.Second)
}

func (this *Element) WaitMin(n ...int) *Element {
	return this.Wait(time.Duration(conv.GetDefaultInt(1, n...)) * time.Minute)
}

func (this *Element) Write(s string) error {
	return this.WebElement.SendKeys(s)
}
