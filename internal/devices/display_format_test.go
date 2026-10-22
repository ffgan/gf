package dev

import "testing"

func TestFormatWinDisplay(t *testing.T) {
	got := formatWinDisplay(winDisplay{
		Name: "DELL P2425D", W: 2560, H: 1440,
		Hz: 60, DPI: 120, Wmm: 531, Hmm: 299,
		Kind: "External", ShowHz: true,
	})
	want := `DELL P2425D: 2560x1440 @ 1.25x in 24", 60 Hz [External]`
	if got != want {
		t.Fatalf("got %q", got)
	}
	primary := formatWinDisplay(winDisplay{
		Name: "DELL P2425D", W: 2560, H: 1440,
		Hz: 60, DPI: 120, Wmm: 531, Hmm: 299,
		Kind: "External", Primary: true, ShowHz: true,
	})
	if primary != want+" *" {
		t.Fatalf("primary %q", primary)
	}
}

func TestEdidKey(t *testing.T) {
	got := edidKeyFromDevicePath(`\\?\DISPLAY#DELA0B8#5&1&0&UID256#{e6f07b5f-ee97-4a90-b076-33f57bf4eaa7}`)
	want := `SYSTEM\CurrentControlSet\Enum\DISPLAY\DELA0B8\5&1&0&UID256\Device Parameters`
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestEdidName(t *testing.T) {
	edid := make([]byte, 128)
	edid[0] = 0x00
	edid[1] = 0xFF
	edid[21] = 53
	edid[22] = 30
	edid[54] = 0
	edid[55] = 0
	edid[56] = 0
	edid[57] = 0xFC
	copy(edid[59:], []byte("DELL P2425D\n "))
	name, wmm, hmm := edidNameSize(edid)
	if name != "DELL P2425D" || wmm != 530 || hmm != 300 {
		t.Fatalf("name %q size %d %d", name, wmm, hmm)
	}
}
