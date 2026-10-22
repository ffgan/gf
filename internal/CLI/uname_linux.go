//go:build linux

package cli

import (
	"sync"
	"syscall"
)

var (
	utsOnce sync.Once
	uts     syscall.Utsname
	utsErr  error
)

func loadUts() (*syscall.Utsname, bool) {
	utsOnce.Do(func() {
		utsErr = syscall.Uname(&uts)
	})
	if utsErr != nil {
		return nil, false
	}
	return &uts, true
}

// syscall.Utsname fields are [65]int8 on amd64 and [65]uint8 on arm,
// ppc64, riscv64, and s390x. Android reuses those linux types.
func utsString[T int8 | uint8](a [65]T) string {
	b := make([]byte, 0, 65)
	for _, c := range a {
		if c == 0 {
			break
		}
		b = append(b, byte(c))
	}
	return string(b)
}

// unameFast answers the flags gf actually uses from the uname syscall.
// Spawning the uname binary costs a process just to read a struct the
// kernel already exposes.
func unameFast(command string) (string, bool) {
	u, ok := loadUts()
	if !ok {
		return "", false
	}
	switch command {
	case "-s":
		return utsString(u.Sysname), true
	case "-r":
		return utsString(u.Release), true
	case "-m":
		return utsString(u.Machine), true
	case "-n":
		return utsString(u.Nodename), true
	case "-v":
		return utsString(u.Version), true
	case "-srm":
		return utsString(u.Sysname) + " " + utsString(u.Release) + " " + utsString(u.Machine), true
	case "-o":
		if utsString(u.Sysname) == "Linux" {
			return "GNU/Linux", true
		}
		return "", false
	default:
		return "", false
	}
}
