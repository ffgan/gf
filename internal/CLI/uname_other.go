//go:build !linux && !windows

package cli

func unameFast(string) (string, bool) { return "", false }
