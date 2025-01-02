package tray

func Name(name string) OptionMenu {
	return func(m *Menu) {
		m.SetName(name)
	}
}

func Icon(ico []byte) OptionMenu {
	return func(m *Menu) {
		m.SetIcon(ico)
	}
}

func Enable(b ...bool) OptionMenu {
	return func(m *Menu) {
		if len(b) == 0 || b[0] {
			m.Enable()
		} else {
			m.Disable()
		}
	}
}

func Show(b ...bool) OptionMenu {
	return func(m *Menu) {
		if len(b) == 0 || b[0] {
			m.Show()
		} else {
			m.Hide()
		}
	}
}

func Hint(hint string) OptionMenu {
	return func(m *Menu) {
		m.SetHint(hint)
	}
}

func OnClick(fn func(m *Menu)) OptionMenu {
	return func(m *Menu) {
		m.OnClick(fn)
	}
}
