package dev

import (
	"os"
	"path/filepath"
	"strings"
)

// GetGPUDriver reads the kernel driver bound to each DRM card.
func GetGPUDriver() string {
	matches, _ := filepath.Glob("/sys/class/drm/card[0-9]/device/driver")
	seen := map[string]bool{}
	var names []string
	for _, match := range matches {
		target, err := os.Readlink(match)
		if err != nil {
			continue
		}
		name := filepath.Base(target)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}
