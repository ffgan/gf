package dev

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	cli "github.com/ffgan/gf/internal/CLI"
	"github.com/ffgan/gf/internal/utils"
)

func GetGPU(osName string) string {
	var gpu string

	switch osName {
	case cli.Linux:
		gpu = getLinuxGPU()

	case cli.MacOSX, cli.MacOS:
		gpu = getMacGPU()

	case cli.Iphone:

	case cli.Windows:
		gpu = getWindowsGPU()

	case cli.FreeBSD, cli.DragonFly:
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

	// Optional: remove brand names if disabled
	if os.Getenv("GPU_BRAND") == utils.OFF {
		gpu = strings.ReplaceAll(gpu, "AMD", "")
		gpu = strings.ReplaceAll(gpu, "NVIDIA", "")
		gpu = strings.ReplaceAll(gpu, "Intel", "")
	}

	return gpu
}

func getLinuxGPU() string {
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
	out := utils.RunCommand("wmic", "path", "Win32_VideoController", "get", "caption")
	var gpus []string
	for _, line := range strings.Split(out, "\n") {
		line = utils.Trim(line)
		if line == "" || strings.Contains(line, "Caption") {
			continue
		}
		gpus = append(gpus, line)
	}
	return strings.Join(gpus, ", ")
}
