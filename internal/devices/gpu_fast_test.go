package dev

import "testing"

func TestFormatWindowsGPUs(t *testing.T) {
	both := `
    DriverDesc    REG_SZ    Microsoft Basic Display Adapter
    DriverDesc    REG_SZ    NVIDIA GeForce RTX 3060
    DriverDesc    REG_SZ    NVIDIA GeForce RTX 3060
`
	if got := formatWindowsGPUs(both); got != "NVIDIA GeForce RTX 3060" {
		t.Fatalf("got %q", got)
	}
	basic := `
    DriverDesc    REG_SZ    Microsoft Basic Display Adapter
`
	if got := formatWindowsGPUs(basic); got != "Microsoft Basic Display Adapter" {
		t.Fatalf("got %q", got)
	}
}
