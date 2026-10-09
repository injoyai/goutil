//go:build !cgo
// +build !cgo

package tray

import (
	"errors"
	"fmt"

	"github.com/injoyai/base/safe"
)

// 非CGO环境兜底实现: 系统托盘依赖CGO(getlantern/systray), 此文件仅保证编译通过

type MenuItem struct {
	ClickedCh chan struct{}
}

func (this *MenuItem) SetTitle(title string)     {}
func (this *MenuItem) SetTooltip(tooltip string) {}
func (this *MenuItem) SetIcon(icon []byte)       {}
func (this *MenuItem) Hide()                     {}
func (this *MenuItem) Show()                     {}
func (this *MenuItem) Enable()                   {}
func (this *MenuItem) Disable()                  {}
func (this *MenuItem) Check()                    {}
func (this *MenuItem) Uncheck()                  {}
func (this *MenuItem) Checked() bool             { return false }
func (this *MenuItem) AddSubMenuItem(title, tooltip string) *MenuItem {
	return &MenuItem{ClickedCh: make(chan struct{})}
}

type (
	Option     func(s *Tray)
	OptionMenu func(m *Menu)
)

func Run(op ...Option) error {
	return errors.New("系统托盘需要启用CGO")
}

type Tray struct {
	*safe.Closer
	OnClose func()
}

func (this *Tray) SetIco(icon []byte) *Tray  { return this }
func (this *Tray) SetHint(hint string) *Tray { return this }
func (this *Tray) SetHintf(format string, args ...any) *Tray {
	return this.SetHint(fmt.Sprintf(format, args...))
}
func (this *Tray) AddSeparator()            {}
func (this *Tray) AddMenu() *Menu           { return NewMenu() }
func (this *Tray) AddMenuCheck() *MenuCheck { return NewMenuCheck() }

type Menu struct {
	*MenuItem
	*safe.Closer
	onClick func(m *Menu)
}

func (this *Menu) SetOptions(op ...OptionMenu) *Menu {
	for _, v := range op {
		v(this)
	}
	return this
}
func (this *Menu) OnClick(fn func(m *Menu)) *Menu { this.onClick = fn; return this }
func (this *Menu) AddMenu() *Menu                 { return NewMenu() }
func (this *Menu) SetName(name string) *Menu      { return this }
func (this *Menu) SetHint(hint string) *Menu      { return this }
func (this *Menu) SetIcon(icon []byte) *Menu      { return this }
func (this *Menu) Hide() *Menu                    { return this }
func (this *Menu) Show() *Menu                    { return this }
func (this *Menu) Enable() *Menu                  { return this }
func (this *Menu) Disable() *Menu                 { return this }

type MenuCheck struct{ *Menu }

func (this *MenuCheck) GetChecked() bool                   { return this.MenuItem.Checked() }
func (this *MenuCheck) SetChecked(checked bool) *MenuCheck { return this }

func NewMenuCheck() *MenuCheck { return &MenuCheck{Menu: newMenu()} }
func NewMenu() *Menu           { return newMenu() }

func newMenu() *Menu {
	return &Menu{
		MenuItem: &MenuItem{ClickedCh: make(chan struct{})},
		Closer:   safe.NewCloser(),
	}
}

func WithLabel(name string, op ...OptionMenu) Option {
	return func(s *Tray) { s.AddMenu().SetName(name).Disable().SetOptions(op...) }
}
func WithShell(name string, cmd string, op ...OptionMenu) Option {
	return func(s *Tray) { s.AddMenu().SetName(name).OnClick(func(m *Menu) {}).SetOptions(op...) }
}
func WithStartup(op ...OptionMenu) Option {
	return func(s *Tray) { s.AddMenuCheck().SetChecked(false).SetName("自启").OnClick(func(m *Menu) {}).SetOptions(op...) }
}
func WithButton(name string, f func(m *Menu), op ...OptionMenu) Option {
	return func(s *Tray) { s.AddMenu().SetName(name).OnClick(f).SetOptions(op...) }
}
func WithShow(f func(m *Menu), op ...OptionMenu) Option {
	return func(s *Tray) { s.AddMenu().SetName("显示").SetIcon(IconShow).OnClick(f).SetOptions(op...) }
}
func WithSetting(f func(m *Menu), op ...OptionMenu) Option {
	return func(s *Tray) { s.AddMenu().SetName("配置").SetIcon(IconSetting).OnClick(f).SetOptions(op...) }
}
func WithSeparator() Option {
	return func(s *Tray) { s.AddSeparator() }
}
func WithExit(op ...OptionMenu) Option {
	return func(s *Tray) { s.AddMenu().SetName("退出").SetIcon(IconExit).OnClick(func(m *Menu) {}).SetOptions(op...) }
}
func WithIco(ico []byte) Option {
	return func(s *Tray) { s.SetIco(ico) }
}
func WithHint(hint string) Option {
	return func(s *Tray) { s.SetHint(hint) }
}
