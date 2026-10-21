package cli

import (
	"os"
	"regexp"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

// GetUsers lists unique names from `who`.
func GetUsers() string {
	out := utils.RunCommand("who")
	if out == "" {
		return ""
	}
	seen := map[string]bool{}
	var names []string
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || seen[fields[0]] {
			continue
		}
		seen[fields[0]] = true
		names = append(names, fields[0])
	}
	return strings.Join(names, ", ")
}

// GetLocale returns the active locale.
func GetLocale() string {
	for _, key := range []string{"LC_ALL", "LANG", "LC_MESSAGES"} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
	}
	for _, path := range []string{"/etc/locale.conf", "/etc/default/locale"} {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "LANG=") || strings.HasPrefix(line, "LC_ALL=") {
				return strings.Trim(strings.TrimPrefix(strings.SplitN(line, "=", 2)[1], ""), `"'`)
			}
		}
	}
	return ""
}

// ToolVersion runs bin --version (or the given args) and returns a short version.
func ToolVersion(bins []string, args ...string) string {
	if len(args) == 0 {
		args = []string{"--version"}
	}
	for _, bin := range bins {
		if !utils.CommandExists(bin) {
			continue
		}
		out := utils.ExecOutput(bin, args...)
		line := ""
		for _, candidate := range strings.Split(out, "\n") {
			candidate = strings.TrimSpace(candidate)
			if candidate != "" {
				line = candidate
				break
			}
		}
		if line == "" {
			continue
		}
		if m := regexp.MustCompile(`v?\d+\.\d+(?:\.\d+)?`).FindString(line); m != "" {
			return strings.TrimPrefix(m, "v")
		}
		return line
	}
	return ""
}
