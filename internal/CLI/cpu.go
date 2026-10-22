package cli

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"os"
	"strings"
	"unsafe"

	"strconv"

	"github.com/ffgan/gf/internal/utils"
)

func GetCPU(OS, Machine, ShowCores, ShowSpeed, ShowTemp, ShowBrand, speedType, SpeedShort, TempUnit string) string {
	var showCores bool
	var showSpeed bool
	var showTemp bool
	var showBrand bool
	var speedShort bool

	coreType := "logical"
	switch ShowCores {
	case "physical":
		showCores = true
		coreType = "physical"
	case utils.OFF:
		showCores = false
	default:
		// "logical", "on", and the neofetch default all show logical cores.
		showCores = true
		coreType = "logical"
	}

	if ShowSpeed == utils.ON {
		showSpeed = true
	}

	tempUnit := "C"
	switch ShowTemp {
	case utils.ON, "C":
		showTemp = true
	case "F":
		showTemp = true
		tempUnit = "F"
	}
	if TempUnit == "F" {
		tempUnit = "F"
	}

	if ShowBrand == utils.ON {
		showBrand = true
	}

	if SpeedShort == utils.ON {
		speedShort = true
	}
	if speedType == "" {
		speedType = "bios_limit"
	}

	config := CPUConfig{
		ShowCores:  showCores,
		ShowSpeed:  showSpeed,
		ShowTemp:   showTemp,
		ShowBrand:  showBrand,
		CoreType:   coreType,
		SpeedType:  speedType,
		SpeedShort: speedShort,
		TempUnit:   tempUnit,
	}

	return getCPU(OS, Machine, config)
}

type CPUInfo struct {
	Name         string
	Cores        int
	Speed        int // MHz
	Temperature  float64
	OS           string
	Machine      string
	CountsInName bool
	// CoreLabel is a hybrid layout such as "8+16". It replaces the single count.
	CoreLabel string
}

type CPUConfig struct {
	ShowCores  bool
	ShowSpeed  bool
	ShowTemp   bool
	ShowBrand  bool
	CoreType   string // "logical" or "physical"
	SpeedType  string
	SpeedShort bool
	TempUnit   string // "C" or "F"
}

func getCPU(OS, Machine string, config CPUConfig) string {
	info, err := GetCPUInfo(OS, Machine, config)
	if err != nil || info == nil || info.Name == "" {
		return ""
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
	case Windows:
		return getCPUWindows(info)
	case Darwin:
		return getCPUDarwin(info, config)
	case FreeBSD, OpenBSD, NetBSD, DragonFly, BSD, Bitrig:
		return getCPUBSD(info, config)
	case Solaris, illumos:
		return getCPUSolaris(info, config)
	default:
		if isBSDFamily(OS) {
			return getCPUBSD(info, config)
		}
		if utils.FileExists("/proc/cpuinfo") {
			return getCPULinux(info, config)
		}
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
	dt, _ := os.ReadFile("/sys/firmware/devicetree/base/model")
	info.Name, info.CountsInName = IdentifyCPU(string(data), string(dt))

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
		if info.Cores == 0 {
			for _, line := range lines {
				if strings.HasPrefix(line, "processor") {
					info.Cores++
				}
			}
		}
	}

	// Get CPU speed. speed_type selects a cpufreq file; bios_limit is neofetch's default.
	speedName := config.SpeedType
	if speedName == "" {
		speedName = "bios_limit"
	}
	candidates := []string{
		"/sys/devices/system/cpu/cpu0/cpufreq/" + speedName,
		"/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq",
		"/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq",
	}
	var speedFile string
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			speedFile = candidate
			break
		}
	}
	if data, err := os.ReadFile(speedFile); err == nil && speedFile != "" {
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
	info.Name = applyFamilyCode(utils.RunCommand("sysctl", "-n", "machdep.cpu.brand_string"))
	if info.Name == "" || genericCPUName(info.Name) {
		family := strings.TrimSpace(utils.RunCommand("sysctl", "-n", "hw.cpufamily"))
		if named := applyFamilyCode(family); named != "" && named != family {
			info.Name = named
		} else {
			info.Name = ""
		}
	}

	// Fallback to system_profiler
	if info.Name == "" {
		out, err := utils.CommandOutput("system_profiler", "SPHardwareDataType")
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

	out, err := utils.CommandOutput("sysctl", "-n", sysctlKey)
	if err == nil {
		if cores, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.Cores = cores
		}
	}

	return info, nil
}

func getCPUBSD(info *CPUInfo, config CPUConfig) (*CPUInfo, error) {
	// Get CPU model
	info.Name = utils.RunCommand("sysctl", "-n", "hw.model")

	// Get CPU speed
	out, err := utils.CommandOutput("sysctl", "-n", "hw.cpuspeed")
	if err == nil {
		if speed, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.Speed = speed
		}
	}

	// Get CPU cores
	out, err = utils.CommandOutput("sysctl", "-n", "hw.ncpu")
	if err == nil {
		if cores, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.Cores = cores
		}
	}

	// Get temperature (FreeBSD/OpenBSD)
	out, _ = utils.CommandOutput("sysctl", "-n", "dev.cpu.0.temperature")
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
	out, err := utils.CommandOutput("psrinfo", "-pv")
	if err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) > 1 {
			info.Name = strings.TrimSpace(lines[1])
		}
	}

	// Get CPU speed
	out, err = utils.CommandOutput("sh", "-c", "psrinfo -v | awk '/operates at/ {print $6; exit}'")
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

	out, err = utils.CommandOutput("sh", "-c", cmd)
	if err == nil {
		if cores, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.Cores = cores
		}
	}

	return info, nil
}

func getCPUWindows(info *CPUInfo) (*CPUInfo, error) {
	const cpu0 = `HARDWARE\DESCRIPTION\System\CentralProcessor\0`
	const cpus = `HARDWARE\DESCRIPTION\System\CentralProcessor`
	name, ok := utils.RegString(utils.HKLM, cpu0, "ProcessorNameString")
	if !ok || name == "" {
		proc := utils.RunCommand("reg.exe", "query", `HKLM\`+cpu0)
		name = regValue(proc, "ProcessorNameString")
		if mhz, err := strconv.Atoi(regNumber(regValue(proc, "~MHz"))); err == nil {
			info.Speed = mhz
		}
		list := utils.RunCommand("reg.exe", "query", `HKLM\`+cpus)
		info.Cores = countRegSubkeys(list)
	} else {
		info.Name = name
		if mhz, ok := utils.RegUint(utils.HKLM, cpu0, "~MHz"); ok {
			info.Speed = int(mhz)
		}
		if subs, ok := utils.RegSubkeys(utils.HKLM, cpus); ok {
			info.Cores = len(subs)
		}
	}
	if info.Name == "" {
		info.Name = name
	}
	if info.Name == "" {
		return nil, fmt.Errorf("cpu name not in the registry")
	}
	info.CoreLabel = windowsCoreLabel()
	return info, nil
}

func FormatCPUInfo(info *CPUInfo, config CPUConfig) string {
	cpu := info.Name

	// Windows keeps the vendor marks fastfetch prints. Other systems follow
	// the shorter neofetch cleanup.
	if info.OS == Windows {
		cpu = strings.Join(strings.Fields(cpu), " ")
	} else {
		cpu = cleanCPUName(cpu)
	}

	// Remove brand if disabled
	if !config.ShowBrand {
		cpu = strings.TrimPrefix(cpu, "AMD ")
		cpu = strings.TrimPrefix(cpu, "Intel ")
		cpu = strings.TrimPrefix(cpu, "Qualcomm ")
	}

	// Add cores. A big.LITTLE name already lists each cluster.
	if config.ShowCores && !info.CountsInName {
		switch {
		case info.CoreLabel != "":
			cpu = fmt.Sprintf("%s (%s)", cpu, info.CoreLabel)
		case info.Cores > 0:
			cpu = fmt.Sprintf("%s (%d)", cpu, info.Cores)
		}
	}

	// Add speed
	if config.ShowSpeed && info.Speed > 0 {
		cpu = fmt.Sprintf("%s @ %s", cpu, formatCPUSpeed(info, config.SpeedShort))
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

func cleanCPUName(cpu string) string {
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
	return strings.Join(strings.Fields(cpu), " ")
}

func formatCPUSpeed(info *CPUInfo, shorthand bool) string {
	if info.OS == Windows && info.Speed >= 1000 {
		return fmt.Sprintf("%.2f GHz", float64(info.Speed)/1000)
	}
	return formatSpeed(info.Speed, shorthand)
}

func formatSpeed(mhz int, shorthand bool) string {
	if mhz < 1000 {
		return fmt.Sprintf("%dMHz", mhz)
	}
	ghz := float64(mhz) / 1000
	if shorthand {
		return fmt.Sprintf("%.1fGHz", ghz)
	}
	return fmt.Sprintf("%.3fGHz", ghz)
}

// coreLabelFromBuffer reads GetLogicalProcessorInformationEx output and
// returns performance+efficiency counts, highest efficiency class first.
// A single class returns an empty string so the plain core count is used.
func coreLabelFromBuffer(buf []byte) string {
	maskSize := int(unsafe.Sizeof(uintptr(0)))
	stride := 16
	if maskSize == 4 {
		stride = 12
	}
	counts := map[byte]int{}
	for off := 0; off+32 <= len(buf); {
		rel := binary.LittleEndian.Uint32(buf[off:])
		size := int(binary.LittleEndian.Uint32(buf[off+4:]))
		if size < 32 || off+size > len(buf) {
			break
		}
		if rel == 0 {
			class := buf[off+9]
			groups := int(binary.LittleEndian.Uint16(buf[off+30:]))
			n := 0
			for g := 0; g < groups; g++ {
				start := off + 32 + g*stride
				if start+maskSize > off+size || start+maskSize > len(buf) {
					break
				}
				var mask uint64
				if maskSize == 8 {
					mask = binary.LittleEndian.Uint64(buf[start:])
				} else {
					mask = uint64(binary.LittleEndian.Uint32(buf[start:]))
				}
				n += bits.OnesCount64(mask)
			}
			if n > 0 {
				counts[class] += n
			}
		}
		off += size
	}
	if len(counts) < 2 {
		return ""
	}
	classes := make([]byte, 0, len(counts))
	for class := range counts {
		classes = append(classes, class)
	}
	for i := 0; i < len(classes); i++ {
		for j := i + 1; j < len(classes); j++ {
			if classes[j] > classes[i] {
				classes[i], classes[j] = classes[j], classes[i]
			}
		}
	}
	parts := make([]string, len(classes))
	for i, class := range classes {
		parts[i] = strconv.Itoa(counts[class])
	}
	return strings.Join(parts, "+")
}
