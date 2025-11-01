package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func getCPU() string {
	b, err := os.ReadFile("/proc/cpuinfo")
	if err == nil {
		for _, line := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(line, "model name") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					name := strings.TrimSpace(parts[1])
					cores := getCPUCoreCount()
					if cores > 0 {
						return fmt.Sprintf("%s (%d cores)", name, cores)
					}
					return name
				}
			}
		}
	}
	return utils.RunFirst("lscpu")
}

func getCPUCoreCount() int {
	out, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return 0
	}
	return strings.Count(string(out), "processor\t:")
}
