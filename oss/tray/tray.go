package tray

import (
	"github.com/getlantern/systray"
	"github.com/injoyai/base/safe"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/goutil/oss/win"
	"path/filepath"
)

type Option func(s *Stray)

// WithLabel 新增Label菜单
func WithLabel(name string) Option {
	return func(s *Stray) {
		s.AddMenu().SetName(name).Disable()
	}
}

// WithStartup 添加自启菜单
func WithStartup() Option {
	return func(s *Stray) {
		s.AddMenuCheck().SetName("自启").OnClick(func(m *Menu) {
			filename := oss.ExecName()
			if !m.Checked() {
				win.CreateStartupShortcut(filename)
			} else {
				oss.Remove(oss.UserStartupDir(filepath.Base(filename) + ".link"))
			}
		})
	}
}

// WithShow 添加显示GUI
func WithShow(f func(m *Menu)) Option {
	return func(s *Stray) {
		s.AddMenu().SetName("显示").OnClick(f)
	}
}

// WithSeparator 添加横线
func WithSeparator() Option {
	return func(ui *Stray) {
		ui.AddSeparator()
	}
}

// WithExit 添加退出菜单
func WithExit() Option {
	return func(s *Stray) {
		s.AddMenu().
			SetName("退出").
			OnClick(func(m *Menu) {
				s.Close()
			})
	}
}

// WithHint 修改提示信息
func WithHint(hint string) Option {
	return func(s *Stray) {
		s.SetHint(hint)
	}
}

func Run(op ...Option) <-chan struct{} {
	s := &Stray{
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
			s.SetIco(DefaultIcon)
			for _, v := range op {
				v(s)
			}
		},
		func() { s.Closer.Close() },
	)
	return s.Closer.Done()
}

type Stray struct {
	*safe.Closer
	OnClose func()
}

// SetIco 设置图标
func (this *Stray) SetIco(icon []byte) *Stray {
	systray.SetIcon(icon)
	return this
}

// SetHint 设置提示
func (this *Stray) SetHint(hint string) *Stray {
	systray.SetTooltip(hint)
	return this
}

// AddSeparator 添加分割线
func (this *Stray) AddSeparator() {
	systray.AddSeparator()
}

// AddMenu 添加普通菜单
func (this *Stray) AddMenu() *Menu {
	return NewMenu()
}

// AddMenuCheck 添加选择菜单
func (this *Stray) AddMenuCheck() *MenuCheck {
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
