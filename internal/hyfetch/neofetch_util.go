package hyfetch

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunFastfetch writes the recolored logo and execs fastfetch --file-raw.
func RunFastfetch(art string, args []string) error {
	path, err := writeTemp(art, "gf-ascii-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(path)
	cmdArgs := append([]string{"--file-raw", path}, args...)
	cmd := exec.Command("fastfetch", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if errorsIsNotFound(err) {
			return fmt.Errorf("fastfetch was not found")
		}
		return err
	}
	return nil
}

// RunMacchina writes a one-off theme that points at the recolored logo.
func RunMacchina(art string, args []string) error {
	asciiPath, err := writeTemp(art, "gf-ascii-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(asciiPath)
	theme := fmt.Sprintf("hide_ascii = false\nseparator = \"\"\n\n[custom_ascii]\npath = %q\n", asciiPath)
	themePath, err := writeTemp(theme, "gf-macchina-*.toml")
	if err != nil {
		return err
	}
	defer os.Remove(themePath)
	// macchina --theme takes the file stem when the theme lives in its theme
	// directory. Passing the file path works on current macchina builds.
	cmdArgs := append([]string{"--theme", themePath}, args...)
	cmd := exec.Command("macchina", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if errorsIsNotFound(err) {
			return fmt.Errorf("macchina was not found")
		}
		return err
	}
	return nil
}

func writeTemp(body, pattern string) (string, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	if _, err := f.WriteString(body); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	if err := f.Close(); err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

func errorsIsNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "executable file not found")
}

// Paint applies scale, lightness, and alignment and returns recolored lines.
func Paint(raw string, cfg Config, opt Options, theme Theme) ([]string, ColorProfile, string, error) {
	custom, err := cfg.CustomProfiles()
	if err != nil {
		return nil, ColorProfile{}, "", err
	}
	spec := cfg.Preset.Pick()
	if opt.PresetSet {
		spec = opt.Preset
	}
	profile, name, err := ResolvePreset(spec, custom)
	if err != nil {
		return nil, ColorProfile{}, "", err
	}
	switch {
	case opt.ScaleSet:
		profile = profile.Lighten(opt.Scale)
	case opt.LightnessSet:
		profile = profile.WithLightness(assignReplace, opt.Lightness)
	default:
		profile = profile.WithLightnessAdaptive(cfg.EffectiveLightness(theme), theme)
	}
	align := cfg.ColorAlign
	if opt.RandomAlignment {
		align.Mode = "random"
	}
	colored, err := Recolor(raw, align, profile, cfg.Mode, theme)
	if err != nil {
		return nil, profile, name, err
	}
	if colored == "" {
		return nil, profile, name, fmt.Errorf("ascii art is empty")
	}
	return strings.Split(colored, "\n"), profile, name, nil
}
