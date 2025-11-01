package gui

func getFont() string {
	ctx := StyleContext{
		Name:      "gtk-font-name",
		GSettings: "font-name",
		GConf:     "font_theme",
		XfConf:    "/Gtk/FontName",
		Kde:       "font",
		LxQt:      "font",
		Qt5ct:     "general",
		Fly:       "DefaultFont",
	}
	return getStyle(ctx)
}
