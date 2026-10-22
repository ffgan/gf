//go:build !windows

package cli

func windowsLongName() string {
	return ""
}

func readWindowsCurrentVersion() string {
	return ""
}

func readWindowsBIOS() string {
	return ""
}
