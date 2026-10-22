package gui

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/ffgan/gf/internal/utils"
)

func GetTheme() string {
	if runtime.GOOS == "windows" {
		return windowsTheme()
	}
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
	var conf string
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		conf = filepath.Join(dir, "gtk-3.0", "settings.ini")
	} else {
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

var (
	gsettingsOnce  sync.Once
	gsettingsOK    bool
	gsettingsCache sync.Map
)

func getGSettings(schema, key string) string {
	gsettingsOnce.Do(func() {
		gsettingsOK = utils.CommandExists("gsettings")
	})
	if !gsettingsOK {
		return ""
	}
	cacheKey := schema + "\x00" + key
	if v, ok := gsettingsCache.Load(cacheKey); ok {
		return v.(string)
	}
	out := utils.RunCommand("gsettings", "get", schema, key)
	out = strings.Trim(strings.TrimSpace(out), "'\"")
	gsettingsCache.Store(cacheKey, out)
	return out
}

type xfProperty struct {
	Name  string       `xml:"name,attr"`
	Value string       `xml:"value,attr"`
	Props []xfProperty `xml:"property"`
}

type xfChannel struct {
	Props []xfProperty `xml:"property"`
}

var (
	xsettingsOnce sync.Once
	xsettings     *xfChannel
	xsettingsOK   bool
)

func loadXsettings() *xfChannel {
	xsettingsOnce.Do(func() {
		base := os.Getenv("XDG_CONFIG_HOME")
		if base == "" {
			base = filepath.Join(os.Getenv("HOME"), ".config")
		}
		path := filepath.Join(base, "xfce4", "xfconf", "xfce-perchannel-xml", "xsettings.xml")
		data, err := os.ReadFile(path)
		if err != nil {
			return
		}
		var ch xfChannel
		if xml.Unmarshal(data, &ch) != nil {
			return
		}
		xsettings = &ch
		xsettingsOK = true
	})
	if !xsettingsOK {
		return nil
	}
	return xsettings
}

func lookupXf(props []xfProperty, path []string) (string, bool) {
	if len(path) == 0 {
		return "", false
	}
	for _, p := range props {
		if p.Name != path[0] {
			continue
		}
		if len(path) == 1 {
			return p.Value, true
		}
		return lookupXf(p.Props, path[1:])
	}
	return "", false
}

// readXfconf reads xfce's xsettings channel from disk. xfconf-query talks
// to xfconfd over dbus; the xml file is the same data already on disk.
func readXfconf(prop string) string {
	if ch := loadXsettings(); ch != nil {
		parts := strings.Split(strings.Trim(prop, "/"), "/")
		if value, ok := lookupXf(ch.Props, parts); ok && value != "" {
			return value
		}
	}
	if prop == "" {
		return ""
	}
	return utils.RunCommand("xfconf-query", "-c", "xsettings", "-p", prop)
}

// themeOrder keeps GTK/Qt fragments stable. Map iteration order changes
// every process, so the same theme used to print in a different order
// on every run.
var themeOrder = []string{"GTK2", "GTK3", "Qt", "KDE"}

func buildThemeString(parts map[string]string) string {
	seen := map[string]bool{}
	var chunks []string
	for _, key := range themeOrder {
		value := strings.Trim(strings.TrimSpace(parts[key]), `"'`)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		chunks = append(chunks, fmt.Sprintf("%s [%s]", value, key))
	}
	return strings.Join(chunks, ", ")
}
