package modules

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// FormatBytes prints a size the way a fetch line does: binary units,
// with a fraction only when the value is not a whole number of units.
func FormatBytes(n uint64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	v := float64(n)
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	if i == 0 || v == math.Trunc(v) {
		return fmt.Sprintf("%.0f %s", v, units[i])
	}
	return fmt.Sprintf("%.2f %s", v, units[i])
}

func joinLines(lines []string) string {
	var out []string
	seen := map[string]bool{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || seen[line] {
			continue
		}
		seen[line] = true
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func firstLine(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	line, _, _ := strings.Cut(text, "\n")
	return strings.TrimSpace(line)
}

// perSecond prints a byte count as a rate. A non-positive span prints the
// count itself, which is the total rather than a per-second sample.
func perSecond(n uint64, seconds float64) string {
	if seconds <= 0 {
		return FormatBytes(n)
	}
	rate := uint64(float64(n) / seconds)
	return FormatBytes(rate) + "/s"
}

func formatDiskRate(name string, read, write uint64, seconds float64) string {
	return fmt.Sprintf("%s: %s (R) - %s (W)", name, perSecond(read, seconds), perSecond(write, seconds))
}

func formatNetRate(name string, rx, tx uint64, seconds float64) string {
	return fmt.Sprintf("%s: %s (rx) - %s (tx)", name, perSecond(rx, seconds), perSecond(tx, seconds))
}

func formatProcesses(n, threads uint64, haveThreads bool) string {
	if n == 0 {
		return ""
	}
	if !haveThreads {
		return fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("%d (%d threads)", n, threads)
}

type cacheCount struct {
	level int
	size  uint64
	kind  string
	num   int
}

func addCache(items []cacheCount, level int, size uint64, kind string) []cacheCount {
	for i := range items {
		if items[i].level == level && items[i].size == size && items[i].kind == kind {
			items[i].num++
			return items
		}
	}
	return append(items, cacheCount{level: level, size: size, kind: kind, num: 1})
}

func summarizeCache(items []cacheCount) string {
	if len(items) == 0 {
		return ""
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].level != items[j].level {
			return items[i].level < items[j].level
		}
		return items[i].kind < items[j].kind
	})
	var lines []string
	for level := 1; level <= 4; level++ {
		var parts []string
		for _, it := range items {
			if it.level != level {
				continue
			}
			size := FormatBytes(it.size)
			if it.num > 1 {
				size = fmt.Sprintf("%dx %s", it.num, size)
			}
			parts = append(parts, fmt.Sprintf("%s (%s)", size, it.kind))
		}
		if len(parts) > 0 {
			lines = append(lines, fmt.Sprintf("L%d: %s", level, strings.Join(parts, ", ")))
		}
	}
	return strings.Join(lines, "\n")
}
