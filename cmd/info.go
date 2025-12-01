package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/ffgan/gf/configs"
	cli "github.com/ffgan/gf/internal/CLI"
	gui "github.com/ffgan/gf/internal/GUI"
	dev "github.com/ffgan/gf/internal/devices"
	"github.com/ffgan/gf/internal/logo"
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

	memory string

	network string

	bluetooth string

	resolution string

	term      string
	term_font string

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
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	h.user = title.User
	h.hostname = title.Hostname
	h.title = title.Title
	h.length = title.Length
}

func (h *hardware) Get_kernel() {
	h.kernel = cli.GetKernel(h.OSArch, h.DistroShorthand, h.KernelShorthand, h.ascii_distro, h.os, h.kernel_version, h.KernelName, h.KernelMachine)
}

func (h *hardware) Get_uptime() {
	h.uptime = cli.GetUptime(h.os, h.UptimeShorthand)
}

func (h *hardware) Get_packages() {
	h.packages = cli.DetectPackages(h.os)
}

func (h *hardware) Get_shell() {
	h.shell = cli.GetShell(h.ShellPath, h.ShellVersion)
}

func (h *hardware) Get_editor() {
	h.editor = cli.GetEditor(h.EditorPath, h.EditorVersion)
}

func (h *hardware) Get_de() {
	h.de = gui.GetDE(h.os, h.distro)
}

func (h *hardware) Get_wm() {
	h.wm = gui.GetWM(h.os, h.KernelName)
}

func (h *hardware) Get_wm_theme() {
	h.wm_theme = gui.GetTheme()
}

func (h *hardware) Get_cpu() {
	if h.CPUTemp == "off" {
		h.cpu = cli.GetCPU(h.os, h.KernelMachine, h.CPUCores, h.CPUSpeed, h.CPUTemp, h.CPUBrand, h.SpeedType, h.SpeedShorthand, h.CPUTemp)
	}
	if h.CPUTemp == "on" {
		h.CPUTemp = "C"
	}
	h.cpu = cli.GetCPU(h.os, h.KernelMachine, h.CPUCores, h.CPUSpeed, "on", h.CPUBrand, h.SpeedType, h.SpeedShorthand, h.CPUTemp)
}

func (h *hardware) Get_gpu() {
	h.gpu = dev.GetGPU(h.os)
}

func (h *hardware) Get_memory() {
	h.memory = cli.PrintMem(h.os, h.Memory, h.ProgressBar)
}

func (h *hardware) Get_network() {
	h.network = dev.GetNetwork(h.os)
}

func (h *hardware) Get_bluetooth() {
	h.bluetooth = dev.GetBluetooth()
}

func (h *hardware) Get_song() {

}

func (h *hardware) Get_resolution() {
	h.resolution = dev.GetResolution(h.os, h.KernelMachine)
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
	h.term_font = gui.GetFont()
}

func (h *hardware) Get_cursor() {
	h.cursor = gui.GetCursor()
}

func (h *hardware) Get_java_ver() {

}

func (h *hardware) Get_python_ver() {

}

func (h *hardware) Get_node_ver() {

}

func (h *hardware) Get_term() {
	h.term = cli.GetTerm(h.os)
}

func (h *hardware) Get_term_font() {
	h.term_font = gui.GetFont()
}

func (h *hardware) Get_disk() {

}

func (h *hardware) Get_power_adapter() {

}

func (h *hardware) Get_battery() {

}

func (h *hardware) Get_local_ip() {

}

func (h *hardware) Get_public_ip() {

}

func (h *hardware) Get_users() {

}

func (h *hardware) Get_locale() {

}

func (h *hardware) Get_gpu_driver() {

}

func (h *hardware) Get_cols() {
	h.cols = cli.Getcols(h.leftMax)
}

func (h *hardware) Get_image_source() {

}

func (h *hardware) Get_wallpaper() {

}

func (h *hardware) Get_w3m_img_path() {

}

func (h *hardware) Get_window_size() {

}

func (h *hardware) Get_term_size() {

}

func (h *hardware) Get_image_size() {

}

func (h *hardware) Get_underline() {
	switch h.UnderlineEnabled {
	case "on":
		h.underline = strings.Repeat("-", 10)
	case "off":
		h.underline = ""
	}
}

func (h *hardware) Get_bold() {
	switch h.ascii_bold {
	case "on":
		h.ascii_bold = "\\e[1m"
	case "off":
		h.ascii_bold = ""
	}

	switch h.bold {
	case "on":
		h.bold = "\\e[1m"
	case "off":
		h.bold = ""
	}

}

func (h *hardware) Get_full_path() {

}

func (h *hardware) Get_user_config() {

}

func get_cache_dir() string {
	cacheDir := os.Getenv("TMPDIR")
	if cacheDir == "" {
		cacheDir = "/tmp"
	}
	return cacheDir
}

func (h *hardware) Get_ppid() {

}

func (h *hardware) Get_process_name() {

}

func (h *hardware) Get_args() {

}

// func info(subtitle, function_name string) {

// }

// func prin() {

// }

// func print_info(asciiData, lines []string) {
// 	info.PrintInfo(asciiData, lines)
// }

func image_backend() {

}
func old_functions() {

}

func dynamic_prompt() {

}

func Get_distro_ascii(ascii_distro, image_source string, ASCIIFiles embed.FS) []string {
	return logo.PrintASCII(ascii_distro, image_source, ASCIIFiles)
}

func cache_uname() (string, string, string, string) {
	uname_arr := strings.Split(cli.UName("-srm"), " ")
	kernel_name := uname_arr[0]
	kernel_version := uname_arr[1]
	kernel_machine := uname_arr[2]

	// TODO: https://github.com/ffgan/hyfetch/blob/master/neofetch#L6384
	var darwin_name string
	return kernel_name, kernel_version, kernel_machine, darwin_name
}

// Global state variables (mimicking bash script state)
var (
	subtitle     string
	prin         bool
	info_height  int
	json         bool
	reset        string = "\x1b[0m"
	colon_color  string = ""
	separator    string = ":"
	info_color   string = ""
	bold         string = "\x1b[1m"
	text_padding string = ""
	zws          string = "\u200B" // zero-width space
)

func NewContext(target hardware) *Context {
	return &Context{
		target: target,
		values: make(map[string]string),
	}
}

type Context struct {
	target hardware
	values map[string]string
}

func (c *Context) info(arg1, arg2 string) {
	if arg2 != "" {
		subtitle = arg1
	}

	prin = false

	var funcName string
	if arg2 != "" {
		funcName = arg2
	} else {
		funcName = arg1
	}

	methodName := "Get_" + funcName

	c.callMethod(methodName)

	// If the Get_func function called 'prin' directly, stop here
	if prin {
		return
	}

	output := c.getField(arg2)
	fmt.Printf("\033[%dC", c.target.leftMax)
	// Print based on conditions
	if arg2 != "" && strings.TrimSpace(output) != "" {
		// c.prin(arg1, output)
		fmt.Println(Info(arg1, output))
	} else if strings.TrimSpace(output) != "" {
		c.prin(output, "")
	} else {
		log.Printf("Info: Couldn't detect %s.", arg1)
	}

	// Unset subtitle
	subtitle = ""
}

func Info(name, data string) string {
	if data == "" {
		return ""
	}
	if name != "" {
		name = "\033[1;34m" + name + ": \033[0m"
	}
	return name + data
}

// prin prints formatted output
func (c *Context) prin(arg1, arg2 string) {
	var str string
	var subtitle_color string

	// If $2 doesn't exist we format $1 as info
	if trim(arg1) != "" && arg2 != "" {
		if json {
			fmt.Printf("    %s\n", Info(arg1, arg2))
			return
		}
		if arg2 != "" {
			str = arg1 + ": " + arg2
		} else {
			str = arg1
		}
	} else {
		if arg2 != "" {
			str = arg2
		} else {
			str = arg1
		}
		subtitle_color = info_color
	}

	// Remove all reset sequences from string
	str = strings.ReplaceAll(trim(str), "\x1b[0m", "")

	// Calculate length without escape sequences
	length := stripSequences(str)
	_ = len(length) // length variable for compatibility

	// Format the output
	str = strings.Replace(str, ":", reset+colon_color+separator+info_color, 1)
	str = subtitle_color + bold + str

	// Print the info
	var padding string
	if text_padding != "" {
		padding = fmt.Sprintf("\x1b[%sC", text_padding)
	}

	// Replace literal \n (not newlines) in output
	str = strings.ReplaceAll(str, "\\n", "")

	fmt.Printf("%s%s%s%s \n", padding, zws, str, reset)

	// Calculate info height
	info_height++

	// Log that prin was used
	prin = true
}

// trim removes leading and trailing whitespace
func trim(s string) string {
	return strings.TrimSpace(s)
}

// stripSequences removes ANSI escape sequences
func stripSequences(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return re.ReplaceAllString(s, "")
}

func (c *Context) getField(name string) string {
	v := reflect.ValueOf(c.target)
	field := v.FieldByName(name)
	return field.String()
}

func (c *Context) callMethod(methodName string) {
	v := reflect.ValueOf(&c.target)

	method := v.MethodByName(methodName)
	if !method.IsValid() {
		fmt.Printf("Error: Method %s not found\n", methodName)
		return
	}

	method.Call(nil)
}
