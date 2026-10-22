package hyfetch

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

//go:embed data/presets.json
var presetJSON []byte

//go:embed data/font_logos.json
var fontLogoJSON []byte

// Presets are the flag palettes shipped with hyfetch.
var Presets map[string]ColorProfile

func init() {
	loaded, err := loadPresets(presetJSON)
	if err != nil {
		panic(err)
	}
	Presets = loaded
}

type presetFile map[string]presetEntry

type presetEntry struct {
	Colors  []string
	Weights []int
}

func (p *presetEntry) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '[' {
		return json.Unmarshal(b, &p.Colors)
	}
	var obj struct {
		Colors  []string `json:"colors"`
		Weights []int    `json:"weights"`
	}
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	p.Colors = obj.Colors
	p.Weights = obj.Weights
	return nil
}

func loadPresets(data []byte) (map[string]ColorProfile, error) {
	var file presetFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("presets: %w", err)
	}
	out := make(map[string]ColorProfile, len(file))
	for name, entry := range file {
		profile, err := ProfileFromHex(entry.Colors)
		if err != nil {
			return nil, fmt.Errorf("preset %s: %w", name, err)
		}
		if len(entry.Weights) > 0 {
			profile, err = profile.WithWeights(entry.Weights)
			if err != nil {
				return nil, fmt.Errorf("preset %s: %w", name, err)
			}
		}
		out[name] = profile
	}
	return out, nil
}

// PresetNames returns the sorted flag names plus the reserved "random" choice.
func PresetNames(extra map[string]ColorProfile) []string {
	names := make([]string, 0, len(Presets)+len(extra)+1)
	seen := map[string]bool{}
	for name := range Presets {
		names = append(names, name)
		seen[name] = true
	}
	for name := range extra {
		if !seen[name] {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return append(names, "random")
}

// ResolvePreset accepts a preset name, "random", a comma-separated list of
// names (one is chosen), or a comma-separated list of hex colors.
func ResolvePreset(spec string, custom map[string]ColorProfile) (ColorProfile, string, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		spec = "rainbow"
	}
	if strings.Contains(spec, "#") {
		parts := splitList(spec)
		profile, err := ProfileFromHex(parts)
		if err != nil {
			return ColorProfile{}, "", err
		}
		return profile, spec, nil
	}
	if strings.Contains(spec, ",") {
		parts := splitList(spec)
		if len(parts) == 0 {
			return ColorProfile{}, "", fmt.Errorf("empty preset list")
		}
		return ResolvePreset(parts[rand.Intn(len(parts))], custom)
	}
	profiles := mergedProfiles(custom)
	if spec == "random" {
		if len(profiles) == 0 {
			return ColorProfile{}, "", fmt.Errorf("preset list is empty")
		}
		names := make([]string, 0, len(profiles))
		for name := range profiles {
			names = append(names, name)
		}
		sort.Strings(names)
		name := names[rand.Intn(len(names))]
		return profiles[name], name, nil
	}
	if profile, ok := profiles[spec]; ok {
		return profile, spec, nil
	}
	return ColorProfile{}, "", fmt.Errorf("PRESET should be comma-separated hex colors or one of {%s}", strings.Join(PresetNames(custom), ","))
}

func mergedProfiles(custom map[string]ColorProfile) map[string]ColorProfile {
	out := make(map[string]ColorProfile, len(Presets)+len(custom))
	for name, profile := range Presets {
		out[name] = profile
	}
	for name, profile := range custom {
		out[name] = profile
	}
	return out
}

func splitList(spec string) []string {
	parts := strings.FieldsFunc(spec, func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\t'
	})
	var out []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
