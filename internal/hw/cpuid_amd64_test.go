//go:build amd64

package hw

import (
	"encoding/binary"
	"testing"
)

func TestCPUIDLeaf0Vendor(t *testing.T) {
	_, b, c, d := cpuidRegs(0, 0)
	var buf [12]byte
	binary.LittleEndian.PutUint32(buf[0:4], b)
	binary.LittleEndian.PutUint32(buf[4:8], d)
	binary.LittleEndian.PutUint32(buf[8:12], c)
	got := string(buf[:])
	switch got {
	case "GenuineIntel", "AuthenticAMD", "HygonGenuine", "CentaurHauls", "Shanghai":
	default:
		t.Fatalf("CPUID vendor %q", got)
	}
}
