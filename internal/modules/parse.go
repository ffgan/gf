package modules

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ParseDNS returns nameserver addresses from a resolv.conf document.
func ParseDNS(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 || !strings.EqualFold(fields[0], "nameserver") {
			continue
		}
		addr := fields[1]
		if i := strings.IndexByte(addr, '#'); i >= 0 {
			addr = addr[:i]
		}
		addr = strings.TrimSpace(addr)
		if addr != "" {
			out = append(out, addr)
		}
	}
	return out
}

// ParseLoadavg reads the three fields from /proc/loadavg.
func ParseLoadavg(text string) (float64, float64, float64, bool) {
	fields := strings.Fields(text)
	if len(fields) < 3 {
		return 0, 0, 0, false
	}
	a, err1 := strconv.ParseFloat(fields[0], 64)
	b, err2 := strconv.ParseFloat(fields[1], 64)
	c, err3 := strconv.ParseFloat(fields[2], 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, false
	}
	return a, b, c, true
}

func formatLoadavg(a, b, c float64) string {
	return fmt.Sprintf("%.2f, %.2f, %.2f", a, b, c)
}

type cpuSample struct {
	inUse uint64
	total uint64
}

// ParseCPUStat reads the aggregate "cpu " line and each "cpuN" line.
// The first return is the aggregate sample.
func ParseCPUStat(text string) (cpuSample, []cpuSample, bool) {
	var agg cpuSample
	var cores []cpuSample
	var saw bool
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 8 || !strings.HasPrefix(fields[0], "cpu") {
			continue
		}
		nums := make([]uint64, 0, 7)
		ok := true
		for _, f := range fields[1:8] {
			n, err := strconv.ParseUint(f, 10, 64)
			if err != nil {
				ok = false
				break
			}
			nums = append(nums, n)
		}
		if !ok || len(nums) < 7 {
			continue
		}
		// user nice system idle iowait irq softirq
		inUse := nums[0] + nums[1] + nums[2] + nums[5] + nums[6]
		sample := cpuSample{inUse: inUse, total: inUse + nums[3] + nums[4]}
		if fields[0] == "cpu" {
			agg = sample
			saw = true
			continue
		}
		cores = append(cores, sample)
	}
	if !saw && len(cores) > 0 {
		for _, c := range cores {
			agg.inUse += c.inUse
			agg.total += c.total
		}
		saw = true
	}
	return agg, cores, saw
}

func usagePercent(prev, next cpuSample) string {
	if next.total <= prev.total {
		return ""
	}
	used := next.inUse - prev.inUse
	span := next.total - prev.total
	pct := float64(used) * 100 / float64(span)
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	return fmt.Sprintf("%.0f%%", pct)
}

// ParseSoundCards reads /proc/asound/cards lines that name a card.
func ParseSoundCards(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !unicode.IsDigit(rune(line[0])) {
			continue
		}
		_, after, ok := strings.Cut(line, " - ")
		if !ok {
			continue
		}
		name := strings.TrimSpace(after)
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}

type wirelessLink struct {
	iface   string
	quality string
}

// ParseWireless reads /proc/net/wireless.
func ParseWireless(text string) []wirelessLink {
	var out []wirelessLink
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Inter") || strings.HasPrefix(line, "face") {
			continue
		}
		iface, rest, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		fields := strings.Fields(rest)
		if len(fields) < 2 {
			continue
		}
		q := strings.TrimRight(fields[1], ".")
		out = append(out, wirelessLink{iface: strings.TrimSpace(iface), quality: q})
	}
	return out
}

// ParseDiskStat reads a sysfs block stat file. Counts are bytes.
func ParseDiskStat(text string) (readBytes, writeBytes uint64, ok bool) {
	fields := strings.Fields(text)
	if len(fields) < 7 {
		return 0, 0, false
	}
	readSectors, err1 := strconv.ParseUint(fields[2], 10, 64)
	writeSectors, err2 := strconv.ParseUint(fields[6], 10, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return readSectors * 512, writeSectors * 512, true
}

// ParseKeyboards returns keyboard names from /proc/bus/input/devices.
// A device counts when it has a kbd handler, EV_REP, and KEY_A, which
// drops power buttons and media keys.
func ParseKeyboards(text string) []string {
	return parseInput(text, func(handlers string) bool {
		return hasToken(handlers, "kbd")
	})
}

// ParseGamepads returns joystick names from the same input document.
func ParseGamepads(text string) []string {
	return parseInput(text, func(handlers string) bool {
		for _, tok := range strings.Fields(handlers) {
			if strings.HasPrefix(tok, "js") && len(tok) > 2 && unicode.IsDigit(rune(tok[2])) {
				return true
			}
		}
		return false
	})
}

func parseInput(text string, want func(handlers string) bool) []string {
	var out []string
	var name, handlers string
	skip := false
	flush := func() {
		if !skip && name != "" && handlers != "" && want(handlers) {
			out = append(out, name)
		}
		name, handlers, skip = "", "", false
	}
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			flush()
			continue
		}
		switch {
		case strings.HasPrefix(line, "N: Name="):
			name = strings.Trim(strings.TrimPrefix(line, "N: Name="), `"`)
		case strings.HasPrefix(line, "H: Handlers="):
			handlers = strings.TrimPrefix(line, "H: Handlers=")
		case strings.HasPrefix(line, "B: EV="):
			val, err := strconv.ParseUint(strings.TrimSpace(strings.TrimPrefix(line, "B: EV=")), 16, 64)
			if err != nil || val&(1<<20) == 0 {
				skip = true
			}
		case strings.HasPrefix(line, "B: KEY="):
			bits := strings.TrimSpace(strings.TrimPrefix(line, "B: KEY="))
			last := bits
			if i := strings.LastIndexByte(bits, ' '); i >= 0 {
				last = bits[i+1:]
			}
			val, err := strconv.ParseUint(last, 16, 64)
			if err != nil || val&(1<<30) == 0 {
				skip = true
			}
		}
	}
	flush()
	return out
}

func hasToken(list, token string) bool {
	for _, tok := range strings.Fields(list) {
		if tok == token {
			return true
		}
	}
	return false
}

// procStat is one /proc/<pid>/stat snapshot.
type procStat struct {
	pid     int
	name    string
	flags   uint64
	cpu     uint64
	threads uint64
}

const pfKThread = 0x00200000

func (p procStat) kernel() bool {
	return p.flags&pfKThread != 0
}

// ParseProcStat parses a /proc/<pid>/stat file. The command name may
// contain spaces, so the fields after the last ')' are what matters.
func ParseProcStat(text string) (procStat, bool) {
	text = strings.TrimSpace(text)
	open := strings.IndexByte(text, '(')
	close := strings.LastIndexByte(text, ')')
	if open < 0 || close < open+1 || close+2 >= len(text) {
		return procStat{}, false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(text[:open]))
	if err != nil {
		return procStat{}, false
	}
	rest := strings.Fields(text[close+2:])
	// rest[0] is field 3 (state). flags is field 9, utime 14, stime 15, threads 20.
	if len(rest) < 18 {
		return procStat{}, false
	}
	flags, err1 := strconv.ParseUint(rest[6], 10, 64)
	utime, err2 := strconv.ParseUint(rest[11], 10, 64)
	stime, err3 := strconv.ParseUint(rest[12], 10, 64)
	threads, err4 := strconv.ParseUint(rest[17], 10, 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return procStat{}, false
	}
	return procStat{
		pid:     pid,
		name:    text[open+1 : close],
		flags:   flags,
		cpu:     utime + stime,
		threads: threads,
	}, true
}

// ParseEDID returns the monitor name and the physical size in millimetres.
func ParseEDID(data []byte) (name string, widthMM, heightMM int) {
	if len(data) < 128 || data[0] != 0x00 || data[1] != 0xFF {
		return "", 0, 0
	}
	widthMM = int(data[21]) * 10
	heightMM = int(data[22]) * 10
	for _, off := range []int{54, 72, 90, 108} {
		if off+18 > len(data) {
			break
		}
		if data[off] != 0 || data[off+1] != 0 || data[off+3] != 0xFC {
			continue
		}
		raw := string(data[off+5 : off+18])
		raw = strings.Trim(raw, " \t\n\r\x00")
		if raw != "" {
			name = raw
			break
		}
	}
	return name, widthMM, heightMM
}

// ParseBootOrder reads BootCurrent. The file starts with a 4-byte EFI
// attribute header; the uint16 after it is the boot entry index.
func ParseBootOrder(data []byte) (uint16, bool) {
	if len(data) < 6 {
		return 0, false
	}
	return uint16(data[4]) | uint16(data[5])<<8, true
}

// ParseBootName reads the UTF-16 description from a Boot#### variable.
func ParseBootName(data []byte) string {
	if len(data) < 4+4+2+2 {
		return ""
	}
	body := data[4:]
	// attributes uint32, file path list length uint16, then UTF-16 description.
	if len(body) < 8 {
		return ""
	}
	desc := body[6:]
	var chars []uint16
	for len(desc) >= 2 {
		u := uint16(desc[0]) | uint16(desc[1])<<8
		desc = desc[2:]
		if u == 0 {
			break
		}
		chars = append(chars, u)
	}
	return string(utf16Decode(chars))
}

func utf16Decode(chars []uint16) []rune {
	out := make([]rune, len(chars))
	for i, c := range chars {
		out[i] = rune(c)
	}
	return out
}

// ParseSecureBoot reports whether the SecureBoot variable is on.
func ParseSecureBoot(data []byte) (enabled bool, ok bool) {
	if len(data) < 5 {
		return false, false
	}
	return data[4] == 1, true
}

// ColorFGBG maps a COLORFGBG value such as "15;0" to a fg/bg pair.
func ColorFGBG(value string) string {
	fg, bg, ok := strings.Cut(value, ";")
	if !ok {
		return ""
	}
	fgN, err1 := strconv.Atoi(strings.TrimSpace(fg))
	bgN, err2 := strconv.Atoi(strings.TrimSpace(bg))
	if err1 != nil || err2 != nil {
		return ""
	}
	return fmt.Sprintf("fg %s / bg %s", ansiColor(fgN), ansiColor(bgN))
}

func ansiColor(n int) string {
	colors := [...]string{
		"#000000", "#cd0000", "#00cd00", "#cdcd00",
		"#0000ee", "#cd00cd", "#00cdcd", "#e5e5e5",
		"#7f7f7f", "#ff0000", "#00ff00", "#ffff00",
		"#5c5cff", "#ff00ff", "#00ffff", "#ffffff",
	}
	if n < 0 || n >= len(colors) {
		return fmt.Sprintf("%d", n)
	}
	return colors[n]
}

// CleanURI turns file:// and quoted gsettings values into a path.
func CleanURI(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"'`)
	value = strings.TrimPrefix(value, "file://")
	return strings.TrimSpace(value)
}

// ParseDMIDecodeMemory reads `dmidecode -t memory` text.
func ParseDMIDecodeMemory(text string) []string {
	var out []string
	var size, kind, speed, vendor, locator string
	flush := func() {
		if size == "" || strings.EqualFold(size, "no module installed") || strings.EqualFold(size, "not installed") {
			size, kind, speed, vendor, locator = "", "", "", "", ""
			return
		}
		parts := []string{size}
		mem := strings.TrimSpace(kind)
		if speed != "" && !strings.EqualFold(speed, "unknown") {
			mem = strings.TrimSpace(mem + " " + speed)
		}
		if mem != "" {
			parts = append(parts, mem)
		}
		if locator != "" {
			parts = append(parts, locator)
		}
		if vendor != "" && !strings.EqualFold(vendor, "unknown") {
			parts = append(parts, "("+vendor+")")
		}
		out = append(out, strings.Join(parts, " "))
		size, kind, speed, vendor, locator = "", "", "", "", ""
	}
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "Handle ") {
			flush()
			continue
		}
		key, val, ok := strings.Cut(trim, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		switch key {
		case "Size":
			size = val
		case "Type":
			kind = val
		case "Speed", "Configured Memory Speed":
			if speed == "" || key == "Configured Memory Speed" {
				speed = val
			}
		case "Manufacturer":
			vendor = val
		case "Locator":
			locator = val
		}
	}
	flush()
	return out
}

// systemdVersion takes the first line of `systemctl --version`.
func systemdVersion(text string) string {
	line := firstLine(text)
	rest, ok := strings.CutPrefix(line, "systemd ")
	if !ok {
		return ""
	}
	ver, _, _ := strings.Cut(rest, " ")
	return ver
}

// ParseLibc prefers getconf output, then the first line of ldd --version.
func ParseLibc(getconf, ldd string) string {
	getconf = firstLine(getconf)
	if strings.HasPrefix(strings.ToLower(getconf), "glibc") {
		return getconf
	}
	line := firstLine(ldd)
	if line == "" {
		return getconf
	}
	if strings.Contains(strings.ToLower(line), "musl") {
		return "musl"
	}
	if i := strings.LastIndex(line, " "); i >= 0 {
		tail := strings.TrimSpace(line[i+1:])
		if _, err := strconv.Atoi(strings.Trim(tail, ".")); err == nil || strings.Contains(tail, ".") {
			if strings.Contains(strings.ToLower(line), "glibc") || strings.Contains(strings.ToLower(line), "gnu") {
				return "glibc " + tail
			}
		}
	}
	return line
}

// BrightnessPercent formats current/max as a percentage plus the device name.
func BrightnessPercent(current, max, name string) string {
	cur, err1 := strconv.ParseFloat(strings.TrimSpace(current), 64)
	total, err2 := strconv.ParseFloat(strings.TrimSpace(max), 64)
	if err1 != nil || err2 != nil || total <= 0 {
		return ""
	}
	pct := cur * 100 / total
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	if name == "" {
		return fmt.Sprintf("%.0f%%", pct)
	}
	return fmt.Sprintf("%.0f%% [%s]", pct, name)
}
