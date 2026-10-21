package dev

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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
	// PCI class 0x03 is a display controller. Reading sysfs is what
	// fastfetch does; lspci and glxinfo each start a process, and glxinfo
	// initialises a GL context.
	if name, ok := gpuFromSysfs(); ok {
		return name
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

func gpuFromSysfs() (string, bool) {
	entries, err := os.ReadDir("/sys/bus/pci/devices")
	if err != nil {
		return "", false
	}

	var found []pciID
	seenID := map[pciID]struct{}{}
	for _, e := range entries {
		base := filepath.Join("/sys/bus/pci/devices", e.Name())
		class := strings.ToLower(strings.TrimSpace(readPCI(filepath.Join(base, "class"))))
		// class is "0x030000" for VGA, "0x030200" for 3D, "0x038000" for other display.
		if !strings.HasPrefix(class, "0x03") {
			continue
		}
		id := pciID{
			vendor: normPCI(readPCI(filepath.Join(base, "vendor"))),
			device: normPCI(readPCI(filepath.Join(base, "device"))),
		}
		if id.vendor == "" || id.device == "" {
			continue
		}
		if _, ok := seenID[id]; ok {
			continue
		}
		seenID[id] = struct{}{}
		found = append(found, id)
	}
	if len(found) == 0 {
		return "", true
	}

	names := lookupPCINames(found)
	if names == nil {
		// sysfs saw a GPU but we have no pci.ids to name it. Let lspci try.
		return "", false
	}
	var out []string
	seen := map[string]struct{}{}
	for _, id := range found {
		name := names[id]
		if name == "" {
			continue
		}
		name = normalizeGPUName(name)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return strings.Join(out, ", "), true
}

func readPCI(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func normPCI(id string) string {
	id = strings.ToLower(strings.TrimSpace(id))
	id = strings.TrimPrefix(id, "0x")
	if len(id) > 4 {
		id = id[len(id)-4:]
	}
	return id
}

type pciID struct {
	vendor string
	device string
}

func lookupPCINames(ids []pciID) map[pciID]string {
	needed := map[string]map[string]struct{}{}
	for _, id := range ids {
		if needed[id.vendor] == nil {
			needed[id.vendor] = map[string]struct{}{}
		}
		needed[id.vendor][id.device] = struct{}{}
	}

	var f *os.File
	for _, path := range []string{
		"/usr/share/hwdata/pci.ids",
		"/usr/share/misc/pci.ids",
		"/usr/share/pci.ids",
	} {
		opened, err := os.Open(path)
		if err == nil {
			f = opened
			break
		}
	}
	if f == nil {
		return nil
	}
	defer f.Close()

	names := map[pciID]string{}
	sc := bufio.NewScanner(f)
	// pci.ids lines are long but bounded. Raise the default 64K just in case.
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var vendorID, vendorName string
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "\t\t") {
			continue
		}
		if strings.HasPrefix(line, "\t") {
			if vendorID == "" {
				continue
			}
			if _, wantVendor := needed[vendorID]; !wantVendor {
				continue
			}
			rest := line[1:]
			if len(rest) < 4 {
				continue
			}
			devID := strings.ToLower(rest[:4])
			if _, want := needed[vendorID][devID]; !want {
				continue
			}
			devName := strings.TrimSpace(rest[4:])
			names[pciID{vendorID, devID}] = vendorName + " " + devName
			continue
		}
		if len(line) < 4 {
			continue
		}
		vendorID = strings.ToLower(line[:4])
		if _, want := needed[vendorID]; want {
			vendorName = strings.TrimSpace(line[4:])
		} else {
			vendorID = ""
			vendorName = ""
		}
	}
	return names
}

func normalizeGPUName(raw string) string {
	gpu := raw
	switch {
	case strings.Contains(gpu, "Advanced") || strings.Contains(gpu, "AMD") || strings.Contains(gpu, "ATI"):
		brand := "AMD"
		if strings.Contains(gpu, "AMD") && strings.Contains(gpu, "ATI") {
			brand = "AMD ATI"
		} else if strings.Contains(gpu, "ATI") && !strings.Contains(gpu, "AMD") {
			brand = "ATI"
		}
		gpu = strings.ReplaceAll(gpu, "[AMD/ATI]", "")
		gpu = strings.ReplaceAll(gpu, "[AMD]", "")
		gpu = strings.ReplaceAll(gpu, "Advanced Micro Devices, Inc.", "")
		if idx := strings.Index(gpu, "["); idx != -1 {
			gpu = gpu[idx+1:]
			if end := strings.Index(gpu, "]"); end != -1 {
				gpu = gpu[:end]
			}
		}
		gpu = brand + " " + strings.TrimSpace(gpu)
	case strings.Contains(gpu, "NVIDIA"):
		if idx := strings.Index(gpu, "["); idx != -1 {
			gpu = gpu[idx+1:]
			if end := strings.Index(gpu, "]"); end != -1 {
				gpu = gpu[:end]
			}
		}
		gpu = "NVIDIA " + strings.TrimSpace(gpu)
	case strings.Contains(gpu, "Intel"):
		if idx := strings.Index(gpu, "Intel"); idx != -1 {
			gpu = gpu[idx:]
		}
		gpu = strings.ReplaceAll(gpu, "(R)", "")
		gpu = strings.ReplaceAll(gpu, "Corporation", "")
		gpu = strings.ReplaceAll(gpu, "Integrated Graphics Controller", "")
		if idx := strings.Index(gpu, "("); idx != -1 {
			gpu = gpu[:idx]
		}
		gpu = strings.TrimSpace(gpu)
		if gpu == "" || gpu == "Intel" {
			gpu = "Intel Integrated Graphics"
		}
	case strings.Contains(gpu, "VirtualBox"):
		gpu = "VirtualBox Graphics Adapter"
	default:
		return strings.TrimSpace(raw)
	}
	return strings.Join(strings.Fields(gpu), " ")
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
