package info

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ffgan/gf/internal/logo"
)

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
