package dev

import (
	"os"
	"path/filepath"
	"strings"
)

func getBios(osName string) string {
	const dmiPath = "/sys/devices/virtual/dmi/id"

	if !strings.Contains(osName, "Linux") {
		return ""
	}

	readFile := func(path string) string {
		data, err := os.ReadFile(path)
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(data))
	}

	biosVendor := readFile(filepath.Join(dmiPath, "bios_vendor"))
	if biosVendor == "" {
		return ""
	}

	bios := biosVendor

	if rel := readFile(filepath.Join(dmiPath, "bios_release")); rel != "" {
		bios += " " + rel
	}

	if date := readFile(filepath.Join(dmiPath, "bios_date")); date != "" {
		bios += " (" + date + ")"
	}

	return bios
}
