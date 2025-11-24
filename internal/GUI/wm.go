package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	cli "github.com/ffgan/gf/internal/CLI"
)

func GetWM(osName, kernelName string) string {
	if wmRun {
		return wm
	}

	// kernelName := getKernelName()
	// osName := cli.GetOS()
	xdgRuntime := os.Getenv("XDG_RUNTIME_DIR")
	waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
	if waylandDisplay == "" {
		waylandDisplay = "wayland-0"
	}

	display := os.Getenv("DISPLAY")

	// --- detect process listing flags ---
	psFlags := []string{"-e"}
	if strings.Contains(kernelName, "OpenBSD") {
		psFlags = []string{"x", "-c"}
	}

	// --- Wayland detection ---
	if xdgRuntime != "" {
		socket := filepath.Join(xdgRuntime, waylandDisplay)
		if fi, err := os.Stat(socket); err == nil && fi.Mode()&os.ModeSocket != 0 {
			tmpPid := tryLsofOrFuser(socket)
			if tmpPid != "" {
				wm = getProcessName(tmpPid)
			} else {
				// fallback: scan known wayland wms
				wm = scanProcesses(psFlags, []string{
					"arcan", "asc", "clayland", "dwc", "dwl", "fireplace",
					"gnome-shell", "greenfield", "grefsen", "hikari", "Hyprland",
					"kwin", "lipstick", "maynard", "mazecompositor", "motorcar",
					"orbital", "orbment", "perceptia", "river", "rustland",
					"sway", "ulubis", "velox", "wavy", "way-cooler",
					"wayfire", "wayhouse", "westeros", "westford", "weston",
				})
			}
		}
	}

	// --- X11 detection ---
	if wm == "" && display != "" && osName != cli.MacOSX && osName != cli.MacOS && osName != cli.FreeMiNT {
		wm = scanProcesses(psFlags, []string{
			"sowm", "catwm", "fvwm", "dwm", "2bwm", "monsterwm", "tinywm", "x11fs", "xmonad",
		})
		if wm == "" && commandExists("xprop") {
			id := getRootWindowID()
			if id != "" {
				wm = getWMNameFromXProp(id)
			}
		}
	}

	// --- macOS detection ---
	if wm == "" {
		switch osName {
		case cli.MacOSX, cli.MacOS:
			psLine := runCmd("ps", "-e")
			switch {
			case strings.Contains(psLine, "chunkwm"):
				wm = "chunkwm"
			case strings.Contains(psLine, "kwm"):
				wm = "Kwm"
			case strings.Contains(psLine, "yabai"):
				wm = "yabai"
			case strings.Contains(psLine, "Amethyst"):
				wm = "Amethyst"
			case strings.Contains(psLine, "Spectacle"):
				wm = "Spectacle"
			case strings.Contains(psLine, "Rectangle"):
				wm = "Rectangle"
			default:
				wm = "Quartz Compositor"
			}

		case cli.Windows:
			tasklist := runCmd("tasklist")
			for _, w := range []string{"bugn", "Windawesome", "blackbox", "emerge", "litestep"} {
				if strings.Contains(tasklist, w) {
					wm = w
					break
				}
			}
			if wm == "blackbox" {
				wm = "bbLean (Blackbox)"
			}
			if wm != "" {
				wm += ", DWM.exe"
			}

		case cli.FreeMiNT:
			procs, _ := filepath.Glob("/proc/*")
			for _, p := range procs {
				switch {
				case strings.Contains(p, "xaaes") || strings.Contains(p, "xaloader"):
					wm = "XaAES"
				case strings.Contains(p, "myaes"):
					wm = "MyAES"
				case strings.Contains(p, "naes"):
					wm = "N.AES"
				case strings.Contains(p, "geneva"):
					wm = "Geneva"
				default:
					wm = "Atari AES"
				}
			}
		}
	}

	// --- normalize names ---
	if strings.Contains(wm, "WINDOWMAKER") {
		wm = "wmaker"
	}
	if strings.Contains(wm, "GNOME") && strings.Contains(wm, "Shell") {
		wm = "Mutter"
	}

	wmRun = true
	return wm
}

func getKernelName() string {
	out := runCmd("uname", "-s")
	return strings.TrimSpace(out)
}

func tryLsofOrFuser(path string) string {
	if commandExists("lsof") {
		out := runCmd("lsof", "-t", path)
		return strings.TrimSpace(out)
	}
	if commandExists("fuser") {
		out := runCmd("fuser", path)
		out = strings.TrimSpace(out)
		parts := strings.Fields(out)
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	return ""
}

func getProcessName(pid string) string {
	out := runCmd("ps", "-p", pid, "-ho", "comm=")
	return strings.TrimSpace(out)
}

func scanProcesses(psFlags []string, names []string) string {
	args := append([]string{}, psFlags...)
	out := runCmd("ps", args...)
	for _, n := range names {
		re := regexp.MustCompile(fmt.Sprintf(`(?m)^%s$`, regexp.QuoteMeta(n)))
		if re.FindString(out) != "" {
			return n
		}
	}
	return ""
}

func getRootWindowID() string {
	out := runCmd("xprop", "-root", "-notype", "_NET_SUPPORTING_WM_CHECK")
	fields := strings.Fields(out)
	if len(fields) > 0 {
		return fields[len(fields)-1]
	}
	return ""
}

func getWMNameFromXProp(id string) string {
	out := runCmd("xprop", "-id", id, "-notype", "-len", "100", "-f", "_NET_WM_NAME", "8t")
	if idx := strings.Index(out, "WM_NAME = "); idx != -1 {
		val := out[idx+len("WM_NAME = "):]
		val = strings.Trim(val, "\" \n")
		return val
	}
	return ""
}
