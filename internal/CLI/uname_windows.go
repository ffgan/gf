//go:build windows

package cli

import (
	"os"
	"runtime"
)

// unameFast answers the flags gf asks for without the uname binary.
// PowerShell and cmd do not ship uname, and UName used to exit when it
// was missing.
func unameFast(command string) (string, bool) {
	switch command {
	case "-s", "-o":
		return "Windows_NT", true
	case "-r", "-v":
		ver := windowsKernelVersion(windowsCurrentVersion())
		if ver == "" {
			ver = "0"
		}
		return ver, true
	case "-m", "-p":
		return unameMachine(runtime.GOARCH), true
	case "-n":
		name, err := os.Hostname()
		if err != nil || name == "" {
			return "localhost", true
		}
		return name, true
	case "-srm":
		sys, _ := unameFast("-s")
		rel, _ := unameFast("-r")
		mach, _ := unameFast("-m")
		return sys + " " + rel + " " + mach, true
	default:
		return "", false
	}
}
