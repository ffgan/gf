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

func TestParsePrin(t *testing.T) {
	cfg := DefaultConfig()
	data := `
print_info() {
    prin "Hello" "world"
    prin "just text"
    info "Board" board
}
`
	if err := ParseConfig(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if len(cfg.InfoList) != 3 {
		t.Fatalf("lines = %#v", cfg.InfoList)
	}
	if cfg.InfoList[0][0] != "Hello" || cfg.InfoList[0][1] != "prin" || cfg.InfoList[0][2] != "world" {
		t.Fatalf("prin = %#v", cfg.InfoList[0])
	}
	if cfg.InfoList[1][0] != "prin" || cfg.InfoList[1][1] != "just text" {
		t.Fatalf("bare prin = %#v", cfg.InfoList[1])
	}
	if cfg.InfoList[2][1] != "board" {
		t.Fatalf("board = %#v", cfg.InfoList[2])
	}
	if err := ParseConfig("print_info() {\n    info title\n}\nweather_location=\"Singapore\"\ncommand=\"uname -r\"\n", &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.WeatherLocation != "Singapore" || cfg.CommandText != "uname -r" {
		t.Fatalf("fastfetch config = %#v", cfg.Fastfetch)
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
