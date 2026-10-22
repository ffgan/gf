package cli

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

var versionRE = regexp.MustCompile(`\d+(?:\.\d+)*`)

// earlyDistro applies the release sources neofetch consults before os-release.
// Those files identify older and embedded distributions that also ship a
// generic os-release.
func earlyDistro(shorthand string, files map[string]string) (display, ascii, name, id, idLike string, ok bool) {
	if len(files) == 0 {
		return
	}
	if strings.EqualFold(files["kernel-name"], "HarmonyOS") {
		name, id = "HarmonyOS", "harmonyos"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			if ver := strings.TrimSpace(files["harmony-version"]); ver != "" {
				display = name + " " + ver
			}
		}
		return display, "HarmonyOS", name, id, "", true
	}
	if text := strings.TrimSpace(files["bedrock"]); text != "" && files["bedrock-restrict"] == "" {
		name, id = "Bedrock Linux", "bedrock"
		display = text
		if shorthand == "on" || shorthand == "tiny" {
			display = name
		}
		return display, "Bedrock", name, id, "", true
	}
	if text := strings.TrimSpace(files["redstar-release"]); text != "" {
		name, id = "Red Star OS", "redstar"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			if ver := firstVersion(text); ver != "" {
				display = name + " " + ver
			}
		}
		return display, "Red Star", name, id, "", true
	}
	if text := strings.TrimSpace(files["krass-version"]); text != "" {
		name, id = "KrassOS", "krassos"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			if ver := releaseWord(text, "release", 1); ver != "" {
				display = name + " " + ver
			}
		}
		return display, "KrassOS", name, id, "", true
	}
	if text := strings.TrimSpace(files["armbian-release"]); text != "" {
		name, id = "Armbian", "armbian"
		return armbianName(text), "Armbian", name, id, "", true
	}
	if text := strings.TrimSpace(files["manjaro-release"]); text != "" {
		name, id = "Manjaro Linux", "manjaro"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			if src := strings.TrimSpace(files["lsb-src"]); src != "" {
				display = name + " " + src
			} else if ver := firstVersion(text); ver != "" {
				display = name + " " + ver
			}
		}
		return display, "Manjaro", name, id, "", true
	}
	if text := strings.TrimSpace(files["siduction-version"]); text != "" {
		name, id = "Siduction", "siduction"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			detail := strings.TrimSpace(files["lsb-sic"])
			if detail == "" {
				detail = strings.TrimSpace(strings.Split(text, "\n")[0])
			}
			if detail != "" {
				display = name + " (" + detail + ")"
			}
		}
		return display, "Siduction", name, id, "", true
	}
	if text := strings.TrimSpace(files["q4os-version"]); text != "" {
		name, id = "Q4OS", "q4os"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			display = strings.TrimSpace(name + " " + strings.TrimSpace(strings.Split(text, "\n")[0]))
		}
		return display, "Q4OS", name, id, "", true
	}
	if text := strings.TrimSpace(files["mcst-version"]); text != "" {
		name, id = "OS Elbrus", "elbrus"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			display = strings.TrimSpace(name + " " + strings.TrimSpace(strings.Split(text, "\n")[0]))
		}
		return display, "OS Elbrus", name, id, "", true
	}
	if _, marked := files["nixos-marker"]; marked {
		info := parseOSReleaseText(files["os-release"])
		name = info["NAME"]
		if name == "" {
			name = "NixOS"
		}
		id = "nixos"
		switch shorthand {
		case "tiny":
			display = name
		case "on":
			display = strings.TrimSpace(name + " " + info["VERSION"])
		default:
			ver := strings.TrimSpace(files["nixos-version"])
			if ver == "" {
				ver = info["VERSION"]
			}
			display = strings.TrimSpace(name + " " + ver)
		}
		return display, "NixOS", name, id, "", true
	}
	if text := strings.TrimSpace(files["linspire-version"]); text != "" {
		name, id = "Linspire", "linspire"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			if ver := releaseWord(text, "version", 1); ver != "" {
				display = name + " " + ver
			}
		}
		return display, "Linspire", name, id, "", true
	}
	if text := strings.TrimSpace(files["lindows-version"]); text != "" {
		name, id = "LindowsOS", "lindows"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			ver := strings.TrimSpace(releaseWord(text, "version", 1) + " " + releaseWord(text, "version", 2))
			if ver != "" {
				display = strings.TrimSpace(name + " " + ver)
			}
		}
		return display, "Lindows", name, id, "", true
	}
	if text := strings.TrimSpace(files["astra-version"]); text != "" {
		name, id = "Astra Linux", "astra"
		ver := firstVersion(text)
		code := "ALSE"
		if strings.HasPrefix(ver, "2.12") {
			code = "ALCE"
		}
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			display = strings.TrimSpace(name + " " + ver + " " + code)
		}
		return display, display, name, id, "", true
	}
	if text := strings.TrimSpace(files["pveversion"]); text != "" {
		name, id = "Proxmox VE", "proxmox"
		return proxmoxName(text, shorthand), "Proxmox", name, id, "", true
	}
	if text := strings.TrimSpace(files["twistver"]); text != "" {
		display = twistName(text, shorthand, files)
		return display, distroArtName(display, "", "", display), display, "", "", true
	}
	if files["anise"] != "" {
		info := parseOSReleaseText(files["os-release"])
		name = firstNonEmpty(info["NAME"], info["DISTRIB_DESCRIPTION"], "MacaroniOS")
		ver := firstNonEmpty(info["VERSION"], info["DISTRIB_RELEASE"])
		id = "macaronios"
		switch shorthand {
		case "on", "tiny":
			display = name
		default:
			display = strings.TrimSpace(name + " " + ver)
		}
		return display, "MacaroniOS", name, id, "", true
	}
	if si, sd := strings.TrimSpace(files["lsb-si"]), strings.TrimSpace(files["lsb-sd"]); si != "" || sd != "" {
		low := strings.ToLower(si)
		info := parseOSReleaseText(files["os-release"])
		if strings.Contains(low, "kslinux") && files["os-release"] != "" {
			name = firstNonEmpty(info["NAME"], "KSLinux")
			display = name
			if shorthand != "on" && shorthand != "tiny" {
				display = strings.TrimSpace(name + " " + info["VERSION"])
			}
			return display, "KSLinux", name, "kslinux", "", true
		}
		if strings.TrimSpace(files["debian_version"]) != "" && files["os-release"] != "" &&
			!strings.Contains(low, "buntu") && !strings.Contains(low, "neon") {
			name = firstNonEmpty(info["NAME"], si, "Debian")
			id = firstNonEmpty(info["ID"], strings.ToLower(si))
			ver := strings.TrimSpace(files["debian_version"])
			switch shorthand {
			case "on":
				display = name
			case "tiny":
				display = firstNonEmpty(si, name)
			default:
				display = strings.TrimSpace(name + " " + ver)
				if code := info["VERSION_CODENAME"]; code != "" {
					display += " (" + code + ")"
				}
			}
			return display, distroArtName(name, id, info["ID_LIKE"], display), name, id, info["ID_LIKE"], true
		}
		name = firstNonEmpty(si, sd)
		switch shorthand {
		case "on", "tiny":
			display = name
		default:
			display = firstNonEmpty(sd, si)
		}
		return display, distroArtName(name, strings.ToLower(si), "", display), name, strings.ToLower(si), "", true
	}
	return
}

// lateDistro covers command output and marker files neofetch reads only after
// os-release is missing.
func lateDistro(osName, shorthand string, files map[string]string) (display, ascii, name, id string, ok bool) {
	if len(files) == 0 {
		return
	}
	if text := strings.TrimSpace(files["crux"]); text != "" {
		display = text
		switch shorthand {
		case "on":
			display = strings.ReplaceAll(display, "version", "")
		case "tiny":
			if i := strings.Index(strings.ToLower(display), "version"); i >= 0 {
				display = display[:i]
			}
		}
		return strings.Join(strings.Fields(display), " "), "CRUX", "CRUX", "crux", true
	}
	if text := strings.TrimSpace(files["slitaz-release"]); text != "" {
		name, id = "SliTaz", "slitaz"
		return strings.TrimSpace(name + " " + strings.TrimSpace(strings.Split(text, "\n")[0])), "SliTaz", name, id, true
	}
	if text := strings.TrimSpace(files["android-release"]); text != "" {
		return strings.TrimSpace("Android " + text), "Android", "Android", "android", true
	}
	if strings.TrimSpace(files["vzlinux-release"]) != "" {
		return "VzLinux", "VzLinux", "VzLinux", "vzlinux", true
	}
	if ver := strings.TrimSpace(files["guix-version"]); ver != "" {
		name = "Guix System"
		display = name
		if shorthand != "on" && shorthand != "tiny" {
			display = strings.TrimSpace(name + " " + ver)
		}
		return display, "Guix", name, "guix", true
	}
	if text := strings.TrimSpace(files["openbsd-version"]); text != "" && (osName == OpenBSD || osName == "" || osName == BSD) {
		fields := strings.Fields(text)
		display = text
		if len(fields) >= 2 {
			display = fields[0] + " " + fields[1]
		}
		return display, "OpenBSD", "OpenBSD", "openbsd", true
	}
	if text := files["synology-version"]; strings.TrimSpace(text) != "" {
		info := parseOSReleaseText(text)
		if pv := info["productversion"]; pv != "" {
			display = "Synology DSM " + pv
			if build := info["buildnumber"]; build != "" {
				display += "-" + build
			}
			if sf := info["smallfixnumber"]; sf != "" && sf != "0" {
				display += " Update " + sf
			}
			return display, "Synology", "Synology DSM", "synology", true
		}
	}
	if files["peropesis"] != "" {
		return "Peropesis", "Peropesis", "Peropesis", "peropesis", true
	}
	if text := strings.TrimSpace(files["opnsense-version"]); text != "" {
		name, id = "OPNsense", "opnsense"
		if shorthand == "on" || shorthand == "tiny" {
			return name, "OPNsense", name, id, true
		}
		ver := text
		if i := strings.Index(ver, "-"); i >= 0 {
			ver = ver[:i]
		}
		return strings.TrimSpace(name + " " + ver), "OPNsense", name, id, true
	}
	if props, ver := files["qnx-props"], strings.TrimSpace(files["qnx-version"]); props != "" || ver != "" {
		return qnxName(props, ver), "QNX", "QNX", "qnx", true
	}
	if files["pcbsd"] != "" {
		return "PCBSD", "PCBSD", "PCBSD", "pcbsd", true
	}
	if files["trueos"] != "" {
		return "TrueOS", "TrueOS", "TrueOS", "trueos", true
	}
	if files["pacbsd"] != "" {
		return "PacBSD", "PacBSD", "PacBSD", "pacbsd", true
	}
	if files["hardenedbsd"] != "" {
		return "HardenedBSD", "HardenedBSD", "HardenedBSD", "hardenedbsd", true
	}
	return
}

// finishDistro applies the rewrites neofetch does after a name is chosen:
// Ubuntu flavors, Fedora immutable variants, WSL and Chrome OS suffixes.
func finishDistro(display, ascii, name, shorthand string, files map[string]string) (string, string) {
	if files == nil {
		files = map[string]string{}
	}
	display = strings.TrimSpace(strings.TrimPrefix(display, "NAME="))
	display = strings.Join(strings.Fields(strings.ReplaceAll(display, "Enterprise Server", "")), " ")
	if strings.HasPrefix(display, "Ubuntu") {
		if strings.TrimSpace(files["floflis-release"]) != "" {
			return "Floflis", "Floflis"
		}
		if flav := ubuntuFlavor(files["xdg-config-dirs"]); flav != "" && !strings.Contains(display, flav) {
			display = strings.Replace(display, "Ubuntu", flav, 1)
			ascii = flav
		}
	}
	if strings.HasPrefix(display, "Fedora") || strings.EqualFold(name, "Fedora") {
		if variant, rewritten, art, ok := fedoraVariant(display, shorthand, files); ok {
			if rewritten != "" {
				display = rewritten
			}
			ascii = art
			_ = variant
		}
	}
	if strings.HasPrefix(display, "FreeBSD") && files["nomadbsd"] != "" {
		return "NomadBSD", "NomadBSD"
	}
	proc := strings.ToLower(files["proc-version"])
	if strings.Contains(proc, "microsoft") {
		win := strings.TrimSpace(files["windows-product"])
		if win != "" && !strings.Contains(display, "Windows") {
			switch shorthand {
			case "on":
				display += " [" + win + "]"
			case "tiny":
				display = win
				ascii = windowsArt(win)
			default:
				display += " on " + win
			}
		}
	}
	if strings.Contains(proc, "chrome-bot") || files["cros-ec"] != "" {
		if !strings.Contains(strings.ToLower(display), "chrome") {
			ver := strings.TrimSpace(files["chrome-milestone"])
			switch shorthand {
			case "tiny":
				display = "Chrome OS"
			case "on":
				display = strings.TrimSpace("Chrome OS " + ver)
			default:
				display = strings.TrimSpace("Chrome OS " + ver + " " + display)
			}
			ascii = "Chrom"
		}
	}
	return strings.TrimSpace(display), ascii
}

func ubuntuFlavor(xdg string) string {
	low := strings.ToLower(xdg)
	switch {
	case strings.Contains(low, "studio"):
		return "Ubuntu Studio"
	case strings.Contains(low, "kde"), strings.Contains(low, "plasma"):
		return "Kubuntu"
	case strings.Contains(low, "mate"):
		return "Ubuntu MATE"
	case strings.Contains(low, "xubuntu"):
		return "Xubuntu"
	case strings.Contains(low, "lubuntu"):
		return "Lubuntu"
	case strings.Contains(low, "budgie"):
		return "Ubuntu Budgie"
	case strings.Contains(low, "cinnamon"):
		return "Ubuntu Cinnamon"
	case strings.Contains(low, "ubuntusway"):
		return "Ubuntu Sway"
	case strings.Contains(low, "ukui"):
		return "Ubuntu Kylin"
	default:
		return ""
	}
}

func fedoraVariant(display, shorthand string, files map[string]string) (variant, rewritten, art string, ok bool) {
	info := parseOSReleaseText(files["os-release"])
	variant = info["VARIANT"]
	switch variant {
	case "CoreOS", "Kinoite", "Sericea", "Silverblue":
	default:
		return "", "", "", false
	}
	art = "Fedora " + variant
	switch shorthand {
	case "tiny":
		rewritten = variant
	case "on":
		rewritten = strings.TrimSpace(art + " " + info["VERSION_ID"])
	default:
		if strings.Contains(display, variant) {
			rewritten = ""
		} else {
			ostree := firstNonEmpty(info["OSTREE_VERSION"], files["ostree-version"], info["VERSION_ID"])
			rewritten = strings.TrimSpace(art + " " + ostree)
		}
	}
	return variant, rewritten, art, true
}

func armbianName(text string) string {
	info := parseOSReleaseText(text)
	code := info["DISTRIBUTION_CODENAME"]
	ver := info["VERSION"]
	if code == "" && ver == "" {
		if pretty := info["PRETTY_NAME"]; pretty != "" {
			return pretty
		}
		return "Armbian"
	}
	display := strings.TrimSpace("Armbian " + code)
	if ver != "" {
		display += " (" + ver + ")"
	}
	return display
}

func proxmoxName(out, shorthand string) string {
	if shorthand == "on" || shorthand == "tiny" {
		return "Proxmox VE"
	}
	out = strings.TrimSpace(strings.TrimPrefix(out, "pve-manager/"))
	if i := strings.Index(out, "/"); i >= 0 {
		out = out[:i]
	}
	return strings.TrimSpace("Proxmox VE " + out)
}

func twistName(base, shorthand string, files map[string]string) string {
	base = strings.TrimSpace(base)
	if arm := strings.TrimSpace(files["armbian-release"]); arm != "" {
		info := parseOSReleaseText(arm)
		switch shorthand {
		case "on":
			return strings.TrimSpace(base + " (Armbian " + info["VERSION"] + ")")
		case "tiny":
			return base + " (Armbian)"
		default:
			return strings.TrimSpace(base + " (Armbian " + info["VERSION"] + " " + info["DISTRIBUTION_CODENAME"] + ")")
		}
	}
	info := parseOSReleaseText(files["os-release"])
	switch shorthand {
	case "on":
		return strings.TrimSpace(base + " (" + info["NAME"] + " " + info["VERSION_ID"] + ")")
	case "tiny":
		return strings.TrimSpace(base + " (" + firstNonEmpty(info["NAME"], "Linux") + ")")
	default:
		return strings.TrimSpace(base + " (" + firstNonEmpty(info["PRETTY_NAME"], info["NAME"]) + ")")
	}
}

func qnxName(props, version string) string {
	name := ""
	for _, line := range strings.Split(props, "\n") {
		key, val, ok := strings.Cut(line, "::")
		if ok && strings.TrimSpace(key) == "device_os" {
			name = strings.TrimSpace(val)
		}
	}
	if name == "" {
		name = "QNX"
	}
	if version != "" {
		return name + " (" + version + ")"
	}
	return name
}

func solarisDisplay(line, shorthand string) string {
	line = strings.TrimSpace(strings.Split(line, "\n")[0])
	if i := strings.Index(line, "("); i >= 0 {
		line = strings.TrimSpace(line[:i])
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return line
	}
	switch shorthand {
	case "on", "tiny":
		if len(fields) >= 3 {
			return fields[0] + " " + fields[2]
		}
		return strings.Join(fields, " ")
	default:
		if len(fields) >= 3 {
			return strings.Join(fields[:3], " ")
		}
		return strings.Join(fields, " ")
	}
}

func firstVersion(s string) string {
	return versionRE.FindString(s)
}

func releaseWord(text, needle string, index int) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(strings.ToLower(line), strings.ToLower(needle)) {
			fields := strings.Fields(line)
			if index < len(fields) {
				return fields[index]
			}
		}
	}
	fields := strings.Fields(text)
	if index < len(fields) && len(fields) > 1 {
		return fields[index]
	}
	if index == 0 {
		return strings.TrimSpace(text)
	}
	return ""
}

func formatApple(product, version, build, shorthand string) (display, ascii string, hideArch bool) {
	low := strings.ToLower(strings.TrimSpace(product))
	version = strings.TrimSpace(version)
	build = strings.TrimSpace(build)
	if strings.Contains(low, "iphone") || strings.Contains(low, "ipad") || low == "ios" {
		return strings.TrimSpace("iOS " + version), "iOS", true
	}
	if product == "" || strings.Contains(low, "mac") {
		return formatMacOS(version, build, shorthand), Darwin, false
	}
	return strings.TrimSpace(product + " " + version), product, false
}

func formatMacOS(version, build, shorthand string) string {
	version = strings.TrimSpace(version)
	build = strings.TrimSpace(build)
	codename := macOSCodename(version)
	switch shorthand {
	case "tiny":
		// neofetch rewrites 10.x and 11 to the family name and leaves later
		// codenames in place.
		if strings.HasPrefix(version, "10.") || strings.HasPrefix(version, "11.") {
			return strings.TrimSpace(macFamily(version) + " " + version)
		}
		return strings.TrimSpace(macOSCodename(version) + " " + version)
	case "on":
		return strings.TrimSpace(codename + " " + version)
	default:
		display := strings.TrimSpace(codename + " " + version)
		if build != "" {
			display += " (" + build + ")"
		}
		return display
	}
}

func macFamily(version string) string {
	switch {
	case strings.HasPrefix(version, "10.0"), strings.HasPrefix(version, "10.1"),
		strings.HasPrefix(version, "10.2"), strings.HasPrefix(version, "10.3"),
		strings.HasPrefix(version, "10.4"), strings.HasPrefix(version, "10.5"),
		strings.HasPrefix(version, "10.6"), strings.HasPrefix(version, "10.7"):
		// 10.10 and later also have the "10.1" prefix, so exclude them first.
		if strings.HasPrefix(version, "10.10") || strings.HasPrefix(version, "10.11") ||
			strings.HasPrefix(version, "10.12") || strings.HasPrefix(version, "10.13") ||
			strings.HasPrefix(version, "10.14") || strings.HasPrefix(version, "10.15") ||
			strings.HasPrefix(version, "10.16") {
			break
		}
		return "Mac OS X"
	}
	switch {
	case strings.HasPrefix(version, "10.8"), strings.HasPrefix(version, "10.9"),
		strings.HasPrefix(version, "10.10"), strings.HasPrefix(version, "10.11"):
		return "OS X"
	default:
		return "macOS"
	}
}

func windowsProduct(product, build string) string {
	product = strings.TrimSpace(product)
	product = strings.TrimPrefix(product, "Microsoft ")
	if strings.Contains(product, "Windows 10") && buildAtLeast(build, 22000) {
		product = strings.Replace(product, "Windows 10", "Windows 11", 1)
	}
	return strings.TrimSpace(product)
}

func windowsArt(product string) string {
	low := strings.ToLower(product)
	switch {
	case strings.Contains(low, "windows 11"), low == "windows11":
		return "Windows 11"
	case strings.Contains(low, "windows 10"), strings.Contains(low, "windows 8"), low == "windows10", low == "windows8":
		return "Windows 10"
	case strings.Contains(low, "windows 95"), strings.Contains(low, "windows95"), strings.Contains(low, "win95"):
		return "Windows95"
	default:
		return "Windows"
	}
}

func buildAtLeast(build string, min int) bool {
	n := 0
	seen := false
	for _, r := range strings.TrimSpace(build) {
		if r < '0' || r > '9' {
			break
		}
		seen = true
		n = n*10 + int(r-'0')
	}
	return seen && n >= min
}

// enrichRelease adds command output and marker files that neofetch probes
// while naming a distribution. FormatDistro stays pure; this only runs on
// the real machine.
func enrichRelease(files map[string]string) {
	if files == nil {
		return
	}
	if os.Getenv("BEDROCK_RESTRICT") != "" {
		files["bedrock-restrict"] = "1"
	}
	if xdg := os.Getenv("XDG_CONFIG_DIRS"); xdg != "" {
		files["xdg-config-dirs"] = xdg
	}
	if v := readSmall("/proc/version"); v != "" && files["proc-version"] == "" {
		files["proc-version"] = v
	}
	if utils.PathExists("/dev/cros_ec") {
		files["cros-ec"] = "1"
	}
	if utils.CommandExists("lsb_release") {
		if files["lsb-si"] == "" {
			files["lsb-si"] = utils.RunCommand("lsb_release", "-si")
		}
		if files["lsb-sd"] == "" {
			files["lsb-sd"] = utils.RunCommand("lsb_release", "-sd")
		}
		if files["lsb-sic"] == "" {
			files["lsb-sic"] = utils.RunCommand("lsb_release", "-sic")
		}
		if files["lsb-src"] == "" {
			files["lsb-src"] = utils.RunCommand("lsb_release", "-src")
		}
	}
	if utils.CommandExists("pveversion") && files["pveversion"] == "" {
		files["pveversion"] = utils.RunCommand("pveversion")
	}
	if utils.CommandExists("nixos-version") && files["nixos-version"] == "" {
		files["nixos-version"] = utils.RunCommand("nixos-version")
	}
	if utils.CommandExists("crux") && files["crux"] == "" {
		files["crux"] = utils.RunCommand("crux")
	}
	if utils.CommandExists("guix") && files["guix-version"] == "" {
		fields := strings.Fields(utils.RunCommand("guix", "-V"))
		if len(fields) > 0 {
			files["guix-version"] = fields[len(fields)-1]
		}
	}
	if utils.CommandExists("anise") {
		files["anise"] = "1"
	}
	if utils.PathExists("/usr/bin/nomadbsd-sysinfo") {
		files["nomadbsd"] = "1"
	}
	if utils.FileExists("/etc/pcbsd-lang") {
		files["pcbsd"] = "1"
	}
	if utils.FileExists("/etc/trueos-lang") {
		files["trueos"] = "1"
	}
	if utils.FileExists("/etc/pacbsd-release") {
		files["pacbsd"] = "1"
	}
	if utils.FileExists("/etc/hbsd-update.conf") {
		files["hardenedbsd"] = "1"
	}
	if matches, _ := filepath.Glob("/usr/bin/*peropesis*"); len(matches) > 0 {
		files["peropesis"] = "1"
	}
	kernel := utils.RunCommand("uname", "-s")
	if strings.EqualFold(kernel, "HarmonyOS") {
		files["kernel-name"] = "HarmonyOS"
		if utils.CommandExists("param") {
			ver := strings.TrimPrefix(utils.RunCommand("param", "get", "const.ohos.fullname"), "OpenHarmony-")
			if strings.TrimSpace(ver) == "" {
				ver = utils.RunCommand("param", "get", "const.product.software.version")
			}
			files["harmony-version"] = strings.TrimSpace(ver)
		}
	}
	if strings.EqualFold(kernel, "OpenBSD") && utils.CommandExists("sysctl") && files["openbsd-version"] == "" {
		files["openbsd-version"] = utils.RunCommand("sysctl", "-n", "kern.version")
	}
	if strings.Contains(strings.ToLower(files["proc-version"]), "microsoft") && utils.CommandExists("reg.exe") {
		product := regValue(utils.RunCommand("reg.exe", "query", `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "/v", "ProductName"), "ProductName")
		build := regValue(utils.RunCommand("reg.exe", "query", `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "/v", "CurrentBuildNumber"), "CurrentBuildNumber")
		if product != "" {
			files["windows-product"] = windowsProduct(product, build)
		}
	}
	if text := files["lsb-release"]; text != "" && files["chrome-milestone"] == "" {
		info := parseOSReleaseText(text)
		files["chrome-milestone"] = info["CHROMEOS_RELEASE_CHROME_MILESTONE"]
	}
}

func regValue(out, name string) string {
	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, name) || !strings.Contains(line, "REG_") {
			continue
		}
		fields := strings.Fields(line)
		for i, field := range fields {
			if strings.HasPrefix(field, "REG_") && i+1 < len(fields) {
				return strings.Join(fields[i+1:], " ")
			}
		}
	}
	return ""
}
