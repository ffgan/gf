package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	cli "github.com/ffgan/gf/internal/CLI"
	"github.com/ffgan/gf/internal/utils"
)

var (
	deOnce    sync.Once
	de        string
	deVersion string
	wmOnce    sync.Once
	wm        string

	desktopOnce   sync.Once
	desktopCached string
)

func detectDesktopEnvironment() string {
	desktopOnce.Do(func() {
		if current := os.Getenv("XDG_CURRENT_DESKTOP"); current != "" {
			desktopCached = current
			return
		}
		// One /proc walk instead of `ps -e`. Theme, icons, font and cursor
		// all ask for this, and each used to list every process.
		var blob string
		if names, ok := utils.ProcessNames(); ok {
			blob = strings.Join(names, "\n")
		} else {
			blob = utils.RunCommand("ps", "-e")
		}
		switch {
		case strings.Contains(blob, "plasmashell"):
			desktopCached = "KDE"
		case strings.Contains(blob, "gnome-shell"):
			desktopCached = "Gnome"
		case strings.Contains(blob, "cinnamon"):
			desktopCached = "Cinnamon"
		case strings.Contains(blob, "xfce"):
			desktopCached = "Xfce"
		case strings.Contains(blob, "lxqt"):
			desktopCached = "LXQt"
		default:
			desktopCached = "Unknown"
		}
	})
	return desktopCached
}

func GetDE(osName, distro, kernel_name string) string {
	deOnce.Do(func() {
		de = detectDE(osName, distro, kernel_name)
	})
	return de
}

func detectDE(osName, distro, kernel_name string) string {
	switch osName {
	case cli.MacOSX, cli.MacOS:
		de = "Aqua"

	case cli.RavynOS:
		de = "Magma"

	case cli.Windows:
		switch {
		case strings.Contains(distro, "Windows 10") || strings.Contains(distro, "Windows 11"):
			de = "Fluent"
		case strings.Contains(distro, "Windows 8"):
			de = "Metro"
		default:
			de = "Aero"
		}

	case cli.FreeMiNT:
		matches, _ := filepath.Glob("/proc/*")
		for _, m := range matches {
			switch {
			case strings.Contains(m, "thing"):
				de = "Thing"
			case strings.Contains(m, "jinnee"):
				de = "Jinnee"
			case strings.Contains(m, "tera"):
				de = "Teradesk"
			case strings.Contains(m, "neod"):
				de = "NeoDesk"
			case strings.Contains(m, "zdesk"):
				de = "zDesk"
			case strings.Contains(m, "mdesk"):
				de = "mDesk"
			}
		}

	default:
		desktopSession := os.Getenv("DESKTOP_SESSION")
		xdgCurrent := os.Getenv("XDG_CURRENT_DESKTOP")

		switch {
		case strings.Contains(desktopSession, "regolith"):
			de = "Regolith"

		case desktopSession == "fly":
			de = "Fly"

		case xdgCurrent != "":
			de = strings.ReplaceAll(xdgCurrent, "X-", "")
			de = strings.ReplaceAll(de, "Budgie:GNOME", "Budgie")
			de = strings.ReplaceAll(de, ":Unity7:ubuntu", "")

		case desktopSession != "":
			de = filepath.Base(desktopSession)
			if strings.EqualFold(de, "trinity") {
				de = "Trinity"
			}

		case os.Getenv("GNOME_DESKTOP_SESSION_ID") != "":
			de = "GNOME"

		case os.Getenv("MATE_DESKTOP_SESSION_ID") != "":
			de = "MATE"

		case os.Getenv("TDE_FULL_SESSION") != "":
			de = "Trinity"
		}

		// Remove DE if it's the same as WM
		wmName := GetWM(osName, kernel_name)
		if de == wmName {
			de = ""
			return de
		}
	}

	// Try xprop as fallback
	if de == "" && os.Getenv("DISPLAY") != "" && utils.CommandExists("xprop") {
		out := utils.RunCommand("xprop", "-root")
		if strings.Contains(out, "KDE_SESSION_VERSION") {
			de = "KDE"
		} else if strings.Contains(out, "xfce4") {
			de = "Xfce4"
		} else if strings.Contains(out, "xfce5") {
			de = "Xfce5"
		}
	}

	// Normalize strings
	switch {
	case strings.Contains(de, "KDE_SESSION_VERSION"):
		de = "KDE" + strings.Split(de, "=")[1]
	case strings.Contains(de, "xfce4"):
		de = "Xfce4"
	case strings.Contains(de, "xfce5"):
		de = "Xfce5"
	case strings.Contains(strings.ToLower(de), "mate"):
		de = "MATE"
	case strings.Contains(de, "GNOME"):
		de = "GNOME"
	case strings.Contains(de, "MUFFIN"):
		de = "Cinnamon"
	}

	// KDE/Plasma version handling
	kdeVer := os.Getenv("KDE_SESSION_VERSION")
	if kdeVer >= "4" {
		de = strings.Replace(de, "KDE", "Plasma", 1)
	}
	if kdeVer >= "6" {
		de = strings.Replace(de, "Plasma", "Plasma6", 1)
	}

	// Add version info if available
	if deVersion == utils.ON && de != "" {
		deVer, kfVer, qtVer := getDEVersion(de)
		if deVer != "" {
			if strings.HasPrefix(de, "Plasma") {
				de = fmt.Sprintf("%s %s [KF %s] [Qt %s]", de, deVer, kfVer, qtVer)
			} else {
				de = fmt.Sprintf("%s %s", de, deVer)
			}
		}
	}

	// Add session type (X11 / Wayland)
	if st := os.Getenv("XDG_SESSION_TYPE"); st != "" {
		de += fmt.Sprintf(" (%s)", st)
	}

	return de
}

func getDEVersion(name string) (string, string, string) {
	var deVer, kfVer, qtVer string
	switch {
	case strings.HasPrefix(name, "Plasma6"):
		deVer = utils.RunCommand("plasmashell", "--version")
		kinfo := utils.RunCommand("kinfo")
		kfVer, qtVer = parseKFQt(kinfo)

	case strings.HasPrefix(name, "Plasma"):
		deVer = utils.RunCommand("plasmashell", "--version")
		kinfo := utils.RunCommand("kf5-config", "--version")
		kfVer, qtVer = parseKFQt(kinfo)

	case strings.HasPrefix(name, "MATE"):
		deVer = utils.RunCommand("mate-session", "--version")
	case strings.HasPrefix(name, "Xfce"):
		deVer = utils.RunCommand("xfce4-session", "--version")
	case strings.HasPrefix(name, "GNOME"):
		deVer = utils.RunCommand("gnome-shell", "--version")
	case strings.HasPrefix(name, "Cinnamon"):
		deVer = utils.RunCommand("cinnamon", "--version")
	case strings.HasPrefix(name, "Budgie"):
		deVer = utils.RunCommand("budgie-desktop", "--version")
	case strings.HasPrefix(name, "LXQt"):
		deVer = utils.RunCommand("lxqt-session", "--version")
	case strings.HasPrefix(name, "Trinity"):
		deVer = utils.RunCommand("tde-config", "--version")
	case strings.HasPrefix(name, "Unity"):
		deVer = utils.RunCommand("unity", "--version")
	}
	return utils.Trim(deVer), utils.Trim(kfVer), utils.Trim(qtVer)
}

func parseKFQt(info string) (kfVer, qtVer string) {
	for _, line := range strings.Split(info, "\n") {
		if strings.Contains(line, "KDE Frameworks") {
			kfVer = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Qt") {
			qtVer = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return
}
