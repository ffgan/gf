//go:build freebsd || openbsd || netbsd || dragonfly

package modules

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

func init() {
	hostPlatform = platformHooks{
		board:        bsdBoard,
		chassis:      bsdChassis,
		bootmgr:      bsdBootmgr,
		brightness:   bsdBrightness,
		cpuCache:     bsdCPUCache,
		cpuUsage:     func() string { return usageFromSamples(readCPTime) },
		loadavg:      bsdLoadavg,
		initSystem:   func() string { return unixInit("init") },
		processes:    unixProcesses,
		top:          unixTop,
		sound:        bsdSound,
		physicalDisk: bsdDisks,
		diskIO:       bsdDiskIO,
		netIO:        unixNetIO,
		wifi:         func() string { return joinLines(ParseIfconfigSSID(runOut("ifconfig"))) },
		loginManager: bsdLogin,
		player:       func() string { return playersFromNames(psNames()) },
	}
}

func bsdBoard() string {
	if runtime.GOOS == "freebsd" {
		if line := joinID(runOut("kenv", "smbios.planar.maker"), runOut("kenv", "smbios.planar.product"), runOut("kenv", "smbios.planar.version")); line != "" {
			return line
		}
	}
	if runtime.GOOS == "netbsd" {
		if line := joinID(
			runOut("sysctl", "-n", "machdep.dmi.board-vendor"),
			runOut("sysctl", "-n", "machdep.dmi.board-product"),
			runOut("sysctl", "-n", "machdep.dmi.board-version"),
		); line != "" {
			return line
		}
	}
	return joinID(runOut("sysctl", "-n", "hw.vendor"), runOut("sysctl", "-n", "hw.product"), runOut("sysctl", "-n", "hw.version"))
}

func bsdChassis() string {
	raw := ""
	switch runtime.GOOS {
	case "freebsd":
		raw = runOut("kenv", "smbios.chassis.type")
	case "netbsd":
		raw = runOut("sysctl", "-n", "machdep.dmi.chassis-type")
	}
	raw = strings.TrimSpace(raw)
	if n, err := strconv.Atoi(raw); err == nil {
		return ChassisType(n)
	}
	return raw
}

func bsdBootmgr() string {
	for _, name := range []string{"/boot/efi/EFI/freebsd/loader.efi", "/boot/loader", "/boot/efi/boot/bootx64.efi"} {
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}
	return ""
}

func bsdBrightness() string {
	cur := firstSysctl("hw.acpi.video.lcd0.brightness", "hw.acpi.video.brightness")
	full := firstSysctl("hw.acpi.video.lcd0.fullpower", "hw.acpi.video.fullpower")
	return BrightnessPercent(cur, full, "lcd0")
}

func firstSysctl(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(runOut("sysctl", "-n", key)); value != "" {
			return value
		}
	}
	return ""
}

func bsdCPUCache() string {
	n, _ := strconv.Atoi(firstSysctl("hw.ncpu"))
	return CacheFromSysctl(n,
		sysctlUint("hw.l1dcachesize"),
		sysctlUint("hw.l1icachesize"),
		sysctlUint("hw.l2cachesize"),
		sysctlUint("hw.l3cachesize"),
	)
}

func bsdLoadavg() string {
	a, b, c, ok := ParseBraceLoadavg(runOut("sysctl", "-n", "vm.loadavg"))
	if !ok {
		return ""
	}
	return formatLoadavg(a, b, c)
}

func bsdSound() string {
	for _, path := range []string{"/dev/sndstat", "/dev/sound"} {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if lines := ParseSndstat(string(data)); len(lines) > 0 {
			return joinLines(lines)
		}
	}
	return ""
}

func bsdDisks() string {
	if lines := ParseGeomDisk(runOut("geom", "disk", "list")); len(lines) > 0 {
		return joinLines(lines)
	}
	if lines := ParseCamcontrol(runOut("camcontrol", "devlist")); len(lines) > 0 {
		return joinLines(lines)
	}
	names := runOut("sysctl", "-n", "hw.disknames")
	var lines []string
	for _, part := range strings.Split(names, ",") {
		name, _, _ := strings.Cut(part, ":")
		name = strings.TrimSpace(name)
		if name != "" {
			lines = append(lines, name)
		}
	}
	return joinLines(lines)
}

func bsdDiskIO() string {
	text := runOut("iostat", "-x", "-c", "2", "-w", "1")
	if strings.TrimSpace(text) == "" {
		text = runOut("iostat", "-c", "2", "-w", "1")
	}
	return formatDiskRates(ParseIostat(text))
}

func bsdLogin() string {
	running := map[string]bool{}
	for _, name := range psNames() {
		running[strings.ToLower(name)] = true
	}
	for _, name := range []string{"sddm", "gdm", "gdm3", "lightdm", "greetd", "lxdm", "slim", "xdm"} {
		if running[name] {
			return name
		}
	}
	return ""
}
