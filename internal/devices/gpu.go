package dev

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func trim(s string) string {
	return strings.TrimSpace(s)
}

func runCommand(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getGPU() string {
	osName := detectOS()
	var gpu string

	switch osName {
	case "Linux":
		gpu = getLinuxGPU()

	case "Mac OS X", "macOS":
		gpu = getMacGPU()

	case "Windows":
		gpu = getWindowsGPU()

	case "FreeBSD", "DragonFly":
		out := runCommand("pciconf", "-lv")
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "device") {
				gpu = trim(strings.SplitN(line, "=", 2)[1])
				break
			}
		}

	default:
		// Fallback: Try OpenGL
		out := runCommand("glxinfo", "-B")
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "OpenGL renderer string") {
				gpu = strings.TrimPrefix(line, "OpenGL renderer string:")
				gpu = trim(gpu)
				break
			}
		}
	}

	// Optional: remove brand names if disabled
	if os.Getenv("GPU_BRAND") == "off" {
		gpu = strings.ReplaceAll(gpu, "AMD", "")
		gpu = strings.ReplaceAll(gpu, "NVIDIA", "")
		gpu = strings.ReplaceAll(gpu, "Intel", "")
	}

	return gpu
}

func getLinuxGPU() string {
	out := runCommand("lspci", "-mm")
	if out == "" {
		return parseGLXInfo()
	}

	// Parse "Display", "3D", "VGA"
	lines := strings.Split(out, "\n")
	seen := map[string]bool{}
	var gpus []string

	re := regexp.MustCompile(`"([^"]+)"`)
	for _, line := range lines {
		if strings.Contains(line, "Display") || strings.Contains(line, "3D") || strings.Contains(line, "VGA") {
			fields := re.FindAllString(line, -1)
			if len(fields) >= 3 {
				info := strings.Join(fields[1:], " ")
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

	// Try to normalize brand names
	for i, gpu := range gpus {
		switch {
		case strings.Contains(gpu, "AMD"), strings.Contains(gpu, "ATI"):
			gpu = strings.ReplaceAll(gpu, "[AMD/ATI]", "")
			gpu = strings.ReplaceAll(gpu, "Advanced Micro Devices, Inc.", "")
			gpus[i] = "AMD " + trim(gpu)
		case strings.Contains(gpu, "NVIDIA"):
			gpu = strings.ReplaceAll(gpu, "[NVIDIA]", "")
			gpus[i] = "NVIDIA " + trim(gpu)
		case strings.Contains(gpu, "Intel"):
			gpu = strings.ReplaceAll(gpu, "Corporation", "")
			gpus[i] = "Intel " + trim(gpu)
		case strings.Contains(gpu, "VirtualBox"):
			gpus[i] = "VirtualBox Graphics Adapter"
		}
	}

	return strings.Join(gpus, ", ")
}

func parseGLXInfo() string {
	out := runCommand("glxinfo", "-B")
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "OpenGL renderer string") {
			return trim(strings.TrimPrefix(line, "OpenGL renderer string:"))
		}
	}
	return ""
}

func getMacGPU() string {
	// macOS ARM (Apple Silicon)
	if runCommand("uname", "-m") == "arm64" {
		chip := runCommand("system_profiler", "SPDisplaysDataType")
		chipset := parseSystemProfiler(chip, "Chipset Model")
		cores := parseSystemProfiler(chip, "Total Number of Cores")
		if chipset != "" {
			return fmt.Sprintf("%s (%s cores)", chipset, cores)
		}
	}

	// Intel mac
	out := runCommand("system_profiler", "SPDisplaysDataType")
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
	out := runCommand("wmic", "path", "Win32_VideoController", "get", "caption")
	var gpus []string
	for _, line := range strings.Split(out, "\n") {
		line = trim(line)
		if line == "" || strings.Contains(line, "Caption") {
			continue
		}
		gpus = append(gpus, line)
	}
	return strings.Join(gpus, ", ")
}
