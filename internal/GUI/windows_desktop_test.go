package gui

import (
	"encoding/binary"
	"testing"
	"unicode/utf16"
)

func TestWindowsWMTheme(t *testing.T) {
	got := formatWindowsWMTheme(`C:\Windows\resources\Themes\custom.theme`, 0x00D47800, true, true, 1, 1, true)
	want := "Custom - Blue (System: Light, Apps: Light)"
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestDesktopIcons(t *testing.T) {
	if desktopIcons(false, nil) != "Recycle Bin" {
		t.Fatal(desktopIcons(false, nil))
	}
	vals := map[string]uint32{
		"{20D04FE0-3AEA-1069-A2D8-08002B30309D}": 0,
		"{645FF040-5081-101B-9F08-00AA002F954E}": 0,
		"{5399E694-6CE5-4D6C-8FCE-1D8870FDCBA0}": 0,
	}
	got := desktopIcons(true, vals)
	if got != "This PC, Recycle Bin, Control Panel" {
		t.Fatalf("got %q", got)
	}
}

func TestWindowsCursor(t *testing.T) {
	if got := formatWindowsCursor("Windows Aero", 32, true); got != "Windows Aero (32px)" {
		t.Fatalf("got %q", got)
	}
}

func TestWindowsFont(t *testing.T) {
	raw := make([]byte, 92)
	var height int32 = -15
	binary.LittleEndian.PutUint32(raw[:4], uint32(height))
	name := utf16.Encode([]rune("Microsoft YaHei UI"))
	for i, c := range name {
		binary.LittleEndian.PutUint16(raw[28+i*2:], c)
	}
	label := logFontLabel(raw)
	var fonts [4]string
	for i := range fonts {
		fonts[i] = label
	}
	got := formatWindowsFonts(fonts)
	want := "Microsoft YaHei UI (15pt) [Caption / Menu / Message / Status]"
	if got != want {
		t.Fatalf("got %q", got)
	}
}
