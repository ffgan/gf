package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCountOpkgStatus(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "status")
	body := "Package: busybox\nVersion: 1\n\nPackage: libc\nVersion: 1\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := countPrefixedLines([]string{path}, "Package:"); got != 2 {
		t.Fatalf("got %d want 2", got)
	}
}

func TestCountManifestNames(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.nix")
	body := "(name \"hello\")\n(name \"coreutils\")\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := countManifestNames(path); got != 2 {
		t.Fatalf("got %d want 2", got)
	}
}
