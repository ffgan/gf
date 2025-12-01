package gui

import (
	"os"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

type StyleContext struct {
	Name      string
	GSettings string
	GConf     string
	XfConf    string
	Kde       string
	LxQt      string
	Qt5ct     string
	Fly       string
}

func getStyle(ctx StyleContext) string {
	var gtk2, gtk3, qt, kde string

	display := os.Getenv("DISPLAY")
	if display == "" {
		return ""
	}

	de := detectDesktopEnvironment()
	switch {
	case strings.Contains(de, "KDE"), strings.Contains(de, "Plasma"):
		kde = readKDETheme(ctx.Kde)
	case strings.Contains(de, "Gnome"), strings.Contains(de, "Unity"), strings.Contains(de, "Budgie"):
		gtk3 = getGSettings("org.gnome.desktop.interface", ctx.GSettings)
		gtk2 = gtk3
	case strings.Contains(de, "Cinnamon"):
		gtk3 = getGSettings("org.cinnamon.desktop.interface", ctx.GSettings)
		gtk2 = gtk3
	case strings.Contains(de, "Mate"):
		gtk3 = getGSettings("org.mate.interface", ctx.GSettings)
		gtk2 = gtk3
	case strings.Contains(de, "Xfce"):
		gtk2 = utils.RunCommand("xfconf-query", "-c", "xsettings", "-p", ctx.XfConf)
	case strings.Contains(de, "LXQt"):
		qt = readLXQtTheme(ctx.LxQt)
	case strings.Contains(de, "Fly"):
		qt = readFlyTheme(ctx.Fly)
	}

	// GTK fallback 检测
	if gtk2 == "" {
		gtk2 = readGtk2Theme(ctx.Name)
	}
	if gtk3 == "" {
		gtk3 = readGtk3Theme(ctx.Name, ctx.GSettings)
	}

	// Qt5ct 特殊处理
	if qt == "" {
		if os.Getenv("QT_QPA_PLATFORMTHEME") == "qt5ct" {
			qt = readQt5ctTheme(ctx.Qt5ct)
		}
	}

	// 拼接结果
	theme := buildThemeString(map[string]string{
		"KDE":  kde,
		"Qt":   qt,
		"GTK2": gtk2,
		"GTK3": gtk3,
	})
	return theme
}
