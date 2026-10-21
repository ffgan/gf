package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetDistro(osName, os_arch, kernel_machine, distroShorthand, ascii_distro string) (string, string) {
	var distro string
	var asciiName string

	switch osName {
	case Linux, BSD, MINIX, Ironclad:
		switch {
		case utils.FileExists("/etc/os-release"):
			info := parseOSRelease("/etc/os-release")
			name := info["NAME"]
			if name == "" {
				name = info["ID"]
			}
			version := info["VERSION"]
			if version == "" {
				version = info["VERSION_ID"]
			}
			pretty := info["PRETTY_NAME"]
			switch distroShorthand {
			case "tiny":
				distro = name
			case utils.ON:
				distro = strings.TrimSpace(name + " " + version)
			default:
				if pretty != "" {
					distro = pretty
				} else {
					distro = strings.TrimSpace(name + " " + version)
				}
			}
			asciiName = firstNonEmpty(info["NAME"], info["ID"], pretty)
		case utils.FileExists("/etc/debian_version"):
			ver := utils.ReadFirstLine("/etc/debian_version")
			distro = "Debian " + ver
		case utils.FileExists("/etc/redhat-release"):
			ver := utils.ReadFirstLine("/etc/redhat-release")
			distro = "Red Hat " + ver
		default:
			distro = "Linux"
		}
	case Darwin:
		out := utils.RunCommand("sw_vers", "-productVersion")
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
		out := utils.RunCommand("powershell", "-Command", psCmd)
		build := utils.RunCommand("powershell", "-Command", "Get-ItemProperty 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion' | Select-Object -ExpandProperty CurrentBuildNumber")
		distro = strings.TrimSpace(out)
		if strings.Contains(distro, "Windows 10") {
			if buildInt := strings.TrimSpace(build); buildInt >= "22000" {
				distro = strings.Replace(distro, "Windows 10", "Windows 11", 1)
			}
		}
	case Solaris:
		distro = utils.ReadFirstLine("/etc/release")
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

	switch ascii_distro {
	case "", utils.AUTO:
		if asciiName != "" {
			ascii_distro = asciiName
		} else {
			ascii_distro = strings.TrimSpace(distro)
		}
	}

	return strings.TrimSpace(distro), strings.TrimSpace(ascii_distro)
}

func parseOSRelease(path string) map[string]string {
	content, _ := os.ReadFile(path)
	info := map[string]string{}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		info[parts[0]] = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	}
	return info
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
