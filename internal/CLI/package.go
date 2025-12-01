package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

type PkgInfo struct {
	Manager string
	Count   int
}

func DetectPackages(osname string) string {
	// TODO: 修复与hyfetch不一致的地方
	// Packages: 227 (pip), 2699 (rpm), 25 (flatpak)
	// Packages: 227 (pip), 2699 (rpm), 19 (flatpak-system), 6 (flatpak-user)
	managers := []PkgInfo{}
	add := func(manager string, count int) {
		if count > 0 {
			managers = append(managers, PkgInfo{manager, count})
		}
	}

	// --- Language package managers ---
	if utils.CommandExists("pipx") {
		add("pipx", utils.GetPkgCount("pipx", "list", "--short"))
	}
	if utils.CommandExists("pip") {
		add("pip", utils.GetPkgCount("pip", "freeze"))
	}
	if utils.CommandExists("cargo") {
		cmd := exec.Command("cargo", "install", "--list")
		out, _ := cmd.Output()
		lines := 0
		sc := bufio.NewScanner(bytes.NewReader(out))
		for sc.Scan() {
			if !strings.HasPrefix(sc.Text(), " ") {
				lines++
			}
		}
		add("cargo", lines)
	}
	if utils.CommandExists("npm") {
		if st, err := os.Stat("/usr/lib/node_modules"); err == nil && st.IsDir() {
			add("npm", utils.FilePathGlobLens("/usr/lib/node_modules/*/"))
		} else if st, err := os.Stat("/usr/local/lib/node_modules"); err == nil && st.IsDir() {
			add("npm", utils.FilePathGlobLens("/usr/local/lib/node_modules/*/"))
		} else {
			out, err := exec.Command("npm", "root", "-g").Output()
			if err == nil {
				root := strings.TrimSpace(string(out))
				add("npm", utils.FilePathGlobLens(filepath.Join(root, "*/")))
			}
		}
	}
	if utils.CommandExists("pnpm") {
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
		}
		if utils.CommandExists("rpm") {
			add("rpm", utils.GetPkgCount("rpm", "-qa", "--nodigest", "--nosignature"))
		}
		if utils.CommandExists("flatpak") {
			add("flatpak", utils.GetPkgCount("flatpak", "list"))
		}
		if utils.CommandExists("snap") {
			add("snap", utils.GetPkgCount("snap", "list"))
		}
		if utils.CommandExists("brew") {
			root, err := exec.Command("brew", "--cellar").Output()
			if err == nil {
				p1 := strings.TrimSpace(string(root)) + "/*"
				add("brew", utils.FilePathGlobLens(p1))
			}
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
		if utils.CommandExists("winget") {
			add("winget", utils.GetPkgCount("winget", "list", "--accept-source-agreements"))
		}
	}

	// --- Total and formatting ---
	total := 0
	for _, m := range managers {
		total += m.Count
	}

	if len(managers) == 0 {
		return "No packages found"
	}
	res := []string{}

	for _, m := range managers {
		res = append(res, fmt.Sprintf("%d (%s)", m.Count, m.Manager))
	}

	return strings.Join(res, ", ")
}
