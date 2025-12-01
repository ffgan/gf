package gui

import (
	"os"
	"strings"
)

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

func detectDesktopEnvironment() string {
	// 优先使用 XDG_CURRENT_DESKTOP
	if de := os.Getenv("XDG_CURRENT_DESKTOP"); de != "" {
		return de
	}
	// Fallback: 检查进程名
	out := runCmd("ps", "-e")
	switch {
	case strings.Contains(out, "plasmashell"):
		return "KDE"
	case strings.Contains(out, "gnome-shell"):
		return "Gnome"
	case strings.Contains(out, "cinnamon"):
		return "Cinnamon"
	case strings.Contains(out, "xfce"):
		return "Xfce"
	case strings.Contains(out, "lxqt"):
		return "LXQt"
	}
	return "Unknown"
}
