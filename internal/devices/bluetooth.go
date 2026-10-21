package dev

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetBluetooth() string {
	if names, ok := bluetoothFromSysfs(); ok {
		return strings.Join(names, ", ")
	}
	if !utils.CommandExists("lsusb") {
		return ""
	}

	out := utils.ExecOutput("lsusb")
	lines := strings.Split(out, "\n")
	var bluetooths []string
	for _, l := range lines {
		if strings.Contains(strings.ToLower(l), "bluetooth") {
			// 提取设备名称部分
			if len(l) >= 34 {
				name := strings.TrimSpace(l[33:])
				name = strings.TrimSuffix(name, "Bluetooth")
				name = strings.TrimSuffix(name, "bluetooth")
				bluetooths = append(bluetooths, strings.TrimSpace(name))
			}
		}
	}
	return strings.Join(bluetooths, ", ")
}

// bluetoothFromSysfs reads USB product strings the kernel already cached.
// lsusb walks usbfs and is the slow path when sysfs is missing.
func bluetoothFromSysfs() ([]string, bool) {
	if entries, err := os.ReadDir("/sys/class/bluetooth"); err == nil {
		var names []string
		for _, e := range entries {
			base := filepath.Join("/sys/class/bluetooth", e.Name())
			name := strings.TrimSpace(readSmall(filepath.Join(base, "device", "product")))
			if name == "" {
				name = strings.TrimSpace(readSmall(filepath.Join(base, "name")))
			}
			name = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(name, "Bluetooth"), "bluetooth"))
			if name != "" {
				names = append(names, name)
			}
		}
		return names, true
	}

	entries, err := os.ReadDir("/sys/bus/usb/devices")
	if err != nil {
		return nil, false
	}
	var names []string
	for _, e := range entries {
		prod := strings.TrimSpace(readSmall(filepath.Join("/sys/bus/usb/devices", e.Name(), "product")))
		if prod == "" || !strings.Contains(strings.ToLower(prod), "bluetooth") {
			continue
		}
		name := strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(prod, "Bluetooth"), "bluetooth"))
		if name != "" {
			names = append(names, name)
		}
	}
	return names, true
}

func readSmall(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}
