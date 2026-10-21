package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

// DiskReport describes how disk usage should be printed.
type DiskReport struct {
	Show      string
	Subtitle  string
	Percent   string
	Display   string
	BarWidth  int
	BarFilled string
	BarEmpty  string
}

// GetDisk returns one line per requested mount, for example "(/) 7.0G / 253G (3%)".
func GetDisk(report DiskReport) string {
	mounts := parseShowList(report.Show)
	if len(mounts) == 0 {
		mounts = []string{"/"}
	}
	var lines []string
	for _, mount := range mounts {
		used, total, percent, err := dfBytes(mount)
		if err != nil || total <= 0 {
			continue
		}
		text := fmt.Sprintf("%s / %s", humanBytes(used), humanBytes(total))
		if report.Percent != utils.OFF {
			text += fmt.Sprintf(" (%d%%)", percent)
		}
		if report.Display == "bar" || report.Display == "infobar" || report.Display == "barinfo" {
			bar := diskBar(used, total, report)
			switch report.Display {
			case "bar":
				text = bar
			case "barinfo":
				text = bar + " " + text
			default:
				text = text + " " + bar
			}
		}
		if label := diskLabel(mount, report.Subtitle); label != "" {
			text = label + " " + text
		}
		lines = append(lines, text)
	}
	return strings.Join(lines, "\n")
}

func parseShowList(show string) []string {
	show = strings.TrimSpace(show)
	show = strings.Trim(show, "()")
	if show == "" {
		return nil
	}
	fields := strings.Fields(show)
	var mounts []string
	for _, field := range fields {
		field = strings.Trim(field, `"'`)
		field = strings.TrimSpace(field)
		if field != "" {
			mounts = append(mounts, field)
		}
	}
	return mounts
}

func diskLabel(mount, subtitle string) string {
	switch subtitle {
	case "none", utils.OFF:
		return ""
	case "dir":
		if mount == "/" {
			return "(/)"
		}
		parts := strings.Split(strings.TrimRight(mount, "/"), "/")
		return "(" + parts[len(parts)-1] + ")"
	case "name":
		return "(" + mount + ")"
	default:
		return "(" + mount + ")"
	}
}

func dfBytes(path string) (int64, int64, int, error) {
	out, err := utils.CommandOutput("df", "-B1", "-P", path)
	if err != nil {
		return 0, 0, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return 0, 0, 0, fmt.Errorf("unexpected df output")
	}
	fields := strings.Fields(lines[len(lines)-1])
	if len(fields) < 4 {
		return 0, 0, 0, fmt.Errorf("unexpected df columns")
	}
	total, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	used, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	percent := 0
	if len(fields) >= 5 && strings.HasSuffix(fields[4], "%") {
		percent, _ = strconv.Atoi(strings.TrimSuffix(fields[4], "%"))
	} else if total > 0 {
		percent = int(used * 100 / total)
	}
	return used, total, percent, nil
}

func humanBytes(n int64) string {
	units := []string{"B", "K", "M", "G", "T", "P"}
	value := float64(n)
	unit := 0
	for value >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}
	if unit == 0 {
		return fmt.Sprintf("%d%s", n, units[0])
	}
	if value >= 100 || value == float64(int64(value)) {
		return fmt.Sprintf("%.0f%s", value, units[unit])
	}
	return fmt.Sprintf("%.1f%s", value, units[unit])
}

func diskBar(used, total int64, report DiskReport) string {
	width := report.BarWidth
	if width <= 0 {
		width = 15
	}
	filled := int(float64(used) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	on := report.BarFilled
	off := report.BarEmpty
	if on == "" {
		on = "="
	}
	if off == "" {
		off = "-"
	}
	return "[" + strings.Repeat(on, filled) + strings.Repeat(off, width-filled) + "]"
}
