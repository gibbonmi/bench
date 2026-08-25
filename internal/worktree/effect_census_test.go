package worktree

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// ambientReadPattern names the process-context reads only the effect boundary may
// perform. Lower owners receive home, root, time, and repository facts explicitly.
var ambientReadPattern = regexp.MustCompile(`\b(os\.Getenv|os\.LookupEnv|os\.Getwd|time\.Now)\(`)

// TestEffectBoundaryCensus fails when any production file below the boundary reads
// the environment, the current directory, or the clock. The one named boundary
// adapter file, effects.go, owns every ambient read; commands resolve there once
// and pass values down. (Coverage row EI1.)
func TestEffectBoundaryCensus(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		if name == "effects.go" {
			continue // the named effect boundary adapter file
		}
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			code := line
			if idx := strings.Index(code, "//"); idx >= 0 {
				code = code[:idx]
			}
			if match := ambientReadPattern.FindString(code); match != "" {
				t.Errorf("%s:%d: ambient read %s below the effect boundary; resolve it in effects.go and pass the value down", name, i+1, match)
			}
		}
	}
}

// TestHarnessCensusRefusesChildStartOutsideHarness proves the harness effect
// census still refuses a child-process start in a test file outside the harness
// file set. The census walks the package directory itself, so this test drives
// its two sources — the effect pattern and the harness file set — over a
// synthetic line instead. The subject text is assembled from parts, because the
// census reads this file too and a literal would report this line.
// (Coverage row WF13.)
func TestHarnessCensusRefusesChildStartOutsideHarness(t *testing.T) {
	t.Parallel()
	line := "\tcmd := " + "exec" + "." + "Command(\"git\", \"status\")"
	if journeyHarnessFiles["clean_test.go"] {
		t.Fatal("clean_test.go is not a harness file")
	}
	if match := outsideHarnessEffect.FindString(line); match == "" {
		t.Fatalf("the harness census accepts a child start outside the harness: %q", line)
	}
}

// TestPoolAtExplicitHome proves the pool path derives from the injected home value
// alone, with no environment read. (Coverage row EI1.)
func TestPoolAtExplicitHome(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "repos", "example")
	a := poolAt(filepath.Join(string(filepath.Separator), "home-a"), root)
	b := poolAt(filepath.Join(string(filepath.Separator), "home-b"), root)
	if !strings.HasPrefix(a, filepath.Join(string(filepath.Separator), "home-a", "worktrees")) {
		t.Fatalf("poolAt ignored the injected home: %s", a)
	}
	if filepath.Base(a) != filepath.Base(b) {
		t.Fatalf("pool key must depend only on the root: %s vs %s", a, b)
	}
}

// TestPoolKeysDirAtExplicitHome proves the pool parent derives from the injected home.
func TestPoolKeysDirAtExplicitHome(t *testing.T) {
	home := filepath.Join(string(filepath.Separator), "elsewhere", ".bench")
	if got, want := poolKeysDirAt(home), filepath.Join(home, "worktrees"); got != want {
		t.Fatalf("poolKeysDirAt = %s, want %s", got, want)
	}
}

// TestLeaseLineExplicitTime proves the lease record carries the injected instant,
// not an ambient clock read.
func TestLeaseLineExplicitTime(t *testing.T) {
	now := time.Date(2026, 8, 23, 10, 30, 0, 0, time.FixedZone("plus2", 2*3600))
	line := string(leaseLine(now))
	if !strings.Contains(line, now.UTC().Format(leaseTimeLayout)) {
		t.Fatalf("lease line %q does not carry the injected UTC instant", line)
	}
}
