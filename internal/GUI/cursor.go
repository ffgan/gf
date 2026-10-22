package gui

import "runtime"

func GetCursor() string {
	if runtime.GOOS == "windows" {
		return windowsCursor()
	}
	ctx := StyleContext{
		Name:      "gtk-cursor-theme-name",
		GSettings: "cursor-theme",
		GConf:     "cursor_theme",
		XfConf:    "/Gtk/CursorThemeName",
		Kde:       "cursorTheme",
	}
	return getStyle(ctx)
}
