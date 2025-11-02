package dev

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	cli "github.com/ffgan/gf/internal/CLI"
)

func getResolution() string {
	var resolution string
	osName := cli.GetOS()

	switch osName {
	case "macOS":
		resolution = getResolutionMac()
	case "iPhone OS":
		resolution = getResolutionIOS()
	case "Windows":
		resolution = getResolutionWindows()
	case "Haiku":
		resolution = getResolutionHaiku()
	default:
		resolution = getResolutionX11()
	}

	resolution = strings.TrimSuffix(resolution, ",")
	resolution = strings.TrimSpace(resolution)
	if !strings.Contains(resolution, "x") {
		return ""
	}
	return resolution
}

// --- macOS ---
func getResolutionMac() string {
	tmp := "/tmp/neofetch_system_profiler_SPDisplaysDataType.xml"
	cmd := exec.Command("system_profiler", "SPDisplaysDataType", "-xml")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	os.WriteFile(tmp, out, 0644)
	defer os.Remove(tmp)

	var resList []string
	for gpu := 0; gpu < 999; gpu++ {
		key := fmt.Sprintf("0:_items:%d", gpu)
		if !plistBuddyExists(tmp, key) {
			break
		}
		for display := 0; display < 999; display++ {
			base := fmt.Sprintf("%s:spdisplays_ndrvs:%d", key, display)
			rKey := base + ":_spdisplays_resolution"
			pKey := base + ":_spdisplays_pixels"

			rOut := runPlistBuddy(tmp, rKey)
			if rOut == "" {
				break
			}
			rOut = stripDecimal(rOut)

			pOut := runPlistBuddy(tmp, pKey)
			if pOut != "" {
				scale := calcScale(rOut, pOut)
				if scale > 1 {
					rOut += fmt.Sprintf(" @ %dx", scale)
				}
			}
			rOut = strings.ReplaceAll(rOut, " x ", "x")
			resList = append(resList, rOut)
		}
	}
	return strings.Join(resList, ", ")
}

func plistBuddyExists(plist, key string) bool {
	cmd := exec.Command("PlistBuddy", "-c", "print "+key, plist)
	return cmd.Run() == nil
}

func runPlistBuddy(plist, key string) string {
	cmd := exec.Command("PlistBuddy", "-c", "print "+key, plist)
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out))
}

func stripDecimal(s string) string {
	re := regexp.MustCompile(`\.[0-9]+`)
	return re.ReplaceAllString(s, "")
}

func calcScale(res, px string) int {
	rParts := strings.Fields(res)
	pParts := strings.Fields(px)
	if len(rParts) == 0 || len(pParts) == 0 {
		return 1
	}
	sx, _ := strconv.Atoi(rParts[0])
	pxv, _ := strconv.Atoi(pParts[0])
	if sx == 0 {
		return 1
	}
	return pxv / sx
}

// --- iOS ---
func getResolutionIOS() string {
	machine := cli.GetKernelMachine()
	// 这里只实现示例映射，可继续扩展
	table := map[string]string{
		"iPhone14,7": "1170x2532",
		"iPhone15,3": "1290x2796",
	}
	if val, ok := table[machine]; ok {
		return val
	}
	return ""
}

// --- Windows ---
func getResolutionWindows() string {
	cmdX := exec.Command("wmic", "path", "Win32_VideoController", "get", "CurrentHorizontalResolution,CurrentVerticalResolution")
	out, err := cmdX.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	var resolutions []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			resolutions = append(resolutions, fmt.Sprintf("%sx%s", fields[0], fields[1]))
		}
	}
	return strings.Join(resolutions, ", ")
}

// --- Haiku ---
func getResolutionHaiku() string {
	cmd := exec.Command("screenmode")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	re := regexp.MustCompile(`([0-9]+)x([0-9]+).*@ ([0-9.]+)`)
	match := re.FindStringSubmatch(string(out))
	if len(match) >= 3 {
		return fmt.Sprintf("%sx%s", match[1], match[2])
	}
	return ""
}

// --- Linux/X11 ---
func getResolutionX11() string {
	// 优先使用 xrandr
	if pathExists("xrandr") && os.Getenv("DISPLAY") != "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		cmd := exec.Command("xrandr", "--nograb", "--current")
		out, err := cmd.Output()
		if err != nil {
			return ""
		}
		re := regexp.MustCompile(`([0-9]+)x([0-9]+)\+`)
		matches := re.FindAllStringSubmatch(string(out), -1)
		var list []string
		for _, m := range matches {
			if len(m) > 2 {
				list = append(list, fmt.Sprintf("%sx%s", m[1], m[2]))
			}
		}
		return strings.Join(list, ", ")
	}

	// fallback: xdpyinfo
	if pathExists("xdpyinfo") && os.Getenv("DISPLAY") != "" {
		cmd := exec.Command("xdpyinfo")
		out, err := cmd.Output()
		if err == nil {
			re := regexp.MustCompile(`dimensions:\s+([0-9]+x[0-9]+)`)
			if m := re.FindStringSubmatch(string(out)); len(m) > 1 {
				return m[1]
			}
		}
	}

	// fallback: /sys/class/drm
	drm := "/sys/class/drm"
	if files, err := os.ReadDir(drm); err == nil {
		var list []string
		for _, f := range files {
			modeFile := filepath.Join(drm, f.Name(), "modes")
			if data, err := os.ReadFile(modeFile); err == nil {
				line := strings.TrimSpace(strings.SplitN(string(data), "\n", 2)[0])
				if line != "" {
					list = append(list, line)
				}
			}
		}
		return strings.Join(list, ", ")
	}

	return ""
}

func pathExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
