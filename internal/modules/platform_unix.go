//go:build darwin || freebsd || openbsd || netbsd || dragonfly

package modules

import (
	"strconv"
	"strings"
	"time"
)

func sysctlUint(key string) uint64 {
	n, _ := strconv.ParseUint(strings.TrimSpace(runOut("sysctl", "-n", key)), 10, 64)
	return n
}

func readCPTime() (cpuSample, bool) {
	return ParseCPTime(runOut("sysctl", "-n", "kern.cp_time"))
}

func unixInit(fallback string) string {
	name := strings.TrimSpace(runOut("ps", "-p", "1", "-o", "comm="))
	if name == "" {
		return fallback
	}
	if i := strings.LastIndexAny(name, `/\`); i >= 0 {
		name = name[i+1:]
	}
	return name
}

func psNames() []string {
	return strings.Fields(runOut("ps", "-axc", "-o", "comm="))
}

func unixProcesses() string {
	text := runOut("ps", "-ax", "-o", "comm=,nlwp=")
	n, threads, have := ParsePSCount(text)
	if n == 0 {
		text = runOut("ps", "-axc", "-o", "comm=")
		n, threads, have = ParsePSCount(text)
	}
	return formatProcesses(n, threads, have)
}

func unixTop() string {
	return joinLines(ParsePSTop(runOut("ps", "-axrc", "-o", "%cpu=,comm=")))
}

func unixNetIO() string {
	read := func() []netCounter { return ParseNetstatIB(runOut("netstat", "-ib")) }
	first := read()
	if len(first) == 0 {
		return ""
	}
	span := ioSample
	second := first
	if span > 0 {
		time.Sleep(span)
		second = read()
	}
	return formatCounterRates(first, second, span.Seconds())
}
