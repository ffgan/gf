//go:build !windows

package cli

import "fmt"

func getMemoryWindows() (int64, int64, error) {
	return 0, 0, fmt.Errorf("Windows support requires additional implementation")
}
