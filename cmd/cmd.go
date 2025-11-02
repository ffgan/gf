package cmd

import (
	"embed"

	"github.com/ffgan/gf/configs"
	"github.com/ffgan/gf/internal/info"
)

func Run(ASCIIFiles embed.FS) {
	config, err := configs.LoadConfig("")
	if err != nil {
		panic(err)
	}

	parseArgs(config)

	ascii, lines := info.GetInfo(config, ASCIIFiles)

	info.PrintInfo(ascii, lines)
}
