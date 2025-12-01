package cli

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetShell(shellPath, shellVersion string) string {
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		return ""
	}

	var shell string
	switch shellPath {
	case utils.ON:
		shell = shellEnv + " "
	case utils.OFF:
		shell = filepath.Base(shellEnv) + " "
	}

	// If version output is disabled, return early
	if shellVersion != utils.ON {
		return shell
	}

	shellName := filepath.Base(shellEnv)

	switch shellName {
	case "bash":
		version := os.Getenv("BASH_VERSION")
		if version == "" {
			version = utils.RunCommand("-c", "printf %s \"$BASH_VERSION\"")
		}
		shell += strings.Split(version, "-")[0]

	case "sh", "ash", "dash", "es":
		// no version info
	case "ksh", "mksh", "pdksh", "lksh":
		version := utils.RunCommand("-c", "printf %s \"$KSH_VERSION\"")
		version = strings.ReplaceAll(version, " * KSH", "")
		version = strings.ReplaceAll(version, "version", "")
		shell += strings.TrimSpace(version)

	case "osh":
		version := os.Getenv("OIL_VERSION")
		if version == "" {
			version = utils.RunCommand("-c", "printf %s \"$OIL_VERSION\"")
		}
		shell += version

	case "tcsh":
		version := utils.RunCommand("-c", "printf %s $tcsh")
		shell += version

	case "yash":
		version := utils.RunCommand("--version")
		version = strings.ReplaceAll(version, " "+shellName, "")
		version = strings.ReplaceAll(version, " Yet another shell", "")
		version = regexp.MustCompile(`Copyright.*`).ReplaceAllString(version, "")
		shell += strings.TrimSpace(version)

	case "nu":
		version := utils.RunCommand("-c", "version | get version")
		version = strings.ReplaceAll(version, " "+shellName, "")
		shell += version

	default:
		version := utils.RunCommand("--version")
		version = strings.ReplaceAll(version, " "+shellName, "")
		shell += version
	}

	// Clean unwanted info similar to Bash
	cleanups := []struct {
		pattern *regexp.Regexp
		repl    string
	}{
		{regexp.MustCompile(`, version`), ""},
		{regexp.MustCompile(`xonsh/`), "xonsh "},
		{regexp.MustCompile(`options.*`), ""},
		{regexp.MustCompile(`\([^)]*\)`), ""},
	}

	for _, c := range cleanups {
		shell = c.pattern.ReplaceAllString(shell, c.repl)
	}

	return strings.TrimSpace(shell)
}
