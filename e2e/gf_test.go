package e2e_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ffgan/gf/cmd"
)

// TestGF builds the real binary and runs it. Unit tests never entered
// main, so a panic while collecting a line (for example a missing Win32
// procedure) shipped without being seen.
func TestGF(t *testing.T) {
	bin := buildGF(t)

	t.Run("version", func(t *testing.T) {
		stdout, _ := runGF(t, bin, "--version")
		if got := strings.TrimSpace(stdout); got != "gf "+cmd.Version {
			t.Fatalf("version %q, want gf %s", got, cmd.Version)
		}
	})

	t.Run("help", func(t *testing.T) {
		stdout, stderr := runGF(t, bin, "--help")
		if !strings.Contains(stdout+stderr, "Usage:") {
			t.Fatalf("help text missing Usage:\n%s%s", stdout, stderr)
		}
	})

	t.Run("print_config", func(t *testing.T) {
		stdout, _ := runGF(t, bin, "--print_config")
		if !strings.Contains(stdout, "print_info()") || !strings.Contains(stdout, "info title") {
			t.Fatalf("config dump missing print_info:\n%s", stdout)
		}
	})

	t.Run("unknown_flag", func(t *testing.T) {
		stdout, stderr, err := runGFResult(t, bin, "--not-a-real-flag")
		if err == nil {
			t.Fatalf("unknown flag exited 0\n%s%s", stdout, stderr)
		}
		if strings.Contains(stdout+stderr, "panic:") {
			t.Fatalf("unknown flag panicked\n%s%s", stdout, stderr)
		}
	})

	t.Run("stdout", func(t *testing.T) {
		stdout, stderr := runGF(t, bin, "--stdout", "--no_config")
		if strings.Contains(stdout, "\x1b") {
			t.Fatalf("--stdout still emitted ANSI:\n%s", stdout)
		}
		if stderr != "" {
			t.Fatalf("stderr:\n%s", stderr)
		}
		assertInfo(t, stdout)
	})
}

func assertInfo(t *testing.T, stdout string) {
	t.Helper()
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(lines) < 4 || !strings.Contains(lines[0], "@") {
		t.Fatalf("first line should be user@host:\n%s", stdout)
	}
	if lines[1] == "" || strings.Trim(lines[1], "-") != "" {
		t.Fatalf("second line should be the title underline:\n%s", stdout)
	}

	labels := []string{"OS", "Kernel", "Uptime", "Shell", "CPU", "Memory", "Disk"}
	pos := map[string]int{}
	for i, line := range lines {
		name, _, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}
		if _, seen := pos[name]; !seen {
			pos[name] = i
		}
	}
	var missing []string
	prev := 0
	for _, label := range labels {
		i, ok := pos[label]
		if !ok || strings.TrimSpace(strings.TrimPrefix(lines[i], label+": ")) == "" {
			missing = append(missing, label)
			continue
		}
		if i < prev {
			t.Errorf("%s out of order in:\n%s", label, stdout)
		}
		prev = i
	}
	if len(missing) > 0 {
		t.Fatalf("missing %s\n%s", strings.Join(missing, ", "), stdout)
	}

	osLine := strings.TrimPrefix(lines[pos["OS"]], "OS: ")
	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(osLine, "Windows") {
			t.Errorf("OS line %q", osLine)
		}
	case "darwin":
		if !strings.Contains(osLine, "macOS") && !strings.Contains(osLine, "Mac OS") && !strings.Contains(osLine, "Darwin") {
			t.Errorf("OS line %q", osLine)
		}
	}
	if !strings.Contains(lines[pos["Memory"]], "/") {
		t.Errorf("memory line %q", lines[pos["Memory"]])
	}
}

func buildGF(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "gf")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	goBin := filepath.Join(runtime.GOROOT(), "bin", "go")
	cmd := exec.Command(goBin, "build", "-o", bin, ".")
	cmd.Dir = moduleRoot(t)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return bin
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func runGF(t *testing.T, bin string, args ...string) (string, string) {
	t.Helper()
	stdout, stderr, err := runGFResult(t, bin, args...)
	if err != nil {
		t.Fatalf("%v\nstdout:\n%s\nstderr:\n%s", err, stdout, stderr)
	}
	return stdout, stderr
}

func runGFResult(t *testing.T, bin string, args ...string) (string, string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	home := t.TempDir()
	cmd.Env = replaceEnv(os.Environ(), map[string]string{
		"HOME":            home,
		"USERPROFILE":     home,
		"XDG_CONFIG_HOME": filepath.Join(home, "config"),
		"XDG_CACHE_HOME":  filepath.Join(home, "cache"),
	})
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("%s timed out\nstdout:\n%s\nstderr:\n%s", bin, stdout.String(), stderr.String())
	}
	combined := stdout.String() + stderr.String()
	if strings.Contains(combined, "panic:") {
		t.Fatalf("panic\n%s", combined)
	}
	return stdout.String(), stderr.String(), err
}

func replaceEnv(base []string, overrides map[string]string) []string {
	out := make([]string, 0, len(base)+len(overrides))
	seen := map[string]bool{}
	for _, entry := range base {
		key, _, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		if value, replace := overrides[key]; replace {
			if !seen[key] {
				out = append(out, key+"="+value)
				seen[key] = true
			}
			continue
		}
		out = append(out, entry)
	}
	for key, value := range overrides {
		if !seen[key] {
			out = append(out, key+"="+value)
		}
	}
	return out
}
