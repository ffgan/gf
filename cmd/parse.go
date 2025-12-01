package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/ffgan/gf/configs"
	"github.com/ffgan/gf/internal/utils"
)

const Version = "0.0.1"

var show_version bool

func setupArgs(config *configs.Config) {
	// Info
	flag.StringVar(&config.TitleFqdn, "title_fqdn", config.TitleFqdn, `# Hide/Show Fully qualified domain name. 
# 
# Default:  'off' 
# Values:   'on', 'off' 
# Flag:     --title_fqdn`)
	flag.StringVar(&config.PackageManagers, "package_managers", config.PackageManagers, `# Show/Hide Package Manager names.
#
# Default: 'tiny'
# Values:  'on', 'tiny' 'off'
# Flag:    --package_managers
#
# Example:
# on:   '998 (pacman), 8 (flatpak), 4 (snap)'
# tiny: '908 (pacman, flatpak, snap)'
# off:  '908'`)
	flag.StringVar(&config.PackageSeparate, "package_separate", config.PackageSeparate, `# Show separate user and system packages for supported package managers
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --package_separate
#
# Example:
# on:  '8 packages (flatpak-system), 9 packages (flatpak-user)'
# off: '17 packages (flatpak)'`)
	flag.StringVar(&config.PackageMinimal, "package_minimal", config.PackageMinimal, `# Reduce output of packages list by not showing programming language package managers or Steam games
#
# Flag:    --package_minimal
#
# Example:
# default:  'Packages: 1 (npm), 991 (emerge), 3 (steam), 23 (flatpak-system)'
# minimal: 'Packages: 991 (emerge), 23 (flatpak-system)'`)
	flag.StringVar(&config.OSArch, "os_arch", config.OSArch, `# Show/Hide OS Architecture.
# Show 'x86_64', 'x86' and etc in 'Distro:' output.
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --os_arch
#
# Example:
# on:  'Arch Linux x86_64'
# off: 'Arch Linux'`)
	flag.StringVar(&config.CPUCores, "cpu_cores", config.CPUCores, `# CPU Cores
# Display CPU cores in output
#
# Default: 'logical'
# Values:  'logical', 'physical', 'off'
# Flag:    --cpu_cores
# Support: 'physical' doesn't work on BSD.
#
# Example:
# logical:  'Intel i7-6500U (4) @ 3.1GHz' (All virtual cores)
# physical: 'Intel i7-6500U (2) @ 3.1GHz' (All physical cores)
# off:      'Intel i7-6500U @ 3.1GHz'`)
	flag.StringVar(&config.CPUSpeed, "cpu_speed", config.CPUSpeed, `# CPU Speed
# Hide/Show CPU speed.
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --cpu_speed
#
# Example:
# on:  'Intel i7-6500U (4) @ 3.1GHz'
# off: 'Intel i7-6500U (4)'`)
	flag.StringVar(&config.SpeedType, "speed_type", config.SpeedType, `# CPU speed type
#
# Default: 'bios_limit'
# Values: 'scaling_cur_freq', 'scaling_min_freq', 'scaling_max_freq', 'bios_limit'.
# Flag:    --speed_type
# Supports: Linux with 'cpufreq'
# NOTE: Any file in '/sys/devices/system/cpu/cpu0/cpufreq' can be used as a value.`)
	flag.StringVar(&config.SpeedShorthand, "speed_shorthand", config.SpeedShorthand, `# CPU speed shorthand
#
# Default: 'off'
# Values: 'on', 'off'.
# Flag:    --speed_shorthand
# NOTE: This flag is not supported in systems with CPU speed less than 1 GHz
#
# Example:
# on:    'i7-6500U (4) @ 3.1GHz'
# off:   'i7-6500U (4) @ 3.100GHz'`)
	flag.StringVar(&config.DistroShorthand, "distro_shorthand", config.DistroShorthand, `# Shorten the output of the distro function
#
# Default:  'off'
# Values:   'on', 'tiny', 'off'
# Flag:     --distro_shorthand
# Supports: Everything except Windows and Haiku`)
	flag.StringVar(&config.KernelShorthand, "kernel_shorthand", config.KernelShorthand, `# Shorten the output of the kernel function.
#
# Default:  'on'
# Values:   'on', 'off'
# Flag:     --kernel_shorthand
# Supports: Everything except *BSDs (except PacBSD and PC-BSD)
#
# Example:
# on:  '4.8.9-1-ARCH'
# off: 'Linux 4.8.9-1-ARCH'`)
	flag.StringVar(&config.UptimeShorthand, "uptime_shorthand", config.UptimeShorthand, `# Shorten the output of the uptime function
#
# Default: 'on'
# Values:  'on', 'tiny', 'off'
# Flag:    --uptime_shorthand
#
# Example:
# on:   '2 days, 10 hours, 3 mins'
# tiny: '2d 10h 3m'
# off:  '2 days, 10 hours, 3 minutes'`)
	flag.StringVar(&config.CPUBrand, "cpu_brand", config.CPUBrand, `# Enable/Disable CPU brand in output.
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --cpu_brand
#
# Example:
# on:   'Intel i7-6500U'
# off:  'i7-6500U (4)'`)
	flag.StringVar(&config.GPUBrand, "gpu_brand", config.GPUBrand, `# Enable/Disable GPU Brand
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --gpu_brand
#
# Example:
# on:  'AMD HD 7950'
# off: 'HD 7950'`)
	flag.StringVar(&config.GPUType, "gpu_type", config.GPUType, `# Which GPU to display
#
# Default: 'all'
# Values:  'all', 'dedicated', 'integrated'
# Flag:    --gpu_type
# Supports: Linux
#
# Example:
# all:
#   GPU1: AMD HD 7950
#   GPU2: Intel Integrated Graphics
#
# dedicated:
#   GPU1: AMD HD 7950
#
# integrated:
#   GPU1: Intel Integrated Graphics`)
	flag.StringVar(&config.RefreshRate, "refresh_rate", config.RefreshRate, `# Display refresh rate next to each monitor
# Default: 'off'
# Values:  'on', 'off'
# Flag:    --refresh_rate
# Supports: Doesn't work on Windows.
#
# Example:
# on:  '1920x1080 @ 60Hz'
# off: '1920x1080'`)
	flag.StringVar(&config.DeVersion, "de_version", config.DeVersion, `# Show Desktop Environment version
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --de_version`)
	flag.StringVar(&config.GTKShorthand, "gtk_shorthand", config.GTKShorthand, `# Shorten output of GTK Theme / Icons / Font
#
# Default: 'off'
# Values:  'on', 'off'
# Flag:    --gtk_shorthand
#
# Example:
# on:  'Numix, Adwaita'
# off: 'Numix [GTK2], Adwaita [GTK3]'`)
	flag.StringVar(&config.GTK2, "gtk2", config.GTK2, `# Enable/Disable gtk2 Theme / Icons / Font
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --gtk2
#
# Example:
# on:  'Numix [GTK2], Adwaita [GTK3]'
# off: 'Adwaita [GTK3]'`)
	flag.StringVar(&config.GTK3, "gtk3", config.GTK3, `# Enable/Disable gtk3 Theme / Icons / Font
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --gtk3
#
# Example:
# on:  'Numix [GTK2], Adwaita [GTK3]'
# off: 'Numix [GTK2]'`)
	flag.StringVar(&config.QT, "qt", config.QT, `# Enable/Disable Qt Theme / Icons / Font
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --qt
#
# Example:
# on:  'Breeze [Qt], Arc [GTK3]'
# off: 'Arc [GTK3]'`)
	flag.StringVar(&config.ShellPath, "shell_path", config.ShellPath, `# Show the path to $SHELL
#
# Default: 'off'
# Values:  'on', 'off'
# Flag:    --shell_path
#
# Example:
# on:  '/bin/bash'
# off: 'bash'`)
	flag.StringVar(&config.ShellVersion, "shell_version", config.ShellVersion, `# Show $SHELL version
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --shell_version
#
# Example:
# on:  'bash 4.4.5'
# off: 'bash'`)
	flag.StringVar(&config.EditorPath, "editor_path", config.EditorPath, `# Show path to $EDITOR
#
# Default: 'off'
# Values:  'on', 'off'
# Flag:    --editor_path
#
# Example:
# on:  '/opt/bin/vim'
# off: 'vim'`)
	flag.StringVar(&config.EditorVersion, "editor_version", config.EditorVersion, `# Show $EDITOR version
#
# Default:  'on'
# Values:   'on', 'off'
# Flag:     '--editor_version'
#
# Example:
# on:  'vim 9.0'
# off: 'vim'`)
	flag.StringVar(&config.PublicIPHost, "ip_host", config.PublicIPHost, `# Website to ping for the public IP
#
# Default: 'http://ident.me'
# Values:  'url'
# Flag:    --ip_host`)
	flag.StringVar(&config.PublicIPTimeout, "ip_timeout", config.PublicIPTimeout, `# Public IP timeout.
#
# Default: '2'
# Values:  'int'
# Flag:    --ip_timeout`)

	// TODO: 需要特殊处理
	flag.StringVar(&config.LocalIPInterface, "ip_interface", config.LocalIPInterface,
		`# Local IP interface
#
# Default: 'auto' (interface of default route)
# Values:  'auto', 'en0', 'en1'
# Flag:    --ip_interface`)

	flag.StringVar(&config.SongFormat, "song_format", config.SongFormat, `# Format to display song information.
#
# Default: '%artist% - %album% - %title%'
# Values:  '%artist%', '%album%', '%title%'
# Flag:    --song_format
#
# Example:
# default: 'Song: Jet - Get Born - Sgt Major'`)
	flag.StringVar(&config.SongShorthand, "song_shorthand", config.SongShorthand, `# Print the Artist, Album and Title on separate lines
#
# Default: 'off'
# Values:  'on', 'off'
# Flag:    --song_shorthand
#
# Example:
# on:  'Artist: The Fratellis'
#      'Album: Costello Music'
#      'Song: Chelsea Dagger'
#
# off: 'Song: The Fratellis - Costello Music - Chelsea Dagger'`)
	flag.StringVar(&config.MusicPlayer, "music_player", config.MusicPlayer, `# Manually specify a music player.
#
# Default: 'auto'
# Values:  'auto', 'player-name'
# Flag:    --music_player
#
# Available values for 'player-name':
#
# amarok
# audacious
# banshee
# bluemindo
# cider
# clementine
# cmus
# deadbeef
# deepin-music
# dragon
# elisa
# exaile
# gnome-music
# gmusicbrowser
# gogglesmm
# guayadeque
# io.elementary.music
# iTunes
# Music
# juk
# lollypop
# MellowPlayer
# mocp
# mopidy
# mpd
# muine
# netease-cloud-music
# olivia
# playerctl
# pogo
# pragha
# qmmp
# quodlibet
# rhythmbox
# sayonara
# smplayer
# spotify
# strawberry
# tauonmb
# tomahawk
# vlc
# xmms2d
# xnoise
# yarock`)
	flag.StringVar(&config.MemoryPercent, "memory_percent", config.MemoryPercent, `# Show memory percentage in output.
#
# Default: 'off'
# Values:  'on', 'off'
# Flag:    --memory_percent
#
# Example:
# on:   '1801MiB / 7881MiB (22%)'
# off:  '1801MiB / 7881MiB'`)
	flag.StringVar(&config.MemoryUnit, "memory_unit", config.MemoryUnit, `# Change memory output unit.
#
# Default: 'mib'
# Values:  'kib', 'mib', 'gib', 'tib'
# Flag:    --memory_unit
#
# Example:
# kib  '1020928KiB / 7117824KiB'
# mib  '1042MiB / 6951MiB'
# gib: ' 0.98GiB / 6.79GiB'`)
	flag.StringVar(&config.MemPrecision, "memory_precision", config.MemPrecision, `# Change memory output precision.
#
# Default: '2'
# Values: integer ≥ 0
# Flag:    --memory_precision`)

	// TODO: 需要特殊处理
	flag.StringVar(&config.CPUTemp, "cpu_temp", config.CPUTemp, `# CPU Temperature
# Hide/Show CPU temperature.
# Note the temperature is added to the regular CPU function.
#
# Default: 'off'
# Values:  'C', 'F', 'off'
# Flag:    --cpu_temp
# Supports: Linux, BSD
# NOTE: For FreeBSD and NetBSD-based systems, you'll need to enable
#       coretemp kernel module. This only supports newer Intel processors.
#
# Example:
# C:   'Intel i7-6500U (4) @ 3.1GHz [27.2°C]'
# F:   'Intel i7-6500U (4) @ 3.1GHz [82.0°F]'
# off: 'Intel i7-6500U (4) @ 3.1GHz'`)
	if config.CPUTemp == utils.ON {
		config.CPUTemp = "C"
	}

	flag.StringVar(&config.DiskSubtitle, "disk_subtitle", config.DiskSubtitle, `# Disk subtitle.
# What to append to the Disk subtitle.
#
# Default: 'mount'
# Values:  'mount', 'name', 'dir', 'none'
# Flag:    --disk_subtitle
#
# Example:
# name:   'Disk (/dev/sda1): 74G / 118G (66%)'
#         'Disk (/dev/sdb2): 74G / 118G (66%)'
#
# mount:  'Disk (/): 74G / 118G (66%)'
#         'Disk (/mnt/Local Disk): 74G / 118G (66%)'
#         'Disk (/mnt/Videos): 74G / 118G (66%)'
#
# dir:    'Disk (/): 74G / 118G (66%)'
#         'Disk (Local Disk): 74G / 118G (66%)'
#         'Disk (Videos): 74G / 118G (66%)'
#
# none:   'Disk: 74G / 118G (66%)'
#         'Disk: 74G / 118G (66%)'
#         'Disk: 74G / 118G (66%)'`)
	flag.StringVar(&config.DiskPercent, "disk_percent", config.DiskPercent, `# Disk percent.
# Show/Hide disk percent.
#
# Default: 'on'
# Values:  'on', 'off'
# Flag:    --disk_percent
#
# Example:
# on:  'Disk (/): 74G / 118G (66%)'
# off: 'Disk (/): 74G / 118G'`)

	// TODO: 需要特殊处理
	flag.StringVar(&config.DiskShow, "disk_show", config.DiskShow, `# Which disks to display.
# The values can be any /dev/sdXX, mount point or directory.
# NOTE: By default we only show the disk info for '/'.
#
# Default: '/'
# Values:  '/', '/dev/sdXX', '/path/to/drive'.
# Flag:    --disk_show
#
# Example:
# disk_show=('/' '/dev/sdb1'):
#      'Disk (/): 74G / 118G (66%)'
#      'Disk (/mnt/Videos): 823G / 893G (93%)'
#
# disk_show=('/'):
#      'Disk (/): 74G / 118G (66%)'
#`)

	// TODO: 需要特殊处理
	// --disable
	// Text Colors
	// --colors

	// Text Formatting
	flag.StringVar(&config.UnderlineEnabled, "underline", config.UnderlineEnabled, `# Enable/Disable Underline
#
# Default:  'on'
# Values:   'on', 'off'
# Flag:     --underline`)
	flag.StringVar(&config.UnderlineChar, "underline_char", config.UnderlineChar, `# Underline character
#
# Default:  '-'
# Values:   'string'
# Flag:     --underline_char`)
	flag.StringVar(&config.Bold, "bold", config.Bold, `# Toggle bold text
#
# Default:  'on'
# Values:   'on', 'off'
# Flag:     --bold`)
	flag.StringVar(&config.Separator, "separator", config.Separator, `# Info Separator
# Replace the default separator with the specified string.
#
# Default:  ':'
# Flag:     --separator
#
# Example:
# separator="->":   'Shell-> bash'
# separator=" =":   'WM = dwm'`)

	//  Color Blocks
	flag.StringVar(&config.Colors, "color_blocks", config.Colors, `# Toggle color blocks
#
# Default:  'on'
# Values:   'on', 'off'
# Flag:     --color_blocks`)
	flag.StringVar(&config.BlockRange, "block_range", config.BlockRange, `# Color block range
# The range of colors to print.
#
# Default:  '0', '15'
# Values:   'num'
# Flag:     --block_range
#
# Example:
#
# Display colors 0-7 in the blocks.  (8 colors)
# neofetch --block_range 0 7
#
# Display colors 0-15 in the blocks. (16 colors)
# neofetch --block_range 0 15`)
	flag.StringVar(&config.BlockWidth, "block_width", config.BlockWidth, `# Color block width in spaces
#
# Default:  '3'
# Values:   'num'
# Flag:     --block_width`)
	flag.StringVar(&config.BlockHeight, "block_height", config.BlockHeight, `# Color block height in lines
#
# Default:  '1'
# Values:   'num'
# Flag:     --block_height`)
	flag.StringVar(&config.ColOffset, "col_offset", config.ColOffset, `# Color Alignment
#
# Default: 'auto'
# Values: 'auto', 'num'
# Flag: --col_offset
#
# Number specifies how far from the left side of the terminal (in spaces) to
# begin printing the columns, in case you want to e.g. center them under your
# text.
# Example:
# col_offset=utils.AUTO - Default behavior of neofetch
# col_offset=7      - Leave 7 spaces then print the colors`)

	// Bars
	// TODO: 需要特殊处理
	flag.StringVar(&config.BarCharElapsed, "bar_char", config.BarCharElapsed, `# Bar characters
#
# Default:  '-', '='
# Values:   'string', 'string'
# Flag:     --bar_char
#
# Example:
# neofetch --bar_char 'elapsed' 'total'
# neofetch --bar_char '-' '='`)

	flag.StringVar(&config.BarBorder, "bar_border", config.BarBorder, `# Toggle Bar border
#
# Default:  'on'
# Values:   'on', 'off'
# Flag:     --bar_border`)
	flag.StringVar(&config.BarLength, "bar_length", config.BarLength, `# Progress bar length in spaces
# Number of chars long to make the progress bars.
#
# Default:  '15'
# Values:   'num'
# Flag:     --bar_length`)

	// TODO: 需要特殊处理
	flag.StringVar(&config.BarColorElapsed, "bar_colors", config.BarColorElapsed, `# Progress bar colors
# When set to distro, uses your distro's logo colors.
#
# Default:  'distro', 'distro'
# Values:   'distro', 'num'
# Flag:     --bar_colors
#
# Example:
# neofetch --bar_colors 3 4
# neofetch --bar_colors distro 5`)

	flag.StringVar(&config.MemoryDisplay, "memory_display", config.MemoryDisplay, `# Display a bar with the info.
#
# Default: 'off'
# Values:  'bar', 'infobar', 'barinfo', 'off'
# Flags:   --memory_display
#          --battery_display
#          --disk_display
#
# Example:
# bar:     '[---=======]'
# infobar: 'info [---=======]'
# barinfo: '[---=======] info'
# off:     'info'`)
	flag.StringVar(&config.BatteryDisplay, "battery_display", config.BatteryDisplay, `# Display a bar with the info.
#
# Default: 'off'
# Values:  'bar', 'infobar', 'barinfo', 'off'
# Flags:   --memory_display
#          --battery_display
#          --disk_display
#
# Example:
# bar:     '[---=======]'
# infobar: 'info [---=======]'
# barinfo: '[---=======] info'
# off:     'info'`)
	flag.StringVar(&config.DiskDisplay, "disk_display", config.DiskDisplay, `# Display a bar with the info.
#
# Default: 'off'
# Values:  'bar', 'infobar', 'barinfo', 'off'
# Flags:   --memory_display
#          --battery_display
#          --disk_display
#
# Example:
# bar:     '[---=======]'
# infobar: 'info [---=======]'
# barinfo: '[---=======] info'
# off:     'info'`)

	// Image backend
	flag.StringVar(&config.ImageBackend, "backend", config.ImageBackend, `# Image backend.
#
# Default:  'ascii'
# Values:   'ascii', 'caca', 'catimg', 'chafa', 'jp2a', 'iterm2', 'off',
#           'pot', 'termpix', 'pixterm', 'tycat', 'w3m', 'kitty', 'ueberzug',
#           'viu'

# Flag:     --backend`)
	flag.StringVar(&config.ImageSource, "source", config.ImageSource, `# Image Source
#
# Which image or ascii file to display.
#
# Default:  'auto'
# Values:   'auto', 'ascii', 'wallpaper', '/path/to/img', '/path/to/ascii', '/path/to/dir/'
#           'command output (neofetch --ascii "$(fortune | cowsay -W 30)")'
# Flag:     --source
#
# NOTE: 'auto' will pick the best image source for whatever image backend is used.
#       In ascii mode, distro ascii art will be used and in an image mode, your
#       wallpaper will be used.`)

	// Image options
	flag.StringVar(&config.ImageLoop, "loop", config.ImageLoop, `# Image loop
# Setting this to on will make neofetch redraw the image constantly until
# Ctrl+C is pressed. This fixes display issues in some terminal emulators.
#
# Default:  'off'
# Values:   'on', 'off'
# Flag:     --loop`)
	flag.StringVar(&config.ImageSize, "image_size", config.ImageSize, `# Image size
# The image is half the terminal width by default.
#
# Default: 'auto'
# Values:  'auto', '00px', '00%', 'none'
# Flags:   --image_size
#          --size`)
	flag.StringVar(&config.CatimgSize, "catimg_size", config.CatimgSize, `# Catimg block size.
# Control the resolution of catimg.
#
# Default: '2'
# Values:  '1', '2'
# Flags:   --catimg_size`)
	flag.StringVar(&config.CropMode, "crop_mode", config.CropMode, `# Crop mode
#
# Default:  'normal'
# Values:   'normal', 'fit', 'fill'
# Flag:     --crop_mode
#
# See this wiki page to learn about the fit and fill options.
# https://github.com/dylanaraps/neofetch/wiki/What-is-Waifu-Crop%3F`)
	flag.StringVar(&config.CropOffset, "crop_offset", config.CropOffset, `# Crop offset
# Note: Only affects 'normal' crop mode.
#
# Default:  'center'
# Values:   'northwest', 'north', 'northeast', 'west', 'center'
#           'east', 'southwest', 'south', 'southeast'
# Flag:     --crop_offset`)
	flag.StringVar(&config.Xoffset, "xoffset", config.Xoffset, `# Image offsets
# Only works with the w3m backend.
#
# Default: '0'
# Values:  'px'
# Flags:   --xoffset
#          --yoffset`)
	flag.StringVar(&config.Yoffset, "yoffset", config.Yoffset, `# Image offsets
# Only works with the w3m backend.
#
# Default: '0'
# Values:  'px'
# Flags:   --xoffset
#          --yoffset`)
	flag.StringVar(&config.BackgroundColor, "background_color", config.BackgroundColor, `# Image background color
# Only works with the w3m backend.
#
# Default: ''
# Values:  'color', 'blue'
# Flag:    --bg_color`)
	flag.StringVar(&config.BackgroundColor, "bg_color", config.BackgroundColor, `# Image background color
# Only works with the w3m backend.
#
# Default: ''
# Values:  'color', 'blue'
# Flag:    --bg_color`)
	flag.StringVar(&config.Gap, "gap", config.Gap, `# Gap between image and text
#
# Default: '3'
# Values:  'num', '-num'
# Flag:    --gap`)
	// TODO: 需要特殊处理
	// --clean
	// --ascii_colors
	// --ascii_distro
	// --ascii_bold
	// --logo
	// -L

	// Other
	// TODO: 需要特殊处理
	// --config
	flag.StringVar(&no_config, "no_config", utils.OFF, "Don't create the user config file.")
	flag.StringVar(&config.Stdout, "stdout", config.Stdout, "Turn off all colors and disables any ASCII/image backend.")
	flag.StringVar(&verbose, "v", utils.ON, "Display error messages.")
	flag.StringVar(&print_config, "print_config", utils.ON, "Print the default config file to stdout.")
	flag.StringVar(&vverbose, "vv", utils.ON, "Display a verbose log for error reporting.")

	// --help

	flag.BoolVar(&show_version, "version", false, "Show gf version")

	// --gen-man
	// --json
	// --travis

}

var (
	no_config    string
	verbose      string
	print_config string
	vverbose     string
)

func parseArgs(config *configs.Config) {
	// TODO: if user setup config file in args. handle it first
	setupArgs(config)

	flag.Parse()
	if show_version {
		fmt.Printf("GF %s\n", Version)
		os.Exit(0)
	}
}
