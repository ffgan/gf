package hyfetch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPresetsLoadRainbow(t *testing.T) {
	if len(Presets) < 100 {
		t.Fatalf("presets = %d", len(Presets))
	}
	rainbow := Presets["rainbow"]
	if len(rainbow.Colors) != 6 || rainbow.Colors[0].R != 0xE5 || rainbow.Colors[0].G != 0 {
		t.Fatalf("rainbow = %+v", rainbow.Colors)
	}
	if _, ok := Presets["transgender"]; !ok {
		t.Fatal("missing transgender")
	}
	if _, ok := Presets["demisexual"]; !ok {
		t.Fatal("weighted preset missing")
	}
}

func TestWithLengthCentersOddExtra(t *testing.T) {
	profile := ColorProfile{Colors: []RGB{{1, 0, 0}, {2, 0, 0}, {3, 0, 0}}}
	spread, err := profile.WithLength(5)
	if err != nil {
		t.Fatal(err)
	}
	if len(spread.Colors) != 5 {
		t.Fatalf("len = %d", len(spread.Colors))
	}
	// repeats 1, extras 2 -> borders gain one each: 2, 1, 2
	got := []uint8{spread.Colors[0].R, spread.Colors[1].R, spread.Colors[2].R, spread.Colors[3].R, spread.Colors[4].R}
	want := []uint8{1, 1, 2, 3, 3}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("spread = %v, want %v", got, want)
		}
	}
}

func TestRecolorHorizontalUsesFlag(t *testing.T) {
	raw := "${c1}####\n####\n####\n"
	profile := Presets["transgender"]
	out, err := Recolor(raw, ColorAlignment{Mode: "horizontal"}, profile, ModeRGB, ThemeDark)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "${c") {
		t.Fatalf("placeholder left behind: %q", out)
	}
	if !strings.Contains(out, "\033[38;2;") {
		t.Fatalf("missing truecolor: %q", out)
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("lines = %d", len(lines))
	}
}

func TestResolveHexPreset(t *testing.T) {
	profile, _, err := ResolvePreset("#ff0000,#00ff00,#0000ff", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(profile.Colors) != 3 || profile.Colors[1].G != 255 {
		t.Fatalf("%+v", profile.Colors)
	}
	if _, _, err := ResolvePreset("not-a-flag", nil); err == nil {
		t.Fatal("expected unknown preset to fail")
	}
}

func TestConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hyfetch.json")
	cfg := DefaultConfig()
	cfg.Preset = PresetValue{Values: []string{"bisexual"}}
	cfg.Mode = ModeRGB
	cfg.ColorAlign = ColorAlignment{Mode: "vertical"}
	cfg.Args = ArgList{"--logo", "Arch"}
	cfg.PrideMonthDisable = true
	if err := cfg.Save(path); err != nil {
		t.Fatal(err)
	}
	loaded, ok, err := Load(path)
	if err != nil || !ok {
		t.Fatalf("load: ok=%v err=%v", ok, err)
	}
	if loaded.Preset.Pick() != "bisexual" || loaded.Mode != ModeRGB || loaded.ColorAlign.Mode != "vertical" {
		t.Fatalf("%+v", loaded)
	}
	if len(loaded.Args) != 2 || loaded.Args[0] != "--logo" {
		t.Fatalf("args = %#v", loaded.Args)
	}
}

func TestPythonStyleArgsString(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hyfetch.json")
	body := `{
  "preset": "rainbow",
  "mode": "8bit",
  "color_align": {"mode": "horizontal", "custom_colors": {}, "fore_back": []},
  "backend": "neofetch",
  "args": "--color_blocks off",
  "pride_month_disable": false
}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, ok, err := Load(path)
	if err != nil || !ok {
		t.Fatalf("load: %v", err)
	}
	if strings.Join(loaded.Args, " ") != "--color_blocks off" {
		t.Fatalf("args = %#v", loaded.Args)
	}
}

func TestSplitArgsLeavesNeofetchFlags(t *testing.T) {
	opt, rest, err := SplitArgs([]string{"-p", "lesbian", "--logo", "Arch", "--backend", "ascii", "-b", "fastfetch", "--args", "--cpu_speed off"})
	if err != nil {
		t.Fatal(err)
	}
	if opt.Preset != "lesbian" || opt.Backend != "fastfetch" {
		t.Fatalf("%+v", opt)
	}
	if strings.Join(rest, " ") != "--logo Arch --backend ascii" {
		t.Fatalf("rest = %#v", rest)
	}
	if strings.Join(opt.Args, " ") != "--cpu_speed off" {
		t.Fatalf("args = %#v", opt.Args)
	}
}

func TestLightnessMovesRed(t *testing.T) {
	in := RGB{R: 255}
	out := in.assignLightness(assignClampMin, 0.65)
	if out.R <= out.G || out.R <= out.B {
		t.Fatalf("not reddish: %+v", out)
	}
	if out == in {
		t.Fatalf("lightness did not change %+v", out)
	}
}

func TestShouldPrideRespectsMonth(t *testing.T) {
	cfg := DefaultConfig()
	if ShouldPride(cfg, false, time.Date(2026, time.October, 21, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("october should not auto-play")
	}
	cfg.PrideMonthDisable = true
	if ShouldPride(cfg, true, time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("disabled config should stay quiet")
	}
}

func TestFontLogoUbuntu(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	glyph, err := FontLogo("Ubuntu 24.04.5 LTS x86_64")
	if err != nil {
		t.Fatal(err)
	}
	if glyph == "" {
		t.Fatal("empty glyph")
	}
}

func TestAnsi256Red(t *testing.T) {
	if got := (RGB{255, 0, 0}).ansi256(); got != 196 {
		t.Fatalf("ansi = %d", got)
	}
}
