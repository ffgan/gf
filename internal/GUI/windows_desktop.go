package gui

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func formatWindowsWMTheme(themePath string, color uint32, hasColor, colorIsBGR bool, system, apps uint32, hasMode bool) string {
	var b strings.Builder
	if base := themeFileBase(themePath); base != "" {
		b.WriteString(base)
	}
	if hasColor {
		if b.Len() > 0 {
			b.WriteString(" - ")
		}
		rgb := color & 0xFFFFFF
		if colorIsBGR {
			rgb = ((color & 0xFF) << 16) | (color & 0xFF00) | ((color >> 16) & 0xFF)
		}
		if name, ok := accentNames[rgb]; ok {
			b.WriteString(name)
		} else {
			b.WriteString("#")
			const hex = "0123456789ABCDEF"
			for shift := 20; shift >= 0; shift -= 4 {
				b.WriteByte(hex[(rgb>>uint(shift))&0xF])
			}
		}
	}
	if hasMode {
		paren := b.Len() > 0
		if paren {
			b.WriteString(" (")
		}
		light := func(v uint32) string {
			if v != 0 {
				return "Light"
			}
			return "Dark"
		}
		b.WriteString("System: ")
		b.WriteString(light(system))
		b.WriteString(", Apps: ")
		b.WriteString(light(apps))
		if paren {
			b.WriteByte(')')
		}
	}
	return b.String()
}

func themeFileBase(path string) string {
	path = strings.ReplaceAll(path, `\`, "/")
	if i := strings.LastIndex(path, "/"); i >= 0 {
		path = path[i+1:]
	}
	if i := strings.LastIndex(path, "."); i > 0 {
		path = path[:i]
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(path)
	if r == utf8.RuneError {
		return path
	}
	return string(unicode.ToUpper(r)) + path[size:]
}

var accentNames = map[uint32]string{
	0x696cc3: "Yellow gold",
	0xff8c00: "Gold",
	0xf7630c: "Orange bright",
	0xca5010: "Orange dark",
	0xda3b01: "Rust",
	0xef6950: "Pale rust",
	0xd13438: "Brick red",
	0xff4343: "Mod red",
	0xe74856: "Pale red",
	0xe81123: "Red",
	0xea005e: "Rose bright",
	0xc30052: "Rose",
	0xe3008c: "Plum light",
	0xbf0077: "Plum",
	0xc239b3: "Orchid light",
	0x9a0089: "Orchid",
	0x0078d4: "Blue",
	0x0063b1: "Navy blue",
	0x8d8bd7: "Purple shadow",
	0x6b69d6: "Purple shadow dark",
	0x8764b8: "Iris pastel",
	0x744da9: "Iris Spring",
	0xb146c2: "Violet red light",
	0x881798: "Violet red",
	0x0099bc: "Cool blue bright",
	0x2d7d9a: "Cool blue",
	0x00b7c3: "Seafoam",
	0x038387: "Seafoam teal",
	0x00b294: "Mint light",
	0x018574: "Mint dark",
	0x00cc6a: "Turf green",
	0x10893e: "Sport green",
	0x7a7574: "Gray",
	0x5d5a58: "Gray brown",
	0x68768a: "Steel blue",
	0x515c6b: "Metal blue",
	0x567c73: "Pale moss",
	0x486860: "Moss",
	0x498205: "Meadow green",
	0x107c10: "Green",
	0x767676: "Overcast",
	0x4c4a48: "Storm",
	0x69797e: "Blue gray",
	0x4a5459: "Gray dark",
	0x647c64: "Liddy green",
	0x4c574e: "Sage",
	0x807143: "Camouflage desert",
	0x766c59: "Camouflage",
	0x000000: "Black",
	0xFFFFFF: "White",
}

func desktopIcons(keyExists bool, vals map[string]uint32) string {
	if !keyExists {
		return "Recycle Bin"
	}
	get := func(id string, def uint32) uint32 {
		for k, v := range vals {
			if strings.EqualFold(k, id) {
				return v
			}
		}
		return def
	}
	thisPC := get("{20D04FE0-3AEA-1069-A2D8-08002B30309D}", 1)
	users := get("{59031a47-3f72-44a7-89c5-5595fe6b30ee}", 1)
	network := get("{F02C1A0D-BE21-4350-88B0-7367FC96EF3C}", 1)
	recycle := get("{645FF040-5081-101B-9F08-00AA002F954E}", 0)
	panel := get("{5399E694-6CE5-4D6C-8FCE-1D8870FDCBA0}", 1)
	if thisPC != 0 && users != 0 && network != 0 && recycle != 0 && panel != 0 {
		return ""
	}
	var parts []string
	if thisPC == 0 {
		parts = append(parts, "This PC")
	}
	if users == 0 {
		parts = append(parts, "User's Files")
	}
	if network == 0 {
		parts = append(parts, "Remote Network")
	}
	if recycle == 0 {
		parts = append(parts, "Recycle Bin")
	}
	if panel == 0 {
		parts = append(parts, "Control Panel")
	}
	return strings.Join(parts, ", ")
}

func formatWindowsCursor(theme string, size uint32, hasSize bool) string {
	theme = strings.TrimSpace(theme)
	if theme == "" {
		return ""
	}
	if hasSize && size > 0 {
		return theme + " (" + itoa(size) + "px)"
	}
	return theme
}

func itoa(n uint32) string {
	if n == 0 {
		return "0"
	}
	var buf [12]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
