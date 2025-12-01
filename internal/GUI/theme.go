package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetTheme() string {
	ctx := StyleContext{
		Name:      "gtk-theme-name",
		GSettings: "gtk-theme",
		GConf:     "gtk_theme",
		XfConf:    "/Net/ThemeName",
		Kde:       "widgetStyle",
		LxQt:      "style",
		Qt5ct:     "style=",
		Fly:       "ColorScheme",
	}
	return getStyle(ctx)
}

func readKDETheme(key string) string {
	path := filepath.Join(os.Getenv("HOME"), ".config", "kdeglobals")
	if !utils.FileExists(path) {
		return ""
	}
	b, _ := os.ReadFile(path)
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		if strings.HasPrefix(l, key) {
			return strings.TrimSpace(strings.SplitN(l, "=", 2)[1])
		}
	}
	return ""
}

func readGtk2Theme(name string) string {
	paths := []string{
		filepath.Join(os.Getenv("HOME"), ".gtkrc-2.0"),
		"/etc/gtk-2.0/gtkrc",
		"/usr/share/gtk-2.0/gtkrc",
	}
	for _, p := range paths {
		if !utils.FileExists(p) {
			continue
		}
		data, _ := os.ReadFile(p)
		for _, line := range strings.Split(string(data), "\n") {
			if strings.Contains(line, name) {
				return strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
			}
		}
	}
	return ""
}

func readGtk3Theme(name, key string) string {
	conf := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "gtk-3.0", "settings.ini")
	if conf == "" {
		conf = filepath.Join(os.Getenv("HOME"), ".config", "gtk-3.0", "settings.ini")
	}
	if utils.FileExists(conf) {
		data, _ := os.ReadFile(conf)
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, name) {
				return strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
			}
		}
	}
	// gsettings fallback
	return getGSettings("org.gnome.desktop.interface", key)
}

func readLXQtTheme(key string) string {
	path := filepath.Join(os.Getenv("HOME"), ".config", "lxqt", "lxqt.conf")
	if !utils.FileExists(path) {
		return ""
	}
	data, _ := os.ReadFile(path)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, key) {
			return strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
		}
	}
	return ""
}

func readFlyTheme(key string) string {
	path := filepath.Join(os.Getenv("HOME"), ".fly", "paletterc")
	if !utils.FileExists(path) {
		return ""
	}
	data, _ := os.ReadFile(path)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, key) {
			return filepath.Base(strings.TrimSpace(strings.SplitN(line, "=", 2)[1]))
		}
	}
	return ""
}

func readQt5ctTheme(key string) string {
	conf := filepath.Join(os.Getenv("HOME"), ".config", "qt5ct", "qt5ct.conf")
	if !utils.FileExists(conf) {
		return ""
	}
	data, _ := os.ReadFile(conf)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, key) {
			return strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
		}
	}
	return ""
}

func getGSettings(schema, key string) string {
	out := utils.RunCommand("gsettings", "get", schema, key)
	return strings.Trim(strings.TrimSpace(out), "'\"")
}

func buildThemeString(parts map[string]string) string {
	var res string
	for k, v := range parts {
		if v != "" {
			if res != "" {
				res += ", "
			}
			res += fmt.Sprintf("%s [%s]", v, k)
		}
	}
	return res
}
