package cli

import "testing"

func TestFormatSpeed(t *testing.T) {
	if got := formatSpeed(2400, true); got != "2.4GHz" {
		t.Fatalf("short = %q", got)
	}
	if got := formatSpeed(2400, false); got != "2.400GHz" {
		t.Fatalf("long = %q", got)
	}
	if got := formatSpeed(800, false); got != "800MHz" {
		t.Fatalf("mhz = %q", got)
	}
}

func TestFormatCPUInfo(t *testing.T) {
	got := FormatCPUInfo(&CPUInfo{Name: "Intel(R) Xeon(R) Processor", Cores: 4, Speed: 2400}, CPUConfig{
		ShowCores:  true,
		ShowSpeed:  true,
		ShowBrand:  true,
		SpeedShort: true,
	})
	if got != "Intel Xeon (4) @ 2.4GHz" {
		t.Fatalf("cpu = %q", got)
	}
}

func TestHumanBytes(t *testing.T) {
	if got := humanBytes(7*1024*1024*1024 + 512*1024*1024); got != "7.5G" {
		t.Fatalf("bytes = %q", got)
	}
}

func TestParseShowList(t *testing.T) {
	got := parseShowList("('/' '/home')")
	if len(got) != 2 || got[0] != "/" || got[1] != "/home" {
		t.Fatalf("mounts = %#v", got)
	}
}

func TestRenderBlocks(t *testing.T) {
	rows := RenderBlocks(3, 1, 0, 15)
	if len(rows) != 2 {
		t.Fatalf("rows = %d", len(rows))
	}
	if rows[0] == "" || rows[1] == "" {
		t.Fatal("empty color row")
	}
}
