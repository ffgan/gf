package cli

import (
	"fmt"
	"regexp"
	"strings"
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

func getModel() string {
	osName := getOS()
	kernelName := getKernelName()
	kernelMachine := getKernelMachine()
	// TODO: correct this
	osxVersion := "xx"
	var model string

	switch osName {
	case Linux:
		if pathExists("/android/system/") || pathExists("/system/app/") {
			model = fmt.Sprintf("%s %s",
				runCmd("getprop", "ro.product.brand"),
				runCmd("getprop", "ro.product.model"))
		} else if pathExists("/sys/devices/virtual/dmi/id/product_name") ||
			pathExists("/sys/devices/virtual/dmi/id/product_version") {
			model = readFileTrim("/sys/devices/virtual/dmi/id/product_name") + " " +
				readFileTrim("/sys/devices/virtual/dmi/id/product_version")
		} else if pathExists("/sys/devices/virtual/dmi/id/board_vendor") ||
			pathExists("/sys/devices/virtual/dmi/id/board_name") {
			model = readFileTrim("/sys/devices/virtual/dmi/id/board_vendor") + " " +
				readFileTrim("/sys/devices/virtual/dmi/id/board_name")
		} else if pathExists("/sys/firmware/devicetree/base/model") {
			model = readFileTrim("/sys/firmware/devicetree/base/model")
		} else if pathExists("/tmp/sysinfo/model") {
			model = readFileTrim("/tmp/sysinfo/model")
		}

	case MacOSX, MacOS, RavynOS:
		arch := runCmd("arch")
		isHackintosh := false
		if arch != "arm64" {
			out := runCmd("kextstat")
			if strings.Contains(out, "FakeSMC") || strings.Contains(out, "VirtualSMC") {
				isHackintosh = true
			}
		}

		if isHackintosh {
			model = fmt.Sprintf("Hackintosh (SMBIOS: %s)", runCmd("sysctl", "-n", "hw.model"))
		} else {
			if strings.HasPrefix(osxVersion, "10.4") || strings.HasPrefix(osxVersion, "10.5") {
				line := runCmd("system_profiler", "SPHardwareDataType")
				re := regexp.MustCompile(`Machine Name:\s*(.+)`)
				m := re.FindStringSubmatch(line)
				if len(m) > 1 {
					model = fmt.Sprintf("%s (%s)", strings.TrimSpace(m[1]), runCmd("sysctl", "-n", "hw.model"))
				}
			} else {
				model = runCmd("sysctl", "-n", "hw.model")
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
			model = runCmd("kenv", "smbios.system.version")
		} else {
			model = runCmd("sysctl", "-n", "hw.vendor", "hw.product")
		}

	case Windows:
		model = runCmd("wmic", "computersystem", "get", "manufacturer,model")
		model = strings.ReplaceAll(model, "Manufacturer", "")
		model = strings.ReplaceAll(model, "Model", "")

	case Solaris, illumos:
		model = runCmd("prtconf", "-b")
		re := regexp.MustCompile(`banner-name:\s*(.+)`)
		m := re.FindStringSubmatch(model)
		if len(m) > 1 {
			model = strings.TrimSpace(m[1])
		}
		virt := runCmd("/usr/bin/uname", "-V")
		if virt != "" && virt != "non-virtualized" {
			if model == "" {
				model = runCmd("uname", "-i")
			}
			model = fmt.Sprintf("%s (%s)", model, virt)
		}

	case AIX:
		model = runCmd("/usr/bin/uname", "-M")

	case FreeMiNT:
		model = runCmd("sysctl", "-n", "hw.model")
		model = strings.ReplaceAll(model, "(_MCH *)", "")

	case Interix:
		model = runCmd("/dev/fs/C/Windows/System32/wbem/WMIC.exe", "computersystem", "get", "manufacturer,model")
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
