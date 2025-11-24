package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func GetShell(shellPath, shellVersion string) string {
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		return ""
	}

	var shell string
	switch shellPath {
	case "on":
		shell = shellEnv + " "
	case "off":
		shell = filepath.Base(shellEnv) + " "
	}

	// If version output is disabled, return early
	if shellVersion != "on" {
		return shell
	}

	shellName := filepath.Base(shellEnv)

	getOutput := func(args ...string) string {
		cmd := exec.Command(shellEnv, args...)
		out, err := cmd.Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(out))
	}

	switch shellName {
	case "bash":
		version := os.Getenv("BASH_VERSION")
		if version == "" {
			version = getOutput("-c", "printf %s \"$BASH_VERSION\"")
		}
		shell += strings.Split(version, "-")[0]

	case "sh", "ash", "dash", "es":
		// no version info
	case "ksh", "mksh", "pdksh", "lksh":
		version := getOutput("-c", "printf %s \"$KSH_VERSION\"")
		version = strings.ReplaceAll(version, " * KSH", "")
		version = strings.ReplaceAll(version, "version", "")
		shell += strings.TrimSpace(version)

	case "osh":
		version := os.Getenv("OIL_VERSION")
		if version == "" {
			version = getOutput("-c", "printf %s \"$OIL_VERSION\"")
		}
		shell += version

	case "tcsh":
		version := getOutput("-c", "printf %s $tcsh")
		shell += version

	case "yash":
		version := getOutput("--version")
		version = strings.ReplaceAll(version, " "+shellName, "")
		version = strings.ReplaceAll(version, " Yet another shell", "")
		version = regexp.MustCompile(`Copyright.*`).ReplaceAllString(version, "")
		shell += strings.TrimSpace(version)

	case "nu":
		version := getOutput("-c", "version | get version")
		version = strings.ReplaceAll(version, " "+shellName, "")
		shell += version

	default:
		version := getOutput("--version")
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
