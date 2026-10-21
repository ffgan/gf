package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ON   = "on"
	OFF  = "off"
	AUTO = "auto"
)

func Unquote(s string) string {
	return strings.Trim(s, `"'`)
}

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ReadFirstLine(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

func ReadFileTrim(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func FilePathGlobLens(glob string) int {
	matches, _ := filepath.Glob(glob)
	return len(matches)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetParentPID(pid int) int {
	out, err := CommandOutput("ps", "-o", "ppid=", "-p", fmt.Sprint(pid))
	if err != nil {
		return 0
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0
	}
	var p int
	fmt.Sscanf(s, "%d", &p)
	return p
}

func GetProcessName(pid string) string {
	out, err := CommandOutput("ps", "-p", pid, "-o", "comm=")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func ScanProcesses(psFlags []string, names []string) string {
	out := RunCommand("ps", "-e", "-o", "comm=")
	if strings.TrimSpace(out) == "" && len(psFlags) > 0 {
		out = RunCommand("ps", psFlags...)
	}
	found := map[string]bool{}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			found[line] = true
			if fields := strings.Fields(line); len(fields) > 0 {
				found[fields[len(fields)-1]] = true
			}
		}
	}
	for _, n := range names {
		if found[n] {
			return n
		}
		if len(n) > 15 && found[n[:15]] {
			return n
		}
	}
	return ""
}
