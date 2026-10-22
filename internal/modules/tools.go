package modules

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ffgan/gf/internal/utils"
)

func commandExists(name string) bool {
	return utils.CommandExists(name)
}

func runOut(name string, args ...string) string {
	return utils.RunCommand(name, args...)
}

func runCombined(name string, args ...string) string {
	out, err := utils.CommandCombinedOutput(name, args...)
	if err != nil && len(out) == 0 {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// RunUserCommand runs the command module. An empty command detects nothing.
func RunUserCommand(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	if runtime.GOOS == "windows" {
		return runCombined("cmd", "/C", text)
	}
	return runCombined("sh", "-c", text)
}

// Libc is the C library linked on this machine.
func Libc() string {
	getconf := ""
	if commandExists("getconf") {
		getconf = runOut("getconf", "GNU_LIBC_VERSION")
	}
	ldd := ""
	if commandExists("ldd") {
		ldd = runCombined("ldd", "--version")
	}
	return ParseLibc(getconf, ldd)
}

// LoginManager follows the display-manager unit when systemd provides one.
func LoginManager() string {
	target := sysPath("etc", "systemd", "system", "display-manager.service")
	if dest, err := os.Readlink(target); err == nil && dest != "" {
		base := dest
		if i := strings.LastIndexAny(dest, `/\`); i >= 0 {
			base = dest[i+1:]
		}
		base = strings.TrimSuffix(base, ".service")
		if base != "" {
			return base
		}
	}
	running := map[string]bool{}
	for _, proc := range readProcs() {
		running[proc.name] = true
	}
	for _, name := range []string{"sddm", "gdm", "gdm3", "lightdm", "greetd", "lxdm", "slim", "xdm"} {
		if running[name] {
			return name
		}
	}
	return ""
}

// Player is the active media player name.
func Player() string {
	if commandExists("playerctl") {
		list := runOut("playerctl", "-l")
		var lines []string
		for _, name := range strings.Split(list, "\n") {
			name = strings.TrimSpace(name)
			if name == "" || strings.Contains(name, "No players") {
				continue
			}
			status := strings.TrimSpace(runOut("playerctl", "-p", name, "status"))
			if status != "" {
				lines = append(lines, name+" ("+status+")")
			} else {
				lines = append(lines, name)
			}
		}
		if len(lines) > 0 {
			return joinLines(lines)
		}
	}
	known := map[string]string{
		"spotify": "Spotify", "vlc": "VLC", "mpv": "mpv",
		"rhythmbox": "Rhythmbox", "clementine": "Clementine",
		"strawberry": "Strawberry", "cmus": "cmus", "elisa": "Elisa",
	}
	var found []string
	for _, proc := range readProcs() {
		if label, ok := known[strings.ToLower(proc.name)]; ok {
			found = append(found, label)
		}
	}
	return joinLines(found)
}

// OpenGL asks glxinfo for the renderer and version.
func OpenGL() string {
	if !commandExists("glxinfo") {
		return ""
	}
	text := runCombined("glxinfo", "-B")
	renderer, version := "", ""
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if v, ok := strings.CutPrefix(line, "OpenGL renderer string:"); ok {
			renderer = strings.TrimSpace(v)
		}
		if v, ok := strings.CutPrefix(line, "OpenGL version string:"); ok {
			version = strings.TrimSpace(v)
		}
	}
	switch {
	case renderer != "" && version != "":
		return renderer + " (" + version + ")"
	default:
		return renderer
	}
}

// Vulkan asks vulkaninfo for each device name.
func Vulkan() string {
	if !commandExists("vulkaninfo") {
		return ""
	}
	text := runCombined("vulkaninfo", "--summary")
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(line)
		if name, ok := strings.CutPrefix(trim, "deviceName"); ok {
			name = strings.TrimSpace(strings.TrimLeft(name, "= "))
			if name != "" {
				lines = append(lines, name)
			}
		}
	}
	return joinLines(lines)
}

// OpenCL asks clinfo for device names.
func OpenCL() string {
	if !commandExists("clinfo") {
		return ""
	}
	text := runCombined("clinfo")
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(line)
		name, ok := strings.CutPrefix(trim, "Device Name")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		if name != "" {
			lines = append(lines, name)
		}
	}
	return joinLines(lines)
}

// Codec summarises VA-API from vainfo.
func Codec() string {
	if !commandExists("vainfo") {
		return ""
	}
	text := runCombined("vainfo")
	driver, api := "", ""
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(line)
		if v, ok := strings.CutPrefix(trim, "vainfo: Driver version:"); ok {
			driver = strings.TrimSpace(v)
		}
		if v, ok := strings.CutPrefix(trim, "vainfo: VA-API version:"); ok {
			api = strings.TrimSpace(v)
		}
	}
	switch {
	case driver != "" && api != "":
		return driver + " (VA-API " + api + ")"
	case driver != "":
		return driver
	default:
		return ""
	}
}

// Zpool lists pool name, health, and size from the zpool command.
func Zpool() string {
	if !commandExists("zpool") {
		return ""
	}
	text := runOut("zpool", "list", "-H", "-o", "name,health,size,alloc")
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		alloc := ""
		if len(fields) >= 4 {
			alloc = ", " + fields[3] + " used"
		}
		lines = append(lines, fmt.Sprintf("%s %s (%s%s)", fields[0], fields[1], fields[2], alloc))
	}
	return joinLines(lines)
}

// Weather asks wttr.in, which is what fastfetch uses. location and format
// follow that service: spaces in the format are sent as '+'.
func Weather(location, format string) string {
	format = strings.TrimSpace(format)
	if format == "" {
		format = "%t - %C (%l)"
	}
	format = strings.ReplaceAll(format, "+", " ")
	path := "/"
	if loc := strings.TrimSpace(location); loc != "" {
		path += url.PathEscape(loc)
	}
	query := url.Values{}
	query.Set("format", format)
	endpoint := "https://wttr.in" + path + "?" + query.Encode()
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "curl/0.0.0")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil || resp.StatusCode >= 400 {
		return ""
	}
	text := strings.TrimSpace(string(body))
	if strings.Contains(strings.ToLower(text), "<html") {
		return ""
	}
	return text
}
