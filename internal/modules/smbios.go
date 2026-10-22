package modules

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// ParseMemoryDevices reads an SMBIOS structure table and returns one
// line per installed type-17 memory device.
func ParseMemoryDevices(table []byte) []string {
	var out []string
	for len(table) >= 4 {
		typ := table[0]
		length := int(table[1])
		if length < 4 || length > len(table) {
			break
		}
		body := table[:length]
		rest := table[length:]
		end := 0
		for end+1 < len(rest) {
			if rest[end] == 0 && rest[end+1] == 0 {
				end += 2
				break
			}
			end++
		}
		stringsRaw := rest[:end]
		table = rest[end:]
		if typ != 17 {
			continue
		}
		if line := formatMemoryDevice(body, stringsRaw); line != "" {
			out = append(out, line)
		}
	}
	return out
}

func formatMemoryDevice(body, rawStrings []byte) string {
	if len(body) < 0x0E {
		return ""
	}
	size := binary.LittleEndian.Uint16(body[0x0C:])
	if size == 0 {
		return ""
	}
	var bytes uint64
	switch {
	case size == 0x7FFF && len(body) >= 0x20:
		bytes = uint64(binary.LittleEndian.Uint32(body[0x1C:])) * 1024 * 1024
	case size&0x8000 != 0:
		bytes = uint64(size&^0x8000) * 1024
	default:
		bytes = uint64(size) * 1024 * 1024
	}
	if bytes == 0 {
		return ""
	}
	kind := ""
	if len(body) > 0x12 {
		kind = memoryTypeName(body[0x12])
	}
	speed := 0
	if len(body) >= 0x17 {
		speed = int(binary.LittleEndian.Uint16(body[0x15:]))
	}
	if speed == 0 && len(body) >= 0x22 {
		speed = int(binary.LittleEndian.Uint16(body[0x20:]))
	}
	locator := smbiosString(rawStrings, byteAt(body, 0x10))
	vendor := smbiosString(rawStrings, byteAt(body, 0x17))
	form := ""
	if len(body) > 0x0E {
		form = formFactorName(body[0x0E])
	}

	parts := []string{FormatBytes(bytes)}
	desc := kind
	if speed > 0 {
		if desc != "" {
			desc += fmt.Sprintf("-%d", speed)
		} else {
			desc = fmt.Sprintf("%d MT/s", speed)
		}
	}
	if desc != "" {
		parts = append(parts, desc)
	}
	if form != "" && locator == "" {
		parts = append(parts, form)
	}
	if locator != "" {
		parts = append(parts, locator)
	}
	if vendor != "" && !strings.EqualFold(vendor, "unknown") {
		parts = append(parts, "("+vendor+")")
	}
	return strings.Join(parts, " ")
}

func byteAt(body []byte, off int) byte {
	if off < 0 || off >= len(body) {
		return 0
	}
	return body[off]
}

func smbiosString(raw []byte, index byte) string {
	if index == 0 {
		return ""
	}
	parts := strings.Split(string(raw), "\x00")
	if int(index) > len(parts) {
		return ""
	}
	return strings.TrimSpace(parts[index-1])
}

func memoryTypeName(n byte) string {
	names := map[byte]string{
		0x01: "Other",
		0x02: "Unknown",
		0x03: "DRAM",
		0x07: "RAM",
		0x0F: "SDRAM",
		0x12: "DDR",
		0x13: "DDR2",
		0x18: "DDR3",
		0x1A: "DDR4",
		0x1B: "LPDDR",
		0x1C: "LPDDR2",
		0x1D: "LPDDR3",
		0x1E: "LPDDR4",
		0x22: "DDR5",
		0x23: "LPDDR5",
	}
	return names[n]
}

func formFactorName(n byte) string {
	names := map[byte]string{
		0x03: "SIMM",
		0x09: "DIMM",
		0x0D: "SODIMM",
		0x0F: "FB-DIMM",
	}
	return names[n]
}

// ChassisType maps an SMBIOS chassis type byte. The high bit is a
// lock flag and is ignored.
func ChassisType(n int) string {
	switch n & 0x7f {
	case 0x01:
		return "Other"
	case 0x02:
		return "Unknown"
	case 0x03:
		return "Desktop"
	case 0x04:
		return "Low Profile Desktop"
	case 0x05:
		return "Pizza Box"
	case 0x06:
		return "Mini Tower"
	case 0x07:
		return "Tower"
	case 0x08:
		return "Portable"
	case 0x09:
		return "Laptop"
	case 0x0A:
		return "Notebook"
	case 0x0B:
		return "Hand Held"
	case 0x0C:
		return "Docking Station"
	case 0x0D:
		return "All in One"
	case 0x0E:
		return "Sub Notebook"
	case 0x0F:
		return "Space-saving"
	case 0x10:
		return "Lunch Box"
	case 0x11:
		return "Main Server Chassis"
	case 0x12:
		return "Expansion Chassis"
	case 0x13:
		return "SubChassis"
	case 0x14:
		return "Bus Expansion Chassis"
	case 0x15:
		return "Peripheral Chassis"
	case 0x16:
		return "RAID Chassis"
	case 0x17:
		return "Rack Mount Chassis"
	case 0x18:
		return "Sealed-case PC"
	case 0x19:
		return "Multi-system chassis"
	case 0x1A:
		return "Compact PCI"
	case 0x1B:
		return "Advanced TCA"
	case 0x1C:
		return "Blade"
	case 0x1D:
		return "Blade Enclosure"
	case 0x1E:
		return "Tablet"
	case 0x1F:
		return "Convertible"
	case 0x20:
		return "Detachable"
	case 0x21:
		return "IoT Gateway"
	case 0x22:
		return "Embedded PC"
	case 0x23:
		return "Mini PC"
	case 0x24:
		return "Stick PC"
	default:
		return ""
	}
}
