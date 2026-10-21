package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ffgan/gf/configs"
	cli "github.com/ffgan/gf/internal/CLI"
	"github.com/ffgan/gf/internal/logo"
	"github.com/ffgan/gf/internal/utils"
)

const neofetchConf = "/etc/neofetch/config.conf"

func Run(ASCIIFiles embed.FS) {
	kernelName, kernelVersion, kernelMachine, darwinName := cacheUname()
	osName := cli.GetOS(kernelName, darwinName)

	explicit, skipFile := configFromArgs(os.Args[1:])
	var filePath string
	switch {
	case explicit != "":
		filePath = explicit
	case !skipFile:
		filePath = findConfig()
	}

	config, err := configs.LoadConfig(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gf: %v\n", err)
		fallback := configs.DefaultConfig()
		config = &fallback
	}

	disabled := parseArgs(config)

	h := &hardware{
		cache_dir:      getCacheDir(),
		KernelName:     kernelName,
		kernel_version: kernelVersion,
		KernelMachine:  kernelMachine,
		DarwinName:     darwinName,
		os:             osName,
		Config:         *config,
		verbose:        verbose,
		disabled:       disabled,
	}
	if h.Stdout == utils.ON {
		h.plain = true
	}

	logo.SetAppearance(h.ASCIIBold != utils.OFF, h.ASCIIColors, h.Colors)
	h.Get_distro()

	var art []string
	if !h.plainMode() && !strings.EqualFold(h.ImageBackend, "off") {
		if lines := imageArt(h.ImageSource); len(lines) > 0 {
			art = lines
			logo.ApplyDistro(h.ascii_distro)
		} else {
			art = logo.PrintASCII(h.ascii_distro, h.ImageSource, ASCIIFiles)
		}
	} else {
		logo.ApplyDistro(h.ascii_distro)
	}

	show(art, h.Rows(), h.gap())
}

func findConfig() string {
	var candidates []string
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		candidates = append(candidates, filepath.Join(dir, "gf", "config.conf"))
	}
	if home := os.Getenv("HOME"); home != "" {
		candidates = append(candidates,
			filepath.Join(home, ".config", "gf", "config.conf"),
			filepath.Join(home, ".config", "neofetch", "config.conf"),
		)
	}
	candidates = append(candidates, "/etc/gf/config.conf", neofetchConf, "/etc/neofetch/default.conf")
	for _, path := range candidates {
		if utils.FileExists(path) {
			return path
		}
	}
	return ""
}

func configFromArgs(args []string) (path string, skip bool) {
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--no_config" || args[i] == "--noconfig":
			skip = true
		case args[i] == "--config" && i+1 < len(args):
			i++
			path = args[i]
		case strings.HasPrefix(args[i], "--config="):
			path = strings.TrimPrefix(args[i], "--config=")
		}
	}
	return path, skip
}

func getCacheDir() string {
	cacheDir := os.Getenv("TMPDIR")
	if cacheDir == "" {
		cacheDir = "/tmp"
	}
	return cacheDir
}

func cacheUname() (string, string, string, string) {
	parts := strings.Fields(cli.UName("-srm"))
	for len(parts) < 3 {
		parts = append(parts, "")
	}
	return parts[0], parts[1], parts[2], ""
}

func imageArt(source string) []string {
	fi, err := os.Stat(source)
	if err != nil || !fi.Mode().IsRegular() {
		return nil
	}
	switch strings.ToLower(filepath.Ext(source)) {
	case ".png", ".jpg", ".jpeg", ".jpe", ".gif", ".webp":
	default:
		return nil
	}
	tools := []struct {
		bin  string
		args []string
	}{
		{"chafa", []string{"--symbols", "block", "--size", "40x40", source}},
		{"catimg", []string{"-w", "40", source}},
		{"viu", []string{"-w", "40", source}},
	}
	for _, tool := range tools {
		if !utils.CommandExists(tool.bin) {
			continue
		}
		out := utils.RunCommand(tool.bin, tool.args...)
		if strings.TrimSpace(out) == "" {
			continue
		}
		return strings.Split(strings.TrimRight(out, "\n"), "\n")
	}
	return nil
}
