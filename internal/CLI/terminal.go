package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

var (
	termRun bool
	term    string
)

// GetTerm 负责检测当前终端类型
func GetTerm(osname string) string {
	if termRun {
		return term
	}
	defer func() { termRun = true }()

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

	// 通过父进程检测
	if term == "" {
		parent := os.Getppid()
		for term == "" && parent > 1 {
			name := utils.GetProcessName(string(parent))
			switch {
			case strings.Contains(name, "gnome-terminal-"):
				term = "gnome-terminal"
			case name == "kgx":
				term = "gnome-console"
			case strings.Contains(name, "urxvtd"):
				term = "urxvt"
			case strings.Contains(name, "nvim"):
				term = "Neovim Terminal"
			case strings.Contains(name, "NeoVimServer"):
				term = "VimR Terminal"
			case name == "tmux" || strings.Contains(name, "systemd") || strings.Contains(name, "sshd"):
				return term
			default:
				if osname == "Linux" {
					if real, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", parent)); err == nil {
						term = filepath.Base(real)
					}
				}
				if term == "" {
					term = filepath.Base(name)
				}
				// 修复 nix wrapper
				if strings.HasPrefix(term, ".") && strings.HasSuffix(term, "-wrapped") {
					term = strings.TrimPrefix(term, ".")
					term = strings.TrimSuffix(term, "-wrapped")
				}
			}
			parent = utils.GetParentPID(parent)
		}
	}

	if os.Getenv("FIG_TERM") == "1" {
		term += " + Fig"
	}
	if os.Getenv("TERMUX_VERSION") != "" && (term == "" || term == "com.termux") {
		term = "Termux " + os.Getenv("TERMUX_VERSION")
	}

	return term
}
