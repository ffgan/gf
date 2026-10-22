package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ffgan/gf/internal/utils"
)

func GetUptime(osName, uptimeShorthand string) string {
	var s int64

	// Any OS that publishes procfs can answer without a process. Linux and
	// Windows keep the uptime -s fallback when that file is missing.
	if sec, ok := uptimeFromProc(); ok {
		s = sec
	} else if sec, ok := uptimeFromPlatform(); ok {
		s = sec
	} else {
		s = uptimeFromCommand(osName)
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
	case utils.ON:
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

func uptimeFromCommand(osName string) int64 {
	switch {
	case osName == Linux || osName == Windows:
		bootStr, err1 := utils.CommandOutput("uptime", "-s")
		if err1 != nil {
			return 0
		}
		bootTime, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(string(bootStr)))
		if err != nil {
			return 0
		}
		return time.Now().Unix() - bootTime.Unix()
	case osName == Darwin || isBSDFamily(osName):
		out, err := utils.CommandOutput("sysctl", "-n", "kern.boottime")
		if err != nil {
			return 0
		}
		line := strings.TrimSpace(string(out))
		line = strings.TrimPrefix(line, "{ sec = ")
		line = strings.Split(line, ",")[0]
		boot, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return 0
		}
		return time.Now().Unix() - boot
	case osName == Solaris || osName == illumos:
		out, err := utils.CommandOutput("kstat", "-p", "unix:0:system_misc:boot_time")
		if err != nil {
			return 0
		}
		fields := strings.Fields(string(out))
		if len(fields) != 2 {
			return 0
		}
		boot, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return 0
		}
		return time.Now().Unix() - boot
	default:
		return uptimeFromPS(osName)
	}
}

func uptimeFromPS(osName string) int64 {
	var out []byte
	var err error
	if osName == AIX || osName == IRIX || osName == Ironclad {
		out, err = utils.CommandOutput("ps", "-o", "etime=", "-p", "1")
	} else {
		out, err = utils.CommandOutput("ps", "-o", "etime=", "-p", "0")
	}
	if err != nil {
		return 0
	}
	t := strings.TrimSpace(string(out))
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
	return d*86400 + h*3600 + m*60 + sec
}

func uptimeFromProc() (int64, bool) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, false
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, false
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, false
	}
	return int64(val), true
}
