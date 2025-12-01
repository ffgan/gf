package dev

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	cli "github.com/ffgan/gf/internal/CLI"
	"github.com/ffgan/gf/internal/utils"
)

func GetNetwork(osName string) string {
	if strings.Contains(osName, cli.Linux) {
		return getNetworkLinux()
	}
	if strings.Contains(osName, cli.MacOSX) || strings.Contains(osName, cli.MacOS) {
		return getNetworkMac()
	}
	return "Unknown OS"
}

func getNetworkLinux() string {
	var networks []string

	entries, _ := filepath.Glob("/sys/class/net/*")
	for _, ifacePath := range entries {
		operstate, err := os.ReadFile(filepath.Join(ifacePath, "operstate"))
		if err != nil || strings.TrimSpace(string(operstate)) != "up" {
			continue
		}
		if _, err := os.Stat(filepath.Join(ifacePath, "device")); err != nil {
			continue
		}

		if _, err := os.Stat(filepath.Join(ifacePath, "phy80211")); err == nil {
			// Wi-Fi device
			networks = append(networks, "Wifi")
			phy, err := os.ReadFile(filepath.Join(ifacePath, "phy80211", "name"))
			if err == nil {
				out, _ := exec.Command("iw", strings.TrimSpace(string(phy)), "info").Output()
				if bytes.Contains(out, []byte("VHT Capabilities")) {
					networks = append(networks, "6")
				}
			}
		} else {
			speed, err := os.ReadFile(filepath.Join(ifacePath, "speed"))
			if err == nil {
				networks = append(networks, strings.TrimSpace(string(speed)))
			}
		}
	}

	return formatNetworks(networks)
}

func getNetworkMac() string {
	var network string

	// 获取当前活动网络接口
	activeIface := utils.ExecOutput(`route`, `get`, `default`)
	reIface := regexp.MustCompile(`interface: (\w+)`)
	matches := reIface.FindStringSubmatch(activeIface)
	if len(matches) < 2 {
		return "Unknown"
	}
	activeNetwork := matches[1]

	// 解析接口名称
	hwPorts := utils.ExecOutput(`networksetup`, `-listallhardwareports`)
	rePort := regexp.MustCompile(`Hardware Port: (.+)\nDevice: ` + activeNetwork)
	portMatch := rePort.FindStringSubmatch(hwPorts)
	var activeName string
	if len(portMatch) > 1 {
		activeName = portMatch[1]
	}

	var phyMode, linkSpeed string

	if activeName == "Wi-Fi" {
		tmp := "/tmp/neofetch_system_profiler_SPAirPortDataType.xml"
		_ = exec.Command("system_profiler", "-detailLevel", "basic", "-xml", "SPAirPortDataType").Run()

		phyMode = plistBuddyPrint(tmp, "0:_items:0:spairport_airport_interfaces:0:spairport_current_network_information:spairport_network_phymode")
		linkSpeed = plistBuddyPrint(tmp, "0:_items:0:spairport_airport_interfaces:0:spairport_current_network_information:spairport_network_rate")
		if linkSpeed != "" {
			linkSpeed += " Mbps"
		}
	} else {
		out := utils.ExecOutput(`ifconfig`, activeNetwork)
		reMedia := regexp.MustCompile(`media:.*\(([^)]+)\)`)
		match := reMedia.FindStringSubmatch(out)
		if len(match) > 1 {
			linkSpeed = match[1]
		}
	}

	network = fmt.Sprintf("%s: %s", activeNetwork, activeName)
	if phyMode != "" {
		network += fmt.Sprintf(" (%s)", phyMode)
	}
	if linkSpeed != "" {
		network += fmt.Sprintf(" @ %s", linkSpeed)
	}

	return network
}

func plistBuddyPrint(file, key string) string {
	out, err := exec.Command("PlistBuddy", "-c", "Print "+key, file).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func formatNetworks(list []string) string {
	counts := make(map[string]int)
	for _, n := range list {
		counts[n]++
	}

	type pair struct {
		name  string
		count int
	}
	var pairs []pair
	for k, v := range counts {
		pairs = append(pairs, pair{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].name > pairs[j].name
	})

	var result strings.Builder
	for _, p := range pairs {
		if p.count > 1 {
			result.WriteString(fmt.Sprintf("%dx ", p.count))
		}

		n := p.name
		switch {
		case n == "Wifi":
			result.WriteString("Wifi; ")
		case n == "6":
			result.WriteString("Wifi6; ")
		case n == "-1" || n == "":
			result.WriteString("Unknown; ")
		default:
			val, _ := strconv.Atoi(n)
			if n == "2500" {
				result.WriteString("2.5 Gbps; ")
			} else if val%1000 == 0 {
				result.WriteString(fmt.Sprintf("%s Gbps; ", strings.TrimSuffix(n, "000")))
			} else {
				result.WriteString(fmt.Sprintf("%d Mbps; ", val))
			}
		}
	}

	res := strings.TrimSuffix(result.String(), "; ")
	return res
}
