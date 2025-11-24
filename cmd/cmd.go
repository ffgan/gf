package cmd

import (
	"embed"

	"github.com/ffgan/gf/configs"
	cli "github.com/ffgan/gf/internal/CLI"
)

func Run(ASCIIFiles embed.FS) {
	kernel_name, kernel_version, kernel_machine, darwin_name := cache_uname()

	os := cli.GetOS(kernel_name, darwin_name)

	cache_dir := get_cache_dir()

	// load Neofetch default config.
	config, err := configs.LoadConfig("")
	if err != nil {
		panic(err)
	}

	// TODO: load neofetch config
	// 结合上面内容，先加载一遍默认的，再尝试加载下面的配置文件
	// # If /etc/neofetch/default.conf exist, set config variable to its content
	// if [[ -f /etc/neofetch/default.conf ]]; then
	//     config="$(< /etc/neofetch/default.conf)"
	// fi

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
	_ = asciiData

	image_backend()
	old_functions()

	ctx := NewContext(h)

	ctx.info("OS", "os")
	ctx.info("Kernel", "kernel")
	ctx.info("BIOS", "bios")

	dynamic_prompt()
}
