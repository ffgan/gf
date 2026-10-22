package gui

import (
	"encoding/binary"
	"strings"
	"unicode/utf16"
)

func parseLogFont(b []byte) (name string, pt int, ok bool) {
	if len(b) < 32 {
		return "", 0, false
	}
	height := int32(binary.LittleEndian.Uint32(b[:4]))
	if height < 0 {
		pt = int(-height)
	}
	var chars []uint16
	for i := 28; i+1 < len(b) && len(chars) < 32; i += 2 {
		c := binary.LittleEndian.Uint16(b[i:])
		if c == 0 {
			break
		}
		chars = append(chars, c)
	}
	name = strings.TrimSpace(string(utf16.Decode(chars)))
	return name, pt, name != ""
}

func logFontLabel(b []byte) string {
	name, pt, ok := parseLogFont(b)
	if !ok {
		return ""
	}
	if pt > 0 {
		return name + " (" + itoa(uint32(pt)) + "pt)"
	}
	return name
}

func formatWindowsFonts(fonts [4]string) string {
	types := [4]string{"Caption", "Menu", "Message", "Status"}
	var b strings.Builder
	for i := 0; i < 4; i++ {
		if fonts[i] == "" {
			continue
		}
		if b.Len() == 0 || fonts[i] != fonts[i-1] {
			if b.Len() > 0 {
				b.WriteString("], ")
			}
			b.WriteString(fonts[i])
			b.WriteString(" [")
			b.WriteString(types[i])
		} else {
			b.WriteString(" / ")
			b.WriteString(types[i])
		}
	}
	if b.Len() == 0 {
		return ""
	}
	b.WriteByte(']')
	return b.String()
}
