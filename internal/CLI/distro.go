package cli

import (
	"os"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetDistro(osName, os_arch, kernel_machine, distroShorthand, ascii_distro string) (string, string) {
	var distro string
	var asciiName string
	var releaseName, releaseID, releaseLike string
	hideArch := false
	files := collectReleaseFiles()
	enrichRelease(files)

	switch osName {
	case Linux, BSD, MINIX, Ironclad, FreeBSD, OpenBSD, NetBSD, DragonFly, Bitrig, "Fiwix", "QNX":
		distro, asciiName, releaseName, releaseID, releaseLike = FormatDistro(osName, distroShorthand, files)
		if isBSDFamily(osName) && (distro == osName || distro == "") {
			if rel := strings.TrimSpace(UName("-r")); rel != "" && !strings.Contains(distro, rel) {
				distro = strings.TrimSpace(osName + " " + rel)
			}
		}
	case Darwin:
		product := strings.TrimSpace(utils.RunCommand("sw_vers", "-productName"))
		version := strings.TrimSpace(utils.RunCommand("sw_vers", "-productVersion"))
		build := strings.TrimSpace(utils.RunCommand("sw_vers", "-buildVersion"))
		distro, asciiName, hideArch = formatApple(product, version, build, distroShorthand)
	case Windows:
		psCmd := `Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion' | Select-Object -ExpandProperty ProductName`
		out := utils.RunCommand("powershell", "-Command", psCmd)
		build := utils.RunCommand("powershell", "-Command", "Get-ItemProperty 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion' | Select-Object -ExpandProperty CurrentBuildNumber")
		distro = windowsProduct(out, build)
		asciiName = windowsArt(distro)
	case Solaris, illumos:
		distro, asciiName, releaseName, releaseID, releaseLike = FormatDistro(osName, distroShorthand, files)
	case Haiku:
		distro, asciiName, releaseName, releaseID, releaseLike = FormatDistro(osName, distroShorthand, files)
	case AIX:
		level := strings.TrimSpace(utils.RunCommand("oslevel"))
		switch {
		case level == "":
			distro = "AIX"
		case strings.HasPrefix(strings.ToLower(level), "aix"):
			distro = level
		default:
			distro = "AIX " + level
		}
		asciiName = "AIX"
	case IRIX:
		distro = strings.TrimSpace("IRIX " + utils.RunCommand("uname", "-r"))
		asciiName = "IRIX"
	case Interix:
		distro = strings.TrimSpace("Interix " + utils.RunCommand("uname", "-r"))
		asciiName = "Interix"
	case digitalUNIX:
		distro = strings.TrimSpace("digital UNIX " + utils.RunCommand("uname", "-r"))
		asciiName = "digitalUNIX"
	default:
		if isBSDFamily(osName) {
			distro, asciiName, releaseName, releaseID, releaseLike = FormatDistro(osName, distroShorthand, files)
			if distro == osName || distro == "" {
				if rel := strings.TrimSpace(UName("-r")); rel != "" && !strings.Contains(distro, rel) {
					distro = strings.TrimSpace(osName + " " + rel)
				}
			}
		}
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
		if isBSDFamily(osName) {
			machine_arch = UName("-p")
		} else {
			machine_arch = kernel_machine
		}
	}

	if os_arch == utils.ON && !hideArch {
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
	return parseOSReleaseText(readSmall(path))
}

func collectReleaseFiles() map[string]string {
	files := map[string]string{}
	put := func(key, path string, keepEmpty bool) {
		if _, ok := files[key]; ok || !utils.FileExists(path) {
			return
		}
		content := readSmall(path)
		if strings.TrimSpace(content) == "" && !keepEmpty {
			return
		}
		files[key] = content
	}
	put("os-release", "/etc/os-release", false)
	put("os-release", "/usr/lib/os-release", false)
	put("os-release", "/etc/initrd-release", false)
	put("lsb-release", "/etc/lsb-release", false)
	put("openwrt_release", "/etc/openwrt_release", false)
	put("alpine-release", "/etc/alpine-release", false)
	put("artix-release", "/etc/artix-release", true)
	put("arch-release", "/etc/arch-release", true)
	put("gentoo-release", "/etc/gentoo-release", false)
	put("slackware-version", "/etc/slackware-version", false)
	put("suse-release", "/etc/SuSE-release", false)
	put("fedora-release", "/etc/fedora-release", false)
	put("centos-release", "/etc/centos-release", false)
	put("redhat-release", "/etc/redhat-release", false)
	put("system-release", "/etc/system-release", false)
	put("debian_version", "/etc/debian_version", false)
	put("bedrock", "/bedrock/etc/bedrock-release", false)
	put("build.prop", "/system/build.prop", false)
	put("os-release", "/system/etc/os-release", false)
	put("haiku-version", "/boot/system/version", false)
	put("solaris-release", "/etc/release", false)
	put("smartos-product", "/etc/product", false)
	put("armbian-release", "/etc/armbian-release", false)
	put("redstar-release", "/etc/redstar-release", false)
	put("siduction-version", "/etc/siduction-version", false)
	put("gobo-version", "/etc/GoboLinuxVersion", false)
	put("sde-version", "/etc/SDE-VERSION", false)
	put("mcst-version", "/etc/mcst_version", false)
	put("angstrom-version", "/etc/angstrom-version", false)
	put("frugalware-release", "/etc/frugalware-release", false)
	put("pardus-release", "/etc/pardus-release", false)
	put("mageia-release", "/etc/mageia-release", false)
	put("mandriva-release", "/etc/mandriva-release", false)
	put("mandriva-release", "/etc/mandrake-release", false)
	put("meego-release", "/etc/meego-release", false)
	put("sabayon-release", "/etc/sabayon-release", false)
	put("bluewhite64-version", "/etc/bluewhite64-version", false)
	put("crux-release", "/etc/crux-release", false)
	put("tizen-release", "/etc/tizen-release", false)
	put("krass-version", "/etc/KrassOS.version", false)
	put("manjaro-release", "/etc/manjaro-release", false)
	put("q4os-version", "/etc/q4os_version", false)
	put("nixos-marker", "/etc/NIXOS", true)
	put("linspire-version", "/etc/linspire-version", false)
	put("lindows-version", "/etc/lindowsos-version", false)
	put("astra-version", "/etc/astra_version", false)
	put("vzlinux-release", "/etc/vzlinux-release", false)
	put("slitaz-release", "/etc/slitaz-release", false)
	put("opnsense-version", "/usr/local/opnsense/version/opnsense", false)
	put("synology-version", "/etc/VERSION", false)
	put("floflis-release", "/etc/floflis-release", true)
	put("qnx-props", "/pps/services/deviceproperties", false)
	put("qnx-version", "/base/etc/os.version", false)
	put("issue", "/etc/issue", false)
	return files
}

func macOSCodename(version string) string {
	switch {
	case strings.HasPrefix(version, "10.16"):
		return "macOS Big Sur"
	case strings.HasPrefix(version, "10.15"):
		return "macOS Catalina"
	case strings.HasPrefix(version, "10.14"):
		return "macOS Mojave"
	case strings.HasPrefix(version, "10.13"):
		return "macOS High Sierra"
	case strings.HasPrefix(version, "10.12"):
		return "macOS Sierra"
	case strings.HasPrefix(version, "10.11"):
		return "OS X El Capitan"
	case strings.HasPrefix(version, "10.10"):
		return "OS X Yosemite"
	case strings.HasPrefix(version, "10.9"):
		return "OS X Mavericks"
	case strings.HasPrefix(version, "10.8"):
		return "OS X Mountain Lion"
	case strings.HasPrefix(version, "10.7"):
		return "Mac OS X Lion"
	case strings.HasPrefix(version, "10.6"):
		return "Mac OS X Snow Leopard"
	case strings.HasPrefix(version, "10.5"):
		return "Mac OS X Leopard"
	case strings.HasPrefix(version, "10.4"):
		return "Mac OS X Tiger"
	case strings.HasPrefix(version, "10.3"):
		return "Mac OS X Panther"
	case strings.HasPrefix(version, "10.2"):
		return "Mac OS X Jaguar"
	case strings.HasPrefix(version, "10.1"):
		return "Mac OS X Puma"
	case strings.HasPrefix(version, "10.0"):
		return "Mac OS X Cheetah"
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
