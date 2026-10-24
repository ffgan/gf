package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ffgan/gf/configs"
	"github.com/ffgan/gf/internal/hyfetch"
	"github.com/ffgan/gf/internal/modules"
	"github.com/ffgan/gf/internal/utils"
)

// Version is reported by --version and -V. Releases override it with
// -ldflags "-X github.com/ffgan/gf/cmd.Version=...".
var Version = "0.3.6"

var (
	showVersion bool
	printConfig bool
	noConfig    bool
	stdoutMode  bool
	verbose     bool
	vverbose    bool
	jsonMode    bool
	listModules bool
	structure   string
)

func setupArgs(config *configs.Config) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), `gf %s
Show system information beside an ASCII logo, in the spirit of neofetch.

Usage:
  gf [options]

Options:
`, Version)
		fs.PrintDefaults()
		fmt.Fprint(fs.Output(), hyfetch.HelpText)
	}
	flag.CommandLine = fs
	fs.StringVar(&config.TitleFqdn, "title_fqdn", config.TitleFqdn, "Show the fully qualified domain name (on, off)")
	fs.StringVar(&config.PackageManagers, "package_managers", config.PackageManagers, "Package line style: on, tiny, off")
	fs.StringVar(&config.PackageSeparate, "package_separate", config.PackageSeparate, "Split user and system packages when supported (on, off)")
	fs.StringVar(&config.PackageMinimal, "package_minimal", config.PackageMinimal, "Hide language package managers (on, off)")
	fs.StringVar(&config.OSArch, "os_arch", config.OSArch, "Append the architecture to the OS line (on, off)")
	fs.StringVar(&config.CPUCores, "cpu_cores", config.CPUCores, "CPU cores to show: logical, physical, off")
	fs.StringVar(&config.CPUSpeed, "cpu_speed", config.CPUSpeed, "Show the CPU speed (on, off)")
	fs.StringVar(&config.SpeedType, "speed_type", config.SpeedType, "cpufreq file to read, such as bios_limit or scaling_max_freq")
	fs.StringVar(&config.SpeedShorthand, "speed_shorthand", config.SpeedShorthand, "Print GHz with one decimal (on, off)")
	fs.StringVar(&config.DistroShorthand, "distro_shorthand", config.DistroShorthand, "Distro name length: off, on, tiny")
	fs.StringVar(&config.KernelShorthand, "kernel_shorthand", config.KernelShorthand, "Hide the kernel name (on, off)")
	fs.StringVar(&config.UptimeShorthand, "uptime_shorthand", config.UptimeShorthand, "Uptime style: on, tiny, off")
	fs.StringVar(&config.CPUBrand, "cpu_brand", config.CPUBrand, "Show the CPU vendor (on, off)")
	fs.StringVar(&config.GPUBrand, "gpu_brand", config.GPUBrand, "Show the GPU vendor (on, off)")
	fs.StringVar(&config.GPUType, "gpu_type", config.GPUType, "Which GPUs to show: all, dedicated, integrated")
	fs.StringVar(&config.RefreshRate, "refresh_rate", config.RefreshRate, "Append the refresh rate to the resolution (on, off)")
	fs.StringVar(&config.DeVersion, "de_version", config.DeVersion, "Show the desktop version (on, off)")
	fs.StringVar(&config.GTKShorthand, "gtk_shorthand", config.GTKShorthand, "Shorten GTK theme names (on, off)")
	fs.StringVar(&config.GTK2, "gtk2", config.GTK2, "Include GTK2 theme info (on, off)")
	fs.StringVar(&config.GTK3, "gtk3", config.GTK3, "Include GTK3 theme info (on, off)")
	fs.StringVar(&config.QT, "qt", config.QT, "Include Qt theme info (on, off)")
	fs.StringVar(&config.ShellPath, "shell_path", config.ShellPath, "Show the path to the shell (on, off)")
	fs.StringVar(&config.ShellVersion, "shell_version", config.ShellVersion, "Show the shell version (on, off)")
	fs.StringVar(&config.EditorPath, "editor_path", config.EditorPath, "Show the path to the editor (on, off)")
	fs.StringVar(&config.EditorVersion, "editor_version", config.EditorVersion, "Show the editor version (on, off)")
	fs.StringVar(&config.PublicIPHost, "ip_host", config.PublicIPHost, "URL used to look up the public IP")
	fs.StringVar(&config.PublicIPTimeout, "ip_timeout", config.PublicIPTimeout, "Public IP lookup timeout in seconds")
	fs.StringVar(&config.LocalIPInterface, "ip_interface", config.LocalIPInterface, "Interface for the local IP (auto, eth0, ...)")
	fs.StringVar(&config.SongFormat, "song_format", config.SongFormat, "Song format using %artist%, %album% and %title%")
	fs.StringVar(&config.SongShorthand, "song_shorthand", config.SongShorthand, "Reserved song layout switch (on, off)")
	fs.StringVar(&config.MusicPlayer, "music_player", config.MusicPlayer, "playerctl player name, or auto")
	fs.StringVar(&config.MemoryPercent, "memory_percent", config.MemoryPercent, "Show memory percentage (on, off)")
	fs.StringVar(&config.MemoryUnit, "memory_unit", config.MemoryUnit, "Memory unit: kib, mib, gib, tib")
	fs.StringVar(&config.MemPrecision, "memory_precision", config.MemPrecision, "Digits after the decimal point for memory")
	fs.StringVar(&config.CPUTemp, "cpu_temp", config.CPUTemp, "CPU temperature: off, C, F")
	fs.StringVar(&config.DiskSubtitle, "disk_subtitle", config.DiskSubtitle, "Disk label: mount, dir, name, none")
	fs.StringVar(&config.DiskPercent, "disk_percent", config.DiskPercent, "Show disk percentage (on, off)")
	fs.StringVar(&config.DiskShow, "disk_show", config.DiskShow, "Mounts to show, for example '/' or '(/ /home)'")
	fs.StringVar(&config.UnderlineEnabled, "underline", config.UnderlineEnabled, "Draw the line under the title (on, off)")
	fs.StringVar(&config.UnderlineChar, "underline_char", config.UnderlineChar, "Character used for the underline")
	fs.StringVar(&config.Bold, "bold", config.Bold, "Bold labels and title (on, off)")
	fs.StringVar(&config.Separator, "separator", config.Separator, "Text between a label and its value")
	fs.StringVar(&config.Enabled, "color_blocks", config.Enabled, "Show the color palette (on, off)")
	fs.StringVar(&config.BlockRange, "block_range", config.BlockRange, "Palette range, for example '0 15'")
	fs.StringVar(&config.BlockWidth, "block_width", config.BlockWidth, "Width of each color block")
	fs.StringVar(&config.BlockHeight, "block_height", config.BlockHeight, "Height of each color block row")
	fs.StringVar(&config.ColOffset, "col_offset", config.ColOffset, "Kept for neofetch config compatibility")
	fs.StringVar(&config.BarBorder, "bar_border", config.BarBorder, "Draw progress bar borders (on, off)")
	fs.StringVar(&config.BarLength, "bar_length", config.BarLength, "Progress bar length")
	fs.StringVar(&config.MemoryDisplay, "memory_display", config.MemoryDisplay, "Memory layout: off, bar, infobar, barinfo")
	fs.StringVar(&config.BatteryDisplay, "battery_display", config.BatteryDisplay, "Battery layout: off, bar, infobar, barinfo")
	fs.StringVar(&config.DiskDisplay, "disk_display", config.DiskDisplay, "Disk layout: off, bar, infobar, barinfo")
	fs.StringVar(&config.ImageBackend, "backend", config.ImageBackend, "Logo backend: ascii or off")
	fs.StringVar(&config.ImageSource, "source", config.ImageSource, "auto, a distro name, an ASCII file, or an image")
	fs.StringVar(&config.ImageLoop, "loop", config.ImageLoop, "Kept for neofetch config compatibility")
	fs.StringVar(&config.ImageSize, "image_size", config.ImageSize, "Kept for neofetch config compatibility")
	fs.StringVar(&config.ImageSize, "size", config.ImageSize, "Alias of --image_size")
	fs.StringVar(&config.CatimgSize, "catimg_size", config.CatimgSize, "catimg block size when an image backend is used")
	fs.StringVar(&config.CropMode, "crop_mode", config.CropMode, "Kept for neofetch config compatibility")
	fs.StringVar(&config.CropOffset, "crop_offset", config.CropOffset, "Kept for neofetch config compatibility")
	fs.StringVar(&config.Xoffset, "xoffset", config.Xoffset, "Kept for neofetch config compatibility")
	fs.StringVar(&config.Yoffset, "yoffset", config.Yoffset, "Kept for neofetch config compatibility")
	fs.StringVar(&config.BackgroundColor, "bg_color", config.BackgroundColor, "Kept for neofetch config compatibility")
	fs.StringVar(&config.Gap, "gap", config.Gap, "Spaces between the logo and the text")
	fs.StringVar(&config.ASCIIDistro, "ascii_distro", config.ASCIIDistro, "Force a logo by distro name")
	fs.StringVar(&config.ASCIIDistro, "logo", config.ASCIIDistro, "Alias of --ascii_distro")
	fs.StringVar(&config.ASCIIColors, "ascii_colors", config.ASCIIColors, "Logo palette: distro, or indexes such as '1 7'")
	fs.StringVar(&config.ASCIIBold, "ascii_bold", config.ASCIIBold, "Bold the logo (on, off)")
	fs.StringVar(&config.Colors, "colors", config.Colors, "Text palette: distro, or six color indexes")
	fs.BoolVar(&stdoutMode, "stdout", false, "Plain text, without a logo or colors")
	fs.BoolVar(&noConfig, "no_config", false, "Ignore config files")
	fs.BoolVar(&verbose, "v", false, "Print a note when a value cannot be detected")
	fs.BoolVar(&vverbose, "vv", false, "Alias of -v")
	fs.BoolVar(&printConfig, "print_config", false, "Print a config file and exit")
	fs.BoolVar(&showVersion, "version", false, "Print the version and exit")
	fs.BoolVar(&jsonMode, "json", false, "Print the configured info lines as JSON")
	fs.BoolVar(&listModules, "list_modules", false, "Print every module name and exit")
	fs.StringVar(&structure, "structure", "", "Modules to show, separated by colons, or all")
	fs.StringVar(&structure, "s", "", "Alias of --structure")
	fs.StringVar(&config.CustomText, "custom", config.CustomText, "Text printed by the custom module")
	fs.StringVar(&config.CommandText, "command", config.CommandText, "Shell command run by the command module")
	fs.StringVar(&config.WeatherLocation, "weather_location", config.WeatherLocation, "City passed to the weather module")
	fs.StringVar(&config.WeatherFormat, "weather_format", config.WeatherFormat, "wttr.in format for the weather module")
	fs.StringVar(&config.SeparatorString, "separator_string", config.SeparatorString, "Character repeated by the separator module")
	fs.StringVar(&configPathFlag, "config", "", "Config file to load")
}

var configPathFlag string

func parseArgs(config *configs.Config) map[string]bool {
	rest, disabled := extractDisable(os.Args[1:])
	setupArgs(config)
	if err := flag.CommandLine.Parse(rest); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		os.Exit(2)
	}
	if stdoutMode {
		config.Stdout = utils.ON
	}
	if vverbose {
		verbose = true
	}
	if showVersion {
		fmt.Printf("gf %s\n", Version)
		os.Exit(0)
	}
	if printConfig {
		configs.WriteDefault(os.Stdout)
		os.Exit(0)
	}
	if listModules {
		for _, name := range modules.Names() {
			fmt.Println(name)
		}
		os.Exit(0)
	}
	if structure != "" {
		lines, err := modules.Structure(structure)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gf: %v\n", err)
			os.Exit(2)
		}
		config.InfoList = lines
	}
	return disabled
}

func extractDisable(args []string) ([]string, map[string]bool) {
	var rest []string
	disabled := map[string]bool{}
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--disable" || args[i] == "--off":
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				addDisabled(disabled, args[i])
			}
		case strings.HasPrefix(args[i], "--disable="), strings.HasPrefix(args[i], "--off="):
			_, value, _ := strings.Cut(args[i], "=")
			addDisabled(disabled, value)
		default:
			rest = append(rest, args[i])
		}
	}
	return rest, disabled
}

func addDisabled(dst map[string]bool, value string) {
	for _, part := range strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ' '
	}) {
		part = strings.ToLower(strings.TrimSpace(part))
		if part != "" {
			dst[part] = true
		}
	}
}
