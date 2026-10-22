//go:build windows

package cli

import (
	"fmt"
	"syscall"
	"unsafe"
)

type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

func getMemoryWindows() (int64, int64, error) {
	var status memoryStatusEx
	status.Length = uint32(unsafe.Sizeof(status))
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx")
	r, _, err := proc.Call(uintptr(unsafe.Pointer(&status)))
	if r == 0 {
		if err != nil && err != syscall.Errno(0) {
			return 0, 0, err
		}
		return 0, 0, fmt.Errorf("GlobalMemoryStatusEx failed")
	}
	if status.TotalPhys == 0 {
		return 0, 0, fmt.Errorf("total physical memory is zero")
	}
	used := (status.TotalPhys - status.AvailPhys) / 1024
	total := status.TotalPhys / 1024
	return int64(used), int64(total), nil
}
