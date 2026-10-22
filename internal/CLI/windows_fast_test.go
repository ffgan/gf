package cli

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleCurrentVersion = `
HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion
    ProductName    REG_SZ    Windows 10 Pro
    CurrentBuildNumber    REG_SZ    22631
    CurrentMajorVersionNumber    REG_DWORD    0xa
    CurrentMinorVersionNumber    REG_DWORD    0x0
    UBR    REG_DWORD    0xf28
    CurrentVersion    REG_SZ    6.3
`

func TestWindowsRegistryParsers(t *testing.T) {
	if got := windowsDistroFromReg(sampleCurrentVersion); got != "Windows 11 Pro" {
		t.Fatalf("distro %q", got)
	}
	if got := windowsKernelVersion(sampleCurrentVersion); got != "10.0.22631.3880" {
		t.Fatalf("kernel %q", got)
	}
	bios := `
HKEY_LOCAL_MACHINE\HARDWARE\DESCRIPTION\System\BIOS
    SystemManufacturer    REG_SZ    Dell Inc.
    SystemProductName    REG_SZ    Latitude 5420
`
	if got := windowsHardwareModel(bios); got != "Dell Inc. Latitude 5420" {
		t.Fatalf("model %q", got)
	}
	uninstall := `
HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall

HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Git
    DisplayName    REG_SZ    Git

HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Git

HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\VSCode
`
	if got := countRegSubkeys(uninstall); got != 2 {
		t.Fatalf("subkeys %d", got)
	}
}

func TestCountCargoInstalls(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CARGO_HOME", dir)
	if _, ok := countCargoInstalls(); ok {
		t.Fatal("missing index should not count")
	}
	body := `{"installs":{"ripgrep 14.0.0":{"bins":["rg"]},"fd-find 9.0.0":{"bins":["fd"]}}}`
	if err := os.WriteFile(filepath.Join(dir, ".crates2.json"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	n, ok := countCargoInstalls()
	if !ok || n != 2 {
		t.Fatalf("crates2 %d %v", n, ok)
	}
}

func TestCountCargoToml(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CARGO_HOME", dir)
	body := "[v1]\n\"ripgrep 14.0.0\" = [\"rg\"]\n"
	if err := os.WriteFile(filepath.Join(dir, ".crates.toml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	n, ok := countCargoInstalls()
	if !ok || n != 1 {
		t.Fatalf("crates.toml %d %v", n, ok)
	}
}
