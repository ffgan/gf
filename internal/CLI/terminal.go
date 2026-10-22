package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ffgan/gf/internal/utils"
)

var (
	termOnce sync.Once
	term     string
)

// GetTerm 负责检测当前终端类型
func GetTerm(osname string) string {
	termOnce.Do(func() {
		detectTerm(osname)
	})
	return term
}

func detectTerm(osname string) {

	termProgram := os.Getenv("TERM_PROGRAM")
	termEnv := os.Getenv("TERM")

	switch termProgram {
	case "iTerm.app":
		term = "iTerm2"
	case "Terminal.app", "Apple_Terminal":
		term = "Apple Terminal"
	case "Hyper":
		term = "HyperTerm"
	default:
		if termProgram != "" {
			term = strings.TrimSuffix(termProgram, ".app")
		}
	}

	// Quick checks
	if termEnv == "tw52" || termEnv == "tw100" {
		term = "TosWin2"
	}
	if os.Getenv("SSH_CONNECTION") != "" {
		term = os.Getenv("SSH_TTY")
	}
	if os.Getenv("WT_SESSION") != "" {
		term = "Windows Terminal"
	}

	// Walk parent processes, skipping shells, until a known terminal shows up.
	if term == "" {
		parent := os.Getppid()
		for parent > 1 {
			name := processBase(parent, osname)
			switch {
			case strings.Contains(name, "gnome-terminal"):
				term = "gnome-terminal"
			case name == "kgx":
				term = "gnome-console"
			case name == "urxvtd", name == "urxvt":
				term = "urxvt"
			case strings.Contains(name, "nvim"):
				term = "Neovim Terminal"
			case strings.Contains(name, "NeoVimServer"):
				term = "VimR Terminal"
			case knownTerminal(name):
				term = name
			}
			next := utils.GetParentPID(parent)
			if term != "" || next == parent || next <= 1 || isStopParent(name) {
				break
			}
			if name == "" || isShell(name) || isIgnoredParent(name) {
				parent = next
				continue
			}
			parent = next
		}
	}

	if term == "" && termEnv != "" && termEnv != "dumb" && termEnv != "linux" {
		term = termEnv
	}

	if os.Getenv("FIG_TERM") == "1" {
		term += " + Fig"
	}
	if os.Getenv("TERMUX_VERSION") != "" && (term == "" || term == "com.termux") {
		term = "Termux " + os.Getenv("TERMUX_VERSION")
	}
}

func processBase(pid int, osname string) string {
	name := utils.GetProcessName(strconv.Itoa(pid))
	if osname == Linux {
		if real, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid)); err == nil {
			base := filepath.Base(real)
			if base != "" && base != "." {
				name = base
			}
		}
	}
	name = filepath.Base(name)
	if strings.HasPrefix(name, ".") && strings.HasSuffix(name, "-wrapped") {
		name = strings.TrimPrefix(name, ".")
		name = strings.TrimSuffix(name, "-wrapped")
	}
	return name
}

func isShell(name string) bool {
	switch name {
	case "sh", "bash", "zsh", "fish", "dash", "ash", "ksh", "mksh", "tcsh", "nu", "elvish", "xonsh", "oil", "osh", "yash", "pwsh", "csh":
		return true
	default:
		return false
	}
}

func isStopParent(name string) bool {
	return name == "tmux" || strings.Contains(name, "sshd")
}

func isIgnoredParent(name string) bool {
	switch name {
	case "sudo", "su", "doas", "login", "init", "systemd", "node", "python", "python3", "go", "ruby", "perl", "java":
		return true
	default:
		return strings.Contains(name, "systemd") || strings.Contains(name, "containerd")
	}
}

func knownTerminal(name string) bool {
	known := []string{
		"xfce4-terminal", "konsole", "xterm", "alacritty", "kitty", "wezterm",
		"tilix", "terminator", "terminology", "st", "foot", "ghostty", "contour",
		"hyper", "tabby", "warp", "code", "cursor", "mate-terminal", "lxterminal",
		"qterminal", "yakuake", "termite", "sakura", "tilda", "guake", "ptyxis",
		"rio", "blackbox-terminal", "gnome-terminal", "gnome-console", "urxvt",
		"rxvt", "cool-retro-term", "x-terminal-emul",
	}
	for _, item := range known {
		if name == item {
			return true
		}
	}
	return false
}
