package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bufio"
	"time"

	"github.com/ffgan/gf/configs"
	cli "github.com/ffgan/gf/internal/CLI"
	"github.com/ffgan/gf/internal/hyfetch"
	"github.com/ffgan/gf/internal/logo"
	"github.com/ffgan/gf/internal/utils"
)

const neofetchConf = "/etc/neofetch/config.conf"

func Run(ASCIIFiles embed.FS) {
	hyOpts, rest, err := hyfetch.SplitArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "gf: %v\n", err)
		os.Exit(2)
	}
	if hyOpts.Version {
		fmt.Printf("gf %s\n", Version)
		os.Exit(0)
	}
	os.Args = append([]string{os.Args[0]}, rest...)

	kernelName, kernelVersion, kernelMachine, darwinName := cacheUname()
	osName := cli.GetOS(kernelName, darwinName)

	hyPath := hyOpts.ConfigFile
	if hyPath == "" {
		hyPath = hyfetch.ConfigPath()
	}
	var hyCfg hyfetch.Config
	hyExists := false
	if hyOpts.Configure {
		hyCfg, err = hyfetch.Configure(hyPath, os.Stdin, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gf: %v\n", err)
			os.Exit(2)
		}
		hyExists = true
	} else {
		hyCfg, hyExists, err = hyfetch.Load(hyPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gf: %v\n", err)
			os.Exit(1)
		}
	}
	hyCfg, err = hyOpts.Apply(hyCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gf: %v\n", err)
		os.Exit(2)
	}
	colorOn := hyExists || hyOpts.WantsColor()
	if colorOn && (hyCfg.Backend == "" || hyCfg.Backend == "neofetch") && len(hyCfg.Args) > 0 {
		os.Args = append(os.Args, hyCfg.Args...)
	}

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
	if colorOn && hyCfg.Distro != "" {
		h.ASCIIDistro = hyCfg.Distro
	}
	h.Get_distro()

	if hyOpts.PrintFontLogo {
		glyph, err := hyfetch.FontLogo(h.distro)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gf: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(glyph)
		return
	}
	if hyOpts.TestPrint {
		fmt.Println(logo.RawASCII(h.ascii_distro, h.ImageSource, ASCIIFiles))
		return
	}

	theme := hyfetch.ResolveTheme(hyCfg, hyOpts)
	now := time.Now()
	if hyfetch.ShouldPride(hyCfg, hyOpts.June, now) {
		mode := hyCfg.Mode
		if mode == "" {
			mode = hyfetch.Mode8bit
		}
		hyfetch.Pride(mode)
		if !hyOpts.June {
			hyfetch.MarkPrideShown(now)
		}
	}

	var art []string
	if colorOn && !h.plainMode() && !strings.EqualFold(h.ImageBackend, "off") {
		art = hyfetchArt(h, hyCfg, hyOpts, theme, ASCIIFiles)
	} else if !h.plainMode() && !strings.EqualFold(h.ImageBackend, "off") {
		if lines := imageArt(h.ImageSource); len(lines) > 0 {
			art = lines
			logo.ApplyDistro(h.ascii_distro)
		} else {
			art = logo.PrintASCII(h.ascii_distro, h.ImageSource, ASCIIFiles)
		}
	} else {
		logo.ApplyDistro(h.ascii_distro)
	}

	if jsonMode {
		h.writeJSON(os.Stdout)
		return
	}
	show(art, h.Rows(), h.gap())
	if hyOpts.AskExit {
		fmt.Fprint(os.Stderr, "Press enter to exit...")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	}
}

func hyfetchArt(h *hardware, cfg hyfetch.Config, opt hyfetch.Options, theme hyfetch.Theme, files embed.FS) []string {
	source := h.ImageSource
	if cfg.CustomASCIIPath != "" {
		source = cfg.CustomASCIIPath
	}
	raw := logo.RawASCII(h.ascii_distro, source, files)
	lines, profile, name, err := hyfetch.Paint(raw, cfg, opt, theme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gf: %v\n", err)
		return logo.PrintASCII(h.ascii_distro, h.ImageSource, files)
	}
	if opt.Debug {
		fmt.Fprintf(os.Stderr, "gf: preset %s, mode %s, align %s, colors %d\n", name, cfg.Mode, cfg.ColorAlign.Mode, len(profile.Colors))
	}
	art := strings.Join(lines, "\n")
	switch cfg.Backend {
	case "fastfetch":
		if err := hyfetch.RunFastfetch(art, []string(cfg.Args)); err != nil {
			fmt.Fprintf(os.Stderr, "gf: %v\n", err)
			break
		}
		os.Exit(0)
	case "macchina":
		if err := hyfetch.RunMacchina(art, []string(cfg.Args)); err != nil {
			fmt.Fprintf(os.Stderr, "gf: %v\n", err)
			break
		}
		os.Exit(0)
	}
	return lines
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
