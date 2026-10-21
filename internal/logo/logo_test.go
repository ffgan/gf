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
		"HarmonyOS":       "Android",
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
