//go:build !linux

package cli

func unameFast(string) (string, bool) { return "", false }
