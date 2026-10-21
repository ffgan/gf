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
	if p, ok := parentPIDFast(pid); ok {
		return p
	}
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
	if name, ok := processNameFast(pid); ok {
		return name
	}
	out, err := CommandOutput("ps", "-p", pid, "-o", "comm=")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// ProcessNames returns the comm of every process. On Linux this is one
// pass over /proc, shared by every module that used to spawn `ps`.
func ProcessNames() ([]string, bool) {
	return processNamesFast()
}

func ScanProcesses(psFlags []string, names []string) string {
	if comms, ok := processNamesFast(); ok {
		running := make(map[string]struct{}, len(comms))
		for _, c := range comms {
			running[c] = struct{}{}
		}
		for _, n := range names {
			if _, found := running[n]; found {
				return n
			}
		}
		return ""
	}
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
