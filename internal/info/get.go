package info

import (
	"embed"
	"strings"

	"github.com/ffgan/gf/configs"
	cli "github.com/ffgan/gf/internal/CLI"
	gui "github.com/ffgan/gf/internal/GUI"
	dev "github.com/ffgan/gf/internal/devices"
	"github.com/ffgan/gf/internal/logo"
)

func GetInfoLines(config *configs.Config) []string {
	return []string{
		cli.GetTitle(),
		strings.Repeat("-", 10),
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
		cli.PrintMem(config.Memory, config.ProgressBar),
		dev.DetectNetwork(),
		dev.DetectBluetooth(),
		dev.GetBIOS(),
	}
}

func GetInfo(config *configs.Config, ASCIIFiles embed.FS) ([]string, []string) {
	_, config.ASCIIDistro = cli.GetDistro(config.OSArch, config.DistroShorthand, config.ASCIIDistro)

	asciiart := logo.PrintASCII(config.ASCIIDistro, config.ImageSource, ASCIIFiles)

	infoLines := GetInfoLines(config)

	return asciiart, infoLines
}
