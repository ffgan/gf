package hyfetch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FontLogo returns the nerd-font glyph hyfetch prints for a distro name.
// A previous match is cached at ~/.cache/hyfetch/font_logo, same as hyfetch.
func FontLogo(distro string) (string, error) {
	cache, err := cachePath("font_logo")
	if err == nil {
		if data, err := os.ReadFile(cache); err == nil && len(data) > 0 {
			return string(data), nil
		}
	}
	var logos map[string]string
	if err := json.Unmarshal(fontLogoJSON, &logos); err != nil {
		return "", err
	}
	lower := map[string]string{}
	for name, glyph := range logos {
		lower[strings.ToLower(name)] = glyph
	}
	distro = strings.ToLower(strings.TrimSpace(distro))
	matched := matchFont(lower, distro)
	if matched == "" {
		return "", fmt.Errorf("no font logo found for distro %q", distro)
	}
	glyph := lower[matched]
	if cache != "" {
		_ = os.MkdirAll(filepath.Dir(cache), 0o755)
		_ = os.WriteFile(cache, []byte(glyph), 0o644)
	}
	return glyph, nil
}

func matchFont(logos map[string]string, distro string) string {
	best := ""
	for name := range logos {
		if strings.Contains(distro, name) && len(name) > len(best) {
			best = name
		}
	}
	if best != "" {
		return best
	}
	for name := range logos {
		if distro != "" && strings.Contains(name, distro) && len(name) > len(best) {
			best = name
		}
	}
	if best != "" {
		return best
	}
	for name := range logos {
		for _, part := range strings.Fields(name) {
			if len(part) >= 3 && strings.Contains(distro, part) && len(name) > len(best) {
				best = name
			}
		}
	}
	return best
}

func cachePath(name string) (string, error) {
	dir := os.Getenv("XDG_CACHE_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil || home == "" {
			return "", fmt.Errorf("no home directory")
		}
		dir = filepath.Join(home, ".cache")
	}
	return filepath.Join(dir, "hyfetch", name), nil
}
