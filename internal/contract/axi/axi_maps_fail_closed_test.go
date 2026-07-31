package axi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestAXIMapsActiveProjectionAndCount(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("decisions/model.md", activeMapDocument("shaping", "- Honest fog.", "— (open)", "— (open)", "— (deferred)"))
	f.WriteFile("specs/compiled/decisions/ignored.md", "# Invalid compiled map\n")

	rows := f.Bench("maps")
	rows.RequireExit(0)
	requireAXIFirstLine(t, rows.Stdout, "maps[3]{map,title,type,state,blockers}:")
	requireAXILine(t, rows.Stdout, "  model,First,Research,frontier,\"\"")
	requireAXILine(t, rows.Stdout, "  model,Second,Task,blocked,First")
	requireAXILine(t, rows.Stdout, "  model,Third,Prototype,deferred,First")
	if strings.Contains(rows.Stdout, "ignored") {
		t.Fatalf("compiled map entered active rows:\n%s", rows.Stdout)
	}

	count := f.Bench("maps", "--count")
	count.RequireExit(0)
	if got := strings.TrimSpace(count.Stdout); got != "1" {
		t.Fatalf("maps --count = %q, want 1", got)
	}

	status := f.Bench("status", "--all")
	status.RequireExit(0)
	requireContainsFold(t, status.Stdout, "1 unresolved map")
}

func TestAXIMapsFogAndInvalidCountOnce(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("decisions/fog.md", activeMapDocument("shaping", "- Honest fog.", "Resolved.", "Resolved.", "Resolved."))
	f.WriteFile("decisions/bad.md", "# Bad\n")
	f.WriteFile("decisions/ready.md", activeMapDocument("ready", "", "Resolved.", "Resolved.", "Resolved."))

	rows := f.Bench("maps")
	rows.RequireExit(1)
	requireAXILine(t, rows.Stdout, "  fog,Not yet specified,fog,shaping,\"\"")
	requireAXILine(t, rows.Stdout, "  bad,invalid,map,invalid,\"decisions/bad.md: missing Status\"")
	requireNoAXILineMatching(t, rows.Stdout, `^  ready,`)

	count := f.Bench("maps", "--count")
	count.RequireExit(0)
	if got := strings.TrimSpace(count.Stdout); got != "2" {
		t.Fatalf("maps --count = %q, want 2", got)
	}
}

func TestAXIMapsEmptyFogShapingCountAgreesWithListing(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("decisions/silent.md", activeMapDocument("shaping", "", "Resolved.", "Resolved.", "Resolved."))

	rows := f.Bench("maps")
	rows.RequireExit(0)
	requireAXIFirstLine(t, rows.Stdout, "maps[1]{map,title,type,state,blockers}:")
	requireAXILine(t, rows.Stdout, "  silent,Not yet specified,fog,shaping,\"\"")

	count := f.Bench("maps", "--count")
	count.RequireExit(0)
	if got := strings.TrimSpace(count.Stdout); got != "1" {
		t.Fatalf("maps --count = %q, want 1", got)
	}
}

func TestAXIMapsFailClosedCandidatesAndAbsentDirectory(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)

	t.Run("absent directory is definitive empty", func(t *testing.T) {
		f := contract.NewFixture(t)
		out := f.Bench("maps")
		out.RequireExit(0)
		requireAXIFirstLine(t, out.Stdout, "maps[0]{map,title,type,state,blockers}:")
	})
	t.Run("malformed and unreadable candidates remain visible", func(t *testing.T) {
		f := contract.NewFixture(t)
		f.WriteFile("decisions/badutf8.md", "\xff\xfe")
		f.WriteFile("decisions/locked.md", activeMapDocument("shaping", "", "— (open)", "Resolved.", "Resolved."))
		f.WriteUnreadable("decisions/locked.md", 0o644)
		out := f.Bench("maps")
		out.RequireExit(1)
		requireContainsFold(t, out.Stdout, "badutf8,invalid,map,invalid")
		requireContainsFold(t, out.Stdout, "locked,invalid,map,invalid")
	})
}

func TestAXIMapsTemplateAndGrammar(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	template := f.Bench("maps", "--template")
	template.RequireExit(0)
	requireContainsFold(t, template.Stdout, "Status: shaping")
	for _, args := range [][]string{{"maps", "--count", "--template"}, {"maps", "--template", "--count"}} {
		out := f.Bench(args...)
		out.RequireExit(2)
		requireContainsFold(t, out.Stdout, "--count and --template are mutually exclusive")
	}
}

func TestAXIMapsQueriesAreReadOnly(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("decisions/model.md", activeMapDocument("shaping", "", "— (open)", "— (open)", "— (deferred)"))
	f.Git("config", "user.email", "bench@local").RequireExit(0)
	f.Git("config", "user.name", "bench").RequireExit(0)
	f.Git("add", ".").RequireExit(0)
	f.Git("commit", "-m", "seed").RequireExit(0)
	f.WriteFile("scratch.txt", "untracked\n")
	tracked, err := os.ReadFile(filepath.Join(f.Root, "decisions", "model.md"))
	if err != nil {
		t.Fatal(err)
	}
	untracked, err := os.ReadFile(filepath.Join(f.Root, "scratch.txt"))
	if err != nil {
		t.Fatal(err)
	}
	before := f.Git("status", "--porcelain").Stdout
	for _, args := range [][]string{{"maps"}, {"maps", "--count"}, {"maps", "--template"}, {"maps"}, {"maps", "--count"}, {"maps", "--template"}} {
		f.Bench(args...).RequireExit(0)
	}
	after := f.Git("status", "--porcelain").Stdout
	got, err := os.ReadFile(filepath.Join(f.Root, "decisions", "model.md"))
	if err != nil {
		t.Fatal(err)
	}
	scratch, err := os.ReadFile(filepath.Join(f.Root, "scratch.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if before != after || string(tracked) != string(got) || string(untracked) != string(scratch) {
		t.Fatalf("maps queries changed repository state: before=%q after=%q", before, after)
	}
}

func TestAXIStatusReportsUnknownForFailedActiveScan(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("decisions/model.md", activeMapDocument("shaping", "", "— (open)", "Resolved.", "Resolved."))
	f.WriteUnreadable("decisions", 0o755)

	status := f.Bench("status", "--all")
	status.RequireExit(0)
	requireContainsFold(t, status.Stdout, "unknown (decisions is unreadable)")
	count := f.Bench("maps", "--count")
	count.RequireExit(1)
	requireContainsFold(t, count.Stdout, "error: decisions is unreadable")
}

func activeMapDocument(status, fog, first, second, third string) string {
	return "# Model\n\nStatus: " + status + "\n\n## Destination\n\nSettle it.\n\n## #1: First\n\nBlocked by: none\nType: Research\n\n### Question\n\nFirst?\n\n### Answer\n\n" + first + "\n\n## #2: Second\n\nBlocked by: #1\nType: Task\n\n### Question\n\nSecond?\n\n### Answer\n\n" + second + "\n\n## #3: Third\n\nBlocked by: #1\nType: Prototype\n\n### Question\n\nThird?\n\n### Answer\n\n" + third + "\n\n## Not yet specified\n\n" + fog + "\n\n## Spec-writer discretion\n\n## Out of scope\n\n## Sources\n"
}
