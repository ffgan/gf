package cli

import (
	"fmt"
	"os"
	"strings"

	"os/exec"
	"strconv"
)

type CPUInfo struct {
	Name        string
	Cores       int
	Speed       int // MHz
	Temperature float64
	OS          string
	Machine     string
}

type CPUConfig struct {
	ShowCores  bool
	ShowSpeed  bool
	ShowTemp   bool
	ShowBrand  bool
	CoreType   string // "logical" or "physical"
	SpeedShort bool
	TempUnit   string // "C" or "F"
}

func getCPU(OS, Machine string, config CPUConfig) string {
	info, err := GetCPUInfo(OS, Machine, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return FormatCPUInfo(info, config)
}

func GetCPUInfo(OS, Machine string, config CPUConfig) (*CPUInfo, error) {
	info := &CPUInfo{
		OS:      OS,
		Machine: Machine,
	}

	switch OS {
	case Linux:
		return getCPULinux(info, config)
	case Darwin:
		return getCPUDarwin(info, config)
	case FreeBSD, OpenBSD, NetBSD, DragonFly:
		return getCPUBSD(info, config)
	case Solaris:
		return getCPUSolaris(info, config)
	default:
		return nil, fmt.Errorf("unsupported OS: %s", OS)
	}
}

func getCPULinux(info *CPUInfo, config CPUConfig) (*CPUInfo, error) {
	// Read /proc/cpuinfo
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	// Parse CPU name
	for _, line := range lines {
		if strings.Contains(line, "model name") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				info.Name = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	// Fallback for ARM and other architectures
	if info.Name == "" {
		for _, line := range lines {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "hardware") || strings.Contains(lower, "processor") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					info.Name = strings.TrimSpace(parts[1])
					break
				}
			}
		}
	}

	// Get CPU cores
	if config.CoreType == "logical" {
		processorCount := 0
		for _, line := range lines {
			if strings.HasPrefix(line, "processor") {
				processorCount++
			}
		}
		info.Cores = processorCount
	} else {
		// Physical cores
		coreIDs := make(map[string]bool)
		for _, line := range lines {
			if strings.HasPrefix(line, "core id") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					coreIDs[strings.TrimSpace(parts[1])] = true
				}
			}
		}
		info.Cores = len(coreIDs)
	}

	// Get CPU speed
	speedFile := "/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq"
	if data, err := os.ReadFile(speedFile); err == nil {
		if speed, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			info.Speed = speed / 1000 // Convert KHz to MHz
		}
	} else {
		// Fallback to cpuinfo
		for _, line := range lines {
			if strings.Contains(line, "cpu MHz") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					speedStr := strings.TrimSpace(parts[1])
					if speed, err := strconv.ParseFloat(speedStr, 64); err == nil {
						info.Speed = int(speed)
						break
					}
				}
			}
		}
	}

	// Get CPU temperature
	tempDirs := []string{
		"/sys/class/hwmon/hwmon0/temp1_input",
		"/sys/class/hwmon/hwmon1/temp1_input",
		"/sys/class/thermal/thermal_zone0/temp",
	}

	for _, tempFile := range tempDirs {
		if data, err := os.ReadFile(tempFile); err == nil {
			if temp, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
				info.Temperature = float64(temp) / 1000.0
				break
			}
		}
	}

	return info, nil
}

func getCPUDarwin(info *CPUInfo, config CPUConfig) (*CPUInfo, error) {
	// Get CPU brand string
	out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	if err == nil {
		info.Name = strings.TrimSpace(string(out))
	}

	// Fallback to system_profiler
	if info.Name == "" {
		out, err := exec.Command("system_profiler", "SPHardwareDataType").Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.Contains(line, "Processor Name") {
					parts := strings.Split(line, ":")
					if len(parts) >= 2 {
						info.Name = strings.TrimSpace(parts[1])
						break
					}
				}
			}
		}
	}

	// Get CPU cores
	var sysctlKey string
	if config.CoreType == "logical" {
		sysctlKey = "hw.logicalcpu_max"
	} else {
		sysctlKey = "hw.physicalcpu_max"
	}

	out, err = exec.Command("sysctl", "-n", sysctlKey).Output()
	if err == nil {
		if cores, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.Cores = cores
		}
	}

	return info, nil
}

func getCPUBSD(info *CPUInfo, config CPUConfig) (*CPUInfo, error) {
	// Get CPU model
	out, err := exec.Command("sysctl", "-n", "hw.model").Output()
	if err == nil {
		info.Name = strings.TrimSpace(string(out))
	}

	// Get CPU speed
	out, err = exec.Command("sysctl", "-n", "hw.cpuspeed").Output()
	if err == nil {
		if speed, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.Speed = speed
		}
	}

	// Get CPU cores
	out, err = exec.Command("sysctl", "-n", "hw.ncpu").Output()
	if err == nil {
		if cores, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.Cores = cores
		}
	}

	// Get temperature (FreeBSD/OpenBSD)
	out, _ = exec.Command("sysctl", "-n", "dev.cpu.0.temperature").Output()
	tempStr := strings.TrimSpace(string(out))
	if tempStr != "" {
		tempStr = strings.TrimSuffix(tempStr, "C")
		if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
			info.Temperature = temp
		}
	}

	return info, nil
}

func getCPUSolaris(info *CPUInfo, config CPUConfig) (*CPUInfo, error) {
	// Get CPU info
	out, err := exec.Command("psrinfo", "-pv").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) > 1 {
			info.Name = strings.TrimSpace(lines[1])
		}
	}

	// Get CPU speed
	out, err = exec.Command("sh", "-c", "psrinfo -v | awk '/operates at/ {print $6; exit}'").Output()
	if err == nil {
		if speed, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.Speed = speed
		}
	}

	// Get CPU cores
	var cmd string
	if config.CoreType == "logical" {
		cmd = "kstat -m cpu_info | grep -c chip_id"
	} else {
		cmd = "psrinfo -p"
	}

	out, err = exec.Command("sh", "-c", cmd).Output()
	if err == nil {
		if cores, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.Cores = cores
		}
	}

	return info, nil
}

func FormatCPUInfo(info *CPUInfo, config CPUConfig) string {
	cpu := info.Name

	// Clean up CPU name
	replacements := map[string]string{
		"(TM)":              "",
		"(tm)":              "",
		"(R)":               "",
		"(r)":               "",
		"CPU":               "",
		"Processor":         "",
		"Dual-Core":         "",
		"Quad-Core":         "",
		"Six-Core":          "",
		"Eight-Core":        "",
		"with Radeon":       "",
		"Graphics":          "",
		"Technologies, Inc": "",
		"Core2":             "Core 2",
	}

	for old, new := range replacements {
		cpu = strings.ReplaceAll(cpu, old, new)
	}

	// Remove extra spaces
	cpu = strings.Join(strings.Fields(cpu), " ")

	// Remove brand if disabled
	if !config.ShowBrand {
		cpu = strings.TrimPrefix(cpu, "AMD ")
		cpu = strings.TrimPrefix(cpu, "Intel ")
		cpu = strings.TrimPrefix(cpu, "Qualcomm ")
	}

	// Add cores
	if config.ShowCores && info.Cores > 0 {
		cpu = fmt.Sprintf("%s (%d)", cpu, info.Cores)
	}

	// Add speed
	if config.ShowSpeed && info.Speed > 0 {
		if info.Speed < 1000 {
			cpu = fmt.Sprintf("%s @ %dMHz", cpu, info.Speed)
		} else {
			speed := float64(info.Speed)
			if config.SpeedShort {
				speed = float64(info.Speed / 100)
			}
			cpu = fmt.Sprintf("%s @ %.1fGHz", cpu, speed/10.0)
		}
	}

	// Add temperature
	if config.ShowTemp && info.Temperature > 0 {
		temp := info.Temperature
		if config.TempUnit == "F" {
			temp = temp*9.0/5.0 + 32.0
		}
		cpu = fmt.Sprintf("%s [%.1f°%s]", cpu, temp, config.TempUnit)
	}

	return cpu
}
