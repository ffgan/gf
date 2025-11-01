package info

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ffgan/gf/configs"
	cli "github.com/ffgan/gf/internal/CLI"
	gui "github.com/ffgan/gf/internal/GUI"
	dev "github.com/ffgan/gf/internal/devices"
	"github.com/ffgan/gf/internal/logo"
	"github.com/ffgan/gf/internal/utils"
)

func Setup(config *configs.Config) {
	_, config.ASCIIDistro = cli.GetDistro(config.OSArch, config.DistroShorthand, config.ASCIIDistro)
}

func PrintInfo(config *configs.Config) {
	asciiart := logo.PrintASCII(config.ASCIIDistro, config.ImageSource)

	infoLines := []string{
		cli.PrintDistro(config.OSArch, config.DistroShorthand, config.ASCIIDistro),
		cli.PrintHost(),
		cli.PrintKernel(config.OSArch, config.DistroShorthand, config.KernelShorthand, config.ASCIIDistro),
		cli.GetUptime(config.UptimeShorthand),
		cli.DetectPackages(),
		cli.PrintShell(config.ShellPath, config.ShellVersion),
		cli.Geteditor(),
		dev.DetectResolution(),
		gui.GetDE(),
		gui.GetWM(),
		gui.GetTheme(),
		gui.GetIcons(),
		gui.GetCursor(),
		cli.GetTerm(),
		cli.GetCPU(),
		dev.DetectGPU(),
		cli.GetMemory(),
		dev.DetectNetwork(),
		dev.DetectBluetooth(),
		dev.GetBIOS(),
	}

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
