package cli

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

type kernelconfig struct {
	osArch, distroShorthand, kernelShorthand, ascii_distro, osName, kernelVersion, kernelName, KernelMachine string
}

func GetKernel(osArch, distroShorthand, kernelShorthand, ascii_distro, osName, kernelVersion, kernelName, KernelMachine string) string {
	kc := kernelconfig{
		osArch:          osArch,
		distroShorthand: distroShorthand,
		kernelShorthand: kernelShorthand,
		ascii_distro:    ascii_distro,
		osName:          osName,
		kernelVersion:   kernelVersion,
		kernelName:      kernelName,
		KernelMachine:   KernelMachine,
	}
	return getKernel(kc)
}

func getKernel(kc kernelconfig) string {
	var kernel string
	// osName := getOS()
	// kernelVersion := getKernelVersion()
	// Skip for integrated systems
	if matched, _ := regexp.MatchString(`AIX|IRIX`, kc.osName); matched {
		return kernel
	}

	switch kc.osName {
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
		kernel = kc.kernelVersion + " " +
			kc.KernelMachine + " " +
			getKernelRelease()
		return kernel
	}
	// kernelName := getKernelName()

	switch kc.kernelShorthand {
	case utils.ON:
		kernel = kc.kernelVersion
	case utils.OFF:
		kernel = kc.kernelName + " " + kc.kernelVersion
	}

	// Hide kernel info if identical to distro
	distro, _ := GetDistro(kc.osName, kc.osArch, kc.KernelMachine, kc.distroShorthand, kc.ascii_distro)
	if matched, _ := regexp.MatchString(`BSD|MINIX`, kc.osName); matched && strings.Contains(distro, kc.kernelName) {
		if kc.distroShorthand == "on" || kc.distroShorthand == "tiny" {
			kernel = kc.kernelVersion
		} else {
			kernel = ""
		}
	}
	return kernel
}

// func getKernelVersion() string {
// 	return UName("-r")
// }

// func GetKernelMachine() string {
// 	return getKernelMachine()
// }

// func getKernelMachine() string {
// 	return UName("-m")
// }

// func getKernelName() string {
// 	return UName("-s")
// }

func getKernelRelease() string {
	return UName("-v")
}

func UName(command string) string {
	out, err := exec.Command("uname", command).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run uname %s: %v\n", command, err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(out))
}
