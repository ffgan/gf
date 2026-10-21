package hyfetch

import (
	"fmt"
	"strconv"
	"strings"
)

// AnsiMode matches hyfetch's color mode. "8bit" is the 256-color palette and
// "rgb" is 24-bit truecolor.
type AnsiMode string

const (
	Mode8bit AnsiMode = "8bit"
	ModeRGB  AnsiMode = "rgb"
)

// Theme is the terminal background hyfetch adapts lightness to.
type Theme string

const (
	ThemeDark  Theme = "dark"
	ThemeLight Theme = "light"
)

// RGB is an 8-bit sRGB color.
type RGB struct {
	R, G, B uint8
}

func ParseHex(s string) (RGB, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "#")
	if len(s) == 3 {
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	}
	if len(s) != 6 {
		return RGB{}, fmt.Errorf("invalid hex color: %s", s)
	}
	for _, c := range s {
		if !isHex(c) {
			return RGB{}, fmt.Errorf("invalid hex color: %s", s)
		}
	}
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color: %s", s)
	}
	return RGB{uint8(n >> 16), uint8(n >> 8), uint8(n)}, nil
}

func isHex(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func (c RGB) Ansi(mode AnsiMode, foreground bool) string {
	ch := "38"
	if !foreground {
		ch = "48"
	}
	if mode == ModeRGB {
		return fmt.Sprintf("\033[%s;2;%d;%d;%dm", ch, c.R, c.G, c.B)
	}
	return fmt.Sprintf("\033[%s;5;%dm", ch, c.ansi256())
}

// ansi256 follows hyfetch's Python quantizer, which is the palette the
// published presets were tuned against.
func (c RGB) ansi256() int {
	r, g, b := float64(c.R), float64(c.G), float64(c.B)
	sep := 42.5
	gray := false
	for {
		if r < sep || g < sep || b < sep {
			gray = r < sep && g < sep && b < sep
			break
		}
		sep += 42.5
	}
	if gray {
		return 232 + int((r+g+b)/33)
	}
	return 16 + int(r/256*6)*36 + int(g/256*6)*6 + int(b/256*6)
}

const (
	assignReplace = iota
	assignClampMin
	assignClampMax
)

func (c RGB) assignLightness(kind int, light float64) RGB {
	h, s, l := srgbToOKHSL(c.R, c.G, c.B)
	switch kind {
	case assignReplace:
		l = light
	case assignClampMin:
		if l < light {
			l = light
		}
	case assignClampMax:
		if l > light {
			l = light
		}
	}
	r, g, b := okhslToSRGB(h, s, clamp01(l))
	return RGB{r, g, b}
}

func (c RGB) lighten(multiplier float64) RGB {
	_, _, l := srgbToOKHSL(c.R, c.G, c.B)
	return c.assignLightness(assignReplace, clamp01(l*multiplier))
}

func (c RGB) IsLight() bool {
	_, _, l := srgbToOKHSL(c.R, c.G, c.B)
	return l > 0.5
}

// DefaultLightness is hyfetch's per-theme OKHSL lightness.
func DefaultLightness(theme Theme) float64 {
	if theme == ThemeLight {
		return 0.4
	}
	return 0.65
}

// Color replaces hyfetch's extended Minecraft color codes.
func Color(msg string, mode AnsiMode) string {
	replacer := strings.NewReplacer(
		"&0", "\033[38;5;0m",
		"&1", "\033[38;5;4m",
		"&2", "\033[38;5;2m",
		"&3", "\033[38;5;6m",
		"&4", "\033[38;5;1m",
		"&5", "\033[38;5;5m",
		"&6", "\033[38;5;3m",
		"&7", "\033[38;5;7m",
		"&8", "\033[38;5;8m",
		"&9", "\033[38;5;12m",
		"&a", "\033[38;5;10m",
		"&b", "\033[38;5;14m",
		"&c", "\033[38;5;9m",
		"&d", "\033[38;5;13m",
		"&e", "\033[38;5;11m",
		"&f", "\033[38;5;15m",
		"&l", "\033[1m",
		"&o", "\033[3m",
		"&n", "\033[4m",
		"&k", "\033[8m",
		"&m", "\033[9m",
		"&r", "\033[0m",
		"&-", "\n",
		"&~", "\033[39m",
		"&*", "\033[49m",
		"&L", "\033[22m",
		"&O", "\033[23m",
		"&N", "\033[24m",
		"&K", "\033[28m",
		"&M", "\033[29m",
	)
	msg = replacer.Replace(msg)
	for _, prefix := range []string{"&gf(", "&gb("} {
		for {
			i := strings.Index(msg, prefix)
			if i < 0 {
				break
			}
			end := strings.Index(msg[i:], ")")
			if end < 0 {
				break
			}
			end += i
			code := msg[i+4 : end]
			fore := msg[i+2] == 'f'
			rgb, err := parseMinecraftRGB(code)
			if err != nil {
				break
			}
			msg = msg[:i] + rgb.Ansi(mode, fore) + msg[end+1:]
		}
	}
	return msg
}

func parseMinecraftRGB(code string) (RGB, error) {
	if strings.HasPrefix(code, "#") {
		return ParseHex(code)
	}
	code = strings.NewReplacer(",", " ", ";", " ").Replace(code)
	fields := strings.Fields(code)
	if len(fields) != 3 {
		return RGB{}, fmt.Errorf("rgb code %q", code)
	}
	var nums [3]uint8
	for i, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil || n < 0 || n > 255 {
			return RGB{}, fmt.Errorf("rgb code %q", code)
		}
		nums[i] = uint8(n)
	}
	return RGB{nums[0], nums[1], nums[2]}, nil
}

func ParseMode(s string) (AnsiMode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "8bit", "256", "ansi256":
		return Mode8bit, nil
	case "rgb", "truecolor", "24bit":
		return ModeRGB, nil
	default:
		return "", fmt.Errorf("MODE should be one of {8bit,rgb}")
	}
}

func ParseLightness(s string) (float64, error) {
	n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil || n < 0 || n > 1 {
		return 0, fmt.Errorf("invalid lightness %q, expected a value between 0 and 1", s)
	}
	return n, nil
}
