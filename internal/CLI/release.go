package cli

import "strings"

// FormatDistro names a distribution from the release files a system actually
// has. os-release wins. Older and embedded systems fall through to the single
// file they do ship, and the result is the OS name when nothing else exists.
func FormatDistro(osName, shorthand string, files map[string]string) (display, ascii, name, id, idLike string) {
	if text := files["os-release"]; text != "" {
		info := parseOSReleaseText(text)
		name = info["NAME"]
		if name == "" {
			name = info["ID"]
		}
		id = info["ID"]
		idLike = info["ID_LIKE"]
		version := info["VERSION"]
		if version == "" {
			version = info["VERSION_ID"]
		}
		pretty := info["PRETTY_NAME"]
		switch shorthand {
		case "tiny":
			display = name
		case "on":
			display = strings.TrimSpace(name + " " + version)
		default:
			if pretty != "" {
				display = pretty
			} else {
				display = strings.TrimSpace(name + " " + version)
			}
		}
		ascii = distroArtName(name, id, idLike, pretty)
		return display, ascii, name, id, idLike
	}

	if text := files["lsb-release"]; text != "" {
		info := parseOSReleaseText(text)
		id = firstNonEmpty(info["DISTRIB_ID"], info["CHROMEOS_RELEASE_NAME"])
		version := firstNonEmpty(info["DISTRIB_RELEASE"], info["CHROMEOS_RELEASE_VERSION"])
		pretty := info["DISTRIB_DESCRIPTION"]
		name = id
		if strings.Contains(strings.ToLower(text), "chromeos") || info["CHROMEOS_RELEASE_NAME"] != "" {
			name = "Chrome OS"
			id = "chromeos"
		}
		switch shorthand {
		case "tiny":
			display = name
		case "on":
			display = strings.TrimSpace(name + " " + version)
		default:
			display = firstNonEmpty(pretty, strings.TrimSpace(name+" "+version))
		}
		ascii = distroArtName(name, id, "", pretty)
		return display, ascii, name, id, ""
	}

	if text := files["openwrt_release"]; text != "" {
		info := parseOSReleaseText(text)
		name = info["DISTRIB_ID"]
		if name == "" {
			name = "OpenWrt"
		}
		id = "openwrt"
		version := info["DISTRIB_RELEASE"]
		pretty := info["DISTRIB_DESCRIPTION"]
		display = releaseDisplay(shorthand, name, version, pretty)
		ascii = distroArtName(name, id, "", pretty)
		return display, ascii, name, id, ""
	}

	if text, ok := files["alpine-release"]; ok && strings.TrimSpace(text) != "" {
		name, id = "Alpine", "alpine"
		display = releaseDisplay(shorthand, name, strings.TrimSpace(text), "")
		return display, "Alpine", name, id, ""
	}
	if _, ok := files["artix-release"]; ok {
		name, id = "Artix Linux", "artix"
		return name, "Artix", name, id, ""
	}
	if _, ok := files["arch-release"]; ok {
		name, id = "Arch Linux", "arch"
		display = name
		return display, "Arch", name, id, ""
	}
	if text := strings.TrimSpace(files["gentoo-release"]); text != "" {
		name, id = "Gentoo", "gentoo"
		display = text
		return display, "Gentoo", name, id, ""
	}
	if text := strings.TrimSpace(files["slackware-version"]); text != "" {
		name, id = "Slackware", "slackware"
		if strings.Contains(strings.ToLower(text), "slackware") {
			display = text
		} else {
			display = "Slackware " + text
		}
		return display, "Slackware", name, id, ""
	}
	if text := strings.TrimSpace(files["suse-release"]); text != "" {
		name, id = "openSUSE", "opensuse"
		art := "openSUSE"
		low := strings.ToLower(text)
		if strings.Contains(low, "enterprise") || strings.Contains(low, "sles") || strings.Contains(low, "sled") {
			name, id, art = "SUSE", "sles", "SUSE"
		}
		return text, art, name, id, ""
	}
	for _, key := range []string{"fedora-release", "centos-release", "redhat-release", "system-release"} {
		if text := strings.TrimSpace(files[key]); text != "" {
			name = text
			id = strings.TrimSuffix(key, "-release")
			return text, redhatArt(text), name, id, ""
		}
	}
	if text := strings.TrimSpace(files["debian_version"]); text != "" {
		name, id = "Debian", "debian"
		display = "Debian " + text
		return display, "Debian", name, id, ""
	}
	if text := strings.TrimSpace(files["bedrock"]); text != "" {
		name, id = "Bedrock Linux", "bedrock"
		display = firstNonEmpty(text, name)
		return display, "Bedrock", name, id, ""
	}
	if text := files["build.prop"]; text != "" {
		props := parseOSReleaseText(text)
		ver := props["ro.build.version.release"]
		name, id = "Android", "android"
		display = strings.TrimSpace(name + " " + ver)
		return display, "Android", name, id, ""
	}
	if text := strings.TrimSpace(files["haiku-version"]); text != "" {
		name, id = "Haiku", "haiku"
		line := strings.TrimSpace(strings.Split(text, "\n")[0])
		if strings.HasPrefix(strings.ToLower(line), "haiku") {
			display = line
		} else {
			display = "Haiku " + line
		}
		return display, "Haiku", name, id, ""
	}
	if text := strings.TrimSpace(files["solaris-release"]); text != "" {
		name = osName
		if name == "" {
			name = "Solaris"
		}
		line := strings.TrimSpace(strings.Split(text, "\n")[0])
		art := solarisArt(line, name)
		return line, art, name, strings.ToLower(name), ""
	}
	if text := strings.TrimSpace(files["smartos-product"]); text != "" {
		name = smartOSName(text)
		return name, "SmartOS", name, "smartos", ""
	}
	if text := strings.TrimSpace(files["freebsd-version"]); text != "" {
		name, id = "FreeBSD", "freebsd"
		display = "FreeBSD " + text
		return display, "FreeBSD", name, id, ""
	}
	if text := files["armbian-release"]; text != "" {
		info := parseOSReleaseText(text)
		name = firstNonEmpty(info["PRETTY_NAME"], info["BOARD_NAME"], "Armbian")
		id = "armbian"
		version := info["VERSION"]
		display = releaseDisplay(shorthand, "Armbian", version, info["PRETTY_NAME"])
		if shorthand == "" || (shorthand != "tiny" && shorthand != "on" && info["PRETTY_NAME"] != "") {
			display = firstNonEmpty(info["PRETTY_NAME"], display)
		}
		return display, "Armbian", name, id, ""
	}
	if display, ascii, name, id, ok := legacyRelease(shorthand, files); ok {
		return display, ascii, name, id, ""
	}
	if text := cleanIssue(files["issue"]); text != "" {
		return text, text, text, "", ""
	}

	if isBSDFamily(osName) || osName == Haiku || osName == Solaris || osName == illumos || osName == AIX || osName == MINIX {
		return osName, osName, osName, strings.ToLower(osName), ""
	}
	return "", "", "", "", ""
}

func legacyRelease(shorthand string, files map[string]string) (display, ascii, name, id string, ok bool) {
	type spec struct {
		key, name, art string
	}
	for _, item := range []spec{
		{"redstar-release", "Red Star OS", "Red Star OS"},
		{"siduction-version", "Siduction", "Siduction"},
		{"gobo-version", "GoboLinux", "GoboLinux"},
		{"sde-version", "SDE", "SDE"},
		{"mcst-version", "OS Elbrus", "OS Elbrus"},
		{"angstrom-version", "Ångström", "Ångström"},
		{"frugalware-release", "Frugalware", "Frugalware"},
		{"pardus-release", "Pardus", "Pardus"},
		{"mageia-release", "Mageia", "Mageia"},
		{"mandriva-release", "Mandriva", "Mandriva"},
		{"meego-release", "MeeGo", "MeeGo"},
		{"sabayon-release", "Sabayon", "Sabayon"},
		{"bluewhite64-version", "Bluewhite64", "Bluewhite64"},
		{"crux-release", "CRUX", "CRUX"},
		{"tizen-release", "Tizen", "Tizen"},
	} {
		text := strings.TrimSpace(files[item.key])
		if text == "" {
			continue
		}
		ver := strings.TrimSpace(strings.Split(text, "\n")[0])
		display = releaseDisplay(shorthand, item.name, ver, "")
		if shorthand != "tiny" && shorthand != "on" && strings.Contains(strings.ToLower(ver), strings.ToLower(item.name)) {
			display = ver
		}
		id = strings.ToLower(strings.ReplaceAll(item.name, " ", ""))
		return display, item.art, item.name, id, true
	}
	return "", "", "", "", false
}

func redhatArt(text string) string {
	low := strings.ToLower(text)
	switch {
	case strings.Contains(low, "centos"):
		return "CentOS"
	case strings.Contains(low, "fedora"):
		return "Fedora"
	case strings.Contains(low, "alma"):
		return "AlmaLinux"
	case strings.Contains(low, "rocky"):
		return "rocky"
	case strings.Contains(low, "amazon"):
		return "Amazon"
	case strings.Contains(low, "oracle"):
		return "Oracle"
	case strings.Contains(low, "red hat"), strings.Contains(low, "redhat"):
		return "RHEL"
	default:
		return text
	}
}

func solarisArt(line, osName string) string {
	low := strings.ToLower(line)
	switch {
	case strings.Contains(low, "omnios"):
		return "OmniOS"
	case strings.Contains(low, "openindiana"):
		return "OpenIndiana"
	case strings.Contains(low, "smartos"):
		return "SmartOS"
	case osName == illumos:
		return osName
	default:
		return Solaris
	}
}

func smartOSName(text string) string {
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "name:") {
			if name := strings.TrimSpace(strings.SplitN(line, ":", 2)[1]); name != "" {
				return name
			}
		}
	}
	return "SmartOS"
}

func cleanIssue(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	line := strings.TrimSpace(strings.Split(text, "\n")[0])
	line = strings.NewReplacer(
		`\n`, "", `\l`, "", `\r`, "", `\m`, "", `\s`, "", `\v`, "", `\S`, "", `\t`, " ", `\\`, `\`,
	).Replace(line)
	return strings.Join(strings.Fields(line), " ")
}

func releaseDisplay(shorthand, name, version, pretty string) string {
	switch shorthand {
	case "tiny":
		return name
	case "on":
		return strings.TrimSpace(name + " " + version)
	default:
		if pretty != "" {
			return pretty
		}
		return strings.TrimSpace(name + " " + version)
	}
}

func distroArtName(name, id, idLike, pretty string) string {
	if key := LogoKey(name, id); key != "" {
		return key
	}
	if key := LogoKey("", firstToken(idLike)); key != "" {
		return key
	}
	return firstNonEmpty(name, id, pretty)
}

func firstToken(s string) string {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func parseOSReleaseText(content string) map[string]string {
	info := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		info[parts[0]] = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	}
	return info
}
