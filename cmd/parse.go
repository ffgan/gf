package cmd

import (
	"flag"

	"github.com/ffgan/gf/configs"
)

func setupArgs(config *configs.Config) {
	flag.StringVar(&config.TitleFqdn, "title_fqdn", config.TitleFqdn, "Set title_fqdn value")
	flag.StringVar(&config.PackageManagers, "package_managers", config.PackageManagers, "Set package managers")
	flag.StringVar(&config.PackageSeparate, "package_separate", config.PackageSeparate, "Set package separate")
	flag.StringVar(&config.PackageMinimal, "package_minimal", config.PackageMinimal, "Enable package minimal mode")

	flag.StringVar(&config.OSArch, "os_arch", config.OSArch, "Set OS architecture")
	flag.StringVar(&config.CPUCores, "cpu_cores", config.CPUCores, "Set CPU cores")
	flag.StringVar(&config.CPUSpeed, "cpu_speed", config.CPUSpeed, "Set CPU speed")
	flag.StringVar(&config.SpeedType, "speed_type", config.SpeedType, "Set CPU speed type")
	flag.StringVar(&config.SpeedShorthand, "speed_shorthand", config.SpeedShorthand, "Set speed shorthand")
	flag.StringVar(&config.DistroShorthand, "distro_shorthand", config.DistroShorthand, "Set distro shorthand")
	flag.StringVar(&config.KernelShorthand, "kernel_shorthand", config.KernelShorthand, "Set kernel shorthand")
	flag.StringVar(&config.UptimeShorthand, "uptime_shorthand", config.UptimeShorthand, "Set uptime shorthand")
	flag.StringVar(&config.CPUBrand, "cpu_brand", config.CPUBrand, "Set CPU brand")
	flag.StringVar(&config.GPUBrand, "gpu_brand", config.GPUBrand, "Set GPU brand")
	flag.StringVar(&config.GPUType, "gpu_type", config.GPUType, "Set GPU type")
	flag.StringVar(&config.RefreshRate, "refresh_rate", config.RefreshRate, "Set refresh rate")
	flag.StringVar(&config.DeVersion, "de_version", config.DeVersion, "Set DE version")

	flag.StringVar(&config.GTKShorthand, "gtk_shorthand", config.GTKShorthand, "GTK shorthand")
	flag.StringVar(&config.GTK2, "gtk2", config.GTK2, "GTK2 version")
	flag.StringVar(&config.GTK3, "gtk3", config.GTK3, "GTK3 version")
	flag.StringVar(&config.QT, "qt", config.QT, "QT version")

	flag.StringVar(&config.ShellPath, "shell_path", config.ShellPath, "Shell path")
	flag.StringVar(&config.ShellVersion, "shell_version", config.ShellVersion, "Shell version")
	flag.StringVar(&config.EditorPath, "editor_path", config.EditorPath, "Editor path")
	flag.StringVar(&config.EditorVersion, "editor_version", config.EditorVersion, "Editor version")

	flag.StringVar(&config.PublicIPHost, "ip_host", config.PublicIPHost, "Public IP host")
	flag.StringVar(&config.PublicIPTimeout, "ip_timeout", config.PublicIPTimeout, "IP timeout")
	flag.StringVar(&config.LocalIPInterface, "ip_interface", config.LocalIPInterface,
		"Local IP interface list (comma separated)")

	flag.StringVar(&config.MemoryPercent, "memory_percent", config.MemoryPercent, "Memory percent")
	flag.StringVar(&config.MemoryUnit, "memory_unit", config.MemoryUnit, "Memory unit")
	flag.StringVar(&config.MemPrecision, "memory_precision", config.MemPrecision, "Memory precision")
	flag.StringVar(&config.CPUTemp, "cpu_temp", config.CPUTemp, "CPU temperature (on or C)")

	flag.StringVar(&config.DiskSubtitle, "disk_subtitle", config.DiskSubtitle, "Disk subtitle")
	flag.StringVar(&config.DiskPercent, "disk_percent", config.DiskPercent, "Disk percent")
	flag.StringVar(&config.DiskShow, "disk_show", config.DiskShow, "Disk show targets (comma separated)")

	flag.StringVar(&config.Colors, "color_blocks", config.Colors, "Enable color blocks")
	flag.StringVar(&config.BlockRange, "block_range", config.BlockRange, "Block range (start,end)")
	flag.StringVar(&config.BlockWidth, "block_width", config.BlockWidth, "Block width")
	flag.StringVar(&config.BlockHeight, "block_height", config.BlockHeight, "Block height")
	flag.StringVar(&config.ColOffset, "col_offset", config.ColOffset, "Column offset")

	flag.StringVar(&config.BarCharElapsed, "bar_char", config.BarCharElapsed, "Bar character elapsed/total (e.g. '#' '-')")
	flag.StringVar(&config.BarBorder, "bar_border", config.BarBorder, "Bar border style")
	flag.StringVar(&config.BarLength, "bar_length", config.BarLength, "Bar length")
	flag.StringVar(&config.BarColorElapsed, "bar_colors", config.BarColorElapsed, "Bar colors elapsed/total")

	flag.StringVar(&config.MemoryDisplay, "memory_display", config.MemoryDisplay, "Memory display mode")
	flag.StringVar(&config.BatteryDisplay, "battery_display", config.BatteryDisplay, "Battery display mode")
	flag.StringVar(&config.DiskDisplay, "disk_display", config.DiskDisplay, "Disk display mode")

	flag.StringVar(&config.ImageBackend, "backend", config.ImageBackend, "Image backend type")
	flag.StringVar(&config.ImageSource, "source", config.ImageSource, "Image source path")
	flag.StringVar(&config.ImageLoop, "loop", config.ImageLoop, "Enable image loop")
	flag.StringVar(&config.ImageSize, "image_size", config.ImageSize, "Image size")
	flag.StringVar(&config.CatimgSize, "catimg_size", config.CatimgSize, "Catimg size")
	flag.StringVar(&config.CropMode, "crop_mode", config.CropMode, "Crop mode")
	flag.StringVar(&config.CropOffset, "crop_offset", config.CropOffset, "Crop offset")
	flag.StringVar(&config.Xoffset, "xoffset", config.Xoffset, "X offset")
	flag.StringVar(&config.Yoffset, "yoffset", config.Yoffset, "Y offset")
	flag.StringVar(&config.BackgroundColor, "bg_color", config.BackgroundColor, "Background color")
	flag.StringVar(&config.Gap, "gap", config.Gap, "Gap between sections")

	// flag.StringVar(&a.NoConfig, "no_config", "", "Disable config loading")
	flag.StringVar(&config.Stdout, "stdout", config.Stdout, "Stdout output mode")
	// flag.StringVar(&a.Verbose, "v", "", "Verbose output")
	// flag.StringVar(&a.JSON, "json", "", "Output JSON")
	// flag.StringVar(&a.Version, "version", "", "Show version")
	// flag.StringVar(&a.Help, "help", "", "Show help")
	// flag.StringVar(&a.Clean, "clean", "", "Clean cache and exit")
}

func parseArgs(config *configs.Config) {
	// TODO: if user setup config file in args. handle it first
	setupArgs(config)

	flag.Parse()

	if config.CPUTemp == "on" {
		config.CPUTemp = "C"
	}
}
