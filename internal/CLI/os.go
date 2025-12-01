package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetOS(kernelName, darwin_name string) string {
	var osName string

	switch {
	case kernelName == Darwin:
		// TODO: https://github.com/ffgan/hyfetch/blob/master/neofetch#L6384

	case kernelName == SunOS:
		out, err := utils.CommandOutput("uname", "-o")
		if err == nil && strings.Contains(string(out), illumos) {
			osName = illumos
		} else {
			osName = Solaris
		}

	case kernelName == Haiku:
		osName = Haiku

	case kernelName == MINIX:
		osName = MINIX

	case kernelName == AIX:
		osName = AIX

	case strings.HasPrefix(kernelName, IRIX):
		osName = IRIX

	case kernelName == FreeMiNT:
		osName = FreeMiNT

	case kernelName == Interix:
		osName = Interix

	case kernelName == Ironclad:
		osName = Ironclad

	case kernelName == OSF1:
		osName = digitalUNIX

	case kernelName == Linux || strings.HasPrefix(kernelName, GNU):
		osName = Linux

	case strings.HasSuffix(kernelName, BSD) ||
		kernelName == DragonFly ||
		kernelName == Bitrig:
		osName = BSD

	case strings.HasPrefix(kernelName, CYGWIN) ||
		strings.HasPrefix(kernelName, MSYS) ||
		strings.HasPrefix(kernelName, MINGW) ||
		kernelName == Windows_NT:
		osName = Windows

	default:
		fmt.Fprintf(os.Stderr, "Unknown OS detected: '%s', aborting...\n", kernelName)
		fmt.Fprintln(os.Stderr, "Open an issue on GitHub to add support for your OS.")
		os.Exit(1)
	}

	return osName
}
