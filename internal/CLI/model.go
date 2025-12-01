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
		if utils.PathExists("/android/system/") || utils.PathExists("/system/app/") {
			model = fmt.Sprintf("%s %s",
				utils.RunCommand("getprop", "ro.product.brand"),
				utils.RunCommand("getprop", "ro.product.model"))
		} else if utils.PathExists("/sys/devices/virtual/dmi/id/product_name") ||
			utils.PathExists("/sys/devices/virtual/dmi/id/product_version") {
			model = utils.ReadFileTrim("/sys/devices/virtual/dmi/id/product_name") + " " +
				utils.ReadFileTrim("/sys/devices/virtual/dmi/id/product_version")
		} else if utils.PathExists("/sys/devices/virtual/dmi/id/board_vendor") ||
			utils.PathExists("/sys/devices/virtual/dmi/id/board_name") {
			model = utils.ReadFileTrim("/sys/devices/virtual/dmi/id/board_vendor") + " " +
				utils.ReadFileTrim("/sys/devices/virtual/dmi/id/board_name")
		} else if utils.PathExists("/sys/firmware/devicetree/base/model") {
			model = utils.ReadFileTrim("/sys/firmware/devicetree/base/model")
		} else if utils.PathExists("/tmp/sysinfo/model") {
			model = utils.ReadFileTrim("/tmp/sysinfo/model")
		}

	case MacOSX, MacOS, RavynOS:
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
		// Simplified mapping (full table omitted for brevity)
		switch kernelMachine {
		case "iPhone15,2":
			model = "iPhone 14 Pro"
		case "iPhone15,3":
			model = "iPhone 14 Pro Max"
		default:
			model = kernelMachine
		}

	case BSD, MINIX:
		if kernelName == "FreeBSD" {
			model = utils.RunCommand("kenv", "smbios.system.version")
		} else {
			model = utils.RunCommand("sysctl", "-n", "hw.vendor", "hw.product")
		}

	case Windows:
		model = utils.RunCommand("wmic", "computersystem", "get", "manufacturer,model")
		model = strings.ReplaceAll(model, "Manufacturer", "")
		model = strings.ReplaceAll(model, "Model", "")

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
	}

	model = cleanModelName(model)

	if strings.HasPrefix(model, "Standard PC") {
		model = fmt.Sprintf("KVM/QEMU (%s)", model)
	} else if strings.HasPrefix(model, "OpenBSD") {
		model = fmt.Sprintf("vmm (%s)", model)
	}

	return strings.TrimSpace(model)
}
