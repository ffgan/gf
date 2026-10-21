//go:build !linux

package utils

func parentPIDFast(int) (int, bool)         { return 0, false }
func processNameFast(string) (string, bool) { return "", false }
func processNamesFast() ([]string, bool)    { return nil, false }
