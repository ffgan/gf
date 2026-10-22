package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/ffgan/gf/internal/utils"
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
			TitleFqdn: utils.OFF,
		},
		Kernel: Kernel{
			KernelShorthand: utils.ON,
		},
		Distro: Distro{
			DistroShorthand: utils.OFF,
			OSArch:          utils.ON,
		},
		Uptime: Uptime{
			UptimeShorthand: utils.ON,
		},
		Memory: Memory{
			MemoryPercent: utils.ON,
			MemoryUnit:    "gib",
			MemPrecision:  "2",
		},
		Packages: Packages{
			PackageManagers: utils.ON,
			PackageSeparate: utils.ON,
			PackageMinimal:  "",
		},
		Shell: Shell{
			ShellPath:    utils.OFF,
			ShellVersion: utils.ON,
		},
		Editor: Editor{
			EditorPath:    utils.OFF,
			EditorVersion: utils.ON,
		},
		CPU: CPU{
			SpeedType:      "bios_limit",
			SpeedShorthand: utils.ON,
			CPUBrand:       utils.ON,
			CPUSpeed:       utils.ON,
			CPUCores:       "logical",
			CPUTemp:        utils.OFF,
		},
		GPU: GPU{
			GPUBrand: utils.ON,
			GPUType:  "all",
		},
		Resolution: Resolution{
			RefreshRate: utils.ON,
		},
		GTK: GTK{
			GTKShorthand: utils.OFF,
			GTK2:         utils.ON,
			GTK3:         utils.ON,
			QT:           utils.ON,
		},
		IPAddress: IPAddress{
			PublicIPHost:     "http://ident.me",
			PublicIPTimeout:  "2",
			LocalIPInterface: "('auto')",
		},
		DE: DE{
			DeVersion: utils.ON,
		},
		Disk: Disk{
			DiskShow:     "('/')",
			DiskSubtitle: "mount",
			DiskPercent:  utils.ON,
		},
		Song: Song{
			MusicPlayer:   utils.AUTO,
			SongFormat:    "%artist% - %album% - %title%",
			SongShorthand: utils.OFF,
			MPCArgs:       "()",
		},
		TextColor: TextColor{
			Colors: "(distro)",
		},
		TextOptions: TextOptions{
			Bold:             utils.ON,
			UnderlineEnabled: utils.ON,
			UnderlineChar:    "-",
			Separator:        ":",
		},
		ColorBlocks: ColorBlocks{
			BlockRange:  "(0 15)",
			Enabled:     utils.ON,
			BlockWidth:  "3",
			BlockHeight: "1",
			ColOffset:   utils.AUTO,
		},
		ProgressBar: ProgressBar{
			BarCharElapsed:  "-",
			BarCharTotal:    "=",
			BarBorder:       utils.ON,
			BarLength:       "15",
			BarColorElapsed: "distro",
			BarColorTotal:   "distro",
			MemoryDisplay:   utils.OFF,
			BatteryDisplay:  utils.OFF,
			DiskDisplay:     utils.OFF,
		},
		BackendSettings: BackendSettings{
			ImageBackend: "ascii",
			ImageSource:  utils.AUTO,
		},
		ASCIIOptions: ASCIIOptions{
			ASCIIDistro: utils.AUTO,
			ASCIIColors: "(distro)",
			ASCIIBold:   utils.ON,
		},
		ImageOptions: ImageOptions{
			ImageLoop:       utils.OFF,
			ThumbnailDir:    "${XDG_CACHE_HOME:-${HOME}/.cache}/thumbnails/neofetch",
			CropMode:        "normal",
			CropOffset:      "center",
			ImageSize:       utils.AUTO,
			CatimgSize:      "2",
			Gap:             "3",
			Yoffset:         "0",
			Xoffset:         "0",
			BackgroundColor: "",
		},
		MiscOptions: MiscOptions{
			Stdout: utils.AUTO,
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
				{"Font", "font"},
				{"Terminal", "term"},
				{"Terminal Font", "term_font"},
				{"CPU", "cpu"},
				{"GPU", "gpu"},
				{"NPU", "npu"},
				{"GPU Driver", "gpu_driver"},
				{"Memory", "memory"},
				{"Swap", "swap"},
				{"Disk", "disk"},
				{"Battery", "battery"},
				{"Power Adapter", "power_adapter"},
				{"Local IP", "local_ip"},
				{"Locale", "locale"},
				{"Users", "users"},
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
	print_info_lines_end := -1

	for i, line := range real_lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "print_info() {" || trimmed == "print_info(){" {
			print_info_lines_start = i
		}
		if trimmed == "}" && print_info_lines_start != -1 && print_info_lines_end == -1 {
			print_info_lines_end = i
		}
	}

	if print_info_lines_start == -1 || print_info_lines_end == -1 || print_info_lines_end < print_info_lines_start {
		config_list := real_lines
		return applyAssignments(config, config_list)
	}

	config.InfoList = nil
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
	return applyAssignments(config, real_lines[print_info_lines_end+1:])
}

func applyAssignments(config *Config, lines []string) error {
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := cleanConfigValue(parts[1])
		_ = SetByGFTag(config, key, value)
	}
	return nil
}

func cleanConfigValue(value string) string {
	value = strings.TrimSpace(value)
	if i := commentIndex(value); i >= 0 {
		value = strings.TrimSpace(value[:i])
	}
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}
	return strings.TrimSpace(value)
}

func commentIndex(value string) int {
	inSingle, inDouble := false, false
	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				return i
			}
		}
	}
	return -1
}
