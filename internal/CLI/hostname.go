package cli

import (
	"os"
	"os/exec"
	"strings"
)

func getHostname() string {
	if h, err := os.Hostname(); err == nil {
		return h
	}
	if b, err := exec.Command("hostname").Output(); err == nil {
		return strings.TrimSpace(string(b))
	}
	return "unknown"
}
