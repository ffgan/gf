package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetOS() string {
	return getOS()
}

// getOS replicates the logic of the Bash get_os() function.
func getOS() string {
	kernelName := getKernelName()
	var osName string

	switch {
	case kernelName == "Darwin":
		osName = getDarwinName()

	case kernelName == "SunOS":
		out, err := exec.Command("uname", "-o").Output()
		if err == nil && strings.Contains(string(out), "illumos") {
			osName = "illumos"
		} else {
			osName = "Solaris"
		}

	case kernelName == "Haiku":
		osName = "Haiku"

	case kernelName == "MINIX":
		osName = "MINIX"

	case kernelName == "AIX":
		osName = "AIX"

	case strings.HasPrefix(kernelName, "IRIX"):
		osName = "IRIX"

	case kernelName == "FreeMiNT":
		osName = "FreeMiNT"

	case kernelName == "Interix":
		osName = "Interix"

	case kernelName == "Ironclad":
		osName = "Ironclad"

	case kernelName == "OSF1":
		osName = "digitalUNIX"

	case kernelName == "Linux" || strings.HasPrefix(kernelName, "GNU"):
		osName = "Linux"

	case strings.HasSuffix(kernelName, "BSD") ||
		kernelName == "DragonFly" ||
		kernelName == "Bitrig":
		osName = "BSD"

	case strings.HasPrefix(kernelName, "CYGWIN") ||
		strings.HasPrefix(kernelName, "MSYS") ||
		strings.HasPrefix(kernelName, "MINGW") ||
		kernelName == "Windows_NT":
		osName = "Windows"

	default:
		fmt.Fprintf(os.Stderr, "Unknown OS detected: '%s', aborting...\n", kernelName)
		fmt.Fprintln(os.Stderr, "Open an issue on GitHub to add support for your OS.")
		os.Exit(1)
	}

	return osName
}

// getDarwinName retrieves a human-friendly macOS name (like "macOS" or "Darwin")
func getDarwinName() string {
	out, err := exec.Command("sw_vers", "-productName").Output()
	if err != nil {
		return "Darwin"
	}
	return strings.TrimSpace(string(out))
}
