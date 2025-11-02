package cli

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

var (
	termRun  bool
	term     string
	termFont string
)

const (
	Linux     string = "Linux"
	MacOSX    string = "Mac OS X"
	MacOS     string = "macOS"
	Iphone    string = "iPhone OS"
	Windows   string = "Windows"
	Interix   string = "Interix"
	Haiku     string = "Haiku"
	FreeBSD   string = "FreeBSD"
	DragonFly string = "DragonFly"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readFirstLine(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

func Trim(s string) string {
	return strings.TrimSpace(s)
}
func RunCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return Trim(string(out))
}

func trimQuotes(s string) string {
	s = strings.Trim(s, `"`)
	return s
}

func readFileTrim(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func runCmd(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}
