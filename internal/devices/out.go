package dev

import cli "github.com/ffgan/gf/internal/CLI"

func DetectGPU() string {
	return "GPU: " + getGPU()
}

func DetectResolution() string {
	return "Resolution: " + getResolution()
}

func DetectNetwork() string {
	return "Network: " + getNetwork(cli.GetOS())
}

func DetectBluetooth() string {
	return "Bluetooth: " + getBluetooth()
}

func GetBIOS() string {
	return "BIOS: " + getBIOS(cli.GetOS())
}
