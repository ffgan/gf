//go:build amd64

package hw

import "encoding/binary"

func cpuidRegs(eax, ecx uint32) (a, b, c, d uint32)

// Hypervisor reads the hypervisor bit and the CPUID vendor string.
// An empty result means the CPU did not advertise a hypervisor.
func Hypervisor() string {
	_, _, c, _ := cpuidRegs(1, 0)
	if c&(1<<31) == 0 {
		return ""
	}
	_, b, c2, d := cpuidRegs(0x40000000, 0)
	return HypervisorName(decodeVendor(b, c2, d))
}

func decodeVendor(b, c, d uint32) string {
	var buf [12]byte
	binary.LittleEndian.PutUint32(buf[0:4], b)
	binary.LittleEndian.PutUint32(buf[4:8], c)
	binary.LittleEndian.PutUint32(buf[8:12], d)
	return string(buf[:])
}
