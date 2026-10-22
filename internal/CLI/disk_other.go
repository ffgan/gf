//go:build !windows

package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func localMount(mount string) string {
	return mount
}

func windowsDiskReport(mounts []string, report DiskReport) (string, bool) {
	return "", false
}

func dfBytes(path string) (int64, int64, int, error) {
	out, err := utils.CommandOutput("df", "-B1", "-P", path)
	if err != nil {
		return 0, 0, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return 0, 0, 0, fmt.Errorf("unexpected df output")
	}
	fields := strings.Fields(lines[len(lines)-1])
	if len(fields) < 4 {
		return 0, 0, 0, fmt.Errorf("unexpected df columns")
	}
	total, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	used, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	percent := 0
	if len(fields) >= 5 && strings.HasSuffix(fields[4], "%") {
		percent, _ = strconv.Atoi(strings.TrimSuffix(fields[4], "%"))
	} else if total > 0 {
		percent = int(used * 100 / total)
	}
	return used, total, percent, nil
}
