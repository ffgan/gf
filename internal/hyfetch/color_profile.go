package hyfetch

import (
	"fmt"
	"strings"
)

// ColorProfile is an ordered flag palette, the same structure hyfetch spreads
// across ASCII art.
type ColorProfile struct {
	Colors []RGB
}

func ProfileFromHex(colors []string) (ColorProfile, error) {
	if len(colors) == 0 {
		return ColorProfile{}, fmt.Errorf("hex color list cannot be empty")
	}
	out := make([]RGB, 0, len(colors))
	for _, raw := range colors {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		if !strings.HasPrefix(raw, "#") && len(raw) != 3 && len(raw) != 6 {
			return ColorProfile{}, fmt.Errorf("invalid hex color: %s", raw)
		}
		c, err := ParseHex(raw)
		if err != nil {
			return ColorProfile{}, err
		}
		out = append(out, c)
	}
	if len(out) == 0 {
		return ColorProfile{}, fmt.Errorf("hex color list cannot be empty")
	}
	return ColorProfile{Colors: out}, nil
}

func (p ColorProfile) WithWeights(weights []int) (ColorProfile, error) {
	if len(weights) != len(p.Colors) {
		return ColorProfile{}, fmt.Errorf("weights should have the same number of elements as colors")
	}
	var out []RGB
	for i, w := range weights {
		if w < 0 {
			return ColorProfile{}, fmt.Errorf("negative weight")
		}
		for n := 0; n < w; n++ {
			out = append(out, p.Colors[i])
		}
	}
	if len(out) == 0 {
		return ColorProfile{}, fmt.Errorf("weights produced an empty profile")
	}
	return ColorProfile{Colors: out}, nil
}

// WithLength spreads the palette to length the way hyfetch does: every color
// is repeated, the center takes a leftover odd slot, and the borders share the
// remaining even slots.
func (p ColorProfile) WithLength(length int) (ColorProfile, error) {
	if length <= 0 {
		return ColorProfile{}, fmt.Errorf("length must be positive")
	}
	if len(p.Colors) == 0 {
		return ColorProfile{}, fmt.Errorf("color profile is empty")
	}
	orig := len(p.Colors)
	repeats := length / orig
	weights := make([]int, orig)
	for i := range weights {
		weights[i] = repeats
	}
	extras := length % orig
	center := orig / 2
	if extras%2 == 1 {
		weights[center]++
		extras--
	}
	for border := 0; border < extras/2; border++ {
		weights[border]++
		weights[orig-1-border]++
	}
	return p.WithWeights(weights)
}

// ColorText paints each character. spaceOnly leaves non-spaces untouched, which
// is how hyfetch draws flag previews on block rows.
func (p ColorProfile) ColorText(txt string, mode AnsiMode, foreground, spaceOnly bool) (string, error) {
	runes := []rune(txt)
	if len(runes) == 0 {
		return "", nil
	}
	spread, err := p.WithLength(len(runes))
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	for i, r := range runes {
		if spaceOnly && r != ' ' {
			if i > 0 && runes[i-1] == ' ' {
				buf.WriteString("\033[39;49m")
			}
			buf.WriteRune(r)
			continue
		}
		buf.WriteString(spread.Colors[i].Ansi(mode, foreground))
		buf.WriteRune(r)
	}
	buf.WriteString("\033[39;49m")
	return buf.String(), nil
}

func (p ColorProfile) Lighten(multiplier float64) ColorProfile {
	out := make([]RGB, len(p.Colors))
	for i, c := range p.Colors {
		out[i] = c.lighten(multiplier)
	}
	return ColorProfile{Colors: out}
}

func (p ColorProfile) WithLightness(kind int, light float64) ColorProfile {
	out := make([]RGB, len(p.Colors))
	for i, c := range p.Colors {
		out[i] = c.assignLightness(kind, light)
	}
	return ColorProfile{Colors: out}
}

func (p ColorProfile) WithLightnessAdaptive(light float64, theme Theme) ColorProfile {
	if theme == ThemeLight {
		return p.WithLightness(assignClampMax, light)
	}
	return p.WithLightness(assignClampMin, light)
}

func (p ColorProfile) Unique() ColorProfile {
	seen := map[RGB]bool{}
	var out []RGB
	for _, c := range p.Colors {
		if seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return ColorProfile{Colors: out}
}
