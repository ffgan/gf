package utils

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
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

func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func CommandExists(cmd string) bool {
	_, err := LookPath(cmd)
	return err == nil
}

// commandTimeout stops a probe that would otherwise sit on a slow
// Windows helper such as winget, wmic, or a store execution alias.
const commandTimeout = 8 * time.Second

func CommandOutput(name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).Output()
	return out, err
}

func CommandCombinedOutput(name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
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
