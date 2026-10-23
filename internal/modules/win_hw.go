//go:build windows

package modules

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

func windowsCPUCache() string {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetLogicalProcessorInformation")
	if proc.Find() != nil {
		return ""
	}
	var size uint32
	_, _, err := proc.Call(0, uintptr(unsafe.Pointer(&size)))
	if size == 0 {
		_ = err
		return ""
	}
	buf := make([]byte, size)
	r, _, _ := proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r == 0 || int(size) > len(buf) {
		return ""
	}
	buf = buf[:size]
	ptr := int(unsafe.Sizeof(uintptr(0)))
	stride := 24
	cacheAt := 8
	if ptr == 8 {
		stride = 32
		cacheAt = 16
	}
	var items []cacheCount
	for off := 0; off+stride <= len(buf); off += stride {
		rel := uint32(buf[off+ptr]) | uint32(buf[off+ptr+1])<<8 | uint32(buf[off+ptr+2])<<16 | uint32(buf[off+ptr+3])<<24
		if rel != 2 {
			continue
		}
		base := off + cacheAt
		if base+12 > len(buf) {
			break
		}
		level := int(buf[base])
		sizeBytes := uint32(buf[base+4]) | uint32(buf[base+5])<<8 | uint32(buf[base+6])<<16 | uint32(buf[base+7])<<24
		kindN := uint32(buf[base+8]) | uint32(buf[base+9])<<8 | uint32(buf[base+10])<<16 | uint32(buf[base+11])<<24
		kind := "U"
		switch kindN {
		case 1:
			kind = "I"
		case 2:
			kind = "D"
		case 3:
			kind = "T"
		}
		if level < 1 || level > 4 || sizeBytes == 0 {
			continue
		}
		items = addCache(items, level, uint64(sizeBytes), kind)
	}
	return summarizeCache(items)
}

type winDisk struct {
	name   string
	size   uint64
	read   uint64
	write  uint64
	haveIO bool
}

func windowsDisks() string {
	disks := queryDisks(false)
	if len(disks) == 0 {
		disks = cimDisks()
	}
	var lines []string
	for _, disk := range disks {
		if disk.size == 0 {
			lines = append(lines, disk.name)
			continue
		}
		lines = append(lines, fmt.Sprintf("%s (%s)", disk.name, FormatBytes(disk.size)))
	}
	return joinLines(lines)
}

func windowsDiskIO() string {
	first := queryDisks(true)
	gotIO := false
	for _, disk := range first {
		if disk.haveIO {
			gotIO = true
			break
		}
	}
	if !gotIO {
		return cimDiskIO()
	}
	span := ioSample
	second := map[string]winDisk{}
	if span > 0 {
		time.Sleep(span)
		for _, disk := range queryDisks(true) {
			second[disk.name] = disk
		}
	}
	var lines []string
	for _, disk := range first {
		read, write := disk.read, disk.write
		if span > 0 {
			nxt, ok := second[disk.name]
			if !ok {
				continue
			}
			read, write = 0, 0
			if nxt.read >= disk.read {
				read = nxt.read - disk.read
			}
			if nxt.write >= disk.write {
				write = nxt.write - disk.write
			}
		}
		lines = append(lines, formatDiskRate(disk.name, read, write, span.Seconds()))
	}
	return joinLines(lines)
}

func queryDisks(counters bool) []winDisk {
	kernel := syscall.NewLazyDLL("kernel32.dll")
	open := kernel.NewProc("CreateFileW")
	ioctl := kernel.NewProc("DeviceIoControl")
	closer := kernel.NewProc("CloseHandle")
	if open.Find() != nil || ioctl.Find() != nil {
		return nil
	}
	var out []winDisk
	for i := 0; i < 32; i++ {
		path, _ := syscall.UTF16PtrFromString(fmt.Sprintf(`\\.\PhysicalDrive%d`, i))
		handle, _, _ := open.Call(uintptr(unsafe.Pointer(path)), 0x80, 3, 0, 3, 0, 0)
		if handle == 0 || handle == ^uintptr(0) {
			if i == 0 {
				continue
			}
			break
		}
		disk := winDisk{name: fmt.Sprintf("PhysicalDrive%d", i)}
		geom := make([]byte, 256)
		var got uint32
		if r, _, _ := ioctl.Call(handle, 0x700A0, 0, 0, uintptr(unsafe.Pointer(&geom[0])), uintptr(len(geom)), uintptr(unsafe.Pointer(&got)), 0); r != 0 && got >= 32 {
			disk.size = uint64(readInt64(geom[24:]))
		}
		desc := make([]byte, 1024)
		query := make([]byte, 12)
		if r, _, _ := ioctl.Call(handle, 0x2D1400, uintptr(unsafe.Pointer(&query[0])), uintptr(len(query)), uintptr(unsafe.Pointer(&desc[0])), uintptr(len(desc)), uintptr(unsafe.Pointer(&got)), 0); r != 0 && int(got) >= 32 && int(got) <= len(desc) {
			if name := storageName(desc[:got]); name != "" {
				disk.name = name
			}
		}
		if counters {
			perf := make([]byte, 88)
			if r, _, _ := ioctl.Call(handle, 0x70020, 0, 0, uintptr(unsafe.Pointer(&perf[0])), uintptr(len(perf)), uintptr(unsafe.Pointer(&got)), 0); r != 0 && got >= 16 {
				disk.read = uint64(readInt64(perf[0:]))
				disk.write = uint64(readInt64(perf[8:]))
				disk.haveIO = true
			}
		}
		closer.Call(handle)
		out = append(out, disk)
	}
	return out
}

func storageName(desc []byte) string {
	if len(desc) < 20 {
		return ""
	}
	vendorAt := int(uint32(desc[12]) | uint32(desc[13])<<8 | uint32(desc[14])<<16 | uint32(desc[15])<<24)
	productAt := int(uint32(desc[16]) | uint32(desc[17])<<8 | uint32(desc[18])<<16 | uint32(desc[19])<<24)
	vendor := cString(desc, vendorAt)
	product := cString(desc, productAt)
	return strings.TrimSpace(vendor + " " + product)
}

func cString(buf []byte, off int) string {
	if off <= 0 || off >= len(buf) {
		return ""
	}
	end := off
	for end < len(buf) && buf[end] != 0 {
		end++
	}
	return strings.TrimSpace(string(buf[off:end]))
}

func cimDisks() []winDisk {
	script := `Get-CimInstance Win32_DiskDrive | ForEach-Object { '{0}|{1}' -f $_.Model, $_.Size }`
	text := strings.ReplaceAll(runOut("powershell", "-NoProfile", "-NonInteractive", "-Command", script), "\x00", "")
	var out []winDisk
	for _, line := range strings.Split(text, "\n") {
		name, sizeText, ok := strings.Cut(strings.TrimSpace(line), "|")
		if !ok {
			continue
		}
		var size uint64
		fmt.Sscanf(strings.TrimSpace(sizeText), "%d", &size)
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		out = append(out, winDisk{name: name, size: size})
	}
	return out
}

func cimDiskIO() string {
	script := `Get-CimInstance Win32_PerfFormattedData_PerfDisk_PhysicalDisk | ForEach-Object { '{0}|{1}|{2}' -f $_.Name, $_.DiskReadBytesPersec, $_.DiskWriteBytesPersec }`
	text := strings.ReplaceAll(runOut("powershell", "-NoProfile", "-NonInteractive", "-Command", script), "\x00", "")
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		parts := strings.Split(strings.TrimSpace(line), "|")
		if len(parts) != 3 || strings.Contains(parts[0], "_Total") {
			continue
		}
		var read, write uint64
		fmt.Sscanf(parts[1], "%d", &read)
		fmt.Sscanf(parts[2], "%d", &write)
		lines = append(lines, formatDiskRate(strings.TrimSpace(parts[0]), read, write, 1))
	}
	return joinLines(lines)
}

type mibIfRow struct {
	Name            [256]uint16
	Index           uint32
	Type            uint32
	Mtu             uint32
	Speed           uint32
	PhysAddrLen     uint32
	PhysAddr        [8]byte
	AdminStatus     uint32
	OperStatus      uint32
	LastChange      uint32
	InOctets        uint32
	InUcastPkts     uint32
	InNUcastPkts    uint32
	InDiscards      uint32
	InErrors        uint32
	InUnknownProtos uint32
	OutOctets       uint32
	OutUcastPkts    uint32
	OutNUcastPkts   uint32
	OutDiscards     uint32
	OutErrors       uint32
	OutQLen         uint32
	DescrLen        uint32
	Descr           [256]byte
}

func windowsNetIO() string {
	read := readIfTable
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

func readIfTable() []netCounter {
	iphlp := syscall.NewLazyDLL("iphlpapi.dll")
	proc := iphlp.NewProc("GetIfTable")
	if proc.Find() != nil {
		return nil
	}
	var size uint32
	proc.Call(0, uintptr(unsafe.Pointer(&size)), 1)
	if size < 4 {
		return nil
	}
	buf := make([]byte, size+64)
	r, _, _ := proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)), 1)
	if r != 0 || int(size) > len(buf) || size < 4 {
		return nil
	}
	buf = buf[:size]
	count := int(uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24)
	stride := int(unsafe.Sizeof(mibIfRow{}))
	var out []netCounter
	for i := 0; i < count; i++ {
		off := 4 + i*stride
		if off+stride > len(buf) {
			break
		}
		row := (*mibIfRow)(unsafe.Pointer(&buf[off]))
		if row.Type == 24 || (row.OperStatus != 4 && row.OperStatus != 5) {
			continue
		}
		n := row.DescrLen
		if n > uint32(len(row.Descr)) {
			n = uint32(len(row.Descr))
		}
		name := strings.TrimRight(string(row.Descr[:n]), "\x00 ")
		if name == "" {
			name = syscall.UTF16ToString(row.Name[:])
		}
		if name == "" || strings.EqualFold(name, "loopback") {
			continue
		}
		out = append(out, netCounter{name: name, rx: uint64(row.InOctets), tx: uint64(row.OutOctets)})
	}
	return out
}

func windowsWifi() string {
	if lines := wifiWLAN(); lines != "" {
		return lines
	}
	return joinLines(ParseNetshWifi(runOut("netsh", "wlan", "show", "interfaces")))
}

func wifiWLAN() string {
	dll := syscall.NewLazyDLL("wlanapi.dll")
	open := dll.NewProc("WlanOpenHandle")
	enum := dll.NewProc("WlanEnumInterfaces")
	query := dll.NewProc("WlanQueryInterface")
	freeMem := dll.NewProc("WlanFreeMemory")
	closeH := dll.NewProc("WlanCloseHandle")
	if open.Find() != nil || enum.Find() != nil || query.Find() != nil {
		return ""
	}
	var negotiated uint32
	var client uintptr
	if r, _, _ := open.Call(2, 0, uintptr(unsafe.Pointer(&negotiated)), uintptr(unsafe.Pointer(&client))); r != 0 {
		return ""
	}
	defer closeH.Call(client, 0)
	var list uintptr
	if r, _, _ := enum.Call(client, 0, uintptr(unsafe.Pointer(&list))); r != 0 || list == 0 {
		return ""
	}
	defer freeMem.Call(list)
	count := *(*uint32)(unsafe.Pointer(list))
	const infoSize = 532
	var lines []string
	for i := uint32(0); i < count && i < 16; i++ {
		info := unsafe.Add(unsafe.Pointer(list), 8+int(i)*infoSize)
		state := *(*uint32)(unsafe.Add(info, 16+512))
		if state != 1 {
			continue
		}
		var dataSize uint32
		var data uintptr
		if r, _, _ := query.Call(client, uintptr(info), 7, 0, uintptr(unsafe.Pointer(&dataSize)), uintptr(unsafe.Pointer(&data)), 0); r != 0 || data == 0 {
			continue
		}
		ssid, quality := readWLAN(data, dataSize)
		freeMem.Call(data)
		if ssid == "" {
			continue
		}
		if quality > 0 {
			lines = append(lines, fmt.Sprintf("%s (%d%%)", ssid, quality))
		} else {
			lines = append(lines, ssid)
		}
	}
	return joinLines(lines)
}

func readWLAN(data uintptr, size uint32) (string, uint32) {
	if size < 580 {
		return "", 0
	}
	raw := unsafe.Slice((*byte)(unsafe.Pointer(data)), size)
	ssidLen := int(uint32(raw[520]) | uint32(raw[521])<<8 | uint32(raw[522])<<16 | uint32(raw[523])<<24)
	if ssidLen <= 0 || ssidLen > 32 || 524+ssidLen > len(raw) {
		return "", 0
	}
	quality := uint32(raw[576]) | uint32(raw[577])<<8 | uint32(raw[578])<<16 | uint32(raw[579])<<24
	return strings.TrimSpace(string(raw[524 : 524+ssidLen])), quality
}

func windowsTPM() string {
	tbs := syscall.NewLazyDLL("tbs.dll")
	proc := tbs.NewProc("Tbsi_GetDeviceInfo")
	if proc.Find() != nil {
		return ""
	}
	var info [4]uint32
	if r, _, _ := proc.Call(uintptr(unsafe.Sizeof(info)), uintptr(unsafe.Pointer(&info[0]))); r != 0 {
		return ""
	}
	version := ""
	switch info[1] {
	case 1:
		version = "1.2"
	case 2:
		version = "2.0"
	default:
		return ""
	}
	kind := ""
	switch info[2] {
	case 1:
		kind = "MMIO"
	case 2:
		kind = "TrustZone"
	case 3:
		kind = "HW"
	case 4:
		kind = "emulator"
	case 5:
		kind = "SPB"
	}
	if kind == "" {
		return "TPM " + version
	}
	return "TPM " + version + " (" + kind + ")"
}
