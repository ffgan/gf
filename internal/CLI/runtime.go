package cli

import (
	"os"
	"strings"

	"github.com/ffgan/gf/internal/hw"
	"github.com/ffgan/gf/internal/utils"
)

// HostInput is everything PickHost needs. Callers read the files; the
// decision itself stays free of the filesystem so empty DMI cannot hide a
// later source.
type HostInput struct {
	AndroidBrand   string
	AndroidModel   string
	DeviceTree     string
	DMIVendor      string
	DMIProduct     string
	DMIVersion     string
	DMIBoardVendor string
	DMIBoardName   string
	SysinfoModel   string
	Kernel         string
	Hypervisor     string
	WSL            bool
}

// PickHost chooses the machine name. The first source that still has text
// after junk strings are removed wins.
func PickHost(in HostInput) string {
	if brand, model := meaningful(in.AndroidBrand), meaningful(in.AndroidModel); brand != "" || model != "" {
		return strings.TrimSpace(brand + " " + model)
	}
	if dt := meaningful(stripCtl(in.DeviceTree)); dt != "" {
		return dt
	}
	vendor := meaningful(in.DMIVendor)
	product := meaningful(in.DMIProduct)
	version := meaningful(in.DMIVersion)
	if vendor != "" || product != "" {
		model := strings.TrimSpace(vendor + " " + product)
		if version != "" && !strings.Contains(model, version) {
			model = strings.TrimSpace(model + " " + version)
		}
		return finishHost(model)
	}
	if bv, bn := meaningful(in.DMIBoardVendor), meaningful(in.DMIBoardName); bv != "" || bn != "" {
		return finishHost(strings.TrimSpace(bv + " " + bn))
	}
	if s := meaningful(in.SysinfoModel); s != "" {
		return finishHost(s)
	}
	if in.WSL || wslFromKernel(in.Kernel) {
		return "Windows Subsystem for Linux"
	}
	if h := strings.TrimSpace(in.Hypervisor); h != "" {
		return h
	}
	return ""
}

func finishHost(model string) string {
	model = cleanModelName(model)
	switch {
	case strings.HasPrefix(model, "Standard PC"):
		return "KVM/QEMU (" + model + ")"
	case strings.HasPrefix(model, "OpenBSD"):
		return "vmm (" + model + ")"
	default:
		return strings.TrimSpace(model)
	}
}

func meaningful(s string) string {
	s = stripCtl(s)
	s = cleanModelName(s)
	if s == "" {
		return ""
	}
	switch strings.ToLower(s) {
	case "none", "n/a", "na", "unknown", "not available", "not specified",
		"to be filled by o.e.m.", "default string", "system product name",
		"system manufacturer", "chassis", "empty", "invalid", "oem":
		return ""
	}
	if strings.Trim(s, "0 ") == "" {
		return ""
	}
	return s
}

func wslFromKernel(kernel string) bool {
	return strings.Contains(strings.ToLower(kernel), "microsoft")
}

// RuntimeSig describes a Linux userspace. Exists and Read may be nil.
type RuntimeSig struct {
	ID     string
	IDLike string
	Name   string
	Kernel string
	Env    map[string]string
	Exists func(string) bool
	Read   func(string) string
}

// RuntimeTags reports layered environments that the distro name itself does
// not already mention. The order is stable.
func RuntimeTags(s RuntimeSig) []string {
	exists := s.Exists
	if exists == nil {
		exists = func(string) bool { return false }
	}
	read := s.Read
	if read == nil {
		read = func(string) string { return "" }
	}
	id := strings.ToLower(s.ID)
	name := strings.ToLower(s.Name + " " + s.ID + " " + s.IDLike)
	var tags []string
	add := func(tag string, on bool) {
		if on {
			tags = append(tags, tag)
		}
	}

	termux := s.Env["TERMUX_VERSION"] != "" || strings.Contains(s.Env["PREFIX"], "com.termux") || exists("/data/data/com.termux/files/usr")
	android := s.Env["ANDROID_ROOT"] != "" || s.Env["ANDROID_DATA"] != "" || exists("/system/build.prop")
	add("Termux", termux)
	add("Android", android && !termux)

	kernel := strings.ToLower(s.Kernel)
	wsl2 := strings.Contains(kernel, "wsl2") || strings.Contains(kernel, "microsoft-standard")
	wsl := wsl2 || strings.Contains(kernel, "microsoft")
	add("WSL2", wsl2)
	add("WSL1", wsl && !wsl2)

	add("Docker", exists("/.dockerenv"))
	add("Podman", exists("/run/.containerenv"))
	add("Kubernetes", exists("/var/run/secrets/kubernetes.io") || exists("/var/run/secrets/kubernetes.io/serviceaccount"))

	lsb := strings.ToLower(read("/etc/lsb-release"))
	chrome := strings.Contains(lsb, "chromeos")
	add("Crostini", exists("/dev/.cros_milestone") && !chrome)
	add("ChromeOS", chrome && !strings.Contains(name, "chrome"))

	add("OpenWrt", (id == "openwrt" || exists("/etc/openwrt_release")) && !strings.Contains(name, "openwrt") && !strings.Contains(name, "wrt"))
	debianish := id == "debian" || strings.Contains(s.IDLike, "debian") || strings.Contains(name, "debian")
	add("Proxmox", exists("/etc/pve") && debianish && !strings.Contains(name, "proxmox"))

	dt := strings.ToLower(read("/sys/firmware/devicetree/base/model"))
	rpi := exists("/etc/rpi-issue") || strings.Contains(dt, "raspberry")
	add("Raspberry Pi OS", rpi && debianish && !strings.Contains(name, "raspberry") && id != "raspbian" && id != "raspios")

	add("Bedrock", exists("/bedrock/etc/bedrock-release") && !strings.Contains(name, "bedrock"))
	return tags
}

// LogoKey returns a radix key whose art we already ship when the distro's
// own NAME would not prefix-match that art.
func LogoKey(name, id string) string {
	key := strings.ToLower(strings.TrimSpace(id))
	lowName := strings.ToLower(strings.TrimSpace(name))
	if art, ok := distroLogoIDs[key]; ok {
		return art
	}
	switch key {
	case "uos", "uniontech":
		return "Uos"
	case "openkylin":
		return "openKylin"
	case "kylin":
		return "Ubuntu Kylin"
	case "openeuler", "euleros":
		return "openEuler"
	case "harmonyos", "openharmony", "hyperos":
		return "Android"
	case "anolis", "opencloudos":
		return "CentOS"
	case "loongnix":
		return "Debian"
	case "openwrt":
		return "Wrt"
	case "raspbian", "raspios", "raspberrypi":
		return "Raspbian"
	case "rocky":
		return "rocky"
	}
	switch {
	case strings.Contains(lowName, "raspberry pi os"):
		return "Raspbian"
	case strings.Contains(lowName, "rocky linux"):
		return "rocky"
	case strings.Contains(lowName, "openharmony"), strings.Contains(lowName, "harmonyos"), strings.Contains(lowName, "hyperos"):
		return "Android"
	case strings.Contains(lowName, "openkylin"):
		return "openKylin"
	case strings.Contains(lowName, "kylin"):
		return "Ubuntu Kylin"
	case strings.Contains(lowName, "openeuler"):
		return "openEuler"
	case strings.Contains(lowName, "openwrt"):
		return "Wrt"
	case lowName == "uos" || strings.HasPrefix(lowName, "uos "):
		return "Uos"
	case strings.Contains(lowName, "loongnix"):
		return "Debian"
	case strings.Contains(lowName, "anolis"), strings.Contains(lowName, "opencloudos"):
		return "CentOS"
	}
	return ""
}

// distroLogoIDs maps os-release IDs whose spelling is not the logo key.
var distroLogoIDs = map[string]string{
	"amzn":                "Amazon",
	"amazonlinux":         "Amazon",
	"ol":                  "Oracle",
	"oraclelinux":         "Oracle",
	"ubuntu-mate":         "Ubuntu MATE",
	"ubuntumate":          "Ubuntu MATE",
	"ubuntu-budgie":       "Ubuntu Budgie",
	"ubuntubudgie":        "Ubuntu Budgie",
	"ubuntu-studio":       "Ubuntu Studio",
	"ubuntustudio":        "Ubuntu Studio",
	"ubuntu-cinnamon":     "Ubuntu Cinnamon",
	"ubuntucinnamon":      "Ubuntu Cinnamon",
	"ubuntu-kylin":        "Ubuntu Kylin",
	"ubuntukylin":         "Ubuntu Kylin",
	"ubuntu-touch":        "Ubuntu Touch",
	"ubuntu-sway":         "Ubuntu Sway",
	"ubuntusway":          "Ubuntu Sway",
	"linuxlite":           "Linux Lite",
	"linux-lite":          "Linux Lite",
	"fedora-silverblue":   "Fedora Silverblue",
	"silverblue":          "Fedora Silverblue",
	"fedora-kinoite":      "Fedora Kinoite",
	"kinoite":             "Fedora Kinoite",
	"fedora-sericea":      "Fedora Sericea",
	"sericea":             "Fedora Sericea",
	"fedora-coreos":       "Fedora CoreOS",
	"clear-linux-os":      "Clear Linux OS",
	"clearlinux":          "Clear Linux OS",
	"cbl-mariner":         "CBL-Mariner",
	"mariner":             "CBL-Mariner",
	"opensuse-slowroll":   "openSUSE Tumbleweed-Slowroll",
	"slowroll":            "openSUSE Tumbleweed-Slowroll",
	"sles":                "SUSE",
	"sled":                "SUSE",
	"suse":                "SUSE",
	"aosc":                "AOSC OS",
	"aosc-os":             "AOSC OS",
	"archlabs":            "ARCHlabs",
	"crystal-linux":       "Crystal Linux",
	"rebornos":            "Reborn OS",
	"ultramarine":         "Ultramarine Linux",
	"sourcemage":          "Source Mage",
	"rhinolinux":          "Rhino Linux",
	"rhino":               "Rhino Linux",
	"vanillaos":           "Vanilla",
	"parrotos":            "Parrot",
	"sailfishos":          "Sailfish",
	"chromeos":            "Chrom",
	"chromiumos":          "Chrom",
	"pmos":                "PostMarketOS",
	"mx-linux":            "MX",
	"midnightbsd":         "FreeBSD",
	"pfsense":             "FreeBSD",
	"astra":               "Astra Linux",
	"astralinux":          "Astra Linux",
	"miraclelinux":        "MIRACLE LINUX",
	"calculate":           "Calculate",
	"calculate-linux":     "Calculate",
	"lfs":                 "LinuxFromScratch",
	"linuxfromscratch":    "LinuxFromScratch",
	"instantos":           "instantOS",
	"kde-neon":            "KDE",
	"linuxmint":           "Linux Mint",
	"pop":                 "Pop!_OS",
	"neon":                "KDE",
	"archlinux":           "Arch",
	"archlinuxarm":        "Arch",
	"archarm":             "Arch",
	"alarm":               "Arch",
	"manjaro-arm":         "Manjaro",
	"sles_sap":            "SUSE",
	"suse_sap":            "SUSE",
	"opensuse-leap":       "openSUSE Leap",
	"opensuse-tumbleweed": "openSUSE Tumbleweed",
	"opensuse-microos":    "openSUSE",
	"suse-microos":        "openSUSE",
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

func AnnotateDistro(distro string, tags []string) string {
	if len(tags) == 0 {
		return distro
	}
	return distro + " (" + strings.Join(tags, ", ") + ")"
}

func readSmall(path string) string {
	b, err := os.ReadFile(path)
	if err != nil || len(b) > 256*1024 {
		return ""
	}
	return string(b)
}

func linuxHostInput() HostInput {
	in := HostInput{
		DeviceTree:     readSmall("/sys/firmware/devicetree/base/model"),
		DMIVendor:      readSmall("/sys/devices/virtual/dmi/id/sys_vendor"),
		DMIProduct:     readSmall("/sys/devices/virtual/dmi/id/product_name"),
		DMIVersion:     readSmall("/sys/devices/virtual/dmi/id/product_version"),
		DMIBoardVendor: readSmall("/sys/devices/virtual/dmi/id/board_vendor"),
		DMIBoardName:   readSmall("/sys/devices/virtual/dmi/id/board_name"),
		SysinfoModel:   readSmall("/tmp/sysinfo/model"),
		Kernel:         readSmall("/proc/version"),
		Hypervisor:     hw.Hypervisor(),
	}
	if utils.PathExists("/system/build.prop") || os.Getenv("ANDROID_ROOT") != "" {
		in.AndroidBrand, in.AndroidModel = androidProduct()
	}
	return in
}

func androidProduct() (string, string) {
	props := map[string]string{}
	for _, path := range []string{"/system/build.prop", "/vendor/build.prop"} {
		for _, line := range strings.Split(readSmall(path), "\n") {
			k, v, ok := strings.Cut(line, "=")
			if ok && props[k] == "" {
				props[k] = strings.TrimSpace(v)
			}
		}
	}
	brand := props["ro.product.brand"]
	model := props["ro.product.model"]
	if brand == "" && model == "" && utils.CommandExists("getprop") {
		brand = utils.RunCommand("getprop", "ro.product.brand")
		model = utils.RunCommand("getprop", "ro.product.model")
	}
	return brand, model
}

func iphoneName(machine string) string {
	if name, ok := iphoneModels[machine]; ok {
		return name
	}
	return machine
}

var iphoneModels = map[string]string{
	"iPhone14,2": "iPhone 13 Pro",
	"iPhone14,3": "iPhone 13 Pro Max",
	"iPhone14,4": "iPhone 13 mini",
	"iPhone14,5": "iPhone 13",
	"iPhone14,6": "iPhone SE (3rd generation)",
	"iPhone14,7": "iPhone 14",
	"iPhone14,8": "iPhone 14 Plus",
	"iPhone15,2": "iPhone 14 Pro",
	"iPhone15,3": "iPhone 14 Pro Max",
	"iPhone15,4": "iPhone 15",
	"iPhone15,5": "iPhone 15 Plus",
	"iPhone16,1": "iPhone 15 Pro",
	"iPhone16,2": "iPhone 15 Pro Max",
	"iPhone17,1": "iPhone 16 Pro",
	"iPhone17,2": "iPhone 16 Pro Max",
	"iPhone17,3": "iPhone 16",
	"iPhone17,4": "iPhone 16 Plus",
}
