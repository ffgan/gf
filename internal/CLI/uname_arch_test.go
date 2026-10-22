package cli

import "testing"

func TestUnameMachine(t *testing.T) {
	cases := map[string]string{
		"amd64": "x86_64",
		"386":   "i686",
		"arm64": "aarch64",
		"arm":   "armv7l",
		"ppc64": "ppc64",
	}
	for arch, want := range cases {
		if got := unameMachine(arch); got != want {
			t.Fatalf("%s: got %q, want %q", arch, got, want)
		}
	}
}

func TestWindowsNTIsWindows(t *testing.T) {
	if got := GetOS("Windows_NT", ""); got != Windows {
		t.Fatalf("GetOS(Windows_NT) = %q", got)
	}
}
