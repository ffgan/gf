package logo

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

var (
	sysLocale       = "C.UTF-8"
	asciiData       string
	asciiLen        int
	asciiLines      int
	c1, c2, c3      = "\033[31m", "\033[32m", "\033[33m"
	c4, c5, c6      = "\033[34m", "\033[35m", "\033[36m"
	gap             = 2
	Reset           = "\033[0m"
	asciiBold       = "\033[1m"
	colorText       = utils.ON
	barColorElapsed = "distro"
	barColorTotal   = "distro"
	colors          = []string{"distro"}

	asciiColorOverride []int
	lowerIndex         map[string]DistroWithColor
)

func init() {
	lowerIndex = make(map[string]DistroWithColor)
	if r == nil {
		return
	}
	r.Walk(func(s string, v interface{}) bool {
		if dc, ok := v.(DistroWithColor); ok {
			lowerIndex[strings.ToLower(s)] = dc
		}
		return false
	})
}

// SetAppearance applies bold and palette overrides before the logo is drawn.
// asciiColors and textColors use neofetch's "distro" keyword or a list of
// color indexes such as "1 7" or "(2 4)".
func SetAppearance(bold bool, asciiColors, textColors string) {
	if bold {
		asciiBold = "\033[1m"
	} else {
		asciiBold = ""
	}
	asciiColorOverride = parseColorInts(asciiColors)
	textColors = strings.Trim(textColors, "() \t")
	if textColors == "" || strings.EqualFold(textColors, "distro") {
		colors = []string{"distro"}
		return
	}
	colors = strings.Fields(strings.NewReplacer(",", " ", "(", " ", ")", " ").Replace(textColors))
}

// ApplyDistro loads the palette for a distro without printing artwork.
func ApplyDistro(name string) {
	dc := Resolve(name)
	setColors(dc.Colors...)
}

func parseColorInts(spec string) []int {
	spec = strings.Trim(spec, "() \t")
	if spec == "" || strings.EqualFold(spec, "distro") {
		return nil
	}
	var nums []int
	for _, field := range strings.Fields(strings.NewReplacer(",", " ", "(", " ", ")", " ").Replace(spec)) {
		n, err := strconv.Atoi(field)
		if err == nil {
			nums = append(nums, n)
		}
	}
	return nums
}

// idAliases maps os-release ID values onto logo names.
var idAliases = map[string]string{
	"ubuntu":              "Ubuntu",
	"debian":              "Debian",
	"arch":                "Arch",
	"archlinux":           "Arch",
	"linuxmint":           "Linux Mint",
	"mint":                "Linux Mint",
	"fedora":              "Fedora",
	"pop":                 "Pop!_OS",
	"popos":               "Pop!_OS",
	"opensuse-tumbleweed": "openSUSE Tumbleweed",
	"opensuse-leap":       "openSUSE Leap",
	"opensuse":            "openSUSE",
	"almalinux":           "AlmaLinux",
	"rocky":               "Rocky",
	"rhel":                "RHEL",
	"centos":              "CentOS",
	"kali":                "Kali",
	"elementary":          "Elementary",
	"zorin":               "Zorin",
	"manjaro":             "Manjaro",
	"endeavouros":         "EndeavourOS",
	"garuda":              "Garuda",
	"nixos":               "NixOS",
	"gentoo":              "Gentoo",
	"alpine":              "Alpine",
	"raspbian":            "Raspbian",
	"linux":               "Linux",
	"macos":               "Darwin",
	"darwin":              "Darwin",
	"cachyos":             "CachyOS",
	"nobara":              "Nobara",
	"void":                "Void",
	"artix":               "Artix",
	"arcolinux":           "ArcoLinux",
	"mx":                  "MX",
	"parrot":              "Parrot",
	"deepin":              "Deepin",
	"kde neon":            "KDE",
	"neon":                "KDE",
	"uos":                 "Uos",
	"uniontech":           "Uos",
	"kylin":               "Ubuntu Kylin",
	"openkylin":           "openKylin",
	"openeuler":           "openEuler",
	"euleros":             "openEuler",
	"harmonyos":           "Android",
	"openharmony":         "Android",
	"hyperos":             "Android",
	"anolis":              "CentOS",
	"opencloudos":         "CentOS",
	"loongnix":            "Debian",
	"openwrt":             "Wrt",
	"raspios":             "Raspbian",
	"raspberrypi":         "Raspbian",
	"raspberry pi os":     "Raspbian",
	"amzn":                "Amazon",
	"amazonlinux":         "Amazon",
	"ol":                  "Oracle",
	"oraclelinux":         "Oracle",
	"ubuntu-mate":         "Ubuntu MATE",
	"ubuntumate":          "Ubuntu MATE",
	"ubuntu-budgie":       "Ubuntu Budgie",
	"ubuntu-studio":       "Ubuntu Studio",
	"ubuntu-cinnamon":     "Ubuntu Cinnamon",
	"ubuntu-kylin":        "Ubuntu Kylin",
	"ubuntu-touch":        "Ubuntu Touch",
	"ubuntu-sway":         "Ubuntu Sway",
	"linuxlite":           "Linux Lite",
	"linux-lite":          "Linux Lite",
	"fedora-silverblue":   "Fedora Silverblue",
	"fedora-kinoite":      "Fedora Kinoite",
	"fedora-sericea":      "Fedora Sericea",
	"fedora-coreos":       "Fedora CoreOS",
	"clear-linux-os":      "Clear Linux OS",
	"cbl-mariner":         "CBL-Mariner",
	"mariner":             "CBL-Mariner",
	"opensuse-slowroll":   "openSUSE Tumbleweed-Slowroll",
	"sles":                "SUSE",
	"sled":                "SUSE",
	"aosc":                "AOSC OS",
	"archlabs":            "ARCHlabs",
	"crystal-linux":       "Crystal Linux",
	"rebornos":            "Reborn OS",
	"ultramarine":         "Ultramarine Linux",
	"sourcemage":          "Source Mage",
	"rhinolinux":          "Rhino Linux",
	"vanillaos":           "Vanilla",
	"parrotos":            "Parrot",
	"sailfishos":          "Sailfish",
	"chromeos":            "Chrom",
	"pmos":                "PostMarketOS",
	"mx-linux":            "MX",
	"midnightbsd":         "FreeBSD",
	"astra":               "Astra Linux",
	"miraclelinux":        "MIRACLE LINUX",
	"calculate":           "Calculate",
	"lfs":                 "LinuxFromScratch",
	"instantos":           "instantOS",
	"archlinuxarm":        "Arch",
	"archarm":             "Arch",
	"alarm":               "Arch",
	"manjaro-arm":         "Manjaro",
	"sles_sap":            "SUSE",
	"opensuse-microos":    "openSUSE",
	"guixsystem":          "Guix",
	"azurelinux":          "CBL-Mariner",
	"immortalwrt":         "Wrt",
	"librewrt":            "Wrt",
	"pidora":              "Fedora",
	"ghostbsd":            "GhostBSD",
	"hardenedbsd":         "HardenedBSD",
	"nomadbsd":            "NomadBSD",
	"opnsense":            "OPNsense",
	"omnios":              "OmniOS",
	"openindiana":         "OpenIndiana",
	"smartos":             "SmartOS",
}

// Resolve finds the logo record for a distro name, ID, or pretty name.
func Resolve(name string) DistroWithColor {
	name = strings.TrimSpace(name)
	if dc, ok := lookup(name); ok {
		return dc
	}
	if dc, ok := lookup("Linux"); ok {
		return dc
	}
	return DistroWithColor{FileName: "Linux", Colors: []int{8, 3}}
}

func lookup(name string) (DistroWithColor, bool) {
	name = strings.TrimSpace(name)
	if name == "" || r == nil {
		return DistroWithColor{}, false
	}
	if v, ok := r.Get(name); ok {
		return v.(DistroWithColor), true
	}
	if dc, ok := lowerIndex[strings.ToLower(name)]; ok {
		return dc, true
	}
	if alias, ok := idAliases[strings.ToLower(name)]; ok {
		if dc, ok := lookup(alias); ok {
			return dc, true
		}
	}
	if prefix, v, ok := r.LongestPrefix(name); ok && validPrefix(name, prefix) {
		return v.(DistroWithColor), true
	}
	fields := strings.Fields(name)
	for n := len(fields); n >= 1; n-- {
		cand := strings.Join(fields[:n], " ")
		if dc, ok := lowerIndex[strings.ToLower(cand)]; ok {
			return dc, true
		}
		if alias, ok := idAliases[strings.ToLower(cand)]; ok {
			if dc, ok := lowerIndex[strings.ToLower(alias)]; ok {
				return dc, true
			}
		}
	}
	if len(fields) > 0 {
		if alias, ok := idAliases[strings.ToLower(fields[0])]; ok {
			if dc, ok := lowerIndex[strings.ToLower(alias)]; ok {
				return dc, true
			}
		}
	}
	return DistroWithColor{}, false
}

func validPrefix(name, prefix string) bool {
	if len([]rune(prefix)) < 3 || len(prefix) > len(name) {
		return false
	}
	if len(name) == len(prefix) {
		return true
	}
	next := name[len(prefix)]
	return next == ' ' || next == '-' || next == '_' || next == '/' || next == '('
}

func loadDistroASCII(ASCIIDistro string, ASCIIFiles embed.FS) {
	dc := Resolve(ASCIIDistro)
	setColors(dc.Colors...)

	filename := "assets/" + dc.FileName
	f, err := ASCIIFiles.Open(filename)
	if err != nil {
		asciiData = fallbackLogo
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		asciiData = fallbackLogo
		return
	}
	asciiData = string(data)
}

const fallbackLogo = `${c1}    .--.
${c1}   |o_o |
${c1}   |:_/ |
${c1}  //   \ \
${c1} (|     | )
${c1}/'\_   _/` + "`\\" + `
${c1}\___)=(___/
`

func stripEscapeCodes(s string) string {
	re := regexp.MustCompile(`\x1B\[[0-9;]*[A-Za-z]`)
	return re.ReplaceAllString(s, "")
}

// RawASCII returns distro art with ${cN} placeholders still in place.
func RawASCII(ASCIIDistro, imageSource string, ASCIIFiles embed.FS) string {
	loadASCII(ASCIIDistro, imageSource, ASCIIFiles)
	return asciiData
}

func PrintASCII(ASCIIDistro, imageSource string, ASCIIFiles embed.FS) []string {
	raw := loadASCII(ASCIIDistro, imageSource, ASCIIFiles)
	lines := colorize(raw)
	if len(lines) <= 1 {
		loadDistroASCII("Linux", ASCIIFiles)
		lines = colorize(asciiData)
	}
	return lines
}

func loadASCII(ASCIIDistro, imageSource string, ASCIIFiles embed.FS) string {
	raw := ""
	if imageSource != "" && imageSource != "ascii" && imageSource != utils.AUTO {
		fi, err := os.Stat(imageSource)
		if err == nil && fi.Mode().IsRegular() {
			ext := strings.ToLower(filepath.Ext(imageSource))
			switch ext {
			case ".png", ".jpg", ".jpeg", ".jpe", ".svg", ".gif", ".webp":
				// Real images are rendered by an external backend when one exists.
			default:
				data, err := os.ReadFile(imageSource)
				if err == nil {
					raw = string(data)
				}
			}
		} else if !strings.Contains(imageSource, "\n") && len(imageSource) < 64 && !strings.Contains(imageSource, "/") {
			// A short token is a distro override, not literal art.
			ASCIIDistro = imageSource
		} else {
			raw = imageSource
		}
	}

	// Palette follows the distro even when the artwork comes from a file.
	loadDistroASCII(ASCIIDistro, ASCIIFiles)
	if strings.TrimSpace(raw) != "" {
		asciiData = raw
	}

	_ = sysLocale
	_ = asciiLen
	_ = asciiLines
	return asciiData
}

func colorize(data string) []string {
	replacer := strings.NewReplacer(
		"${c1}", c1, "${c2}", c2, "${c3}", c3,
		"${c4}", c4, "${c5}", c5, "${c6}", c6,
		"\\033", "\033",
	)
	data = replacer.Replace(data)
	data = strings.TrimRight(data, "\n")
	if strings.TrimSpace(stripEscapeCodes(data)) == "" {
		return nil
	}
	return propagateColor(splitKeep(data))
}

var colorCode = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// propagateColor carries an open color onto following lines. Distro art usually
// sets a color once and relies on it staying on, but each output row is reset
// so the info text does not inherit it.
func propagateColor(lines []string) []string {
	active := ""
	for i, line := range lines {
		if active != "" && !strings.HasPrefix(line, "\x1b") {
			line = active + line
			lines[i] = line
		}
		active = openColor(line)
	}
	return lines
}

func openColor(line string) string {
	active := ""
	for _, code := range colorCode.FindAllString(line, -1) {
		if code == "\x1b[0m" {
			active = ""
			continue
		}
		active += code
	}
	return active
}

func splitKeep(data string) []string {
	sc := bufio.NewScanner(strings.NewReader(data))
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func color(input string) string {
	if input == "" {
		return Reset
	}
	switch {
	case len(input) == 1 && input[0] >= '0' && input[0] <= '7':
		return fmt.Sprintf("%s\033[3%sm", Reset, input)
	case input == "fg":
		return Reset
	case strings.HasPrefix(input, "#"):
		rgb, err := strconv.ParseInt(strings.TrimPrefix(input, "#"), 16, 32)
		if err != nil {
			return Reset
		}
		r := (rgb >> 16) & 0xff
		g := (rgb >> 8) & 0xff
		b := rgb & 0xff
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	default:
		return fmt.Sprintf("\033[38;5;%sm", input)
	}
}

func setColors(argsInts ...int) {
	nums := argsInts
	if len(asciiColorOverride) > 0 {
		nums = asciiColorOverride
	}
	args := make([]string, len(nums))
	for i, n := range nums {
		args[i] = strconv.Itoa(n)
	}

	slots := []*string{&c1, &c2, &c3, &c4, &c5, &c6}
	for i, slot := range slots {
		if i < len(args) {
			*slot = color(args[i])
			if asciiBold != "" {
				*slot += asciiBold
			}
		} else {
			*slot = ""
		}
	}

	if colorText != utils.OFF {
		setTextColors(args...)
	}
}

var (
	titleColor, atColor, underlineColor  string
	subtitleColor, colonColor, infoColor string
)

func setTextColors(args ...string) {
	if len(colors) == 0 || colors[0] == "distro" {
		titleColor = color(firstColor(args, 0, "7"))
		if len(args) > 1 && args[1] != "" {
			subtitleColor = color(args[1])
		} else {
			subtitleColor = titleColor
		}
		atColor = Reset
		underlineColor = titleColor
		colonColor = Reset
		infoColor = Reset

		if getArgInt(args, 0, 1) == 8 {
			titleColor = Reset
		}
		if getArgInt(args, 1, 7) == 8 {
			subtitleColor = Reset
		}
		if len(args) > 1 && getArgInt(args, 1, 7) == 7 {
			subtitleColor = color(getArg(args, 0))
		}
		if getArgInt(args, 0, 1) == 7 {
			titleColor = Reset
		}
	} else {
		titleColor = color(getColor(0))
		atColor = color(getColor(1))
		underlineColor = color(getColor(2))
		subtitleColor = color(getColor(3))
		colonColor = color(getColor(4))
		infoColor = color(getColor(5))
	}

	if barColorElapsed == "distro" || barColorElapsed == "" {
		barColorElapsed = color("fg")
	} else if !strings.Contains(barColorElapsed, "\033") {
		barColorElapsed = color(barColorElapsed)
	}

	if barColorTotal == "distro" || barColorTotal == "" {
		barColorTotal = color("fg")
	} else if !strings.Contains(barColorTotal, "\033") {
		barColorTotal = color(barColorTotal)
	}
}

func firstColor(args []string, index int, def string) string {
	if v := getArg(args, index); v != "" {
		return v
	}
	return def
}

func getArg(args []string, index int) string {
	if index < len(args) {
		return args[index]
	}
	return ""
}

func getArgInt(args []string, index int, def int) int {
	if index >= len(args) {
		return def
	}
	n, err := strconv.Atoi(args[index])
	if err != nil {
		return def
	}
	return n
}

func getColor(i int) string {
	if i < len(colors) {
		return colors[i]
	}
	return "7"
}

func TitleColor() string {
	if titleColor == "" {
		return "\033[1m"
	}
	return titleColor
}

func SubtitleColor() string {
	if subtitleColor == "" {
		return TitleColor()
	}
	return subtitleColor
}

func UnderlineColor() string {
	if underlineColor == "" {
		return TitleColor()
	}
	return underlineColor
}

func AtColor() string {
	if atColor == "" {
		return Reset
	}
	return atColor
}

func ResetCode() string { return Reset }
