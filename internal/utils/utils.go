package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	ON   = "on"
	OFF  = "off"
	AUTO = "auto"
)

func Unquote(s string) string {
	return strings.Trim(s, `"'`)
}

func RunFirst(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func RunCmdCount(cmd string, args ...string) int {
	c := exec.Command(cmd, args...)
	out, err := c.Output()
	if err != nil {
		return 0
	}
	lines := strings.Count(string(out), "\n")
	if strings.HasSuffix(string(out), "\n") {
		return lines
	}
	return lines + 1
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ReadFirstLine(path string) string {
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

func ReadFileTrim(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func RunCmd(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func GetPkgCount(cmd ...string) int {
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return 0
	}
	lines := bytes.Count(out, []byte{'\n'})
	return lines
}

func FilePathGlobLens(glob string) int {
	matches, _ := filepath.Glob(glob)
	return len(matches)
}

func GetParentPID(pid int) int {
	out, err := exec.Command("ps", "-o", "ppid=", "-p", fmt.Sprint(pid)).Output()
	if err != nil {
		return 0
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0
	}
	var p int
	fmt.Sscanf(s, "%d", &p)
	return p
}

func GetProcessName(pid string) string {
	out, err := exec.Command("ps", "-p", pid, "-o", "comm=").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func ExecOutput(name string, args ...string) string {
	out, _ := exec.Command(name, args...).CombinedOutput()
	return string(out)
}

func CommandOutput(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func ScanProcesses(psFlags []string, names []string) string {
	args := append([]string{}, psFlags...)
	out := RunCmd("ps", args...)
	for _, n := range names {
		re := regexp.MustCompile(fmt.Sprintf(`(?m)^%s$`, regexp.QuoteMeta(n)))
		if re.FindString(out) != "" {
			return n
		}
	}
	return ""
}
