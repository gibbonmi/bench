package axi

import (
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func runBenchInDir(t testing.TB, f contract.Fixture, dir string, args ...string) contract.Probe {
	t.Helper()
	selection := contract.SelectedBench(t)
	env := map[string]string{"BENCH_KIT": selection.SourceRoot, "BENCH_RUN_BINARY": selection.Path}
	return contract.RunAtWithInput(t, f, dir, env, "", selection.Path, args...)
}

func bashPath(t testing.TB) string {
	t.Helper()
	bin, err := exec.LookPath("bash")
	if err != nil {
		t.Fatalf("find bash: %v", err)
	}
	return bin
}

func writeExecutableAt(t testing.TB, dir, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o755); err != nil {
		t.Fatalf("write executable %s: %v", name, err)
	}
}

// writeGitFailureShim writes a "git" executable to dir that fails one chosen
// post-resolution `bench diff` call while delegating every other git invocation to
// the real binary — the PATH-stub pattern from testAXIBlockDangerousGitCoreErrored,
// generalized here so base resolution (symbolic-ref, config, cat-file, merge-base)
// still succeeds and only the named call breaks. Which call fails is selected at run
// time via the FAIL_GIT_CALL env var ("files", "log", or "body"), read by the shim
// itself so one script serves all three probes.
func writeGitFailureShim(t testing.TB, dir string) {
	t.Helper()
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("find real git: %v", err)
	}
	script := "#!/usr/bin/env bash\n" +
		"case \"$FAIL_GIT_CALL\" in\n" +
		"  files) [ \"$1\" = diff ] && [ \"$2\" = --name-status ] && exit 17 ;;\n" +
		"  log) [ \"$1\" = log ] && exit 17 ;;\n" +
		"  body) [ \"$1\" = diff ] && [ \"$2\" != --name-status ] && exit 17 ;;\n" +
		"esac\n" +
		"exec \"" + realGit + "\" \"$@\"\n"
	writeExecutableAt(t, dir, "git", script)
}

func requireAXIFirstLine(t testing.TB, out, want string) {
	t.Helper()
	if got := axiFirstLine(out); got != want {
		t.Fatalf("first line = %q, want %q\noutput:\n%s", got, want, out)
	}
}

func axiFirstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func requireAXILine(t testing.TB, out, want string) {
	t.Helper()
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if line == want {
			return
		}
	}
	t.Fatalf("missing line %q\noutput:\n%s", want, out)
}

func requireIndentedRowCount(t testing.TB, out string, want int) {
	t.Helper()
	got := 0
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if strings.HasPrefix(line, "  ") {
			got++
		}
	}
	if got != want {
		t.Fatalf("indented row count = %d, want %d\noutput:\n%s", got, want, out)
	}
}

func requireNoAXILineMatching(t testing.TB, out, expr string) {
	t.Helper()
	re := regexp.MustCompile(expr)
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if re.MatchString(line) {
			t.Fatalf("unexpected line matching %q: %q\noutput:\n%s", expr, line, out)
		}
	}
}

func requireContainsFold(t testing.TB, haystack, needle string) {
	t.Helper()
	if !strings.Contains(strings.ToLower(haystack), strings.ToLower(needle)) {
		t.Fatalf("missing %q\noutput:\n%s", needle, haystack)
	}
}

// requireLogRow asserts a `log` table row for sha/renderedSubject, tolerating the
// TOON encoder's own choice to quote the sha cell (a short sha can start with a
// digit and read as numeric-ish, e.g. "05af9b0") — that quoting is independent of
// anything this seam controls, so the row match accepts either form.
func requireLogRow(t *testing.T, probe contract.Probe, sha, renderedSubject string) {
	t.Helper()
	plain := "  " + sha + "," + renderedSubject
	quoted := "  \"" + sha + "\"," + renderedSubject
	for _, got := range strings.Split(strings.TrimSuffix(probe.Stdout, "\n"), "\n") {
		if got == plain || got == quoted {
			return
		}
	}
	t.Fatalf("missing log row for sha %q\nstdout:\n%s\nstderr:\n%s", sha, probe.Stdout, probe.Stderr)
}

func requireOutputLine(t *testing.T, probe contract.Probe, line string) {
	t.Helper()
	for _, got := range strings.Split(strings.TrimSuffix(probe.Stdout, "\n"), "\n") {
		if got == line {
			return
		}
	}
	t.Fatalf("missing output line %q\nstdout:\n%s\nstderr:\n%s", line, probe.Stdout, probe.Stderr)
}

func requireOutputPrefix(t *testing.T, probe contract.Probe, prefix string) {
	t.Helper()
	for _, got := range strings.Split(strings.TrimSuffix(probe.Stdout, "\n"), "\n") {
		if strings.HasPrefix(got, prefix) {
			return
		}
	}
	t.Fatalf("missing output line with prefix %q\nstdout:\n%s\nstderr:\n%s", prefix, probe.Stdout, probe.Stderr)
}

func requireNoOutput(t *testing.T, probe contract.Probe) {
	t.Helper()
	if probe.Stdout != "" || probe.Stderr != "" {
		t.Fatalf("expected no output\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
}
