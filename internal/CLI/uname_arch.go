package cli

// unameMachine maps a Go architecture onto the name uname -m prints.
func unameMachine(goarch string) string {
	switch goarch {
	case "amd64":
		return "x86_64"
	case "386":
		return "i686"
	case "arm64":
		return "aarch64"
	case "arm":
		return "armv7l"
	default:
		return goarch
	}
}
