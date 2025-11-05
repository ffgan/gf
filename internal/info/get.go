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
	// TODO: 有些行的信息无法读取的时候，应当不显示
	// TOOD: 读取配置文件，控制显示的行
	lines := []string{
		cli.GetTitle(),
		// TODO: 这里数字应该与title长度有关
		strings.Repeat("-", 10),
		Info("OS", cli.PrintDistro(config.OSArch, config.DistroShorthand, config.ASCIIDistro)),
		Info("Host", cli.PrintHost()),
		Info("Kernel", cli.PrintKernel(config.OSArch, config.DistroShorthand, config.KernelShorthand, config.ASCIIDistro)),
		Info("Uptime", cli.GetUptime(config.UptimeShorthand)),
		Info("Packages", cli.DetectPackages()),
		Info("Shell", cli.PrintShell(config.ShellPath, config.ShellVersion)),
		Info("Editor", cli.Geteditor()),
		Info("Resolution", dev.DetectResolution()),
		Info("DE", gui.GetDE()),
		Info("WM", gui.GetWM()),
		Info("Theme", gui.GetTheme()),
		Info("Icons", gui.GetIcons()),
		Info("Cursor", gui.GetCursor()),
		Info("Terminal", cli.GetTerm()),
		Info("CPU", cli.GetCPU(config.CPUCores, config.CPUSpeed, config.CPUTemp, config.CPUBrand, config.SpeedType, config.SpeedShorthand, config.CPUTemp)),
		Info("GPU", dev.DetectGPU()),
		Info("Memory", cli.PrintMem(config.Memory, config.ProgressBar)),
		Info("Network", dev.DetectNetwork()),
		Info("Bluetooth", dev.DetectBluetooth()),
		Info("BIOS", dev.GetBIOS()),
	}
	var res []string
	for _, line := range lines {
		if line == "" {
			continue
		}
		res = append(res, line)
	}
	return res
}

func Info(name, data string) string {
	if data == "" {
		return ""
	}
	name = "\033[1;34m" + name + ": \033[0m"
	return name + data
}

func GetInfo(config *configs.Config, ASCIIFiles embed.FS) ([]string, []string) {
	_, config.ASCIIDistro = cli.GetDistro(config.OSArch, config.DistroShorthand, config.ASCIIDistro)

	asciiart := logo.PrintASCII(config.ASCIIDistro, config.ImageSource, ASCIIFiles)

	infoLines := GetInfoLines(config)

	return asciiart, infoLines
}
