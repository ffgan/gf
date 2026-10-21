package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// IdentifyCPU names a CPU from /proc/cpuinfo text and an optional device-tree
// model. countsInName is set when the name already lists each cluster, so
// the formatter does not append another core count.
func IdentifyCPU(cpuinfo, deviceTree string) (name string, countsInName bool) {
	blocks := parseCPUInfo(cpuinfo)
	var model, vendor, hardware, uarch, isa string
	type partKey struct {
		impl uint32
		part uint32
	}
	var order []partKey
	counts := map[partKey]int{}
	sawPart := false

	for _, b := range blocks {
		if m := b["model name"]; m != "" && (model == "" || genericCPUName(model)) && !genericCPUName(m) {
			model = m
		} else if model == "" {
			model = m
		}
		if vendor == "" {
			vendor = b["vendor_id"]
		}
		if hardware == "" {
			hardware = b["hardware"]
		}
		if uarch == "" {
			uarch = b["uarch"]
		}
		if isa == "" {
			isa = b["isa"]
		}
		implS, partS := b["cpu implementer"], b["cpu part"]
		if implS == "" || partS == "" {
			continue
		}
		impl, ok1 := parseHexID(implS)
		part, ok2 := parseHexID(partS)
		if !ok1 || !ok2 {
			continue
		}
		sawPart = true
		k := partKey{impl, part}
		if counts[k] == 0 {
			order = append(order, k)
		}
		counts[k]++
	}

	if model != "" {
		mapped := applyFamilyCode(model)
		if !genericCPUName(mapped) {
			return applyVendorPrefix(vendor, mapped), false
		}
	}
	if sawPart {
		multi := len(order) > 1
		pieces := make([]string, 0, len(order))
		for _, k := range order {
			partName, ok := armPartName(k.impl, k.part)
			if !ok {
				partName = unknownPart(k.impl, k.part)
			}
			if multi {
				pieces = append(pieces, fmt.Sprintf("%dx %s", counts[k], partName))
			} else {
				pieces = append(pieces, partName)
			}
		}
		return strings.Join(pieces, " + "), multi
	}
	if uarch != "" {
		return riscvName(uarch), false
	}
	if soc := socFromText(hardware + "\n" + deviceTree); soc != "" {
		return soc, false
	}
	if hardware != "" && !genericCPUName(hardware) {
		return hardware, false
	}
	if dt := stripCtl(deviceTree); dt != "" && !genericCPUName(dt) {
		return dt, false
	}
	if v := applyVendorPrefix(vendor, ""); v != "" {
		return v, false
	}
	if isa != "" {
		return "RISC-V " + isa, false
	}
	return "", false
}

func unknownPart(impl, part uint32) string {
	if v := armVendor(impl); v != "" {
		return fmt.Sprintf("%s 0x%x", v, part)
	}
	return fmt.Sprintf("Unknown 0x%x:0x%x", impl, part)
}

func genericCPUName(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	low := strings.ToLower(s)
	switch low {
	case "unknown", "processor", "arm processor", "cpu", "n/a", "none", "not specified":
		return true
	}
	if strings.HasPrefix(low, "armv") && strings.Contains(low, "processor") {
		return true
	}
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func parseCPUInfo(text string) []map[string]string {
	var blocks []map[string]string
	cur := map[string]string{}
	flush := func() {
		if len(cur) == 0 {
			return
		}
		blocks = append(blocks, cur)
		cur = map[string]string{}
	}
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			flush()
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(k))
		val := strings.TrimSpace(v)
		if key == "processor" && len(cur) > 0 {
			flush()
		}
		cur[key] = val
	}
	flush()
	return blocks
}

func parseHexID(s string) (uint32, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.TrimPrefix(s, "0x")
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0, false
	}
	return uint32(n), true
}

func applyFamilyCode(model string) string {
	s := strings.TrimSpace(model)
	var n uint64
	var err error
	if strings.HasPrefix(strings.ToLower(s), "0x") {
		n, err = strconv.ParseUint(s[2:], 16, 32)
	} else {
		n, err = strconv.ParseUint(s, 10, 32)
	}
	if err != nil {
		return model
	}
	if name, ok := appleFamily[uint32(n)]; ok {
		return name
	}
	if name, ok := qcomFamily[uint32(n)]; ok {
		return name
	}
	return model
}

func applyVendorPrefix(vendor, model string) string {
	switch vendor {
	case "HygonGenuine":
		if model == "" || genericCPUName(model) {
			return "Hygon"
		}
		if !strings.Contains(strings.ToLower(model), "hygon") {
			return "Hygon " + model
		}
	case "CentaurHauls", "Shanghai":
		if model == "" || genericCPUName(model) {
			return "Zhaoxin"
		}
		if !strings.Contains(strings.ToLower(model), "zhaoxin") {
			return "Zhaoxin " + model
		}
	}
	return model
}

// Apple and Qualcomm family codes are the decimal values machdep.cpu.brand_string
// and Asahi publish. Names match the fastfetch table, including the shared
// M5 Pro / Max id.
var appleFamily = map[uint32]string{
	8103: "Apple M1",
	6000: "Apple M1 Pro",
	6001: "Apple M1 Max",
	6002: "Apple M1 Ultra",
	8112: "Apple M2",
	6020: "Apple M2 Pro",
	6021: "Apple M2 Max",
	6022: "Apple M2 Ultra",
	8122: "Apple M3",
	6030: "Apple M3 Pro",
	6031: "Apple M3 Max",
	6034: "Apple M3 Max",
	8132: "Apple M4",
	6040: "Apple M4 Pro",
	6041: "Apple M4 Max",
	8140: "Apple A18 Pro",
	8142: "Apple M5",
	6050: "Apple M5 Pro / Max",
}

var qcomFamily = map[uint32]string{
	7180: "Qualcomm Snapdragon 7c",
	7280: "Qualcomm Snapdragon 7c+ Gen 3",
	8180: "Qualcomm Snapdragon 8cx Gen 2 5G",
	8280: "Qualcomm Snapdragon 8cx Gen 3",
}

func riscvName(uarch string) string {
	u := strings.ToLower(strings.TrimSpace(uarch))
	u = strings.TrimSuffix(u, "-mc")
	if n, ok := riscvUarch[u]; ok {
		return n
	}
	vendor, core, ok := strings.Cut(u, ",")
	if !ok {
		return strings.TrimSpace(uarch)
	}
	return titleWord(vendor) + " " + strings.ToUpper(core)
}

var riscvUarch = map[string]string{
	"sifive,u74":          "SiFive U74",
	"sifive,u54":          "SiFive U54",
	"sifive,u54-mc":       "SiFive U54",
	"sifive,e31":          "SiFive E31",
	"sifive,e51":          "SiFive E51",
	"thead,c910":          "T-Head C910",
	"thead,c906":          "T-Head C906",
	"xiangshan,nanhu":     "XiangShan Nanhu",
	"xiangshan,kunminghu": "XiangShan Kunminghu",
	"spacemit,x60":        "SpacemiT X60",
	"sophgo,sg2042":       "Sophgo SG2042",
	"starfive,jh7110":     "StarFive JH7110",
	"andes,ax45":          "Andes AX45",
}

func titleWord(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func socFromText(s string) string {
	low := strings.ToLower(stripCtl(s))
	switch {
	case strings.Contains(low, "raspberry pi 5"), strings.Contains(low, "bcm2712"):
		return "Raspberry Pi 5"
	case strings.Contains(low, "raspberry pi 4"), strings.Contains(low, "bcm2711"):
		return "Raspberry Pi 4"
	case strings.Contains(low, "raspberry pi 3"), strings.Contains(low, "bcm2837"):
		return "Raspberry Pi 3"
	case strings.Contains(low, "raspberry pi zero 2"):
		return "Raspberry Pi Zero 2"
	case strings.Contains(low, "raspberry pi zero"):
		return "Raspberry Pi Zero"
	case strings.Contains(low, "rk3588"):
		return "Rockchip RK3588"
	case strings.Contains(low, "rk3399"):
		return "Rockchip RK3399"
	case strings.Contains(low, "rk3568"):
		return "Rockchip RK3568"
	case strings.Contains(low, "rk3566"):
		return "Rockchip RK3566"
	case strings.Contains(low, "rk3328"):
		return "Rockchip RK3328"
	case strings.Contains(low, "rk3288"):
		return "Rockchip RK3288"
	case strings.Contains(low, "th1520"):
		return "T-Head TH1520"
	case strings.Contains(low, "sg2042"):
		return "Sophgo SG2042"
	case strings.Contains(low, "jh7110"):
		return "StarFive JH7110"
	case strings.Contains(low, "jh7100"):
		return "StarFive JH7100"
	case strings.Contains(low, "spacemit"):
		if strings.Contains(low, "k1") {
			return "SpacemiT K1"
		}
		return "SpacemiT"
	case strings.Contains(low, "kunpeng"):
		return "HiSilicon Kunpeng"
	case strings.Contains(low, "kirin"):
		return "HiSilicon Kirin"
	case strings.Contains(low, "phytium"), strings.Contains(low, "ft-2000"), strings.Contains(low, "ft2000"), strings.Contains(low, "d2000"):
		if strings.Contains(low, "d2000") {
			return "Phytium D2000"
		}
		if strings.Contains(low, "ft-2000") || strings.Contains(low, "ft2000") {
			return "Phytium FT-2000"
		}
		return "Phytium"
	case strings.Contains(low, "s912"):
		return "Amlogic S912"
	case strings.Contains(low, "s905"):
		return "Amlogic S905"
	case strings.Contains(low, "amlogic"):
		return "Amlogic"
	case strings.Contains(low, "allwinner"), strings.Contains(low, "sun50i"), strings.Contains(low, "sun8i"):
		return "Allwinner"
	}
	return ""
}

func stripCtl(s string) string {
	if i := strings.IndexByte(s, 0); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}
