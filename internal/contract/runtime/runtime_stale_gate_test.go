package runtime

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testRuntimeStatusStaleGateDriftClassification(t *testing.T) {
	benign := statusRowExpectation{
		contains:    []string{"stale (reduced-scope drift)", "re-run when convenient"},
		notContains: []string{"stale (gated tree"},
	}
	strong := statusRowExpectation{
		contains:    []string{"stale (gated tree", "re-run the gate"},
		notContains: []string{"reduced-scope drift"},
	}
	invalid := statusRowExpectation{
		contains:    []string{"invalid verdict", "re-run the gate"},
		notContains: []string{"reduced-scope drift", "stale (gated tree"},
	}
	cases := []staleGateStatusCase{
		{name: "added ROADMAP.md is capture-only", mutate: writeRuntimeFile("ROADMAP.md", "- 2026-07-05  parked idea\n"), want: benign},
		{name: "modified tracked session handoff is capture-only", seed: writeRuntimeFile("capture/session-handoff.md", "# Session handoff\n\nfirst capture\n"), mutate: writeRuntimeFile("capture/session-handoff.md", "# Session handoff\n\nrewritten capture\n"), want: benign},
		{name: "modified ROADMAP.md is capture-only", seed: writeRuntimeFile("ROADMAP.md", "- 2026-07-04  old idea\n"), mutate: writeRuntimeFile("ROADMAP.md", "- 2026-07-05  parked idea\n"), want: benign},
		{name: "deleted ROADMAP.md is capture-only", seed: writeRuntimeFile("ROADMAP.md", "- 2026-07-04  old idea\n"), mutate: removeRuntimePath("ROADMAP.md"), want: benign},
		{name: "added capture/IDEAS.md is capture-only", mutate: writeRuntimeFile("capture/IDEAS.md", "- 2026-07-05  parked idea\n"), want: benign},
		{name: "modified capture/IDEAS.md is capture-only", seed: writeRuntimeFile("capture/IDEAS.md", "- 2026-07-04  old idea\n"), mutate: writeRuntimeFile("capture/IDEAS.md", "- 2026-07-05  parked idea\n"), want: benign},
		{name: "added decision map is reduced-scope", mutate: writeRuntimeFile("decisions/map.md", "# Decision map\n"), want: benign},
		{name: "nested IDEAS lookalike is strong stale", mutate: writeRuntimeFile("notes/capture/IDEAS.md", "doc drift\n"), want: strong},
		{name: "added .bench-notes.md is capture-only", mutate: writeRuntimeFile(".bench-notes.md", "scratch\n"), want: benign},
		{name: "modified .bench-notes.md is capture-only", seed: writeRuntimeFile(".bench-notes.md", "old\n"), mutate: writeRuntimeFile(".bench-notes.md", "new\n"), want: benign},
		{name: "deleted .bench-notes.md is capture-only", seed: writeRuntimeFile(".bench-notes.md", "old\n"), mutate: removeRuntimePath(".bench-notes.md"), want: benign},
		{name: "docs ROADMAP lookalike is strong stale", mutate: writeRuntimeFile("docs/ROADMAP.md", "doc drift\n"), want: strong},
		{name: "ROADMAP backup lookalike is strong stale", mutate: writeRuntimeFile("ROADMAP.md.bak", "doc drift\n"), want: strong},
		{name: ".bench notes lookalike is strong stale", mutate: writeRuntimeFile(".bench/notes.md", "doc drift\n"), want: strong},
		{name: "nested ROADMAP lookalike is strong stale", mutate: writeRuntimeFile("notes/ROADMAP.md", "doc drift\n"), want: strong},
		{name: "mixed capture-only and real drift is strong stale", mutate: func(_ testing.TB, f contract.Fixture) {
			f.WriteFile("ROADMAP.md", "- 2026-07-05  parked idea\n")
			f.WriteFile("docs/x.md", "doc drift\n")
		}, want: strong},
		{name: "cache missing tree is invalid", cache: literalGateCache("green\n"), want: invalid},
		{name: "cache tree none is invalid", cache: literalGateCache("green none 2026-06-30T00:00:00Z\n"), want: invalid},
		{name: "cache missing timestamp is strong stale", cache: func(t testing.TB, f contract.Fixture) string {
			t.Helper()
			return fmt.Sprintf("green %s\n", strings.TrimSpace(f.Bench("tree-hash").Stdout))
		}, want: invalid},
		{name: "untrusted cache status is invalid", cache: func(t testing.TB, f contract.Fixture) string {
			t.Helper()
			return gateCacheLine(t, f, "yellow")
		}, want: invalid},
		{name: "deep cwd ROADMAP drift is capture-only", seed: writeRuntimeFile("sub/.keep", ""), mutate: writeRuntimeFile("ROADMAP.md", "- 2026-07-05  parked idea\n"), runDir: "sub", want: benign},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertStaleGateStatusCase(t, tc)
		})
	}
}

type staleGateStatusCase struct {
	name   string
	seed   runtimeFixtureEdit
	mutate runtimeFixtureEdit
	cache  func(testing.TB, contract.Fixture) string
	runDir string
	want   statusRowExpectation
}

type runtimeFixtureEdit func(testing.TB, contract.Fixture)

type statusRowExpectation struct {
	contains    []string
	notContains []string
}

func assertStaleGateStatusCase(t testing.TB, tc staleGateStatusCase) {
	t.Helper()
	f := contract.NewFixture(t)
	commitAllowEmpty(t, f, "init")
	if tc.seed != nil {
		tc.seed(t, f)
		f.CommitAll("seed")
	}
	cache := gateCacheLine(t, f, "green")
	if tc.cache != nil {
		cache = tc.cache(t, f)
	}
	writeGateCache(t, f, cache)
	if tc.mutate != nil {
		tc.mutate(t, f)
	}
	out := f.Bench("status").Stdout
	if tc.runDir != "" {
		out = contract.RunAt(t, f, filepath.Join(f.Root, filepath.FromSlash(tc.runDir)), nil, "bash", benchPath(t), "status").Stdout
	}
	requireStatusRow(t, out, tc.want)
}

func requireStatusRow(t testing.TB, out string, want statusRowExpectation) {
	t.Helper()
	for _, needle := range want.contains {
		contract.RequireContains(t, out, needle)
	}
	for _, needle := range want.notContains {
		contract.RequireNotContains(t, out, needle)
	}
}

func writeRuntimeFile(path, contents string) runtimeFixtureEdit {
	return func(_ testing.TB, f contract.Fixture) {
		f.WriteFile(path, contents)
	}
}

func removeRuntimePath(path string) runtimeFixtureEdit {
	return func(t testing.TB, f contract.Fixture) {
		t.Helper()
		contract.Remove(t, filepath.Join(f.Root, filepath.FromSlash(path)))
	}
}

func literalGateCache(line string) func(testing.TB, contract.Fixture) string {
	return func(t testing.TB, _ contract.Fixture) string {
		t.Helper()
		return line
	}
}

func writeGateCache(t testing.TB, f contract.Fixture, line string) {
	t.Helper()
	path := filepath.Join(gitDir(t, f), "bench-last-gate")
	contract.WriteFileAbs(t, path, line)
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatal(err)
	}
}

func gateCacheLine(t testing.TB, f contract.Fixture, status string) string {
	t.Helper()
	tree := strings.TrimSpace(f.Bench("tree-hash").Stdout)
	return fmt.Sprintf(`{"schema":1,"state":"ready","status":%q,"tree":%q,"oracle":"%s","recorded_at":%q}`+"\n", status, tree, strings.Repeat("0", 64), time.Now().UTC().Truncate(time.Second).Format(time.RFC3339))
}
