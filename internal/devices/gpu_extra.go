package dev

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
		if fallback := namesFromVendors(found); fallback != "" {
			return fallback, true
		}
		return "", false
	}
	var out []string
	seen := map[string]struct{}{}
	for _, id := range found {
		name := names[id]
		if name == "" {
			continue
		}
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	if joined := strings.Join(out, ", "); joined != "" {
		return joined, true
	}
	if fallback := namesFromVendors(found); fallback != "" {
		return fallback, true
	}
	return "", false
}

type pciID struct {
	vendor string
	device string
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
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var vendorID, vendorName string
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "\t\t") {
			continue
		}
		if strings.HasPrefix(line, "\t") {
			if vendorID == "" || needed[vendorID] == nil {
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
			names[pciID{vendorID, devID}] = vendorName + " " + strings.TrimSpace(rest[4:])
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

func namesFromVendors(ids []pciID) string {
	var out []string
	seen := map[string]struct{}{}
	for _, id := range ids {
		vendor, ok := knownPCIVendors[id.vendor]
		if !ok {
			continue
		}
		label := vendor
		switch vendor {
		case "Intel", "AMD", "NVIDIA":
			label = vendor + " Graphics"
		case "VirtIO":
			label = "VirtIO GPU"
		case "QEMU":
			label = "QEMU VGA"
		case "VMware":
			label = "VMware SVGA"
		case "VirtualBox":
			label = "VirtualBox Graphics Adapter"
		case "Hyper-V":
			label = "Hyper-V Video"
		}
		label = fmt.Sprintf("%s (%s:%s)", label, id.vendor, id.device)
		if _, ok := seen[label]; ok {
			continue
		}
		seen[label] = struct{}{}
		out = append(out, label)
	}
	return strings.Join(out, ", ")
}

var knownPCIVendors = map[string]string{
	"8086": "Intel",
	"1002": "AMD",
	"1022": "AMD",
	"10de": "NVIDIA",
	"1af4": "VirtIO",
	"1234": "QEMU",
	"15ad": "VMware",
	"80ee": "VirtualBox",
	"1414": "Hyper-V",
	"1ed5": "Moore Threads",
	"1ee0": "Biren",
	"0014": "Loongson",
	"19e5": "HiSilicon",
	"1d94": "Hygon",
	"1d17": "Zhaoxin",
	"14e4": "Broadcom",
	"17cb": "Qualcomm",
}

func gpuFromDRM() string {
	entries, err := os.ReadDir("/sys/class/drm")
	if err != nil {
		return ""
	}
	var names []string
	seen := map[string]struct{}{}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "card") || strings.Contains(e.Name(), "-") {
			continue
		}
		link, err := os.Readlink(filepath.Join("/sys/class/drm", e.Name(), "device", "driver"))
		if err != nil {
			continue
		}
		name, ok := drmDriverNames[filepath.Base(link)]
		if !ok {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

var drmDriverNames = map[string]string{
	"i915":       "Intel Graphics",
	"xe":         "Intel Graphics",
	"amdgpu":     "AMD Graphics",
	"radeon":     "AMD Graphics",
	"nouveau":    "NVIDIA Graphics",
	"nvidia":     "NVIDIA Graphics",
	"panfrost":   "Mali",
	"panthor":    "Mali",
	"lima":       "Mali",
	"mali":       "Mali",
	"vc4":        "VideoCore",
	"v3d":        "VideoCore",
	"msm":        "Adreno",
	"virtio_gpu": "VirtIO GPU",
	"vmwgfx":     "VMware SVGA",
	"bochs-drm":  "Bochs",
	"qxl":        "QXL",
	"apple":      "Apple Graphics",
	"loongson":   "Loongson Graphics",
	"etnaviv":    "Vivante",
	"imx-drm":    "i.MX",
	"rockchip":   "Rockchip",
	"meson":      "Amlogic",
	"sun4i-drm":  "Allwinner",
	"ast":        "ASPEED",
	"mgag200":    "Matrox",
	"hibmc":      "HiSilicon",
	"phytium":    "Phytium",
}
