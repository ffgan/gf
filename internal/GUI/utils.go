package gui

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

var (
	deRun     bool
	de        string
	deVersion string
	wmRun     bool
	wm        string
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

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func commandOutput(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func runCmd(cmd string, args ...string) string {
	c := exec.Command(cmd, args...)
	var buf bytes.Buffer
	c.Stdout = &buf
	_ = c.Run()
	return buf.String()
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
		gtk2 = runCmd("xfconf-query", "-c", "xsettings", "-p", ctx.XfConf)
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
