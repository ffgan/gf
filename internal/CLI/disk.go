package cli

import (
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func getDisk() string {
	out, err := utils.CommandOutput("bash", "-lc", "df -h / | tail -1 | awk '{print $3\"/\"$2\" (\"$5\")\"}'")
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	return "unknown"
}
