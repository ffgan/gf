package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

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

// DetectPackages counts installed packages from each manager's files.
// Spawning pip, rpm, npm, or cargo is the slow path fastfetch avoids.
func DetectPackages(osname string, cfg PackageConfig) string {
	if cfg.Managers == "" {
		cfg.Managers = utils.ON
	}
	type probe struct {
		name string
		fn   func() int
	}
	var probes []probe
	add := func(name string, fn func() int) {
		probes = append(probes, probe{name, fn})
	}

	if !cfg.Minimal {
		add("pipx", countPipx)
		add("pip", countPip)
		add("cargo", func() int {
			n, ok := countCargoInstalls()
			if !ok {
				return 0
			}
			return n
		})
		add("npm", countNpm)
		add("pnpm", countPnpm)
	}

	switch osname {
	case Linux:
		add("dpkg", countDpkg)
		add("pacman", countPacman)
		add("apk", countApkOrDB)
		add("rpm", countRpm)
		if cfg.Separate {
			add("flatpak-system", func() int { return countDirs("/var/lib/flatpak/app") })
			add("flatpak-user", countFlatpakUser)
		} else {
			add("flatpak", countFlatpak)
		}
		add("snap", countSnap)
		add("brew", countBrew)
		add("opkg", countOpkg)
		add("emerge", countEmerge)
		add("xbps", countXbps)
		add("slackpkg", countSlack)
		add("kiss", countKiss)
		add("eopkg", countEopkg)
		add("linglong", countLinglong)
		add("guix", countGuix)
		add("nix", countNix)
		add("conda", countConda)
		add("pkgsrc", func() int { return countDirs("/usr/pkg/pkgdb") })
		add("pacstall", func() int { return countDirs("/var/lib/pacstall/metadata") })
		add("paludis", func() int { return countNestedDirs("/var/db/paludis/repositories") })
		add("tazpkg", func() int { return countDirs("/var/lib/tazpkg/installed") })
		add("sorcery", func() int { return countFiles("/var/state/sorcery/packages") })
		add("cards", func() int { return countFiles("/var/lib/pkg/db") })
		add("scratch", func() int { return countDirs("/var/lib/scratchpkg/db") })
		add("ipkg", func() int {
			return countPrefixedLines([]string{"/usr/lib/ipkg/status", "/var/lib/ipkg/status"}, "Package:")
		})
		add("moss", func() int { return countDirs("/var/db/moss/pkgs") })
		if home := os.Getenv("HOME"); home != "" {
			amHome := home
			add("am", func() int { return countDirs(filepath.Join(amHome, ".local/share/am")) })
		}
	case FreeBSD, DragonFly, BSD, "MidnightBSD", "HardenedBSD", "GhostBSD", "NomadBSD":
		add("pkg", countFreeBSDPkg)
	case OpenBSD:
		add("pkg_add", func() int { return countDirs("/var/db/pkg") })
	case NetBSD:
		add("pkgin", func() int { return countDirs("/usr/pkg/pkgdb") })
	case Haiku:
		add("pkgman", func() int { return countSuffix("/boot/system/packages", ".hpkg") })
	case Solaris, illumos:
		add("pkg", func() int { return countDirs("/var/sadm/pkg") })
	case MINIX:
		add("pkgin", func() int { return countDirs("/usr/pkg/pkgdb") })
	case Darwin:
		add("brew", countBrew)
		add("macports", countMacports)
	default:
		if isBSDFamily(osname) {
			add("pkg", countFreeBSDPkg)
		}
	}

	counts := make([]int, len(probes))
	var wg sync.WaitGroup
	for i := range probes {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			counts[i] = probes[i].fn()
		}(i)
	}
	wg.Wait()

	managers := make([]PkgInfo, 0, len(probes))
	for i, p := range probes {
		if counts[i] > 0 {
			managers = append(managers, PkgInfo{p.name, counts[i]})
		}
	}
	if osname == Windows {
		addWindowsPackages(func(name string, n int) {
			if n > 0 {
				managers = append(managers, PkgInfo{name, n})
			}
		})
	}

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

func countApkOrDB() int {
	if n := countApk(); n > 0 {
		return n
	}
	return countApkDB()
}

func countFlatpakUser() int {
	data := os.Getenv("XDG_DATA_HOME")
	if data == "" {
		data = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	return countDirs(filepath.Join(data, "flatpak", "app"))
}

func countMacports() int {
	if !utils.CommandExists("port") {
		return 0
	}
	return utils.GetPkgCount("port", "installed")
}

func countDpkg() int {
	return countDpkgFile("/var/lib/dpkg/status")
}

func countDpkgFile(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()

	n := 0
	sc := bufio.NewScanner(f)
	// The status file is one logical record per package. dpkg-query -W
	// lists packages whose status is exactly "install ok installed".
	for sc.Scan() {
		if sc.Text() == "Status: install ok installed" {
			n++
		}
	}
	return n
}

func countPacman() int {
	return countDirs("/var/lib/pacman/local")
}

func countApk() int {
	for _, path := range []string{"/lib/apk/db/installed", "/var/lib/apk/db/installed"} {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		n := 0
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if strings.HasPrefix(sc.Text(), "P:") {
				n++
			}
		}
		f.Close()
		return n
	}
	return 0
}

func countRpm() int {
	for _, db := range []string{
		"/usr/lib/sysimage/rpm/rpmdb.sqlite",
		"/var/lib/rpm/rpmdb.sqlite",
	} {
		if !utils.FileExists(db) {
			continue
		}
		// A single COUNT is what fastfetch asks sqlite for. `rpm -qa`
		// opens the same database and then formats every package name.
		out := utils.RunCommand("sqlite3", db, "SELECT COUNT(*) FROM Packages;")
		if n, err := strconv.Atoi(strings.TrimSpace(out)); err == nil && n > 0 {
			return n
		}
	}
	if utils.CommandExists("rpm") {
		return utils.GetPkgCount("rpm", "-qa", "--nodigest", "--nosignature")
	}
	return 0
}

func countFlatpak() int {
	dirs := []string{"/var/lib/flatpak/app"}
	data := os.Getenv("XDG_DATA_HOME")
	if data == "" {
		data = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	dirs = append(dirs, filepath.Join(data, "flatpak", "app"))

	n := 0
	seen := map[string]struct{}{}
	for _, dir := range dirs {
		real, err := filepath.EvalSymlinks(dir)
		if err != nil {
			real = dir
		}
		if _, ok := seen[real]; ok {
			continue
		}
		seen[real] = struct{}{}
		n += countDirs(dir)
	}
	return n
}

func countSnap() int {
	entries, err := os.ReadDir("/snap")
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		switch e.Name() {
		case "bin", "README":
			continue
		}
		n++
	}
	return n
}

func countBrew() int {
	dirs := []string{}
	if cellar := os.Getenv("HOMEBREW_CELLAR"); cellar != "" {
		dirs = append(dirs, cellar)
	}
	dirs = append(dirs,
		"/home/linuxbrew/.linuxbrew/Cellar",
		"/opt/homebrew/Cellar",
		"/usr/local/Cellar",
		filepath.Join(os.Getenv("HOME"), ".linuxbrew", "Cellar"),
	)
	for _, dir := range dirs {
		if n := countDirs(dir); n > 0 {
			return n
		}
	}
	if !utils.CommandExists("brew") {
		return 0
	}
	root, err := utils.CommandOutput("brew", "--cellar")
	if err != nil {
		return 0
	}
	return countDirs(strings.TrimSpace(string(root)))
}

func countPipx() int {
	return countDirs(filepath.Join(os.Getenv("HOME"), ".local", "share", "pipx", "venvs"))
}

func countPnpm() int {
	return utils.FilePathGlobLens(filepath.Join(os.Getenv("HOME"), ".local/share/pnpm/global/5/node_modules/*/"))
}

func countNpm() int {
	dirs := []string{
		"/usr/lib/node_modules",
		"/usr/local/lib/node_modules",
		"/usr/lib64/node_modules",
	}
	if prefix := os.Getenv("NPM_CONFIG_PREFIX"); prefix != "" {
		dirs = append(dirs, filepath.Join(prefix, "lib", "node_modules"))
	}
	dirs = append(dirs, filepath.Join(os.Getenv("HOME"), ".npm-global", "lib", "node_modules"))
	for _, dir := range dirs {
		n := countDirs(dir)
		if n == 0 {
			continue
		}
		// `.bin` is a directory of shims, not a package.
		if utils.PathExists(filepath.Join(dir, ".bin")) {
			n--
		}
		if n > 0 {
			return n
		}
	}
	return 0
}

// countPip counts distributions the way `pip freeze` does: one name per
// .dist-info or .egg-info, with pip itself omitted. Reading the
// site-packages directories avoids starting the interpreter.
func countPip() int {
	globs := []string{
		"/usr/lib/python3/dist-packages",
		"/usr/lib/python3.*/dist-packages",
		"/usr/lib/python3.*/site-packages",
		"/usr/local/lib/python3.*/dist-packages",
		"/usr/local/lib/python3.*/site-packages",
		filepath.Join(os.Getenv("HOME"), ".local/lib/python3.*/site-packages"),
	}
	names := map[string]struct{}{}
	seenDir := map[string]struct{}{}
	for _, pattern := range globs {
		matches, _ := filepath.Glob(pattern)
		for _, dir := range matches {
			real, err := filepath.EvalSymlinks(dir)
			if err != nil {
				real = dir
			}
			if _, ok := seenDir[real]; ok {
				continue
			}
			seenDir[real] = struct{}{}
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, e := range entries {
				name := pipMetaName(e.Name())
				if name == "" || name == "pip" {
					continue
				}
				names[name] = struct{}{}
			}
		}
	}
	return len(names)
}

func pipMetaName(base string) string {
	var name string
	switch {
	case strings.HasSuffix(base, ".dist-info"):
		name = strings.TrimSuffix(base, ".dist-info")
	case strings.HasSuffix(base, ".egg-info"):
		name = strings.TrimSuffix(base, ".egg-info")
	default:
		return ""
	}
	if i := strings.LastIndex(name, "-"); i > 0 {
		rest := name[i+1:]
		if rest != "" && rest[0] >= '0' && rest[0] <= '9' {
			name = name[:i]
		}
	}
	name = strings.ToLower(strings.ReplaceAll(name, "_", "-"))
	return name
}
