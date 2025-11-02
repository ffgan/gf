package cli

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func GetKernel(osArch, distroShorthand, kernelShorthand, ascii_distro string) string {
	return getKernel(osArch, distroShorthand, kernelShorthand, ascii_distro)
}

func getKernel(osArch, distroShorthand, kernelShorthand, ascii_distro string) string {
	var kernel string
	osName := getOS()
	kernelVersion := getKernelVersion()
	// Skip for integrated systems
	if matched, _ := regexp.MatchString(`AIX|IRIX`, osName); matched {
		return kernel
	}

	switch osName {
	case Haiku:
		return getKernelRelease()

	case Windows:
		out, _ := exec.Command("wmic", "os", "get", "Version").Output()
		// remove "Version" header and trim
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) > 1 {
			kernel = strings.TrimSpace(lines[1])
		}

		return kernel

	case Solaris:
		ver := getKernelRelease()
		if matched, _ := regexp.MatchString(`^11\.[0123]$`, ver); !matched {
			kernel = ver
		}
		return kernel

	case Interix:
		kernel = getKernelVersion() + " " +
			getKernelMachine() + " " +
			getKernelRelease()
		return kernel
	}
	kernelName := getKernelName()
	// Default case
	if kernelShorthand == "on" {
		kernel = kernelVersion
	} else if kernelShorthand == "off" {
		kernel = kernelName + " " + kernelVersion
	}

	// Hide kernel info if identical to distro
	distro, _ := getDistro(osArch, distroShorthand, ascii_distro)
	if matched, _ := regexp.MatchString(`BSD|MINIX`, osName); matched && strings.Contains(distro, kernelName) {
		if distroShorthand == "on" || distroShorthand == "tiny" {
			kernel = kernelVersion
		} else {
			kernel = ""
		}
	}
	return kernel
}

func getKernelVersion() string {
	return uname("-r")
}

func GetKernelMachine() string {
	return getKernelMachine()
}

func getKernelMachine() string {
	return uname("-m")
}

func getKernelName() string {
	return uname("-s")
}

func getKernelRelease() string {
	return uname("-v")
}

func uname(command string) string {
	out, err := exec.Command("uname", command).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run uname %s: %v\n", command, err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(out))
}
