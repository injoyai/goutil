package tray

func Easy(name string, ico []byte, onSetting func(m *Menu), op ...Option) error {
	return Run(
		WithIco(ico),
		WithHint(name),
		func(s *Tray) {
			s.AddMenu().SetName("配置").SetIcon(IconSetting).OnClick(onSetting)
			for _, v := range op {
				v(s)
			}
		},
		WithStartup(),
		WithSeparator(),
		WithExit(),
	)
}
