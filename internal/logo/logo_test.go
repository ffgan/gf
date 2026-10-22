package logo

import "testing"

func TestResolveUbuntu(t *testing.T) {
	got := Resolve("Ubuntu 24.04.5 LTS x86_64")
	if got.FileName != "Ubuntu" {
		t.Fatalf("filename = %q", got.FileName)
	}
	id := Resolve("ubuntu")
	if id.FileName != "Ubuntu" {
		t.Fatalf("id filename = %q", id.FileName)
	}
}

func TestResolveMoreDistros(t *testing.T) {
	cases := map[string]string{
		"OpenWrt":         "Wrt",
		"uos":             "Uos",
		"Kylin":           "Ubuntu-Kylin",
		"openKylin":       "openKylin",
		"HarmonyOS":       "HarmonyOS",
		"OpenHarmony":     "HarmonyOS",
		"HyperOS":         "Android",
		"Raspberry Pi OS": "Raspbian",
		"Loongnix":        "Debian",
		"Anolis OS":       "CentOS",
		"Rocky Linux":     "rocky",
		"archlinux":       "Arch",
		"opensuse-leap":   "openSUSE_Leap",
		"ghostbsd":        "GhostBSD",
		"omnios":          "OmniOS",
		"azurelinux":      "CBL-Mariner",
	}
	for name, file := range cases {
		if got := Resolve(name).FileName; got != file {
			t.Errorf("%s: got %q want %q", name, got, file)
		}
	}
}

func TestResolveNeofetchHistorical(t *testing.T) {
	cases := map[string]string{
		"Windows 11 Pro":       "windows11",
		"windows11":            "windows11",
		"Windows 10 Home":      "windows10",
		"Windows 8.1 Pro":      "windows10",
		"Windows7":             "Windows",
		"Windows 95":           "Windows95",
		"Windows XP":           "Windows",
		"mac_small":            "mac_small",
		"Mac OS X Tiger":       "mac",
		"AOSC OS/Retro":        "AOSC_Retro",
		"AlterLinux":           "Alter",
		"MX_Linux":             "MX",
		"iOS 17.4":             "mac",
		"Chrome OS":            "Chrom",
		"TUXEDO OS":            "Tuxedo",
		"OpenSolaris":          "OpenSolaris",
		"FiwixOS":              "Fiwix",
		"Blackberry 10":        "QNX",
		"whiite-linux":         "whiite-linux",
		"wii-linux-ngx":        "whiite-linux",
		"NixOwOS":              "NixOwOS",
		"Gloire":               "Gloire",
		"void_old":             "void_old",
		"Ubuntu on Windows 11": "Ubuntu",
		"SomeOS on Windows 10": "windows10",
	}
	for name, file := range cases {
		if got := Resolve(name).FileName; got != file {
			t.Errorf("%s: got %q want %q", name, got, file)
		}
	}
	if got := Resolve("Astra Linux 2.12.0 ALCE").Colors; len(got) != 2 || got[0] != 160 || got[1] != 231 {
		t.Fatalf("ALCE colors %v", got)
	}
	if got := Resolve("Astra Linux 1.7 ALSE").Colors; len(got) != 2 || got[0] != 32 || got[1] != 231 {
		t.Fatalf("ALSE colors %v", got)
	}
	if got := Resolve("NixOwOS").RawColors; len(got) != 5 || got[0] != "#5277c3" {
		t.Fatalf("nixowos colors %v", got)
	}
}

func TestResolveFallback(t *testing.T) {
	got := Resolve("not-a-real-distro")
	if got.FileName != "Linux" {
		t.Fatalf("filename = %q", got.FileName)
	}
}

func TestColorBasic(t *testing.T) {
	got := color("1")
	if got != "\033[0m\033[31m" {
		t.Fatalf("color = %q", got)
	}
	if color("") != Reset {
		t.Fatal("empty color should reset")
	}
}

func TestPropagateColor(t *testing.T) {
	lines := propagateColor([]string{
		"\033[31m\033[1mfirst",
		"second",
	})
	if lines[1] != "\033[31m\033[1msecond" {
		t.Fatalf("propagated = %q", lines[1])
	}
}

func TestSetColorsUsesDecimal(t *testing.T) {
	asciiColorOverride = nil
	asciiBold = ""
	setColors(4)
	if c1 != "\033[0m\033[34m" {
		t.Fatalf("c1 = %q", c1)
	}
}
