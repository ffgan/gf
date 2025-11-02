package cli

import (
	"fmt"
	"os"
	"strings"
)

func GetDistro(osArch, distroShorthand, ascii_distro string) (string, string) {
	return getDistro(osArch, distroShorthand, ascii_distro)
}

func getDistro(osArch, distroShorthand, ascii_distro string) (string, string) {
	var distro string

	osName := getOS()
	kernelMachine := getKernelMachine()
	kernelVersion := getKernelVersion()

	switch osName {
	case "Linux", "BSD", "MINIX", "Ironclad":
		switch {
		case fileExists("/etc/os-release"):
			// Try to parse /etc/os-release
			content, _ := os.ReadFile("/etc/os-release")
			lines := strings.Split(string(content), "\n")
			info := map[string]string{}
			for _, line := range lines {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					info[parts[0]] = strings.Trim(parts[1], `"`)
				}
			}
			name := info["NAME"]
			version := info["VERSION_ID"]
			switch distroShorthand {
			case "on":
				distro = fmt.Sprintf("%s %s", name, version)
			case "tiny":
				distro = name
			default:
				distro = fmt.Sprintf("%s %s", name, version)
			}
		case fileExists("/etc/debian_version"):
			ver := readFirstLine("/etc/debian_version")
			distro = "Debian " + ver
		case fileExists("/etc/redhat-release"):
			ver := readFirstLine("/etc/redhat-release")
			distro = "Red Hat " + ver
		default:
			distro = "Linux"
		}
	case "darwin":
		out := RunCommand("sw_vers", "-productVersion")
		codename := "macOS"
		version := out
		switch {
		case strings.HasPrefix(out, "10.15"):
			codename = "macOS Catalina"
		case strings.HasPrefix(out, "11"):
			codename = "macOS Big Sur"
		case strings.HasPrefix(out, "12"):
			codename = "macOS Monterey"
		case strings.HasPrefix(out, "13"):
			codename = "macOS Ventura"
		case strings.HasPrefix(out, "14"):
			codename = "macOS Sonoma"
		}
		distro = fmt.Sprintf("%s %s", codename, version)
	case "windows":
		// For Windows, use PowerShell to query ProductName
		psCmd := `Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion' | Select-Object -ExpandProperty ProductName`
		out := RunCommand("powershell", "-Command", psCmd)
		build := RunCommand("powershell", "-Command", "Get-ItemProperty 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion' | Select-Object -ExpandProperty CurrentBuildNumber")
		distro = strings.TrimSpace(out)
		if strings.Contains(distro, "Windows 10") {
			if buildInt := strings.TrimSpace(build); buildInt >= "22000" {
				distro = strings.Replace(distro, "Windows 10", "Windows 11", 1)
			}
		}
	case "solaris":
		distro = readFirstLine("/etc/release")
	case "haiku":
		distro = "Haiku"
	case "aix":
		distro = "AIX " + kernelVersion
	default:
		distro = osName + " " + kernelVersion
	}

	if osArch == "on" {
		distro += " " + kernelMachine
	}

	if ascii_distro == "auto" {
		ascii_distro = strings.TrimSpace(distro)
	}

	return distro, ascii_distro
}
