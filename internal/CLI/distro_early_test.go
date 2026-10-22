package cli

import (
	"strings"
	"testing"
)

func TestMacOSHistoricalNames(t *testing.T) {
	if macOSCodename("10.16") != "macOS Big Sur" || macOSCodename("10.16.1") != "macOS Big Sur" {
		t.Fatal(macOSCodename("10.16"), macOSCodename("10.16.1"))
	}
	if got := formatMacOS("10.16", "20A5343j", "off"); got != "macOS Big Sur 10.16 (20A5343j)" {
		t.Fatal(got)
	}
	if got := formatMacOS("10.4.11", "8S2167", "tiny"); got != "Mac OS X 10.4.11" {
		t.Fatal(got)
	}
	if got := formatMacOS("10.11.6", "", "tiny"); got != "OS X 10.11.6" {
		t.Fatal(got)
	}
	if got := formatMacOS("12.6.1", "21G217", "on"); got != "macOS Monterey 12.6.1" {
		t.Fatal(got)
	}
	if got := formatMacOS("12.6.1", "21G217", "tiny"); got != "macOS Monterey 12.6.1" {
		t.Fatal(got)
	}
	if got := formatMacOS("10.0.4", "", "off"); got != "Mac OS X Cheetah 10.0.4" {
		t.Fatal(got)
	}
	display, ascii, hide := formatApple("iPhone OS", "17.4", "21E219", "off")
	if display != "iOS 17.4" || ascii != "iOS" || !hide {
		t.Fatalf("ios %q %q %v", display, ascii, hide)
	}
}

func TestWindowsHistoricalNames(t *testing.T) {
	if got := windowsProduct("Microsoft Windows 10 Pro", "22621"); got != "Windows 11 Pro" {
		t.Fatal(got)
	}
	if got := windowsProduct("Windows 10 Home", "19045"); got != "Windows 10 Home" {
		t.Fatal(got)
	}
	if got := windowsProduct("Windows 10 Pro", "9999"); got != "Windows 10 Pro" {
		t.Fatal(got)
	}
	cases := map[string]string{
		"Windows 11 Pro": "Windows 11",
		"Windows 10 Pro": "Windows 10",
		"Windows 8.1":    "Windows 10",
		"Windows 95":     "Windows95",
		"Windows XP":     "Windows",
		"Windows 7":      "Windows",
		"Windows 2000":   "Windows",
	}
	for in, want := range cases {
		if got := windowsArt(in); got != want {
			t.Errorf("%s: %s", in, got)
		}
	}
}

func TestEarlyDistrosBeatOSRelease(t *testing.T) {
	display, ascii, _, id, _ := FormatDistro(Linux, "off", map[string]string{
		"os-release":      "NAME=Debian\nID=debian\nPRETTY_NAME=\"Debian GNU/Linux\"\n",
		"redstar-release": "Red Star OS release 3.0\n",
	})
	if display != "Red Star OS 3.0" || ascii != "Red Star" || id != "redstar" {
		t.Fatalf("redstar %q %q %q", display, ascii, id)
	}
	display, _, _, _, _ = FormatDistro(Linux, "tiny", map[string]string{
		"redstar-release": "Red Star OS release 3.0\n",
	})
	if display != "Red Star OS" {
		t.Fatal(display)
	}

	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"os-release":       "PRETTY_NAME=\"Ubuntu 22.04.2 LTS\"\nNAME=Ubuntu\n",
		"bedrock":          "Bedrock Linux 0.7.20 Poki\n",
		"bedrock-restrict": "",
	})
	if !strings.Contains(display, "Bedrock Linux 0.7.20") || ascii != "Bedrock" {
		t.Fatalf("bedrock %q %q", display, ascii)
	}
	display, _, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"os-release":       "PRETTY_NAME=\"Ubuntu 22.04\"\nNAME=Ubuntu\nID=ubuntu\n",
		"bedrock":          "Bedrock Linux\n",
		"bedrock-restrict": "1",
	})
	if display != "Ubuntu 22.04" {
		t.Fatal(display)
	}

	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"os-release":      "PRETTY_NAME=\"Debian GNU/Linux 12\"\nNAME=Debian\n",
		"armbian-release": "VERSION=23.02.2\nDISTRIBUTION_CODENAME=jammy\n",
	})
	if display != "Armbian jammy (23.02.2)" || ascii != "Armbian" {
		t.Fatalf("armbian %q %q", display, ascii)
	}

	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"astra-version": "2.12.45\n",
	})
	if display != "Astra Linux 2.12.45 ALCE" || !strings.Contains(ascii, "ALCE") {
		t.Fatalf("astra %q %q", display, ascii)
	}
	display, _, _, _, _ = FormatDistro(Linux, "on", map[string]string{"astra-version": "1.7.4\n"})
	if display != "Astra Linux" {
		t.Fatal(display)
	}

	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"os-release":    "NAME=NixOS\nID=nixos\nVERSION=23.11\n",
		"nixos-marker":  "",
		"nixos-version": "23.11.1234.abcdef",
	})
	if display != "NixOS 23.11.1234.abcdef" || ascii != "NixOS" {
		t.Fatalf("nixos %q %q", display, ascii)
	}

	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"pveversion": "pve-manager/8.1.4/ec5affc9e41f1d1d\n",
	})
	if display != "Proxmox VE 8.1.4" || ascii != "Proxmox" {
		t.Fatalf("pve %q %q", display, ascii)
	}
	display, _, _, _, _ = FormatDistro(Linux, "tiny", map[string]string{"pveversion": "pve-manager/8.1.4/x"})
	if display != "Proxmox VE" {
		t.Fatal(display)
	}
}

func TestLSBAndLateFiles(t *testing.T) {
	display, ascii, _, id, _ := FormatDistro(Linux, "off", map[string]string{
		"os-release":     "NAME=Debian\nID=debian\nVERSION_CODENAME=bookworm\nPRETTY_NAME=\"Debian GNU/Linux 12\"\n",
		"debian_version": "12.7\n",
		"lsb-si":         "Debian",
		"lsb-sd":         "Debian GNU/Linux 12 (bookworm)",
	})
	if display != "Debian 12.7 (bookworm)" || ascii != "Debian" || id != "debian" {
		t.Fatalf("debian lsb %q %q %q", display, ascii, id)
	}
	display, _, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"os-release": "NAME=Ubuntu\nID=ubuntu\nPRETTY_NAME=\"Ubuntu Custom\"\n",
		"lsb-si":     "Ubuntu",
		"lsb-sd":     "Ubuntu 22.04.3 LTS",
	})
	if display != "Ubuntu 22.04.3 LTS" {
		t.Fatal(display)
	}

	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"os-release":      "NAME=Ubuntu\nID=ubuntu\nPRETTY_NAME=\"Ubuntu 24.04 LTS\"\n",
		"xdg-config-dirs": "/etc/xdg/xdg-plasma:/etc/xdg",
	})
	if display != "Kubuntu 24.04 LTS" || ascii != "Kubuntu" {
		t.Fatalf("flavor %q %q", display, ascii)
	}

	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"os-release": "NAME=Fedora Linux\nID=fedora\nVERSION_ID=40\nVARIANT=Silverblue\nPRETTY_NAME=\"Fedora Linux 40 (Silverblue)\"\n",
	})
	if display != "Fedora Linux 40 (Silverblue)" || ascii != "Fedora Silverblue" {
		t.Fatalf("silverblue %q %q", display, ascii)
	}

	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"os-release":      "NAME=Ubuntu\nID=ubuntu\nPRETTY_NAME=\"Ubuntu 22.04 LTS\"\n",
		"proc-version":    "Linux version 5.15.90.1-microsoft-standard",
		"windows-product": "Windows 11 Pro",
	})
	if display != "Ubuntu 22.04 LTS on Windows 11 Pro" || ascii != "Ubuntu" {
		t.Fatalf("wsl %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro(Linux, "tiny", map[string]string{
		"os-release":      "NAME=Ubuntu\nID=ubuntu\nPRETTY_NAME=\"Ubuntu 22.04 LTS\"\n",
		"proc-version":    "Linux version 5.15-Microsoft",
		"windows-product": "Windows 10 Home",
	})
	if display != "Windows 10 Home" || ascii != "Windows 10" {
		t.Fatalf("wsl tiny %q %q", display, ascii)
	}

	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{"slitaz-release": "5.0\n"})
	if display != "SliTaz 5.0" || ascii != "SliTaz" {
		t.Fatalf("slitaz %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"synology-version": "productversion=\"7.2\"\nbuildnumber=\"64570\"\nsmallfixnumber=\"1\"\n",
	})
	if display != "Synology DSM 7.2-64570 Update 1" || ascii != "Synology" {
		t.Fatalf("synology %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro("QNX", "off", map[string]string{
		"qnx-props":   "device_os::QNX Neutrino\n",
		"qnx-version": "7.1.0",
	})
	if display != "QNX Neutrino (7.1.0)" || ascii != "QNX" {
		t.Fatalf("qnx %q %q", display, ascii)
	}
	display, ascii, _, _, _ = FormatDistro(OpenBSD, "off", map[string]string{
		"openbsd-version": "OpenBSD 7.5 (GENERIC.MP) #123: Fri Mar 1",
	})
	if display != "OpenBSD 7.5" || ascii != "OpenBSD" {
		t.Fatalf("openbsd %q %q", display, ascii)
	}
	display, _, _, _, _ = FormatDistro(Linux, "off", map[string]string{
		"redhat-release": "Red Hat Enterprise Server release 5\n",
	})
	if display != "Red Hat release 5" {
		t.Fatal(display)
	}
	line := "Oracle Solaris 11.4 SPARC"
	if got := solarisDisplay(line, "on"); got != "Oracle 11.4" {
		t.Fatal(got)
	}
	if got := solarisDisplay(line, "off"); got != "Oracle Solaris 11.4" {
		t.Fatal(got)
	}
}

func TestClassifyHarmonyAndFiwix(t *testing.T) {
	if got, known := ClassifyKernel("HarmonyOS"); got != Linux || !known {
		t.Fatalf("harmony %q %v", got, known)
	}
	if got, known := ClassifyKernel("Fiwix"); got != "Fiwix" || !known {
		t.Fatalf("fiwix %q %v", got, known)
	}
}
