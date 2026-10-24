package dev

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// LocalIP returns the IPv4 address of the requested interface.
// "auto" follows the default route, then the first global address.
func LocalIP(ifaceName string) string {
	// Windows net.Interface* calls GetAdaptersAddresses with gateway
	// data. That walk is several seconds on machines with virtual NICs,
	// and LocalIP used to do it twice.
	if text, ok := platformLocalIP(ifaceName); ok {
		return text
	}
	ifaceName = strings.Trim(ifaceName, "() '\"")
	name := ifaceName
	if name == "" || name == "auto" {
		name = defaultRouteIface()
		if name == "" {
			name = platformRouteIface()
		}
	}
	cidr := ""
	if name != "" && name != "auto" {
		cidr = ifaceCIDR(name)
	}
	if cidr == "" {
		fallback, ip := firstGlobalCIDR()
		cidr = ip
		if name == "" || name == "auto" {
			name = fallback
		}
	}
	return formatLocalIP(name, cidr)
}

// PublicIP fetches the address from host, which defaults to ident.me.
func PublicIP(host, timeoutSec string) string {
	host = strings.Trim(host, `"'`)
	if host == "" {
		host = "https://ident.me"
	}
	seconds, _ := strconv.Atoi(strings.TrimSpace(timeoutSec))
	if seconds <= 0 {
		seconds = 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host, nil)
	if err != nil {
		return ""
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 256))
	if err != nil {
		return ""
	}
	text := strings.TrimSpace(string(body))
	if strings.ContainsAny(text, " \n<") {
		return ""
	}
	return text
}

func defaultRouteIface() string {
	data, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return ""
	}
	bestMetric := int(^uint(0) >> 1)
	best := ""
	lines := strings.Split(string(data), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 7 || fields[1] != "00000000" {
			continue
		}
		metric, _ := strconv.ParseInt(fields[6], 16, 64)
		if best == "" || int(metric) < bestMetric {
			best = fields[0]
			bestMetric = int(metric)
		}
	}
	return best
}

func ifaceCIDR(name string) string {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if cidr := addrCIDR(addr); cidr != "" {
			return cidr
		}
	}
	return ""
}

func firstGlobalCIDR() (string, string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if cidr := addrCIDR(addr); cidr != "" {
				return iface.Name, cidr
			}
		}
	}
	return "", ""
}

// ipv4Prefix formats a MIB_IPADDRROW address and mask. Windows stores
// both DWORDs in network byte order.
func ipv4Prefix(addr, mask uint32) string {
	ip := net.IPv4(byte(addr), byte(addr>>8), byte(addr>>16), byte(addr>>24)).To4()
	if ip == nil || !ip.IsGlobalUnicast() {
		return ""
	}
	m := net.IPv4Mask(byte(mask), byte(mask>>8), byte(mask>>16), byte(mask>>24))
	ones, bits := m.Size()
	if bits != 32 || ones <= 0 {
		return ip.String()
	}
	return fmt.Sprintf("%s/%d", ip, ones)
}

func addrCIDR(addr net.Addr) string {
	ipnet, ok := addr.(*net.IPNet)
	if !ok || ipnet.IP == nil {
		return ""
	}
	ip4 := ipnet.IP.To4()
	if ip4 == nil || !ip4.IsGlobalUnicast() {
		return ""
	}
	ones, bits := ipnet.Mask.Size()
	if bits != 32 || ones <= 0 {
		return ip4.String()
	}
	return fmt.Sprintf("%s/%d", ip4, ones)
}
