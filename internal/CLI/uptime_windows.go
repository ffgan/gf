//go:build windows

package cli

import "syscall"

func uptimeFromPlatform() (int64, bool) {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetTickCount64")
	ms, _, _ := proc.Call()
	if ms == 0 {
		return 0, false
	}
	return int64(ms / 1000), true
}
