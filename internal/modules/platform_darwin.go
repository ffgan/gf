//go:build darwin

package modules

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func init() {
	hostPlatform = platformHooks{
		board:          darwinBoard,
		chassis:        darwinChassis,
		bootmgr:        darwinBootmgr,
		camera:         func() string { return joinLines(FormatCameraProfile(profiler("SPCameraDataType"))) },
		bluetoothRadio: func() string { return FormatBluetoothProfile(profiler("SPBluetoothDataType")) },
		cpuCache:       darwinCPUCache,
		cpuUsage:       func() string { return usageFromSamples(readCPTime) },
		loadavg:        darwinLoadavg,
		initSystem:     func() string { return unixInit("launchd") },
		keyboard:       darwinKeyboard,
		gamepad:        darwinGamepad,
		mouse:          darwinMouse,
		processes:      unixProcesses,
		top:            unixTop,
		sound:          func() string { return joinLines(FormatAudioProfile(profiler("SPAudioDataType"))) },
		tpm:            darwinTPM,
		physicalMemory: darwinMemory,
		physicalDisk:   func() string { return joinLines(ParseDiskutil(runOut("diskutil", "info", "-all"))) },
		diskIO:         darwinDiskIO,
		netIO:          unixNetIO,
		monitor:        func() string { return joinLines(FormatDisplayProfile(profiler("SPDisplaysDataType"))) },
		wifi:           darwinWifi,
		wallpaper:      darwinWallpaper,
		loginManager:   func() string { return "loginwindow" },
		player:         func() string { return playersFromNames(psNames()) },
		libc:           darwinLibc,
	}
}

var profilerCache sync.Map

func profiler(dataType string) string {
	if value, ok := profilerCache.Load(dataType); ok {
		return value.(string)
	}
	text := runOut("system_profiler", dataType)
	profilerCache.Store(dataType, text)
	return text
}

func darwinBoard() string {
	vendor, model, board := ParseIORegPlatform(runOut("ioreg", "-rd1", "-c", "IOPlatformExpertDevice"))
	if model == "" {
		model = runOut("sysctl", "-n", "hw.model")
	}
	if vendor == "" && model != "" {
		vendor = "Apple"
	}
	return joinID(vendor, model, board)
}

func darwinChassis() string {
	_, model, _ := ParseIORegPlatform(runOut("ioreg", "-rd1", "-c", "IOPlatformExpertDevice"))
	if model == "" {
		model = runOut("sysctl", "-n", "hw.model")
	}
	return AppleChassis(model)
}

func darwinBootmgr() string {
	if runtime.GOARCH == "arm64" || runOut("sysctl", "-n", "hw.optional.arm64") == "1" {
		return "iBoot"
	}
	return "Apple EFI"
}

func darwinCPUCache() string {
	n, _ := strconv.Atoi(runOut("sysctl", "-n", "hw.physicalcpu"))
	if n == 0 {
		n, _ = strconv.Atoi(runOut("sysctl", "-n", "hw.ncpu"))
	}
	return CacheFromSysctl(n,
		sysctlUint("hw.l1dcachesize"),
		sysctlUint("hw.l1icachesize"),
		sysctlUint("hw.l2cachesize"),
		sysctlUint("hw.l3cachesize"),
	)
}

func darwinLoadavg() string {
	a, b, c, ok := ParseBraceLoadavg(runOut("sysctl", "-n", "vm.loadavg"))
	if !ok {
		return ""
	}
	return formatLoadavg(a, b, c)
}

func darwinMemory() string {
	if lines := FormatMemoryProfile(profiler("SPMemoryDataType")); len(lines) > 0 {
		return joinLines(lines)
	}
	n := sysctlUint("hw.memsize")
	if n == 0 {
		return ""
	}
	if runtime.GOARCH == "arm64" {
		return "Unified " + FormatBytes(n)
	}
	return FormatBytes(n)
}

func darwinDiskIO() string {
	text := runOut("iostat", "-d", "-c", "2", "-w", "1")
	return formatDiskRates(ParseIostat(text))
}

func darwinWifi() string {
	device := "en0"
	text := runOut("networksetup", "-listallhardwareports")
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.Contains(line, "Wi-Fi") && i+1 < len(lines) {
			if dev, ok := strings.CutPrefix(strings.TrimSpace(lines[i+1]), "Device:"); ok {
				if name := strings.TrimSpace(dev); name != "" {
					device = name
				}
			}
		}
	}
	if ssid := ParseAirportSSID(runOut("networksetup", "-getairportnetwork", device)); ssid != "" {
		return ssid + " (" + device + ")"
	}
	return ""
}

func darwinWallpaper() string {
	path := strings.TrimSpace(runOut("osascript", "-e", `tell application "System Events" to get picture of current desktop`))
	if strings.Contains(path, "/") {
		return path
	}
	return ""
}

func darwinTPM() string {
	if strings.Contains(profiler("SPiBridgeDataType"), "T2") {
		return "Apple T2"
	}
	return ""
}

func darwinLibc() string {
	if _, err := os.Stat("/usr/lib/libSystem.B.dylib"); err == nil {
		return "libSystem"
	}
	return ""
}

func darwinKeyboard() string {
	return joinLines(filterNames(darwinUSB(), "keyboard"))
}

func darwinMouse() string {
	names := filterNames(darwinUSB(), "mouse", "trackpad")
	if ioregContains("Apple Internal Keyboard / Trackpad") || ioregContains("AppleMultitouchDevice") {
		names = append([]string{"Apple Trackpad"}, names...)
	}
	return joinLines(names)
}

func darwinGamepad() string {
	var out []string
	for _, name := range darwinUSB() {
		if LooksLikeGamepad(name) {
			out = append(out, name)
		}
	}
	return joinLines(out)
}

func darwinUSB() []string {
	names := USBProductNames(profiler("SPUSBDataType"))
	if ioregContains("Apple Internal Keyboard / Trackpad") || ioregContains("AppleHIDKeyboard") {
		names = append([]string{"Apple Internal Keyboard"}, names...)
	}
	return names
}

func filterNames(names []string, words ...string) []string {
	var out []string
	for _, name := range names {
		lower := strings.ToLower(name)
		for _, word := range words {
			if strings.Contains(lower, word) {
				out = append(out, name)
				break
			}
		}
	}
	return out
}

var ioregText string
var ioregOnce sync.Once

func ioregContains(name string) bool {
	ioregOnce.Do(func() {
		ioregText = runOut("ioreg", "-r", "-c", "IOHIDDevice")
	})
	return strings.Contains(ioregText, name)
}
