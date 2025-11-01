package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getMemory() string {
	b, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "unknown"
	}
	var memTotal, memFree int64
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			memTotal, _ = strconv.ParseInt(fields[1], 10, 64)
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			memFree, _ = strconv.ParseInt(fields[1], 10, 64)
		}
	}
	if memTotal == 0 {
		return "unknown"
	}
	used := (memTotal - memFree) / 1024
	total := memTotal / 1024
	return fmt.Sprintf("%dMiB / %dMiB", used, total)
}
