package dev

import (
	"net"
	"strings"
)

// parseRouteInterface reads `route -n get default`. -n matters: without
// it, route may block on a reverse lookup of the gateway.
func parseRouteInterface(text string) string {
	for _, line := range strings.Split(text, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		if strings.TrimSpace(key) != "interface" {
			continue
		}
		fields := strings.Fields(value)
		if len(fields) == 0 {
			return ""
		}
		return fields[0]
	}
	return ""
}

// parseIfconfigIPv4 returns the first global IPv4 on an interface.
// Link-local 169.254 addresses are skipped.
func parseIfconfigIPv4(text string) string {
	for _, line := range strings.Split(text, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[0] != "inet" {
			continue
		}
		if ip := usableIPv4Text(strings.TrimPrefix(fields[1], "addr:")); ip != "" {
			return ip
		}
	}
	return ""
}

func usableIPv4Text(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	ipText, _, _ := strings.Cut(text, "/")
	parsed := net.ParseIP(ipText)
	if parsed == nil {
		return ""
	}
	ip4 := parsed.To4()
	if ip4 == nil || !ip4.IsGlobalUnicast() {
		return ""
	}
	return ip4.String()
}
