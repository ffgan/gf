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
		shell = trimExe(filepath.Base(shellEnv)) + " "
	}

	// If version output is disabled, return early
	if shellVersion != utils.ON {
		return shell
	}

	shellName := trimExe(filepath.Base(shellEnv))

	switch shellName {
	case "sh", "ash", "dash", "es":
		// these shells don't report a useful version
	default:
		shell += versionOfShell(shellEnv, shellName)
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

func trimExe(name string) string {
	if len(name) > 4 && strings.EqualFold(name[len(name)-4:], ".exe") {
		return name[:len(name)-4]
	}
	return name
}

func versionOfShell(shellPath, shellName string) string {
	if shellPath == "" {
		return ""
	}
	out := utils.RunCommand(shellPath, "--version")
	if out == "" {
		out = utils.ExecOutput(shellPath, "--version")
	}
	line := out
	if i := strings.IndexAny(line, "\r\n"); i >= 0 {
		line = line[:i]
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	re := regexp.MustCompile(`([0-9]+\.[0-9]+(?:\.[0-9]+)?)`)
	if m := re.FindString(line); m != "" {
		return m
	}
	line = strings.ReplaceAll(line, shellName, "")
	return strings.TrimSpace(line)
}
