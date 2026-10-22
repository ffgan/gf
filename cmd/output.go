package cmd

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/ffgan/gf/internal/info"
	"github.com/ffgan/gf/internal/logo"
	"github.com/ffgan/gf/internal/modules"
	"github.com/ffgan/gf/internal/utils"
)

func (h *hardware) plainMode() bool {
	return h.plain || h.Stdout == utils.ON
}

type infoEntry struct {
	label  string
	key    string
	value  string
	blank  bool
	blocks bool
}

func (h *hardware) Rows() []string {
	var rows []string
	for _, entry := range h.entries() {
		switch {
		case entry.blank:
			rows = append(rows, "")
		case entry.blocks:
			for _, block := range strings.Split(entry.value, "\n") {
				if block != "" {
					rows = append(rows, block)
				}
			}
		case entry.label == "":
			rows = append(rows, entry.value)
		default:
			for _, part := range strings.Split(entry.value, "\n") {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				rows = append(rows, h.formatPair(entry.label, part))
			}
		}
	}
	return rows
}

func (h *hardware) entries() []infoEntry {
	h.prefetch()
	var out []infoEntry
	for _, line := range h.InfoList {
		if h.skipped(line) {
			continue
		}
		if label, text, ok := prinLine(line); ok {
			if label == "" {
				out = append(out, infoEntry{value: text})
			} else {
				out = append(out, infoEntry{label: label, value: text})
			}
			continue
		}
		label, key := splitInfo(line)
		key = modules.Canonical(key)
		if key == "" || key == "line_break" || key == "blank" {
			out = append(out, infoEntry{blank: true})
			continue
		}
		if key == "cols" {
			out = append(out, infoEntry{key: key, value: h.cols, blocks: true})
			continue
		}
		value := strings.TrimSpace(h.field(key))
		if value == "" {
			if h.verbose {
				fmt.Fprintf(os.Stderr, "gf: couldn't detect %s\n", key)
			}
			continue
		}
		out = append(out, infoEntry{label: label, key: key, value: value})
	}
	return out
}

func prinLine(line []string) (label, text string, ok bool) {
	n := len(line)
	if n >= 2 && line[n-2] == "prin" {
		if n >= 3 {
			return line[0], line[n-1], true
		}
		return "", line[n-1], true
	}
	return "", "", false
}

func (h *hardware) skipped(line []string) bool {
	if len(h.disabled) == 0 {
		return false
	}
	for _, part := range line {
		if h.disabled[strings.ToLower(part)] {
			return true
		}
	}
	_, key := splitInfo(line)
	key = modules.Canonical(key)
	return h.disabled[strings.ToLower(key)] || h.disabled[key]
}

func splitInfo(line []string) (label, key string) {
	switch len(line) {
	case 0:
		return "", ""
	case 1:
		return "", line[0]
	default:
		return line[0], line[len(line)-1]
	}
}

// prefetch detects every info line concurrently. fastfetch's modules do
// not wait on each other, and the slow ones overlap. A second wave waits
// for values the first wave writes.
func (h *hardware) prefetch() {
	seen := map[string]bool{}
	var first, second []string
	for _, line := range h.InfoList {
		if h.skipped(line) {
			continue
		}
		if _, _, ok := prinLine(line); ok {
			continue
		}
		_, key := splitInfo(line)
		key = modules.Canonical(key)
		if key == "" || key == "line_break" || key == "blank" || key == "distro" || seen[key] {
			continue
		}
		seen[key] = true
		switch key {
		case "wm_theme", "term_font", "underline", "separator":
			second = append(second, key)
		default:
			first = append(first, key)
		}
	}
	run := func(keys []string) {
		var wg sync.WaitGroup
		for _, key := range keys {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				h.invoke(key)
			}(key)
		}
		wg.Wait()
	}
	run(first)
	run(second)
}

func (h *hardware) invoke(key string) {
	key = modules.Canonical(key)
	if detect, ok := fastfetchDetectors[key]; ok {
		h.remember(key, detect(h))
		return
	}
	method := reflect.ValueOf(h).MethodByName("Get_" + key)
	if !method.IsValid() {
		if h.verbose {
			fmt.Fprintf(os.Stderr, "gf: no info function %q\n", key)
		}
		return
	}
	method.Call(nil)
}

func (h *hardware) field(name string) string {
	name = modules.Canonical(name)
	if value, ok := h.recall(name); ok {
		return value
	}
	switch name {
	case "distro":
		return h.distro
	case "model":
		return h.model
	case "kernel":
		return h.kernel
	case "uptime":
		return h.uptime
	case "packages":
		return h.packages
	case "shell":
		return h.shell
	case "editor":
		return h.editor
	case "de":
		return h.de
	case "wm":
		return h.wm
	case "wm_theme":
		return h.wm_theme
	case "cpu":
		return h.cpu
	case "gpu":
		return h.gpu
	case "npu":
		return h.npu
	case "memory":
		return h.memory
	case "swap":
		return h.swap
	case "network":
		return h.network
	case "bluetooth":
		return h.bluetooth
	case "resolution":
		return h.resolution
	case "term":
		return h.term
	case "term_font":
		return h.term_font
	case "font":
		return h.font
	case "theme":
		return h.theme
	case "icons":
		return h.icons
	case "cursor":
		return h.cursor
	case "cols":
		return h.cols
	case "title":
		return h.title
	case "underline":
		return h.underline
	case "bios":
		return h.bios
	case "disk":
		return h.disk
	case "battery":
		return h.battery
	case "local_ip":
		return h.local_ip
	case "public_ip":
		return h.public_ip
	case "users":
		return h.users
	case "locale":
		return h.locale
	case "gpu_driver":
		return h.gpu_driver
	case "song":
		return h.song
	case "java_ver":
		return h.java_ver
	case "python_ver":
		return h.python_ver
	case "node_ver":
		return h.node_ver
	case "term_size":
		return h.term_size
	case "window_size":
		return h.window_size
	case "power_adapter":
		return h.power_adapter
	default:
		return ""
	}
}

func (h *hardware) formatPair(label, value string) string {
	sep := h.Separator
	if sep == "" {
		sep = ":"
	}
	if h.plainMode() {
		return label + sep + " " + value
	}
	bold := ""
	if h.Bold != utils.OFF {
		bold = "\033[1m"
	}
	// Color sequences start by resetting, so bold has to come after them.
	return logo.SubtitleColor() + bold + label + logo.ResetCode() + sep + " " + value
}

func (h *hardware) gap() int {
	n, err := strconv.Atoi(strings.TrimSpace(h.Gap))
	if err != nil {
		return 3
	}
	if n < 0 {
		return 0
	}
	return n
}

func show(art, rows []string, gap int) {
	info.Show(art, rows, gap)
}
