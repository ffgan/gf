package hyfetch

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// Config is a hyfetch.json document.
type Config struct {
	Preset              PresetValue         `json:"preset"`
	Mode                AnsiMode            `json:"mode"`
	AutoDetectLightDark *bool               `json:"auto_detect_light_dark,omitempty"`
	LightDark           Theme               `json:"light_dark,omitempty"`
	Lightness           *float64            `json:"lightness,omitempty"`
	LightnessPerTheme   *LightnessPerTheme  `json:"lightness_per_theme,omitempty"`
	ColorAlign          ColorAlignment      `json:"color_align"`
	Backend             string              `json:"backend"`
	Args                ArgList             `json:"args,omitempty"`
	Distro              string              `json:"distro,omitempty"`
	PrideMonthDisable   bool                `json:"pride_month_disable"`
	CustomASCIIPath     string              `json:"custom_ascii_path,omitempty"`
	CustomPresets       map[string][]string `json:"custom_presets,omitempty"`
}

// LightnessPerTheme stores one OKHSL lightness for each background.
type LightnessPerTheme struct {
	Dark  float64 `json:"dark"`
	Light float64 `json:"light"`
}

func (l LightnessPerTheme) Get(theme Theme) float64 {
	if theme == ThemeLight {
		return l.Light
	}
	return l.Dark
}

// ArgList accepts hyfetch's string form ("--logo Arch") and a JSON array.
type ArgList []string

func (a ArgList) MarshalJSON() ([]byte, error) {
	if a == nil {
		return []byte("null"), nil
	}
	return json.Marshal(shellJoin([]string(a)))
}

func (a *ArgList) UnmarshalJSON(b []byte) error {
	text := strings.TrimSpace(string(b))
	if text == "" || text == "null" {
		*a = nil
		return nil
	}
	if text[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if strings.TrimSpace(s) == "" {
			*a = nil
			return nil
		}
		words, err := shellSplit(s)
		if err != nil {
			return err
		}
		*a = words
		return nil
	}
	var list []string
	if err := json.Unmarshal(b, &list); err != nil {
		return err
	}
	*a = list
	return nil
}

func shellJoin(args []string) string {
	parts := make([]string, len(args))
	for i, arg := range args {
		if arg == "" || strings.ContainsAny(arg, " \t\n\"'\\") {
			parts[i] = "'" + strings.ReplaceAll(arg, "'", `'\''`) + "'"
		} else {
			parts[i] = arg
		}
	}
	return strings.Join(parts, " ")
}

// PresetValue is either one preset or a list hyfetch picks from at random.
type PresetValue struct {
	Values []string
}

func (p PresetValue) Pick() string {
	switch len(p.Values) {
	case 0:
		return "rainbow"
	case 1:
		return p.Values[0]
	default:
		return p.Values[rand.Intn(len(p.Values))]
	}
}

func (p PresetValue) MarshalJSON() ([]byte, error) {
	if len(p.Values) == 1 {
		return json.Marshal(p.Values[0])
	}
	return json.Marshal(p.Values)
}

func (p *PresetValue) UnmarshalJSON(b []byte) error {
	b = []byte(strings.TrimSpace(string(b)))
	if len(b) == 0 || string(b) == "null" {
		p.Values = nil
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		p.Values = []string{s}
		return nil
	}
	var list []string
	if err := json.Unmarshal(b, &list); err != nil {
		return err
	}
	p.Values = list
	return nil
}

// ColorAlignment is how flag colors land on the logo.
type ColorAlignment struct {
	Mode         string      `json:"mode"`
	CustomColors map[int]int `json:"custom_colors,omitempty"`
	ForeBack     []int       `json:"fore_back,omitempty"`
}

func (c *ColorAlignment) UnmarshalJSON(b []byte) error {
	var raw struct {
		Mode         string          `json:"mode"`
		CustomColors json.RawMessage `json:"custom_colors"`
		ForeBack     []int           `json:"fore_back"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	c.Mode = strings.ToLower(raw.Mode)
	c.ForeBack = raw.ForeBack
	c.CustomColors = map[int]int{}
	if len(raw.CustomColors) == 0 || string(raw.CustomColors) == "null" {
		return nil
	}
	var asMap map[string]int
	if err := json.Unmarshal(raw.CustomColors, &asMap); err == nil {
		for k, v := range asMap {
			var n int
			if _, err := fmt.Sscan(k, &n); err == nil {
				c.CustomColors[n] = v
			}
		}
		return nil
	}
	var asList []int
	if err := json.Unmarshal(raw.CustomColors, &asList); err == nil {
		for i, v := range asList {
			c.CustomColors[i+1] = v
		}
	}
	return nil
}

func DefaultConfig() Config {
	return Config{
		Preset:     PresetValue{Values: []string{"rainbow"}},
		Mode:       Mode8bit,
		LightDark:  ThemeDark,
		ColorAlign: ColorAlignment{Mode: "horizontal"},
		Backend:    "neofetch",
	}
}

// ConfigPath is ~/.config/hyfetch.json, or $XDG_CONFIG_HOME/hyfetch.json.
func ConfigPath() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "hyfetch.json")
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "hyfetch.json"
	}
	return filepath.Join(home, ".config", "hyfetch.json")
}

// Load reads a hyfetch config. exists is false when the file is absent.
func Load(path string) (Config, bool, error) {
	if path == "" {
		path = ConfigPath()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), false, nil
		}
		return Config{}, false, fmt.Errorf("open %s: %w", path, err)
	}
	cfg := DefaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, true, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := cfg.normalize(); err != nil {
		return Config{}, true, err
	}
	return cfg, true, nil
}

func (c *Config) normalize() error {
	if len(c.Preset.Values) == 0 {
		c.Preset.Values = []string{"rainbow"}
	}
	if c.Mode == "" {
		c.Mode = Mode8bit
	}
	mode, err := ParseMode(string(c.Mode))
	if err != nil {
		return err
	}
	c.Mode = mode
	switch strings.ToLower(string(c.LightDark)) {
	case "", "dark":
		c.LightDark = ThemeDark
	case "light":
		c.LightDark = ThemeLight
	default:
		return fmt.Errorf("light_dark should be light or dark")
	}
	if c.ColorAlign.Mode == "" {
		c.ColorAlign.Mode = "horizontal"
	}
	switch c.ColorAlign.Mode {
	case "horizontal", "vertical", "custom", "random":
	default:
		return fmt.Errorf("color_align.mode should be horizontal, vertical, custom, or random")
	}
	if c.Backend == "" {
		c.Backend = "neofetch"
	}
	switch strings.ToLower(c.Backend) {
	case "neofetch", "fastfetch", "macchina":
		c.Backend = strings.ToLower(c.Backend)
	default:
		return fmt.Errorf("backend should be neofetch, fastfetch, or macchina")
	}
	if c.Lightness != nil && (*c.Lightness < 0 || *c.Lightness > 1) {
		return fmt.Errorf("lightness must be between 0 and 1")
	}
	for name := range c.CustomPresets {
		if name == "random" {
			return fmt.Errorf("custom preset key `random` is reserved")
		}
	}
	return nil
}

func (c Config) CustomProfiles() (map[string]ColorProfile, error) {
	out := map[string]ColorProfile{}
	for name, colors := range c.CustomPresets {
		profile, err := ProfileFromHex(colors)
		if err != nil {
			return nil, fmt.Errorf("custom preset %s: %w", name, err)
		}
		out[name] = profile
	}
	return out, nil
}

// Save writes the config as indented JSON, creating parent directories.
func (c Config) Save(path string) error {
	if path == "" {
		path = ConfigPath()
	}
	if err := c.normalize(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

// EffectiveLightness is the OKHSL lightness hyfetch applies when the command
// line did not pass --c-scale or --c-set-l.
func (c Config) EffectiveLightness(theme Theme) float64 {
	if c.LightnessPerTheme != nil {
		return c.LightnessPerTheme.Get(theme)
	}
	if c.Lightness != nil {
		return *c.Lightness
	}
	return DefaultLightness(theme)
}
