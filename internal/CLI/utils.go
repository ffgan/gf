package cli

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

var (
	termRun  bool
	term     string
	termFont string
)

const (
	Linux       string = "Linux"
	MacOSX      string = "Mac OS X"
	MacOS       string = "macOS"
	Iphone      string = "iPhone OS"
	Windows     string = "Windows"
	Interix     string = "Interix"
	Haiku       string = "Haiku"
	FreeBSD     string = "FreeBSD"
	DragonFly   string = "DragonFly"
	Darwin      string = "Darwin"
	SunOS       string = "SunOS"
	illumos     string = "illumos"
	Solaris     string = "Solaris"
	MINIX       string = "MINIX"
	AIX         string = "AIX"
	IRIX        string = "IRIX"
	FreeMiNT    string = "FreeMiNT"
	Ironclad    string = "Ironclad"
	OSF1        string = "OSF1"
	digitalUNIX string = "digitalUNIX"
	GNU         string = "GNU"
	Bitrig      string = "Bitrig"
	BSD         string = "BSD"
	CYGWIN      string = "CYGWIN"
	MSYS        string = "MSYS"
	MINGW       string = "MINGW"
	Windows_NT  string = "Windows_NT"
	RavynOS     string = "ravynOS"
	NetBSD      string = "NetBSD"
	OpenBSD     string = "OpenBSD"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readFirstLine(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

func Trim(s string) string {
	return strings.TrimSpace(s)
}
func RunCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return Trim(string(out))
}

func trimQuotes(s string) string {
	s = strings.Trim(s, `"`)
	return s
}

func readFileTrim(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func runCmd(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}

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
