package cmd

import (
	"embed"

	"github.com/ffgan/gf/configs"
	cli "github.com/ffgan/gf/internal/CLI"
	"github.com/ffgan/gf/internal/info"
)

const neofetch_conf = "/etc/neofetch/default.conf"

func Run(ASCIIFiles embed.FS) {
	kernel_name, kernel_version, kernel_machine, darwin_name := cache_uname()

	os := cli.GetOS(kernel_name, darwin_name)

	cache_dir := get_cache_dir()

	var file_path string
	if cli.FileExists(neofetch_conf) == true {
		// load Neofetch default config.
		file_path = neofetch_conf
	}

	config, err := configs.LoadConfig(file_path)
	if err != nil {
		panic(err)
	}

	parseArgs(config)

	h := hardware{
		cache_dir:      cache_dir,
		KernelName:     kernel_name,
		kernel_version: kernel_version,
		KernelMachine:  kernel_machine,
		os:             os,
		Config:         *config,
	}

	// getsimple()

	// distro, ascii_distro := get_distro(os, config.OSArch, kernel_machine, config.DistroShorthand, config.ASCIIDistro)
	h.Get_distro()
	// ascii_bold, bold := get_bold(config.ASCIIBold, config.Bold)
	h.Get_bold()

	asciiData := Get_distro_ascii(h.ascii_distro, config.ImageSource, ASCIIFiles)
	h.leftMax = info.PrintInfo(asciiData)

	image_backend()
	old_functions()

	ctx := NewContext(h)

	for _, line := range config.InfoList {
		var prefix, func_name string
		if len(line) == 1 {
			func_name = line[0]
		} else if len(line) == 2 {
			prefix = line[0]
			func_name = line[1]
		}
		ctx.info(prefix, func_name)
	}

	dynamic_prompt()
}

func Info(name, data string) string {
	if data == "" {
		return ""
	}
	name = "\033[1;34m" + name + ": \033[0m"
	return name + data
}
