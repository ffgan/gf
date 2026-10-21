package dev

import (
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

// GetSong asks playerctl for the current track. format understands
// %artist%, %album% and %title%.
func GetSong(format, player string) string {
	if !utils.CommandExists("playerctl") {
		return ""
	}
	template := "{{artist}} - {{title}}"
	switch {
	case strings.Contains(format, "%album%") && strings.Contains(format, "%artist%") && strings.Contains(format, "%title%"):
		template = "{{artist}} - {{album}} - {{title}}"
	case strings.Contains(format, "%title%") && !strings.Contains(format, "%artist%"):
		template = "{{title}}"
	}
	args := []string{"metadata", "--format", template}
	player = strings.TrimSpace(player)
	if player != "" && player != "auto" && player != "playerctl" {
		args = append([]string{"-p", player}, args...)
	}
	out := utils.RunCommand("playerctl", args...)
	if strings.Contains(out, "No players found") || strings.Contains(out, "No player could handle this command") {
		return ""
	}
	return strings.TrimSpace(out)
}
