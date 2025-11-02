package cli

func GetCPU() string {
	return getCPU()
}

func GetDisk() string {
	return "Disk: " + getDisk()
}

// func GetMemory() string {
// 	return "Memory: " + getMemory()
// }

func DetectPackages() string {
	// TODO: 修复与hyfetch不一致的地方
	// Packages: 227 (pip), 2699 (rpm), 25 (flatpak)
	// Packages: 227 (pip), 2699 (rpm), 19 (flatpak-system), 6 (flatpak-user)
	return get_packages()
}

func GetUptime(uptimeShorthand string) string {
	return getUptime(uptimeShorthand)
}

func Geteditor() string {
	return getEditor("off", "on")
}

func GetHostname() string {
	return "Host: " + getHostname()
}

func GetTerm() string {
	return getTerm()
}

func PrintOS() string {
	return "OS: " + getOS()
}

func PrintKernel(osArch, distroShorthand, kernelShorthand, ascii_distro string) string {
	return GetKernel(osArch, distroShorthand, kernelShorthand, ascii_distro)
}

func PrintHost() string {
	return getModel()
}

func PrintDistro(osArch, distroShorthand, ascii_distro string) string {
	// TODO: 修复与hyfetch不一致的地方
	//  OS: Fedora Linux 43 x86_64
	//  OS: Fedora Linux 43 (COSMIC) x86_64
	distro, _ := getDistro(osArch, distroShorthand, ascii_distro)
	return distro
}

func PrintShell(shellPath, shellVersion string) string {
	return getShell(shellPath, shellVersion)
}
