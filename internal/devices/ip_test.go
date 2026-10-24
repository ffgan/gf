package dev

import (
	"net"
	"testing"
)

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
	linkLocal := uint32(169) | uint32(254)<<8 | uint32(1)<<16 | uint32(8)<<24
	if got := ipv4Prefix(linkLocal, 0xffff0000); got != "" {
		t.Fatalf("link-local %q", got)
	}
}

func TestParseRouteInterface(t *testing.T) {
	text := `
   route to: default
destination: default
    gateway: 192.168.1.1
  interface: en0
      flags: <UP,GATEWAY,DONE,STATIC,PRCLONING,GLOBAL>
`
	if got := parseRouteInterface(text); got != "en0" {
		t.Fatalf("interface %q", got)
	}
	if got := parseRouteInterface("interface:\n"); got != "" {
		t.Fatalf("empty interface %q", got)
	}
}

func TestParseIfconfigIPv4(t *testing.T) {
	text := `
en0: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
	inet 169.254.22.8 netmask 0xffff0000 broadcast 169.254.255.255
	inet 10.1.2.3 netmask 0xffffff00 broadcast 10.1.2.255
`
	if got := parseIfconfigIPv4(text); got != "10.1.2.3" {
		t.Fatalf("got %q", got)
	}
	if got := usableIPv4Text("169.254.1.1"); got != "" {
		t.Fatalf("link-local %q", got)
	}
	if got := usableIPv4Text("10.0.0.8"); got != "10.0.0.8" {
		t.Fatalf("got %q", got)
	}
}

func TestAddrCIDRSkipsLinkLocal(t *testing.T) {
	link := &net.IPNet{IP: net.ParseIP("169.254.4.5"), Mask: net.CIDRMask(16, 32)}
	if got := addrCIDR(link); got != "" {
		t.Fatalf("link-local %q", got)
	}
	lan := &net.IPNet{IP: net.ParseIP("192.168.0.20"), Mask: net.CIDRMask(24, 32)}
	if got := addrCIDR(lan); got != "192.168.0.20/24" {
		t.Fatalf("got %q", got)
	}
}
