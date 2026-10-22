//go:build windows

package dev

import (
	"syscall"
	"testing"
)

// TestWindowsDisplayAPIs calls the Win32 path behind the Resolution line.
// GetDeviceCaps lives in gdi32. Resolving it from user32 panics in
// LazyProc.Call, which is what killed gf.exe on Windows.
func TestWindowsDisplayAPIs(t *testing.T) {
	want := logPixelsX(t)
	got := systemDPI(syscall.NewLazyDLL("user32.dll"))
	if got != want {
		t.Fatalf("systemDPI() = %d, gdi32 GetDeviceCaps(LOGPIXELSX) = %d", got, want)
	}

	// Same call the binary makes. An empty result is fine when the session
	// has no monitor; a panic is not.
	_ = getResolutionWindows(true)
}

func logPixelsX(t *testing.T) int {
	t.Helper()
	user32 := syscall.NewLazyDLL("user32.dll")
	getDC := user32.NewProc("GetDC")
	release := user32.NewProc("ReleaseDC")
	caps := syscall.NewLazyDLL("gdi32.dll").NewProc("GetDeviceCaps")
	if err := caps.Find(); err != nil {
		t.Fatalf("gdi32 GetDeviceCaps: %v", err)
	}
	if err := user32.NewProc("GetDeviceCaps").Find(); err == nil {
		t.Fatal("user32 exported GetDeviceCaps; the Resolution line must keep using gdi32")
	}
	if err := getDC.Find(); err != nil {
		t.Fatal(err)
	}
	if err := release.Find(); err != nil {
		t.Fatal(err)
	}
	dc, _, _ := getDC.Call(0)
	if dc == 0 {
		t.Fatal("GetDC(0) returned 0")
	}
	defer release.Call(0, dc)
	dpi, _, _ := caps.Call(dc, 88) // LOGPIXELSX
	if dpi == 0 {
		t.Fatal("GetDeviceCaps returned 0")
	}
	return int(dpi)
}
