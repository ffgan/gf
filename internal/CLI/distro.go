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
	var releaseName, releaseID, releaseLike string

	switch osName {
	case Linux, BSD, MINIX, Ironclad:
		switch {
		case utils.FileExists("/etc/os-release"):
			info := parseOSRelease("/etc/os-release")
			releaseName = info["NAME"]
			releaseID = info["ID"]
			releaseLike = info["ID_LIKE"]
			name := releaseName
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
			if key := LogoKey(releaseName, releaseID); key != "" {
				asciiName = key
			}
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
		distro = fmt.Sprintf("%s %s", macOSCodename(out), out)
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
	}
	if strings.TrimSpace(distro) == "" {
		distro = osName
	}
	if osName == Linux {
		distro = AnnotateDistro(strings.TrimSpace(distro), RuntimeTags(RuntimeSig{
			ID:     releaseID,
			IDLike: releaseLike,
			Name:   releaseName,
			Kernel: readSmall("/proc/version"),
			Env: map[string]string{
				"TERMUX_VERSION": os.Getenv("TERMUX_VERSION"),
				"PREFIX":         os.Getenv("PREFIX"),
				"ANDROID_ROOT":   os.Getenv("ANDROID_ROOT"),
				"ANDROID_DATA":   os.Getenv("ANDROID_DATA"),
			},
			Exists: utils.PathExists,
			Read:   readSmall,
		}))
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

func macOSCodename(version string) string {
	switch {
	case strings.HasPrefix(version, "10.15"):
		return "macOS Catalina"
	case strings.HasPrefix(version, "11"):
		return "macOS Big Sur"
	case strings.HasPrefix(version, "12"):
		return "macOS Monterey"
	case strings.HasPrefix(version, "13"):
		return "macOS Ventura"
	case strings.HasPrefix(version, "14"):
		return "macOS Sonoma"
	case strings.HasPrefix(version, "15"):
		return "macOS Sequoia"
	case strings.HasPrefix(version, "26"):
		return "macOS Tahoe"
	default:
		return "macOS"
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
