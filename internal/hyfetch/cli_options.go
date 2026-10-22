package hyfetch

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Options are the flags hyfetch accepts in front of the neofetch backend.
type Options struct {
	Configure       bool
	ConfigFile      string
	Preset          string
	PresetSet       bool
	Mode            AnsiMode
	ModeSet         bool
	Backend         string
	BackendSet      bool
	Args            []string
	ArgsSet         bool
	Scale           float64
	ScaleSet        bool
	Lightness       float64
	LightnessSet    bool
	June            bool
	Debug           bool
	Distro          string
	DistroSet       bool
	AsciiFile       string
	AsciiFileSet    bool
	PrintFontLogo   bool
	TestPrint       bool
	AskExit         bool
	AutoDetect      bool
	AutoDetectSet   bool
	RandomAlignment bool
	Version         bool
}

// WantsColor reports whether the command line asked for flag coloring even
// when hyfetch.json is absent.
func (o Options) WantsColor() bool {
	return o.Configure || o.PresetSet || o.ModeSet || o.BackendSet || o.ScaleSet || o.LightnessSet || o.DistroSet || o.AsciiFileSet || o.RandomAlignment || o.AutoDetectSet || o.June
}

// HelpText is appended to gf's help.
const HelpText = `
Hyfetch-compatible options:
  -c, --config                  create or edit ~/.config/hyfetch.json
  -C, --config-file FILE        use another hyfetch.json
  -p, --preset PRESET           flag preset, comma-separated list, or #hex colors
  -m, --mode MODE               8bit or rgb
  -b, --backend BACKEND         neofetch (built-in), fastfetch, or macchina
      --args ARGS               extra arguments for the fetch backend
      --c-scale SCALE           multiply OKHSL lightness
      --c-set-l LIGHTNESS       set OKHSL lightness (0-1)
      --june                    play the pride month animation
      --distro, --test-distro   draw a specific distro logo
      --ascii-file FILE         use an ASCII file instead of the distro logo
      --print-font-logo         print the nerd-font logo and exit
      --auto-detect-light-dark  detect a light or dark terminal from COLORFGBG
  -r, --random-alignment        shuffle flag colors onto the logo slots
      --debug                   print the resolved preset on stderr
  -V                            print the version and exit
`

// SplitArgs pulls hyfetch flags out of argv. Everything else is left for the
// neofetch-style parser. --backend is only consumed when its value is a fetch
// backend; image backends such as "ascii" stay with gf.
func SplitArgs(args []string) (Options, []string, error) {
	var opt Options
	var rest []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		take := func() (string, error) {
			if i+1 >= len(args) {
				return "", fmt.Errorf("flag %s needs a value", a)
			}
			i++
			return args[i], nil
		}
		switch {
		case a == "-c":
			opt.Configure = true
		case a == "--config":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				opt.Configure = true
			} else {
				rest = append(rest, a, args[i+1])
				i++
			}
		case a == "-C" || a == "--config-file":
			v, err := take()
			if err != nil {
				return opt, nil, err
			}
			opt.ConfigFile = v
		case strings.HasPrefix(a, "--config-file="):
			opt.ConfigFile = strings.TrimPrefix(a, "--config-file=")
		case a == "-p" || a == "--preset":
			v, err := take()
			if err != nil {
				return opt, nil, err
			}
			opt.Preset, opt.PresetSet = v, true
		case strings.HasPrefix(a, "--preset="):
			opt.Preset, opt.PresetSet = strings.TrimPrefix(a, "--preset="), true
		case a == "-m" || a == "--mode":
			v, err := take()
			if err != nil {
				return opt, nil, err
			}
			mode, err := ParseMode(v)
			if err != nil {
				return opt, nil, err
			}
			opt.Mode, opt.ModeSet = mode, true
		case strings.HasPrefix(a, "--mode="):
			mode, err := ParseMode(strings.TrimPrefix(a, "--mode="))
			if err != nil {
				return opt, nil, err
			}
			opt.Mode, opt.ModeSet = mode, true
		case a == "-b" || a == "--backend":
			if i+1 < len(args) && isFetchBackend(args[i+1]) {
				i++
				opt.Backend, opt.BackendSet = strings.ToLower(args[i]), true
			} else if strings.HasPrefix(a, "--backend=") && isFetchBackend(strings.TrimPrefix(a, "--backend=")) {
				opt.Backend, opt.BackendSet = strings.ToLower(strings.TrimPrefix(a, "--backend=")), true
			} else {
				rest = append(rest, a)
			}
		case strings.HasPrefix(a, "--backend="):
			v := strings.TrimPrefix(a, "--backend=")
			if isFetchBackend(v) {
				opt.Backend, opt.BackendSet = strings.ToLower(v), true
			} else {
				rest = append(rest, a)
			}
		case a == "--args":
			v, err := take()
			if err != nil {
				return opt, nil, err
			}
			words, err := shellSplit(v)
			if err != nil {
				return opt, nil, fmt.Errorf("ARGS should be valid command-line arguments: %w", err)
			}
			opt.Args, opt.ArgsSet = words, true
		case strings.HasPrefix(a, "--args="):
			words, err := shellSplit(strings.TrimPrefix(a, "--args="))
			if err != nil {
				return opt, nil, fmt.Errorf("ARGS should be valid command-line arguments: %w", err)
			}
			opt.Args, opt.ArgsSet = words, true
		case a == "--c-scale":
			v, err := take()
			if err != nil {
				return opt, nil, err
			}
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return opt, nil, fmt.Errorf("SCALE should be a number")
			}
			opt.Scale, opt.ScaleSet = n, true
		case strings.HasPrefix(a, "--c-scale="):
			n, err := strconv.ParseFloat(strings.TrimPrefix(a, "--c-scale="), 64)
			if err != nil {
				return opt, nil, fmt.Errorf("SCALE should be a number")
			}
			opt.Scale, opt.ScaleSet = n, true
		case a == "--c-set-l":
			v, err := take()
			if err != nil {
				return opt, nil, err
			}
			n, err := ParseLightness(v)
			if err != nil {
				return opt, nil, err
			}
			opt.Lightness, opt.LightnessSet = n, true
		case strings.HasPrefix(a, "--c-set-l="):
			n, err := ParseLightness(strings.TrimPrefix(a, "--c-set-l="))
			if err != nil {
				return opt, nil, err
			}
			opt.Lightness, opt.LightnessSet = n, true
		case a == "--june":
			opt.June = true
		case a == "--debug":
			opt.Debug = true
		case a == "--distro" || a == "--test-distro":
			v, err := take()
			if err != nil {
				return opt, nil, err
			}
			opt.Distro, opt.DistroSet = v, true
		case strings.HasPrefix(a, "--distro="):
			opt.Distro, opt.DistroSet = strings.TrimPrefix(a, "--distro="), true
		case strings.HasPrefix(a, "--test-distro="):
			opt.Distro, opt.DistroSet = strings.TrimPrefix(a, "--test-distro="), true
		case a == "--ascii-file":
			v, err := take()
			if err != nil {
				return opt, nil, err
			}
			opt.AsciiFile, opt.AsciiFileSet = v, true
		case strings.HasPrefix(a, "--ascii-file="):
			opt.AsciiFile, opt.AsciiFileSet = strings.TrimPrefix(a, "--ascii-file="), true
		case a == "--print-font-logo":
			opt.PrintFontLogo = true
		case a == "--test-print":
			opt.TestPrint = true
		case a == "--ask-exit":
			opt.AskExit = true
		case a == "--auto-detect-light-dark":
			v, err := take()
			if err != nil {
				return opt, nil, err
			}
			on, err := parseBool(v)
			if err != nil {
				return opt, nil, err
			}
			opt.AutoDetect, opt.AutoDetectSet = on, true
		case strings.HasPrefix(a, "--auto-detect-light-dark="):
			on, err := parseBool(strings.TrimPrefix(a, "--auto-detect-light-dark="))
			if err != nil {
				return opt, nil, err
			}
			opt.AutoDetect, opt.AutoDetectSet = on, true
		case a == "-r" || a == "--random-alignment":
			opt.RandomAlignment = true
		case a == "-V":
			opt.Version = true
		default:
			rest = append(rest, a)
		}
	}
	return opt, rest, nil
}

func isFetchBackend(s string) bool {
	switch strings.ToLower(s) {
	case "neofetch", "fastfetch", "macchina":
		return true
	default:
		return false
	}
}

func parseBool(s string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "t", "true", "on", "yes":
		return true, nil
	case "0", "f", "false", "off", "no":
		return false, nil
	default:
		return false, fmt.Errorf("BOOL should be true or false")
	}
}

// shellSplit splits a POSIX-ish argument string. Hyfetch passes --args through
// the same way.
func shellSplit(s string) ([]string, error) {
	var out []string
	var buf strings.Builder
	var quote rune
	escape := false
	flush := func() {
		if buf.Len() > 0 || quote != 0 {
			out = append(out, buf.String())
			buf.Reset()
		}
	}
	for _, r := range s {
		switch {
		case escape:
			buf.WriteRune(r)
			escape = false
		case r == '\\' && quote != '\'':
			escape = true
		case quote != 0 && r == quote:
			quote = 0
			out = append(out, buf.String())
			buf.Reset()
		case quote != 0:
			buf.WriteRune(r)
		case r == '\'' || r == '"':
			quote = r
		case r == ' ' || r == '\t' || r == '\n':
			if buf.Len() > 0 {
				flush()
			}
		default:
			buf.WriteRune(r)
		}
	}
	if escape || quote != 0 {
		return nil, fmt.Errorf("unterminated string")
	}
	if buf.Len() > 0 {
		out = append(out, buf.String())
	}
	return out, nil
}

// Apply overlays command-line options onto a loaded config.
func (o Options) Apply(cfg Config) (Config, error) {
	if o.PresetSet {
		cfg.Preset = PresetValue{Values: []string{o.Preset}}
	}
	if o.ModeSet {
		cfg.Mode = o.Mode
	}
	if o.BackendSet {
		cfg.Backend = o.Backend
	}
	if o.ArgsSet {
		cfg.Args = ArgList(o.Args)
	}
	if o.DistroSet {
		cfg.Distro = o.Distro
	}
	if o.AsciiFileSet {
		cfg.CustomASCIIPath = o.AsciiFile
	}
	if o.AutoDetectSet {
		cfg.AutoDetectLightDark = &o.AutoDetect
	}
	if o.RandomAlignment {
		cfg.ColorAlign.Mode = "random"
	}
	if o.LightnessSet {
		cfg.Lightness = &o.Lightness
	}
	return cfg, cfg.normalize()
}

// Configure writes a hyfetch.json. A terminal gets the short preset/mode/alignment
// prompts; otherwise the defaults are saved so the command can continue.
func Configure(path string, in io.Reader, out io.Writer) (Config, error) {
	cfg := DefaultConfig()
	if !isTerminalFile(os.Stdin) || !isTerminalFile(os.Stdout) {
		if err := cfg.Save(path); err != nil {
			return cfg, err
		}
		fmt.Fprintf(os.Stderr, "Wrote %s\n", path)
		return cfg, nil
	}
	names := PresetNames(nil)
	fmt.Fprintln(out, Color("&l&bhyfetch&~&L - neofetch with flags <3", Mode8bit))
	fmt.Fprintln(out)
	fmt.Fprintf(out, "Presets: %s\n", strings.Join(names, ", "))
	preset := literalInput(out, in, "Which flag preset?", "rainbow")
	if _, _, err := ResolvePreset(preset, nil); err != nil {
		return cfg, err
	}
	cfg.Preset = PresetValue{Values: []string{preset}}
	mode := literalInput(out, in, "Color mode? (8bit/rgb)", "8bit")
	parsed, err := ParseMode(mode)
	if err != nil {
		return cfg, err
	}
	cfg.Mode = parsed
	align := literalInput(out, in, "Color alignment? (horizontal/vertical/random)", "horizontal")
	switch strings.ToLower(align) {
	case "horizontal", "vertical", "random", "custom":
		cfg.ColorAlign.Mode = strings.ToLower(align)
	default:
		return cfg, fmt.Errorf("color alignment should be horizontal, vertical, or random")
	}
	if err := cfg.Save(path); err != nil {
		return cfg, err
	}
	fmt.Fprintf(os.Stderr, "Saved %s\n", path)
	return cfg, nil
}

func literalInput(out io.Writer, in io.Reader, prompt, def string) string {
	fmt.Fprintf(out, "%s (default: %s)\n> ", prompt, def)
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return def
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}

func isTerminalFile(f *os.File) bool {
	if f == nil {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// ThemeFromEnv reads COLORFGBG. The last number is the background; 7 and 15
// are the usual light backgrounds.
func ThemeFromEnv() (Theme, bool) {
	raw := os.Getenv("COLORFGBG")
	if raw == "" {
		return "", false
	}
	parts := strings.FieldsFunc(raw, func(r rune) bool { return r == ';' || r == ':' })
	if len(parts) == 0 {
		return "", false
	}
	n, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return "", false
	}
	if n == 7 || n == 15 {
		return ThemeLight, true
	}
	return ThemeDark, true
}

// ResolveTheme picks the background hyfetch should adapt to.
func ResolveTheme(cfg Config, opt Options) Theme {
	auto := false
	if cfg.AutoDetectLightDark != nil {
		auto = *cfg.AutoDetectLightDark
	}
	if opt.AutoDetectSet {
		auto = opt.AutoDetect
	}
	if auto {
		if theme, ok := ThemeFromEnv(); ok {
			return theme
		}
	}
	if cfg.LightDark == ThemeLight {
		return ThemeLight
	}
	return ThemeDark
}
