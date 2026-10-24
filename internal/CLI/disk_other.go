//go:build !windows

package cli

import (
	"fmt"
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
	// -k -P is POSIX. GNU df accepts -B1, but macOS df rejects it and the
	// disk line disappears. -k blocks are 1024 bytes.
	out, err := utils.CommandOutput("df", "-k", "-P", path)
	if err != nil {
		return 0, 0, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return 0, 0, 0, fmt.Errorf("unexpected df output")
	}
	used, total, avail, percent, err := parseDFPOSIX(lines[len(lines)-1])
	if err != nil {
		return 0, 0, 0, err
	}
	used *= 1024
	total *= 1024
	avail *= 1024
	if adjusted, pct, ok := applyContainerUsage(path, used, total, avail); ok {
		return adjusted, total, pct, nil
	}
	return used, total, percent, nil
}
