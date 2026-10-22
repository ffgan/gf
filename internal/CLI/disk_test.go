package cli

import "testing"

func TestParseDFPOSIX(t *testing.T) {
	tests := []struct {
		line    string
		used    int64
		total   int64
		percent int
	}{
		{
			line:    "/dev/vdc         265592036 29023376 236552276      11% /",
			used:    29023376,
			total:   265592036,
			percent: 11,
		},
		{
			line:    "My Disk 1000 200 800 20% /Volumes/My Disk",
			used:    200,
			total:   1000,
			percent: 20,
		},
	}
	for _, tt := range tests {
		used, total, percent, err := parseDFPOSIX(tt.line)
		if err != nil {
			t.Fatalf("%q: %v", tt.line, err)
		}
		if used != tt.used || total != tt.total || percent != tt.percent {
			t.Fatalf("%q: got used=%d total=%d percent=%d", tt.line, used, total, percent)
		}
	}
	if _, _, _, err := parseDFPOSIX("not df output"); err == nil {
		t.Fatal("expected an error")
	}
}
