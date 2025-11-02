package main

import (
	"embed"

	"github.com/ffgan/gf/cmd"
)

//go:embed assets/*
var assets embed.FS

func main() {
	cmd.Run(assets)
}
