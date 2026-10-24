//go:build !windows && !darwin

package cli

func applyContainerUsage(string, int64, int64, int64) (int64, int, bool) {
	return 0, 0, false
}
