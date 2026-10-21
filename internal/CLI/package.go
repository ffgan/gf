package cli

import (
	"bufio"
	"encoding/json"
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

type pkgProbe struct {
	name string
	fn   func() int
}

// DetectPackages counts installed packages by reading each manager's
// database directly. fastfetch does the same: a directory listing or a
// status file is milliseconds, while `pip freeze`, `rpm -qa`, `npm root`
// and `cargo install --list` each start a whole runtime.
func DetectPackages(osname string) string {
	probes := []pkgProbe{
		{"pipx", countPipx},
		{"pip", countPip},
		{"cargo", countCargo},
		{"npm", countNpm},
		{"pnpm", countPnpm},
	}

	switch osname {
	case Linux:
		probes = append(probes,
			pkgProbe{"dpkg", countDpkg},
			pkgProbe{"pacman", countPacman},
			pkgProbe{"apk", countApk},
			pkgProbe{"rpm", countRpm},
			pkgProbe{"flatpak", countFlatpak},
			pkgProbe{"snap", countSnap},
			pkgProbe{"brew", countBrew},
		)
	case Darwin:
		probes = append(probes,
			pkgProbe{"brew", func() int { return countDirs("/usr/local/Cellar") }},
			pkgProbe{"macports", func() int {
				if !utils.CommandExists("port") {
					return 0
				}
				return utils.GetPkgCount("port", "installed")
			}},
		)
	case Windows:
		probes = append(probes,
			pkgProbe{"choco", countChoco},
			pkgProbe{"winget", func() int {
				if !utils.CommandExists("winget") {
					return 0
				}
				return utils.GetPkgCount("winget", "list", "--accept-source-agreements")
			}},
		)
	}

	counts := make([]int, len(probes))
	var wg sync.WaitGroup
	for i, p := range probes {
		wg.Add(1)
		go func(i int, fn func() int) {
			defer wg.Done()
			counts[i] = fn()
		}(i, p.fn)
	}
	wg.Wait()

	managers := make([]PkgInfo, 0, len(probes))
	for i, p := range probes {
		if counts[i] > 0 {
			managers = append(managers, PkgInfo{p.name, counts[i]})
		}
	}
	if len(managers) == 0 {
		return "No packages found"
	}
	res := make([]string, 0, len(managers))
	for _, m := range managers {
		res = append(res, fmt.Sprintf("%d (%s)", m.Count, m.Manager))
	}
	return strings.Join(res, ", ")
}

func countDirs(path string) int {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if e.IsDir() {
			n++
		}
	}
	return n
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

func countChoco() int {
	pdata := os.Getenv("ProgramData")
	if pdata == "" {
		pdata = `C:\ProgramData`
	}
	return countDirs(filepath.Join(pdata, "chocolatey", "lib"))
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

func countCargo() int {
	home := os.Getenv("HOME")
	if b, err := os.ReadFile(filepath.Join(home, ".cargo", ".crates2.json")); err == nil {
		var doc struct {
			Installs map[string]json.RawMessage `json:"installs"`
		}
		if json.Unmarshal(b, &doc) == nil {
			return len(doc.Installs)
		}
	}
	b, err := os.ReadFile(filepath.Join(home, ".cargo", ".crates.toml"))
	if err != nil {
		return 0
	}
	n := 0
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\"") && strings.Contains(line, "=") {
			n++
		}
	}
	return n
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
