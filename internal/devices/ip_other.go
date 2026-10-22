//go:build !windows

package dev

import "strings"

func formatLocalIP(name, cidr string) string {
	ip, _, _ := strings.Cut(cidr, "/")
	return ip
}

func platformRouteIface() string {
	return ""
}
