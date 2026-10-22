//go:build windows

package cli

import (
	"os"
	"strings"
)

func platformUsers() []string {
	name := os.Getenv("USERNAME")
	if name == "" {
		name = os.Getenv("USER")
	}
	if i := strings.LastIndex(name, `\`); i >= 0 {
		name = name[i+1:]
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	return []string{name}
}
