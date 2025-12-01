package configs

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func LoadConfig(fileName string) (*Config, error) {
	config := DefaultConfig()
	if fileName == "" {
		return &config, nil
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist, %w", err)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	err = ParseConfig(string(data), &config)
	if err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return &config, nil
}

func DefaultConfig() Config {
	return Config{
		Title: Title{
			TitleFqdn: "off",
		},
		Kernel: Kernel{
			KernelShorthand: "on",
		},
		Distro: Distro{
			DistroShorthand: "off",
			OSArch:          "on",
		},
		Uptime: Uptime{
			UptimeShorthand: "on",
		},
		Memory: Memory{
			MemoryPercent: "on",
			MemoryUnit:    "gib",
			MemPrecision:  "2",
		},
		Packages: Packages{
			PackageManagers: "on",
			PackageSeparate: "on",
			PackageMinimal:  "",
		},
		Shell: Shell{
			ShellPath:    "off",
			ShellVersion: "on",
		},
		Editor: Editor{
			EditorPath:    "off",
			EditorVersion: "on",
		},
		CPU: CPU{
			SpeedType:      "bios_limit",
			SpeedShorthand: "on",
			CPUBrand:       "on",
			CPUSpeed:       "on",
			CPUCores:       "logical",
			CPUTemp:        "off",
		},
		GPU: GPU{
			GPUBrand: "on",
			GPUType:  "all",
		},
		Resolution: Resolution{
			RefreshRate: "on",
		},
		GTK: GTK{
			GTKShorthand: "off",
			GTK2:         "on",
			GTK3:         "on",
			QT:           "on",
		},
		IPAddress: IPAddress{
			PublicIPHost:     "http://ident.me",
			PublicIPTimeout:  "2",
			LocalIPInterface: "('auto')",
		},
		DE: DE{
			DeVersion: "on",
		},
		Disk: Disk{
			DiskShow:     "('/')",
			DiskSubtitle: "mount",
			DiskPercent:  "on",
		},
		Song: Song{
			MusicPlayer:   "auto",
			SongFormat:    "%artist% - %album% - %title%",
			SongShorthand: "off",
			MPCArgs:       "()",
		},
		TextColor: TextColor{
			Colors: "(distro)",
		},
		TextOptions: TextOptions{
			Bold:             "on",
			UnderlineEnabled: "on",
			UnderlineChar:    "-",
			Separator:        ":",
		},
		ColorBlocks: ColorBlocks{
			BlockRange:  "(0 15)",
			ColorBlocks: "on",
			BlockWidth:  "3",
			BlockHeight: "1",
			ColOffset:   "auto",
		},
		ProgressBar: ProgressBar{
			BarCharElapsed:  "-",
			BarCharTotal:    "=",
			BarBorder:       "on",
			BarLength:       "15",
			BarColorElapsed: "distro",
			BarColorTotal:   "distro",
			MemoryDisplay:   "off",
			BatteryDisplay:  "off",
			DiskDisplay:     "off",
		},
		BackendSettings: BackendSettings{
			ImageBackend: "ascii",
			ImageSource:  "auto",
		},
		ASCIIOptions: ASCIIOptions{
			ASCIIDistro: "auto",
			ASCIIColors: "(distro)",
			ASCIIBold:   "on",
		},
		ImageOptions: ImageOptions{
			ImageLoop:       "off",
			ThumbnailDir:    "${XDG_CACHE_HOME:-${HOME}/.cache}/thumbnails/neofetch",
			CropMode:        "normal",
			CropOffset:      "center",
			ImageSize:       "auto",
			CatimgSize:      "2",
			Gap:             "3",
			Yoffset:         "0",
			Xoffset:         "0",
			BackgroundColor: "",
		},
		MiscOptions: MiscOptions{
			Stdout: "auto",
		},
		PrintInfo: PrintInfo{
			InfoList: [][]string{
				{"title"},
				{"underline"},
				{"OS", "distro"},
				{"Host", "model"},
				{"Kernel", "kernel"},
				{"Uptime", "uptime"},
				{"Packages", "packages"},
				{"Shell", "shell"},
				{"Editor", "editor"},
				{"Resolution", "resolution"},
				{"DE", "de"},
				{"WM", "wm"},
				{"WM Theme", "wm_theme"},
				{"Theme", "theme"},
				{"Icons", "icons"},
				{"Cursor", "cursor"},
				{"Terminal", "term"},
				{"Terminal Font", "term_font"},
				{"CPU", "cpu"},
				{"GPU", "gpu"},
				{"Memory", "memory"},
				{"Network", "network"},
				{"Bluetooth", "bluetooth"},
				{"BIOS", "bios"},
				{"cols"},
			},
		},
	}
}

func ParseConfig(data string, config *Config) error {
	lines := strings.Split(data, "\n")
	real_lines := []string{}

	// preprocess lines: remove comments and empty lines
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "#") || len(line) == 0 {
			continue
		}
		real_lines = append(real_lines, line)
	}

	print_info_lines_start := -1
	for i, line := range real_lines {
		if line == "print_info() {" {
			print_info_lines_start = i
			break
		}
	}
	if print_info_lines_start == -1 {
		return fmt.Errorf("invalid config format,cannot find print_info block's start")
	}

	print_info_lines_end := -1
	for i, line := range real_lines {
		if line == "}" {
			print_info_lines_end = i
			break
		}
	}

	if print_info_lines_end == -1 {
		return fmt.Errorf("invalid config format,cannot find print_info block's end")
	}
	if config.InfoList == nil {
		for i := print_info_lines_start + 1; i < print_info_lines_end; i++ {
			var info_line_list []string
			if strings.Contains(real_lines[i], "\"") {
				var tmp []string
				info_line_list = strings.Split(real_lines[i], "\"")
				for _, item := range info_line_list {
					item = strings.TrimSpace(item)
					tmp = append(tmp, item)
				}
				info_line_list = tmp
			} else {
				info_line_list = strings.Split(real_lines[i], " ")
			}
			var tmp []string
			for _, item := range info_line_list {
				if item == "" || item == "info" {
					continue
				}
				tmp = append(tmp, strings.ReplaceAll(item, "\"", ""))
			}

			config.InfoList = append(config.InfoList, tmp)
		}
	}

	config_list := real_lines[print_info_lines_end+1:]

	t := reflect.ValueOf(config)

	for _, line := range config_list {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		_ = SetByGFTag(&t, key, value)
	}
	return nil
}
