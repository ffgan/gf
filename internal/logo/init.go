package logo

import "github.com/armon/go-radix"

var r *radix.Tree

func init() {
	r = radix.New()

	r.Insert("Adélie", DistroWithColor{
		FileName: "adelie",
		Colors:   []int{4, 7, 6},
	})
	r.Insert("Adelie", DistroWithColor{
		FileName: "adelie",
		Colors:   []int{4, 7, 6},
	})

	r.Insert("AIX", DistroWithColor{
		FileName: "AIX",
		Colors:   []int{2, 7},
	})

	r.Insert("Aperio GNU/Linux", DistroWithColor{
		FileName: "Aperio",
		Colors:   []int{255},
	})

	r.Insert("Aperture", DistroWithColor{
		FileName: "Aperture",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("Asahi", DistroWithColor{
		FileName: "Asahi",
		Colors:   []int{3, 2, 1, 8, 7, 6, 4},
	})

	r.Insert("Hash", DistroWithColor{
		FileName: "Hash",
		Colors:   []int{123},
	})

	r.Insert("HarDClanZ", DistroWithColor{
		FileName: "HarDClanZ",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("AlmaLinux", DistroWithColor{
		FileName: "AlmaLinux",
		Colors:   []int{1, 3, 4, 2, 6},
	})

	r.Insert("Exodia Predator OS", DistroWithColor{
		FileName: "exodia-predator",
		Colors:   []int{5, 5},
	})
	r.Insert("exodia-predator", DistroWithColor{
		FileName: "exodia-predator",
		Colors:   []int{5, 5},
	})
	r.Insert("Predator", DistroWithColor{
		FileName: "exodia-predator",
		Colors:   []int{5, 5},
	})

	r.Insert("alpine_small", DistroWithColor{
		FileName: "alpine_small",
		Colors:   []int{4, 7},
	})

	r.Insert("Alpine", DistroWithColor{
		FileName: "Alpine",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("Alter", DistroWithColor{
		FileName: "Alter",
		Colors:   []int{6, 6},
	})

	r.Insert("Amazon", DistroWithColor{
		FileName: "Amazon",
		Colors:   []int{3, 7},
	})

	r.Insert("AmogOS", DistroWithColor{
		FileName: "AmogOS",
		Colors:   []int{15, 6},
	})

	r.Insert("Anarchy", DistroWithColor{
		FileName: "Anarchy",
		Colors:   []int{7, 4},
	})

	r.Insert("android_small", DistroWithColor{
		FileName: "android_small",
		Colors:   []int{2, 7},
	})

	r.Insert("Android", DistroWithColor{
		FileName: "Android",
		Colors:   []int{2, 7},
	})

	r.Insert("ArseLinux", DistroWithColor{
		FileName: "ArseLinux",
		Colors:   []int{4, 7},
	})

	r.Insert("instantOS", DistroWithColor{
		FileName: "instantOS",
		Colors:   []int{4, 6},
	})

	r.Insert("Antergos", DistroWithColor{
		FileName: "Antergos",
		Colors:   []int{4, 6},
	})

	r.Insert("antiX", DistroWithColor{
		FileName: "antiX",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("Afterglow", DistroWithColor{
		FileName: "Afterglow",
		Colors:   []int{5, 1, 3, 4},
	})

	r.Insert("AOSC OS", DistroWithColor{
		FileName: "AOSC",
		Colors:   []int{4, 0, 1, 3},
	})

	r.Insert("Apricity", DistroWithColor{
		FileName: "Apricity",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("archcraft_ascii", DistroWithColor{
		FileName: "archcraft_ascii",
		Colors:   []int{6, 1, 2, 3, 4, 5},
	})

	r.Insert("archcraft_minimal", DistroWithColor{
		FileName: "archcraft_minimal",
		Colors:   []int{6, 1, 2, 3, 4, 5},
	})

	r.Insert("Archcraft", DistroWithColor{
		FileName: "Archcraft",
		Colors:   []int{6, 1, 2, 3, 4, 5},
	})

	r.Insert("arcolinux_small", DistroWithColor{
		FileName: "arcolinux_small",
		Colors:   []int{7, 4},
	})

	r.Insert("ArcoLinux", DistroWithColor{
		FileName: "ArcoLinux",
		Colors:   []int{7, 4},
	})

	r.Insert("Arkane", DistroWithColor{
		FileName: "Arkane",
		Colors:   []int{7, 130, 237},
	})

	r.Insert("arch_small", DistroWithColor{
		FileName: "arch_small",
		Colors:   []int{6, 7, 1},
	})

	r.Insert("arch_old", DistroWithColor{
		FileName: "arch_old",
		Colors:   []int{6, 7, 1},
	})

	r.Insert("ArchBox", DistroWithColor{
		FileName: "ArchBox",
		Colors:   []int{2, 7, 1},
	})

	r.Insert("ARCHlabs", DistroWithColor{
		FileName: "ARCHlabs",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("ArchStrike", DistroWithColor{
		FileName: "ArchStrike",
		Colors:   []int{8, 6},
	})

	r.Insert("astOS", DistroWithColor{
		FileName: "astOS",
		Colors:   []int{8},
	})

	r.Insert("XFerience", DistroWithColor{
		FileName: "XFerience",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("Stock Linux", DistroWithColor{
		FileName: "Stock",
		Colors:   []int{4, 7},
	})

	r.Insert("ArchMerge", DistroWithColor{
		FileName: "ArchMerge",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("aerOS", DistroWithColor{
		FileName: "aerOS",
		// TODO: need to fix ,  "fg 0 0 0"
		Colors: []int{6, 6, 7, 1},
	})

	r.Insert("Arch", DistroWithColor{
		FileName: "Arch",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("unicodearch", DistroWithColor{
		FileName: "unicodearch",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("artix_small", DistroWithColor{
		FileName: "artix_small",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("Artix", DistroWithColor{
		FileName: "Artix",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("Cobalt", DistroWithColor{
		FileName: "Cobalt",
		Colors:   []int{33, 33, 59, 31, 8},
	})

	r.Insert("Arya", DistroWithColor{
		FileName: "Arya",
		Colors:   []int{2, 1},
	})

	r.Insert("AsteroidOS", DistroWithColor{
		FileName: "AsteroidOS",
		Colors:   []int{160, 208, 202, 214},
	})

	r.Insert("Athena", DistroWithColor{
		FileName: "Athena",
		Colors:   []int{7, 3},
	})

	r.Insert("azos", DistroWithColor{
		FileName: "azos",
		Colors:   []int{6, 1},
	})

	r.Insert("Bedrock", DistroWithColor{
		FileName: "Bedrock",
		Colors:   []int{8, 7},
	})

	r.Insert("Bazzite", DistroWithColor{
		FileName: "Bazzite",
		Colors:   []int{5, 5},
	})

	r.Insert("BigLinux_large", DistroWithColor{
		FileName: "BigLinux_large",
		Colors:   []int{2, 3, 4, 7},
	})

	r.Insert("BigLinux", DistroWithColor{
		FileName: "BigLinux",
		Colors:   []int{6, 11, 4},
	})

	r.Insert("Bitrig", DistroWithColor{
		FileName: "Bitrig",
		Colors:   []int{2, 7},
	})

	r.Insert("BlackArch", DistroWithColor{
		FileName: "BlackArch",
		Colors:   []int{1, 1, 0, 1},
	})

	r.Insert("BlackMesa", DistroWithColor{
		FileName: "BlackMesa",
		Colors:   []int{1},
	})

	r.Insert("blackPanther", DistroWithColor{
		FileName: "blackPanther",
		Colors:   []int{1, 11, 12},
	})
	r.Insert("blackpanther", DistroWithColor{
		FileName: "blackPanther",
		Colors:   []int{1, 11, 12},
	})

	r.Insert("MatuusOS", DistroWithColor{
		FileName: "MatuusOS",
		Colors:   []int{9, 11, 0},
	})

	r.Insert("BLAG", DistroWithColor{
		FileName: "BLAG",
		Colors:   []int{5, 7},
	})

	r.Insert("BlankOn", DistroWithColor{
		FileName: "BlankOn",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("BlueLight", DistroWithColor{
		FileName: "BlueLight",
		Colors:   []int{7, 4},
	})

	r.Insert("Bodhi", DistroWithColor{
		FileName: "Bodhi",
		Colors:   []int{7, 11, 2},
	})

	r.Insert("bonsai", DistroWithColor{
		FileName: "bonsai",
		Colors:   []int{6, 2, 3},
	})

	r.Insert("BSD", DistroWithColor{
		FileName: "BSD",
		Colors:   []int{1, 7, 4, 3, 6},
	})

	r.Insert("BunsenLabs", DistroWithColor{
		FileName: "BunsenLabs",
		// TODO: need to fix ,  "fg 7"
		Colors: []int{7},
	})

	r.Insert("CachyOS", DistroWithColor{
		FileName: "CachyOS",
		Colors:   []int{2, 8, 6},
	})

	r.Insert("Calculate", DistroWithColor{
		FileName: "Calculate",
		Colors:   []int{7, 3},
	})

	r.Insert("Carbs", DistroWithColor{
		FileName: "Carbs",
		Colors:   []int{4, 5, 4, 4, 4, 4},
	})

	r.Insert("CalinixOS", DistroWithColor{
		FileName: "CalinixOS",
		Colors:   []int{5, 4},
	})

	r.Insert("CalinixOS_small", DistroWithColor{
		FileName: "CalinixOS_small",
		Colors:   []int{5, 4},
	})

	r.Insert("CBL-Mariner", DistroWithColor{
		FileName: "CBL-Mariner",
		Colors:   []int{6},
	})

	r.Insert("CelOS", DistroWithColor{
		FileName: "CelOS",
		Colors:   []int{4, 6, 0, 5},
	})

	r.Insert("centos_small", DistroWithColor{
		FileName: "centos_small",
		Colors:   []int{3, 2, 4, 5, 7},
	})

	r.Insert("CentOS", DistroWithColor{
		FileName: "CentOS",
		Colors:   []int{3, 2, 4, 5, 7},
	})

	r.Insert("Center", DistroWithColor{
		FileName: "Center",
		Colors:   []int{7, 7},
	})

	r.Insert("Chakra", DistroWithColor{
		FileName: "Chakra",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("ChaletOS", DistroWithColor{
		FileName: "ChaletOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Chapeau", DistroWithColor{
		FileName: "Chapeau",
		Colors:   []int{2, 7},
	})
	r.Insert("Chimera", DistroWithColor{
		FileName: "Chimera",
		Colors:   []int{1, 5, 4, 1},
	})
	r.Insert("ChonkySealOS", DistroWithColor{
		FileName: "ChonkySealOS",
		Colors:   []int{7},
	})

	r.Insert("Chrom", DistroWithColor{
		FileName: "Chrom",
		Colors:   []int{2, 1, 3, 4, 7},
	})
	r.Insert("cleanjaro_small", DistroWithColor{
		FileName: "cleanjaro_small",
		Colors:   []int{7, 7},
	})
	r.Insert("Cleanjaro", DistroWithColor{
		FileName: "Cleanjaro",
		Colors:   []int{7, 7},
	})
	r.Insert("ClearOS", DistroWithColor{
		FileName: "ClearOS",
		Colors:   []int{2},
	})
	r.Insert("Clear Linux OS", DistroWithColor{
		FileName: "Clear_Linux",
		Colors:   []int{4, 3, 7, 6},
	})
	r.Insert("Clear_Linux", DistroWithColor{
		FileName: "Clear_Linux",
		Colors:   []int{4, 3, 7, 6},
	})
	r.Insert("Clover", DistroWithColor{
		FileName: "Clover",
		Colors:   []int{2, 6},
	})
	r.Insert("Condres", DistroWithColor{
		FileName: "Condres",
		Colors:   []int{2, 3, 6},
	})

	r.Insert("Container Linux by CoreOS", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Container_Linux", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Fedora CoreOS", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("CoreOS", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("crux_small", DistroWithColor{
		FileName: "crux_small",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("KISS", DistroWithColor{
		FileName: "crux_small",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("CRUX", DistroWithColor{
		FileName: "CRUX",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("Crystal Linux", DistroWithColor{
		FileName: "Crystal",
		Colors:   []int{13, 5},
	})

	r.Insert("Cucumber", DistroWithColor{
		FileName: "Cucumber",
		Colors:   []int{2, 3},
	})

	r.Insert("CutefishOS", DistroWithColor{
		FileName: "CutefishOS",
		Colors:   []int{6, 7, 4},
	})

	r.Insert("CuteOS", DistroWithColor{
		FileName: "CuteOS",
		Colors:   []int{33, 50, 57},
	})

	r.Insert("CyberOS", DistroWithColor{
		FileName: "CyberOS",
		Colors:   []int{50, 32, 57},
	})

	r.Insert("dahlia", DistroWithColor{
		FileName: "dahlia",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("debian_small", DistroWithColor{
		FileName: "debian_small",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("Debian", DistroWithColor{
		FileName: "Debian",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("digital UNIX", DistroWithColor{
		FileName: "digital",
		Colors:   []int{1, 6, 7},
	})

	r.Insert("BLAG", DistroWithColor{
		FileName: "BLAG",
		Colors:   []int{5, 7},
	})

	r.Insert("BlankOn", DistroWithColor{
		FileName: "BlankOn",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("BlueLight", DistroWithColor{
		FileName: "BlueLight",
		Colors:   []int{7, 4},
	})

	r.Insert("Bodhi", DistroWithColor{
		FileName: "Bodhi",
		Colors:   []int{7, 11, 2},
	})

	r.Insert("bonsai", DistroWithColor{
		FileName: "bonsai",
		Colors:   []int{6, 2, 3},
	})

	r.Insert("BSD", DistroWithColor{
		FileName: "BSD",
		Colors:   []int{1, 7, 4, 3, 6},
	})

	r.Insert("BunsenLabs", DistroWithColor{
		FileName: "BunsenLabs",
		// TODO: need to fix ,  "fg 7"
		Colors: []int{7},
	})

	r.Insert("CachyOS", DistroWithColor{
		FileName: "CachyOS",
		Colors:   []int{2, 8, 6},
	})

	r.Insert("Calculate", DistroWithColor{
		FileName: "Calculate",
		Colors:   []int{7, 3},
	})

	r.Insert("Carbs", DistroWithColor{
		FileName: "Carbs",
		Colors:   []int{4, 5, 4, 4, 4, 4},
	})

	r.Insert("CalinixOS", DistroWithColor{
		FileName: "CalinixOS",
		Colors:   []int{5, 4},
	})

	r.Insert("CalinixOS_small", DistroWithColor{
		FileName: "CalinixOS_small",
		Colors:   []int{5, 4},
	})

	r.Insert("CBL-Mariner", DistroWithColor{
		FileName: "CBL-Mariner",
		Colors:   []int{6},
	})

	r.Insert("CelOS", DistroWithColor{
		FileName: "CelOS",
		Colors:   []int{4, 6, 0, 5},
	})

	r.Insert("centos_small", DistroWithColor{
		FileName: "centos_small",
		Colors:   []int{3, 2, 4, 5, 7},
	})

	r.Insert("CentOS", DistroWithColor{
		FileName: "CentOS",
		Colors:   []int{3, 2, 4, 5, 7},
	})

	r.Insert("Center", DistroWithColor{
		FileName: "Center",
		Colors:   []int{7, 7},
	})

	r.Insert("Chakra", DistroWithColor{
		FileName: "Chakra",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("ChaletOS", DistroWithColor{
		FileName: "ChaletOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Chapeau", DistroWithColor{
		FileName: "Chapeau",
		Colors:   []int{2, 7},
	})
	r.Insert("Chimera", DistroWithColor{
		FileName: "Chimera",
		Colors:   []int{1, 5, 4, 1},
	})
	r.Insert("ChonkySealOS", DistroWithColor{
		FileName: "ChonkySealOS",
		Colors:   []int{7},
	})

	r.Insert("Chrom", DistroWithColor{
		FileName: "Chrom",
		Colors:   []int{2, 1, 3, 4, 7},
	})
	r.Insert("cleanjaro_small", DistroWithColor{
		FileName: "cleanjaro_small",
		Colors:   []int{7, 7},
	})
	r.Insert("Cleanjaro", DistroWithColor{
		FileName: "Cleanjaro",
		Colors:   []int{7, 7},
	})
	r.Insert("ClearOS", DistroWithColor{
		FileName: "ClearOS",
		Colors:   []int{2},
	})
	r.Insert("Clear Linux OS", DistroWithColor{
		FileName: "Clear_Linux",
		Colors:   []int{4, 3, 7, 6},
	})
	r.Insert("Clear_Linux", DistroWithColor{
		FileName: "Clear_Linux",
		Colors:   []int{4, 3, 7, 6},
	})
	r.Insert("Clover", DistroWithColor{
		FileName: "Clover",
		Colors:   []int{2, 6},
	})
	r.Insert("Condres", DistroWithColor{
		FileName: "Condres",
		Colors:   []int{2, 3, 6},
	})

	r.Insert("Container Linux by CoreOS", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Container_Linux", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Fedora CoreOS", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("CoreOS", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("crux_small", DistroWithColor{
		FileName: "crux_small",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("KISS", DistroWithColor{
		FileName: "crux_small",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("CRUX", DistroWithColor{
		FileName: "CRUX",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("Crystal Linux", DistroWithColor{
		FileName: "Crystal",
		Colors:   []int{13, 5},
	})

	r.Insert("Cucumber", DistroWithColor{
		FileName: "Cucumber",
		Colors:   []int{2, 3},
	})

	r.Insert("CutefishOS", DistroWithColor{
		FileName: "CutefishOS",
		Colors:   []int{6, 7, 4},
	})

	r.Insert("CuteOS", DistroWithColor{
		FileName: "CuteOS",
		Colors:   []int{33, 50, 57},
	})

	r.Insert("CyberOS", DistroWithColor{
		FileName: "CyberOS",
		Colors:   []int{50, 32, 57},
	})

	r.Insert("dahlia", DistroWithColor{
		FileName: "dahlia",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("debian_small", DistroWithColor{
		FileName: "debian_small",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("Debian", DistroWithColor{
		FileName: "Debian",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("digital UNIX", DistroWithColor{
		FileName: "digital",
		Colors:   []int{1, 6, 7},
	})
	r.Insert("BLAG", DistroWithColor{
		FileName: "BLAG",
		Colors:   []int{5, 7},
	})

	r.Insert("BlankOn", DistroWithColor{
		FileName: "BlankOn",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("BlueLight", DistroWithColor{
		FileName: "BlueLight",
		Colors:   []int{7, 4},
	})

	r.Insert("Bodhi", DistroWithColor{
		FileName: "Bodhi",
		Colors:   []int{7, 11, 2},
	})

	r.Insert("bonsai", DistroWithColor{
		FileName: "bonsai",
		Colors:   []int{6, 2, 3},
	})

	r.Insert("BSD", DistroWithColor{
		FileName: "BSD",
		Colors:   []int{1, 7, 4, 3, 6},
	})

	r.Insert("BunsenLabs", DistroWithColor{
		FileName: "BunsenLabs",
		// TODO: need to fix ,  "fg 7"
		Colors: []int{7},
	})

	r.Insert("CachyOS", DistroWithColor{
		FileName: "CachyOS",
		Colors:   []int{2, 8, 6},
	})

	r.Insert("Calculate", DistroWithColor{
		FileName: "Calculate",
		Colors:   []int{7, 3},
	})

	r.Insert("Carbs", DistroWithColor{
		FileName: "Carbs",
		Colors:   []int{4, 5, 4, 4, 4, 4},
	})

	r.Insert("CalinixOS", DistroWithColor{
		FileName: "CalinixOS",
		Colors:   []int{5, 4},
	})

	r.Insert("CalinixOS_small", DistroWithColor{
		FileName: "CalinixOS_small",
		Colors:   []int{5, 4},
	})

	r.Insert("CBL-Mariner", DistroWithColor{
		FileName: "CBL-Mariner",
		Colors:   []int{6},
	})

	r.Insert("CelOS", DistroWithColor{
		FileName: "CelOS",
		Colors:   []int{4, 6, 0, 5},
	})

	r.Insert("centos_small", DistroWithColor{
		FileName: "centos_small",
		Colors:   []int{3, 2, 4, 5, 7},
	})

	r.Insert("CentOS", DistroWithColor{
		FileName: "CentOS",
		Colors:   []int{3, 2, 4, 5, 7},
	})

	r.Insert("Center", DistroWithColor{
		FileName: "Center",
		Colors:   []int{7, 7},
	})

	r.Insert("Chakra", DistroWithColor{
		FileName: "Chakra",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("ChaletOS", DistroWithColor{
		FileName: "ChaletOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Chapeau", DistroWithColor{
		FileName: "Chapeau",
		Colors:   []int{2, 7},
	})
	r.Insert("Chimera", DistroWithColor{
		FileName: "Chimera",
		Colors:   []int{1, 5, 4, 1},
	})
	r.Insert("ChonkySealOS", DistroWithColor{
		FileName: "ChonkySealOS",
		Colors:   []int{7},
	})

	r.Insert("Chrom", DistroWithColor{
		FileName: "Chrom",
		Colors:   []int{2, 1, 3, 4, 7},
	})
	r.Insert("cleanjaro_small", DistroWithColor{
		FileName: "cleanjaro_small",
		Colors:   []int{7, 7},
	})
	r.Insert("Cleanjaro", DistroWithColor{
		FileName: "Cleanjaro",
		Colors:   []int{7, 7},
	})
	r.Insert("ClearOS", DistroWithColor{
		FileName: "ClearOS",
		Colors:   []int{2},
	})
	r.Insert("Clear Linux OS", DistroWithColor{
		FileName: "Clear_Linux",
		Colors:   []int{4, 3, 7, 6},
	})
	r.Insert("Clear_Linux", DistroWithColor{
		FileName: "Clear_Linux",
		Colors:   []int{4, 3, 7, 6},
	})
	r.Insert("Clover", DistroWithColor{
		FileName: "Clover",
		Colors:   []int{2, 6},
	})
	r.Insert("Condres", DistroWithColor{
		FileName: "Condres",
		Colors:   []int{2, 3, 6},
	})

	r.Insert("Container Linux by CoreOS", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Container_Linux", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Fedora CoreOS", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("CoreOS", DistroWithColor{
		FileName: "CoreOS",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("crux_small", DistroWithColor{
		FileName: "crux_small",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("KISS", DistroWithColor{
		FileName: "crux_small",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("CRUX", DistroWithColor{
		FileName: "CRUX",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("Crystal Linux", DistroWithColor{
		FileName: "Crystal",
		Colors:   []int{13, 5},
	})

	r.Insert("Cucumber", DistroWithColor{
		FileName: "Cucumber",
		Colors:   []int{2, 3},
	})

	r.Insert("CutefishOS", DistroWithColor{
		FileName: "CutefishOS",
		Colors:   []int{6, 7, 4},
	})

	r.Insert("CuteOS", DistroWithColor{
		FileName: "CuteOS",
		Colors:   []int{33, 50, 57},
	})

	r.Insert("CyberOS", DistroWithColor{
		FileName: "CyberOS",
		Colors:   []int{50, 32, 57},
	})

	r.Insert("dahlia", DistroWithColor{
		FileName: "dahlia",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("debian_small", DistroWithColor{
		FileName: "debian_small",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("Debian", DistroWithColor{
		FileName: "Debian",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("digital UNIX", DistroWithColor{
		FileName: "digital",
		Colors:   []int{1, 6, 7},
	})

	r.Insert("Droidian", DistroWithColor{
		FileName: "Droidian",
		Colors:   []int{2, 10},
	})

	r.Insert("Deepin", DistroWithColor{
		FileName: "Deepin",
		Colors:   []int{2, 7},
	})

	r.Insert("openKylin", DistroWithColor{
		FileName: "openKylin",
		Colors:   []int{2, 7},
	})

	r.Insert("DesaOS", DistroWithColor{
		FileName: "DesaOS",
		Colors:   []int{2, 7},
	})

	r.Insert("Devuan", DistroWithColor{
		FileName: "Devuan",
		Colors:   []int{5, 7},
	})

	r.Insert("DietPi", DistroWithColor{
		FileName: "DietPi",
		Colors:   []int{2, 0},
	})

	r.Insert("DracOS", DistroWithColor{
		FileName: "DracOS",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("DarkOs", DistroWithColor{
		FileName: "DarkOs",
		Colors:   []int{1, 6, 5, 3, 2},
	})

	r.Insert("Itc", DistroWithColor{
		FileName: "Itc",
		Colors:   []int{1},
	})

	r.Insert("dragonfly_old", DistroWithColor{
		FileName: "dragonfly_old",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("dragonfly_small", DistroWithColor{
		FileName: "dragonfly_small",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("DragonFly", DistroWithColor{
		FileName: "DragonFly",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("Drauger", DistroWithColor{
		FileName: "Drauger",
		Colors:   []int{1, 7},
	})

	r.Insert("elementary_small", DistroWithColor{
		FileName: "elementary_small",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Elementary", DistroWithColor{
		FileName: "Elementary",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Elive", DistroWithColor{
		FileName: "Elive",
		Colors:   []int{7, 6, 6},
	})

	r.Insert("endeavouros_small", DistroWithColor{
		FileName: "endeavouros_small",
		Colors:   []int{1, 5, 4},
	})

	r.Insert("EndeavourOS", DistroWithColor{
		FileName: "EndeavourOS",
		Colors:   []int{1, 5, 4},
	})

	r.Insert("EncryptOS", DistroWithColor{
		FileName: "EncryptOS",
		Colors:   []int{2, 5, 6},
	})

	r.Insert("Endless", DistroWithColor{
		FileName: "Endless",
		Colors:   []int{1, 7},
	})
	r.Insert("Enso", DistroWithColor{
		FileName: "Enso",
		Colors:   []int{8, 7},
	})

	r.Insert("EuroLinux", DistroWithColor{
		FileName: "EuroLinux",
		Colors:   []int{4, 7},
	})

	r.Insert("EvolutionOS", DistroWithColor{
		FileName: "EvolutionOS",
		Colors:   []int{4, 7},
	})
	r.Insert("eweOS", DistroWithColor{
		FileName: "eweOS",
		Colors:   []int{7, 11, 9, 8, 1},
	})

	r.Insert("Exherbo", DistroWithColor{
		FileName: "Exherbo",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("fedora_small", DistroWithColor{
		FileName: "fedora_small",
		Colors:   []int{12},
	})
	r.Insert("Fedora_old", DistroWithColor{
		FileName: "Fedora_old",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("RFRemix", DistroWithColor{
		FileName: "Fedora_old",
		Colors:   []int{4, 3, 7, 6},
	})
	r.Insert("Fedora_unicode", DistroWithColor{
		FileName: "Fedora_unicode",
		Colors:   []int{12, 7},
	})

	r.Insert("Fedora Silverblue", DistroWithColor{
		FileName: "Silverblue",
		Colors:   []int{4, 7, 12},
	})

	r.Insert("Silverblue", DistroWithColor{
		FileName: "Silverblue",
		Colors:   []int{4, 7, 12},
	})
	r.Insert("Fedora Kinoite", DistroWithColor{
		FileName: "Kinoite",
		Colors:   []int{12, 7},
	})
	r.Insert("Kinoite", DistroWithColor{
		FileName: "Kinoite",
		Colors:   []int{12, 7},
	})

	r.Insert("Fedora Sericea", DistroWithColor{
		FileName: "Sericea",
		Colors:   []int{12, 7},
	})
	r.Insert("Sericea", DistroWithColor{
		FileName: "Sericea",
		Colors:   []int{12, 7},
	})

	r.Insert("Fedora", DistroWithColor{
		FileName: "Fedora",
		Colors:   []int{12, 7},
	})

	r.Insert("Feren", DistroWithColor{
		FileName: "Feren",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("FemboyOS", DistroWithColor{
		FileName: "FemboyOS",
		Colors:   []int{4},
	})

	r.Insert("Finnix", DistroWithColor{
		FileName: "Finnix",
		Colors:   []int{4, 7, 7},
	})

	r.Insert("Furreto", DistroWithColor{
		FileName: "Furreto",
		Colors:   []int{211, 255, 225, 199},
	})

	r.Insert("freebsd_small", DistroWithColor{
		FileName: "freebsd_small",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("FreeBSD", DistroWithColor{
		FileName: "FreeBSD",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("HardenedBSD", DistroWithColor{
		FileName: "FreeBSD",
		Colors:   []int{4, 7, 3},
	})
	r.Insert("FreeMiNT", DistroWithColor{
		FileName: "FreeMiNT",
		Colors:   []int{7},
	})
	r.Insert("Frugalware", DistroWithColor{
		FileName: "Frugalware",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Funtoo", DistroWithColor{
		FileName: "Funtoo",
		Colors:   []int{5, 7},
	})

	r.Insert("GalliumOS", DistroWithColor{
		FileName: "GalliumOS",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("garuda_small", DistroWithColor{
		FileName: "garuda_small",
		Colors:   []int{7, 7, 3, 7, 2, 4},
	})

	r.Insert("Garuda", DistroWithColor{
		FileName: "Garuda",
		Colors:   []int{7, 7, 3, 7, 2, 4},
	})
	r.Insert("gentoo_small", DistroWithColor{
		FileName: "gentoo_small",
		Colors:   []int{5, 7},
	})

	r.Insert("Gentoo", DistroWithColor{
		FileName: "Gentoo",
		Colors:   []int{5, 7},
	})

	r.Insert("Pentoo", DistroWithColor{
		FileName: "Pentoo",
		Colors:   []int{5, 7},
	})

	r.Insert("glaucus", DistroWithColor{
		FileName: "glaucus",
		Colors:   []int{5},
	})

	r.Insert("gNewSense", DistroWithColor{
		FileName: "gNewSense",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("GNOME", DistroWithColor{
		FileName: "GNOME",
		Colors:   []int{4},
	})

	r.Insert("GNU", DistroWithColor{
		FileName: "GNU",
		// TODO: need to fix, "fg 7"
		Colors: []int{7},
	})

	r.Insert("GoboLinux", DistroWithColor{
		FileName: "GoboLinux",
		Colors:   []int{5, 4, 6, 2},
	})

	r.Insert("GrapheneOS", DistroWithColor{
		FileName: "GrapheneOS",
		Colors:   []int{7, 4},
	})

	r.Insert("Grombyang", DistroWithColor{
		FileName: "Grombyang",
		Colors:   []int{4, 2, 1},
	})

	r.Insert("guix_small", DistroWithColor{
		FileName: "guix_small",
		Colors:   []int{3, 7, 6, 1, 8},
	})
	r.Insert("Guix", DistroWithColor{
		FileName: "Guix",
		Colors:   []int{3, 7, 6, 1, 8},
	})

	r.Insert("haiku_small", DistroWithColor{
		FileName: "haiku_small",
		Colors:   []int{2, 8},
	})

	r.Insert("Haiku", DistroWithColor{
		FileName: "Haiku",
		Colors:   []int{1, 3, 7, 2},
	})

	r.Insert("HamoniKR", DistroWithColor{
		FileName: "HamoniKR",
		Colors:   []int{4, 7, 99, 1, 3, 7},
	})

	r.Insert("Huayra", DistroWithColor{
		FileName: "Huayra",
		Colors:   []int{4, 7},
	})

	r.Insert("HydroOS", DistroWithColor{
		FileName: "HydroOS",
		Colors:   []int{1, 2, 3, 4, 5},
	})

	r.Insert("hyperbola_small", DistroWithColor{
		FileName: "hyperbola_small",
		Colors:   []int{8},
	})

	r.Insert("Hyperbola", DistroWithColor{
		FileName: "Hyperbola",
		Colors:   []int{8},
	})

	r.Insert("Hybrid", DistroWithColor{
		FileName: "Hybrid",
		Colors:   []int{4, 12},
	})

	r.Insert("iglunix", DistroWithColor{
		FileName: "iglunix",
		Colors:   []int{8},
	})

	r.Insert("iglu", DistroWithColor{
		FileName: "iglunix",
		Colors:   []int{8},
	})

	r.Insert("Interix", DistroWithColor{
		FileName: "Interix",
		Colors:   []int{1, 7, 4, 0, 3},
	})
	r.Insert("januslinux", DistroWithColor{
		FileName: "januslinux",
		Colors:   []int{4, 5, 6, 2},
	})
	r.Insert("janus", DistroWithColor{
		FileName: "januslinux",
		Colors:   []int{4, 5, 6, 2},
	})
	r.Insert("Ataraxia Linux", DistroWithColor{
		FileName: "januslinux",
		Colors:   []int{4, 5, 6, 2},
	})
	r.Insert("Ataraxia", DistroWithColor{
		FileName: "januslinux",
		Colors:   []int{4, 5, 6, 2},
	})
	r.Insert("Kaisen", DistroWithColor{
		FileName: "Kaisen",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("kali_small", DistroWithColor{
		FileName: "kali_small",
		Colors:   []int{4, 8},
	})
	r.Insert("kalilinux_small", DistroWithColor{
		FileName: "kali_small",
		Colors:   []int{4, 8},
	})
	r.Insert("kali_linux_small", DistroWithColor{
		FileName: "kali_small",
		Colors:   []int{4, 8},
	})

	r.Insert("Kali", DistroWithColor{
		FileName: "Kali",
		Colors:   []int{4, 8},
	})

	r.Insert("KaOS", DistroWithColor{
		FileName: "KaOS",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("KDE", DistroWithColor{
		FileName: "KDE",
		Colors:   []int{2, 7},
	})

	r.Insert("Kibojoe", DistroWithColor{
		FileName: "Kibojoe",
		Colors:   []int{2, 7, 4},
	})

	r.Insert("Kogaion", DistroWithColor{
		FileName: "Kogaion",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Korora", DistroWithColor{
		FileName: "Korora",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("KSLinux", DistroWithColor{
		FileName: "KSLinux",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Kubuntu", DistroWithColor{
		FileName: "Kubuntu",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("LainOS", DistroWithColor{
		FileName: "LainOS",
		Colors:   []int{4, 14, 7, 2, 3, 5},
	})

	r.Insert("LEDE", DistroWithColor{
		FileName: "LEDE",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("LangitKetujuh_old", DistroWithColor{
		FileName: "LangitKetujuh_old",
		Colors:   []int{7, 4},
	})

	r.Insert("LangitKetujuh", DistroWithColor{
		FileName: "LangitKetujuh",
		Colors:   []int{4, 2},
	})

	r.Insert("LaxerOS", DistroWithColor{
		FileName: "LaxerOS",
		Colors:   []int{7, 4},
	})

	r.Insert("LibreELEC", DistroWithColor{
		FileName: "LibreELEC",
		Colors:   []int{2, 3, 7, 14, 13},
	})

	r.Insert("Linux", DistroWithColor{
		FileName: "Linux",
		Colors:   []int{8, 3},
	})

	r.Insert("LinuxFromScratch", DistroWithColor{
		FileName: "LinuxFromScratch",
		Colors:   []int{8, 3},
	})
	r.Insert("LFS", DistroWithColor{
		FileName: "LinuxFromScratch",
		Colors:   []int{8, 3},
	})
	r.Insert("linux-from-scratch", DistroWithColor{
		FileName: "LinuxFromScratch",
		Colors:   []int{8, 3},
	})
	r.Insert("linux_from_scratch", DistroWithColor{
		FileName: "LinuxFromScratch",
		Colors:   []int{8, 3},
	})

	r.Insert("Linux_small", DistroWithColor{
		FileName: "Linux_small",
		Colors:   []int{8, 3},
	})

	r.Insert("linuxlite_small", DistroWithColor{
		FileName: "linuxlite_small",
		Colors:   []int{3, 7},
	})

	r.Insert("Linux Lite", DistroWithColor{
		FileName: "Linux_Lite",
		// TODO: need to fix ,  "fg 0 0 0"
		Colors: []int{3, 7},
	})

	r.Insert("Linux_Lite", DistroWithColor{
		FileName: "Linux_Lite",
		Colors:   []int{3, 7},
	})

	r.Insert("LMDE", DistroWithColor{
		FileName: "LMDE",
		Colors:   []int{2, 7},
	})

	r.Insert("Lubuntu", DistroWithColor{
		FileName: "Lubuntu",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Lunar", DistroWithColor{
		FileName: "Lunar",
		Colors:   []int{4, 7, 3},
	})

	r.Insert("MacaroniOS", DistroWithColor{
		FileName: "MacaroniOS",
		Colors:   []int{2, 3, 15, 14},
	})

	r.Insert("Macaroni OS", DistroWithColor{
		FileName: "MacaroniOS",
		Colors:   []int{2, 3, 15, 14},
	})

	r.Insert("Arya", DistroWithColor{
		FileName: "Arya",
		Colors:   []int{2, 1},
	})

	// FIXME: "mac"*"_small")

	r.Insert("mac", DistroWithColor{
		FileName: "mac",
		Colors:   []int{2, 3, 1, 1, 5, 4},
	})

	r.Insert("Darwin", DistroWithColor{
		FileName: "mac",
		Colors:   []int{2, 3, 1, 1, 5, 4},
	})

	r.Insert("mageia_small", DistroWithColor{
		FileName: "mageia_small",
		Colors:   []int{6, 7},
	})

	r.Insert("Mageia", DistroWithColor{
		FileName: "Mageia",
		Colors:   []int{6, 7},
	})

	r.Insert("Magix", DistroWithColor{
		FileName: "Magix",
		Colors:   []int{206, 39},
	})

	r.Insert("MagpieOS", DistroWithColor{
		FileName: "MagpieOS",
		Colors:   []int{2, 1, 3, 5},
	})

	r.Insert("MainsailOS_small", DistroWithColor{
		FileName: "MainsailOS_small",
		Colors:   []int{1},
	})

	r.Insert("MainsailOS", DistroWithColor{
		FileName: "MainsailOS",
		Colors:   []int{1},
	})

	r.Insert("Mandriva", DistroWithColor{
		FileName: "Mandriva",
		Colors:   []int{4, 3},
	})

	r.Insert("manjaro_small", DistroWithColor{
		FileName: "manjaro_small",
		Colors:   []int{2, 7},
	})

	r.Insert("Manjaro", DistroWithColor{
		FileName: "Manjaro",
		Colors:   []int{2, 7},
	})

	r.Insert("MassOS", DistroWithColor{
		FileName: "MassOS",
		Colors:   []int{7},
	})
	r.Insert("TeArch", DistroWithColor{
		FileName: "TeArch",
		Colors:   []int{39, 7, 1},
	})

	r.Insert("Maui", DistroWithColor{
		FileName: "Maui",
		Colors:   []int{6, 7},
	})

	r.Insert("Mauna", DistroWithColor{
		FileName: "Mauna",
		Colors:   []int{4, 6},
	})

	r.Insert("Meowix", DistroWithColor{
		FileName: "Meowix",
		Colors:   []int{1, 3, 3, 4},
	})

	r.Insert("Mer", DistroWithColor{
		FileName: "Mer",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Minix", DistroWithColor{
		FileName: "Minix",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("MIRACLE LINUX", DistroWithColor{
		FileName: "MIRACLE_LINUX",
		Colors:   []int{29},
	})

	r.Insert("MIRACLE_LINUX", DistroWithColor{
		FileName: "MIRACLE_LINUX",
		Colors:   []int{29},
	})

	r.Insert("Linspire", DistroWithColor{
		FileName: "Linspire",
		Colors:   []int{4, 2},
	})
	r.Insert("Freespire", DistroWithColor{
		FileName: "Linspire",
		Colors:   []int{4, 2},
	})
	r.Insert("Lindows", DistroWithColor{
		FileName: "Linspire",
		Colors:   []int{4, 2},
	})

	r.Insert("linuxmint_small", DistroWithColor{
		FileName: "linuxmint_small",
		Colors:   []int{2, 7},
	})

	r.Insert("Linux Mint Old", DistroWithColor{
		FileName: "mint_old",
		Colors:   []int{2, 7},
	})
	r.Insert("LinuxMintOld", DistroWithColor{
		FileName: "mint_old",
		Colors:   []int{2, 7},
	})
	r.Insert("mint_old", DistroWithColor{
		FileName: "mint_old",
		Colors:   []int{2, 7},
	})

	r.Insert("Linux Mint", DistroWithColor{
		FileName: "mint",
		Colors:   []int{2, 7},
	})
	r.Insert("LinuxMint", DistroWithColor{
		FileName: "mint",
		Colors:   []int{2, 7},
	})
	r.Insert("mint", DistroWithColor{
		FileName: "mint",
		Colors:   []int{2, 7},
	})

	r.Insert("Live Raizo", DistroWithColor{
		FileName: "Live_Raizo",
		Colors:   []int{3},
	})
	r.Insert("Live_Raizo", DistroWithColor{
		FileName: "Live_Raizo",
		Colors:   []int{3},
	})

	r.Insert("mx_small", DistroWithColor{
		FileName: "mx_small",
		Colors:   []int{4, 6, 7},
	})

	r.Insert("MX", DistroWithColor{
		FileName: "MX",
		Colors:   []int{4, 6, 7},
	})

	r.Insert("Namib", DistroWithColor{
		FileName: "Namib",
		Colors:   []int{1},
	})

	r.Insert("NekOS", DistroWithColor{
		FileName: "NekOS",
		Colors:   []int{3, 7, 1},
	})

	r.Insert("Neptune", DistroWithColor{
		FileName: "Neptune",
		Colors:   []int{7},
	})

	r.Insert("netbsd_small", DistroWithColor{
		FileName: "netbsd_small",
		Colors:   []int{5, 7},
	})

	r.Insert("NetBSD", DistroWithColor{
		FileName: "NetBSD",
		Colors:   []int{5, 7},
	})

	r.Insert("Netrunner", DistroWithColor{
		FileName: "Netrunner",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Nitrux", DistroWithColor{
		FileName: "Nitrux",
		Colors:   []int{4},
	})
	r.Insert("nixos_small", DistroWithColor{
		FileName: "nixos_small",
		Colors:   []int{4, 6},
	})
	r.Insert("nixos_old", DistroWithColor{
		FileName: "nixos_old",
		Colors:   []int{4, 6},
	})

	r.Insert("nixos_colorful", DistroWithColor{
		FileName: "nixos_colorful",
		Colors:   []int{4, 6, 4, 6, 4, 6},
	})
	r.Insert("NixOS", DistroWithColor{
		FileName: "NixOS",
		Colors:   []int{4, 6},
	})
	r.Insert("Nobara", DistroWithColor{
		FileName: "Nobara",
		// FIXME: "#999999" "#d80a0a" "#e5b90b"
		Colors: []int{7, 7},
	})
	r.Insert("NomadBSD", DistroWithColor{
		FileName: "NomadBSD",
		Colors:   []int{4},
	})
	r.Insert("GhostBSD", DistroWithColor{
		FileName: "GhostBSD",
		Colors:   []int{4},
	})
	r.Insert("Nurunner", DistroWithColor{
		FileName: "Nurunner",
		Colors:   []int{4},
	})
	r.Insert("NuTyX", DistroWithColor{
		FileName: "NuTyX",
		Colors:   []int{4, 1},
	})
	r.Insert("OBRevenge", DistroWithColor{
		FileName: "OBRevenge",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("OmniOS", DistroWithColor{
		FileName: "OmniOS",
		Colors:   []int{15, 3, 8},
	})
	r.Insert("openbsd_small", DistroWithColor{
		FileName: "openbsd_small",
		Colors:   []int{3, 7, 6, 1, 8},
	})
	r.Insert("OpenBSD", DistroWithColor{
		FileName: "OpenBSD",
		Colors:   []int{3, 7, 6, 1, 8},
	})
	r.Insert("openEuler", DistroWithColor{
		FileName: "openEuler",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("OpenIndiana", DistroWithColor{
		FileName: "OpenIndiana",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("openmamba", DistroWithColor{
		FileName: "openmamba",
		Colors:   []int{7, 2},
	})

	r.Insert("OpenMandriva", DistroWithColor{
		FileName: "OpenMandriva",
		Colors:   []int{4},
	})

	r.Insert("OpenStage", DistroWithColor{
		FileName: "OpenStage",
		Colors:   []int{2},
	})

	// FIXME:*Wrt*
	r.Insert("Wrt", DistroWithColor{
		FileName: "Wrt",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Open Source Media Center", DistroWithColor{
		FileName: "osmc",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("osmc", DistroWithColor{
		FileName: "osmc",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("OPNsense", DistroWithColor{
		FileName: "OPNsense",
		Colors:   []int{7, 202},
	})

	r.Insert("Oracle", DistroWithColor{
		FileName: "Oracle",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("orchid_small", DistroWithColor{
		FileName: "orchid_small",
		Colors:   []int{255, 127},
	})
	r.Insert("orchid", DistroWithColor{
		FileName: "orchid",
		Colors:   []int{255, 127, 127},
	})
	r.Insert("OS Elbrus", DistroWithColor{
		FileName: "Elbrus",
		Colors:   []int{4, 7, 3},
	})
	r.Insert("PacBSD", DistroWithColor{
		FileName: "PacBSD",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("Panwah", DistroWithColor{
		FileName: "Panwah",
		Colors:   []int{7, 1, 0},
	})

	r.Insert("parabola_small", DistroWithColor{
		FileName: "parabola_small",
		Colors:   []int{5, 7},
	})

	r.Insert("Parabola", DistroWithColor{
		FileName: "Parabola",
		Colors:   []int{5, 7},
	})

	r.Insert("Pardus", DistroWithColor{
		FileName: "Pardus",
		Colors:   []int{3, 7, 6, 1, 8},
	})

	r.Insert("Parrot", DistroWithColor{
		FileName: "Parrot",
		Colors:   []int{6, 7},
	})

	r.Insert("Parsix", DistroWithColor{
		FileName: "Parsix",
		Colors:   []int{3, 1, 7, 8},
	})

	r.Insert("PCBSD", DistroWithColor{
		FileName: "PCBSD",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("TrueOS", DistroWithColor{
		FileName: "PCBSD",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("PCLinuxOS", DistroWithColor{
		FileName: "PCLinuxOS",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("pearOS", DistroWithColor{
		FileName: "pearOS",
		Colors:   []int{2, 3, 1, 1, 5, 4},
	})

	r.Insert("Pengwin", DistroWithColor{
		FileName: "Pengwin",
		Colors:   []int{5, 5, 13},
	})

	r.Insert("Peppermint", DistroWithColor{
		FileName: "Peppermint",
		Colors:   []int{1, 15, 3},
	})

	r.Insert("Pisi", DistroWithColor{
		FileName: "Pisi",
		Colors:   []int{12, 7, 6, 1, 8},
	})

	r.Insert("PikaOS", DistroWithColor{
		FileName: "PikaOS",
		Colors:   []int{11},
	})

	r.Insert("phyOS", DistroWithColor{
		FileName: "phyOS",
		Colors:   []int{33, 33, 7, 1},
	})

	r.Insert("PNM Linux", DistroWithColor{
		FileName: "WHPNM",
		Colors:   []int{33, 9, 15, 202},
	})
	r.Insert("WHPNM Linux", DistroWithColor{
		FileName: "WHPNM",
		Colors:   []int{33, 9, 15, 202},
	})

	r.Insert("popos_small", DistroWithColor{
		FileName: "popos_small",
		Colors:   []int{6, 7},
	})
	r.Insert("pop_os_small", DistroWithColor{
		FileName: "popos_small",
		Colors:   []int{6, 7},
	})

	r.Insert("Pop!_OS", DistroWithColor{
		FileName: "pop_os",
		Colors:   []int{6, 7},
	})
	r.Insert("popos", DistroWithColor{
		FileName: "pop_os",
		Colors:   []int{6, 7},
	})
	r.Insert("pop_os", DistroWithColor{
		FileName: "pop_os",
		Colors:   []int{6, 7},
	})

	r.Insert("Porteus", DistroWithColor{
		FileName: "Porteus",
		Colors:   []int{6, 7},
	})

	r.Insert("postmarketos_small", DistroWithColor{
		FileName: "postmarketos_small",
		Colors:   []int{2, 7},
	})
	r.Insert("PostMarketOS", DistroWithColor{
		FileName: "PostMarketOS",
		Colors:   []int{2, 7},
	})
	r.Insert("PuffOS", DistroWithColor{
		FileName: "PuffOS",
		Colors:   []int{3},
	})
	r.Insert("Proxmox", DistroWithColor{
		FileName: "Proxmox",
		Colors:   []int{7, 202},
	})

	r.Insert("Puppy", DistroWithColor{
		FileName: "Puppy",
		Colors:   []int{4, 7},
	})
	r.Insert("Quirky Werewolf", DistroWithColor{
		FileName: "Puppy",
		Colors:   []int{4, 7},
	})
	r.Insert("Precise Puppy", DistroWithColor{
		FileName: "Puppy",
		Colors:   []int{4, 7},
	})
	r.Insert("pureos_small", DistroWithColor{
		FileName: "pureos_small",
		Colors:   []int{2, 7, 7},
	})
	r.Insert("PureOS", DistroWithColor{
		FileName: "PureOS",
		Colors:   []int{2, 7, 7},
	})
	r.Insert("Peropesis", DistroWithColor{
		FileName: "Peropesis",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Q4OS", DistroWithColor{
		FileName: "Q4OS",
		Colors:   []int{4, 1},
	})
	r.Insert("Qubes", DistroWithColor{
		FileName: "Qubes",
		Colors:   []int{4, 5, 7, 6},
	})
	r.Insert("Qubyt", DistroWithColor{
		FileName: "Qubyt",
		Colors:   []int{4, 5, 0, 4},
	})
	r.Insert("Quibian", DistroWithColor{
		FileName: "Quibian",
		Colors:   []int{3, 7},
	})

	r.Insert("Radix", DistroWithColor{
		FileName: "Radix",
		Colors:   []int{1, 2},
	})
	r.Insert("Raspbian_small", DistroWithColor{
		FileName: "Raspbian_small",
		Colors:   []int{2, 1},
	})
	r.Insert("Raspbian", DistroWithColor{
		FileName: "Raspbian",
		Colors:   []int{2, 1},
	})
	r.Insert("ravynOS", DistroWithColor{
		FileName: "ravynOS",
		Colors:   []int{15},
	})

	r.Insert("Reborn OS", DistroWithColor{
		FileName: "Reborn",
		Colors:   []int{0, 4, 6},
	})
	r.Insert("Reborn", DistroWithColor{
		FileName: "Reborn",
		Colors:   []int{0, 4, 6},
	})

	r.Insert("Red Star", DistroWithColor{
		FileName: "Redstar",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("Redstar", DistroWithColor{
		FileName: "Redstar",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("Redcore", DistroWithColor{
		FileName: "Redcore",
		Colors:   []int{1},
	})

	r.Insert("redhat_old", DistroWithColor{
		FileName: "redhat_old",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("rhel_old", DistroWithColor{
		FileName: "redhat_old",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("Redhat", DistroWithColor{
		FileName: "Redhat",
		Colors:   []int{1},
	})

	r.Insert("Red Hat", DistroWithColor{
		FileName: "Redhat",
		Colors:   []int{1},
	})

	r.Insert("rhel", DistroWithColor{
		FileName: "Redhat",
		Colors:   []int{1},
	})
	r.Insert("Refracted Devuan", DistroWithColor{
		FileName: "Refracted_Devuan",
		Colors:   []int{8, 7},
	})
	r.Insert("Refracted_Devuan", DistroWithColor{
		FileName: "Refracted_Devuan",
		Colors:   []int{8, 7},
	})
	r.Insert("Regata", DistroWithColor{
		FileName: "Regata",
		Colors:   []int{7, 1, 4, 5, 3, 2},
	})
	r.Insert("Regolith", DistroWithColor{
		FileName: "Regolith",
		Colors:   []int{1},
	})

	r.Insert("RhaymOS", DistroWithColor{
		FileName: "RhaymOS",
		Colors:   []int{1},
	})

	r.Insert("Rhino Linux", DistroWithColor{
		FileName: "Rhino",
		Colors:   []int{5, 4},
	})

	r.Insert("rocky_small", DistroWithColor{
		FileName: "rocky_small",
		Colors:   []int{2},
	})

	r.Insert("rocky", DistroWithColor{
		FileName: "rocky",
		Colors:   []int{35},
	})

	r.Insert("Rosa", DistroWithColor{
		FileName: "Rosa",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Astra Linux", DistroWithColor{
		FileName: "Astra",
		// TODO: need to fix
		//  if [[ $ascii_distro == *"ALCE"* ]]; then
		// set_colors 160 231
		// else
		// set_colors 32 231
		// fi
		Colors: []int{160, 231},
	})

	r.Insert("sabotage", DistroWithColor{
		FileName: "sabotage",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Sabayon", DistroWithColor{
		FileName: "Sabayon",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Sailfish", DistroWithColor{
		FileName: "Sailfish",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("SalentOS", DistroWithColor{
		FileName: "SalentOS",
		Colors:   []int{2, 1, 3, 7},
	})

	r.Insert("ShastraOS", DistroWithColor{
		FileName: "ShastraOS",
		Colors:   []int{6},
	})

	r.Insert("Sasanqua", DistroWithColor{
		FileName: "Sasanqua",
		Colors:   []int{5, 1, 5},
	})

	r.Insert("Salient OS", DistroWithColor{
		FileName: "salientos",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("SalientOS", DistroWithColor{
		FileName: "salientos",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("salientos", DistroWithColor{
		FileName: "salientos",
		Colors:   []int{6, 6, 7, 1},
	})

	r.Insert("Salix", DistroWithColor{
		FileName: "Salix",
		Colors:   []int{2},
	})

	r.Insert("Scientific", DistroWithColor{
		FileName: "Scientific",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Septor", DistroWithColor{
		FileName: "Septor",
		Colors:   []int{4, 7, 4},
	})
	r.Insert("Serene", DistroWithColor{
		FileName: "Serene",
		Colors:   []int{6, 6},
	})
	r.Insert("SharkLinux", DistroWithColor{
		FileName: "SharkLinux",
		Colors:   []int{4, 7},
	})
	r.Insert("Siduction", DistroWithColor{
		FileName: "Siduction",
		Colors:   []int{4, 4},
	})

	r.Insert("Slackel", DistroWithColor{
		FileName: "Slackel",
		Colors:   []int{3},
	})
	r.Insert("slackware_small", DistroWithColor{
		FileName: "slackware_small",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Slackware", DistroWithColor{
		FileName: "Slackware",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("SliTaz", DistroWithColor{
		FileName: "SliTaz",
		Colors:   []int{3, 3},
	})
	r.Insert("SmartOS", DistroWithColor{
		FileName: "SmartOS",
		Colors:   []int{6, 7},
	})
	r.Insert("SkiffOS", DistroWithColor{
		FileName: "SkiffOS",
		Colors:   []int{4, 7},
	})
	r.Insert("Solus", DistroWithColor{
		FileName: "Solus",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Sulin", DistroWithColor{
		FileName: "Sulin",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Source Mage", DistroWithColor{
		FileName: "Source_Mage",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Source_Mage", DistroWithColor{
		FileName: "Source_Mage",
		Colors:   []int{4, 7, 1},
	})
	r.Insert("Sparky", DistroWithColor{
		FileName: "Sparky",
		Colors:   []int{1, 7},
	})
	r.Insert("Star", DistroWithColor{
		FileName: "Star",
		Colors:   []int{7},
	})

	r.Insert("SteamOS", DistroWithColor{
		FileName: "SteamOS",
		Colors:   []int{5, 7},
	})

	r.Insert("sunos_small", DistroWithColor{
		FileName: "sunos_small",
		Colors:   []int{3, 7},
	})

	r.Insert("solaris_small", DistroWithColor{
		FileName: "sunos_small",
		Colors:   []int{3, 7},
	})

	r.Insert("SunOS", DistroWithColor{
		FileName: "SunOS",
		Colors:   []int{3, 7},
	})

	r.Insert("Solaris", DistroWithColor{
		FileName: "SunOS",
		Colors:   []int{2, 3},
	})

	r.Insert("openSUSE Leap", DistroWithColor{
		FileName: "openSUSE_Leap",
		Colors:   []int{2, 7},
	})

	r.Insert("openSUSE_Leap", DistroWithColor{
		FileName: "openSUSE_Leap",
		Colors:   []int{2, 7},
	})

	r.Insert("t2", DistroWithColor{
		FileName: "t2",
		Colors:   []int{7, 4},
	})

	r.Insert("openSUSE Tumbleweed-Slowroll", DistroWithColor{
		FileName: "openSUSE_Slowroll",
		Colors:   []int{2, 7},
	})
	r.Insert("openSUSE_Slowroll", DistroWithColor{
		FileName: "openSUSE_Slowroll",
		Colors:   []int{2, 7},
	})
	r.Insert("openSUSE Tumbleweed", DistroWithColor{
		FileName: "openSUSE_Tumbleweed",
		Colors:   []int{2, 7},
	})
	r.Insert("openSUSE_Tumbleweed", DistroWithColor{
		FileName: "openSUSE_Tumbleweed",
		Colors:   []int{2, 7},
	})

	r.Insert("opensuse_small", DistroWithColor{
		FileName: "opensuse_small",
		Colors:   []int{2, 7},
	})

	r.Insert("suse_small", DistroWithColor{
		FileName: "opensuse_small",
		Colors:   []int{2, 7},
	})

	r.Insert("openSUSE", DistroWithColor{
		FileName: "SUSE",
		Colors:   []int{2, 7},
	})

	r.Insert("open SUSE", DistroWithColor{
		FileName: "SUSE",
		Colors:   []int{2, 7},
	})

	r.Insert("SUSE", DistroWithColor{
		FileName: "SUSE",
		Colors:   []int{2, 7},
	})

	r.Insert("parch", DistroWithColor{
		FileName: "parch",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Parch", DistroWithColor{
		FileName: "parch",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("PARCH", DistroWithColor{
		FileName: "parch",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("SwagArch", DistroWithColor{
		FileName: "SwagArch",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Tails", DistroWithColor{
		FileName: "Tails",
		Colors:   []int{5, 7},
	})

	r.Insert("TorizonCore", DistroWithColor{
		FileName: "TorizonCore",
		Colors:   []int{3, 4, 8, 5},
	})

	r.Insert("Trisquel", DistroWithColor{
		FileName: "Trisquel",
		Colors:   []int{4, 6},
	})

	r.Insert("Twister", DistroWithColor{
		FileName: "Twister",
		Colors:   []int{2, 1, 4, 5, 7},
	})

	r.Insert("Ubuntu Cinnamon", DistroWithColor{
		FileName: "Ubuntu-Cinnamon",
		Colors:   []int{1, 7},
	})

	r.Insert("Ubuntu-Cinnamon", DistroWithColor{
		FileName: "Ubuntu-Cinnamon",
		Colors:   []int{1, 7},
	})

	r.Insert("Ubuntu Budgie", DistroWithColor{
		FileName: "Ubuntu-Budgie",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Ubuntu-Budgie", DistroWithColor{
		FileName: "Ubuntu-Budgie",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Ubuntu-GNOME", DistroWithColor{
		FileName: "Ubuntu-GNOME",
		Colors:   []int{4, 5, 7, 6},
	})

	r.Insert("Ubuntu Kylin", DistroWithColor{
		FileName: "Ubuntu-Kylin",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("Ubuntu-Kylin", DistroWithColor{
		FileName: "Ubuntu-Kylin",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("Ubuntu Touch", DistroWithColor{
		FileName: "Ubuntu-Touch",
		Colors:   []int{3, 7},
	})

	r.Insert("Ubuntu MATE", DistroWithColor{
		FileName: "Ubuntu-MATE",
		Colors:   []int{2, 7},
	})

	r.Insert("Ubuntu-MATE", DistroWithColor{
		FileName: "Ubuntu-MATE",
		Colors:   []int{2, 7},
	})
	r.Insert("ubuntu_old02", DistroWithColor{
		FileName: "ubuntu_old02",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("Ubuntu Studio", DistroWithColor{
		FileName: "Ubuntu-Studio",
		Colors:   []int{6, 7},
	})

	r.Insert("Ubuntu-Studio", DistroWithColor{
		FileName: "Ubuntu-Studio",
		Colors:   []int{6, 7},
	})
	r.Insert("Ubuntu Sway", DistroWithColor{
		FileName: "Ubuntu-Sway",
		Colors:   []int{6, 7},
	})
	r.Insert("Ubuntu-Sway", DistroWithColor{
		FileName: "Ubuntu-Sway",
		Colors:   []int{6, 7},
	})
	r.Insert("ubuntu_small", DistroWithColor{
		FileName: "ubuntu_small",
		Colors:   []int{1},
	})

	r.Insert("Ubuntu_old", DistroWithColor{
		FileName: "Ubuntu_old",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("i3buntu", DistroWithColor{
		FileName: "Ubuntu_old",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("Floflis", DistroWithColor{
		FileName: "Floflis",
		Colors:   []int{1, 7, 3},
	})
	r.Insert("Ubuntu", DistroWithColor{
		FileName: "Ubuntu",
		Colors:   []int{1},
	})

	r.Insert("Ultramarine Linux", DistroWithColor{
		FileName: "ultramarine",
		Colors:   []int{4, 7},
	})
	r.Insert("ultramarine", DistroWithColor{
		FileName: "ultramarine",
		Colors:   []int{4, 7},
	})

	r.Insert("Univalent", DistroWithColor{
		FileName: "Univalent",
		Colors:   []int{6, 6},
	})

	r.Insert("Uos", DistroWithColor{
		FileName: "Uos",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("Univention", DistroWithColor{
		FileName: "Univention",
		Colors:   []int{1, 7},
	})

	r.Insert("uwuntu", DistroWithColor{
		FileName: "uwuntu",
		Colors:   []int{225, 206, 52},
	})

	r.Insert("UrukOS", DistroWithColor{
		FileName: "UrukOS",
		Colors:   []int{12, 12, 7, 12, 4},
	})

	r.Insert("venom_small", DistroWithColor{
		FileName: "venom_small",
		Colors:   []int{8, 4},
	})

	r.Insert("Venom", DistroWithColor{
		FileName: "Venom",
		Colors:   []int{8, 4},
	})

	r.Insert("void_small", DistroWithColor{
		FileName: "void_small",
		Colors:   []int{2, 8},
	})
	r.Insert("void_old", DistroWithColor{
		FileName: "void_old",
		Colors:   []int{2, 8},
	})
	r.Insert("Void", DistroWithColor{
		FileName: "Void",
		Colors:   []int{8, 2, 7},
	})
	r.Insert("VNux", DistroWithColor{
		FileName: "VNux",
		Colors:   []int{11, 8, 15, 1, 7},
	})

	r.Insert("VzLinux", DistroWithColor{
		FileName: "VzLinux",
		Colors:   []int{1, 7, 3},
	})

	r.Insert("yiffOS", DistroWithColor{
		FileName: "yiffOS",
		Colors:   []int{93, 92},
	})

	r.Insert("semc", DistroWithColor{
		FileName: "semc",
		Colors:   []int{2, 8, 1},
	})
	r.Insert("Vanilla", DistroWithColor{
		FileName: "Vanilla",
		Colors:   []int{11, 15},
	})

	r.Insert("Obarun", DistroWithColor{
		FileName: "Obarun",
		Colors:   []int{6, 6, 7, 1},
	})

	// FIXME:  *"wii-linux-ngx"*|*"whiite-linux"*|\
	// *"gc-linux"*)

	// *"[Windows 11]"*|*"on Windows 11"*|\
	// "Windows 11"* |"windows11")
	//
	//   *"[Windows 10]"*|*"on Windows 10"*|"Windows 8"*|\
	// "Windows 10"* |"windows10"|"windows8")
	r.Insert("Windows95", DistroWithColor{
		FileName: "Windows95",
		Colors:   []int{6, 4, 3, 2, 1, 0},
	})

	r.Insert("Windows", DistroWithColor{
		FileName: "Windows",
		Colors:   []int{1, 2, 4, 3},
	})

	r.Insert("Xray_OS", DistroWithColor{
		FileName: "Xray_OS",
		Colors:   []int{15, 14, 16, 24},
	})

	r.Insert("Xubuntu", DistroWithColor{
		FileName: "Xubuntu",
		Colors:   []int{4, 7, 1},
	})

	r.Insert("Soda", DistroWithColor{
		FileName: "Soda",
		Colors:   []int{1, 7},
	})

	r.Insert("secureblue", DistroWithColor{
		FileName: "secureblue",
		Colors:   []int{12, 7},
	})

	r.Insert("KrassOS", DistroWithColor{
		FileName: "KrassOS",
		Colors:   []int{4, 7},
	})

	r.Insert("KrassOS", DistroWithColor{
		FileName: "Krass",
		Colors:   []int{4, 7},
	})

	r.Insert("IRIX", DistroWithColor{
		FileName: "IRIX",
		Colors:   []int{4, 7},
	})
	r.Insert("Xenia2", DistroWithColor{
		FileName: "Xenia2",
		// FIXME: '#55CDFD' '#F6AAB7' '#FFFFFF'
		Colors: []int{3, 7, 6, 1, 8},
	})

	r.Insert("Xenia", DistroWithColor{
		FileName: "Xenia",
		// FIXME: '#55CDFD' '#F6AAB7' '#FFFFFF'
		Colors: []int{2, 8},
	})

	r.Insert("Tatra", DistroWithColor{
		FileName: "Tatra",
		Colors:   []int{4, 7},
	})

	r.Insert("Zorin", DistroWithColor{
		FileName: "Zorin",
		Colors:   []int{4, 6},
	})

	r.Insert("BSD", DistroWithColor{
		FileName: "BSD",
		Colors:   []int{1, 7, 4, 3, 6},
	})

	r.Insert("Darwin", DistroWithColor{
		FileName: "Darwin",
		Colors:   []int{2, 3, 1, 1, 5, 4},
	})

	r.Insert("GNU", DistroWithColor{
		FileName: "GNU",
		// FIXME: fg 7
		Colors: []int{7},
	})

	r.Insert("Linux", DistroWithColor{
		FileName: "Linux",
		// FIXME: fg 8,3
		Colors: []int{8, 3},
	})

	r.Insert("Profelis SambaBOX", DistroWithColor{
		FileName: "SambaBOX",
		Colors:   []int{3, 6},
	})

	r.Insert("SambaBOX", DistroWithColor{
		FileName: "SambaBOX",
		Colors:   []int{3, 6},
	})

	r.Insert("SunOS", DistroWithColor{
		FileName: "SunOS",
		Colors:   []int{3, 7},
	})

}

type DistroWithColor struct {
	FileName string
	Colors   []int
}
