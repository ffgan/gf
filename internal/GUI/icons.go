package gui

func GetIcons() string {
	ctx := StyleContext{
		Name:      "gtk-icon-theme-name",
		GSettings: "icon-theme",
		GConf:     "icon_theme",
		XfConf:    "/Net/IconThemeName",
		Kde:       "Theme",
		LxQt:      "icon_theme",
		Qt5ct:     "icon_theme",
		Fly:       "IconTheme",
	}
	return getStyle(ctx)
}
