package cli

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ffgan/gf/configs"
	"github.com/ffgan/gf/internal/utils"
)

func GetMemory(osName string, config *MemoryConfig) (*MemoryInfo, error) {
	var memUsed, memTotal int64
	var err error

	switch osName {
	case Linux:
		memUsed, memTotal, err = getMemoryLinux()
	case Darwin:
		memUsed, memTotal, err = getMemoryDarwin()
	case FreeBSD, OpenBSD, NetBSD, DragonFly:
		memUsed, memTotal, err = getMemoryBSD(osName)
	case Solaris:
		memUsed, memTotal, err = getMemorySolaris()
	case Windows:
		memUsed, memTotal, err = getMemoryWindows()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", osName)
	}

	if err != nil {
		return nil, err
	}

	info := &MemoryInfo{
		UsedKiB:  memUsed,
		TotalKiB: memTotal,
	}

	if config.ShowPercent && memTotal > 0 {
		info.Percent = int(memUsed * 100 / memTotal)
	}

	info.Formatted = formatMemory(memUsed, memTotal, info.Percent, config)

	return info, nil
}

func getMemoryLinux() (int64, int64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var memTotal, memFree, buffers, cached, sReclaimable, shmem, memAvail int64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, _ := strconv.ParseInt(fields[1], 10, 64)

		switch key {
		case "MemTotal":
			memTotal = value
		case "MemFree":
			memFree = value
		case "Buffers":
			buffers = value
		case "Cached":
			cached = value
		case "SReclaimable":
			sReclaimable = value
		case "Shmem":
			shmem = value
		case "MemAvailable":
			memAvail = value
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	var memUsed int64
	if memAvail > 0 {
		// Use MemAvailable if available (Linux 3.14+)
		memUsed = memTotal - memAvail
	} else {
		// MemUsed = MemTotal + Shmem - MemFree - Buffers - Cached - SReclaimable
		memUsed = memTotal + shmem - memFree - buffers - cached - sReclaimable
	}

	return memUsed, memTotal, nil
}

func getMemoryDarwin() (int64, int64, error) {
	// Get page size
	pageSize, err := sysctlInt64("hw.pagesize")
	if err != nil {
		return 0, 0, err
	}

	// Get total memory
	memTotal, err := sysctlInt64("hw.memsize")
	if err != nil {
		return 0, 0, err
	}
	memTotal /= 1024 // Convert to KiB

	// Get pageable and purgeable counts
	pageable, err := sysctlInt64("vm.page_pageable_internal_count")
	if err != nil {
		return 0, 0, err
	}

	purgeable, err := sysctlInt64("vm.page_purgeable_count")
	if err != nil {
		return 0, 0, err
	}

	// Get vm_stat output for wired and compressed pages
	output, err := utils.CommandOutput("vm_stat")
	if err != nil {
		return 0, 0, err
	}

	var wired, compressed int64
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "wired") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				wired, _ = strconv.ParseInt(strings.TrimSuffix(fields[3], "."), 10, 64)
			}
		}
		if strings.Contains(line, "occupied") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				compressed, _ = strconv.ParseInt(strings.TrimSuffix(fields[4], "."), 10, 64)
			}
		}
	}

	pagesApp := pageable - purgeable
	memUsed := (pagesApp + wired + compressed) * pageSize / 1024

	return memUsed, memTotal, nil
}

func getMemoryBSD(osName string) (int64, int64, error) {
	// Get total memory
	var memTotal int64
	var err error
	if osName == NetBSD {
		memTotal, err = sysctlInt64("hw.physmem64")
	} else {
		memTotal, err = sysctlInt64("hw.physmem")
	}
	if err != nil {
		return 0, 0, err
	}
	memTotal /= 1024 // Convert to KiB

	var memFree int64

	switch osName {
	case FreeBSD, DragonFly:
		pageSize, err := sysctlInt64("hw.pagesize")
		if err != nil {
			return 0, 0, err
		}

		inactive, _ := sysctlInt64("vm.stats.vm.v_inactive_count")
		unused, _ := sysctlInt64("vm.stats.vm.v_free_count")
		cache, _ := sysctlInt64("vm.stats.vm.v_cache_count")

		memFree = (inactive + unused + cache) * pageSize / 1024

	case OpenBSD:
		// OpenBSD uses vmstat for memory info
		memFree = getVmstatValue(3) // Used is in column 3
		return memFree * 1024, memTotal, nil

	default:
		// Generic BSD
		memFree = getVmstatValue(5)
	}

	memUsed := memTotal - memFree
	return memUsed, memTotal, nil
}

func getMemorySolaris() (int64, int64, error) {
	pageSize, err := getPageSize()
	if err != nil {
		return 0, 0, err
	}

	// Get total pages
	fields := strings.Fields(utils.RunCommand("kstat", "-p", "unix:0:system_pages:pagestotal"))
	pagesTotal, _ := strconv.ParseInt(fields[len(fields)-1], 10, 64)

	// Get free pages
	fields = strings.Fields(utils.RunCommand("kstat", "-p", "unix:0:system_pages:pagesfree"))
	pagesFree, _ := strconv.ParseInt(fields[len(fields)-1], 10, 64)

	memTotal := pagesTotal * pageSize / 1024
	memFree := pagesFree * pageSize / 1024
	memUsed := memTotal - memFree

	return memUsed, memTotal, nil
}

func getMemoryWindows() (int64, int64, error) {
	// This is a simplified version for Windows
	// In practice, you'd use golang.org/x/sys/windows for proper Windows API calls
	return 0, 0, fmt.Errorf("Windows support requires additional implementation")
}

func sysctlInt64(name string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(utils.RunCommand("sysctl", "-n", name)), 10, 64)
	return value, err
}

func getPageSize() (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(utils.RunCommand("pagesize")), 10, 64)
	return value, err
}

func getVmstatValue(column int) int64 {
	lines := strings.Split(utils.RunCommand("vmstat"), "\n")
	if len(lines) < 3 {
		return 0
	}
	fields := strings.Fields(lines[len(lines)-2])
	if len(fields) > column {
		value, _ := strconv.ParseInt(fields[column], 10, 64)
		return value
	}
	return 0
}

func formatMemory(usedKiB, totalKiB int64, percent int, config *MemoryConfig) string {
	var label string
	var divider int64 = 1

	switch config.Unit {
	case "tib":
		label = "TiB"
		divider = 1024 * 1024 * 1024
	case "gib":
		label = "GiB"
		divider = 1024 * 1024
	case "kib":
		label = "KiB"
		divider = 1
	default: // mib
		label = "MiB"
		divider = 1024
	}

	var usedStr, totalStr string
	if config.Precision == 0 {
		usedStr = fmt.Sprintf("%d", usedKiB/divider)
		totalStr = fmt.Sprintf("%d", totalKiB/divider)
	} else {
		multiplier := int64(math.Pow(10, float64(config.Precision)))
		usedInt := usedKiB / divider
		usedDec := (usedKiB % divider) * multiplier / divider
		totalInt := totalKiB / divider
		totalDec := (totalKiB % divider) * multiplier / divider

		format := fmt.Sprintf("%%d.%%0%dd", config.Precision)
		usedStr = fmt.Sprintf(format, usedInt, usedDec)
		totalStr = fmt.Sprintf(format, totalInt, totalDec)
	}

	info := fmt.Sprintf("%s %s / %s %s", usedStr, label, totalStr, label)
	if config.ShowPercent {
		info += fmt.Sprintf(" (%d%%)", percent)
	}

	switch config.DisplayMode {
	case "bar":
		return generateBar(usedKiB, totalKiB, config)
	case "infobar":
		return info + " " + generateBar(usedKiB, totalKiB, config)
	case "barinfo":
		return generateBar(usedKiB, totalKiB, config) + " " + info
	default:
		return info
	}
}

func generateBar(used, total int64, config *MemoryConfig) string {
	if total == 0 {
		return ""
	}

	filled := int(float64(used) / float64(total) * float64(config.BarWidth))
	if filled > config.BarWidth {
		filled = config.BarWidth
	}

	bar := strings.Repeat(config.BarCharFilled, filled)
	bar += strings.Repeat(config.BarCharEmpty, config.BarWidth-filled)

	return "[" + bar + "]"
}

func PrintMem(osName string, mem configs.Memory, bar configs.ProgressBar) string {
	precision, err := strconv.Atoi(mem.MemPrecision)
	if err != nil {
		fmt.Printf("failed to parse MemPrecision, %v\n", err)
		return ""
	}

	BarLength, err := strconv.Atoi(bar.BarLength)
	if err != nil {
		fmt.Printf("failed to parse MemPrecision, %v\n", err)
		return ""
	}

	var ShowPercent bool
	switch mem.MemoryPercent {
	case utils.ON:
		ShowPercent = true
	case utils.OFF:
		ShowPercent = false
	}

	config := &MemoryConfig{
		Unit:          mem.MemoryUnit,
		Precision:     precision,
		ShowPercent:   ShowPercent,
		DisplayMode:   bar.MemoryDisplay,
		BarWidth:      BarLength,
		BarCharFilled: bar.BarCharTotal,
		BarCharEmpty:  bar.BarCharElapsed,
	}

	info, err := GetMemory(osName, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting memory info: %v\n", err)
		os.Exit(1)
	}
	return info.Formatted
}

type MemoryInfo struct {
	UsedKiB   int64
	TotalKiB  int64
	Percent   int
	Formatted string
}

type MemoryConfig struct {
	Unit          string // "tib", "gib", "mib", "kib"
	Precision     int
	ShowPercent   bool
	DisplayMode   string // "bar", "infobar", "barinfo", ""
	BarWidth      int
	BarCharFilled string
	BarCharEmpty  string
}
