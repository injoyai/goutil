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
	"github.com/tebeka/selenium/firefox"
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
)

// New
// 新建实例需要下载chromedriver
// 查看浏览器版本Chrome://version
// http://chromedriver.storage.googleapis.com/index.html
// https://www.chromedownloads.net/chrome64win/
func New(browserPath, driverPath string, option ...Option) *Entity {
	e := &Entity{
		browser:       Chrome,
		browserPath:   browserPath,
		driverPath:    driverPath,
		seleniumPort:  20165,
		seleniumDebug: false,
		retry:         1,
		Prefs:         map[string]interface{}{},
	}
	e.ShowWindow(oss.IsWindows())
	e.ShowImg(true)
	e.SetUserAgent(http.UserAgentDefault)
	for _, v := range option {
		v(e)
	}
	return e
}

type Option func(e *Entity)

type Entity struct {
	browser       Browser //浏览器
	browserPath   string  //浏览器目录
	driverPath    string  //chromedriver路径
	seleniumPort  int     //selenium端口
	seleniumDebug bool    //selenium调试模式
	retry         uint    //重试次数

	Prefs map[string]interface{}
	Args  []string
}

func (this *Entity) SetProxy(u string) *Entity {
	this.Args = append(this.Args, "--proxy-server="+u)
	return this
}

// ShowWindow 显示窗口linux系统无效
func (this *Entity) ShowWindow(b ...bool) *Entity {
	if !oss.IsWindows() || (len(b) > 0 && !b[0]) {
		this.Args = append(this.Args, "--headless")
	} else {
		for i, v := range this.Args {
			if v == "--headless" {
				this.Args = append(this.Args[:i], this.Args[i+1:]...)
				break
			}
		}
	}
	return this
}

// ShowImg 是否加载图片
func (this *Entity) ShowImg(b ...bool) *Entity {
	show := oss.IsWindows() && !(len(b) > 0 && !b[0])
	this.Prefs["profile.managed_default_content_settings.images"] = conv.SelectInt(show, 1, 2)
	return this
}

// SetRetry 设置重试次数
func (this *Entity) SetRetry(n uint) *Entity {
	this.retry = n
	return this
}

// SetBrowser 设置浏览器,目前只支持chrome
func (this *Entity) SetBrowser(b Browser) *Entity {
	this.browser = b
	return this
}

// SetBrowserPath 设置浏览器目录
func (this *Entity) SetBrowserPath(p string) *Entity {
	this.browserPath = p
	return this
}

// SetUserAgent 设置UserAgent
func (this *Entity) SetUserAgent(ua string) *Entity {
	this.Args = append(this.Args, "--user-agent="+ua)
	return this
}

// SetUserAgentDefault 设置UserAgent到默认值
func (this *Entity) SetUserAgentDefault() *Entity {
	return this.SetUserAgent(http.UserAgentDefault)
}

// SetUserAgentRand 设置随机UserAgent
func (this *Entity) SetUserAgentRand() *Entity {
	idx := g.RandInt(0, len(http.UserAgentList)-1)
	return this.SetUserAgent(http.UserAgentList[idx])
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
		selenium.Output(logs.DefaultErr),
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
	caps := selenium.Capabilities{"browserName": string(Chrome)}
	switch this.browser {
	case Chrome:
		caps.AddChrome(chrome.Capabilities{
			Path:  this.browserPath,
			Prefs: this.Prefs,
			Args:  this.Args,
		})
	case Firefox:
		caps.AddFirefox(firefox.Capabilities{
			Binary: this.browserPath,
			Prefs:  this.Prefs,
			Args:   this.Args,
		})
	}

	// 调起浏览器
	web, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", this.seleniumPort))
	if err != nil {
		return err
	}
	defer web.Close()

	return g.Retry(func() error { return f(&WebDriver{web}) }, this.retry)
}

/*

–user-data-dir=”[PATH]” 指定用户文件夹User Data路径，可以把书签这样的用户数据保存在系统分区以外的分区。
–disk-cache-dir=”[PATH]“ 指定缓存Cache路径
–disk-cache-size= 指定Cache大小，单位Byte
–first run 重置到初始状态，第一次运行
–incognito 隐身模式启动
–disable-javascript 禁用Javascript
--omnibox-popup-count="num" 将地址栏弹出的提示菜单数量改为num个。我都改为15个了。
--user-agent="xxxxxxxx" 修改HTTP请求头部的Agent字符串，可以通过about:version页面查看修改效果
--disable-plugins 禁止加载所有插件，可以增加速度。可以通过about:plugins页面查看效果
--disable-javascript 禁用JavaScript，如果觉得速度慢在加上这个
--disable-java 禁用java
--start-maximized 启动就最大化
--no-sandbox 取消沙盒模式
--single-process 单进程运行
--process-per-tab 每个标签使用单独进程
--process-per-site 每个站点使用单独进程
--in-process-plugins 插件不启用单独进程
--disable-popup-blocking 禁用弹出拦截
--disable-plugins 禁用插件
--disable-images 禁用图像
--incognito 启动进入隐身模式
--enable-udd-profiles 启用账户切换菜单
--proxy-pac-url 使用pac代理 [via 1/2]
--lang=zh-CN 设置语言为简体中文
--disk-cache-dir 自定义缓存目录
--disk-cache-size 自定义缓存最大值（单位byte）
--media-cache-size 自定义多媒体缓存最大值（单位byte）
--bookmark-menu 在工具 栏增加一个书签按钮
--enable-sync 启用书签同步

*/
