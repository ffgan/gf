package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPipMetaName(t *testing.T) {
	cases := map[string]string{
		"ansible-9.2.0.dist-info":             "ansible",
		"ansible_core-2.16.3.dist-info":       "ansible-core",
		"cryptography.egg-info":               "cryptography",
		"python_apt-2.7.7.egg-info":           "python-apt",
		"lazr.restfulclient-0.14.6.dist-info": "lazr.restfulclient",
		"setuptools":                          "",
		"pip-24.0.dist-info":                  "pip",
	}
	for in, want := range cases {
		if got := pipMetaName(in); got != want {
			t.Errorf("%s: got %q want %q", in, got, want)
		}
	}
}

func TestCountDpkgFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "status")
	body := "" +
		"Package: a\nStatus: install ok installed\n\n" +
		"Package: b\nStatus: deinstall ok config-files\n\n" +
		"Package: c\nStatus: not-installed\n\n" +
		"Package: d\nStatus: install ok installed\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := countDpkgFile(path); got != 2 {
		t.Fatalf("got %d want 2", got)
	}
}

func TestCountDirsIgnoresFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ALPM_DB_VERSION"), []byte("9"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := countDirs(dir); got != 1 {
		t.Fatalf("got %d want 1", got)
	}
}
