package dev

import "testing"

func TestIPv4Prefix(t *testing.T) {
	// MIB_IPADDRROW stores 172.19.66.187 and 255.255.252.0 in network order.
	addr := uint32(172) | uint32(19)<<8 | uint32(66)<<16 | uint32(187)<<24
	mask := uint32(255) | uint32(255)<<8 | uint32(252)<<16
	if got := ipv4Prefix(addr, mask); got != "172.19.66.187/22" {
		t.Fatalf("got %q", got)
	}
	loopback := uint32(127) | uint32(1)<<24
	if got := ipv4Prefix(loopback, 0xffffffff); got != "" {
		t.Fatalf("loopback %q", got)
	}
	if got := ipv4Prefix(0, 0); got != "" {
		t.Fatalf("unspecified %q", got)
	}
}
