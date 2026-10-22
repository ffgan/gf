package dev

import "testing"

func TestBatteryFromStatus(t *testing.T) {
	if got := batteryFromStatus(1, 128, 255, 0xFFFFFFFF); got != "" {
		t.Fatalf("desktop = %q", got)
	}
	if got := batteryFromStatus(255, 255, 255, 0xFFFFFFFF); got != "" {
		t.Fatalf("unknown = %q", got)
	}
	got := batteryFromStatus(0, 2, 45, 3720)
	want := "45% (1 hour, 2 mins remaining) [Discharging]"
	if got != want {
		t.Fatalf("discharging = %q, want %q", got, want)
	}
	got = batteryFromStatus(1, 8|4, 12, 0xFFFFFFFF)
	want = "12% [AC Connected, Charging, Critical]"
	if got != want {
		t.Fatalf("charging = %q, want %q", got, want)
	}
	got = batteryFromStatus(1, 1, 100, 0)
	want = "100% [AC Connected]"
	if got != want {
		t.Fatalf("full = %q, want %q", got, want)
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		seconds uint64
		want    string
	}{
		{1, "1 second"},
		{42, "42 seconds"},
		{59, "59 seconds"},
		{60, "1 min"},
		{89, "1 min"},
		{90, "2 mins"},
		{3600, "1 hour"},
		{3630, "1 hour, 1 min"},
		{3660, "1 hour, 1 min"},
		{3720, "1 hour, 2 mins"},
		{90061, "1 day, 1 hour, 1 min"},
	}
	for _, tc := range cases {
		if got := formatDuration(tc.seconds); got != tc.want {
			t.Fatalf("%d seconds = %q, want %q", tc.seconds, got, tc.want)
		}
	}
}
