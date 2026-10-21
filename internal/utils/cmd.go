package utils

import (
	"bytes"
	"os/exec"
	"strings"
	"sync"
)

func RunFirst(cmd string, args ...string) string {
	out, err := CommandOutput(cmd, args...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func RunCommand(name string, args ...string) string {
	out, err := CommandOutput(name, args...)
	if err != nil {
		return ""
	}
	return Trim(string(out))
}

func RunCommandCount(name string, args ...string) int {
	out, err := CommandOutput(name, args...)
	if err != nil {
		return 0
	}
	lines := strings.Count(string(out), "\n")
	if strings.HasSuffix(string(out), "\n") {
		return lines
	}
	return lines + 1
}

type lookResult struct {
	path string
	err  error
}

// lookCache remembers PATH lookups for this process. A fetch run asks
// about the same tools from several modules; exec.LookPath walks PATH
// every time otherwise.
var lookCache sync.Map

func LookPath(file string) (string, error) {
	if v, ok := lookCache.Load(file); ok {
		r := v.(lookResult)
		return r.path, r.err
	}
	p, err := exec.LookPath(file)
	lookCache.Store(file, lookResult{path: p, err: err})
	return p, err
}

func CommandExists(cmd string) bool {
	_, err := LookPath(cmd)
	return err == nil
}

func CommandOutput(name string, args ...string) ([]byte, error) {
	out, err := exec.Command(name, args...).Output()
	return out, err
}

func CommandCombinedOutput(name string, args ...string) ([]byte, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return out, err
}

func ExecOutput(name string, args ...string) string {
	out, _ := CommandCombinedOutput(name, args...)
	return string(out)
}

func GetPkgCount(cmd ...string) int {
	out, err := CommandOutput(cmd[0], cmd[1:]...)
	if err != nil {
		return 0
	}
	lines := bytes.Count(out, []byte{'\n'})
	return lines
}
