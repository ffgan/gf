package cli

import "strings"

// ClassifyKernel maps uname -s onto the name the rest of gf switches on.
// BSD kernels keep their own name so FreeBSD, OpenBSD, and NetBSD can use
// their own userland instead of one shared "BSD" bucket.
func ClassifyKernel(kernelName string) (osName string, known bool) {
	raw := strings.TrimSpace(kernelName)
	key := strings.ToLower(raw)
	switch key {
	case "darwin", "xnu":
		return Darwin, true
	case "sunos":
		return SunOS, true
	case "haiku":
		return Haiku, true
	case "minix", "minix3":
		return MINIX, true
	case "aix":
		return AIX, true
	case "irix", "irix64":
		return IRIX, true
	case "freemint", "mint", "tos":
		return FreeMiNT, true
	case "interix":
		return Interix, true
	case "ironclad":
		return Ironclad, true
	case "osf1", "tru64", "digitalunix", "digital unix":
		return digitalUNIX, true
	case "linux":
		return Linux, true
	case "ravynos":
		return RavynOS, true
	case "dragonfly", "dragonflybsd":
		return DragonFly, true
	case "freebsd", "kfreebsd", "gnukfreebsd":
		return FreeBSD, true
	case "openbsd":
		return OpenBSD, true
	case "netbsd":
		return NetBSD, true
	case "bitrig":
		return Bitrig, true
	case "bsd":
		return BSD, true
	}
	// Debian GNU/kFreeBSD reports a GNU-prefixed kernel but is a FreeBSD userland.
	if strings.Contains(key, "kfreebsd") {
		return FreeBSD, true
	}
	if strings.HasPrefix(key, "gnu") {
		return Linux, true
	}
	if strings.HasSuffix(key, "bsd") {
		return raw, true
	}
	upper := strings.ToUpper(raw)
	if strings.HasPrefix(upper, "CYGWIN") ||
		strings.HasPrefix(upper, "MSYS") ||
		strings.HasPrefix(upper, "MINGW") ||
		key == "windows_nt" || key == "windows" {
		return Windows, true
	}
	if canon, ok := otherOS[key]; ok {
		return canon, true
	}
	return raw, false
}

// BSDFamily reports whether osName is a BSD userland gf knows how to query.
func BSDFamily(osName string) bool {
	return isBSDFamily(osName)
}

func isBSDFamily(osName string) bool {
	switch strings.ToLower(osName) {
	case "freebsd", "openbsd", "netbsd", "dragonfly", "bitrig", "bsd",
		"midnightbsd", "hardenedbsd", "ghostbsd", "nomadbsd":
		return true
	default:
		return strings.HasSuffix(strings.ToLower(osName), "bsd")
	}
}

// otherOS marks kernels that have no dedicated probe but should still print
// under their own name instead of an unknown-OS warning.
var otherOS = map[string]string{
	"hp-ux":       "HP-UX",
	"hpux":        "HP-UX",
	"qnx":         "QNX",
	"qnxnto":      "QNX",
	"neutrino":    "QNX",
	"os/2":        "OS/2",
	"os2":         "OS/2",
	"ecomstation": "eComStation",
	"arcaos":      "ArcaOS",
	"reactos":     "ReactOS",
	"morphos":     "MorphOS",
	"aros":        "AROS",
	"amigaos":     "AmigaOS",
	"syllable":    "Syllable",
	"serenity":    "Serenity",
	"serenityos":  "Serenity",
	"redox":       "Redox",
	"fuchsia":     "Fuchsia",
	"zircon":      "Fuchsia",
	"sco_sv":      "SCO_SV",
	"unixware":    "UnixWare",
	"openunix":    "UnixWare",
	"openserver":  "SCO_SV",
	"nonstop":     "NonStop",
	"plan9":       "Plan 9",
	"9front":      "Plan 9",
	"sortix":      "Sortix",
	"helenos":     "HelenOS",
	"managarm":    "managarm",
	"toaru":       "ToaruOS",
	"misaka":      "ToaruOS",
	"genode":      "Genode",
	"sel4":        "seL4",
	"nuttx":       "NuttX",
	"zephyr":      "Zephyr",
	"vxworks":     "VxWorks",
	"integrity":   "INTEGRITY",
	"lynxos":      "LynxOS",
	"pikeos":      "PikeOS",
	"rtems":       "RTEMS",
	"contiki":     "Contiki",
	"riot":        "RIOT",
	"freedos":     "FreeDOS",
	"kolibrios":   "KolibriOS",
	"menuetos":    "MenuetOS",
	"templeos":    "TempleOS",
	"risc os":     "RISC OS",
	"riscos":      "RISC OS",
	"beos":        "BeOS",
	"zeta":        "Zeta",
	"skyos":       "SkyOS",
	"visopsys":    "Visopsys",
	"inferno":     "Inferno",
	"atheos":      "AtheOS",
	"symbian":     "Symbian",
	"palmos":      "PalmOS",
}
