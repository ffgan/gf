package dev

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
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
	for _, part := range strings.Split(gpu, ", ") {
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
	return strings.Join(keep, ", ")
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
	out := utils.RunCommand("reg.exe", "query", `HKLM\SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}`, "/s", "/v", "DriverDesc")
	return formatWindowsGPUs(out)
}

func formatWindowsGPUs(out string) string {
	var gpus []string
	seen := map[string]bool{}
	for _, name := range regAll(out, "DriverDesc") {
		name = strings.TrimSpace(name)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		gpus = append(gpus, name)
	}
	gpus = preferRealGPUs(gpus)
	return strings.Join(gpus, ", ")
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

func regAll(out, name string) []string {
	var vals []string
	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, name) || !strings.Contains(line, "REG_") {
			continue
		}
		fields := strings.Fields(line)
		for i, field := range fields {
			if strings.HasPrefix(field, "REG_") && i+1 < len(fields) {
				vals = append(vals, strings.Join(fields[i+1:], " "))
				break
			}
		}
	}
	return vals
}
