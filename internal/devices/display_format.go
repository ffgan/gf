package dev

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type winDisplay struct {
	Name    string
	W, H    int
	Hz      float64
	DPI     int
	Wmm     int
	Hmm     int
	Kind    string
	Primary bool
	ShowHz  bool
}

func formatWinDisplay(d winDisplay) string {
	if d.W <= 0 || d.H <= 0 {
		return ""
	}
	text := fmt.Sprintf("%dx%d", d.W, d.H)
	if scale := displayScale(d.DPI); scale != "" {
		text += " @ " + scale + "x"
	}
	if inches := diagonalInches(d.Wmm, d.Hmm); inches > 0 {
		text += fmt.Sprintf(" in %d\"", inches)
	}
	if d.ShowHz && d.Hz > 0 {
		if math.Abs(d.Hz-math.Round(d.Hz)) < 0.05 {
			text += fmt.Sprintf(", %.0f Hz", d.Hz)
		} else {
			text += fmt.Sprintf(", %.2f Hz", d.Hz)
		}
	}
	if d.Kind != "" {
		text += " [" + d.Kind + "]"
	}
	if d.Primary {
		text += " *"
	}
	if name := strings.TrimSpace(d.Name); name != "" {
		text = name + ": " + text
	}
	return text
}

func displayScale(dpi int) string {
	if dpi <= 0 || dpi == 96 {
		return ""
	}
	v := math.Round(float64(dpi)/96*100) / 100
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func diagonalInches(wmm, hmm int) int {
	if wmm <= 0 || hmm <= 0 {
		return 0
	}
	return int(math.Round(math.Hypot(float64(wmm), float64(hmm)) / 25.4))
}

func edidNameSize(edid []byte) (name string, wmm, hmm int) {
	if len(edid) < 128 || edid[0] != 0x00 || edid[1] != 0xFF {
		return "", 0, 0
	}
	wmm = int(edid[21]) * 10
	hmm = int(edid[22]) * 10
	for i := 0; i < 4; i++ {
		off := 54 + i*18
		if off+18 > len(edid) {
			break
		}
		if edid[off] == 0 && edid[off+1] == 0 && edid[off+2] == 0 && edid[off+3] == 0xFC {
			raw := strings.Trim(string(edid[off+5:off+18]), " \x00\n")
			if raw != "" {
				name = raw
			}
			break
		}
	}
	return name, wmm, hmm
}

func edidKeyFromDevicePath(path string) string {
	path = strings.TrimPrefix(path, `\\?\`)
	path = strings.TrimPrefix(path, `//?/`)
	if i := strings.Index(path, "{"); i >= 0 {
		path = path[:i]
	}
	path = strings.TrimRight(path, `#\`)
	path = strings.ReplaceAll(path, "#", `\`)
	path = strings.Trim(path, `\`)
	if path == "" {
		return ""
	}
	return `SYSTEM\CurrentControlSet\Enum\` + path + `\Device Parameters`
}
