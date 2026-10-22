// Package modules is the fastfetch module list and the detectors that
// back the names gf did not already implement as neofetch info lines.
package modules

import (
	"fmt"
	"strings"
)

// Mod is one fastfetch module. Key is the gf info function when the
// fastfetch name does not match it. An empty label prints the value
// alone, which is how title, separator, break, and the color blocks work.
type Mod struct {
	Name  string
	Key   string
	Label string
}

func (m Mod) line() []string {
	if m.Label == "" {
		return []string{m.Name}
	}
	return []string{m.Label, m.Name}
}

func (m Mod) id() string {
	if m.Key != "" {
		return m.Key
	}
	return m.Name
}

// fastfetchModules follows src/modules in fastfetch. Names match the
// module type a fastfetch config would use.
var fastfetchModules = []Mod{
	{Name: "title", Label: ""},
	{Name: "separator", Label: ""},
	{Name: "os", Key: "distro", Label: "OS"},
	{Name: "host", Key: "model", Label: "Host"},
	{Name: "bios", Label: "BIOS"},
	{Name: "board", Label: "Board"},
	{Name: "chassis", Label: "Chassis"},
	{Name: "bootmgr", Label: "Boot Manager"},
	{Name: "kernel", Label: "Kernel"},
	{Name: "initsystem", Label: "Init System"},
	{Name: "uptime", Label: "Uptime"},
	{Name: "loadavg", Label: "Load Average"},
	{Name: "processes", Label: "Processes"},
	{Name: "packages", Label: "Packages"},
	{Name: "shell", Label: "Shell"},
	{Name: "editor", Label: "Editor"},
	{Name: "display", Key: "resolution", Label: "Display"},
	{Name: "monitor", Label: "Monitor"},
	{Name: "brightness", Label: "Brightness"},
	{Name: "de", Label: "Desktop Environment"},
	{Name: "wm", Label: "Window Manager"},
	{Name: "wmtheme", Key: "wm_theme", Label: "WM Theme"},
	{Name: "theme", Label: "Theme"},
	{Name: "icons", Label: "Icons"},
	{Name: "font", Label: "Font"},
	{Name: "cursor", Label: "Cursor"},
	{Name: "wallpaper", Label: "Wallpaper"},
	{Name: "terminal", Key: "term", Label: "Terminal"},
	{Name: "terminalfont", Key: "term_font", Label: "Terminal Font"},
	{Name: "terminalsize", Key: "term_size", Label: "Terminal Size"},
	{Name: "terminaltheme", Label: "Terminal Theme"},
	{Name: "cpu", Label: "CPU"},
	{Name: "cpucache", Label: "CPU Cache"},
	{Name: "cpuusage", Label: "CPU Usage"},
	{Name: "gpu", Label: "GPU"},
	{Name: "opengl", Label: "OpenGL"},
	{Name: "opencl", Label: "OpenCL"},
	{Name: "vulkan", Label: "Vulkan"},
	{Name: "codec", Label: "Codec"},
	{Name: "memory", Label: "Memory"},
	{Name: "physicalmemory", Label: "Physical Memory"},
	{Name: "swap", Label: "Swap"},
	{Name: "disk", Label: "Disk"},
	{Name: "diskio", Label: "Disk I/O"},
	{Name: "physicaldisk", Label: "Physical Disk"},
	{Name: "btrfs", Label: "BTRFS"},
	{Name: "zpool", Label: "Zpool"},
	{Name: "battery", Label: "Battery"},
	{Name: "poweradapter", Key: "power_adapter", Label: "Power Adapter"},
	{Name: "lm", Label: "Login Manager"},
	{Name: "locale", Label: "Locale"},
	{Name: "datetime", Label: "Date & Time"},
	{Name: "localip", Key: "local_ip", Label: "Local IP"},
	{Name: "publicip", Key: "public_ip", Label: "Public IP"},
	{Name: "dns", Label: "DNS"},
	{Name: "wifi", Label: "Wi-Fi"},
	{Name: "netio", Label: "Network I/O"},
	{Name: "bluetooth", Label: "Bluetooth"},
	{Name: "bluetoothradio", Label: "Bluetooth Radio"},
	{Name: "sound", Label: "Sound"},
	{Name: "media", Key: "song", Label: "Media"},
	{Name: "player", Label: "Media Player"},
	{Name: "camera", Label: "Camera"},
	{Name: "gamepad", Label: "Gamepad"},
	{Name: "keyboard", Label: "Keyboard"},
	{Name: "mouse", Label: "Mouse"},
	{Name: "tpm", Label: "TPM"},
	{Name: "weather", Label: "Weather"},
	{Name: "users", Label: "Users"},
	{Name: "top", Label: "Top Processes"},
	{Name: "libc", Label: "Libc"},
	{Name: "version", Label: "Version"},
	{Name: "logo", Label: "Logo"},
	{Name: "custom", Label: "Custom"},
	{Name: "command", Label: "Command"},
	{Name: "break", Key: "line_break", Label: ""},
	{Name: "colors", Key: "cols", Label: ""},
}

// extraModules are gf info lines that fastfetch does not ship.
var extraModules = []Mod{
	{Name: "npu", Label: "NPU"},
	{Name: "gpu_driver", Label: "GPU Driver"},
	{Name: "java_ver", Label: "Java"},
	{Name: "python_ver", Label: "Python"},
	{Name: "node_ver", Label: "Node"},
	{Name: "window_size", Label: "Window Size"},
	{Name: "network", Label: "Network"},
}

var (
	byName map[string]Mod
	order  []Mod
)

func init() {
	order = append(append([]Mod{}, fastfetchModules...), extraModules...)
	byName = make(map[string]Mod, len(order)*2)
	for _, m := range order {
		byName[normalize(m.Name)] = m
		if m.Key != "" {
			byName[normalize(m.Key)] = m
		}
	}
}

func normalize(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, "-", "")
	return name
}

// Canonical maps a fastfetch module name onto the gf info key.
func Canonical(name string) string {
	if m, ok := byName[normalize(name)]; ok {
		return m.id()
	}
	return normalize(name)
}

// Names lists every module --structure and print_info accept.
func Names() []string {
	out := make([]string, len(order))
	for i, m := range order {
		out[i] = m.Name
	}
	return out
}

// FastfetchNames is the module set copied from fastfetch, without gf-only lines.
func FastfetchNames() []string {
	out := make([]string, len(fastfetchModules))
	for i, m := range fastfetchModules {
		out[i] = m.Name
	}
	return out
}

// Structure builds a print_info list. "all" is every module. Otherwise
// the spec is a fastfetch-style list separated by colons, commas, or spaces.
func Structure(spec string) ([][]string, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, fmt.Errorf("empty structure")
	}
	if strings.EqualFold(spec, "all") {
		lines := make([][]string, len(order))
		for i, m := range order {
			lines[i] = m.line()
		}
		return lines, nil
	}
	parts := strings.FieldsFunc(spec, func(r rune) bool {
		return r == ':' || r == ',' || r == ' ' || r == '\n'
	})
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty structure")
	}
	lines := make([][]string, 0, len(parts))
	for _, part := range parts {
		m, ok := byName[normalize(part)]
		if !ok {
			return nil, fmt.Errorf("unknown module %q", part)
		}
		lines = append(lines, m.line())
	}
	return lines, nil
}
