//go:build windows

package cli

import (
	"strings"
	"testing"
)

func TestWindowsUnameWithoutBinary(t *testing.T) {
	got, ok := unameFast("-srm")
	if !ok {
		t.Fatal("unameFast(-srm) declined")
	}
	fields := strings.Fields(got)
	if len(fields) != 3 || fields[0] != "Windows_NT" || fields[2] == "" {
		t.Fatalf("uname -srm = %q", got)
	}
	if _, ok := unameFast("-m"); !ok {
		t.Fatal("unameFast(-m) declined")
	}
}
