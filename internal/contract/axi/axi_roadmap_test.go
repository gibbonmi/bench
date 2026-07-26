package axi

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestAXIRoadmapUnsupportedSchema pins story 7's roadmap half: a ROADMAP.md whose
// bytes read cleanly but carry no recognizable `**ID**` row reports the parser's
// own unsupported-schema state rather than printing as a working document. The
// fixture is written without a trailing newline — the hand-edited-file class this
// edge case is assigned to.
func TestAXIRoadmapUnsupportedSchema(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("ROADMAP.md", "# Roadmap\n\njust prose, no roadmap rows here")

	out := f.Bench("roadmap")

	out.RequireExit(1)
	out.RequireContains(out.Stdout, "error: ROADMAP.md is unsupported-schema — no roadmap rows recognized")
}

// TestAXIRoadmapFailClosed pins story 12: an unreadable ROADMAP.md exits 1 with a
// structured error naming the state, never a silent empty-document print.
func TestAXIRoadmapFailClosed(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n**FT1 — x.** body.\n")
	f.WriteUnreadable("ROADMAP.md", 0o644)

	out := f.Bench("roadmap")

	out.RequireExit(1)
	out.RequireContains(out.Stdout, "error: ROADMAP.md is unreadable")
}

// TestAXIRoadmapAbsentIsEmpty is TestAXIRoadmapFailClosed's pair: no ROADMAP.md at
// all exits 0 and renders the maintenance-prompt empty state, forbidding an
// always-exit-1 stub from satisfying the fail-closed row above.
func TestAXIRoadmapAbsentIsEmpty(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)

	out := f.Bench("roadmap")

	out.RequireExit(0)
	out.RequireContains(out.Stdout, "no ROADMAP.md")
	out.RequireContains(out.Stdout, "/bench-what-next")
}

// TestAXIRoadmapEmptyIsNotAbsent asserts the two states in one run, because the bug they
// forbid is their collapse: a zero-byte ROADMAP.md is a file someone created and left
// unwritten, not a repository that never had one, and only absence is the authoritative
// empty state. The empty file therefore takes the same non-absent posture every other
// query command takes — exit 1 naming the state — while the absent case keeps the
// maintenance prompt on exit 0.
func TestAXIRoadmapEmptyIsNotAbsent(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("ROADMAP.md", "")

	empty := f.Bench("roadmap")

	empty.RequireExit(1)
	empty.RequireContains(empty.Stdout, "error: ROADMAP.md is empty")
	if strings.Contains(empty.Stdout, "no ROADMAP.md") {
		t.Fatalf("a present-but-empty ROADMAP.md rendered the absent-file prompt:\n%s", empty.Stdout)
	}
}

// TestAXIRoadmapContextDegrades pins story 13: an unreadable IDEAS.md costs
// `roadmap --context` only that source's state, not the whole snapshot — the
// roadmap rows and learnings blocks from unrelated, readable sources still render,
// and the command still exits 0.
func TestAXIRoadmapContextDegrades(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.Git("branch", "-M", "main")
	f.WriteFile("ROADMAP.md", "# Roadmap\n\n**FT1 — one.** Body.\n\n## Recommended sequence\n\n1. `/bench-implement-spec` — one\n")
	f.WriteFile("IDEAS.md", "- 2026-07-10  retain me\n")
	f.WriteFile(".bench/learnings.md", "## 2026-07-10 — lesson  [open]\n- body\n")
	f.WriteFile(".bench/structure.budgets", "")
	f.WriteFile(".bench/structure-accept", "")
	f.Git("add", "-A")
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "fixture")
	f.WriteUnreadable("IDEAS.md", 0o644)

	out := f.Bench("roadmap", "--context")

	out.RequireExit(0)
	out.RequireContains(out.Stdout, "IDEAS.md,unreadable")
	out.RequireContains(out.Stdout, "roadmap_rows[1]{")
	out.RequireContains(out.Stdout, "learnings[1]{")
}
