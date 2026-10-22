//go:build windows

package dev

import (
	"net"
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

func platformRouteIface() string {
	proc := syscall.NewLazyDLL("iphlpapi.dll").NewProc("GetBestInterface")
	if proc.Find() != nil {
		return ""
	}
	var idx uint32
	// 8.8.8.8 in network byte order, used only to pick the default route.
	dest := uint32(8)<<24 | uint32(8)<<16 | uint32(8)<<8 | 8
	r, _, _ := proc.Call(uintptr(dest), uintptr(unsafe.Pointer(&idx)))
	if r != 0 || idx == 0 {
		return ""
	}
	iface, err := net.InterfaceByIndex(int(idx))
	if err != nil {
		return ""
	}
	return iface.Name
}
