package dev

import (
	"os"
	"path/filepath"
	"strings"
)

// GetBattery reads Linux power_supply batteries. Other systems return empty.
func GetBattery() string {
	matches, _ := filepath.Glob("/sys/class/power_supply/BAT*")
	var parts []string
	for _, dir := range matches {
		capacity := readTrim(filepath.Join(dir, "capacity"))
		if capacity == "" {
			continue
		}
		status := readTrim(filepath.Join(dir, "status"))
		if status == "" {
			parts = append(parts, capacity+"%")
			continue
		}
		parts = append(parts, capacity+"% ["+status+"]")
	}
	return strings.Join(parts, ", ")
}

// GetPowerAdapter reports whether a mains adapter is online.
func GetPowerAdapter() string {
	matches, _ := filepath.Glob("/sys/class/power_supply/*/online")
	var parts []string
	for _, path := range matches {
		state := readTrim(path)
		name := filepath.Base(filepath.Dir(path))
		switch state {
		case "1":
			parts = append(parts, name+" (online)")
		case "0":
			parts = append(parts, name+" (offline)")
		}
	}
	return strings.Join(parts, ", ")
}

func readTrim(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
