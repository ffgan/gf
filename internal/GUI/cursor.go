package gui

func GetCursor() string {
	ctx := StyleContext{
		Name:      "gtk-cursor-theme-name",
		GSettings: "cursor-theme",
		GConf:     "cursor_theme",
		XfConf:    "/Gtk/CursorThemeName",
		Kde:       "cursorTheme",
	}
	return getStyle(ctx)
}
