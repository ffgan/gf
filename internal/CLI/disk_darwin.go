//go:build darwin

package cli

import "os"

// applyContainerUsage rewrites the sealed system volume. /System/Volumes/Data
// is the other half of the APFS volume group; when it exists, df's Used
// column for / is not the storage bar.
func applyContainerUsage(path string, used, total, avail int64) (int64, int, bool) {
	data := false
	if path == "/" {
		fi, err := os.Stat("/System/Volumes/Data")
		data = err == nil && fi.IsDir()
	}
	adjusted, percent := containerUsage(path, used, total, avail, data)
	return adjusted, percent, true
}
