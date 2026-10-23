package modules

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStructureCoversFastfetch(t *testing.T) {
	want := []string{
		"battery", "bios", "bluetooth", "bluetoothradio", "board", "bootmgr", "break",
		"brightness", "btrfs", "camera", "chassis", "codec", "colors", "command",
		"cpu", "cpucache", "cpuusage", "cursor", "custom", "datetime", "de", "disk",
		"diskio", "display", "dns", "editor", "font", "gamepad", "gpu", "host",
		"icons", "initsystem", "kernel", "keyboard", "lm", "loadavg", "locale",
		"localip", "logo", "media", "memory", "monitor", "mouse", "netio", "opencl",
		"opengl", "os", "packages", "physicaldisk", "physicalmemory", "player",
		"poweradapter", "processes", "publicip", "separator", "shell", "sound",
		"swap", "terminal", "terminalfont", "terminalsize", "terminaltheme", "theme",
		"title", "top", "tpm", "uptime", "users", "version", "vulkan", "wallpaper",
		"weather", "wifi", "wm", "wmtheme", "zpool",
	}
	got := map[string]bool{}
	for _, name := range FastfetchNames() {
		got[name] = true
	}
	for _, name := range want {
		if !got[name] {
			t.Errorf("missing fastfetch module %s", name)
		}
	}
	lines, err := Structure("all")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) < len(want) {
		t.Fatalf("all = %d lines, want at least %d", len(lines), len(want))
	}
	if _, err := Structure("not-a-module"); err == nil {
		t.Fatal("expected an error for an unknown module")
	}
	one, err := Structure("OS:Host:CPU")
	if err != nil {
		t.Fatal(err)
	}
	if one[0][0] != "OS" || one[0][1] != "os" || Canonical(one[0][1]) != "distro" {
		t.Fatalf("os line = %#v canonical %s", one[0], Canonical("os"))
	}
	if Canonical("terminalfont") != "term_font" || Canonical("Wi-Fi") != "wifi" {
		t.Fatalf("aliases: %s %s", Canonical("terminalfont"), Canonical("Wi-Fi"))
	}
}

func TestParsers(t *testing.T) {
	dns := ParseDNS("# comment\nnameserver 1.1.1.1\nnameserver 2606:4700:4700::1111 # quad9\n")
	if len(dns) != 2 || dns[0] != "1.1.1.1" || dns[1] != "2606:4700:4700::1111" {
		t.Fatalf("dns = %#v", dns)
	}
	a, b, c, ok := ParseLoadavg("0.15 0.10 0.08 1/123 99")
	if !ok || formatLoadavg(a, b, c) != "0.15, 0.10, 0.08" {
		t.Fatalf("load = %v %v %v %v", a, b, c, ok)
	}
	agg, cores, ok := ParseCPUStat("cpu  100 0 50 1000 10 5 5 0\ncpu0 40 0 20 400 4 2 2 0\ncpu1 60 0 30 600 6 3 3 0\n")
	if !ok || len(cores) != 2 {
		t.Fatalf("cpu ok=%v cores=%d", ok, len(cores))
	}
	next := cpuSample{inUse: agg.inUse + 25, total: agg.total + 100}
	if usagePercent(agg, next) != "25%" {
		t.Fatalf("usage = %s", usagePercent(agg, next))
	}
	cards := ParseSoundCards(" 0 [PCH            ]: HDA-Intel - HDA Intel PCH\n                      ALC255 Analog\n")
	if len(cards) != 1 || cards[0] != "HDA Intel PCH" {
		t.Fatalf("cards = %#v", cards)
	}
	links := ParseWireless("Inter-| sta-\n face | tus\nwlp3s0: 0000   70.  -40.  -256\n")
	if len(links) != 1 || links[0].iface != "wlp3s0" || links[0].quality != "70" {
		t.Fatalf("wifi = %#v", links)
	}
	read, write, ok := ParseDiskStat("123 0 456 0 789 0 1000 0")
	if !ok || read != 456*512 || write != 1000*512 {
		t.Fatalf("disk %d %d %v", read, write, ok)
	}
	kb := ParseKeyboards("N: Name=\"AT Translated Set 2 keyboard\"\nH: Handlers=sysrq kbd event0\nB: EV=120013\nB: KEY=402000000 3803078f800d001 feffffdfffefffff fffffffffffffffe\n\nN: Name=\"Power Button\"\nH: Handlers=kbd event1\nB: EV=3\nB: KEY=0\n\n")
	if len(kb) != 1 || kb[0] != "AT Translated Set 2 keyboard" {
		t.Fatalf("keyboards = %#v", kb)
	}
	stat, ok := ParseProcStat("12 (my proc) S 1 1 1 0 -1 4202496 0 0 0 0 10 5 0 0 20 0 3 0 0")
	if !ok || stat.name != "my proc" || stat.cpu != 15 || stat.threads != 3 || stat.kernel() {
		t.Fatalf("stat = %#v %v", stat, ok)
	}
	if FormatBytes(8192*1024*1024) != "8 GiB" {
		t.Fatalf("bytes = %s", FormatBytes(8192*1024*1024))
	}
	if BrightnessPercent("80", "100", "intel_backlight") != "80% [intel_backlight]" {
		t.Fatal(BrightnessPercent("80", "100", "intel_backlight"))
	}
	if ColorFGBG("15;0") != "fg #ffffff / bg #000000" {
		t.Fatal(ColorFGBG("15;0"))
	}
	if CleanURI("'file:///usr/share/bg.jpg'") != "/usr/share/bg.jpg" {
		t.Fatal(CleanURI("'file:///usr/share/bg.jpg'"))
	}
	if systemdVersion("systemd 255 (255.4)\n+PAM\n") != "255" {
		t.Fatal(systemdVersion("systemd 255 (255.4)\n"))
	}
	if ParseLibc("glibc 2.39", "") != "glibc 2.39" {
		t.Fatal(ParseLibc("glibc 2.39", ""))
	}
	mem := ParseDMIDecodeMemory("Handle 0x0001\n\tSize: 8 GB\n\tLocator: DIMM A\n\tType: DDR4\n\tSpeed: 3200 MT/s\n\tManufacturer: Samsung\n")
	if len(mem) != 1 || !strings.Contains(mem[0], "8 GB") || !strings.Contains(mem[0], "Samsung") {
		t.Fatalf("dmi text = %#v", mem)
	}
	if ChassisType(9) != "Laptop" || ChassisType(0x89) != "Laptop" {
		t.Fatal(ChassisType(9))
	}
}

func TestBootAndMemory(t *testing.T) {
	desc := []uint16{'B', 'o', 'o', 't'}
	body := make([]byte, 6+len(desc)*2+2)
	binary.LittleEndian.PutUint32(body[0:], 1)
	binary.LittleEndian.PutUint16(body[4:], 0)
	for i, c := range desc {
		binary.LittleEndian.PutUint16(body[6+i*2:], c)
	}
	file := append([]byte{1, 0, 0, 0}, body...)
	if ParseBootName(file) != "Boot" {
		t.Fatalf("boot name %q", ParseBootName(file))
	}
	current := []byte{0, 0, 0, 0, 0x34, 0x12}
	order, ok := ParseBootOrder(current)
	if !ok || order != 0x1234 {
		t.Fatalf("order %x %v", order, ok)
	}
	if on, ok := ParseSecureBoot([]byte{0, 0, 0, 0, 1}); !ok || !on {
		t.Fatal("secure boot")
	}

	dev := make([]byte, 0x22)
	dev[0] = 17
	dev[1] = 0x22
	binary.LittleEndian.PutUint16(dev[0x0C:], 8192)
	dev[0x0E] = 0x09
	dev[0x10] = 1
	dev[0x12] = 0x1A
	binary.LittleEndian.PutUint16(dev[0x15:], 3200)
	dev[0x17] = 2
	raw := append(dev, []byte("DIMM A\x00Samsung\x00\x00")...)
	raw = append(raw, []byte{127, 4, 0, 0, 0, 0}...)
	lines := ParseMemoryDevices(raw)
	if len(lines) != 1 || !strings.Contains(lines[0], "8 GiB") || !strings.Contains(lines[0], "DDR4-3200") || !strings.Contains(lines[0], "Samsung") {
		t.Fatalf("memory = %#v", lines)
	}
}

func TestEDIDName(t *testing.T) {
	edid := make([]byte, 128)
	edid[0] = 0x00
	edid[1] = 0xFF
	edid[21] = 60
	edid[22] = 34
	edid[54] = 0
	edid[55] = 0
	edid[56] = 0
	edid[57] = 0xFC
	copy(edid[59:], []byte("Dell U2720\n  "))
	name, w, h := ParseEDID(edid)
	if name != "Dell U2720" || w != 600 || h != 340 {
		t.Fatalf("edid %q %d %d", name, w, h)
	}
}

func TestCPUCacheFixture(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "sys", "devices", "system", "cpu", "cpu0", "cache", "index0")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, value string) {
		if err := os.WriteFile(filepath.Join(base, name), []byte(value), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("level", "1\n")
	write("size", "32K\n")
	write("type", "Data\n")
	write("shared_cpu_list", "0\n")
	other := filepath.Join(root, "sys", "devices", "system", "cpu", "cpu1", "cache", "index0")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"level", "size", "type"} {
		data, err := os.ReadFile(filepath.Join(base, name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(other, name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(other, "shared_cpu_list"), []byte("1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	old := sysRoot
	sysRoot = root
	t.Cleanup(func() { sysRoot = old })
	got := CPUCache()
	if got != "L1: 2x 32 KiB (D)" {
		t.Fatalf("cache = %q", got)
	}
}

func TestHostParsers(t *testing.T) {
	if got := ParseDNSList("1.1.1.1, 8.8.8.8; 2606:4700:4700::1111"); len(got) != 3 || got[2] != "2606:4700:4700::1111" {
		t.Fatalf("dns list %#v", got)
	}
	a, b, c, ok := ParseBraceLoadavg("{ 1.50 0.80 0.40 }")
	if !ok || formatLoadavg(a, b, c) != "1.50, 0.80, 0.40" {
		t.Fatalf("brace load %v %v", ok, formatLoadavg(a, b, c))
	}
	sample, ok := ParseCPTime("10 1 5 2 82")
	if !ok || sample.inUse != 18 || sample.total != 100 {
		t.Fatalf("cp_time %#v %v", sample, ok)
	}
	links := ParseNetstatIB("Name Mtu Network Address Ipkts Ierrs Ibytes Opkts Oerrs Obytes Coll\nen0 1500 <Link> aa 1 0 2048 1 0 1024 0\nlo0 16384 <Link> 1 0 10 1 0 10 0\n")
	if len(links) != 1 || links[0].name != "en0" || links[0].rx != 2048 || links[0].tx != 1024 {
		t.Fatalf("netstat %#v", links)
	}
	rates := ParseIostat("device r/s w/s kr/s kw/s\nada0 1 2 4 8\n")
	if len(rates) != 1 || rates[0].name != "ada0" || rates[0].read != 4*1024 || rates[0].write != 8*1024 {
		t.Fatalf("iostat %#v", rates)
	}
	mac := ParseIostat("              disk0\n    KB/t tps  MB/s\n   0.00   0  1.00\n")
	if len(mac) != 1 || !mac[0].combined || mac[0].read != 1024*1024 {
		t.Fatalf("darwin iostat %#v", mac)
	}
	disks := ParseDiskutil("**********\nWhole: Yes\nVirtual: No\nDevice / Media Name: APPLE SSD\nDisk Size: 500.3 GB (1048576 Bytes)\n")
	if len(disks) != 1 || disks[0] != "APPLE SSD (1 MiB)" {
		t.Fatalf("diskutil %#v", disks)
	}
	geom := ParseGeomDisk("Geom name: ada0\nMediasize: 1048576 (1M)\ndescr: Samsung SSD\n")
	if len(geom) != 1 || geom[0] != "Samsung SSD ada0 (1 MiB)" {
		t.Fatalf("geom %#v", geom)
	}
	if cams := ParseCamcontrol("<Samsung SSD 860> at scbus0 (ada0,pass0)\n"); len(cams) != 1 || cams[0] != "Samsung SSD 860" {
		t.Fatalf("cam %#v", cams)
	}
	wifi := ParseNetshWifi("Name : Wi-Fi\nState : connected\nSSID : Home\nSignal : 90%\n")
	if len(wifi) != 1 || wifi[0] != "Home (90%)" {
		t.Fatalf("netsh %#v", wifi)
	}
	if ParseAirportSSID("Current Wi-Fi Network: Cafe") != "Cafe" || ParseAirportSSID("You are not associated with an AirPort network.") != "" {
		t.Fatal("airport")
	}
	if CleanPNPDesc(`@keyboard.inf,%std%;USB Keyboard`) != "USB Keyboard" {
		t.Fatal(CleanPNPDesc(`@keyboard.inf,%std%;USB Keyboard`))
	}
	if !LooksLikeGamepad("Xbox Wireless Controller") || AppleChassis("MacBookPro18,3") != "Laptop" {
		t.Fatal("names")
	}
	vendor, model, board := ParseIORegPlatform("\"manufacturer\" = <\"Apple Inc.\">\n\"model\" = <4d6163426f6f6b50726f31382c33>\n\"board-id\" = <\"Mac-ABC\">\n")
	if vendor != "Apple Inc." || model != "MacBookPro18,3" || board != "Mac-ABC" {
		t.Fatalf("ioreg %q %q %q", vendor, model, board)
	}
	mem := FormatMemoryProfile("Memory:\n\n    BANK 0:\n\n      Size: 8 GB\n      Type: DDR4\n      Speed: 3200 MHz\n      Manufacturer: Samsung\n")
	if len(mem) != 1 || mem[0] != "8 GB DDR4-3200 (Samsung)" {
		t.Fatalf("memory profile %#v", mem)
	}
	if cache := CacheFromSysctl(4, 32*1024, 32*1024, 256*1024, 8*1024*1024); !strings.Contains(cache, "L1:") || !strings.Contains(cache, "4x") || !strings.Contains(cache, "L3:") {
		t.Fatalf("cache sysctl %q", cache)
	}
	top := ParsePSTop(" 20.0 WindowServer\n  1.0 kernel_task\n  0.4 Finder\n")
	if len(top) != 2 || top[0] != "WindowServer" || top[1] != "Finder" {
		t.Fatalf("ps top %#v", top)
	}
	n, threads, have := ParsePSCount("launchd 2\nFinder 4\n")
	if n != 2 || threads != 6 || !have {
		t.Fatalf("ps count %d %d %v", n, threads, have)
	}
	boardBody := append([]byte{2, 8, 0, 0, 1, 2, 3, 0}, []byte("Vendor\x00Board\x00v1\x00\x00")...)
	chassisBody := append([]byte{3, 6, 0, 0, 1, 9}, []byte("Frame\x00\x00")...)
	id := ParseSMBIOSIDs(append(boardBody, chassisBody...))
	if id.Board() != "Vendor Board v1" || id.Chassis() != "Laptop Frame" {
		t.Fatalf("smbios %q %q", id.Board(), id.Chassis())
	}
	if snd := ParseSndstat("pcm0: <Realtek ALC1220> (play) default\n"); len(snd) != 1 || snd[0] != "Realtek ALC1220 (default)" {
		t.Fatalf("sndstat %#v", snd)
	}
	if ssid := ParseIfconfigSSID("wlan0: flags=8843\n\tssid Home\n"); len(ssid) != 1 || ssid[0] != "Home (wlan0)" {
		t.Fatalf("ssid %#v", ssid)
	}
}

func TestDiskIOTotals(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "sys", "block", "sda", "device")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "model"), []byte("Samsung SSD\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sys", "block", "sda", "size"), []byte("2048\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sys", "block", "sda", "stat"), []byte("0 0 100 0 0 0 50 0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldRoot, oldSample := sysRoot, ioSample
	sysRoot, ioSample = root, 0
	t.Cleanup(func() {
		sysRoot = oldRoot
		ioSample = oldSample
	})
	if got := PhysicalDisk(); got != "Samsung SSD sda (1 MiB)" {
		t.Fatalf("physical %q", got)
	}
	if got := DiskIO(); got != "Samsung SSD sda: 50 KiB (R) - 25 KiB (W)" {
		t.Fatalf("diskio %q", got)
	}
}

func TestHostLoadavg(t *testing.T) {
	if _, err := os.Stat("/proc/loadavg"); err != nil {
		t.Skip(err)
	}
	if sysRoot != "/" {
		t.Fatal("sysRoot leaked")
	}
	got := Loadavg()
	if strings.Count(got, ",") != 2 {
		t.Fatalf("loadavg = %q", got)
	}
}
