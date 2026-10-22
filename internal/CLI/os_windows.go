//go:build windows

package cli

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/ffgan/gf/internal/utils"
)

func readWindowsCurrentVersion() string {
	return regText(utils.HKLM, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		[]string{"ProductName", "DisplayVersion", "CSDVersion", "ReleaseId", "CurrentBuildNumber", "CurrentVersion"},
		[]string{"CurrentMajorVersionNumber", "CurrentMinorVersionNumber", "UBR"},
	)
}

func readWindowsBIOS() string {
	return regText(utils.HKLM, `HARDWARE\DESCRIPTION\System\BIOS`,
		[]string{"SystemManufacturer", "SystemProductName"},
		nil,
	)
}

func regText(root uintptr, path string, strings_, dwords []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "HKEY_%s\\%s\n", hiveName(root), path)
	wrote := false
	for _, name := range strings_ {
		v, ok := utils.RegString(root, path, name)
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "    %s    REG_SZ    %s\n", name, v)
		wrote = true
	}
	for _, name := range dwords {
		v, ok := utils.RegUint(root, path, name)
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "    %s    REG_DWORD    0x%x\n", name, v)
		wrote = true
	}
	if !wrote {
		return ""
	}
	return b.String()
}

func hiveName(root uintptr) string {
	if root == utils.HKCU {
		return "CURRENT_USER"
	}
	return "LOCAL_MACHINE"
}

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
