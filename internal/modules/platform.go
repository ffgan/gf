package modules

import (
	"fmt"
	"time"
)

// hostPlatform is filled in by the Windows, macOS, and BSD files.
// Linux leaves it empty: the sysfs readers already cover those modules.
var hostPlatform platformHooks

type platformHooks struct {
	board          func() string
	chassis        func() string
	bootmgr        func() string
	brightness     func() string
	camera         func() string
	bluetoothRadio func() string
	cpuCache       func() string
	cpuUsage       func() string
	loadavg        func() string
	dns            func() string
	initSystem     func() string
	keyboard       func() string
	gamepad        func() string
	mouse          func() string
	processes      func() string
	top            func() string
	sound          func() string
	tpm            func() string
	physicalMemory func() string
	physicalDisk   func() string
	diskIO         func() string
	netIO          func() string
	monitor        func() string
	wifi           func() string
	wallpaper      func() string
	loginManager   func() string
	player         func() string
	libc           func() string
}

func callPlatform(fn func() string) string {
	if fn == nil {
		return ""
	}
	return fn()
}

func finish(value string, fallback func() string) string {
	if value != "" {
		return value
	}
	return callPlatform(fallback)
}

func formatCounterRates(first, second []netCounter, seconds float64) string {
	prev := map[string]netCounter{}
	for _, item := range first {
		prev[item.name] = item
	}
	var lines []string
	for _, item := range second {
		old, ok := prev[item.name]
		rx, tx := item.rx, item.tx
		if seconds > 0 {
			if !ok {
				continue
			}
			rx, tx = 0, 0
			if item.rx >= old.rx {
				rx = item.rx - old.rx
			}
			if item.tx >= old.tx {
				tx = item.tx - old.tx
			}
		}
		lines = append(lines, formatNetRate(item.name, rx, tx, seconds))
	}
	return joinLines(lines)
}

func formatDiskRates(rates []diskRate) string {
	var lines []string
	for _, rate := range rates {
		if rate.combined {
			lines = append(lines, fmt.Sprintf("%s: %s", rate.name, perSecond(uint64(rate.read), 1)))
			continue
		}
		lines = append(lines, formatDiskRate(rate.name, uint64(rate.read), uint64(rate.write), 1))
	}
	return joinLines(lines)
}

func usageFromSamples(read func() (cpuSample, bool)) string {
	first, ok := read()
	if !ok {
		return ""
	}
	if usagePause <= 0 {
		return ""
	}
	time.Sleep(usagePause)
	second, ok := read()
	if !ok {
		return ""
	}
	return usagePercent(first, second)
}
