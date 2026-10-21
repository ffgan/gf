package dev

import (
	"os"
	"path/filepath"
	"strings"
)

// GetAccel reports NPUs and other accelerators. An empty string means none
// were visible; callers should not treat that as a failed detection.
func GetAccel() string {
	var out []string
	seen := map[string]struct{}{}
	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		if _, ok := seen[s]; ok {
			return
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}

	if entries, err := os.ReadDir("/sys/class/accel"); err == nil {
		for _, e := range entries {
			link, err := os.Readlink(filepath.Join("/sys/class/accel", e.Name(), "device", "driver"))
			if err != nil {
				continue
			}
			add(accelDriverName(filepath.Base(link)))
		}
	}
	if globExists("/dev/davinci*") {
		add("Huawei Ascend")
	}
	if globExists("/dev/cambricon*") || globExists("/dev/cambricon_dev*") {
		add("Cambricon")
	}
	if globExists("/dev/rknpu*") {
		add("Rockchip NPU")
	}
	if driverBound("intel_vpu") {
		add("Intel NPU")
	}
	if driverBound("amdxdna") {
		add("AMD XDNA")
	}
	return strings.Join(out, ", ")
}

func accelDriverName(driver string) string {
	switch driver {
	case "intel_vpu":
		return "Intel NPU"
	case "amdxdna":
		return "AMD XDNA"
	case "rknpu":
		return "Rockchip NPU"
	case "davinci", "devdrv":
		return "Huawei Ascend"
	case "cambricon", "cambricon_drv":
		return "Cambricon"
	case "":
		return ""
	default:
		return driver
	}
}

func globExists(pattern string) bool {
	matches, err := filepath.Glob(pattern)
	return err == nil && len(matches) > 0
}

func driverBound(name string) bool {
	entries, err := os.ReadDir("/sys/bus/pci/devices")
	if err != nil {
		return false
	}
	for _, e := range entries {
		link, err := os.Readlink(filepath.Join("/sys/bus/pci/devices", e.Name(), "driver"))
		if err == nil && filepath.Base(link) == name {
			return true
		}
	}
	if st, err := os.Stat(filepath.Join("/sys/module", name)); err == nil && st.IsDir() {
		return true
	}
	return false
}
