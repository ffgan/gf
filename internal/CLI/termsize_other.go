//go:build !linux

package cli

func ioctlSize(fd int) (int, int, bool) {
	return 0, 0, false
}
