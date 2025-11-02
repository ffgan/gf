package info

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ffgan/gf/internal/logo"
	"github.com/ffgan/gf/internal/utils"
)

func PrintInfo(asciiart, infoLines []string) {
	// TODO: 修复配色问题
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
		right := ""
		if i < len(infoLines) {
			right = infoLines[i]
		}
		// TODO: 处理转义字符
		visLen := visibleLen(left)
		padding := leftMax - visLen
		if padding < 0 {
			padding = 0
		}
		fmt.Printf("%s%s  %s\n", left, strings.Repeat(" ", padding), right)
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
