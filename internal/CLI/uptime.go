package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getUptime(uptimeShorthand string) string {
	var s int64
	osName := getOS()

	switch osName {
	case "Linux", "windows":
		// Try reading /proc/uptime
		if data, err := os.ReadFile("/proc/uptime"); err == nil {
			fields := strings.Fields(string(data))
			if len(fields) > 0 {
				if val, err := strconv.ParseFloat(fields[0], 64); err == nil {
					s = int64(val)
					break
				}
			}
		}

		// Fallback using uptime -s
		bootStr, err1 := exec.Command("uptime", "-s").Output()
		now := time.Now().Unix()
		if err1 == nil {
			if bootTime, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(string(bootStr))); err == nil {
				s = now - bootTime.Unix()
			}
		}

	case "darwin", "freebsd", "netbsd", "openbsd":
		out, err := exec.Command("sysctl", "-n", "kern.boottime").Output()
		if err == nil {
			line := strings.TrimSpace(string(out))
			line = strings.TrimPrefix(line, "{ sec = ")
			line = strings.Split(line, ",")[0]
			if boot, err := strconv.ParseInt(line, 10, 64); err == nil {
				s = time.Now().Unix() - boot
			}
		}

	case "solaris":
		out, err := exec.Command("kstat", "-p", "unix:0:system_misc:boot_time").Output()
		if err == nil {
			fields := strings.Fields(string(out))
			if len(fields) == 2 {
				if boot, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
					s = time.Now().Unix() - boot
				}
			}
		}

	default:
		// Try a generic fallback with "ps -o etime"
		var cmd *exec.Cmd
		if osName == "aix" || osName == "irix" || osName == "ironclad" {
			cmd = exec.Command("ps", "-o", "etime=", "-p", "1")
		} else {
			cmd = exec.Command("ps", "-o", "etime=", "-p", "0")
		}

		out, err := cmd.Output()
		if err == nil {
			t := strings.TrimSpace(string(out))
			t = strings.TrimSuffix(t, "\n")
			t = strings.TrimSuffix(t, "\r")

			var d, h, m, sec int64
			if strings.Contains(t, "-") {
				parts := strings.SplitN(t, "-", 2)
				d, _ = strconv.ParseInt(parts[0], 10, 64)
				t = parts[1]
			}
			timeParts := strings.Split(t, ":")
			switch len(timeParts) {
			case 3:
				h, _ = strconv.ParseInt(timeParts[0], 10, 64)
				m, _ = strconv.ParseInt(timeParts[1], 10, 64)
				sec, _ = strconv.ParseInt(timeParts[2], 10, 64)
			case 2:
				m, _ = strconv.ParseInt(timeParts[0], 10, 64)
				sec, _ = strconv.ParseInt(timeParts[1], 10, 64)
			}
			s = d*86400 + h*3600 + m*60 + sec
		}
	}

	// Convert uptime to days/hours/minutes
	d := s / 86400
	h := (s / 3600) % 24
	m := (s / 60) % 60

	// Format readable uptime
	var parts []string
	if d > 0 {
		part := fmt.Sprintf("%d days", d)
		if d == 1 {
			part = "1 day"
		}
		if d > 100 {
			part += "(!)"
		}
		parts = append(parts, part)
	}
	if h > 0 {
		if h == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", h))
		}
	}
	if m > 0 {
		if m == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", m))
		}
	}

	uptime := strings.Join(parts, ", ")
	if uptime == "" {
		uptime = fmt.Sprintf("%d seconds", s)
	}

	// Handle shorthand formats
	switch uptimeShorthand {
	case "on":
		uptime = strings.ReplaceAll(uptime, " minutes", " mins")
		uptime = strings.ReplaceAll(uptime, " minute", " min")
		uptime = strings.ReplaceAll(uptime, " seconds", " secs")
	case "tiny":
		r := strings.NewReplacer(
			" days", "d", " day", "d",
			" hours", "h", " hour", "h",
			" minutes", "m", " minute", "m",
			" seconds", "s", " second", "s",
			",", "",
		)
		uptime = r.Replace(uptime)
	}

	return strings.TrimSpace(uptime)
}
