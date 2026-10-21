package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetOS(kernelName, darwin_name string) string {
	_ = darwin_name
	osName, known := ClassifyKernel(kernelName)
	if osName == SunOS {
		return refineSunOS()
	}
	if !known {
		fmt.Fprintf(os.Stderr, "gf: unknown OS %q, showing the kernel name\n", kernelName)
	}
	return osName
}

// refineSunOS separates Solaris from illumos. uname -s is SunOS for both;
// uname -o carries the distribution (illumos, SmartOS, OmniOS, Tribblix, …).
func refineSunOS() string {
	out, err := utils.CommandOutput("uname", "-o")
	if err != nil {
		return Solaris
	}
	low := strings.ToLower(strings.TrimSpace(string(out)))
	switch {
	case low == "", low == "sunos", strings.Contains(low, "solaris"):
		return Solaris
	case strings.Contains(low, "illumos"),
		strings.Contains(low, "smartos"),
		strings.Contains(low, "omnios"),
		strings.Contains(low, "tribblix"),
		strings.Contains(low, "openindiana"),
		strings.Contains(low, "dilos"),
		strings.Contains(low, "napp-it"),
		strings.Contains(low, "xstreamos"):
		return illumos
	default:
		return Solaris
	}
}
