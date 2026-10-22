//go:build linux

package cli

import (
	"syscall"
	"unsafe"
)

func ioctlSize(fd int) (int, int, bool) {
	var ws struct {
		Row, Col, Xpixel, Ypixel uint16
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(0x5413), uintptr(unsafe.Pointer(&ws)))
	if errno != 0 || ws.Col == 0 || ws.Row == 0 {
		return 0, 0, false
	}
	return int(ws.Col), int(ws.Row), true
}
