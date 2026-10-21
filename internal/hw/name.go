package hw

import "strings"

// HypervisorName turns a 12-byte CPUID 0x40000000 vendor into a short label.
func HypervisorName(raw string) string {
	s := strings.TrimRight(raw, "\x00 ")
	switch {
	case s == "KVMKVMKVM", strings.HasPrefix(s, "KVMKVM"):
		return "KVM"
	case s == "Microsoft Hv":
		return "Microsoft Hyper-V"
	case strings.HasPrefix(s, "VMware"):
		return "VMware"
	case strings.HasPrefix(s, "XenVMM"):
		return "Xen"
	case strings.HasPrefix(s, "TCGTCG"):
		return "QEMU TCG"
	case strings.HasPrefix(s, "bhyve"):
		return "bhyve"
	case strings.HasPrefix(s, "ACRN"):
		return "ACRN"
	case strings.HasPrefix(s, "prl hyperv"):
		return "Parallels"
	case s == "":
		return "Virtual Machine"
	default:
		return s
	}
}
