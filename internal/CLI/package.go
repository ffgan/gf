package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

type PkgInfo struct {
	Manager string
	Count   int
}

// PackageConfig selects how package counts are gathered and printed.
type PackageConfig struct {
	Managers string
	Separate bool
	Minimal  bool
}

func DetectPackages(osname string, cfg PackageConfig) string {
	// on:   '227 (pip), 2699 (rpm), 19 (flatpak-system), 6 (flatpak-user)'
	// tiny: '2945 (pip, rpm, flatpak)'
	// off:  '2945'
	if cfg.Managers == "" {
		cfg.Managers = utils.ON
	}
	managers := []PkgInfo{}
	add := func(manager string, count int) {
		if count > 0 {
			managers = append(managers, PkgInfo{manager, count})
		}
	}

	// --- Language package managers ---
	if !cfg.Minimal && utils.CommandExists("pipx") {
		add("pipx", utils.GetPkgCount("pipx", "list", "--short"))
	}
	if !cfg.Minimal && utils.CommandExists("pip") {
		add("pip", utils.GetPkgCount("pip", "freeze"))
	}
	if !cfg.Minimal {
		if n, ok := countCargoInstalls(); ok {
			add("cargo", n)
		} else if utils.CommandExists("cargo") {
			out, _ := utils.CommandOutput("cargo", "install", "--list")
			lines := 0
			sc := bufio.NewScanner(bytes.NewReader(out))
			for sc.Scan() {
				if !strings.HasPrefix(sc.Text(), " ") {
					lines++
				}
			}
			add("cargo", lines)
		}
	}
	if !cfg.Minimal && utils.CommandExists("npm") {
		if st, err := os.Stat("/usr/lib/node_modules"); err == nil && st.IsDir() {
			add("npm", utils.FilePathGlobLens("/usr/lib/node_modules/*/"))
		} else if st, err := os.Stat("/usr/local/lib/node_modules"); err == nil && st.IsDir() {
			add("npm", utils.FilePathGlobLens("/usr/local/lib/node_modules/*/"))
		} else {
			out, err := utils.CommandOutput("npm", "root", "-g")
			if err == nil {
				root := strings.TrimSpace(string(out))
				add("npm", utils.FilePathGlobLens(filepath.Join(root, "*/")))
			}
		}
	}
	if !cfg.Minimal && utils.CommandExists("pnpm") {
		p := filepath.Join(os.Getenv("HOME"), ".local/share/pnpm/global/5/node_modules/*/")
		add("pnpm", utils.FilePathGlobLens(p))
	}

	// --- System package managers by OS ---
	switch osname {
	case Linux:
		if utils.CommandExists("dpkg") {
			add("dpkg", utils.GetPkgCount("dpkg-query", "-f", ".\\n", "-W"))
		}
		if utils.CommandExists("pacman") {
			add("pacman", utils.GetPkgCount("pacman", "-Qq", "--color", "never"))
		}
		if utils.CommandExists("apk") {
			add("apk", utils.GetPkgCount("apk", "info"))
		} else {
			add("apk", countApkDB())
		}
		if utils.CommandExists("rpm") {
			add("rpm", utils.GetPkgCount("rpm", "-qa", "--nodigest", "--nosignature"))
		}
		if utils.CommandExists("flatpak") {
			if cfg.Separate {
				add("flatpak-system", utils.GetPkgCount("flatpak", "list", "--system"))
				add("flatpak-user", utils.GetPkgCount("flatpak", "list", "--user"))
			} else {
				add("flatpak", utils.GetPkgCount("flatpak", "list"))
			}
		}
		if utils.CommandExists("snap") {
			add("snap", utils.GetPkgCount("snap", "list"))
		}
		if utils.CommandExists("brew") {
			root, err := utils.CommandOutput("brew", "--cellar")
			if err == nil {
				p1 := strings.TrimSpace(string(root)) + "/*"
				add("brew", utils.FilePathGlobLens(p1))
			}
		}
		add("opkg", countOpkg())
		add("emerge", countEmerge())
		add("xbps", countXbps())
		add("slackpkg", countSlack())
		add("kiss", countKiss())
		add("eopkg", countEopkg())
		add("linglong", countLinglong())
		add("guix", countGuix())
		add("nix", countNix())
		add("conda", countConda())
		add("pkgsrc", countDirs("/usr/pkg/pkgdb"))
		add("pacstall", countDirs("/var/lib/pacstall/metadata"))
		add("paludis", countNestedDirs("/var/db/paludis/repositories"))
		add("tazpkg", countDirs("/var/lib/tazpkg/installed"))
		add("sorcery", countFiles("/var/state/sorcery/packages"))
		add("cards", countFiles("/var/lib/pkg/db"))
		add("scratch", countDirs("/var/lib/scratchpkg/db"))
		add("ipkg", countPrefixedLines([]string{"/usr/lib/ipkg/status", "/var/lib/ipkg/status"}, "Package:"))
		add("moss", countDirs("/var/db/moss/pkgs"))
		if home := os.Getenv("HOME"); home != "" {
			add("am", countDirs(filepath.Join(home, ".local/share/am")))
		}
	case FreeBSD, DragonFly, BSD, "MidnightBSD", "HardenedBSD", "GhostBSD", "NomadBSD":
		add("pkg", countFreeBSDPkg())
	case OpenBSD:
		add("pkg_add", countDirs("/var/db/pkg"))
	case NetBSD:
		add("pkgin", countDirs("/usr/pkg/pkgdb"))
	case Haiku:
		add("pkgman", countSuffix("/boot/system/packages", ".hpkg"))
	case Solaris, illumos:
		add("pkg", countDirs("/var/sadm/pkg"))
	case MINIX:
		add("pkgin", countDirs("/usr/pkg/pkgdb"))
	default:
		if isBSDFamily(osname) {
			add("pkg", countFreeBSDPkg())
		}
	case Darwin:
		if utils.CommandExists("brew") {
			add("brew", utils.FilePathGlobLens("/usr/local/Cellar/*/"))
		}
		if utils.CommandExists("port") {
			add("macports", utils.GetPkgCount("port", "installed"))
		}
	case Windows:
		if utils.CommandExists("choco") {
			pdata := os.Getenv("ProgramData")
			if pdata == "" {
				pdata = `C:\ProgramData`
			}
			add("choco", utils.FilePathGlobLens(filepath.Join(pdata, "chocolatey/lib/*")))
		}
		// winget list talks to the network and often takes tens of seconds.
		// The uninstall registry is the local program list.
		add("programs", countInstalledPrograms())
	}

	// --- Total and formatting ---
	total := 0
	for _, m := range managers {
		total += m.Count
	}

	if len(managers) == 0 || total == 0 {
		return ""
	}

	switch cfg.Managers {
	case utils.OFF:
		return strconv.Itoa(total)
	case "tiny":
		names := make([]string, 0, len(managers))
		for _, m := range managers {
			names = append(names, m.Manager)
		}
		return fmt.Sprintf("%d (%s)", total, strings.Join(names, ", "))
	default:
		res := make([]string, 0, len(managers))
		for _, m := range managers {
			res = append(res, fmt.Sprintf("%d (%s)", m.Count, m.Manager))
		}
		return strings.Join(res, ", ")
	}
}
