//go:build windows

package dev

import (
	"encoding/binary"
	"strings"
	"syscall"
	"unsafe"
)

func formatLocalIP(name, cidr string) string {
	if cidr == "" {
		return ""
	}
	if name == "" {
		return cidr
	}
	return "(" + name + ") " + cidr
}

func platformLocalIP(ifaceName string) (string, bool) {
	return windowsLocalIP(ifaceName), true
}

func platformRouteIface() string {
	return ""
}

func windowsLocalIP(requested string) string {
	requested = strings.Trim(requested, "() '\"")
	rows := ipAddrRows()
	if requested != "" && !strings.EqualFold(requested, "auto") {
		for _, row := range rows {
			alias := aliasForIndex(row.idx)
			if strings.EqualFold(alias, requested) {
				return formatLocalIP(alias, row.cidr)
			}
		}
		return ""
	}
	if idx, ok := bestInterfaceIndex(); ok {
		for _, row := range rows {
			if row.idx == idx {
				return formatLocalIP(aliasForIndex(idx), row.cidr)
			}
		}
	}
	for _, row := range rows {
		return formatLocalIP(aliasForIndex(row.idx), row.cidr)
	}
	return ""
}

type ipRow struct {
	idx  uint32
	cidr string
}

func ipAddrRows() []ipRow {
	iphlpapi := syscall.NewLazyDLL("iphlpapi.dll")
	proc := iphlpapi.NewProc("GetIpAddrTable")
	if proc.Find() != nil {
		return nil
	}
	size := uint32(0)
	proc.Call(0, uintptr(unsafe.Pointer(&size)), 0)
	if size < 4 || size > 1<<20 {
		return nil
	}
	buf := make([]byte, size)
	r, _, _ := proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)), 0)
	if r != 0 || len(buf) < 4 {
		return nil
	}
	n := binary.LittleEndian.Uint32(buf[:4])
	const rowSize = 24
	var rows []ipRow
	for i := uint32(0); i < n; i++ {
		off := 4 + int(i)*rowSize
		if off+12 > len(buf) {
			break
		}
		addr := binary.LittleEndian.Uint32(buf[off:])
		idx := binary.LittleEndian.Uint32(buf[off+4:])
		mask := binary.LittleEndian.Uint32(buf[off+8:])
		cidr := ipv4Prefix(addr, mask)
		if cidr == "" {
			continue
		}
		rows = append(rows, ipRow{idx: idx, cidr: cidr})
	}
	return rows
}

func bestInterfaceIndex() (uint32, bool) {
	proc := syscall.NewLazyDLL("iphlpapi.dll").NewProc("GetBestInterface")
	if proc.Find() != nil {
		return 0, false
	}
	var idx uint32
	// 8.8.8.8, network byte order. Every octet is 8, so endianness does not matter.
	dest := uint32(8)<<24 | uint32(8)<<16 | uint32(8)<<8 | 8
	r, _, _ := proc.Call(uintptr(dest), uintptr(unsafe.Pointer(&idx)))
	if r != 0 || idx == 0 {
		return 0, false
	}
	return idx, true
}

func aliasForIndex(idx uint32) string {
	iphlpapi := syscall.NewLazyDLL("iphlpapi.dll")
	toLuid := iphlpapi.NewProc("ConvertInterfaceIndexToLuid")
	toAlias := iphlpapi.NewProc("ConvertInterfaceLuidToAlias")
	if toLuid.Find() != nil || toAlias.Find() != nil {
		return ""
	}
	var luid uint64
	r, _, _ := toLuid.Call(uintptr(idx), uintptr(unsafe.Pointer(&luid)))
	if r != 0 {
		return ""
	}
	buf := make([]uint16, 256)
	r, _, _ = toAlias.Call(uintptr(unsafe.Pointer(&luid)), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if r != 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}
