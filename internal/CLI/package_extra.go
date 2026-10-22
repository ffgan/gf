package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func countPrefixedLines(paths []string, prefix string) int {
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		n := 0
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if strings.HasPrefix(sc.Text(), prefix) {
				n++
			}
		}
		f.Close()
		return n
	}
	return 0
}

func countApkDB() int {
	return countPrefixedLines([]string{
		"/lib/apk/db/installed",
		"/var/lib/apk/db/installed",
		"/usr/lib/apk/db/installed",
	}, "P:")
}

// countPkgDB counts FreeBSD/DragonFly packages. pkgng stores them in
// local.sqlite; older installs are one directory per package.
func countPkgDB(root string) int {
	db := filepath.Join(root, "local.sqlite")
	if st, err := os.Stat(db); err == nil && !st.IsDir() && utils.CommandExists("sqlite3") {
		out, err := utils.CommandOutput("sqlite3", db, "SELECT COUNT(*) FROM packages;")
		if err == nil {
			if n, conv := strconv.Atoi(strings.TrimSpace(string(out))); conv == nil {
				return n
			}
		}
	}
	return countDirs(root)
}

func countFreeBSDPkg() int {
	return countPkgDB("/var/db/pkg")
}

func countOpkg() int {
	return countPrefixedLines([]string{"/var/lib/opkg/status", "/usr/lib/opkg/status"}, "Package:")
}

func countEmerge() int {
	return countNestedDirs("/var/db/pkg")
}

func countNestedDirs(root string) int {
	cats, err := os.ReadDir(root)
	if err != nil {
		return 0
	}
	n := 0
	for _, c := range cats {
		if !c.IsDir() || strings.HasPrefix(c.Name(), ".") {
			continue
		}
		n += countDirs(filepath.Join(root, c.Name()))
	}
	return n
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

func countXbps() int {
	return countSuffix("/var/db/xbps", ".files")
}

func countSlack() int {
	return countFiles("/var/log/packages")
}

func countKiss() int {
	return countDirs("/var/db/kiss/installed")
}

func countEopkg() int {
	for _, dir := range []string{"/var/lib/eopkg/package", "/var/lib/pisi/package"} {
		if n := countFiles(dir); n > 0 {
			return n
		}
	}
	return 0
}

func countFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() {
			n++
		}
	}
	return n
}

func countLinglong() int {
	return countDirs("/var/lib/linglong/layers")
}

func countGuix() int {
	home := os.Getenv("HOME")
	user := os.Getenv("USER")
	paths := []string{
		filepath.Join(home, ".guix-profile", "manifest"),
		filepath.Join("/var/guix/profiles/per-user", user, "guix-profile", "manifest"),
	}
	for _, path := range paths {
		if n := countManifestNames(path); n > 0 {
			return n
		}
	}
	return 0
}

func countNix() int {
	home := os.Getenv("HOME")
	user := os.Getenv("USER")
	paths := []string{
		filepath.Join(home, ".nix-profile", "manifest.nix"),
		filepath.Join("/nix/var/nix/profiles/per-user", user, "profile", "manifest.nix"),
	}
	for _, path := range paths {
		if n := countManifestNames(path); n > 0 {
			return n
		}
	}
	return 0
}

func countManifestNames(path string) int {
	info, err := os.Stat(path)
	if err != nil || info.Size() > 2<<20 || info.IsDir() {
		return 0
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	return strings.Count(string(b), "(name ")
}

func countConda() int {
	home := os.Getenv("HOME")
	dirs := []string{os.Getenv("CONDA_PREFIX")}
	for _, name := range []string{"miniconda3", "anaconda3", "miniforge3", "mambaforge"} {
		dirs = append(dirs, filepath.Join(home, name))
	}
	dirs = append(dirs, "/opt/conda")
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		if n := countSuffix(filepath.Join(dir, "conda-meta"), ".json"); n > 0 {
			return n
		}
	}
	return 0
}

// countCargoInstalls reads Cargo's install index. Spawning
// `cargo install --list` starts the compiler toolchain.
func countCargoInstalls() (int, bool) {
	home := os.Getenv("CARGO_HOME")
	if home == "" {
		if h := os.Getenv("HOME"); h != "" {
			home = filepath.Join(h, ".cargo")
		} else if h := os.Getenv("USERPROFILE"); h != "" {
			home = filepath.Join(h, ".cargo")
		}
	}
	if home == "" {
		return -1, false
	}
	if data, err := os.ReadFile(filepath.Join(home, ".crates2.json")); err == nil {
		return bytes.Count(data, []byte(`"bins"`)), true
	}
	data, err := os.ReadFile(filepath.Join(home, ".crates.toml"))
	if err != nil {
		return -1, false
	}
	n := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), `"`) {
			n++
		}
	}
	return n, true
}

// countInstalledPrograms counts uninstall-registry keys. `winget list`
// contacts the network and is the multi-second stall on Windows.
func countInstalledPrograms() int {
	keys := []string{
		`HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`HKLM\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
		`HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
	}
	n := 0
	for _, key := range keys {
		n += countRegSubkeys(utils.RunCommand("reg.exe", "query", key))
	}
	return n
}

func addWindowsPackages(add func(string, int)) {
	home := os.Getenv("USERPROFILE")
	if home == "" {
		home = os.Getenv("HOME")
	}
	pdata := os.Getenv("ProgramData")
	cfg := ""
	if home != "" {
		cfg = filepath.Join(home, ".config", "scoop", "config.json")
	}
	user, global := scoopAppCounts(cfg, home, pdata)
	if global > 0 {
		add("scoop-user", user)
		add("scoop-global", global)
	} else {
		add("scoop", user)
	}

	chocoRoot := os.Getenv("ChocolateyInstall")
	if chocoRoot == "" && utils.CommandExists("choco") {
		if pdata == "" {
			pdata = `C:\ProgramData`
		}
		chocoRoot = filepath.Join(pdata, "chocolatey")
	}
	if chocoRoot != "" {
		add("choco", countDirsExcept(filepath.Join(chocoRoot, "lib"), "choco"))
	}
	add("pacman", countMsysPacman())
	// winget list talks to the network and often takes tens of seconds.
	// The uninstall registry is the local program list.
	add("programs", countInstalledPrograms())
}

func scoopAppCounts(configPath, home, programData string) (user, global int) {
	root, globalRoot := "", ""
	if configPath != "" {
		if data, err := os.ReadFile(configPath); err == nil {
			var cfg struct {
				Root   string `json:"root_path"`
				Global string `json:"global_path"`
			}
			if json.Unmarshal(data, &cfg) == nil {
				root = cfg.Root
				globalRoot = cfg.Global
			}
		}
	}
	if root == "" && home != "" {
		root = filepath.Join(home, "scoop")
	}
	if root != "" {
		user = countDirsExcept(filepath.Join(root, "apps"), "scoop")
	}
	if globalRoot == "" && programData != "" {
		globalRoot = filepath.Join(programData, "scoop")
	}
	if globalRoot != "" {
		global = countDirsExcept(filepath.Join(globalRoot, "apps"), "scoop")
	}
	return user, global
}

func countMsysPacman() int {
	prefix := os.Getenv("MSYSTEM_PREFIX")
	if prefix == "" {
		return 0
	}
	return countDirs(filepath.Join(prefix, "..", "var", "lib", "pacman", "local"))
}

func countDirsExcept(path, skip string) int {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if skip != "" && strings.EqualFold(e.Name(), skip) {
			continue
		}
		n++
	}
	return n
}

func countSuffix(dir, suffix string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), suffix) {
			n++
		}
	}
	return n
}
