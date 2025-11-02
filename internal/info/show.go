package info

import (
	"fmt"
	"regexp"

	"github.com/ffgan/gf/internal/logo"
	"github.com/ffgan/gf/internal/utils"
)

func PrintInfo(asciiart, infoLines []string) {
	lines := utils.MaxLen(asciiart, len(infoLines))

	var leftMax int
	for _, left := range asciiart {
		visLen := visibleLen(left)
		if visLen > leftMax {
			leftMax = visLen
		}
	}

	for i := 0; i < lines; i++ {
		left := ""
		if i < len(asciiart) {
			left = asciiart[i]
		} else {
			left = ""
		}
		fmt.Printf("%s\n", left)
	}

	// move cursor to top
	fmt.Printf("\033[%dA", lines)
	leftMax += 2

	for i := 0; i < lines; i++ {
		fmt.Printf("\033[%dC", leftMax)
		right := ""
		if i < len(infoLines) {
			right = infoLines[i]
		}
		// TODO: 处理转义字符
		fmt.Printf("%s\n", right)
	}

	fmt.Println(logo.Reset)
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func visibleLen(s string) int {
	return len([]rune(stripAnsi(s)))
}
