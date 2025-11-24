package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
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
			name := getProcessName(parent)
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
			parent = getParentPID(parent)
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

// GetTermFont 根据 term 读取配置文件提取字体
func GetTermFont(osname string) string {
	if !termRun {
		GetTerm(osname)
	}

	switch {
	case strings.HasPrefix(term, "alacritty"):
		conf := findFirstFile([]string{
			filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "alacritty/alacritty.yml"),
			filepath.Join(os.Getenv("HOME"), ".config/alacritty/alacritty.yml"),
		})
		if conf != "" {
			data, _ := os.ReadFile(conf)
			re := regexp.MustCompile(`family:\s*"?([^"\n]+)"?`)
			if m := re.FindStringSubmatch(string(data)); len(m) > 1 {
				termFont = strings.TrimSpace(m[1])
			}
		}
	case strings.HasPrefix(term, "Apple Terminal"):
		out, _ := exec.Command("osascript", "-e", `tell application "Terminal" to font name of window frontmost & " " & font size of window frontmost`).Output()
		termFont = strings.TrimSpace(string(out))
	case strings.HasPrefix(term, "kitty"):
		cmd := exec.Command("kitty", "+runpy", "from kitty.cli import *;o=create_default_opts();print(f'{o.font_family} {o.font_size}')")
		if out, err := cmd.Output(); err == nil {
			termFont = strings.TrimSpace(string(out))
		}
	case strings.HasPrefix(term, "Hyper"):
		conf := filepath.Join(os.Getenv("HOME"), ".hyper.js")
		data, _ := os.ReadFile(conf)
		re := regexp.MustCompile(`fontFamily:\s*['"]([^'"]+)['"]`)
		if m := re.FindStringSubmatch(string(data)); len(m) > 1 {
			termFont = strings.TrimSpace(m[1])
		}
	case strings.HasPrefix(term, "xfce4-terminal"):
		conf := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "xfce4/terminal/terminalrc")
		data, _ := os.ReadFile(conf)
		re := regexp.MustCompile(`FontName=(.*)`)
		if m := re.FindStringSubmatch(string(data)); len(m) > 1 {
			termFont = strings.TrimSpace(m[1])
		} else {
			termFont = "Monospace 12"
		}
	default:
		termFont = ""
	}
	return termFont
}

func getParentPID(pid int) int {
	out, err := exec.Command("ps", "-o", "ppid=", "-p", fmt.Sprint(pid)).Output()
	if err != nil {
		return 0
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0
	}
	var p int
	fmt.Sscanf(s, "%d", &p)
	return p
}

func getProcessName(pid int) string {
	out, err := exec.Command("ps", "-p", fmt.Sprint(pid), "-o", "comm=").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func findFirstFile(paths []string) string {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
