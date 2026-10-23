package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// sysRoot is "/" on a real machine. Tests point it at a fixture tree.
var sysRoot = "/"

// usagePause is the gap between the two /proc samples used for CPU rates.
var usagePause = 100 * time.Millisecond

// ioSample is the gap fastfetch uses before reporting disk and network rates.
// Zero prints the counters accumulated since boot.
var ioSample = 250 * time.Millisecond

func sysPath(elem ...string) string {
	parts := make([]string, 0, len(elem)+1)
	parts = append(parts, sysRoot)
	parts = append(parts, elem...)
	return filepath.Join(parts...)
}

func readTrim(elem ...string) string {
	data, err := os.ReadFile(sysPath(elem...))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func readBytes(elem ...string) []byte {
	data, err := os.ReadFile(sysPath(elem...))
	if err != nil {
		return nil
	}
	return data
}

// Board is the baseboard vendor, name, and version from DMI.
func Board() string {
	name := firstNonEmpty(
		readTrim("sys", "devices", "virtual", "dmi", "id", "board_name"),
		readTrim("sys", "class", "dmi", "id", "board_name"),
		readTrim("sys", "firmware", "devicetree", "base", "model"),
	)
	vendor := firstNonEmpty(
		readTrim("sys", "devices", "virtual", "dmi", "id", "board_vendor"),
		readTrim("sys", "class", "dmi", "id", "board_vendor"),
	)
	version := firstNonEmpty(
		readTrim("sys", "devices", "virtual", "dmi", "id", "board_version"),
		readTrim("sys", "class", "dmi", "id", "board_version"),
	)
	return finish(joinID(vendor, name, version), hostPlatform.board)
}

// Chassis is the enclosure type, plus the vendor when DMI has one.
func Chassis() string {
	raw := firstNonEmpty(
		readTrim("sys", "devices", "virtual", "dmi", "id", "chassis_type"),
		readTrim("sys", "class", "dmi", "id", "chassis_type"),
	)
	kind := ""
	if n, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil {
		kind = ChassisType(n)
	}
	if kind == "" {
		kind = strings.Trim(raw, "\x00 ")
	}
	vendor := firstNonEmpty(
		readTrim("sys", "devices", "virtual", "dmi", "id", "chassis_vendor"),
		readTrim("sys", "class", "dmi", "id", "chassis_vendor"),
	)
	return finish(joinID(kind, vendor, ""), hostPlatform.chassis)
}

func joinID(a, b, c string) string {
	var parts []string
	for _, p := range []string{a, b, c} {
		p = strings.TrimSpace(strings.Trim(p, "\x00"))
		if p == "" || strings.EqualFold(p, "none") || strings.EqualFold(p, "to be filled by o.e.m.") {
			continue
		}
		parts = append(parts, p)
	}
	return strings.Join(parts, " ")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		v = strings.TrimSpace(strings.Trim(v, "\x00"))
		if v != "" {
			return v
		}
	}
	return ""
}

// Bootmgr is the current EFI boot entry name, with Secure Boot when known.
func Bootmgr() string {
	order, ok := ParseBootOrder(readBytes("sys", "firmware", "efi", "efivars", "BootCurrent-8be4df61-93ca-11d2-aa0d-00e098032b8c"))
	if !ok {
		return callPlatform(hostPlatform.bootmgr)
	}
	name := ParseBootName(readBytes("sys", "firmware", "efi", "efivars", fmt.Sprintf("Boot%04X-8be4df61-93ca-11d2-aa0d-00e098032b8c", order)))
	if name == "" {
		name = fmt.Sprintf("Boot%04X", order)
	}
	if on, ok := ParseSecureBoot(readBytes("sys", "firmware", "efi", "efivars", "SecureBoot-8be4df61-93ca-11d2-aa0d-00e098032b8c")); ok && on {
		return name + " (Secure Boot)"
	}
	return name
}

// Brightness reads every backlight class device.
func Brightness() string {
	dir := sysPath("sys", "class", "backlight")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return callPlatform(hostPlatform.brightness)
	}
	var lines []string
	for _, entry := range entries {
		cur := readTrim("sys", "class", "backlight", entry.Name(), "brightness")
		max := readTrim("sys", "class", "backlight", entry.Name(), "max_brightness")
		if line := BrightnessPercent(cur, max, entry.Name()); line != "" {
			lines = append(lines, line)
		}
	}
	return finish(joinLines(lines), hostPlatform.brightness)
}

// Btrfs lists mounted filesystems under /sys/fs/btrfs.
func Btrfs() string {
	dir := sysPath("sys", "fs", "btrfs")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var lines []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		label := readTrim("sys", "fs", "btrfs", entry.Name(), "label")
		devs := dirNames(sysPath("sys", "fs", "btrfs", entry.Name(), "devices"))
		title := label
		if title == "" {
			title = entry.Name()
		} else {
			title = fmt.Sprintf("%s (%s)", label, entry.Name())
		}
		if len(devs) > 0 {
			title += " on " + strings.Join(devs, ", ")
		}
		lines = append(lines, title)
	}
	return joinLines(lines)
}

func dirNames(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if entry.Name() == "." || entry.Name() == ".." {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}

// Camera lists V4L capture devices by their sysfs name.
func Camera() string {
	dir := sysPath("sys", "class", "video4linux")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return callPlatform(hostPlatform.camera)
	}
	var lines []string
	for _, entry := range entries {
		name := readTrim("sys", "class", "video4linux", entry.Name(), "name")
		if name == "" || strings.Contains(strings.ToLower(name), "metadata") {
			continue
		}
		lines = append(lines, name)
	}
	return finish(joinLines(lines), hostPlatform.camera)
}

// BluetoothRadio lists host controllers.
func BluetoothRadio() string {
	dir := sysPath("sys", "class", "bluetooth")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return callPlatform(hostPlatform.bluetoothRadio)
	}
	var lines []string
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "hci") {
			continue
		}
		addr := readTrim("sys", "class", "bluetooth", entry.Name(), "address")
		name := readTrim("sys", "class", "bluetooth", entry.Name(), "name")
		switch {
		case name != "" && addr != "":
			lines = append(lines, name+" ("+addr+")")
		case name != "":
			lines = append(lines, name)
		case addr != "":
			lines = append(lines, entry.Name()+" ("+addr+")")
		default:
			lines = append(lines, entry.Name())
		}
	}
	return finish(joinLines(lines), hostPlatform.bluetoothRadio)
}

// CPUCache summarises each cache level from sysfs.
func CPUCache() string {
	cpuDir := sysPath("sys", "devices", "system", "cpu")
	entries, err := os.ReadDir(cpuDir)
	if err != nil {
		return callPlatform(hostPlatform.cpuCache)
	}
	seen := map[string]bool{}
	var items []cacheCount
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "cpu") || len(name) < 4 || name[3] < '0' || name[3] > '9' {
			continue
		}
		cacheDir := filepath.Join(cpuDir, name, "cache")
		indexes, err := os.ReadDir(cacheDir)
		if err != nil {
			continue
		}
		for _, index := range indexes {
			if !strings.HasPrefix(index.Name(), "index") {
				continue
			}
			base := filepath.Join(cacheDir, index.Name())
			levelText := strings.TrimSpace(readFile(filepath.Join(base, "level")))
			sizeText := strings.TrimSpace(readFile(filepath.Join(base, "size")))
			kind := strings.TrimSpace(readFile(filepath.Join(base, "type")))
			shared := strings.TrimSpace(readFile(filepath.Join(base, "shared_cpu_list")))
			level, _ := strconv.Atoi(levelText)
			size := parseCacheSize(sizeText)
			if level < 1 || level > 4 || size == 0 || kind == "" {
				continue
			}
			key := fmt.Sprintf("%s|%d|%d|%s", shared, level, size, kind)
			if seen[key] {
				continue
			}
			seen[key] = true
			items = addCache(items, level, size, cacheKind(kind))
		}
	}
	return finish(summarizeCache(items), hostPlatform.cpuCache)
}

func cacheKind(raw string) string {
	switch strings.ToLower(raw) {
	case "instruction":
		return "I"
	case "data":
		return "D"
	case "unified":
		return "U"
	case "trace":
		return "T"
	default:
		return raw
	}
}

func parseCacheSize(text string) uint64 {
	text = strings.TrimSpace(strings.ToUpper(text))
	mult := uint64(1)
	switch {
	case strings.HasSuffix(text, "K"):
		mult = 1024
		text = strings.TrimSuffix(text, "K")
	case strings.HasSuffix(text, "M"):
		mult = 1024 * 1024
		text = strings.TrimSuffix(text, "M")
	case strings.HasSuffix(text, "G"):
		mult = 1024 * 1024 * 1024
		text = strings.TrimSuffix(text, "G")
	}
	n, err := strconv.ParseUint(strings.TrimSpace(text), 10, 64)
	if err != nil {
		return 0
	}
	return n * mult
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// CPUUsage samples /proc/stat twice and returns the aggregate busy percent.
func CPUUsage() string {
	first, _, ok := ParseCPUStat(string(readBytes("proc", "stat")))
	if !ok {
		return callPlatform(hostPlatform.cpuUsage)
	}
	if usagePause > 0 {
		time.Sleep(usagePause)
	}
	second, _, ok := ParseCPUStat(string(readBytes("proc", "stat")))
	if !ok {
		return ""
	}
	return usagePercent(first, second)
}

// Loadavg is the 1, 5, and 15 minute load.
func Loadavg() string {
	a, b, c, ok := ParseLoadavg(readTrim("proc", "loadavg"))
	if !ok {
		return callPlatform(hostPlatform.loadavg)
	}
	return formatLoadavg(a, b, c)
}

// DNS reads resolv.conf, following the systemd-resolved and NetworkManager stubs.
func DNS() string {
	primary := sysPath("etc", "resolv.conf")
	list := ParseDNS(readFile(primary))
	if len(list) == 1 && (list[0] == "127.0.0.53" || list[0] == "127.0.0.1") {
		alts := []string{
			sysPath("run", "systemd", "resolve", "resolv.conf"),
			sysPath("var", "run", "NetworkManager", "resolv.conf"),
			sysPath("run", "resolvconf", "resolv.conf"),
			sysPath("var", "run", "dnsmasq", "resolv.conf"),
		}
		for _, alt := range alts {
			if next := ParseDNS(readFile(alt)); len(next) > 0 {
				list = next
				break
			}
		}
	}
	if len(list) == 0 {
		for _, alt := range []string{
			sysPath("run", "resolvconf", "resolv.conf"),
			sysPath("etc", "resolv.conf.openresolv"),
		} {
			if next := ParseDNS(readFile(alt)); len(next) > 0 {
				list = next
				break
			}
		}
	}
	return finish(joinLines(list), hostPlatform.dns)
}

// InitSystem is PID 1's command name, plus a systemd version when it is cheap to read.
func InitSystem() string {
	name := readTrim("proc", "1", "comm")
	if name == "" {
		return callPlatform(hostPlatform.initSystem)
	}
	if name == "systemd" {
		if ver := systemdVersion(runCombined("systemctl", "--version")); ver != "" {
			return "systemd " + ver
		}
	}
	return name
}

// Keyboard lists real keyboards from the input device list.
func Keyboard() string {
	return finish(joinLines(ParseKeyboards(readFile(sysPath("proc", "bus", "input", "devices")))), hostPlatform.keyboard)
}

// Gamepad lists joystick devices.
func Gamepad() string {
	fromProc := ParseGamepads(readFile(sysPath("proc", "bus", "input", "devices")))
	if len(fromProc) > 0 {
		return joinLines(fromProc)
	}
	dir := sysPath("sys", "class", "input")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return callPlatform(hostPlatform.gamepad)
	}
	var lines []string
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "js") {
			continue
		}
		name := readTrim("sys", "class", "input", entry.Name(), "device", "name")
		if name != "" {
			lines = append(lines, name)
		}
	}
	return finish(joinLines(lines), hostPlatform.gamepad)
}

// Mouse lists mice from sysfs.
func Mouse() string {
	dir := sysPath("sys", "class", "input")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return callPlatform(hostPlatform.mouse)
	}
	var lines []string
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "mouse") {
			continue
		}
		name := readTrim("sys", "class", "input", entry.Name(), "device", "name")
		if name != "" {
			lines = append(lines, name)
		}
	}
	return finish(joinLines(lines), hostPlatform.mouse)
}

// Processes counts userspace tasks and their threads.
func Processes() string {
	procs := readProcs()
	if len(procs) == 0 && !dirExists(sysPath("proc")) {
		return callPlatform(hostPlatform.processes)
	}
	var n, threads uint64
	for _, p := range procs {
		if p.kernel() {
			continue
		}
		n++
		threads += p.threads
	}
	return formatProcesses(n, threads, true)
}

// Top is the busiest userspace processes over a short sample.
func Top() string {
	first := readProcs()
	if len(first) == 0 {
		return callPlatform(hostPlatform.top)
	}
	if usagePause > 0 {
		time.Sleep(usagePause)
	}
	second := readProcs()
	byPID := map[int]procStat{}
	for _, p := range first {
		byPID[p.pid] = p
	}
	type rank struct {
		name  string
		delta uint64
		cpu   uint64
	}
	var ranks []rank
	use := second
	if len(use) == 0 {
		use = first
	}
	for _, p := range use {
		if p.kernel() || p.name == "" {
			continue
		}
		prev, ok := byPID[p.pid]
		delta := p.cpu
		if ok && p.cpu >= prev.cpu {
			delta = p.cpu - prev.cpu
		}
		ranks = append(ranks, rank{name: p.name, delta: delta, cpu: p.cpu})
	}
	sort.Slice(ranks, func(i, j int) bool {
		if ranks[i].delta != ranks[j].delta {
			return ranks[i].delta > ranks[j].delta
		}
		return ranks[i].cpu > ranks[j].cpu
	})
	if len(ranks) > 5 {
		ranks = ranks[:5]
	}
	var lines []string
	for _, r := range ranks {
		lines = append(lines, r.name)
	}
	return joinLines(lines)
}

func readProcs() []procStat {
	dir := sysPath("proc")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []procStat
	for _, entry := range entries {
		if entry.Name() == "" || entry.Name()[0] < '0' || entry.Name()[0] > '9' {
			continue
		}
		text := readFile(filepath.Join(dir, entry.Name(), "stat"))
		stat, ok := ParseProcStat(text)
		if !ok {
			continue
		}
		out = append(out, stat)
	}
	return out
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// Sound lists ALSA cards, and the codec chip when the driver publishes one.
func Sound() string {
	cards := ParseSoundCards(readFile(sysPath("proc", "asound", "cards")))
	codecs := codecNames()
	if len(codecs) == 0 {
		return finish(joinLines(cards), hostPlatform.sound)
	}
	if len(cards) == 0 {
		return finish(joinLines(codecs), hostPlatform.sound)
	}
	return joinLines(append(cards, codecs...))
}

func codecNames() []string {
	dir := sysPath("proc", "asound")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "card") {
			continue
		}
		codecs, err := os.ReadDir(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		for _, codec := range codecs {
			if !strings.HasPrefix(codec.Name(), "codec") {
				continue
			}
			text := readFile(filepath.Join(dir, entry.Name(), codec.Name()))
			for _, line := range strings.Split(text, "\n") {
				if name, ok := strings.CutPrefix(strings.TrimSpace(line), "Codec:"); ok {
					name = strings.TrimSpace(name)
					if name != "" {
						out = append(out, name)
					}
					break
				}
			}
		}
	}
	return out
}

// TPM is the version and description of the first TPM device.
func TPM() string {
	version := readTrim("sys", "class", "tpm", "tpm0", "tpm_version_major")
	if version == "2" {
		version = "2.0"
	}
	desc := readTrim("sys", "class", "tpm", "tpm0", "device", "description")
	switch {
	case version != "" && desc != "":
		return "TPM " + version + " (" + desc + ")"
	case version != "":
		return "TPM " + version
	case desc != "":
		return desc
	default:
		return callPlatform(hostPlatform.tpm)
	}
}

// PhysicalMemory prefers the SMBIOS table and falls back to dmidecode text.
func PhysicalMemory() string {
	if lines := ParseMemoryDevices(readBytes("sys", "firmware", "dmi", "tables", "DMI")); len(lines) > 0 {
		return joinLines(lines)
	}
	if commandExists("dmidecode") {
		if lines := ParseDMIDecodeMemory(runCombined("dmidecode", "-t", "memory")); len(lines) > 0 {
			return joinLines(lines)
		}
	}
	return callPlatform(hostPlatform.physicalMemory)
}

type blockDev struct {
	name  string
	model string
	size  uint64
	read  uint64
	write uint64
}

func blockDevices() []blockDev {
	dir := sysPath("sys", "block")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []blockDev
	for _, entry := range entries {
		name := entry.Name()
		if skipBlock(name) {
			continue
		}
		base := filepath.Join(dir, name)
		if _, err := os.Stat(filepath.Join(base, "device")); err != nil {
			continue
		}
		dev := blockDev{name: name, model: blockModel(base, name)}
		if n, err := strconv.ParseUint(strings.TrimSpace(readFile(filepath.Join(base, "size"))), 10, 64); err == nil {
			dev.size = n * 512
		}
		dev.read, dev.write, _ = ParseDiskStat(readFile(filepath.Join(base, "stat")))
		out = append(out, dev)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

func skipBlock(name string) bool {
	switch {
	case strings.HasPrefix(name, "loop"),
		strings.HasPrefix(name, "ram"),
		strings.HasPrefix(name, "fd"),
		strings.HasPrefix(name, "dm-"),
		strings.HasPrefix(name, "md"),
		strings.HasPrefix(name, "zram"),
		strings.HasPrefix(name, "sr"),
		strings.HasPrefix(name, "nbd"):
		return true
	default:
		return false
	}
}

func blockModel(base, name string) string {
	vendor := strings.TrimSpace(readFile(filepath.Join(base, "device", "vendor")))
	model := strings.TrimSpace(readFile(filepath.Join(base, "device", "model")))
	if isPCIID(vendor) {
		vendor = ""
	}
	text := strings.TrimSpace(vendor + " " + model)
	if text == "" {
		if link, err := os.Readlink(filepath.Join(base, "device", "driver")); err == nil {
			if driver := filepath.Base(link); driver != "" && driver != "." {
				return driver + " " + name
			}
		}
		return name
	}
	return text + " " + name
}

func isPCIID(value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))
	if !strings.HasPrefix(value, "0x") || len(value) < 3 {
		return false
	}
	for _, r := range value[2:] {
		switch {
		case r >= '0' && r <= '9', r >= 'a' && r <= 'f':
		default:
			return false
		}
	}
	return true
}

// PhysicalDisk lists real block devices and their size.
func PhysicalDisk() string {
	devs := blockDevices()
	if len(devs) == 0 {
		return callPlatform(hostPlatform.physicalDisk)
	}
	var lines []string
	for _, dev := range devs {
		if dev.size == 0 {
			lines = append(lines, dev.model)
			continue
		}
		lines = append(lines, fmt.Sprintf("%s (%s)", dev.model, FormatBytes(dev.size)))
	}
	return joinLines(lines)
}

// DiskIO reports the read and write rate of each physical disk.
// With ioSample at zero it reports the totals since boot.
func DiskIO() string {
	first := blockDevices()
	if len(first) == 0 {
		return callPlatform(hostPlatform.diskIO)
	}
	second := map[string]blockDev{}
	span := ioSample
	if span > 0 {
		time.Sleep(span)
		for _, dev := range blockDevices() {
			second[dev.name] = dev
		}
	}
	var lines []string
	for _, dev := range first {
		read, write := dev.read, dev.write
		if span > 0 {
			nxt, ok := second[dev.name]
			if !ok {
				continue
			}
			read, write = 0, 0
			if nxt.read >= dev.read {
				read = nxt.read - dev.read
			}
			if nxt.write >= dev.write {
				write = nxt.write - dev.write
			}
		}
		lines = append(lines, formatDiskRate(dev.model, read, write, span.Seconds()))
	}
	return joinLines(lines)
}

// NetIO lists received and transmitted bytes on interfaces that are up.
func NetIO() string {
	dir := sysPath("sys", "class", "net")
	if !dirExists(dir) {
		return callPlatform(hostPlatform.netIO)
	}
	type link struct {
		rx uint64
		tx uint64
	}
	read := func() map[string]link {
		out := map[string]link{}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return out
		}
		for _, entry := range entries {
			name := entry.Name()
			if name == "lo" {
				continue
			}
			state := readTrim("sys", "class", "net", name, "operstate")
			if state != "" && state != "up" && state != "unknown" {
				continue
			}
			rx, err1 := strconv.ParseUint(readTrim("sys", "class", "net", name, "statistics", "rx_bytes"), 10, 64)
			tx, err2 := strconv.ParseUint(readTrim("sys", "class", "net", name, "statistics", "tx_bytes"), 10, 64)
			if err1 != nil || err2 != nil {
				continue
			}
			out[name] = link{rx: rx, tx: tx}
		}
		return out
	}
	first := read()
	span := ioSample
	second := first
	if span > 0 {
		time.Sleep(span)
		second = read()
	}
	var names []string
	for name := range second {
		names = append(names, name)
	}
	sort.Strings(names)
	var lines []string
	for _, name := range names {
		cur := second[name]
		prev, ok := first[name]
		rx, tx := cur.rx, cur.tx
		if span > 0 {
			if !ok {
				continue
			}
			rx, tx = 0, 0
			if cur.rx >= prev.rx {
				rx = cur.rx - prev.rx
			}
			if cur.tx >= prev.tx {
				tx = cur.tx - prev.tx
			}
		}
		lines = append(lines, formatNetRate(name, rx, tx, span.Seconds()))
	}
	return joinLines(lines)
}

// Monitor lists connected DRM connectors, with the EDID name when it is readable.
func Monitor() string {
	dir := sysPath("sys", "class", "drm")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return callPlatform(hostPlatform.monitor)
	}
	var lines []string
	for _, entry := range entries {
		name := entry.Name()
		if !strings.Contains(name, "-") {
			continue
		}
		status := readTrim("sys", "class", "drm", name, "status")
		if status != "" && status != "connected" {
			continue
		}
		if status == "" {
			continue
		}
		modes := readTrim("sys", "class", "drm", name, "modes")
		mode := firstLine(modes)
		edidName, w, h := ParseEDID(readBytes("sys", "class", "drm", name, "edid"))
		conn := name
		if i := strings.IndexByte(name, '-'); i >= 0 {
			conn = name[i+1:]
		}
		var parts []string
		if edidName != "" {
			parts = append(parts, edidName)
		}
		parts = append(parts, conn)
		if mode != "" {
			parts = append(parts, mode)
		}
		if w > 0 && h > 0 {
			inch := (float64(w)*float64(w) + float64(h)*float64(h))
			inch = sqrt(inch) / 25.4
			if inch >= 1 {
				parts = append(parts, fmt.Sprintf("%.0f inch", inch))
			}
		}
		lines = append(lines, strings.Join(parts, " "))
	}
	return finish(joinLines(lines), hostPlatform.monitor)
}

func sqrt(v float64) float64 {
	if v <= 0 {
		return 0
	}
	g := v
	for i := 0; i < 20; i++ {
		g = (g + v/g) / 2
	}
	return g
}

// Wifi reports the SSID and link quality of each wireless interface.
func Wifi() string {
	links := ParseWireless(readFile(sysPath("proc", "net", "wireless")))
	if len(links) == 0 {
		return callPlatform(hostPlatform.wifi)
	}
	var lines []string
	for _, link := range links {
		ssid := wifiSSID(link.iface)
		switch {
		case ssid != "" && link.quality != "":
			lines = append(lines, fmt.Sprintf("%s (%s, quality %s)", ssid, link.iface, link.quality))
		case ssid != "":
			lines = append(lines, fmt.Sprintf("%s (%s)", ssid, link.iface))
		default:
			lines = append(lines, fmt.Sprintf("%s (quality %s)", link.iface, link.quality))
		}
	}
	return joinLines(lines)
}

func wifiSSID(iface string) string {
	if commandExists("iwgetid") {
		if ssid := strings.TrimSpace(runOut("iwgetid", iface, "-r")); ssid != "" {
			return ssid
		}
	}
	if commandExists("iw") {
		text := runCombined("iw", "dev", iface, "link")
		for _, line := range strings.Split(text, "\n") {
			if name, ok := strings.CutPrefix(strings.TrimSpace(line), "SSID:"); ok {
				return strings.TrimSpace(name)
			}
		}
	}
	return ""
}

// Wallpaper looks through the desktop settings gf can read without a daemon.
func Wallpaper() string {
	if commandExists("gsettings") {
		for _, key := range []string{"picture-uri-dark", "picture-uri"} {
			if path := CleanURI(runOut("gsettings", "get", "org.gnome.desktop.background", key)); path != "" && path != "none" {
				return path
			}
		}
	}
	home, _ := os.UserHomeDir()
	if home == "" {
		return ""
	}
	candidates := []string{
		filepath.Join(home, ".config", "plasma-org.kde.plasma.desktop-appletsrc"),
		filepath.Join(home, ".config", "hypr", "hyprpaper.conf"),
		filepath.Join(home, ".fehbg"),
		filepath.Join(home, ".config", "xfce4", "xfconf", "xfce-perchannel-xml", "xfce4-desktop.xml"),
	}
	for _, path := range candidates {
		if found := wallpaperFile(path); found != "" {
			return found
		}
	}
	return callPlatform(hostPlatform.wallpaper)
}

func wallpaperFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		trim := strings.TrimSpace(line)
		for _, prefix := range []string{"Image=", "wallpaper =", "wallpaper=", "last-image"} {
			if strings.HasPrefix(trim, prefix) || strings.Contains(trim, prefix) {
				_, val, ok := strings.Cut(trim, "=")
				if !ok {
					continue
				}
				val = strings.Trim(strings.TrimSpace(val), `"'`)
				if strings.Contains(val, "last-image") {
					continue
				}
				if val != "" && val != "None" && (strings.Contains(val, "/") || strings.Contains(val, ".")) {
					return CleanURI(val)
				}
			}
		}
		if strings.Contains(trim, "feh") && strings.Contains(trim, "--bg") {
			fields := strings.Fields(trim)
			for i := len(fields) - 1; i >= 0; i-- {
				if strings.Contains(fields[i], "/") {
					return strings.Trim(fields[i], `"'`)
				}
			}
		}
	}
	return ""
}

// TerminalTheme uses COLORFGBG, which terminals export as "fg;bg".
func TerminalTheme() string {
	return ColorFGBG(os.Getenv("COLORFGBG"))
}

// DateTime is the local clock.
func DateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
