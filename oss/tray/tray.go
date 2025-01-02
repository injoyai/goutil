package tray

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/oss/shell"
	"github.com/injoyai/goutil/oss/win"
	"github.com/injoyai/logs"
	"path/filepath"
	"strings"
)

type (
	Option     func(s *Tray)
	OptionMenu func(m *Menu)
)

// WithLabel 新增Label菜单
func WithLabel(name string, op ...OptionMenu) Option {
	return func(s *Tray) {
		s.AddMenu().SetName(name).Disable().SetOptions(op...)
	}
}

// WithShell 执行脚本
func WithShell(name string, cmd string, op ...OptionMenu) Option {
	return func(s *Tray) {
		s.AddMenu().SetName(name).OnClick(func(m *Menu) {
			shell.Run(cmd)
		}).SetOptions(op...)
	}
}

// WithStartup 添加自启菜单
func WithStartup(op ...OptionMenu) Option {
	return func(s *Tray) {
		filename := oss.ExecName()
		_, name := filepath.Split(filename)
		name = strings.Split(name, ".")[0]
		startupFilename := oss.UserStartupDir(name + ".lnk")
		s.AddMenuCheck().
			SetChecked(oss.Exists(startupFilename)).
			SetName("自启").
			OnClick(func(m *Menu) {
				if !m.Checked() {
					logs.PrintErr(win.CreateStartupShortcut(filename))
					m.Check()
				} else {
					logs.PrintErr(oss.Remove(startupFilename))
					m.Uncheck()
				}
			}).
			SetOptions(op...)
	}
}

// WithShow 添加显示GUI
func WithShow(f func(m *Menu), op ...OptionMenu) Option {
	return func(s *Tray) {
		s.AddMenu().SetName("显示").OnClick(f).SetOptions(op...)
	}
}

// WithSeparator 添加横线
func WithSeparator() Option {
	return func(ui *Tray) {
		ui.AddSeparator()
	}
}

// WithExit 添加退出菜单
func WithExit(op ...OptionMenu) Option {
	return func(s *Tray) {
		s.AddMenu().
			SetName("退出").
			SetIcon(IconExit).
			OnClick(func(m *Menu) {
				s.Close()
			}).
			SetOptions(op...)
	}
}

// WithIco 设置图标
func WithIco(ico []byte) Option {
	return func(s *Tray) {
		s.SetIco(ico)
	}
}

// WithHint 修改提示信息
func WithHint(hint string) Option {
	return func(s *Tray) {
		s.SetHint(hint)
	}
}

func Run(op ...Option) <-chan struct{} {
	s := &Tray{
		Closer: safe.NewCloser(),
	}
	s.Closer.SetCloseFunc(func(err error) error {
		if s.OnClose != nil {
			s.OnClose()
		}
		systray.Quit()
		return nil
	})
	systray.Run(
		func() {
			s.SetHint("Go 程序")
			s.SetIco(IconGo)
			for _, v := range op {
				v(s)
			}
		},
		func() { s.Closer.Close() },
	)
	return s.Closer.Done()
}

type Tray struct {
	*safe.Closer
	OnClose func()
}

// SetIco 设置图标
func (this *Tray) SetIco(icon []byte) *Tray {
	systray.SetIcon(icon)
	return this
}

// SetHint 设置提示
func (this *Tray) SetHint(hint string) *Tray {
	systray.SetTooltip(hint)
	return this
}

// SetHintf 设置提示,格式化
func (this *Tray) SetHintf(format string, args ...interface{}) *Tray {
	return this.SetHint(fmt.Sprintf(format, args...))
}

// AddSeparator 添加分割线
func (this *Tray) AddSeparator() {
	systray.AddSeparator()
}

// AddMenu 添加普通菜单
func (this *Tray) AddMenu() *Menu {
	return NewMenu()
}

// AddMenuCheck 添加选择菜单
func (this *Tray) AddMenuCheck() *MenuCheck {
	return NewMenuCheck()
}

type MenuCheck struct {
	*Menu
}

func (this *MenuCheck) GetChecked() bool {
	return this.MenuItem.Checked()
}

func (this *MenuCheck) SetChecked(checked bool) *MenuCheck {
	if checked {
		this.MenuItem.Check()
	} else {
		this.MenuItem.Uncheck()
	}
	return this
}

func NewMenuCheck() *MenuCheck {
	mu := systray.AddMenuItemCheckbox("", "", false)
	return &MenuCheck{Menu: newMenu(mu)}
}

func NewMenu() *Menu {
	mu := systray.AddMenuItem("", "")
	return newMenu(mu)
}

func newMenu(mu *systray.MenuItem) *Menu {
	m := &Menu{
		MenuItem: mu,
		Closer:   safe.NewCloser(),
	}
	go m.run()
	return m
}

type Menu struct {
	*systray.MenuItem
	*safe.Closer
	onClick func(m *Menu)
}

func (this *Menu) run() {
	for {
		select {
		case <-this.Closer.Done():
			return
		case <-this.MenuItem.ClickedCh:
			if this.onClick != nil {
				this.onClick(this)
			}
		}
	}
}

func (this *Menu) SetOptions(op ...OptionMenu) *Menu {
	for _, v := range op {
		v(this)
	}
	return this
}

func (this *Menu) OnClick(fn func(m *Menu)) *Menu {
	this.onClick = fn
	return this
}

func (this *Menu) SetIco(icon []byte) *Menu {
	this.MenuItem.SetIcon(icon)
	return this
}

func (this *Menu) AddMenu() *Menu {
	mu := this.MenuItem.AddSubMenuItem("", "")
	return newMenu(mu)
}

func (this *Menu) SetName(name string) *Menu {
	this.MenuItem.SetTitle(name)
	return this
}

func (this *Menu) SetHint(hint string) *Menu {
	this.MenuItem.SetTooltip(hint)
	return this
}

func (this *Menu) SetIcon(icon []byte) *Menu {
	this.MenuItem.SetIcon(icon)
	return this
}

func (this *Menu) Hide() *Menu {
	this.MenuItem.Hide()
	return this
}

func (this *Menu) Show() *Menu {
	this.MenuItem.Show()
	return this
}

func (this *Menu) Enable() *Menu {
	this.MenuItem.Enable()
	return this
}

func (this *Menu) Disable() *Menu {
	this.MenuItem.Disable()
	return this
}
