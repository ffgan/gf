package hyfetch

import (
	"math/rand"
	"regexp"
	"strings"
)

var neofetchColor = regexp.MustCompile(`\$\{c([1-6])\}`)

// Normalize pads every line to the same visible width, ignoring ${cN} placeholders.
func Normalize(asc string) (lines []string, width int) {
	raw := strings.Split(strings.TrimRight(asc, "\n"), "\n")
	for _, line := range raw {
		if n := visiblePlaceholders(line); n > width {
			width = n
		}
	}
	lines = make([]string, len(raw))
	for i, line := range raw {
		pad := width - visiblePlaceholders(line)
		if pad < 0 {
			pad = 0
		}
		lines[i] = line + strings.Repeat(" ", pad)
	}
	return lines, width
}

func visiblePlaceholders(s string) int {
	return len([]rune(neofetchColor.ReplaceAllString(s, "")))
}

// FillStarting copies the previous color placeholder onto lines that do not
// open with one, matching hyfetch's ascii fill.
func FillStarting(lines []string) []string {
	last := "${c1}"
	out := make([]string, len(lines))
	for i, line := range lines {
		locs := neofetchColor.FindAllStringIndex(line, -1)
		starts := false
		if len(locs) > 0 {
			head := strings.TrimRight(line[:locs[0][0]], " ")
			starts = head == ""
		}
		if starts {
			out[i] = line
		} else {
			out[i] = last + line
		}
		if len(locs) > 0 {
			// Recompute on the original line; the last placeholder still applies.
			orig := neofetchColor.FindAllString(line, -1)
			if len(orig) > 0 {
				last = orig[len(orig)-1]
			}
		}
	}
	return out
}

func colorSlots(lines []string) []int {
	seen := map[int]bool{}
	var slots []int
	for _, line := range lines {
		for _, m := range neofetchColor.FindAllStringSubmatch(line, -1) {
			n := int(m[1][0] - '0')
			if !seen[n] {
				seen[n] = true
				slots = append(slots, n)
			}
		}
	}
	// Keep ascending order so random assignment is stable aside from the shuffle.
	for i := 1; i < len(slots); i++ {
		for j := i; j > 0 && slots[j] < slots[j-1]; j-- {
			slots[j], slots[j-1] = slots[j-1], slots[j]
		}
	}
	return slots
}

// RandomAlignment assigns preset colors to the ${cN} slots in the art.
func RandomAlignment(profile ColorProfile, lines []string) ColorAlignment {
	slots := colorSlots(lines)
	colors := profile.Unique().Colors
	align := ColorAlignment{Mode: "custom", CustomColors: map[int]int{}}
	if len(slots) == 0 || len(colors) == 0 {
		return align
	}
	order := make([]int, 0, len(slots))
	for len(order) < len(slots) {
		for i := range colors {
			order = append(order, i)
			if len(order) == len(slots) {
				break
			}
		}
	}
	rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
	for i, slot := range slots {
		align.CustomColors[slot] = order[i]
	}
	return align
}

// Recolor paints distro ASCII with a flag palette.
func Recolor(asc string, align ColorAlignment, profile ColorProfile, mode AnsiMode, theme Theme) (string, error) {
	if strings.TrimSpace(asc) == "" {
		return "", nil
	}
	switch align.Mode {
	case "", "horizontal":
		align.Mode = "horizontal"
	case "random":
		lines, _ := Normalize(asc)
		return Recolor(asc, RandomAlignment(profile, lines), profile, mode, theme)
	}
	lines, _ := Normalize(asc)
	lines = FillStarting(lines)
	reset := Color("&~&*", mode)
	switch align.Mode {
	case "horizontal":
		return recolorHorizontal(lines, align, profile, mode, theme, reset)
	case "vertical":
		return recolorVertical(lines, align, profile, mode, theme, reset)
	case "custom":
		return recolorCustom(lines, align, profile, mode, reset), nil
	default:
		return recolorHorizontal(lines, align, profile, mode, theme, reset)
	}
}

func recolorHorizontal(lines []string, align ColorAlignment, profile ColorProfile, mode AnsiMode, theme Theme, reset string) (string, error) {
	joined := strings.Join(lines, "\n")
	if len(align.ForeBack) >= 1 {
		fore := themeForeground(theme)
		joined = strings.ReplaceAll(joined, placeholder(align.ForeBack[0]), fore)
	}
	lines = strings.Split(joined, "\n")
	spread, err := profile.WithLength(len(lines))
	if err != nil {
		return "", err
	}
	var out []string
	for i, line := range lines {
		color := spread.Colors[i].Ansi(mode, true)
		if len(align.ForeBack) >= 2 && align.ForeBack[1] != 0 {
			line = strings.ReplaceAll(line, placeholder(align.ForeBack[1]), color)
			line = neofetchColor.ReplaceAllString(line, "")
		} else {
			line = neofetchColor.ReplaceAllString(line, color)
		}
		out = append(out, line+reset)
	}
	return strings.Join(out, "\n"), nil
}

func recolorVertical(lines []string, align ColorAlignment, profile ColorProfile, mode AnsiMode, theme Theme, reset string) (string, error) {
	_ = align
	_ = theme
	plain := neofetchColor.ReplaceAllString(strings.Join(lines, "\n"), "")
	var out []string
	for _, line := range strings.Split(plain, "\n") {
		if strings.TrimSpace(line) == "" && line == "" {
			out = append(out, "")
			continue
		}
		colored, err := profile.ColorText(line, mode, true, false)
		if err != nil {
			return "", err
		}
		out = append(out, colored+reset)
	}
	return strings.Join(out, "\n"), nil
}

func recolorCustom(lines []string, align ColorAlignment, profile ColorProfile, mode AnsiMode, reset string) string {
	colors := profile.Unique().Colors
	if len(colors) == 0 {
		return strings.Join(lines, "\n")
	}
	slots := colorSlots(lines)
	replace := map[int]string{}
	for _, slot := range slots {
		pi, ok := align.CustomColors[slot]
		if !ok {
			if alt, found := align.CustomColors[(slot-1)%2+1]; found {
				pi = alt
			} else {
				pi = slot - 1
			}
		}
		if pi < 0 {
			pi = 0
		}
		replace[slot] = colors[pi%len(colors)].Ansi(mode, true)
	}
	joined := strings.Join(lines, "\n")
	joined = neofetchColor.ReplaceAllStringFunc(joined, func(token string) string {
		m := neofetchColor.FindStringSubmatch(token)
		if m == nil {
			return ""
		}
		return replace[int(m[1][0]-'0')]
	})
	var out []string
	for _, line := range strings.Split(joined, "\n") {
		out = append(out, line+reset)
	}
	return strings.Join(out, "\n")
}

func placeholder(n int) string {
	if n < 1 || n > 6 {
		return ""
	}
	return "${c" + string(rune('0'+n)) + "}"
}

func themeForeground(theme Theme) string {
	if theme == ThemeLight {
		return Color("&0", Mode8bit)
	}
	return Color("&f", Mode8bit)
}
