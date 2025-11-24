package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PkgInfo struct {
	Manager string
	Count   int
}

func has(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func tot(cmd ...string) int {
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return 0
	}
	lines := bytes.Count(out, []byte{'\n'})
	return lines
}

func totSafe(cmd ...string) int {
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return 0
	}
	lines := bytes.Count(out, []byte{'\n'})
	return lines
}

func dir(glob string) int {
	matches, _ := filepath.Glob(glob)
	return len(matches)
}

func get_packages(osname string) string {
	return countPackages(osname)
}

func countPackages(osname string) string {
	managers := []PkgInfo{}
	add := func(manager string, count int) {
		if count > 0 {
			managers = append(managers, PkgInfo{manager, count})
		}
	}

	// --- Language package managers ---
	if has("pipx") {
		add("pipx", tot("pipx", "list", "--short"))
	}
	if has("pip") {
		add("pip", tot("pip", "freeze"))
	}
	if has("cargo") {
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
	if has("npm") {
		if st, err := os.Stat("/usr/lib/node_modules"); err == nil && st.IsDir() {
			add("npm", dir("/usr/lib/node_modules/*/"))
		} else if st, err := os.Stat("/usr/local/lib/node_modules"); err == nil && st.IsDir() {
			add("npm", dir("/usr/local/lib/node_modules/*/"))
		} else {
			out, err := exec.Command("npm", "root", "-g").Output()
			if err == nil {
				root := strings.TrimSpace(string(out))
				add("npm", dir(filepath.Join(root, "*/")))
			}
		}
	}
	if has("pnpm") {
		p := filepath.Join(os.Getenv("HOME"), ".local/share/pnpm/global/5/node_modules/*/")
		add("pnpm", dir(p))
	}

	// --- System package managers by OS ---
	switch osname {
	case Linux:
		if has("dpkg") {
			add("dpkg", tot("dpkg-query", "-f", ".\\n", "-W"))
		}
		if has("pacman") {
			add("pacman", tot("pacman", "-Qq", "--color", "never"))
		}
		if has("apk") {
			add("apk", tot("apk", "info"))
		}
		if has("rpm") {
			add("rpm", tot("rpm", "-qa", "--nodigest", "--nosignature"))
		}
		if has("flatpak") {
			add("flatpak", tot("flatpak", "list"))
		}
		if has("snap") {
			add("snap", tot("snap", "list"))
		}
		if has("brew") {
			root, err := exec.Command("brew", "--cellar").Output()
			if err == nil {
				p1 := strings.TrimSpace(string(root)) + "/*"
				add("brew", dir(p1))
			}
		}
	case Darwin:
		if has("brew") {
			add("brew", dir("/usr/local/Cellar/*/"))
		}
		if has("port") {
			add("macports", tot("port", "installed"))
		}
	case Windows:
		if has("choco") {
			pdata := os.Getenv("ProgramData")
			if pdata == "" {
				pdata = `C:\ProgramData`
			}
			add("choco", dir(filepath.Join(pdata, "chocolatey/lib/*")))
		}
		if has("winget") {
			add("winget", tot("winget", "list", "--accept-source-agreements"))
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
