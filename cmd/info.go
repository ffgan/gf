package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ffgan/gf/configs"
	cli "github.com/ffgan/gf/internal/CLI"
	gui "github.com/ffgan/gf/internal/GUI"
	dev "github.com/ffgan/gf/internal/devices"
	"github.com/ffgan/gf/internal/logo"
	"github.com/ffgan/gf/internal/utils"
)

type hardware struct {
	leftMax   int
	cache_dir string

	KernelName     string
	kernel_version string
	KernelMachine  string
	DarwinName     string
	os             string
	distro         string
	ascii_distro   string

	bios string

	model string

	user     string
	hostname string
	title    string
	length   int

	kernel string

	uptime string

	packages string

	shell string

	editor string

	de string

	wm string

	wm_theme string

	cpu string
	gpu string
	npu string

	memory string
	swap   string

	network string

	bluetooth string

	disk          string
	battery       string
	local_ip      string
	public_ip     string
	users         string
	locale        string
	gpu_driver    string
	song          string
	java_ver      string
	python_ver    string
	node_ver      string
	term_size     string
	window_size   string
	power_adapter string

	plain    bool
	verbose  bool
	disabled map[string]bool
	extraMu  sync.Mutex
	extra    map[string]string

	resolution string

	term      string
	term_font string
	font      string

	ascii_bold string
	bold       string

	underline string
	theme     string
	icons     string
	cursor    string
	cols      string
	configs.Config
}

func (h *hardware) Get_os() {
	h.os = cli.GetOS(h.KernelName, h.DarwinName)
}

func (h *hardware) Get_distro() {
	h.distro, h.ascii_distro = cli.GetDistro(h.os, h.OSArch, h.KernelMachine, h.DistroShorthand, h.ASCIIDistro)
}

func (h *hardware) Get_bios() {
	h.bios = dev.GetBIOS(h.os)
}

func (h *hardware) Get_model() {
	h.model = cli.GetModel(h.os, h.KernelName, h.KernelMachine)
}
func (h *hardware) Get_title() {
	title, err := cli.GetTitle()
	if err != nil || title == nil {
		return
	}
	h.user = title.User
	h.hostname = title.Hostname
	h.length = title.Length
	if h.plainMode() {
		h.title = title.User + "@" + title.Hostname
		return
	}
	bold := ""
	if h.Bold != utils.OFF {
		bold = "\033[1m"
	}
	h.title = logo.TitleColor() + bold + title.User + logo.AtColor() + "@" + logo.TitleColor() + bold + title.Hostname + logo.ResetCode()
}

func (h *hardware) Get_kernel() {
	h.kernel = cli.GetKernel(h.OSArch, h.DistroShorthand, h.KernelShorthand, h.ascii_distro, h.os, h.kernel_version, h.KernelName, h.KernelMachine)
}

func (h *hardware) Get_uptime() {
	h.uptime = cli.GetUptime(h.os, h.UptimeShorthand)
}

func (h *hardware) Get_packages() {
	h.packages = cli.DetectPackages(h.os, cli.PackageConfig{
		Managers: h.PackageManagers,
		Separate: h.PackageSeparate != utils.OFF,
		Minimal:  h.PackageMinimal == utils.ON,
	})
}

func (h *hardware) Get_shell() {
	h.shell = cli.GetShell(h.ShellPath, h.ShellVersion)
}

func (h *hardware) Get_editor() {
	h.editor = cli.GetEditor(h.EditorPath, h.EditorVersion)
}

func (h *hardware) Get_de() {
	h.de = gui.GetDE(h.os, h.distro, h.KernelName, h.DeVersion)
}

func (h *hardware) Get_wm() {
	h.wm = gui.GetWM(h.os, h.KernelName)
}

func (h *hardware) Get_wm_theme() {
	if h.wm == "" {
		h.Get_wm()
	}
	h.wm_theme = gui.GetWMTheme(h.wm)
}

func (h *hardware) Get_cpu() {
	temp := h.CPUTemp
	if temp == utils.ON {
		temp = "C"
	}
	h.cpu = cli.GetCPU(h.os, h.KernelMachine, h.CPUCores, h.CPUSpeed, temp, h.CPUBrand, h.SpeedType, h.SpeedShorthand, temp)
}

func (h *hardware) Get_gpu() {
	h.gpu = dev.GetGPU(h.os, h.GPUBrand, h.GPUType)
}

func (h *hardware) Get_npu() {
	h.npu = dev.GetAccel()
}

func (h *hardware) Get_memory() {
	h.memory = cli.PrintMem(h.os, h.Memory, h.ProgressBar)
}

func (h *hardware) Get_swap() {
	h.swap = cli.PrintSwap(h.os, h.Memory.MemoryPercent == utils.ON)
}

func (h *hardware) Get_network() {
	h.network = dev.GetNetwork(h.os)
}

func (h *hardware) Get_bluetooth() {
	h.bluetooth = dev.GetBluetooth()
}

func (h *hardware) Get_song() {
	h.song = dev.GetSong(h.SongFormat, h.MusicPlayer)
}

func (h *hardware) Get_resolution() {
	h.resolution = dev.GetResolution(h.os, h.KernelMachine, h.RefreshRate == utils.ON)
}

func (h *hardware) Get_style() {
}

func (h *hardware) Get_theme() {
	h.theme = gui.GetTheme()
}

func (h *hardware) Get_icons() {
	h.icons = gui.GetIcons()
}

func (h *hardware) Get_font() {
	h.font = gui.GetFont()
}

func (h *hardware) Get_cursor() {
	h.cursor = gui.GetCursor()
}

func (h *hardware) Get_java_ver() {
	h.java_ver = cli.ToolVersion([]string{"java"}, "-version")
}

func (h *hardware) Get_python_ver() {
	h.python_ver = cli.ToolVersion([]string{"python3", "python"}, "--version")
}

func (h *hardware) Get_node_ver() {
	h.node_ver = cli.ToolVersion([]string{"node"}, "--version")
}

func (h *hardware) Get_term() {
	h.term = cli.GetTerm(h.os)
}

func (h *hardware) Get_term_font() {
	if h.term == "" {
		h.Get_term()
	}
	h.term_font = gui.GetTermFont(h.term)
}

func (h *hardware) Get_disk() {
	width, _ := strconv.Atoi(h.BarLength)
	h.disk = cli.GetDisk(cli.DiskReport{
		Show:      h.DiskShow,
		Subtitle:  h.DiskSubtitle,
		Percent:   h.DiskPercent,
		Display:   h.DiskDisplay,
		BarWidth:  width,
		BarFilled: h.BarCharTotal,
		BarEmpty:  h.BarCharElapsed,
	})
}

func (h *hardware) Get_power_adapter() {
	h.power_adapter = dev.GetPowerAdapter()
}

func (h *hardware) Get_battery() {
	h.battery = dev.GetBattery()
}

func (h *hardware) Get_local_ip() {
	h.local_ip = dev.LocalIP(h.LocalIPInterface)
}

func (h *hardware) Get_public_ip() {
	h.public_ip = dev.PublicIP(h.PublicIPHost, h.PublicIPTimeout)
}

func (h *hardware) Get_users() {
	h.users = cli.GetUsers()
}

func (h *hardware) Get_locale() {
	h.locale = cli.GetLocale()
}

func (h *hardware) Get_gpu_driver() {
	h.gpu_driver = dev.GetGPUDriver()
}

func (h *hardware) Get_cols() {
	if h.plainMode() || h.Enabled == utils.OFF {
		h.cols = ""
		return
	}
	start, end := 0, 15
	fmt.Sscanf(strings.NewReplacer("(", " ", ")", " ", ",", " ").Replace(h.BlockRange), "%d %d", &start, &end)
	width, _ := strconv.Atoi(strings.TrimSpace(h.BlockWidth))
	height, _ := strconv.Atoi(strings.TrimSpace(h.BlockHeight))
	h.cols = strings.Join(cli.RenderBlocks(width, height, start, end), "\n")
}

func (h *hardware) Get_image_source() {

}

func (h *hardware) Get_wallpaper() {

}

func (h *hardware) Get_w3m_img_path() {

}

func (h *hardware) Get_window_size() {
	if os.Getenv("DISPLAY") == "" || !utils.CommandExists("xdotool") {
		h.window_size = ""
		return
	}
	out := utils.RunCommand("xdotool", "getactivewindow", "getwindowgeometry")
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Geometry:") {
			h.window_size = strings.TrimSpace(strings.TrimPrefix(line, "Geometry:"))
			return
		}
	}
}

func (h *hardware) Get_term_size() {
	h.term_size = cli.GetTermSize()
}

func (h *hardware) Get_image_size() {

}

func (h *hardware) Get_underline() {
	if h.UnderlineEnabled == utils.OFF {
		h.underline = ""
		return
	}
	if h.length == 0 {
		h.Get_title()
	}
	n := h.length
	if n <= 0 {
		n = 8
	}
	ch := h.UnderlineChar
	if ch == "" {
		ch = "-"
	}
	line := strings.Repeat(ch, n)
	if h.plainMode() {
		h.underline = line
		return
	}
	h.underline = logo.UnderlineColor() + line + logo.ResetCode()
}

func (h *hardware) Get_bold() {}

func (h *hardware) Get_full_path() {

}

func (h *hardware) Get_user_config() {

}
