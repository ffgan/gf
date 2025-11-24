package cli

import "github.com/ffgan/gf/internal/utils"

func GetCPU(OS, Machine, ShowCores, ShowSpeed, ShowTemp, ShowBrand, CoreType, SpeedShort, TempUnit string) string {
	var showCores bool
	var showSpeed bool
	var showTemp bool
	var showBrand bool
	var speedShort bool
	var tempUnit string // "C" or "F"

	if ShowCores == utils.ON {
		showCores = true
	}

	if ShowSpeed == utils.ON {
		showSpeed = true
	}

	if ShowTemp == utils.ON {
		showTemp = true
	}

	if ShowBrand == utils.ON {
		showBrand = true
	}

	if SpeedShort == utils.ON {
		speedShort = true
	}

	config := CPUConfig{
		ShowCores:  showCores,
		ShowSpeed:  showSpeed,
		ShowTemp:   showTemp,
		ShowBrand:  showBrand,
		CoreType:   CoreType,
		SpeedShort: speedShort,
		TempUnit:   tempUnit,
	}

	return getCPU(OS, Machine, config)
}

func GetDisk() string {
	return "Disk: " + getDisk()
}

func DetectPackages(osname string) string {
	// TODO: 修复与hyfetch不一致的地方
	// Packages: 227 (pip), 2699 (rpm), 25 (flatpak)
	// Packages: 227 (pip), 2699 (rpm), 19 (flatpak-system), 6 (flatpak-user)
	return get_packages(osname)
}

func GetUptime(osName, uptimeShorthand string) string {
	return getUptime(osName, uptimeShorthand)
}

func Geteditor() string {
	return GetEditor("off", "on")
}

func GetHostname() string {
	return "Host: " + getHostname()
}

// func GetTerm() string {
// 	return getTerm()
// }

// func PrintOS() string {
// 	return "OS: " + getOS()
// }

var PrintKernel = GetKernel

// func PrintHost() string {
// 	return GetModel()
// }

// func PrintDistro(osArch, distroShorthand, ascii_distro string) string {
// 	// TODO: 修复与hyfetch不一致的地方
// 	//  OS: Fedora Linux 43 x86_64
// 	//  OS: Fedora Linux 43 (COSMIC) x86_64
// 	distro, _ := getDistro(osArch, distroShorthand, ascii_distro)
// 	return distro
// }

func PrintShell(shellPath, shellVersion string) string {
	return GetShell(shellPath, shellVersion)
}
