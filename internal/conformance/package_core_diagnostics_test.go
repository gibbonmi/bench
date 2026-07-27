package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// probeLogName is the spill target for full probe output, written beside the
// gate cache in the graded root's git dir so the bounded diagnostic tail can
// name where the untruncated output lives.
const probeLogName = "bench-conformance-probe.log"

var (
	probeLogMu    sync.Mutex
	probeLogFresh = map[string]bool{}
)

// spillProbeOutput writes the full, unsanitized probe output to the graded
// root's git-dir log and returns the path. A root with no resolvable git dir
// (a bare fixture tree) skips the spill and the truncated diagnostic stays
// the only surface. The first spill of a process truncates the previous
// run's log; later spills append, so one run's failures share one file.
func spillProbeOutput(root, label, output string) string {
	if root == "" || output == "" {
		return ""
	}
	probe := runAt(root, "git", "rev-parse", "--absolute-git-dir")
	if probe == nil || probe.ExitCode != 0 || strings.TrimSpace(probe.Stdout) == "" {
		return ""
	}
	path := filepath.Join(strings.TrimSpace(probe.Stdout), probeLogName)
	probeLogMu.Lock()
	defer probeLogMu.Unlock()
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if !probeLogFresh[path] {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	f, err := os.OpenFile(path, flags, 0o644)
	if err != nil {
		return ""
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "=== %s\n%s\n", label, output); err != nil {
		return ""
	}
	probeLogFresh[path] = true
	return path
}

func formatProbeFailure(label string, probe *Probe, root string) string {
	var output string
	if probe != nil {
		output = probe.Stdout
		if output != "" && probe.Stderr != "" {
			output += "\n"
		}
		output += probe.Stderr
	}
	logPath := spillProbeOutput(root, label, output)
	output = strings.Map(func(r rune) rune {
		if (r < 0x20 && r != '\n') || r == 0x7f {
			return -1
		}
		return r
	}, output)
	lines := strings.Split(output, "\n")
	if len(lines) > 40 {
		lines = lines[len(lines)-40:]
	}
	tail := strings.Join(lines, "\n")
	if len(tail) > 4*1024 {
		tail = tail[len(tail)-4*1024:]
	}
	diag := label + ":\n" + tail
	if logPath != "" {
		diag += "\nfull probe output: " + logPath
	}
	return diag
}

func TestFormatProbeFailureBoundsAndSanitizesTail(t *testing.T) {
	lines := make([]string, 10_000)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %05d\x1b[31m\x07\x7f", i)
	}
	got := formatProbeFailure("go test failed", &Probe{Stderr: strings.Join(lines, "\n")}, "")

	if !strings.HasPrefix(got, "go test failed:\n") {
		t.Fatalf("formatProbeFailure prefix = %q", got)
	}
	if len(got) > 4*1024+len("go test failed:\n") {
		t.Fatalf("formatProbeFailure length = %d, want at most %d", len(got), 4*1024+len("go test failed:\n"))
	}
	if strings.ContainsAny(got, "\x1b\x07\x7f") {
		t.Fatalf("formatProbeFailure retained control bytes: %q", got)
	}
	if strings.Contains(got, "line 09959") || !strings.Contains(got, "line 09960") || !strings.Contains(got, "line 09999") {
		t.Fatalf("formatProbeFailure did not retain exactly the final 40 lines: %q", got)
	}
}

func TestFormatProbeFailureTruncatesOversizedFinalLinesFromTheTail(t *testing.T) {
	lines := make([]string, 40)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %02d %s\x1b\x07\x7f", i, strings.Repeat("x", 200))
	}
	lines[len(lines)-1] += " FINAL-END"
	const prefix = "go build failed:\n"

	got := formatProbeFailure("go build failed", &Probe{Stdout: strings.Join(lines, "\n")}, "")

	if len(got) > len(prefix)+4*1024 {
		t.Fatalf("formatProbeFailure length = %d, want at most %d", len(got), len(prefix)+4*1024)
	}
	if !strings.HasSuffix(got, " FINAL-END") {
		t.Fatalf("formatProbeFailure lost the final expected bytes: %q", got)
	}
	if strings.ContainsAny(got, "\x1b\x07\x7f") {
		t.Fatalf("formatProbeFailure retained control bytes after truncation: %q", got)
	}
}

func TestFormatProbeFailureLabelsMissingOutput(t *testing.T) {
	for _, tc := range []struct {
		name  string
		probe *Probe
	}{
		{name: "nil probe"},
		{name: "empty probe", probe: &Probe{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := formatProbeFailure("go vet failed", tc.probe, "")
			if got != "go vet failed:\n" {
				t.Fatalf("formatProbeFailure = %q, want labeled empty diagnostic", got)
			}
			if strings.ContainsAny(got, "\x1b\x07\x7f") {
				t.Fatalf("formatProbeFailure retained control bytes: %q", got)
			}
		})
	}
}

func TestFormatProbeFailureSpillsFullOutputToGitDirLog(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	stale := filepath.Join(root, ".git", probeLogName)
	if err := os.WriteFile(stale, []byte("stale previous-run entry\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lines := make([]string, 200)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %05d", i)
	}

	got := formatProbeFailure("go test failed", &Probe{Stdout: strings.Join(lines, "\n")}, root)

	const marker = "\nfull probe output: "
	idx := strings.LastIndex(got, marker)
	if idx < 0 {
		t.Fatalf("diagnostic does not name the spill log: %q", got)
	}
	path := strings.TrimSpace(got[idx+len(marker):])
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the named spill log: %v", err)
	}
	log := string(data)
	if strings.Contains(log, "stale previous-run entry") {
		t.Fatalf("first spill did not truncate the previous run's log: %q", log)
	}
	if !strings.Contains(log, "=== go test failed") || !strings.Contains(log, "line 00000") {
		t.Fatalf("spill log lost the head the diagnostic tail drops: %q", log)
	}
	if strings.Contains(got[:idx], "line 00000") {
		t.Fatalf("diagnostic tail is no longer bounded: %q", got)
	}

	second := formatProbeFailure("go vet failed", &Probe{Stderr: "vet says no"}, root)
	if !strings.Contains(second, marker) {
		t.Fatalf("second diagnostic does not name the spill log: %q", second)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "=== go test failed") || !strings.Contains(string(data), "=== go vet failed") {
		t.Fatalf("second spill did not append to the run's log: %q", string(data))
	}
}

func TestFormatProbeFailureWithoutGitDirKeepsTruncatedDiagnosticOnly(t *testing.T) {
	got := formatProbeFailure("go test failed", &Probe{Stdout: "boom"}, t.TempDir())
	if got != "go test failed:\nboom" {
		t.Fatalf("formatProbeFailure = %q, want the plain truncated diagnostic for a root with no git dir", got)
	}
}

func TestCoreSubprocessFailuresUseProbeFormatter(t *testing.T) {
	h := NewHarness(t)
	source, err := os.ReadFile(h.KitPath("internal", "conformance", "package_core_checks_test.go"))
	if err != nil {
		t.Fatal(err)
	}
	// One label survives the toolchain split: every Go step that used to run here now
	// streams its own output from a gate phase, and only the npm pack probe is still a
	// subprocess this file has to frame.
	labels := map[string]int{
		"npm pack --dry-run failed": 1,
	}
	text := string(source)
	total := 0
	for label, want := range labels {
		if strings.Contains(text, `append(diags, "`+label) {
			t.Errorf("bare subprocess failure diag survives: %q", label)
		}
		if got := strings.Count(text, `formatProbeFailure("`+label+`",`); got != want {
			t.Errorf("formatter call sites for %q = %d, want %d", label, got, want)
		}
		total += want
	}
	if got := strings.Count(text, `formatProbeFailure("`); got != total {
		t.Errorf("formatter call sites = %d, want %d gate-reachable subprocess failures", got, total)
	}
}
