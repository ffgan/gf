package modules

import (
	"strconv"
	"strings"
	"unicode"
)

type netCounter struct {
	name string
	rx   uint64
	tx   uint64
}

type diskRate struct {
	name     string
	read     float64
	write    float64
	combined bool
}

// ParseDNSList splits a Windows or BSD resolver list. Separators are
// commas, spaces, and semicolons.
func ParseDNSList(text string) []string {
	text = strings.ReplaceAll(text, ",", " ")
	text = strings.ReplaceAll(text, ";", " ")
	var out []string
	for _, field := range strings.Fields(text) {
		field = strings.Trim(field, "[]")
		if i := strings.IndexByte(field, '%'); i >= 0 {
			field = field[:i]
		}
		if field == "" || strings.EqualFold(field, "none") {
			continue
		}
		out = append(out, field)
	}
	return out
}

// ParseBraceLoadavg reads the "{ 0.15 0.10 0.08 }" form printed by sysctl.
func ParseBraceLoadavg(text string) (float64, float64, float64, bool) {
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "{}")
	return ParseLoadavg(text)
}

// ParseCPTime reads five cumulative tick counters: user, nice, system,
// interrupt, idle.
func ParseCPTime(text string) (cpuSample, bool) {
	fields := strings.Fields(text)
	if len(fields) < 5 {
		return cpuSample{}, false
	}
	var n [5]uint64
	for i := 0; i < 5; i++ {
		v, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			return cpuSample{}, false
		}
		n[i] = v
	}
	inUse := n[0] + n[1] + n[2] + n[3]
	return cpuSample{inUse: inUse, total: inUse + n[4]}, true
}

// ParseNetstatIB reads the Ibytes and Obytes columns of `netstat -ib`.
func ParseNetstatIB(text string) []netCounter {
	ibytes, obytes := -1, -1
	seen := map[string]bool{}
	var out []netCounter
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if ibytes < 0 {
			for i, field := range fields {
				switch strings.ToLower(field) {
				case "ibytes":
					ibytes = i
				case "obytes":
					obytes = i
				}
			}
			continue
		}
		if ibytes >= len(fields) || obytes >= len(fields) {
			continue
		}
		name := fields[0]
		if name == "lo" || name == "lo0" || seen[name] {
			continue
		}
		rx, err1 := strconv.ParseUint(fields[ibytes], 10, 64)
		tx, err2 := strconv.ParseUint(fields[obytes], 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		seen[name] = true
		out = append(out, netCounter{name: name, rx: rx, tx: tx})
	}
	return out
}

// ParseIostat reads a rate sample from FreeBSD `iostat -x` or the
// single-column macOS `iostat -d` table.
func ParseIostat(text string) []diskRate {
	if rates := parseIostatColumns(text); len(rates) > 0 {
		return rates
	}
	return parseIostatDarwin(text)
}

func parseIostatColumns(text string) []diskRate {
	lines := strings.Split(text, "\n")
	kr, kw := -1, -1
	header := -1
	for i, line := range lines {
		fields := strings.Fields(line)
		r, w := -1, -1
		for n, field := range fields {
			switch strings.ToLower(field) {
			case "kr/s", "kb/s", "kbs":
				if r < 0 {
					r = n
				} else if w < 0 {
					w = n
				}
			case "kw/s":
				w = n
			}
		}
		if r >= 0 && w >= 0 {
			kr, kw, header = r, w, i
		}
	}
	if header < 0 {
		return nil
	}
	hfields := strings.Fields(lines[header])
	if len(hfields) > 0 && !strings.EqualFold(hfields[0], "device") {
		kr++
		kw++
	}
	var out []diskRate
	for _, line := range lines[header+1:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			if len(out) > 0 {
				break
			}
			continue
		}
		if _, err := strconv.ParseFloat(fields[0], 64); err == nil {
			continue
		}
		if kr >= len(fields) || kw >= len(fields) {
			continue
		}
		readKB, err1 := strconv.ParseFloat(fields[kr], 64)
		writeKB, err2 := strconv.ParseFloat(fields[kw], 64)
		if err1 != nil || err2 != nil {
			continue
		}
		out = append(out, diskRate{
			name:  fields[0],
			read:  readKB * 1024,
			write: writeKB * 1024,
		})
	}
	return out
}

func parseIostatDarwin(text string) []diskRate {
	lines := strings.Split(text, "\n")
	header := -1
	for i, line := range lines {
		fields := strings.Fields(strings.ToLower(line))
		if len(fields) >= 3 && fields[0] == "kb/t" && fields[len(fields)-1] == "mb/s" {
			header = i
		}
	}
	if header < 1 {
		return nil
	}
	names := strings.Fields(lines[header-1])
	if len(names) == 0 {
		return nil
	}
	var last []string
	for _, line := range lines[header+1:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			if last != nil {
				break
			}
			continue
		}
		if _, err := strconv.ParseFloat(fields[0], 64); err != nil {
			break
		}
		last = fields
	}
	if len(last) < len(names)*3 {
		return nil
	}
	var out []diskRate
	for i, name := range names {
		mb, err := strconv.ParseFloat(last[i*3+2], 64)
		if err != nil {
			continue
		}
		out = append(out, diskRate{name: name, read: mb * 1024 * 1024, combined: true})
	}
	return out
}

// ParseDiskutil reads `diskutil info -all` and keeps whole physical disks.
func ParseDiskutil(text string) []string {
	chunks := strings.Split(text, "**********")
	var out []string
	for _, chunk := range chunks {
		whole := false
		virtual := false
		name := ""
		var size uint64
		for _, line := range strings.Split(chunk, "\n") {
			key, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			switch key {
			case "Whole":
				whole = strings.EqualFold(value, "yes")
			case "Virtual":
				virtual = strings.EqualFold(value, "yes")
			case "Device / Media Name":
				name = value
			case "Disk Size":
				if n, ok := bytesInParens(value); ok {
					size = n
				}
			}
		}
		if !whole || virtual || name == "" {
			continue
		}
		if size == 0 {
			out = append(out, name)
			continue
		}
		out = append(out, name+" ("+FormatBytes(size)+")")
	}
	return out
}

func bytesInParens(value string) (uint64, bool) {
	open := strings.LastIndex(value, "(")
	close := strings.LastIndex(value, "Bytes")
	if open < 0 || close < 0 || close < open {
		return 0, false
	}
	digits := strings.TrimSpace(value[open+1 : close])
	digits = strings.ReplaceAll(digits, ",", "")
	fields := strings.Fields(digits)
	if len(fields) == 0 {
		return 0, false
	}
	n, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}

// ParseGeomDisk reads `geom disk list`.
func ParseGeomDisk(text string) []string {
	var out []string
	var name, descr string
	var size uint64
	flush := func() {
		if name == "" && descr == "" {
			return
		}
		title := strings.TrimSpace(descr)
		if title == "" {
			title = name
		} else if name != "" {
			title += " " + name
		}
		if size > 0 {
			title += " (" + FormatBytes(size) + ")"
		}
		if title != "" {
			out = append(out, title)
		}
		name, descr, size = "", "", 0
	}
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "Geom name:") {
			flush()
			name = strings.TrimSpace(strings.TrimPrefix(trim, "Geom name:"))
			continue
		}
		if strings.HasPrefix(trim, "Mediasize:") {
			fields := strings.Fields(strings.TrimPrefix(trim, "Mediasize:"))
			if len(fields) > 0 {
				if n, err := strconv.ParseUint(fields[0], 10, 64); err == nil {
					size = n
				}
			}
			continue
		}
		if strings.HasPrefix(trim, "descr:") {
			descr = strings.TrimSpace(strings.TrimPrefix(trim, "descr:"))
		}
	}
	flush()
	return out
}

// ParseCamcontrol reads `camcontrol devlist` device names.
func ParseCamcontrol(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		start := strings.IndexByte(line, '<')
		end := strings.IndexByte(line, '>')
		if start < 0 || end <= start+1 {
			continue
		}
		name := strings.TrimSpace(line[start+1 : end])
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}

// ParseSndstat reads FreeBSD /dev/sndstat playback devices.
func ParseSndstat(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		open := strings.IndexByte(line, '<')
		close := strings.IndexByte(line, '>')
		if open < 0 || close <= open+1 {
			continue
		}
		name := strings.TrimSpace(line[open+1 : close])
		if name == "" {
			continue
		}
		if strings.Contains(strings.ToLower(line), "default") {
			name += " (default)"
		}
		out = append(out, name)
	}
	return out
}

// ParseIfconfigSSID reads the associated network from ifconfig output.
func ParseIfconfigSSID(text string) []string {
	var out []string
	iface := ""
	for _, line := range strings.Split(text, "\n") {
		if line != "" && !unicode.IsSpace(rune(line[0])) {
			name, _, _ := strings.Cut(line, ":")
			iface = strings.TrimSpace(name)
			continue
		}
		trim := strings.TrimSpace(line)
		for _, key := range []string{"ssid ", "nwid "} {
			value, ok := strings.CutPrefix(strings.ToLower(trim), key)
			if !ok {
				continue
			}
			ssid := strings.Trim(strings.TrimSpace(trim[len(key):]), `"`)
			_ = value
			if ssid == "" || strings.EqualFold(ssid, "none") {
				continue
			}
			if iface != "" {
				out = append(out, ssid+" ("+iface+")")
			} else {
				out = append(out, ssid)
			}
		}
	}
	return out
}

// ParseNetshWifi reads `netsh wlan show interfaces`.
func ParseNetshWifi(text string) []string {
	var out []string
	var ssid, signal, state string
	flush := func() {
		if !strings.EqualFold(state, "connected") || ssid == "" {
			ssid, signal, state = "", "", ""
			return
		}
		line := ssid
		if signal != "" {
			line += " (" + signal + ")"
		}
		out = append(out, line)
		ssid, signal, state = "", "", ""
	}
	for _, line := range strings.Split(text, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch strings.ToLower(key) {
		case "name":
			flush()
		case "state":
			state = value
		case "ssid":
			if ssid == "" && !strings.EqualFold(key, "bssid") {
				ssid = value
			}
		case "signal":
			signal = value
		}
	}
	flush()
	return out
}

// ParsePSTop reads `ps -o %cpu=,comm=` lines, busiest first.
func ParsePSTop(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if _, err := strconv.ParseFloat(fields[0], 64); err != nil {
			continue
		}
		name := strings.Join(fields[1:], " ")
		switch strings.ToLower(name) {
		case "kernel_task", "idle", "kernel", "system":
			continue
		}
		out = append(out, name)
		if len(out) == 5 {
			break
		}
	}
	return out
}

// ParsePSCount counts process names. Threads are summed when every line
// has a trailing integer.
func ParsePSCount(text string) (uint64, uint64, bool) {
	var n, threads uint64
	haveThreads := true
	saw := false
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		saw = true
		n++
		if len(fields) == 1 {
			haveThreads = false
			continue
		}
		v, err := strconv.ParseUint(fields[len(fields)-1], 10, 64)
		if err != nil {
			haveThreads = false
			continue
		}
		threads += v
	}
	if !saw {
		return 0, 0, false
	}
	return n, threads, haveThreads
}

// CleanPNPDesc keeps the human-readable half of a Windows DeviceDesc.
func CleanPNPDesc(value string) string {
	value = strings.TrimSpace(value)
	if i := strings.LastIndex(value, ";"); i >= 0 {
		value = strings.TrimSpace(value[i+1:])
	}
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "@") || strings.Contains(value, `\`) {
		return ""
	}
	return value
}

// LooksLikeGamepad reports controller-shaped device names.
func LooksLikeGamepad(name string) bool {
	n := strings.ToLower(name)
	for _, word := range []string{"gamepad", "controller", "joystick", "xbox", "dualshock", "dualsense"} {
		if strings.Contains(n, word) {
			return true
		}
	}
	return false
}

// AppleChassis maps a Mac model identifier onto an SMBIOS-style enclosure.
func AppleChassis(model string) string {
	m := strings.ToLower(strings.ReplaceAll(model, " ", ""))
	switch {
	case strings.Contains(m, "macbook"):
		return "Laptop"
	case strings.Contains(m, "macmini"):
		return "Mini PC"
	case strings.Contains(m, "imac"):
		return "All in One"
	case strings.Contains(m, "macpro"):
		return "Tower"
	case strings.Contains(m, "macstudio"):
		return "Desktop"
	default:
		return ""
	}
}

// ParseIORegPlatform reads IOPlatformExpertDevice properties.
func ParseIORegPlatform(text string) (vendor, model, board string) {
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(line)
		key, value, ok := strings.Cut(trim, "=")
		if !ok {
			continue
		}
		key = strings.Trim(strings.TrimSpace(key), `"`)
		value = decodeIOReg(value)
		switch key {
		case "manufacturer":
			vendor = value
		case "model":
			if model == "" {
				model = value
			}
		case "board-id":
			board = value
		}
	}
	return vendor, model, board
}

func decodeIOReg(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "<>")
	value = strings.Trim(value, `"`)
	value = strings.TrimSpace(value)
	if value == "" || strings.ContainsAny(value, " \t") {
		return value
	}
	if len(value)%2 != 0 || len(value) < 4 {
		return value
	}
	for _, r := range value {
		if !unicode.Is(unicode.ASCII_Hex_Digit, r) {
			return value
		}
	}
	var b strings.Builder
	for i := 0; i < len(value); i += 2 {
		n, err := strconv.ParseUint(value[i:i+2], 16, 8)
		if err != nil || n == 0 {
			break
		}
		if n < 32 || n > 126 {
			return strings.Trim(strings.TrimSpace(value), `"`)
		}
		b.WriteByte(byte(n))
	}
	if b.Len() == 0 {
		return value
	}
	return b.String()
}

type profileDevice struct {
	name  string
	props map[string]string
}

// ParseProfileDevices reads system_profiler stanzas whose headers end
// with a colon and whose properties are "Key: value" lines.
func ParseProfileDevices(text string) []profileDevice {
	var out []profileDevice
	var cur *profileDevice
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		key, value, ok := strings.Cut(trim, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if value == "" {
			out = append(out, profileDevice{name: key, props: map[string]string{}})
			cur = &out[len(out)-1]
			continue
		}
		if cur != nil {
			cur.props[key] = value
		}
	}
	return out
}

// FormatMemoryProfile turns SPMemoryDataType devices into fetch lines.
func FormatMemoryProfile(text string) []string {
	var out []string
	for _, dev := range ParseProfileDevices(text) {
		size := dev.props["Size"]
		if size == "" || strings.EqualFold(size, "empty") {
			continue
		}
		kind := dev.props["Type"]
		speed := dev.props["Speed"]
		if speed != "" {
			speed = strings.TrimSuffix(strings.TrimSpace(strings.TrimSuffix(speed, "MHz")), "MT/s")
			speed = strings.TrimSpace(speed)
			if digits, err := strconv.Atoi(strings.Fields(speed)[0]); err == nil {
				if kind != "" {
					kind += "-" + strconv.Itoa(digits)
				} else {
					kind = strconv.Itoa(digits) + " MT/s"
				}
			}
		}
		line := size
		if kind != "" {
			line += " " + kind
		}
		if vendor := dev.props["Manufacturer"]; vendor != "" && !strings.EqualFold(vendor, "unknown") {
			line += " (" + vendor + ")"
		}
		out = append(out, line)
	}
	return out
}

// FormatDisplayProfile turns SPDisplaysDataType devices into fetch lines.
func FormatDisplayProfile(text string) []string {
	var out []string
	for _, dev := range ParseProfileDevices(text) {
		res := dev.props["Resolution"]
		if res == "" {
			continue
		}
		res = strings.TrimSpace(strings.TrimSuffix(res, "Retina"))
		out = append(out, strings.TrimSpace(dev.name+" "+res))
	}
	return out
}

// FormatAudioProfile prefers the default output, then every named device.
func FormatAudioProfile(text string) []string {
	var all, defaults []string
	for _, dev := range ParseProfileDevices(text) {
		if dev.name == "Audio" || dev.name == "Devices" {
			continue
		}
		if _, ok := dev.props["Manufacturer"]; !ok && len(dev.props) == 0 {
			continue
		}
		if strings.EqualFold(dev.props["Default Output Device"], "yes") {
			defaults = append(defaults, dev.name)
			continue
		}
		if len(dev.props) > 0 {
			all = append(all, dev.name)
		}
	}
	if len(defaults) > 0 {
		return defaults
	}
	return all
}

// FormatBluetoothProfile reads the controller chipset and address.
func FormatBluetoothProfile(text string) string {
	chip, addr := "", ""
	for _, dev := range ParseProfileDevices(text) {
		if dev.props["Chipset"] != "" {
			chip = dev.props["Chipset"]
		}
		if dev.props["Address"] != "" && addr == "" {
			addr = dev.props["Address"]
		}
	}
	switch {
	case chip != "" && addr != "":
		return chip + " (" + addr + ")"
	case chip != "":
		return chip
	default:
		return addr
	}
}

// FormatCameraProfile lists camera device headers.
func FormatCameraProfile(text string) []string {
	var out []string
	for _, dev := range ParseProfileDevices(text) {
		if dev.name == "Camera" || dev.name == "Cameras" {
			continue
		}
		if len(dev.props) == 0 {
			continue
		}
		out = append(out, dev.name)
	}
	return out
}

// USBProductNames lists device headers from SPUSBDataType.
func USBProductNames(text string) []string {
	skip := map[string]bool{
		"USB": true, "USB 3.1 Bus": true, "USB 3.0 Bus": true, "USB 2.0 Bus": true,
	}
	var out []string
	for _, dev := range ParseProfileDevices(text) {
		if skip[dev.name] || strings.Contains(dev.name, "Bus") || strings.Contains(dev.name, "Hub") {
			continue
		}
		if len(dev.props) == 0 {
			continue
		}
		out = append(out, dev.name)
	}
	return out
}

// ParseAirportSSID reads `networksetup -getairportnetwork`.
func ParseAirportSSID(text string) string {
	text = strings.TrimSpace(text)
	if text == "" || strings.Contains(strings.ToLower(text), "not associated") {
		return ""
	}
	if _, value, ok := strings.Cut(text, ":"); ok {
		return strings.TrimSpace(value)
	}
	return text
}

// CacheFromSysctl builds the L1/L2/L3 summary sysctl sizes describe.
// l1 sizes are per core. A large L2 or L3 is treated as one shared cache.
func CacheFromSysctl(ncpu int, l1d, l1i, l2, l3 uint64) string {
	if ncpu < 1 {
		ncpu = 1
	}
	var items []cacheCount
	if l1d > 0 {
		items = addCache(items, 1, l1d, "D")
		items[len(items)-1].num = ncpu
	}
	if l1i > 0 {
		items = addCache(items, 1, l1i, "I")
		items[len(items)-1].num = ncpu
	}
	addLevel := func(level int, size uint64) {
		if size == 0 {
			return
		}
		items = addCache(items, level, size, "U")
		if size < 1024*1024 {
			items[len(items)-1].num = ncpu
		}
	}
	addLevel(2, l2)
	addLevel(3, l3)
	return summarizeCache(items)
}
