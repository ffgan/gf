//go:build !windows

package cli

func uptimeFromPlatform() (int64, bool) {
	return 0, false
}
