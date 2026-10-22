//go:build windows

package cli

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func localMount(mount string) string {
	if mount == "/" {
		drive := os.Getenv("SystemDrive")
		if drive == "" {
			return "C:"
		}
		return drive
	}
	return mount
}

func dfBytes(path string) (int64, int64, int, error) {
	if path == "" {
		return 0, 0, 0, fmt.Errorf("empty path")
	}
	if !strings.HasSuffix(path, `\`) && !strings.HasSuffix(path, "/") {
		path += `\`
	}
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, 0, 0, err
	}
	var freeAvail, total, totalFree uint64
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetDiskFreeSpaceExW")
	r, _, callErr := proc.Call(
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(&freeAvail)),
		uintptr(unsafe.Pointer(&total)),
		uintptr(unsafe.Pointer(&totalFree)),
	)
	if r == 0 {
		if callErr != nil && callErr != syscall.Errno(0) {
			return 0, 0, 0, callErr
		}
		return 0, 0, 0, fmt.Errorf("GetDiskFreeSpaceExW failed")
	}
	if total == 0 {
		return 0, 0, 0, fmt.Errorf("disk total is zero")
	}
	used := int64(total - totalFree)
	percent := int(used * 100 / int64(total))
	return used, int64(total), percent, nil
}
