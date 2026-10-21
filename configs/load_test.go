package configs

import "testing"

func TestParseConfigOverridesDefaults(t *testing.T) {
	cfg := DefaultConfig()
	data := `
# comment
print_info() {
    info title
    info "OS" distro
    info cols
}

kernel_shorthand="off" # keep the name
memory_unit="mib"
color_blocks="off"
disk_show=('/')
`
	if err := ParseConfig(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.KernelShorthand != "off" {
		t.Fatalf("kernel = %q", cfg.KernelShorthand)
	}
	if cfg.MemoryUnit != "mib" {
		t.Fatalf("unit = %q", cfg.MemoryUnit)
	}
	if cfg.Enabled != "off" {
		t.Fatalf("blocks = %q", cfg.Enabled)
	}
	if len(cfg.InfoList) != 3 || cfg.InfoList[1][0] != "OS" || cfg.InfoList[1][1] != "distro" {
		t.Fatalf("info = %#v", cfg.InfoList)
	}
}

func TestParseConfigWithoutPrintInfo(t *testing.T) {
	cfg := DefaultConfig()
	before := len(cfg.InfoList)
	if err := ParseConfig("gap=\"4\"\n", &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Gap != "4" {
		t.Fatalf("gap = %q", cfg.Gap)
	}
	if len(cfg.InfoList) != before {
		t.Fatalf("info list changed: %d", len(cfg.InfoList))
	}
}
