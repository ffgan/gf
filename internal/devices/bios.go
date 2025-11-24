package dev

import (
	"os"
	"path/filepath"
	"strings"

	cli "github.com/ffgan/gf/internal/CLI"
)

func GetBIOS(osName string) string {
	const dmiPath = "/sys/devices/virtual/dmi/id"

	if !strings.Contains(osName, cli.Linux) {
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
