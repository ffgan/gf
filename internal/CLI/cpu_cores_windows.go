//go:build windows

package cli

import (
	"syscall"
	"unsafe"
)

func windowsCoreLabel() string {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetLogicalProcessorInformationEx")
	if err := proc.Find(); err != nil {
		return ""
	}
	var n uint32
	// RelationProcessorCore. The first call sizes the buffer.
	_, _, _ = proc.Call(0, 0, uintptr(unsafe.Pointer(&n)))
	if n == 0 || n > 1<<20 {
		return ""
	}
	buf := make([]byte, n)
	r, _, _ := proc.Call(0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n)))
	if r == 0 || int(n) > len(buf) {
		return ""
	}
	return coreLabelFromBuffer(buf[:n])
}
