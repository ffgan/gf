package dev

import cli "github.com/ffgan/gf/internal/CLI"

func DetectGPU() string {
	return getGPU()
}

func DetectResolution() string {
	return getResolution()
}

func DetectNetwork() string {
	return getNetwork(cli.GetOS())
}

func DetectBluetooth() string {
	return getBluetooth()
}

func GetBIOS() string {
	return getBIOS(cli.GetOS())
}
