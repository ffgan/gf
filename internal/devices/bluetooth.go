package dev

import (
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetBluetooth() string {
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
