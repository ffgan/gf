package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func cleanModelName(model string) string {
	replacements := []string{
		"To be filled by O.E.M.", "To Be Filled", "OEM", "Not Applicable",
		"System Product Name", "System Version", "Undefined", "Default string",
		"Not Specified", "Type1ProductConfigId", "INVALID", "All Series", "�",
	}
	for _, bad := range replacements {
		model = strings.ReplaceAll(model, bad, "")
	}
	model = strings.TrimSpace(model)
	return model
}

func GetModel(osName, kernelName, kernelMachine string) string {
	// osName := getOS()
	// kernelName := getKernelName()
	// kernelMachine := getKernelMachine()
	// TODO: correct this
	osxVersion := "xx"
	var model string

	switch osName {
	case Linux:
		// Empty DMI nodes still exist on many VMs. Skip junk values and
		// keep going through the device tree, WSL, and the CPUID hypervisor.
		model = PickHost(linuxHostInput())
		if strings.TrimSpace(model) == "" {
			virt := utils.RunCommand("systemd-detect-virt")
			if virt != "" && virt != "none" && virt != "container" {
				model = strings.ToUpper(virt)
			}
		}

	case MacOSX, MacOS, RavynOS, Darwin:
		arch := utils.RunCommand("arch")
		isHackintosh := false
		if arch != "arm64" {
			out := utils.RunCommand("kextstat")
			if strings.Contains(out, "FakeSMC") || strings.Contains(out, "VirtualSMC") {
				isHackintosh = true
			}
		}

		if isHackintosh {
			model = fmt.Sprintf("Hackintosh (SMBIOS: %s)", utils.RunCommand("sysctl", "-n", "hw.model"))
		} else {
			if strings.HasPrefix(osxVersion, "10.4") || strings.HasPrefix(osxVersion, "10.5") {
				line := utils.RunCommand("system_profiler", "SPHardwareDataType")
				re := regexp.MustCompile(`Machine Name:\s*(.+)`)
				m := re.FindStringSubmatch(line)
				if len(m) > 1 {
					model = fmt.Sprintf("%s (%s)", strings.TrimSpace(m[1]), utils.RunCommand("sysctl", "-n", "hw.model"))
				}
			} else {
				model = utils.RunCommand("sysctl", "-n", "hw.model")
			}
		}

	case Iphone:
		model = iphoneName(kernelMachine)

	case BSD, MINIX, FreeBSD, OpenBSD, NetBSD, DragonFly, Bitrig:
		if kernelName == "FreeBSD" {
			model = utils.RunCommand("kenv", "smbios.system.version")
		} else {
			model = utils.RunCommand("sysctl", "-n", "hw.vendor", "hw.product")
		}

	case Windows:
		text := readWindowsBIOS()
		if text == "" {
			text = utils.RunCommand("reg.exe", "query", `HKLM\HARDWARE\DESCRIPTION\System\BIOS`)
		}
		model = windowsHardwareModel(text)

	case Solaris, illumos:
		model = utils.RunCommand("prtconf", "-b")
		re := regexp.MustCompile(`banner-name:\s*(.+)`)
		m := re.FindStringSubmatch(model)
		if len(m) > 1 {
			model = strings.TrimSpace(m[1])
		}
		virt := utils.RunCommand("/usr/bin/uname", "-V")
		if virt != "" && virt != "non-virtualized" {
			if model == "" {
				model = utils.RunCommand("uname", "-i")
			}
			model = fmt.Sprintf("%s (%s)", model, virt)
		}

	case AIX:
		model = utils.RunCommand("/usr/bin/uname", "-M")

	case FreeMiNT:
		model = utils.RunCommand("sysctl", "-n", "hw.model")
		model = strings.ReplaceAll(model, "(_MCH *)", "")

	case Interix:
		model = utils.RunCommand("/dev/fs/C/Windows/System32/wbem/WMIC.exe", "computersystem", "get", "manufacturer,model")
		model = strings.ReplaceAll(model, "Manufacturer", "")
		model = strings.ReplaceAll(model, "Model", "")
	default:
		if isBSDFamily(osName) {
			if kernelName == FreeBSD || osName == FreeBSD {
				model = utils.RunCommand("kenv", "smbios.system.version")
			} else {
				model = utils.RunCommand("sysctl", "-n", "hw.vendor", "hw.product")
			}
		}
	}

	model = cleanModelName(model)

	if strings.HasPrefix(model, "Standard PC") {
		model = fmt.Sprintf("KVM/QEMU (%s)", model)
	} else if strings.HasPrefix(model, "OpenBSD") {
		model = fmt.Sprintf("vmm (%s)", model)
	}

	return strings.TrimSpace(model)
}
