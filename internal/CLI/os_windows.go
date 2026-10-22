//go:build windows

package cli

import (
	"syscall"
	"unsafe"
)

func windowsLongName() string {
	proc := syscall.NewLazyDLL("winbrand.dll").NewProc("BrandingFormatString")
	if err := proc.Find(); err != nil {
		return ""
	}
	format, err := syscall.UTF16PtrFromString("%WINDOWS_LONG%")
	if err != nil {
		return ""
	}
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(format)))
	if ret == 0 {
		return ""
	}
	name := utf16PtrString(ret)
	syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalFree").Call(ret)
	return name
}

func utf16PtrString(p uintptr) string {
	if p == 0 {
		return ""
	}
	// lstrcpynW writes into Go memory, so the foreign pointer stays a uintptr.
	dst := make([]uint16, 512)
	r, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("lstrcpynW").Call(
		uintptr(unsafe.Pointer(&dst[0])),
		p,
		uintptr(len(dst)),
	)
	if r == 0 {
		return ""
	}
	return syscall.UTF16ToString(dst)
}
