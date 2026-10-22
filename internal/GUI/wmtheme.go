package gui

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

// GetWMTheme returns the window manager's own theme, separate from the GTK theme.
func GetWMTheme(wm string) string {
	if runtime.GOOS == "windows" {
		return windowsWMTheme()
	}
	low := strings.ToLower(wm)
	switch {
	case strings.Contains(low, "xfwm"):
		return utils.RunCommand("xfconf-query", "-c", "xfwm4", "-p", "/general/theme")
	case strings.Contains(low, "mutter"), strings.Contains(low, "gnome"):
		return strings.Trim(utils.RunCommand("gsettings", "get", "org.gnome.desktop.wm.preferences", "theme"), "'\"")
	case strings.Contains(low, "openbox"):
		return openboxTheme()
	case strings.Contains(low, "kwin"), strings.Contains(low, "plasma"):
		return readKDETheme("theme")
	}
	if utils.CommandExists("xfconf-query") && strings.Contains(utils.RunCommand("ps", "-e", "-o", "comm="), "xfwm4") {
		return utils.RunCommand("xfconf-query", "-c", "xfwm4", "-p", "/general/theme")
	}
	return ""
}

func openboxTheme() string {
	path := filepath.Join(os.Getenv("HOME"), ".config", "openbox", "rc.xml")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "<name>") && strings.Contains(line, "</name>") {
			line = strings.TrimPrefix(line, "<name>")
			line = strings.TrimSuffix(line, "</name>")
			return strings.TrimSpace(line)
		}
	}
	return ""
}

// GetTermFont reads the font configured by the detected terminal emulator.
func GetTermFont(term string) string {
	if font := windowsTermFont(term); font != "" {
		return font
	}
	low := strings.ToLower(term)
	switch {
	case strings.Contains(low, "xfce4-terminal"):
		return iniValue(filepath.Join(os.Getenv("HOME"), ".config", "xfce4", "terminal", "terminalrc"), "FontName")
	case strings.Contains(low, "kitty"):
		return kittyFont()
	case strings.Contains(low, "gnome-terminal"), strings.Contains(low, "gnome-console"):
		font := utils.RunCommand("gsettings", "get", "org.gnome.desktop.interface", "monospace-font-name")
		return strings.Trim(font, "'\"")
	default:
		return ""
	}
}

func kittyFont() string {
	path := filepath.Join(os.Getenv("HOME"), ".config", "kitty", "kitty.conf")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var family, size string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "font_family") {
			family = strings.TrimSpace(strings.TrimPrefix(line, "font_family"))
		}
		if strings.HasPrefix(line, "font_size") {
			size = strings.TrimSpace(strings.TrimPrefix(line, "font_size"))
		}
	}
	return strings.TrimSpace(family + " " + size)
}

func iniValue(path, key string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key) {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
