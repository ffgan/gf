package dev

import "testing"

func TestFormatWindowsGPUs(t *testing.T) {
	both := `
    DriverDesc    REG_SZ    Microsoft Basic Display Adapter
    DriverDesc    REG_SZ    NVIDIA GeForce RTX 3060
    DriverDesc    REG_SZ    NVIDIA GeForce RTX 3060
`
	if got := formatWindowsGPUs(both); got != "NVIDIA GeForce RTX 3060 [Discrete]" {
		t.Fatalf("got %q", got)
	}
	basic := `
    DriverDesc    REG_SZ    Microsoft Basic Display Adapter
`
	if got := formatWindowsGPUs(basic); got != "Microsoft Basic Display Adapter" {
		t.Fatalf("got %q", got)
	}
	intel := `
HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}\0000
    DriverDesc    REG_SZ    Intel(R) Graphics
    HardwareInformation.qwMemorySize    REG_QWORD    0x8000000

HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}\0001
    DriverDesc    REG_SZ    MuMu Virtual Display Adapter
`
	got := formatWindowsGPUs(intel)
	want := "Intel(R) Graphics (128.00 MiB) [Integrated]\nMuMu Virtual Display Adapter"
	if got != want {
		t.Fatalf("got %q", got)
	}
}
