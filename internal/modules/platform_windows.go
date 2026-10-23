//go:build windows

package modules

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"
	"unsafe"

	"github.com/ffgan/gf/internal/utils"
)

func init() {
	hostPlatform = platformHooks{
		board:          func() string { return windowsSMBIOS().Board() },
		chassis:        func() string { return windowsSMBIOS().Chassis() },
		bootmgr:        windowsBootmgr,
		brightness:     windowsBrightness,
		camera:         func() string { return joinLines(pnpDevices(classCamera, nil)) },
		bluetoothRadio: func() string { return joinLines(pnpDevices(classBluetooth, nil)) },
		cpuCache:       windowsCPUCache,
		cpuUsage:       windowsCPUUsage,
		dns:            windowsDNS,
		initSystem:     windowsInit,
		keyboard:       func() string { return joinLines(pnpDevices(classKeyboard, nil)) },
		gamepad:        func() string { return joinLines(pnpDevices(classHID, LooksLikeGamepad)) },
		mouse:          func() string { return joinLines(pnpDevices(classMouse, nil)) },
		processes:      windowsProcesses,
		top:            windowsTop,
		sound:          windowsSound,
		tpm:            windowsTPM,
		physicalMemory: func() string { return joinLines(ParseMemoryDevices(rawSMBIOS())) },
		physicalDisk:   windowsDisks,
		diskIO:         windowsDiskIO,
		netIO:          windowsNetIO,
		monitor:        windowsMonitors,
		wifi:           windowsWifi,
		wallpaper:      windowsWallpaper,
		loginManager:   func() string { return "LogonUI" },
		player:         windowsPlayer,
		libc:           windowsLibc,
	}
}

const (
	classKeyboard  = "{4d36e96b-e325-11ce-bfc1-08002be10318}"
	classMouse     = "{4d36e96f-e325-11ce-bfc1-08002be10318}"
	classMonitor   = "{4d36e96e-e325-11ce-bfc1-08002be10318}"
	classHID       = "{745a17a0-74d3-11d0-b6fe-00a0c90f57da}"
	classBluetooth = "{e0cbf06c-cd8b-4647-bb8a-263b43f0f974}"
	classCamera    = "{ca3e7ab9-b4c3-4ae6-8251-579ef933890f}"
	classImage     = "{6bdd1fc6-810f-11d0-bec7-08002be2092f}"
)

func rawSMBIOS() []byte {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetSystemFirmwareTable")
	const rsmb = 0x52534D42
	n, _, _ := proc.Call(rsmb, 0, 0, 0)
	if n < 8 {
		return nil
	}
	buf := make([]byte, n)
	r, _, _ := proc.Call(rsmb, 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if r < 8 || int(r) > len(buf) {
		return nil
	}
	buf = buf[:r]
	return buf[8:]
}

func windowsSMBIOS() SMBIOSIDs {
	return ParseSMBIOSIDs(rawSMBIOS())
}

func firmwareVar(name string) []byte {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetFirmwareEnvironmentVariableW")
	n, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil
	}
	g, err := syscall.UTF16PtrFromString("{8be4df61-93ca-11d2-aa0d-00e098032b8c}")
	if err != nil {
		return nil
	}
	buf := make([]byte, 4096)
	r, _, _ := proc.Call(uintptr(unsafe.Pointer(n)), uintptr(unsafe.Pointer(g)), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if r == 0 || int(r) > len(buf) {
		return nil
	}
	return buf[:r]
}

func windowsBootmgr() string {
	current := firmwareVar("BootCurrent")
	if len(current) < 2 {
		return ""
	}
	order, ok := ParseBootOrder(append([]byte{0, 0, 0, 0}, current...))
	if !ok {
		return ""
	}
	name := ParseBootName(append([]byte{0, 0, 0, 0}, firmwareVar(fmt.Sprintf("Boot%04X", order))...))
	if name == "" {
		name = fmt.Sprintf("Boot%04X", order)
	}
	secure := firmwareVar("SecureBoot")
	if len(secure) > 0 && secure[0] == 1 {
		return name + " (Secure Boot)"
	}
	return name
}

func windowsBrightness() string {
	script := `(Get-CimInstance -Namespace root/WMI -ClassName WmiMonitorBrightness -ErrorAction SilentlyContinue | ForEach-Object { "$($_.CurrentBrightness)%" }) -join [char]10`
	text := strings.ReplaceAll(runOut("powershell", "-NoProfile", "-NonInteractive", "-Command", script), "\x00", "")
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasSuffix(line, "%") && line != "%" {
			lines = append(lines, line)
		}
	}
	return joinLines(lines)
}

func windowsDNS() string {
	var list []string
	for _, root := range []string{
		`SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces`,
		`SYSTEM\CurrentControlSet\Services\Tcpip6\Parameters\Interfaces`,
	} {
		ids, ok := utils.RegSubkeys(utils.HKLM, root)
		if !ok {
			continue
		}
		for _, id := range ids {
			path := root + `\` + id
			raw, _ := utils.RegString(utils.HKLM, path, "NameServer")
			if strings.TrimSpace(raw) == "" {
				raw, _ = utils.RegString(utils.HKLM, path, "DhcpNameServer")
			}
			list = append(list, ParseDNSList(raw)...)
		}
	}
	return joinLines(list)
}

func windowsWallpaper() string {
	path, ok := utils.RegString(utils.HKCU, `Control Panel\Desktop`, "WallPaper")
	if !ok {
		return ""
	}
	return strings.TrimSpace(path)
}

func windowsSound() string {
	const base = `SOFTWARE\Microsoft\Windows\CurrentVersion\MMDevices\Audio\Render`
	ids, ok := utils.RegSubkeys(utils.HKLM, base)
	if !ok {
		return ""
	}
	var lines []string
	for _, id := range ids {
		state, ok := utils.RegUint(utils.HKLM, base+`\`+id, "DeviceState")
		if !ok || state != 1 {
			continue
		}
		props := base + `\` + id + `\Properties`
		desc, _ := utils.RegString(utils.HKLM, props, `{b3f8fa53-0004-438e-9003-51a46e139bfc},6`)
		if desc == "" {
			desc, _ = utils.RegString(utils.HKLM, props, `{a45c254e-df1c-4efd-8020-67d146a850e0},2`)
		}
		if desc = CleanPNPDesc(desc); desc != "" {
			lines = append(lines, desc)
		}
	}
	return joinLines(lines)
}

func windowsLibc() string {
	root := os.Getenv("SystemRoot")
	if root == "" {
		root = `C:\Windows`
	}
	path := root + `\System32\ucrtbase.dll`
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	if ver := fileVersion(path); ver != "" {
		return "UCRT " + ver
	}
	return "UCRT"
}

func pnpDevices(classGUID string, accept func(string) bool) []string {
	classes := []string{classGUID}
	if strings.EqualFold(classGUID, classCamera) {
		classes = append(classes, classImage)
	}
	roots := []string{"HID", "USB", "ACPI", "PCI", "DISPLAY", "BTHENUM", "BTH", "SWD", "ROOT", "USBVIDEO"}
	const base = `SYSTEM\CurrentControlSet\Enum\`
	seen := map[string]bool{}
	var out []string
	for _, root := range roots {
		devs, ok := utils.RegSubkeys(utils.HKLM, base+root)
		if !ok {
			continue
		}
		for _, dev := range devs {
			insts, ok := utils.RegSubkeys(utils.HKLM, base+root+`\`+dev)
			if !ok {
				continue
			}
			for _, inst := range insts {
				path := base + root + `\` + dev + `\` + inst
				class, ok := utils.RegString(utils.HKLM, path, "ClassGUID")
				if !ok || !equalAny(class, classes) {
					continue
				}
				desc, _ := utils.RegString(utils.HKLM, path, "FriendlyName")
				if desc == "" {
					desc, _ = utils.RegString(utils.HKLM, path, "DeviceDesc")
				}
				desc = CleanPNPDesc(desc)
				if desc == "" || seen[strings.ToLower(desc)] {
					continue
				}
				if accept != nil && !accept(desc) {
					continue
				}
				seen[strings.ToLower(desc)] = true
				out = append(out, desc)
			}
		}
	}
	sort.Strings(out)
	return out
}

func equalAny(value string, options []string) bool {
	for _, option := range options {
		if strings.EqualFold(value, option) {
			return true
		}
	}
	return false
}

func windowsMonitors() string {
	const base = `SYSTEM\CurrentControlSet\Enum\`
	var lines []string
	seen := map[string]bool{}
	for _, root := range []string{"DISPLAY", "MONITOR"} {
		devs, ok := utils.RegSubkeys(utils.HKLM, base+root)
		if !ok {
			continue
		}
		for _, dev := range devs {
			insts, ok := utils.RegSubkeys(utils.HKLM, base+root+`\`+dev)
			if !ok {
				continue
			}
			for _, inst := range insts {
				path := base + root + `\` + dev + `\` + inst
				class, _ := utils.RegString(utils.HKLM, path, "ClassGUID")
				if class != "" && !strings.EqualFold(class, classMonitor) && root != "DISPLAY" {
					continue
				}
				name := ""
				if edid, ok := utils.RegBytes(utils.HKLM, path+`\Device Parameters`, "EDID"); ok {
					edidName, _, _ := ParseEDID(edid)
					name = edidName
				}
				if name == "" {
					desc, _ := utils.RegString(utils.HKLM, path, "FriendlyName")
					if desc == "" {
						desc, _ = utils.RegString(utils.HKLM, path, "DeviceDesc")
					}
					name = CleanPNPDesc(desc)
				}
				if name == "" || seen[strings.ToLower(name)] {
					continue
				}
				seen[strings.ToLower(name)] = true
				lines = append(lines, name)
			}
		}
	}
	sort.Strings(lines)
	return joinLines(lines)
}

func fileVersion(path string) string {
	version := syscall.NewLazyDLL("version.dll")
	sizeProc := version.NewProc("GetFileVersionInfoSizeW")
	infoProc := version.NewProc("GetFileVersionInfoW")
	queryProc := version.NewProc("VerQueryValueW")
	if sizeProc.Find() != nil || infoProc.Find() != nil || queryProc.Find() != nil {
		return ""
	}
	name, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return ""
	}
	size, _, _ := sizeProc.Call(uintptr(unsafe.Pointer(name)), 0)
	if size == 0 {
		return ""
	}
	buf := make([]byte, size)
	if r, _, _ := infoProc.Call(uintptr(unsafe.Pointer(name)), 0, uintptr(len(buf)), uintptr(unsafe.Pointer(&buf[0]))); r == 0 {
		return ""
	}
	sub, _ := syscall.UTF16PtrFromString(`\`)
	var block uintptr
	var blockLen uint32
	if r, _, _ := queryProc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(sub)), uintptr(unsafe.Pointer(&block)), uintptr(unsafe.Pointer(&blockLen))); r == 0 || block == 0 || blockLen < 16 {
		return ""
	}
	data := unsafe.Slice((*byte)(unsafe.Pointer(block)), 16)
	ms := uint32(data[8]) | uint32(data[9])<<8 | uint32(data[10])<<16 | uint32(data[11])<<24
	ls := uint32(data[12]) | uint32(data[13])<<8 | uint32(data[14])<<16 | uint32(data[15])<<24
	return fmt.Sprintf("%d.%d.%d.%d", ms>>16, ms&0xffff, ls>>16, ls&0xffff)
}
