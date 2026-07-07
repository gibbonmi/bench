package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestRuntimeDashboardContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "bench dashboard write contract", testRuntimeDashboardWrite)
	contract.RunParallel(t, "bench dashboard --stdout contract", testRuntimeDashboardStdout)
	contract.RunParallel(t, "bench dashboard atomic-write contract", testRuntimeDashboardAtomicWrite)
	contract.RunParallel(t, "bench dashboard idempotence contract", testRuntimeDashboardIdempotent)
	contract.RunParallel(t, "bench dashboard error-posture contract", testRuntimeDashboardErrors)
}

// testRuntimeDashboardWrite pins story 1: the command writes a self-contained HTML page to
// <git-dir>/bench-dashboard.html and prints that exact path on stdout, exit 0.
func testRuntimeDashboardWrite(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("ROADMAP.md", "# Roadmap\n\nFT1 row\n\n## Recommended sequence\n\n1. a - /bench-shape-idea\n2. b - /bench-implement-spec\n")
	f.CommitAll("base")

	p := f.Bench("dashboard")
	p.RequireExit(0)
	want := filepath.Join(gitDir(t, f), "bench-dashboard.html")
	if got := strings.TrimSpace(p.Stdout); got != want {
		t.Fatalf("printed path %q, want %q", got, want)
	}
	data := contract.ReadFileAbs(t, want)
	if data == "" {
		t.Fatalf("dashboard wrote an empty file at %s", want)
	}
	contract.RequireContains(t, data, "<!DOCTYPE html>")
	contract.RequireContains(t, data, "FT1 row")
}

// testRuntimeDashboardStdout pins story 2: --stdout emits the document on stdout and writes
// no file under the git dir.
func testRuntimeDashboardStdout(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "base")

	p := f.Bench("dashboard", "--stdout")
	p.RequireExit(0)
	if !strings.HasPrefix(p.Stdout, "<!DOCTYPE html>") {
		t.Fatalf("--stdout did not emit an HTML document:\n%.60s", p.Stdout)
	}
	target := filepath.Join(gitDir(t, f), "bench-dashboard.html")
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("--stdout wrote a file at %s (err=%v), want none", target, err)
	}
}

// testRuntimeDashboardAtomicWrite pins story 12: after a successful run the target is a
// complete document and no sibling temp file is left behind.
func testRuntimeDashboardAtomicWrite(t *testing.T) {
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "base")
	f.Bench("dashboard").RequireExit(0)

	dir := gitDir(t, f)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read git dir %s: %v", dir, err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "bench-dashboard-") || strings.HasSuffix(e.Name(), ".tmp") {
			t.Fatalf("leftover temp artifact in git dir: %s", e.Name())
		}
	}
	page := contract.ReadFileAbs(t, filepath.Join(dir, "bench-dashboard.html"))
	if !strings.HasSuffix(strings.TrimRight(page, "\n"), "</html>") {
		t.Fatalf("target is not a complete document (no closing </html>):\n%.80s", page)
	}
}

// testRuntimeDashboardIdempotent pins edge-of-1: re-running overwrites the same single path
// with content equal to the first run once the generation timestamp is set aside.
func testRuntimeDashboardIdempotent(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("IDEAS.md", "- 2026-07-01  parked idea\n")
	f.CommitAll("base")

	first := f.Bench("dashboard")
	first.RequireExit(0)
	target := strings.TrimSpace(first.Stdout)
	firstPage := stripGenerated(contract.ReadFileAbs(t, target))

	second := f.Bench("dashboard")
	second.RequireExit(0)
	if got := strings.TrimSpace(second.Stdout); got != target {
		t.Fatalf("second run wrote a different path %q, want %q", got, target)
	}
	if secondPage := stripGenerated(contract.ReadFileAbs(t, target)); secondPage != firstPage {
		t.Fatalf("re-run changed the non-timestamp content:\nfirst:\n%s\nsecond:\n%s", firstPage, secondPage)
	}
	// Exactly one dashboard file at the path — no per-run variant names.
	entries, _ := os.ReadDir(filepath.Dir(target))
	n := 0
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "bench-dashboard") {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("want exactly one bench-dashboard file, got %d", n)
	}
}

// testRuntimeDashboardErrors pins story 13: outside a repo is the structured not-in-repo
// error (exit 1); an unknown argument is a usage line (exit 2); neither writes a file.
func testRuntimeDashboardErrors(t *testing.T) {
	noRepo := contract.NewFixture(t, contract.WithNoRepo())
	outside := noRepo.Bench("dashboard")
	outside.RequireExit(1)
	contract.RequireContains(t, outside.Stdout, "not in a git repository")

	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "base")
	bad := f.Bench("dashboard", "--bogus")
	bad.RequireExit(2)
	contract.RequireContains(t, bad.Stdout, "usage:")
	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-dashboard.html")); !os.IsNotExist(err) {
		t.Fatalf("a usage error still wrote the dashboard file")
	}
}

// stripGenerated removes the one per-run-varying line (the generation timestamp) so two
// runs' pages can be compared for structural equivalence.
func stripGenerated(page string) string {
	var kept []string
	for _, line := range strings.Split(page, "\n") {
		if strings.Contains(line, "class=\"generated\"") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}
