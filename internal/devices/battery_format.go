package dev

import (
	"strconv"
	"strings"
)

const (
	batteryFlagCritical  = 4
	batteryFlagCharging  = 8
	batteryFlagNoBattery = 128
	batteryFlagUnknown   = 255
)

// batteryFromStatus formats a Windows SYSTEM_POWER_STATUS the way fastfetch
// prints a battery. No battery, or an unknown status with no percentage,
// returns an empty string so the line stays hidden.
func batteryFromStatus(ac, flag, percent byte, lifeSeconds uint32) string {
	if flag&batteryFlagNoBattery != 0 {
		return ""
	}
	if percent == batteryFlagUnknown && (flag == batteryFlagUnknown || flag == 0) {
		return ""
	}
	capacity := -1
	if percent <= 100 {
		capacity = int(percent)
	}
	var status []string
	if ac == 1 {
		status = append(status, "AC Connected")
	}
	if flag&batteryFlagCharging != 0 {
		status = append(status, "Charging")
	}
	if ac == 0 && flag&batteryFlagCharging == 0 && capacity >= 0 {
		status = append(status, "Discharging")
	}
	if flag&batteryFlagCritical != 0 {
		status = append(status, "Critical")
	}
	seconds := int64(-1)
	if lifeSeconds != 0xFFFFFFFF && lifeSeconds > 0 {
		seconds = int64(lifeSeconds)
	}
	return formatBattery(capacity, seconds, status)
}

func formatBattery(percent int, seconds int64, status []string) string {
	var b strings.Builder
	if percent >= 0 {
		b.WriteString(strconv.Itoa(percent))
		b.WriteByte('%')
	}
	if seconds > 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteByte('(')
		b.WriteString(formatDuration(uint64(seconds)))
		b.WriteString(" remaining)")
	}
	if len(status) == 0 {
		return b.String()
	}
	if b.Len() == 0 {
		return strings.Join(status, ", ")
	}
	b.WriteString(" [")
	b.WriteString(strings.Join(status, ", "))
	b.WriteByte(']')
	return b.String()
}

// formatDuration matches fastfetch's default duration text. Totals of a
// minute or more drop leftover seconds, and 30 seconds or more rounds up.
func formatDuration(totalSeconds uint64) string {
	if totalSeconds < 60 {
		if totalSeconds == 1 {
			return "1 second"
		}
		return strconv.FormatUint(totalSeconds, 10) + " seconds"
	}
	seconds := totalSeconds % 60
	totalSeconds /= 60
	if seconds >= 30 {
		totalSeconds++
	}
	minutes := totalSeconds % 60
	totalSeconds /= 60
	hours := totalSeconds % 24
	days := totalSeconds / 24

	var parts []string
	if days > 0 {
		parts = append(parts, unitWord(days, "day"))
	}
	if hours > 0 {
		parts = append(parts, unitWord(hours, "hour"))
	}
	if minutes > 0 {
		word := "min"
		if minutes > 1 {
			word = "mins"
		}
		parts = append(parts, strconv.FormatUint(minutes, 10)+" "+word)
	}
	if len(parts) == 0 {
		return "0 seconds"
	}
	return strings.Join(parts, ", ")
}

func unitWord(n uint64, word string) string {
	if n == 1 {
		return "1 " + word
	}
	return strconv.FormatUint(n, 10) + " " + word + "s"
}
