package dev

import (
	"context"
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
	ifaceName = strings.Trim(ifaceName, "() '\"")
	if ifaceName == "" || ifaceName == "auto" {
		ifaceName = defaultRouteIface()
	}
	if ifaceName != "" && ifaceName != "auto" {
		if ip := ifaceIPv4(ifaceName); ip != "" {
			return ip
		}
	}
	return firstGlobalIPv4()
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

func ifaceIPv4(name string) string {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ip := addrIP(addr)
		if ip != "" {
			return ip
		}
	}
	return ""
}

func firstGlobalIPv4() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
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
			if ip := addrIP(addr); ip != "" {
				return ip
			}
		}
	}
	return ""
}

func addrIP(addr net.Addr) string {
	ipnet, ok := addr.(*net.IPNet)
	if !ok || ipnet.IP == nil || ipnet.IP.IsLoopback() {
		return ""
	}
	ip4 := ipnet.IP.To4()
	if ip4 == nil {
		return ""
	}
	return ip4.String()
}
