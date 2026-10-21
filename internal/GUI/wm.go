package gui

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	cli "github.com/ffgan/gf/internal/CLI"
	"github.com/ffgan/gf/internal/utils"
)

// x11WMs are process names fastfetch/neofetch recognise. Matching comm
// from /proc skips the xprop round trip on the common desktops.
var x11WMs = []string{
	"sowm", "catwm", "fvwm", "dwm", "2bwm", "monsterwm", "tinywm", "x11fs", "xmonad",
	"awesome", "bspwm", "compiz", "cwm", "enlightenment", "fluxbox", "herbstluftwm",
	"i3", "icewm", "jwm", "kwin", "kwin_x11", "leftwm", "marco", "metacity", "muffin",
	"mutter", "openbox", "pekwm", "qtile", "spectrwm", "stumpwm", "xfwm4",
}

func GetWM(osName, kernelName string) string {
	wmOnce.Do(func() {
		wm = detectWM(osName, kernelName)
	})
	return wm
}

func detectWM(osName, kernelName string) string {
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
			tmpPid := ""
			if peer := waylandPeerPID(socket); peer > 0 {
				tmpPid = strconv.Itoa(peer)
			}
			if tmpPid == "" {
				tmpPid = tryLsofOrFuser(socket)
			}
			if tmpPid != "" {
				wm = utils.GetProcessName(tmpPid)
			} else {
				// fallback: scan known wayland wms
				wm = utils.ScanProcesses(psFlags, []string{
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
		wm = utils.ScanProcesses(psFlags, x11WMs)
		if wm == "" && utils.CommandExists("xprop") {
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
			psLine := utils.RunCommand("ps", "-e")
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
			tasklist := utils.RunCommand("tasklist")
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
	if (strings.Contains(wm, "GNOME") && strings.Contains(wm, "Shell")) || wm == "gnome-shell" {
		wm = "Mutter"
	}
	switch wm {
	case "xfwm4":
		wm = "Xfwm4"
	case "kwin", "kwin_x11", "kwin_wayland":
		wm = "KWin"
	}

	return wm
}

func tryLsofOrFuser(path string) string {
	if utils.CommandExists("lsof") {
		out := utils.RunCommand("lsof", "-t", path)
		return strings.TrimSpace(out)
	}
	if utils.CommandExists("fuser") {
		out := utils.RunCommand("fuser", path)
		out = strings.TrimSpace(out)
		parts := strings.Fields(out)
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	return ""
}

func getRootWindowID() string {
	out := utils.RunCommand("xprop", "-root", "-notype", "_NET_SUPPORTING_WM_CHECK")
	fields := strings.Fields(out)
	if len(fields) > 0 {
		return fields[len(fields)-1]
	}
	return ""
}

func getWMNameFromXProp(id string) string {
	// Ask for the one property. Dumping every property and searching for
	// "WM_NAME" also matches the raw WM_NAME value and the _NET_ prefix.
	out := utils.RunCommand("xprop", "-id", id, "-notype", "_NET_WM_NAME")
	i := strings.Index(out, "\"")
	if i < 0 {
		return ""
	}
	rest := out[i+1:]
	j := strings.Index(rest, "\"")
	if j < 0 {
		return strings.TrimSpace(rest)
	}
	return rest[:j]
}
