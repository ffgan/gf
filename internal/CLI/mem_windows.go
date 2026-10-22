//go:build windows

package cli

import (
	"encoding/binary"
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

func getSwapWindows() (int64, int64, error) {
	if used, total, ok := pagefileBytes(); ok && total > 0 {
		return used, total, nil
	}
	var status memoryStatusEx
	status.Length = uint32(unsafe.Sizeof(status))
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx")
	r, _, err := proc.Call(uintptr(unsafe.Pointer(&status)))
	if r == 0 || status.TotalPageFile < status.TotalPhys {
		if err != nil && err != syscall.Errno(0) {
			return 0, 0, err
		}
		return 0, 0, fmt.Errorf("page file is unavailable")
	}
	total := int64(status.TotalPageFile - status.TotalPhys)
	used := int64((status.TotalPageFile - status.AvailPageFile) - (status.TotalPhys - status.AvailPhys))
	if used < 0 {
		used = 0
	}
	return used, total, nil
}

func pagefileBytes() (int64, int64, bool) {
	ntdll := syscall.NewLazyDLL("ntdll.dll")
	query := ntdll.NewProc("NtQuerySystemInformation")
	if query.Find() != nil {
		return 0, 0, false
	}
	buf := make([]byte, 8192)
	var need uint32
	// SystemPagefileInformation = 18
	status, _, _ := query.Call(18, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), uintptr(unsafe.Pointer(&need)))
	if status != 0 {
		return 0, 0, false
	}
	page := uint64(systemPageSize())
	var used, total uint64
	off := 0
	for off+16 <= len(buf) {
		next := int(binary.LittleEndian.Uint32(buf[off:]))
		// mingw SYSTEM_PAGEFILE_INFORMATION, the layout fastfetch reads:
		// NextEntryOffset, TotalUsed, PeakUsed, CurrentSize.
		pagesUsed := binary.LittleEndian.Uint32(buf[off+4:])
		pagesTotal := binary.LittleEndian.Uint32(buf[off+12:])
		used += uint64(pagesUsed) * page
		total += uint64(pagesTotal) * page
		if next < 16 || off+next > len(buf) {
			break
		}
		off += next
	}
	return int64(used), int64(total), total > 0
}

func systemPageSize() uint32 {
	var info [64]byte
	syscall.NewLazyDLL("kernel32.dll").NewProc("GetSystemInfo").Call(uintptr(unsafe.Pointer(&info[0])))
	n := binary.LittleEndian.Uint32(info[4:8])
	if n == 0 {
		return 4096
	}
	return n
}
