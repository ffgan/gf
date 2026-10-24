package e2e_test

import (
	"bytes"
	"context"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
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

	labels := []string{"OS", "Kernel", "Uptime"}
	// Windows CI does not set SHELL, and GetShell stays empty without it.
	if os.Getenv("SHELL") != "" {
		labels = append(labels, "Shell")
	}
	labels = append(labels, "CPU", "Memory", "Disk")
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
	if disk, ok := pos["Disk"]; ok {
		assertDiskLine(t, strings.TrimPrefix(lines[disk], "Disk: "))
	}
	if ip, ok := pos["Local IP"]; ok {
		assertLocalIP(t, strings.TrimPrefix(lines[ip], "Local IP: "))
	}
}

func assertDiskLine(t *testing.T, value string) {
	t.Helper()
	if !diskLinePattern.MatchString(strings.TrimSpace(value)) {
		t.Fatalf("disk line %q", value)
	}
}

func assertLocalIP(t *testing.T, value string) {
	t.Helper()
	value = strings.TrimSpace(value)
	// Windows prints "(Ethernet) 10.1.0.22/20". Other systems print the address.
	if strings.HasPrefix(value, "(") {
		_, rest, ok := strings.Cut(value, ") ")
		if !ok {
			t.Fatalf("local IP %q", value)
		}
		value = rest
	}
	ipText, _, _ := strings.Cut(value, "/")
	ip := net.ParseIP(ipText)
	if ip == nil || ip.To4() == nil || !ip.IsGlobalUnicast() {
		t.Fatalf("local IP %q", value)
	}
}

// Unix df prints "7.0G". Windows prints "12.34 GiB" and may append " - NTFS".
var diskLinePattern = regexp.MustCompile(`^\(.+\) [\d.]+ ?(?:[KMGTP](?:iB)?|B) / [\d.]+ ?(?:[KMGTP](?:iB)?|B) \(\d+%\)(?: - .+)?$`)

// TestMacOSDiskAndLocalIP compares the binary with df and the default
// route on a real Mac. Hosted runners do not share a disk image with
// Linux CI, so this stays in the e2e job instead of a parsed fixture.
func TestMacOSDiskAndLocalIP(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS storage and local IP need a Mac")
	}
	bin := buildGF(t)
	stdout, _ := runGF(t, bin, "--stdout", "--no_config")

	used, total, avail := dfBytes(t, "/")
	wantUsed := used
	if fi, err := os.Stat("/System/Volumes/Data"); err == nil && fi.IsDir() &&
		avail >= 0 && used >= 0 && total > avail && uint64(used)+uint64(avail) < uint64(total) {
		wantUsed = total - avail
	}
	if total <= 0 {
		t.Fatalf("df / total is %d", total)
	}
	wantPct := int(uint64(wantUsed) * 100 / uint64(total))
	disk := infoField(stdout, "Disk")
	gotUsed, gotTotal, gotPct := parseDiskLine(t, disk)
	if gotPct < wantPct-1 || gotPct > wantPct+1 {
		t.Fatalf("disk %q percent %d, want %d (df used %d total %d avail %d)\n%s", disk, gotPct, wantPct, used, total, avail, stdout)
	}
	if !nearBytes(gotUsed, float64(wantUsed)) || !nearBytes(gotTotal, float64(total)) {
		t.Fatalf("disk %q is %f / %f bytes, want about %d / %d\n%s", disk, gotUsed, gotTotal, wantUsed, total, stdout)
	}

	iface := routeInterface(t)
	wantIP := commandTry("ipconfig", "getifaddr", iface)
	if !globalIPv4(wantIP) {
		wantIP = ifconfigIPv4(t, iface)
	}
	if !globalIPv4(wantIP) {
		t.Fatalf("no global IPv4 on %s", iface)
	}
	if got := infoField(stdout, "Local IP"); got != wantIP {
		t.Fatalf("local IP %q, want %q via %s\n%s", got, wantIP, iface, stdout)
	}
}

func infoField(stdout, label string) string {
	prefix := label + ": "
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}

func dfBytes(t *testing.T, path string) (used, total, avail int64) {
	t.Helper()
	text := commandOut(t, "df", "-k", "-P", path)
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 2 {
		t.Fatalf("df %s: %q", path, text)
	}
	fields := strings.Fields(lines[len(lines)-1])
	pct := -1
	for i, field := range fields {
		if strings.HasSuffix(field, "%") {
			pct = i
			break
		}
	}
	if pct < 3 {
		t.Fatalf("df columns %q", lines[len(lines)-1])
	}
	total = parseInt(t, fields[pct-3]) * 1024
	used = parseInt(t, fields[pct-2]) * 1024
	avail = parseInt(t, fields[pct-1]) * 1024
	return used, total, avail
}

func parseInt(t *testing.T, text string) int64 {
	t.Helper()
	n, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

var shownSize = regexp.MustCompile(`([\d.]+) ?([KMGTP](?:iB)?|B)`)

func parseDiskLine(t *testing.T, disk string) (used, total float64, percent int) {
	t.Helper()
	sizes := shownSize.FindAllStringSubmatch(disk, 2)
	if len(sizes) < 2 {
		t.Fatalf("disk sizes %q", disk)
	}
	pct := regexp.MustCompile(`\((\d+)%\)`).FindStringSubmatch(disk)
	if pct == nil {
		t.Fatalf("disk percent %q", disk)
	}
	p, err := strconv.Atoi(pct[1])
	if err != nil {
		t.Fatal(err)
	}
	return sizeBytes(sizes[0][1], sizes[0][2]), sizeBytes(sizes[1][1], sizes[1][2]), p
}

func sizeBytes(num, unit string) float64 {
	value, _ := strconv.ParseFloat(num, 64)
	unit = strings.TrimSuffix(unit, "iB")
	unit = strings.TrimSuffix(unit, "B")
	mul := 1.0
	switch unit {
	case "K":
		mul = 1024
	case "M":
		mul = 1024 * 1024
	case "G":
		mul = 1024 * 1024 * 1024
	case "T":
		mul = 1024 * 1024 * 1024 * 1024
	case "P":
		mul = 1024 * 1024 * 1024 * 1024 * 1024
	}
	return value * mul
}

func nearBytes(got, want float64) bool {
	if want == 0 {
		return got == 0
	}
	delta := got - want
	if delta < 0 {
		delta = -delta
	}
	return delta/want < 0.015
}

func routeInterface(t *testing.T) string {
	t.Helper()
	text := commandOut(t, "route", "-n", "get", "default")
	for _, line := range strings.Split(text, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if ok && strings.TrimSpace(key) == "interface" {
			fields := strings.Fields(value)
			if len(fields) > 0 {
				return fields[0]
			}
		}
	}
	t.Fatalf("default route has no interface:\n%s", text)
	return ""
}

func ifconfigIPv4(t *testing.T, iface string) string {
	t.Helper()
	text, err := exec.Command("ifconfig", iface).Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(text), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "inet" && globalIPv4(fields[1]) {
			return fields[1]
		}
	}
	return ""
}

func globalIPv4(text string) bool {
	ip := net.ParseIP(strings.TrimSpace(text))
	return ip != nil && ip.To4() != nil && ip.IsGlobalUnicast()
}

func commandTry(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func commandOut(t *testing.T, name string, args ...string) string {
	t.Helper()
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
	}
	return string(out)
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
