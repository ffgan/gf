package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	args := append([]string{}, psFlags...)
	out := RunCommand("ps", args...)
	for _, n := range names {
		re := regexp.MustCompile(fmt.Sprintf(`(?m)^%s$`, regexp.QuoteMeta(n)))
		if re.FindString(out) != "" {
			return n
		}
	}
	return ""
}
