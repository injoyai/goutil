package tray

import (
	"github.com/getlantern/systray"
	"github.com/injoyai/base/safe"
)

type Option func(ui *UI)

func WithSeparator() Option {
	return func(ui *UI) {
		ui.AddSeparator()
	}
}

func WithExit() Option {
	return func(ui *UI) {
		ui.AddMenu().
			SetName("退出").
			OnClick(func() {
				ui.Close()
			})
	}
}

func Run(op ...Option) <-chan struct{} {
	ui := &UI{
		Closer: safe.NewCloser(),
	}
	ui.Closer.SetCloseFunc(func(err error) error {
		if ui.OnClose != nil {
			ui.OnClose()
		}
		systray.Quit()
		return nil
	})
	systray.Run(
		func() {
			ui.SetHint("Go 程序")
			ui.SetIcon(DefaultIcon)
			for _, v := range op {
				v(ui)
			}
		},
		func() { ui.Closer.Close() },
	)
	return ui.Closer.Done()
}

type UI struct {
	*safe.Closer
	OnReady []func()
	OnClose func()
}

// SetIcon 设置图标
func (this *UI) SetIcon(icon []byte) *UI {
	systray.SetIcon(icon)
	return this
}

// SetHint 设置提示
func (this *UI) SetHint(hint string) *UI {
	systray.SetTooltip(hint)
	return this
}

// AddSeparator 添加分割线
func (this *UI) AddSeparator() {
	systray.AddSeparator()
}

// AddMenu 添加普通菜单
func (this *UI) AddMenu() *Menu {
	return NewMenu()
}

// AddMenuCheck 添加选择菜单
func (this *UI) AddMenuCheck() *MenuCheck {
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
	onClick func()
}

func (this *Menu) run() {
	for {
		select {
		case <-this.Closer.Done():
			return
		case <-this.MenuItem.ClickedCh:
			if this.onClick != nil {
				this.onClick()
			}
		}
	}
}

func (this *Menu) OnClick(fn func()) *Menu {
	this.onClick = fn
	return this
}

func (this *Menu) SetIcon(icon []byte) *Menu {
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
