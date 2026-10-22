package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScoopAndChocoCounts(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	global := filepath.Join(root, "global")
	for _, dir := range []string{
		filepath.Join(home, "apps", "git"),
		filepath.Join(home, "apps", "scoop"),
		filepath.Join(global, "apps", "7zip"),
		filepath.Join(global, "apps", "scoop"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	cfg := filepath.Join(root, "config.json")
	body := `{"root_path":"` + filepath.ToSlash(home) + `","global_path":"` + filepath.ToSlash(global) + `"}`
	if err := os.WriteFile(cfg, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	user, glob := scoopAppCounts(cfg, "", "")
	if user != 1 || glob != 1 {
		t.Fatalf("scoop = %d user, %d global", user, glob)
	}

	lib := filepath.Join(root, "choco", "lib")
	for _, name := range []string{"choco", "git", "notes.txt"} {
		path := filepath.Join(lib, name)
		if name == "notes.txt" {
			if err := os.MkdirAll(lib, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if n := countDirsExcept(lib, "choco"); n != 1 {
		t.Fatalf("choco packages = %d", n)
	}
}
