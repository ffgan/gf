package cli

import "testing"

func TestParseDFPOSIX(t *testing.T) {
	tests := []struct {
		line    string
		used    int64
		total   int64
		avail   int64
		percent int
	}{
		{
			line:    "/dev/vdc         265592036 29023376 236552276      11% /",
			used:    29023376,
			total:   265592036,
			avail:   236552276,
			percent: 11,
		},
		{
			line:    "My Disk 1000 200 800 20% /Volumes/My Disk",
			used:    200,
			total:   1000,
			avail:   800,
			percent: 20,
		},
		{
			// macOS df without -P adds an inode percentage. Capacity is the first %.
			line:    "/dev/disk3s1s1 489620184 11694280 228926040 5% 404862 1144630200 0% /",
			used:    11694280,
			total:   489620184,
			avail:   228926040,
			percent: 5,
		},
	}
	for _, tt := range tests {
		used, total, avail, percent, err := parseDFPOSIX(tt.line)
		if err != nil {
			t.Fatalf("%q: %v", tt.line, err)
		}
		if used != tt.used || total != tt.total || avail != tt.avail || percent != tt.percent {
			t.Fatalf("%q: got used=%d total=%d avail=%d percent=%d", tt.line, used, total, avail, percent)
		}
	}
	if _, _, _, _, err := parseDFPOSIX("not df output"); err == nil {
		t.Fatal("expected an error")
	}
}

func TestContainerUsage(t *testing.T) {
	const gib = int64(1024 * 1024 * 1024)
	// Sealed system volume: 10 GiB used, 55 GiB free, 466 GiB container.
	// df's own capacity is about 16%, which does not match either figure.
	used, pct := containerUsage("/", 10*gib, 466*gib, 55*gib, true)
	if used != 411*gib || pct != 88 {
		t.Fatalf("system volume: used=%d pct=%d", used, pct)
	}
	// The data volume's Used column is already that volume. Leave it alone.
	used, pct = containerUsage("/System/Volumes/Data", 398*gib, 466*gib, 55*gib, true)
	if used != 398*gib || pct != 85 {
		t.Fatalf("data volume: used=%d pct=%d", used, pct)
	}
	// No volume group: keep df's used bytes.
	used, pct = containerUsage("/", 10*gib, 466*gib, 55*gib, false)
	if used != 10*gib || pct != 2 {
		t.Fatalf("no data volume: used=%d pct=%d", used, pct)
	}
	used, pct = containerUsage("/", 20*gib, 100*gib, 80*gib, true)
	if used != 20*gib || pct != 20 {
		t.Fatalf("full accounting: used=%d pct=%d", used, pct)
	}
}
