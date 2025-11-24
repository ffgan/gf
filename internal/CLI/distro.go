package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetDistro(osName, os_arch, kernel_machine, distroShorthand, ascii_distro string) (string, string) {
	var distro string

	switch osName {
	case Linux, BSD, MINIX, Ironclad:
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
	case Darwin:
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
	case Windows:
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
	case Solaris:
		distro = readFirstLine("/etc/release")
	case Haiku:
		distro = Haiku
		// case AIX:
		// 	distro = AIX + " " + kernelVersion
		// default:
		// 	distro = osName + " " + kernelVersion
	}
	var machine_arch string
	switch osName {
	case Solaris, illumos, AIX, Haiku, IRIX, FreeMiNT, BSD, digitalUNIX:
		machine_arch = UName("-p")
	default:
		machine_arch = kernel_machine
	}

	if os_arch == utils.ON {
		distro += " " + machine_arch
	}

	if ascii_distro == utils.AUTO {
		ascii_distro = strings.TrimSpace(distro)
	}

	return distro, ascii_distro
}
