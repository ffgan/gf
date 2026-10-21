package cli

import "testing"

func TestIdentifyCPUARMClusters(t *testing.T) {
	cpuinfo := "" +
		"processor : 0\nCPU implementer : 0x41\nCPU part : 0xd0b\nmodel name : ARMv8 Processor rev 1 (v8l)\n\n" +
		"processor : 1\nCPU implementer : 0x41\nCPU part : 0xd0b\n\n" +
		"processor : 2\nCPU implementer : 0x41\nCPU part : 0xd0b\n\n" +
		"processor : 3\nCPU implementer : 0x41\nCPU part : 0xd0b\n\n" +
		"processor : 4\nCPU implementer : 0x41\nCPU part : 0xd05\n\n" +
		"processor : 5\nCPU implementer : 0x41\nCPU part : 0xd05\n\n" +
		"processor : 6\nCPU implementer : 0x41\nCPU part : 0xd05\n\n" +
		"processor : 7\nCPU implementer : 0x41\nCPU part : 0xd05\n"
	name, counts := IdentifyCPU(cpuinfo, "")
	if name != "4x Cortex-A76 + 4x Cortex-A55" || !counts {
		t.Fatalf("got %q counts=%v", name, counts)
	}
}

func TestIdentifyCPUKeepsX86Model(t *testing.T) {
	cpuinfo := "vendor_id : GenuineIntel\nmodel name : Intel(R) Xeon(R) Processor\nflags : hypervisor\n"
	name, counts := IdentifyCPU(cpuinfo, "")
	if name != "Intel(R) Xeon(R) Processor" || counts {
		t.Fatalf("got %q counts=%v", name, counts)
	}
}

func TestIdentifyCPUHygonAndZhaoxin(t *testing.T) {
	name, _ := IdentifyCPU("vendor_id : HygonGenuine\n", "")
	if name != "Hygon" {
		t.Fatalf("hygon: %q", name)
	}
	name, _ = IdentifyCPU("vendor_id : Shanghai\nmodel name : KaiXian KX-U6780A\n", "")
	if name != "Zhaoxin KaiXian KX-U6780A" {
		t.Fatalf("zhaoxin: %q", name)
	}
}

func TestIdentifyCPURISCVAndLoongson(t *testing.T) {
	name, counts := IdentifyCPU("processor : 0\nuarch : sifive,u74-mc\nisa : rv64imafdc\n", "")
	if name != "SiFive U74" || counts {
		t.Fatalf("riscv: %q", name)
	}
	name, _ = IdentifyCPU("Model Name : Loongson-3A5000\n", "")
	if name != "Loongson-3A5000" {
		t.Fatalf("loongson: %q", name)
	}
}

func TestIdentifyCPUUnknownPartAndDeviceTree(t *testing.T) {
	name, _ := IdentifyCPU("processor : 0\nCPU implementer : 0x70\nCPU part : 0x999\nmodel name : ARMv8 Processor rev 0 (v8l)\n", "")
	if name != "Phytium 0x999" {
		t.Fatalf("part: %q", name)
	}
	name, _ = IdentifyCPU("processor : 0\nmodel name : ARMv7 Processor rev 4 (v7l)\n", "Raspberry Pi 4 Model B\x00")
	if name != "Raspberry Pi 4" {
		t.Fatalf("dt: %q", name)
	}
	name, _ = IdentifyCPU("Hardware : Rockchip RK3588 EVB\n", "")
	if name != "Rockchip RK3588" {
		t.Fatalf("hw: %q", name)
	}
}

func TestIdentifyCPUAppleCode(t *testing.T) {
	name, _ := IdentifyCPU("model name : 8103\n", "")
	if name != "Apple M1" {
		t.Fatalf("got %q", name)
	}
	if got := applyFamilyCode("6050"); got != "Apple M5 Pro / Max" {
		t.Fatalf("m5: %q", got)
	}
}

func TestFormatCPUInfoSkipsClusterCount(t *testing.T) {
	got := FormatCPUInfo(&CPUInfo{Name: "4x Cortex-A76 + 4x Cortex-A55", Cores: 8, CountsInName: true}, CPUConfig{ShowCores: true})
	if got != "4x Cortex-A76 + 4x Cortex-A55" {
		t.Fatalf("got %q", got)
	}
}

func TestMacOSCodename(t *testing.T) {
	if macOSCodename("15.1") != "macOS Sequoia" || macOSCodename("26.0") != "macOS Tahoe" || macOSCodename("14.5") != "macOS Sonoma" {
		t.Fatal(macOSCodename("15.1"), macOSCodename("26.0"), macOSCodename("14.5"))
	}
	if macOSCodename("16.0") != "macOS" {
		t.Fatal(macOSCodename("16.0"))
	}
}

func TestGetOSDoesNotAbort(t *testing.T) {
	if got := GetOS("Fuchsia", ""); got != "Fuchsia" {
		t.Fatalf("unknown: %q", got)
	}
	if got := GetOS("Darwin", ""); got != "Darwin" {
		t.Fatalf("darwin: %q", got)
	}
	if got := GetOS("FreeBSD", ""); got != "FreeBSD" {
		t.Fatalf("freebsd: %q", got)
	}
	if got := GetOS("OpenBSD", ""); got != "OpenBSD" {
		t.Fatalf("openbsd: %q", got)
	}
	if got := GetOS("HP-UX", ""); got != "HP-UX" {
		t.Fatalf("hpux: %q", got)
	}
}
