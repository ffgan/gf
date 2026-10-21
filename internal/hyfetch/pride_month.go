package hyfetch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const prideASCII = "" +
	".======================================================.\n" +
	"| .  .              .__       .     .  .       , .   | |\n" +
	"| |__| _.._ ._   .  [__)._.* _| _   |\\/| _ ._ -+-|_  | |\n" +
	"| |  |(_][_)[_)\\_|  |   [  |(_](/,  |  |(_)[ ) | [ ) * |\n" +
	"|        |  |  ._|                                     |\n" +
	"'======================================================'\n"

const prideASCIISmall = "" +
	".====================.\n" +
	"| Happy Pride Month! |\n" +
	"'===================='\n"

// ShouldPride matches hyfetch: --june always, and once each June when both
// standard streams are terminals and the animation has not been shown.
func ShouldPride(cfg Config, force bool, now time.Time) bool {
	if cfg.PrideMonthDisable {
		return false
	}
	if force {
		return true
	}
	if now.Month() != time.June {
		return false
	}
	if !isTerminalFile(os.Stdin) || !isTerminalFile(os.Stdout) {
		return false
	}
	path, err := cachePath(fmt.Sprintf("animation-displayed-%d", now.Year()))
	if err != nil {
		return true
	}
	if _, err := os.Stat(path); err == nil {
		return false
	}
	return true
}

// MarkPrideShown records that the automatic June animation already played.
func MarkPrideShown(now time.Time) {
	path, err := cachePath(fmt.Sprintf("animation-displayed-%d", now.Year()))
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(path, []byte("1\n"), 0o644)
}

// Pride prints a short flag animation and the pride-month note.
func Pride(mode AnsiMode) {
	text := prideASCII
	if !isTerminalFile(os.Stdout) {
		text = prideASCIISmall
	}
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	names := []string{"rainbow", "transgender", "bisexual", "pansexual", "lesbian", "nonbinary", "asexual", "progress"}
	frames := len(names)
	if !isTerminalFile(os.Stdout) {
		frames = 1
	}
	for frame := 0; frame < frames; frame++ {
		if frame > 0 {
			fmt.Printf("\033[%dA", len(lines))
		}
		profile := Presets[names[frame%len(names)]]
		spread, err := profile.WithLength(len(lines))
		if err != nil {
			fmt.Println(text)
			break
		}
		for i, line := range lines {
			fmt.Printf("%s%s\033[0m\n", spread.Colors[i].Ansi(mode, true), line)
		}
		if frame+1 < frames {
			time.Sleep(90 * time.Millisecond)
		}
	}
	fmt.Println()
	fmt.Println("Happy pride month!")
	fmt.Println("(You can always view the animation again with `gf --june`)")
	fmt.Println()
}
