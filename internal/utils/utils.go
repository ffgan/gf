package utils

import (
	"os/exec"
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

func RunFirst(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func RunCmdCount(cmd string, args ...string) int {
	c := exec.Command(cmd, args...)
	out, err := c.Output()
	if err != nil {
		return 0
	}
	lines := strings.Count(string(out), "\n")
	if strings.HasSuffix(string(out), "\n") {
		return lines
	}
	return lines + 1
}

func MaxLen(a []string, min int) int {
	if len(a) > min {
		return len(a)
	}
	return min
}
