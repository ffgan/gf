//go:build linux

package gui

import (
	"net"
	"syscall"
	"time"
)

// waylandPeerPID asks the compositor for its pid over the wayland socket.
// fastfetch does the same SO_PEERCRED lookup instead of spawning lsof,
// which walks every file descriptor on the system.
func waylandPeerPID(path string) int {
	conn, err := net.DialTimeout("unix", path, 50*time.Millisecond)
	if err != nil {
		return 0
	}
	defer conn.Close()
	uc, ok := conn.(*net.UnixConn)
	if !ok {
		return 0
	}
	raw, err := uc.SyscallConn()
	if err != nil {
		return 0
	}
	var pid int
	_ = raw.Control(func(fd uintptr) {
		cred, err := syscall.GetsockoptUcred(int(fd), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
		if err == nil && cred != nil {
			pid = int(cred.Pid)
		}
	})
	return pid
}
