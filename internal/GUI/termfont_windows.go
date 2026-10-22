//go:build windows

package gui

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/ffgan/gf/internal/utils"
)

func windowsTermFont(term string) string {
	low := strings.ToLower(term)
	switch {
	case strings.Contains(low, "windows terminal"), strings.Contains(low, "windowsterminal"):
		return windowsTerminalFont()
	case strings.Contains(low, "mintty"):
		return minttyFont()
	case strings.Contains(low, "conhost"), low == "cmd":
		return conhostFont()
	case strings.Contains(low, "conemu"):
		return conemuFont()
	case strings.Contains(low, "warp"):
		return warpFont()
	default:
		return ""
	}
}

func windowsTerminalFont() string {
	id := strings.Trim(os.Getenv("WT_PROFILE_ID"), "'\" ")
	var chosen resolvedFont
	var have bool
	for _, path := range wtSettingsPaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		font, matched, ok := parseWTDocument(data, id)
		if !ok {
			continue
		}
		if matched || !have {
			chosen = font
			have = true
		}
		if matched {
			break
		}
	}
	if !have {
		return ""
	}
	if id != "" && (chosen.face == "" || !chosen.hasSize) {
		mergeWTFragments(id, &chosen)
	}
	return finishTermFont(chosen)
}

func wtSettingsPaths() []string {
	local := os.Getenv("LOCALAPPDATA")
	if local == "" {
		return nil
	}
	return []string{
		filepath.Join(local, `Packages`, `Microsoft.WindowsTerminal_8wekyb3d8bbwe`, `LocalState`, `settings.json`),
		filepath.Join(local, `Packages`, `Microsoft.WindowsTerminalPreview_8wekyb3d8bbwe`, `LocalState`, `settings.json`),
		filepath.Join(local, `Microsoft`, `Windows Terminal`, `settings.json`),
	}
}

func mergeWTFragments(id string, font *resolvedFont) {
	var dirs []string
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		dirs = append(dirs, filepath.Join(local, "Microsoft", "Windows Terminal", "Fragments"))
	}
	if pdata := os.Getenv("ProgramData"); pdata != "" {
		dirs = append(dirs, filepath.Join(pdata, "Microsoft", "Windows Terminal", "Fragments"))
	}
	for _, dir := range dirs {
		apps, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, app := range apps {
			if !app.IsDir() || strings.HasPrefix(app.Name(), ".") {
				continue
			}
			files, err := os.ReadDir(filepath.Join(dir, app.Name()))
			if err != nil {
				continue
			}
			for _, file := range files {
				if file.IsDir() || !strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
					continue
				}
				data, err := os.ReadFile(filepath.Join(dir, app.Name(), file.Name()))
				if err != nil {
					continue
				}
				applyWTFragment(data, id, font)
			}
		}
	}
}

func minttyFont() string {
	name, size := "", ""
	for _, path := range minttyConfigPaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		n, s := parseMinttyConfig(string(data))
		if name == "" {
			name = n
		}
		if size == "" {
			size = s
		}
	}
	if name == "" {
		name = "Lucida Console"
	}
	if size == "" {
		size = "9"
	}
	return formatTermFont(name, size, "pt")
}

func minttyConfigPaths() []string {
	var homes []string
	if h := os.Getenv("HOME"); h != "" {
		homes = append(homes, h)
	}
	if h := os.Getenv("USERPROFILE"); h != "" && !strings.EqualFold(h, os.Getenv("HOME")) {
		homes = append(homes, h)
	}
	var paths []string
	for _, home := range homes {
		paths = append(paths,
			filepath.Join(home, ".config", "mintty", "config"),
			filepath.Join(home, ".minttyrc"),
		)
	}
	return paths
}

func conemuFont() string {
	name, size := "", ""
	for _, env := range []string{"ConEmuDir", "ConEmuBaseDir", "APPDATA"} {
		dir := os.Getenv(env)
		if dir == "" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, "ConEmu.xml"))
		if err != nil {
			continue
		}
		name, size = parseConEmuFont(string(data))
		break
	}
	if name == "" {
		name = "Consolas"
	}
	if size == "" {
		size = "14"
	}
	return formatTermFont(name, size, "pt")
}

func conhostFont() string {
	kernel := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel.NewProc("GetCurrentConsoleFontEx")
	if proc.Find() != nil {
		return ""
	}
	buf := make([]byte, 84)
	binary.LittleEndian.PutUint32(buf[0:4], 84)
	handle, _, _ := kernel.NewProc("GetStdHandle").Call(^uintptr(10))
	r, _, _ := proc.Call(handle, 0, uintptr(unsafe.Pointer(&buf[0])))
	if r == 0 {
		return ""
	}
	height := binary.LittleEndian.Uint16(buf[10:12])
	var chars []uint16
	for i := 20; i+1 < len(buf) && len(chars) < 32; i += 2 {
		c := binary.LittleEndian.Uint16(buf[i:])
		if c == 0 {
			break
		}
		chars = append(chars, c)
	}
	name := strings.TrimSpace(string(utf16.Decode(chars)))
	if name == "" || height == 0 {
		return ""
	}
	return formatTermFont(name, strconv.Itoa(int(height)), "pt")
}

func warpFont() string {
	name, nameOK := utils.RegString(utils.HKCU, `Software\Warp.dev\Warp`, "FontName")
	size, sizeOK := utils.RegString(utils.HKCU, `Software\Warp.dev\Warp`, "FontSize")
	name = strings.Trim(name, `"`)
	size = strings.Trim(size, `"`)
	if !nameOK || name == "" {
		name = "Hack"
	}
	if !sizeOK || size == "" {
		size = "13.0"
	}
	size = strings.TrimSuffix(strings.ToLower(size), "px")
	return formatTermFont(name, size, "px")
}
