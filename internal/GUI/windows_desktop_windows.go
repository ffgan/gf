//go:build windows

package gui

import (
	"bytes"
	"encoding/binary"
	"os"
	"syscall"
	"unsafe"

	cli "github.com/ffgan/gf/internal/CLI"
	"github.com/ffgan/gf/internal/utils"
)

func windowsTheme() string {
	return cli.WindowsDesktop()
}

func windowsIcons() string {
	const key = `Software\Microsoft\Windows\CurrentVersion\Explorer\HideDesktopIcons\NewStartPanel`
	ids := []string{
		"{20D04FE0-3AEA-1069-A2D8-08002B30309D}",
		"{59031a47-3f72-44a7-89c5-5595fe6b30ee}",
		"{F02C1A0D-BE21-4350-88B0-7367FC96EF3C}",
		"{645FF040-5081-101B-9F08-00AA002F954E}",
		"{5399E694-6CE5-4D6C-8FCE-1D8870FDCBA0}",
	}
	vals := map[string]uint32{}
	for _, id := range ids {
		if v, ok := utils.RegUint(utils.HKCU, key, id); ok {
			vals[id] = v
		}
	}
	exists := utils.RegKeyExists(utils.HKCU, key) ||
		utils.RegKeyExists(utils.HKCU, `Software\Microsoft\Windows\CurrentVersion\Explorer\HideDesktopIcons\ClassicStartMenu`)
	return desktopIcons(exists, vals)
}

func windowsCursor() string {
	theme, ok := utils.RegString(utils.HKCU, `Control Panel\Cursors`, "")
	if !ok || theme == "" {
		return ""
	}
	size, has := utils.RegUint(utils.HKCU, `Control Panel\Cursors`, "CursorBaseSize")
	return formatWindowsCursor(theme, size, has)
}

func windowsUIFont() string {
	const key = `Control Panel\Desktop\WindowMetrics`
	names := [4]string{"CaptionFont", "MenuFont", "MessageFont", "StatusFont"}
	var fonts [4]string
	for i, name := range names {
		if b, ok := utils.RegBytes(utils.HKCU, key, name); ok {
			fonts[i] = logFontLabel(b)
		}
	}
	return formatWindowsFonts(fonts)
}

func windowsWMTheme() string {
	theme, _ := utils.RegString(utils.HKCU, `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes`, "CurrentTheme")
	color, hasColor := utils.RegUint(utils.HKCU, `Software\Microsoft\Windows\DWM`, "AccentColor")
	bgr := true
	if !hasColor {
		color, hasColor = utils.RegUint(utils.HKCU, `Software\Microsoft\Windows\DWM`, "ColorizationColor")
		bgr = false
	}
	const personalize = `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`
	system, sysOK := utils.RegUint(utils.HKCU, personalize, "SystemUsesLightTheme")
	apps, appsOK := utils.RegUint(utils.HKCU, personalize, "AppsUseLightTheme")
	if !sysOK {
		system = 1
	}
	if !appsOK {
		apps = 1
	}
	return formatWindowsWMTheme(theme, color, hasColor, bgr, system, apps, sysOK || appsOK)
}

func windowsWMVersion() string {
	root := os.Getenv("SystemRoot")
	if root == "" {
		root = `C:\Windows`
	}
	return fileVersion(root + `\System32\dwm.exe`)
}

func fileVersion(path string) string {
	version := syscall.NewLazyDLL("version.dll")
	sizeProc := version.NewProc("GetFileVersionInfoSizeW")
	infoProc := version.NewProc("GetFileVersionInfoW")
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return ""
	}
	size, _, _ := sizeProc.Call(uintptr(unsafe.Pointer(p)), 0)
	if size == 0 {
		return ""
	}
	buf := make([]byte, size)
	r, _, _ := infoProc.Call(uintptr(unsafe.Pointer(p)), 0, size, uintptr(unsafe.Pointer(&buf[0])))
	if r == 0 {
		return ""
	}
	// VS_FIXEDFILEINFO starts with signature 0xFEEF04BD.
	idx := bytes.Index(buf, []byte{0xBD, 0x04, 0xEF, 0xFE})
	if idx < 0 || idx+16 > len(buf) {
		return ""
	}
	ms := binary.LittleEndian.Uint32(buf[idx+8:])
	ls := binary.LittleEndian.Uint32(buf[idx+12:])
	return itoa(ms>>16) + "." + itoa(ms&0xFFFF) + "." + itoa(ls>>16) + "." + itoa(ls&0xFFFF)
}
