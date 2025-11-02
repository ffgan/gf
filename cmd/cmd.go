package cmd

import (
	"github.com/ffgan/gf/configs"
	"github.com/ffgan/gf/internal/info"
)

func Run() {
	config, err := configs.LoadConfig("")
	if err != nil {
		panic(err)
	}

	parseArgs(config)

	ascii, lines := info.GetInfo(config)

	info.PrintInfo(ascii, lines)
}
