//go:build windows

package cli

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/ffgan/gf/internal/utils"
)

type winVolume struct {
	Mount   string
	Used    int64
	Total   int64
	Percent int
	FS      string
}

func windowsDiskReport(mounts []string, report DiskReport) (string, bool) {
	vols := listWindowsVolumes(mounts)
	var lines []string
	for _, vol := range vols {
		if vol.Total <= 0 {
			continue
		}
		text := flexibleBytes(vol.Used) + " / " + flexibleBytes(vol.Total)
		if report.Percent != utils.OFF {
			text += fmt.Sprintf(" (%d%%)", vol.Percent)
		}
		if label := diskLabel(vol.Mount, report.Subtitle); label != "" {
			text = label + " " + text
		}
		if vol.FS != "" {
			text += " - " + vol.FS
		}
		lines = append(lines, text)
	}
	return strings.Join(lines, "\n"), true
}

func listWindowsVolumes(mounts []string) []winVolume {
	if len(mounts) == 0 || (len(mounts) == 1 && mounts[0] == "/") {
		return fixedWindowsVolumes()
	}
	var vols []winVolume
	for _, mount := range mounts {
		mount = localMount(mount)
		if vol, ok := volumeInfo(mount); ok {
			vols = append(vols, vol)
		}
	}
	return vols
}

func fixedWindowsVolumes() []winVolume {
	kernel := syscall.NewLazyDLL("kernel32.dll")
	mask, _, _ := kernel.NewProc("GetLogicalDrives").Call()
	var vols []winVolume
	for i := 0; i < 26; i++ {
		if mask&(1<<uint(i)) == 0 {
			continue
		}
		root := string(rune('A'+i)) + `:\`
		if driveType(root) != 3 && driveType(root) != 2 {
			continue
		}
		if vol, ok := volumeInfo(root); ok {
			vols = append(vols, vol)
		}
	}
	return vols
}

func driveType(root string) uintptr {
	p, err := syscall.UTF16PtrFromString(root)
	if err != nil {
		return 0
	}
	r, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("GetDriveTypeW").Call(uintptr(unsafe.Pointer(p)))
	return r
}

func volumeInfo(root string) (winVolume, bool) {
	if !strings.HasSuffix(root, `\`) && !strings.HasSuffix(root, "/") {
		root += `\`
	}
	used, total, percent, err := dfBytes(strings.TrimRight(root, `\`))
	if err != nil || total <= 0 {
		return winVolume{}, false
	}
	return winVolume{
		Mount:   root,
		Used:    used,
		Total:   total,
		Percent: percent,
		FS:      volumeFS(root),
	}, true
}

func volumeFS(root string) string {
	p, err := syscall.UTF16PtrFromString(root)
	if err != nil {
		return ""
	}
	fs := make([]uint16, 64)
	r, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("GetVolumeInformationW").Call(
		uintptr(unsafe.Pointer(p)),
		0, 0, 0, 0, 0,
		uintptr(unsafe.Pointer(&fs[0])),
		uintptr(len(fs)),
	)
	if r == 0 {
		return ""
	}
	return syscall.UTF16ToString(fs)
}

func localMount(mount string) string {
	if mount == "/" {
		drive := os.Getenv("SystemDrive")
		if drive == "" {
			return "C:"
		}
		return drive
	}
	return mount
}

func dfBytes(path string) (int64, int64, int, error) {
	if path == "" {
		return 0, 0, 0, fmt.Errorf("empty path")
	}
	if !strings.HasSuffix(path, `\`) && !strings.HasSuffix(path, "/") {
		path += `\`
	}
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, 0, 0, err
	}
	var freeAvail, total, totalFree uint64
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetDiskFreeSpaceExW")
	r, _, callErr := proc.Call(
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(&freeAvail)),
		uintptr(unsafe.Pointer(&total)),
		uintptr(unsafe.Pointer(&totalFree)),
	)
	if r == 0 {
		if callErr != nil && callErr != syscall.Errno(0) {
			return 0, 0, 0, callErr
		}
		return 0, 0, 0, fmt.Errorf("GetDiskFreeSpaceExW failed")
	}
	if total == 0 {
		return 0, 0, 0, fmt.Errorf("disk total is zero")
	}
	used := int64(total - totalFree)
	percent := int(used * 100 / int64(total))
	return used, int64(total), percent, nil
}
