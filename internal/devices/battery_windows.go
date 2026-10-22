//go:build windows

package dev

import (
	"encoding/binary"
	"syscall"
	"unsafe"
)

func getBatteryWindows() string {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetSystemPowerStatus")
	if proc.Find() != nil {
		return ""
	}
	var raw [12]byte
	r, _, _ := proc.Call(uintptr(unsafe.Pointer(&raw[0])))
	if r == 0 {
		return ""
	}
	life := binary.LittleEndian.Uint32(raw[4:8])
	return batteryFromStatus(raw[0], raw[1], raw[2], life)
}
