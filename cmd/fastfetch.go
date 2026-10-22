package cmd

import (
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/ffgan/gf/internal/logo"
	"github.com/ffgan/gf/internal/modules"
)

func (h *hardware) remember(key, value string) {
	h.extraMu.Lock()
	defer h.extraMu.Unlock()
	if h.extra == nil {
		h.extra = map[string]string{}
	}
	h.extra[key] = value
}

func (h *hardware) recall(key string) (string, bool) {
	h.extraMu.Lock()
	defer h.extraMu.Unlock()
	value, ok := h.extra[key]
	return value, ok
}

// fastfetchDetectors are the modules that do not already have a Get_ method.
// Keys are the canonical info names.
var fastfetchDetectors = map[string]func(*hardware) string{
	"board":          func(*hardware) string { return modules.Board() },
	"chassis":        func(*hardware) string { return modules.Chassis() },
	"bootmgr":        func(*hardware) string { return modules.Bootmgr() },
	"brightness":     func(*hardware) string { return modules.Brightness() },
	"btrfs":          func(*hardware) string { return modules.Btrfs() },
	"camera":         func(*hardware) string { return modules.Camera() },
	"bluetoothradio": func(*hardware) string { return modules.BluetoothRadio() },
	"codec":          func(*hardware) string { return modules.Codec() },
	"cpucache":       func(*hardware) string { return modules.CPUCache() },
	"cpuusage":       func(*hardware) string { return modules.CPUUsage() },
	"datetime":       func(*hardware) string { return modules.DateTime() },
	"diskio":         func(*hardware) string { return modules.DiskIO() },
	"dns":            func(*hardware) string { return modules.DNS() },
	"gamepad":        func(*hardware) string { return modules.Gamepad() },
	"initsystem":     func(*hardware) string { return modules.InitSystem() },
	"keyboard":       func(*hardware) string { return modules.Keyboard() },
	"libc":           func(*hardware) string { return modules.Libc() },
	"lm":             func(*hardware) string { return modules.LoginManager() },
	"loadavg":        func(*hardware) string { return modules.Loadavg() },
	"monitor":        func(*hardware) string { return modules.Monitor() },
	"mouse":          func(*hardware) string { return modules.Mouse() },
	"netio":          func(*hardware) string { return modules.NetIO() },
	"opencl":         func(*hardware) string { return modules.OpenCL() },
	"opengl":         func(*hardware) string { return modules.OpenGL() },
	"vulkan":         func(*hardware) string { return modules.Vulkan() },
	"physicaldisk":   func(*hardware) string { return modules.PhysicalDisk() },
	"physicalmemory": func(*hardware) string { return modules.PhysicalMemory() },
	"player":         func(*hardware) string { return modules.Player() },
	"processes":      func(*hardware) string { return modules.Processes() },
	"sound":          func(*hardware) string { return modules.Sound() },
	"terminaltheme":  func(*hardware) string { return modules.TerminalTheme() },
	"top":            func(*hardware) string { return modules.Top() },
	"tpm":            func(*hardware) string { return modules.TPM() },
	"wallpaper":      func(*hardware) string { return modules.Wallpaper() },
	"weather":        func(h *hardware) string { return modules.Weather(h.WeatherLocation, h.WeatherFormat) },
	"wifi":           func(*hardware) string { return modules.Wifi() },
	"zpool":          func(*hardware) string { return modules.Zpool() },
	"version":        func(*hardware) string { return "gf " + Version },
	"custom":         func(h *hardware) string { return strings.TrimSpace(h.CustomText) },
	"command":        func(h *hardware) string { return modules.RunUserCommand(h.CommandText) },
	"logo":           func(h *hardware) string { return h.ascii_distro },
	"separator":      separatorLine,
}

func separatorLine(h *hardware) string {
	if h.length == 0 {
		h.Get_title()
	}
	n := h.length
	if n <= 0 {
		n = 16
	}
	ch := h.SeparatorString
	if ch == "" {
		ch = h.UnderlineChar
	}
	if ch == "" {
		ch = "-"
	}
	line := strings.Repeat(ch, n)
	if h.plainMode() {
		return line
	}
	return logo.UnderlineColor() + line + logo.ResetCode()
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func (h *hardware) writeJSON(w io.Writer) {
	type row struct {
		Key   string `json:"key"`
		Label string `json:"label,omitempty"`
		Value string `json:"value"`
	}
	rows := []row{}
	for _, entry := range h.entries() {
		if entry.blank {
			continue
		}
		rows = append(rows, row{
			Key:   entry.key,
			Label: entry.label,
			Value: ansiPattern.ReplaceAllString(entry.value, ""),
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(map[string]any{
		"version": Version,
		"result":  rows,
	})
}
