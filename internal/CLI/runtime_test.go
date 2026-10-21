package cli

import "testing"

func TestPickHostSkipsEmptyDMI(t *testing.T) {
	got := PickHost(HostInput{
		DMIVendor:  "",
		DMIProduct: "To be filled by O.E.M.",
		Hypervisor: "KVM",
	})
	if got != "KVM" {
		t.Fatalf("got %q", got)
	}
}

func TestPickHostDeviceTreeNULAndAndroid(t *testing.T) {
	got := PickHost(HostInput{DeviceTree: "Raspberry Pi 4 Model B\x00extra", Hypervisor: "KVM"})
	if got != "Raspberry Pi 4 Model B" {
		t.Fatalf("dt: %q", got)
	}
	got = PickHost(HostInput{AndroidBrand: "Xiaomi", AndroidModel: "2201123G", DMIProduct: "Standard PC"})
	if got != "Xiaomi 2201123G" {
		t.Fatalf("android: %q", got)
	}
}

func TestPickHostWSLAndStandardPC(t *testing.T) {
	got := PickHost(HostInput{Kernel: "Linux version 5.15.90.1-microsoft-standard-WSL2", Hypervisor: "KVM"})
	if got != "Windows Subsystem for Linux" {
		t.Fatalf("wsl: %q", got)
	}
	got = PickHost(HostInput{DMIProduct: "Standard PC (i440FX + PIIX, 1996)"})
	if got != "KVM/QEMU (Standard PC (i440FX + PIIX, 1996))" {
		t.Fatalf("pc: %q", got)
	}
}

func TestRuntimeTags(t *testing.T) {
	none := RuntimeTags(RuntimeSig{ID: "ubuntu", Name: "Ubuntu", Kernel: "Linux version 6.12.94+", Exists: func(string) bool { return false }})
	if len(none) != 0 {
		t.Fatalf("plain ubuntu: %v", none)
	}
	wsl := RuntimeTags(RuntimeSig{ID: "ubuntu", Name: "Ubuntu", Kernel: "5.15.90.1-microsoft-standard-WSL2"})
	if len(wsl) != 1 || wsl[0] != "WSL2" {
		t.Fatalf("wsl: %v", wsl)
	}
	wsl1 := RuntimeTags(RuntimeSig{ID: "ubuntu", Name: "Ubuntu", Kernel: "4.4.0-19041-Microsoft"})
	if len(wsl1) != 1 || wsl1[0] != "WSL1" {
		t.Fatalf("wsl1: %v", wsl1)
	}
	docker := RuntimeTags(RuntimeSig{
		ID: "debian", Name: "Debian",
		Exists: func(p string) bool { return p == "/.dockerenv" || p == "/etc/pve" },
	})
	if len(docker) != 2 || docker[0] != "Docker" || docker[1] != "Proxmox" {
		t.Fatalf("docker+pve: %v", docker)
	}
	termux := RuntimeTags(RuntimeSig{Env: map[string]string{"TERMUX_VERSION": "0.118", "ANDROID_ROOT": "/system"}})
	if len(termux) != 1 || termux[0] != "Termux" {
		t.Fatalf("termux: %v", termux)
	}
}

func TestLogoKey(t *testing.T) {
	cases := map[[2]string]string{
		{"UOS", "uos"}:                 "Uos",
		{"Kylin", "kylin"}:             "Ubuntu Kylin",
		{"openKylin", "openkylin"}:     "openKylin",
		{"openEuler", "openeuler"}:     "openEuler",
		{"HarmonyOS", "harmonyos"}:     "Android",
		{"OpenHarmony", "openharmony"}: "Android",
		{"HyperOS", "hyperos"}:         "Android",
		{"Anolis OS", "anolis"}:        "CentOS",
		{"OpenCloudOS", "opencloudos"}: "CentOS",
		{"Loongnix", "loongnix"}:       "Debian",
		{"OpenWrt", "openwrt"}:         "Wrt",
		{"Raspberry Pi OS", "debian"}:  "Raspbian",
		{"Rocky Linux", "rocky"}:       "rocky",
		{"Ubuntu", "ubuntu"}:           "",
		{"Arch Linux", "archlinux"}:    "Arch",
		{"openSUSE", "opensuse-leap"}:  "openSUSE Leap",
		{"GhostBSD", "ghostbsd"}:       "GhostBSD",
	}
	for in, want := range cases {
		if got := LogoKey(in[0], in[1]); got != want {
			t.Errorf("%v: got %q want %q", in, got, want)
		}
	}
}

func TestAnnotateDistro(t *testing.T) {
	if got := AnnotateDistro("Ubuntu 24.04", []string{"WSL2"}); got != "Ubuntu 24.04 (WSL2)" {
		t.Fatalf("%q", got)
	}
	if got := AnnotateDistro("Ubuntu 24.04", nil); got != "Ubuntu 24.04" {
		t.Fatalf("%q", got)
	}
}

func TestIPhoneName(t *testing.T) {
	if iphoneName("iPhone17,1") != "iPhone 16 Pro" {
		t.Fatal(iphoneName("iPhone17,1"))
	}
	if iphoneName("iPhone99,1") != "iPhone99,1" {
		t.Fatal(iphoneName("iPhone99,1"))
	}
}
