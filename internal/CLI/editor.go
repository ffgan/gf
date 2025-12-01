package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

func GetEditor(editorPathOpt, editorVersionOpt string) string {
	// function from neofetch's  --> get_editor()
	editorFullPath := os.Getenv("VISUAL")
	if editorFullPath == "" {
		editorFullPath = os.Getenv("EDITOR")
	}
	if editorFullPath == "" {
		return ""
	}

	cmdPath, err := utils.LookPath(editorFullPath)
	if err != nil {
		return ""
	}
	editorFullPath = cmdPath

	if fi, err := os.Lstat(editorFullPath); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		if resolved, err := filepath.EvalSymlinks(editorFullPath); err == nil {
			editorFullPath = resolved
		}
	}

	var editor string
	switch editorPathOpt {
	case utils.ON:
		editor = editorFullPath + " "
	case utils.OFF:
		editor = filepath.Base(editorFullPath) + " "
	default:
		editor = editorFullPath + " "
	}

	if editorVersionOpt != utils.ON {
		return strings.TrimSpace(editor)
	}

	editorName := filepath.Base(editorFullPath)

	var args []string
	switch editorName {
	case "nano", "vim", "nvim", "micro", "emacs":
		args = []string{"--version"}
	case "kak":
		args = []string{"-version"}
	case "ne":
		args = []string{"-h"}
	default:
		return strings.TrimSpace(editor)
	}

	editorV := utils.RunCommand(editorFullPath, args...)
	if idx := strings.Index(editorV, "\n"); idx != -1 {
		editorV = editorV[:idx]
	}

	editorV = strings.Replace(editorV, "Version: ", "", 1)

	if strings.Contains(strings.ToLower(editorV), strings.ToLower(editorName)) {
		editorName = ""
	}

	editor = strings.TrimSpace(fmt.Sprintf("%s %s", editorName, editorV))

	editor = strings.Replace(editor, ", version", "", 1)
	if idx := strings.Index(editor, "options"); idx != -1 {
		editor = editor[:idx]
	}
	if idx := strings.Index(editor, "("); idx != -1 {
		editor = editor[:idx]
	}

	return strings.TrimSpace(editor)
}
