//go:build darwin

package dev

import (
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func formatLocalIP(name, cidr string) string {
	ip, _, _ := strings.Cut(cidr, "/")
	return ip
}

func platformRouteIface() string {
	return parseRouteInterface(utils.RunCommand("route", "-n", "get", "default"))
}

// platformLocalIP uses the default route, then ipconfig. The first
// interface from getifaddrs is often awdl, bridge, or a link-local
// address, which is not the machine's LAN address.
func platformLocalIP(ifaceName string) (string, bool) {
	requested := strings.Trim(ifaceName, "() '\"")
	auto := requested == "" || strings.EqualFold(requested, "auto")
	name := requested
	if auto {
		name = parseRouteInterface(utils.RunCommand("route", "-n", "get", "default"))
		if name == "" {
			return "", false
		}
	}
	if ip := usableIPv4Text(utils.RunCommand("ipconfig", "getifaddr", name)); ip != "" {
		return ip, true
	}
	// ipconfig refuses utun interfaces. ifconfig still has the address.
	if ip := parseIfconfigIPv4(utils.RunCommand("ifconfig", name)); ip != "" {
		return ip, true
	}
	if auto {
		return "", false
	}
	return "", true
}
