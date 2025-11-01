package cli

func GetCPU() string {
	return "CPU: " + getCPU()
}

func GetDisk() string {
	return "Disk: " + getDisk()
}

func GetMemory() string {
	return "Memory: " + getMemory()
}

func DetectPackages() string {
	return "Packages: " + get_packages()
}

func GetUptime(uptimeShorthand string) string {
	return "Uptime: " + getUptime(uptimeShorthand)
}

func Geteditor() string {
	return "Editor: " + getEditor("off", "on")
}

func GetHostname() string {
	return "Host: " + getHostname()
}

func GetTerm() string {
	return "Terminal: " + getTerm()
}

func PrintOS() string {
	return "OS: " + getOS()
}

func PrintKernel(osArch, distroShorthand, kernelShorthand, ascii_distro string) string {
	return "Kernel: " + GetKernel(osArch, distroShorthand, kernelShorthand, ascii_distro)
}

func PrintHost() string {
	return "Host: " + getModel()
}

func PrintDistro(osArch, distroShorthand, ascii_distro string) string {
	distro, _ := getDistro(osArch, distroShorthand, ascii_distro)
	return "OS: " + distro
}

func PrintShell(shellPath, shellVersion string) string {
	return "Shell: " + getShell(shellPath, shellVersion)
}
