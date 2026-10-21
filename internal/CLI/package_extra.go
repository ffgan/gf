package cli

import (
	"bufio"
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
