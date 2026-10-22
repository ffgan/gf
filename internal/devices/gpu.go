package dev

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	cli "github.com/ffgan/gf/internal/CLI"
	"github.com/ffgan/gf/internal/utils"
)

func GetGPU(osName, showBrand, gpuType string) string {
	var gpu string

	switch osName {
	case cli.Linux:
		gpu = getLinuxGPU()

	case cli.MacOSX, cli.MacOS, cli.Darwin:
		gpu = getMacGPU()

	case cli.Iphone:

	case cli.Windows:
		gpu = getWindowsGPU()

	case cli.FreeBSD, cli.DragonFly, "MidnightBSD", "HardenedBSD", "GhostBSD", "NomadBSD":
		out := utils.RunCommand("pciconf", "-lv")
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "device") {
				gpu = utils.Trim(strings.SplitN(line, "=", 2)[1])
				break
			}
		}

	default:
		// Fallback: Try OpenGL
		out := utils.RunCommand("glxinfo", "-B")
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "OpenGL renderer string") {
				gpu = strings.TrimPrefix(line, "OpenGL renderer string:")
				gpu = utils.Trim(gpu)
				break
			}
		}
	}

	if showBrand == utils.OFF {
		gpu = strings.ReplaceAll(gpu, "AMD", "")
		gpu = strings.ReplaceAll(gpu, "NVIDIA", "")
		gpu = strings.ReplaceAll(gpu, "Intel", "")
		gpu = strings.TrimSpace(gpu)
	}
	return filterGPUType(gpu, gpuType)
}

func filterGPUType(gpu, kind string) string {
	if gpu == "" || kind == "" || kind == "all" {
		return gpu
	}
	var keep []string
	for _, part := range splitGPUList(gpu) {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		integrated := strings.Contains(part, "Intel") || strings.Contains(strings.ToLower(part), "integrated")
		switch kind {
		case "integrated":
			if integrated {
				keep = append(keep, part)
			}
		case "dedicated":
			if !integrated {
				keep = append(keep, part)
			}
		default:
			keep = append(keep, part)
		}
	}
	return strings.Join(keep, "\n")
}

func splitGPUList(gpu string) []string {
	var parts []string
	for _, line := range strings.Split(gpu, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "[") {
			parts = append(parts, line)
			continue
		}
		for _, part := range strings.Split(line, ", ") {
			part = strings.TrimSpace(part)
			if part != "" {
				parts = append(parts, part)
			}
		}
	}
	return parts
}

func getLinuxGPU() string {
	// SoCs such as the Raspberry Pi expose a GPU through DRM and not PCI.
	// Name those from sysfs before starting lspci or a GL context.
	if name, noPCIDisplay := gpuFromSysfs(); name != "" {
		return name
	} else if drm := gpuFromDRM(); drm != "" {
		return drm
	} else if noPCIDisplay {
		return ""
	}
	out := utils.RunCommand("lspci", "-mm")
	if out == "" {
		return parseGLXInfo()
	}

	lines := strings.Split(out, "\n")
	seen := map[string]bool{}
	var gpus []string

	re := regexp.MustCompile(`"([^"]+)"`)

	for _, line := range lines {
		if strings.Contains(line, "Display") || strings.Contains(line, "3D") || strings.Contains(line, "VGA") {
			matches := re.FindAllStringSubmatch(line, -1)
			if len(matches) >= 3 {
				// matches[0][1] = 第1个引号内容 (类型,如 "VGA compatible controller")
				// matches[1][1] = 第2个引号内容 (厂商,如 "Advanced Micro Devices...")
				// matches[2][1] = 第3个引号内容 (设备名,如 "Cezanne...")

				vendor := matches[1][1] // $3 in awk
				var device string

				// 检查倒数第二个字段是否为空或是 "Device xxxx" 格式
				if len(matches) >= 4 {
					lastField := matches[len(matches)-1][1]
					// 如果最后一个字段是 "Device [hex]" 格式,使用 $4
					if strings.HasPrefix(lastField, "Device ") {
						device = matches[2][1] // $4
					} else {
						device = lastField // $(NF-1)
					}
				} else {
					device = matches[2][1]
				}

				info := vendor + " " + device

				if !seen[info] {
					seen[info] = true
					gpus = append(gpus, info)
				}
			}
		}
	}

	if len(gpus) == 0 {
		return parseGLXInfo()
	}

	// Normalize brand names
	var normalizedGPUs []string
	for _, gpu := range gpus {
		original := gpu
		var brand string

		switch {
		case strings.Contains(gpu, "Advanced"):
			// Determine AMD/ATI brand
			if strings.Contains(gpu, "AMD") && strings.Contains(gpu, "ATI") {
				brand = "AMD ATI"
			} else if strings.Contains(gpu, "AMD") {
				brand = "AMD"
			} else if strings.Contains(gpu, "ATI") {
				brand = "ATI"
			}

			// Clean up GPU name
			gpu = strings.ReplaceAll(gpu, "[AMD/ATI]", "")
			gpu = strings.ReplaceAll(gpu, "[AMD]", "")
			gpu = strings.ReplaceAll(gpu, "OEM", "")
			gpu = strings.ReplaceAll(gpu, "Advanced Micro Devices, Inc.", "")

			// Extract content within brackets
			if idx := strings.Index(gpu, "["); idx != -1 {
				gpu = gpu[idx+1:]
				if idx := strings.Index(gpu, "]"); idx != -1 {
					gpu = gpu[:idx]
				}
			}

			gpu = brand + " " + strings.TrimSpace(gpu)

		case strings.Contains(gpu, "NVIDIA"):
			// Remove everything before "NVIDIA"
			if idx := strings.Index(gpu, "NVIDIA"); idx != -1 {
				gpu = gpu[idx:]
			}

			// Extract content within brackets
			if idx := strings.Index(gpu, "["); idx != -1 {
				gpu = gpu[idx+1:]
				if idx := strings.Index(gpu, "]"); idx != -1 {
					gpu = gpu[:idx]
				}
			}

			gpu = "NVIDIA " + strings.TrimSpace(gpu)

		case strings.Contains(gpu, "Intel"):
			// Keep from "Intel" onwards
			if idx := strings.Index(gpu, "Intel"); idx != -1 {
				gpu = gpu[idx:]
			}

			gpu = strings.ReplaceAll(gpu, "(R)", "")
			gpu = strings.ReplaceAll(gpu, "Corporation", "")
			gpu = strings.ReplaceAll(gpu, "Integrated Graphics Controller", "")

			// Remove content in parentheses
			if idx := strings.Index(gpu, "("); idx != -1 {
				gpu = gpu[:idx]
			}

			// Handle Xeon case
			if strings.Contains(original, "Xeon") {
				gpu = "Intel HD Graphics"
			}

			gpu = strings.TrimSpace(gpu)
			if gpu == "" || gpu == "Intel" {
				gpu = "Intel Integrated Graphics"
			}

		case strings.Contains(gpu, "MCST") && strings.Contains(gpu, "MGA2"):
			gpu = "MCST MGA2"

		case strings.Contains(gpu, "VirtualBox"):
			gpu = "VirtualBox Graphics Adapter"

		default:
			// Skip unrecognized GPUs (like the bash 'continue')
			continue
		}

		// Optional: Remove brand prefix if gpu_brand is utils.OFF
		// if c.config.GPUBrand == utils.OFF {
		//     gpu = strings.TrimPrefix(gpu, "AMD ")
		//     gpu = strings.TrimPrefix(gpu, "NVIDIA ")
		//     gpu = strings.TrimPrefix(gpu, "Intel ")
		// }

		normalizedGPUs = append(normalizedGPUs, gpu)
	}

	if len(normalizedGPUs) == 0 {
		return parseGLXInfo()
	}

	return strings.Join(normalizedGPUs, ", ")
}

func parseGLXInfo() string {
	out := utils.RunCommand("glxinfo", "-B")
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "OpenGL renderer string") {
			return utils.Trim(strings.TrimPrefix(line, "OpenGL renderer string:"))
		}
	}
	return ""
}

func getMacGPU() string {
	// macOS ARM (Apple Silicon)
	if utils.RunCommand("uname", "-m") == "arm64" {
		chip := utils.RunCommand("system_profiler", "SPDisplaysDataType")
		chipset := parseSystemProfiler(chip, "Chipset Model")
		cores := parseSystemProfiler(chip, "Total Number of Cores")
		if chipset != "" {
			return fmt.Sprintf("%s (%s cores)", chipset, cores)
		}
	}

	// Intel mac
	out := utils.RunCommand("system_profiler", "SPDisplaysDataType")
	return parseSystemProfiler(out, "Chipset Model")
}

func parseSystemProfiler(data, key string) string {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, key) {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				buf.WriteString(strings.TrimSpace(parts[1]))
				buf.WriteString(", ")
			}
		}
	}
	result := buf.String()
	return strings.TrimSuffix(result, ", ")
}

func getWindowsGPU() string {
	// reg.exe /s prints every binary value under the display class. On a
	// machine with virtual adapters that dump is large enough to dominate
	// startup. The registry API reads DriverDesc and the memory size only.
	if items, ok := windowsRegistryGPUs(); ok {
		return formatWindowsGPUItems(items)
	}
	out := utils.RunCommand("reg.exe", "query", `HKLM\SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}`, "/s")
	return formatWindowsGPUs(out)
}

func windowsRegistryGPUs() ([]windowsGPU, bool) {
	const class = `SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}`
	subs, ok := utils.RegSubkeys(utils.HKLM, class)
	if !ok {
		return nil, false
	}
	var items []windowsGPU
	for _, sub := range subs {
		if !digitsOnly(sub) {
			continue
		}
		path := class + `\` + sub
		name, ok := utils.RegString(utils.HKLM, path, "DriverDesc")
		if !ok || name == "" {
			continue
		}
		mem, _ := utils.RegUint64(utils.HKLM, path, "HardwareInformation.qwMemorySize")
		if mem == 0 {
			mem, _ = utils.RegUint64(utils.HKLM, path, "HardwareInformation.MemorySize")
		}
		items = append(items, windowsGPU{name: name, mem: mem})
	}
	return items, true
}

func digitsOnly(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func formatWindowsGPUs(out string) string {
	return formatWindowsGPUItems(parseWindowsGPUs(out))
}

func formatWindowsGPUItems(items []windowsGPU) string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.name)
	}
	keep := map[string]bool{}
	for _, name := range preferRealGPUs(names) {
		keep[name] = true
	}
	var lines []string
	seen := map[string]bool{}
	for _, item := range items {
		if !keep[item.name] || seen[item.name] {
			continue
		}
		seen[item.name] = true
		lines = append(lines, formatWindowsGPULine(item.name, item.mem))
	}
	return strings.Join(lines, "\n")
}

type windowsGPU struct {
	name string
	mem  uint64
}

func parseWindowsGPUs(out string) []windowsGPU {
	var items []windowsGPU
	var cur *windowsGPU
	var pending uint64
	flush := func() {
		if cur != nil && cur.name != "" {
			items = append(items, *cur)
		}
		cur = nil
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "HKEY_") {
			flush()
			pending = 0
			continue
		}
		if name := sectionValue(line, "DriverDesc"); name != "" {
			flush()
			cur = &windowsGPU{name: name, mem: pending}
			pending = 0
			continue
		}
		mem := regUint64(sectionValue(line, "HardwareInformation.qwMemorySize"))
		if mem == 0 {
			mem = regUint64(sectionValue(line, "HardwareInformation.MemorySize"))
		}
		if mem == 0 {
			continue
		}
		if cur != nil && cur.mem == 0 {
			cur.mem = mem
		} else if cur == nil {
			pending = mem
		}
	}
	flush()
	return items
}

func sectionValue(section, name string) string {
	for _, line := range strings.Split(section, "\n") {
		fields := strings.Fields(line)
		reg := -1
		for i, field := range fields {
			if strings.HasPrefix(field, "REG_") {
				reg = i
				break
			}
		}
		if reg <= 0 || reg+1 >= len(fields) {
			continue
		}
		if strings.Join(fields[:reg], " ") == name {
			return strings.Join(fields[reg+1:], " ")
		}
	}
	return ""
}

func regUint64(s string) uint64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	base := 10
	if strings.HasPrefix(strings.ToLower(s), "0x") {
		s = s[2:]
		base = 16
	}
	n, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		return 0
	}
	return n
}

func formatWindowsGPULine(name string, mem uint64) string {
	kind := windowsGPUKind(name)
	if kind == "Virtual" {
		return name
	}
	line := name
	if mem > 0 {
		line += " (" + formatGPUMem(mem) + ")"
	}
	if kind != "" {
		line += " [" + kind + "]"
	}
	return line
}

func windowsGPUKind(name string) string {
	low := strings.ToLower(name)
	switch {
	case strings.Contains(low, "virtual"), strings.Contains(low, "mumu"),
		strings.Contains(low, "basic display"), strings.Contains(low, "remote display"),
		strings.Contains(low, "indirect"):
		return "Virtual"
	case strings.Contains(low, "arc"):
		return "Discrete"
	case strings.Contains(low, "intel"):
		return "Integrated"
	case strings.Contains(low, "nvidia"), strings.Contains(low, "geforce"),
		strings.Contains(low, "quadro"), strings.Contains(low, "radeon"),
		strings.Contains(low, "amd"):
		return "Discrete"
	default:
		return ""
	}
}

func formatGPUMem(n uint64) string {
	const giB = 1024 * 1024 * 1024
	const miB = 1024 * 1024
	if n >= giB {
		return fmt.Sprintf("%.2f GiB", float64(n)/float64(giB))
	}
	if n >= miB {
		return fmt.Sprintf("%.2f MiB", float64(n)/float64(miB))
	}
	if n == 0 {
		return ""
	}
	return fmt.Sprintf("%d KiB", n/1024)
}

func preferRealGPUs(gpus []string) []string {
	var real []string
	for _, gpu := range gpus {
		low := strings.ToLower(gpu)
		if strings.Contains(low, "microsoft basic display") || strings.Contains(low, "microsoft remote display") {
			continue
		}
		real = append(real, gpu)
	}
	if len(real) > 0 {
		return real
	}
	return gpus
}
