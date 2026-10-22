//go:build windows

package dev

import (
	"encoding/binary"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/ffgan/gf/internal/utils"
)

func getResolutionWindows(refresh bool) string {
	if lines := queryDisplays(refresh); lines != "" {
		return lines
	}
	user32 := syscall.NewLazyDLL("user32.dll")
	metrics := user32.NewProc("GetSystemMetrics")
	cx, _, _ := metrics.Call(0)
	cy, _, _ := metrics.Call(1)
	if cx == 0 || cy == 0 {
		return ""
	}
	return formatWinDisplay(winDisplay{W: int(cx), H: int(cy), ShowHz: refresh})
}

func queryDisplays(refresh bool) string {
	user32 := syscall.NewLazyDLL("user32.dll")
	sizes := user32.NewProc("GetDisplayConfigBufferSizes")
	query := user32.NewProc("QueryDisplayConfig")
	info := user32.NewProc("DisplayConfigGetDeviceInfo")
	if sizes.Find() != nil || query.Find() != nil {
		return ""
	}
	var pathCount, modeCount uint32
	if r, _, _ := sizes.Call(2, uintptr(unsafe.Pointer(&pathCount)), uintptr(unsafe.Pointer(&modeCount))); r != 0 {
		return ""
	}
	if pathCount == 0 || modeCount == 0 || pathCount > 128 || modeCount > 256 {
		return ""
	}
	paths := make([]byte, int(pathCount)*72)
	modes := make([]byte, int(modeCount)*64)
	r, _, _ := query.Call(
		2,
		uintptr(unsafe.Pointer(&pathCount)),
		uintptr(unsafe.Pointer(&paths[0])),
		uintptr(unsafe.Pointer(&modeCount)),
		uintptr(unsafe.Pointer(&modes[0])),
		0,
	)
	if r != 0 {
		return ""
	}
	sysDPI := systemDPI(user32)
	var lines []string
	for i := 0; i < int(pathCount); i++ {
		path := paths[i*72 : (i+1)*72]
		idx := binary.LittleEndian.Uint32(path[12:16])
		if idx == 0xFFFFFFFF || int(idx) >= int(modeCount) {
			idx = binary.LittleEndian.Uint32(path[12:16]) & 0xFFFF
		}
		if int(idx) >= int(modeCount) {
			continue
		}
		mode := modes[int(idx)*64 : int(idx)*64+64]
		if binary.LittleEndian.Uint32(mode[0:4]) != 1 {
			continue
		}
		w := int(binary.LittleEndian.Uint32(mode[16:20]))
		h := int(binary.LittleEndian.Uint32(mode[20:24]))
		x := int32(binary.LittleEndian.Uint32(mode[28:32]))
		y := int32(binary.LittleEndian.Uint32(mode[32:36]))
		rot := binary.LittleEndian.Uint32(path[40:44])
		if rot == 2 || rot == 4 {
			w, h = h, w
		}
		hz := 0.0
		den := binary.LittleEndian.Uint32(path[52:56])
		if den != 0 {
			hz = float64(binary.LittleEndian.Uint32(path[48:52])) / float64(den)
		}
		tech := binary.LittleEndian.Uint32(path[36:40])
		name, wmm, hmm := targetDisplay(info, path)
		dpi := monitorDPI(user32, x, y)
		if dpi == 0 {
			dpi = sysDPI
		}
		line := formatWinDisplay(winDisplay{
			Name:    name,
			W:       w,
			H:       h,
			Hz:      hz,
			DPI:     dpi,
			Wmm:     wmm,
			Hmm:     hmm,
			Kind:    displayKind(tech),
			Primary: x == 0 && y == 0,
			ShowHz:  refresh,
		})
		if line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func displayKind(tech uint32) string {
	switch tech {
	case 0xFFFFFFFF:
		return ""
	case 0x80000000, 11, 13:
		return "Built-in"
	default:
		return "External"
	}
}

func targetDisplay(info *syscall.LazyProc, path []byte) (string, int, int) {
	if info.Find() != nil {
		return "", 0, 0
	}
	buf := make([]byte, 420)
	binary.LittleEndian.PutUint32(buf[0:4], 2) // DISPLAYCONFIG_DEVICE_INFO_GET_TARGET_NAME
	binary.LittleEndian.PutUint32(buf[4:8], 420)
	copy(buf[8:16], path[20:28]) // target adapter LUID
	copy(buf[16:20], path[28:32])
	if r, _, _ := info.Call(uintptr(unsafe.Pointer(&buf[0]))); r != 0 {
		return "", 0, 0
	}
	friendly := utf16BytesToString(buf[36:164])
	devicePath := utf16BytesToString(buf[164:420])
	name, wmm, hmm := "", 0, 0
	if key := edidKeyFromDevicePath(devicePath); key != "" {
		if edid, ok := utils.RegBytes(utils.HKLM, key, "EDID"); ok {
			name, wmm, hmm = edidNameSize(edid)
		}
	}
	if name == "" && !strings.ContainsAny(friendly, `\#`) {
		name = friendly
	}
	return name, wmm, hmm
}

func monitorDPI(user32 *syscall.LazyDLL, x, y int32) int {
	// POINT is passed by value. On 64-bit it occupies one register.
	if unsafe.Sizeof(uintptr(0)) != 8 {
		return 0
	}
	mon := user32.NewProc("MonitorFromPoint")
	if mon.Find() != nil {
		return 0
	}
	var pt struct{ x, y int32 }
	pt.x, pt.y = x, y
	h, _, _ := mon.Call(*(*uintptr)(unsafe.Pointer(&pt)), 0)
	if h == 0 {
		return 0
	}
	shcore := syscall.NewLazyDLL("shcore.dll")
	get := shcore.NewProc("GetDpiForMonitor")
	if get.Find() != nil {
		return 0
	}
	var dx, dy uint32
	if r, _, _ := get.Call(h, 0, uintptr(unsafe.Pointer(&dx)), uintptr(unsafe.Pointer(&dy))); r != 0 {
		return 0
	}
	return int(dx)
}

func systemDPI(user32 *syscall.LazyDLL) int {
	getDC := user32.NewProc("GetDC")
	release := user32.NewProc("ReleaseDC")
	if getDC.Find() != nil || release.Find() != nil {
		return 96
	}
	dc, _, _ := getDC.Call(0)
	if dc == 0 {
		return 96
	}
	defer release.Call(0, dc)
	// GetDeviceCaps is exported by gdi32. Looking it up in user32 panics,
	// because LazyProc.Call calls mustFind.
	caps := syscall.NewLazyDLL("gdi32.dll").NewProc("GetDeviceCaps")
	if caps.Find() != nil {
		return 96
	}
	dpi, _, _ := caps.Call(dc, 88) // LOGPIXELSX
	if dpi == 0 {
		return 96
	}
	return int(dpi)
}

func utf16BytesToString(b []byte) string {
	if len(b) < 2 {
		return ""
	}
	u := make([]uint16, 0, len(b)/2)
	for i := 0; i+1 < len(b); i += 2 {
		c := binary.LittleEndian.Uint16(b[i:])
		if c == 0 {
			break
		}
		u = append(u, c)
	}
	return strings.TrimSpace(string(utf16.Decode(u)))
}
