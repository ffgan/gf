package info

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ffgan/gf/internal/logo"
)

// Show prints the logo on the left and info rows on the right.
func Show(art, rows []string, gap int) {
	if gap < 0 {
		gap = 0
	}
	width := 0
	for _, line := range art {
		if n := visibleLen(line); n > width {
			width = n
		}
	}
	n := len(rows)
	if len(art) > n {
		n = len(art)
	}
	if n == 0 {
		return
	}
	gapText := strings.Repeat(" ", gap)
	for i := 0; i < n; i++ {
		left := ""
		if i < len(art) {
			left = art[i]
		}
		right := ""
		if i < len(rows) {
			right = rows[i]
		}
		if len(art) == 0 {
			fmt.Println(right)
			continue
		}
		pad := width - visibleLen(left)
		if pad < 0 {
			pad = 0
		}
		fmt.Printf("%s%s%s%s\n", left, "\033[0m", strings.Repeat(" ", pad)+gapText, right)
	}
}

func PrintInfo(asciiart []string) int {
	lines := len(asciiart)

	var leftMax int
	for _, left := range asciiart {
		visLen := visibleLen(left)
		if visLen > leftMax {
			leftMax = visLen
		}
	}

	fmt.Println(strings.Join(asciiart, "\n"))

	// move cursor to top
	fmt.Printf("\033[%dA", lines+1)

	fmt.Println(logo.Reset)

	return leftMax + 2

	// for i := 0; i < lines; i++ {
	// 	fmt.Printf("\033[%dC", leftMax)
	// 	right := ""
	// 	if i < len(infoLines) {
	// 		right = infoLines[i]
	// 	}
	// 	// TODO: 处理转义字符
	// 	fmt.Printf("%s\n", right)
	// }

}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func visibleLen(s string) int {
	return len([]rune(stripAnsi(s)))
}

func Info(name, data string) string {
	if data == "" {
		return ""
	}
	name = "\033[1;34m" + name + ": \033[0m"
	return name + data
}
