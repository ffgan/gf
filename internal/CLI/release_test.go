package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatDistroOSReleaseWins(t *testing.T) {
	files := map[string]string{
		"os-release": `NAME="Ubuntu"
VERSION="24.04.5 LTS (Noble Numbat)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 24.04.5 LTS"
`,
		"debian_version": "trixie/sid\n",
		"issue":          `Ubuntu 24.04.5 LTS \n \l`,
	}
	display, ascii, _, id, like := FormatDistro(Linux, "off", files)
	if display != "Ubuntu 24.04.5 LTS" || id != "ubuntu" || like != "debian" {
		t.Fatalf("display=%q id=%q like=%q", display, id, like)
	}
	if ascii != "Ubuntu" {
		t.Fatalf("ascii=%q", ascii)
	}
}

func TestFormatDistroShorthand(t *testing.T) {
	files := map[string]string{
		"os-release": "NAME=Alpine\nID=alpine\nVERSION_ID=3.20.3\nPRETTY_NAME=\"Alpine Linux 3.20.3\"\n",
	}
	if got, _, _, _, _ := FormatDistro(Linux, "tiny", files); got != "Alpine" {
		t.Fatalf("tiny: %q", got)
	}
	if got, _, _, _, _ := FormatDistro(Linux, "on", files); got != "Alpine 3.20.3" {
		t.Fatalf("on: %q", got)
	}
}

func TestFormatDistroLegacyFiles(t *testing.T) {
	display, ascii, _, _, _ := FormatDistro(FreeBSD, "off", nil)
	if display != "FreeBSD" || ascii != "FreeBSD" {
		t.Fatalf("empty freebsd: %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{"alpine-release": "3.19.1\n"})
	if display != "Alpine 3.19.1" || ascii != "Alpine" {
		t.Fatalf("alpine: %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{"debian_version": "12.7\n"})
	if display != "Debian 12.7" || ascii != "Debian" {
		t.Fatalf("debian: %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{"arch-release": ""})
	if display != "Arch Linux" || ascii != "Arch" {
		t.Fatalf("arch: %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"artix-release": "",
		"arch-release":  "",
	})
	if display != "Artix Linux" || ascii != "Artix" {
		t.Fatalf("artix: %q %q", display, ascii)
	}
	display, ascii, _, id, _ := FormatDistro(Linux, "off", map[string]string{
		"openwrt_release": "DISTRIB_ID='OpenWrt'\nDISTRIB_RELEASE='23.05.3'\nDISTRIB_DESCRIPTION='OpenWrt 23.05.3'\n",
	})
	if display != "OpenWrt 23.05.3" || ascii != "Wrt" || id != "openwrt" {
		t.Fatalf("openwrt: %q %q %q", display, ascii, id)
	}
	display, ascii, _, id, _ = FormatDistro(Linux, "off", map[string]string{
		"lsb-release": "DISTRIB_ID=ChromeOS\nDISTRIB_RELEASE=123\nCHROMEOS_RELEASE_NAME=Chrome OS\n",
	})
	if !strings.Contains(display, "Chrome OS") || ascii != "Chrom" || id != "chromeos" {
		t.Fatalf("chrome: %q %q %q", display, ascii, id)
	}
	display, ascii, _, id, _ = FormatDistro(Linux, "off", map[string]string{
		"os-release": "NAME=\"Amazon Linux\"\nID=amzn\nVERSION_ID=2023\nPRETTY_NAME=\"Amazon Linux 2023\"\n",
	})
	if display != "Amazon Linux 2023" || ascii != "Amazon" || id != "amzn" {
		t.Fatalf("amazon: %q %q %q", display, ascii, id)
	}
	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{"gobo-version": "017\n"})
	if display != "GoboLinux 017" || ascii != "GoboLinux" {
		t.Fatalf("gobo: %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"redhat-release": "Red Hat Enterprise Linux release 6.10\n",
	})
	if ascii != "RHEL" || !strings.Contains(display, "Red Hat") {
		t.Fatalf("rhel: %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro(illumos, "off", map[string]string{
		"smartos-product": "Name: SmartOS\nImage: 20240101\n",
	})
	if display != "SmartOS" || ascii != "SmartOS" {
		t.Fatalf("smartos: %q %q", display, ascii)
	}
	display, _, _, _, _ = FormatDistro(Linux, "off", map[string]string{"issue": `Welcome \n \l`})
	if display != "Welcome" {
		t.Fatalf("issue: %q", display)
	}
}

func TestClassifyKernelFamilies(t *testing.T) {
	cases := map[string]string{
		"Linux":          Linux,
		"GNU":            Linux,
		"GNU/Linux":      Linux,
		"GNU/kFreeBSD":   FreeBSD,
		"FreeBSD":        FreeBSD,
		"MidnightBSD":    "MidnightBSD",
		"DragonFlyBSD":   DragonFly,
		"Minix":          MINIX,
		"IRIX64":         IRIX,
		"HP-UX":          "HP-UX",
		"SerenityOS":     "Serenity",
		"Fuchsia":        "Fuchsia",
		"CYGWIN_NT-10.0": Windows,
		"SunOS":          SunOS,
		"ravynos":        RavynOS,
	}
	for in, want := range cases {
		got, known := ClassifyKernel(in)
		if got != want || !known {
			t.Errorf("%s: got %q known=%v want %q", in, got, known, want)
		}
	}
	if got, known := ClassifyKernel("NotAnOS"); known || got != "NotAnOS" {
		t.Fatalf("unknown: %q %v", got, known)
	}
	if !isBSDFamily("GhostBSD") || !isBSDFamily("DragonFly") || isBSDFamily(Linux) {
		t.Fatal("bsd family")
	}
}

func TestMacOSOlderCodenames(t *testing.T) {
	if macOSCodename("10.11.6") != "OS X El Capitan" || macOSCodename("10.4.11") != "Mac OS X Tiger" {
		t.Fatal(macOSCodename("10.11.6"), macOSCodename("10.4.11"))
	}
	if macOSCodename("10.0.4") != "Mac OS X Cheetah" || macOSCodename("11.7") != "macOS Big Sur" {
		t.Fatal(macOSCodename("10.0.4"), macOSCodename("11.7"))
	}
	if macOSCodename("16.0") != "macOS" || macOSCodename("15.2") != "macOS Sequoia" {
		t.Fatal(macOSCodename("16.0"), macOSCodename("15.2"))
	}
}

func TestCountPkgDBDirs(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"curl", "pkg"} {
		if err := os.Mkdir(filepath.Join(root, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "local.sqlite"), []byte("not a database"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := countPkgDB(root); got != 2 {
		t.Fatalf("got %d want 2", got)
	}
}
