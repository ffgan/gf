package cli

import (
	"os/exec"
	"strings"
)

func getDisk() string {
	out, err := exec.Command("bash", "-lc", "df -h / | tail -1 | awk '{print $3\"/\"$2\" (\"$5\")\"}'").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	return "unknown"
}
