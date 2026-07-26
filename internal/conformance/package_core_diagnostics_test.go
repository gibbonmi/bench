package conformance

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func formatProbeFailure(label string, probe *Probe) string {
	var output string
	if probe != nil {
		output = probe.Stdout
		if output != "" && probe.Stderr != "" {
			output += "\n"
		}
		output += probe.Stderr
	}
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
	return label + ":\n" + tail
}

func TestFormatProbeFailureBoundsAndSanitizesTail(t *testing.T) {
	lines := make([]string, 10_000)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %05d\x1b[31m\x07\x7f", i)
	}
	got := formatProbeFailure("go test failed", &Probe{Stderr: strings.Join(lines, "\n")})

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

	got := formatProbeFailure("go build failed", &Probe{Stdout: strings.Join(lines, "\n")})

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
			got := formatProbeFailure("go vet failed", tc.probe)
			if got != "go vet failed:\n" {
				t.Fatalf("formatProbeFailure = %q, want labeled empty diagnostic", got)
			}
			if strings.ContainsAny(got, "\x1b\x07\x7f") {
				t.Fatalf("formatProbeFailure retained control bytes: %q", got)
			}
		})
	}
}

func TestCoreSubprocessFailuresUseProbeFormatter(t *testing.T) {
	h := NewHarness(t)
	source, err := os.ReadFile(h.KitPath("internal", "conformance", "package_core_checks_test.go"))
	if err != nil {
		t.Fatal(err)
	}
	labels := map[string]int{
		"npm pack --dry-run failed":              1,
		"go build setup failed":                  1,
		"go build failed":                        2,
		"go vet failed":                          1,
		"go list failed":                         1,
		"go test failed":                         2,
		"worktree cleanup race test failed":      1,
		"worktree cleanup race test did not run": 1,
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
