package hw

import "testing"

func TestHypervisorName(t *testing.T) {
	cases := map[string]string{
		"KVMKVMKVM":     "KVM",
		"Microsoft Hv":  "Microsoft Hyper-V",
		"VMwareVMware":  "VMware",
		"XenVMMXenVMM":  "Xen",
		"TCGTCGTCGTCG":  "QEMU TCG",
		"bhyve bhyve ":  "bhyve",
		"ACRNACRNACRN":  "ACRN",
		"prl hyperv  ":  "Parallels",
		"":              "Virtual Machine",
		"SomethingElse": "SomethingElse",
	}
	for in, want := range cases {
		if got := HypervisorName(in); got != want {
			t.Errorf("%q: got %q want %q", in, got, want)
		}
	}
}
