package cli

import (
	"fmt"
	"os"
)

// GetTermSize returns "columns x rows" for the controlling terminal.
func GetTermSize() string {
	for _, file := range []*os.File{os.Stdout, os.Stderr} {
		if cols, rows, ok := ioctlSize(int(file.Fd())); ok {
			return fmt.Sprintf("%dx%d", cols, rows)
		}
	}
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return ""
	}
	defer tty.Close()
	if cols, rows, ok := ioctlSize(int(tty.Fd())); ok {
		return fmt.Sprintf("%dx%d", cols, rows)
	}
	return ""
}
