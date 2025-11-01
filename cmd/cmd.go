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

	// fmt.Printf("%+v", config)
	// fmt.Println()

	info.Setup(config)

	info.PrintInfo(config)
}
