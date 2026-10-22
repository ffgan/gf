//go:build !linux

package gui

func waylandPeerPID(string) int { return 0 }
