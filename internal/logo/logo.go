package logo

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	sysLocale       = "C.UTF-8"
	asciiData       string
	asciiLen        int
	asciiLines      int
	c1, c2, c3      = "\033[31m", "\033[32m", "\033[33m"
	c4, c5, c6      = "\033[34m", "\033[35m", "\033[36m"
	gap             = 2
	Reset           = "\033[0m" // reset terminal color
	asciiBold       = "\033[1m"
	colorText       = "on"
	barColorElapsed = "distro"
	barColorTotal   = "distro"
	colors          = []string{"distro"}
)

func getDistroASCII(ASCIIDistro string, ASCIIFiles embed.FS) {
	filename := "assets/"

	_, v, ok := r.LongestPrefix(ASCIIDistro)
	if !ok {
		panic("not found distro")
	}
	dc := v.(DistroWithColor)
	filename += dc.FileName
	setColors(dc.Colors...)

	f, err := ASCIIFiles.Open(filename)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	asciiData = string(data)
}

func stripEscapeCodes(s string) string {
	re := regexp.MustCompile(`\x1B\[[0-9;]*[A-Za-z]`)
	return re.ReplaceAllString(s, "")
}

func PrintASCII(ASCIIDistro, imageSource string, ASCIIFiles embed.FS) (logo []string) {
	fi, err := os.Stat(imageSource)
	if err == nil && fi.Mode().IsRegular() {
		ext := strings.ToLower(filepath.Ext(imageSource))
		switch ext {
		case ".png", ".jpg", ".jpeg", ".jpe", ".svg", ".gif":
		default:
			data, err := os.ReadFile(imageSource)
			if err == nil {
				asciiData = string(data)
			}
		}
	} else if imageSource == "ascii" || imageSource == "auto" {
		getDistroASCII(ASCIIDistro, ASCIIFiles)
	} else {
		asciiData = imageSource
	}

	// 设置 locale（Go 无法全局修改 LC_ALL，这里仅做标识）
	_ = sysLocale

	lines := 0

	if asciiLen > 0 && asciiLines > 0 {
		lines = asciiLines
	} else {
		scanner := bufio.NewScanner(strings.NewReader(asciiData))
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.ReplaceAll(line, "█", " ")

			line = stripEscapeCodes(line)
			line = strings.ReplaceAll(line, "${??}", "")
			line = strings.ReplaceAll(line, "\\033", "\033")

			lines++
		}
	}

	// Fallback
	if lines == 1 {
		lines = 0
		imageSource = "auto"
		getDistroASCII(ASCIIDistro, ASCIIFiles)
		PrintASCII(ASCIIDistro, imageSource, ASCIIFiles)
		return nil
	}

	replacements := map[string]string{
		"${c1}": c1, "${c2}": c2, "${c3}": c3,
		"${c4}": c4, "${c5}": c5, "${c6}": c6,
	}
	for k, v := range replacements {
		asciiData = strings.ReplaceAll(asciiData, k, v)
	}

	sysLocale = "C"
	return strings.Split(asciiData, "\n")
}

func color(input string) string {
	switch {
	case len(input) == 1 && input[0] >= '0' && input[0] <= '7':
		// 8色基础颜色
		return fmt.Sprintf("%s\033[3%sm", Reset, input)
	case input == "fg":
		// 恢复前景色
		return Reset
	case strings.HasPrefix(input, "#"):
		// RGB颜色（#RRGGBB）
		rgb, err := strconv.ParseInt(strings.TrimPrefix(input, "#"), 16, 32)
		if err != nil {
			return Reset
		}
		r := (rgb >> 16) & 0xff
		g := (rgb >> 8) & 0xff
		b := rgb & 0xff
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	default:
		// 256色模式
		return fmt.Sprintf("\033[38;5;%sm", input)
	}
}

func setColors(args_ints ...int) {
	var args []string
	for _, arg := range args_ints {
		args = append(args, string(arg))
	}
	c := make([]string, 6)
	for i := range c {
		if i < len(args) {
			c[i] = color(args[i]) + asciiBold
		} else {
			c[i] = ""
		}
	}

	if colorText != "off" {
		setTextColors(args...)
	}
}

var (
	titleColor, atColor, underlineColor  string
	subtitleColor, colonColor, infoColor string
)

func setTextColors(args ...string) {
	if len(colors) > 0 && colors[0] == "distro" {
		titleColor = color(getArg(args, 0))
		atColor = Reset
		underlineColor = Reset
		subtitleColor = color(getArg(args, 1))
		colonColor = Reset
		infoColor = Reset

		// Bash 中 ((x == 8)) 检查
		if getArgInt(args, 0, 1) == 8 {
			titleColor = Reset
		}
		if getArgInt(args, 1, 7) == 8 {
			subtitleColor = Reset
		}
		if getArgInt(args, 1, 7) == 7 {
			subtitleColor = color(getArg(args, 0))
		}
		if getArgInt(args, 0, 1) == 7 {
			titleColor = Reset
		}
	} else {
		titleColor = color(getColor(0))
		atColor = color(getColor(1))
		underlineColor = color(getColor(2))
		subtitleColor = color(getColor(3))
		colonColor = color(getColor(4))
		infoColor = color(getColor(5))
	}

	// Bar colors
	if barColorElapsed == "distro" {
		barColorElapsed = color("fg")
	} else {
		barColorElapsed = color(barColorElapsed)
	}

	if barColorTotal == "distro" {
		barColorTotal = color("fg")
	} else {
		barColorTotal = color(barColorTotal)
	}

	// fmt.Println("titleColor:", titleColor)
	// fmt.Println("subtitleColor:", subtitleColor)
	// fmt.Println("barColorElapsed:", barColorElapsed)
	// fmt.Println("barColorTotal:", barColorTotal)
}

func getArg(args []string, index int) string {
	if index < len(args) {
		return args[index]
	}
	return ""
}

func getArgInt(args []string, index int, def int) int {
	if index >= len(args) {
		return def
	}
	n, err := strconv.Atoi(args[index])
	if err != nil {
		return def
	}
	return n
}

func getColor(i int) string {
	if i < len(colors) {
		return colors[i]
	}
	return "7"
}
