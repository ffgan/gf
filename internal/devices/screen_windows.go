//go:build windows

package dev

import (
	"fmt"
	"syscall"
)

func getResolutionWindows() string {
	metrics := syscall.NewLazyDLL("user32.dll").NewProc("GetSystemMetrics")
	cx, _, _ := metrics.Call(0) // SM_CXSCREEN
	cy, _, _ := metrics.Call(1) // SM_CYSCREEN
	if cx == 0 || cy == 0 {
		return ""
	}
	return fmt.Sprintf("%dx%d", cx, cy)
}
