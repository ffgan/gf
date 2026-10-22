//go:build linux

package utils

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

func parentPIDFast(pid int) (int, bool) {
	b, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		if _, statErr := os.Stat("/proc"); statErr != nil {
			return 0, false
		}
		return 0, true
	}
	// pid (comm) state ppid ...
	// comm may contain spaces and parentheses, so split after the last ')'.
	s := string(b)
	i := strings.LastIndex(s, ")")
	if i < 0 || i+2 >= len(s) {
		return 0, true
	}
	fields := strings.Fields(s[i+2:])
	if len(fields) < 2 {
		return 0, true
	}
	ppid, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, true
	}
	return ppid, true
}

func processNameFast(pid string) (string, bool) {
	b, err := os.ReadFile("/proc/" + pid + "/comm")
	if err != nil {
		if _, statErr := os.Stat("/proc"); statErr != nil {
			return "", false
		}
		return "", true
	}
	return strings.TrimSpace(string(b)), true
}

var (
	procNamesOnce sync.Once
	procNames     []string
	procNamesOK   bool
)

func processNamesFast() ([]string, bool) {
	procNamesOnce.Do(func() {
		entries, err := os.ReadDir("/proc")
		if err != nil {
			return
		}
		procNamesOK = true
		for _, e := range entries {
			name := e.Name()
			if name == "" || name[0] < '0' || name[0] > '9' {
				continue
			}
			b, err := os.ReadFile("/proc/" + name + "/comm")
			if err != nil {
				continue
			}
			comm := strings.TrimSpace(string(b))
			if comm != "" {
				procNames = append(procNames, comm)
			}
		}
	})
	return procNames, procNamesOK
}
