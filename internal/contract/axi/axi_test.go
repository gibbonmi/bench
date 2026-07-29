package axi

import (
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAXIQuerySurfaceContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI learnings two-entry contract", testAXILearningsTwoEntry)
	contract.RunParallel(t, "AXI learnings empty/template contract", testAXILearningsEmptyTemplate)
	contract.RunParallel(t, "AXI maps unresolved-ticket contract", testAXIMapsUnresolvedTicket)
	contract.RunParallel(t, "AXI TOON field-escaping contract", testAXITOONFieldEscaping)
	contract.RunParallel(t, "AXI usage/exit-2 contract", testAXIUsageExit2)
	contract.RunParallel(t, "AXI subdirectory root-resolution contract", testAXISubdirectoryRootResolution)
	contract.RunParallel(t, "AXI path-with-spaces contract", testAXIPathWithSpaces)
	contract.RunParallel(t, "AXI maps over-match anchoring/fence contract", testAXIMapsOverMatchAnchoringFence)
	contract.RunParallel(t, "AXI maps CRLF-stripping contract", testAXIMapsCRLFStripping)
	contract.RunParallel(t, "AXI maps no-Type-ticket contract", testAXIMapsNoTypeTicket)
	contract.RunParallel(t, "AXI learnings ascii-separator title contract", testAXILearningsASCIISeparatorTitle)
	contract.RunParallel(t, "AXI maps handoff close-readiness contract", testAXIMapsHandoffCloseReadiness)
	contract.RunParallel(t, "AXI maps --count adapter contract", testAXIMapsCountAdapter)
}

func testAXILearningsTwoEntry(t *testing.T) {
	contract.NoteContractFailure(t, "AXI learnings two-entry contract failed")
	f := contract.NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 — first learning  [open]\n- body\n## 2026-02-02 — second learning  [open]\n- body\n")

	out := f.Bench("learnings")

	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "learnings[2]{date,title}:")
	requireIndentedRowCount(t, out.Stdout, 2)
	requireAXILine(t, out.Stdout, "  2026-01-01,first learning")
	requireAXILine(t, out.Stdout, "  2026-02-02,second learning")
}

func testAXILearningsEmptyTemplate(t *testing.T) {
	f := contract.NewFixture(t)

	absent := f.Bench("learnings")
	absent.RequireExit(0)
	requireAXILine(t, absent.Stdout, "learnings[0]{date,title}:")

	f.WriteFile(".bench/learnings.md", "## <date> — <short title>  [open]\n")
	template := f.Bench("learnings")
	template.RequireExit(0)
	requireAXILine(t, template.Stdout, "learnings[0]{date,title}:")

	f.WriteFile(".bench/learnings.md", "# Learnings — usage journal\n\nFormat per entry. Heading: `## YYYY-MM-DD — short title  [open]`\n")
	drained := f.Bench("learnings")
	drained.RequireExit(0)
	requireAXILine(t, drained.Stdout, "learnings[0]{date,title}:")
}

func testAXIMapsUnresolvedTicket(t *testing.T) {
	f := contract.NewFixture(t)

	absent := f.Bench("maps")
	absent.RequireExit(0)
	requireAXILine(t, absent.Stdout, "maps[0]{map,ticket,type,state}:")

	f.WriteFile("decisions/m.md", `# A map

## #1: first?

Type: Grill

### Answer
— (open)

## #2: second?

Type: Research

### Answer
— (deferred)

## #3: third?

Type: Grill

### Answer
GRILL DEFERRED — waiting on X

## #4: resolved?

Type: Grill

### Answer
Decided: yes, do the thing.
`)
	f.WriteFile("specs/x/decisions/compiled.md", "## #1: settled?\nType: Grill\n### Answer\n— (open)\n")
	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "maps[3]{map,ticket,type,state}:")
	requireAXILine(t, out.Stdout, "  m,1,Grill,open")
	requireAXILine(t, out.Stdout, "  m,2,Research,deferred")
	requireAXILine(t, out.Stdout, "  m,3,Grill,grill-deferred")
	requireNoAXILineMatching(t, out.Stdout, `^  m,4,`)
}

func testAXITOONFieldEscaping(t *testing.T) {
	contract.NoteContractFailure(t, "AXI TOON field-escaping contract failed")
	f := contract.NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-03-03 — a, \"b\"  [open]\n")

	out := f.Bench("learnings")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, `  2026-03-03,"a, \"b\""`)
}

func testAXIUsageExit2(t *testing.T) {
	f := contract.NewFixture(t)

	learnings := f.Bench("learnings", "bogusarg")
	learnings.RequireExit(2)
	requireContainsFold(t, learnings.Stdout, "usage")

	maps := f.Bench("maps", "bogusarg")
	maps.RequireExit(2)
	requireContainsFold(t, maps.Stdout, "usage")
}

func testAXISubdirectoryRootResolution(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-04-04 — sub check  [open]\n")
	subdir := filepath.Join(f.Root, "sub", "deeper")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatalf("create subdir: %v", err)
	}

	out := runBenchInDir(t, f, subdir, "learnings")

	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "learnings[1]{date,title}:")
	requireAXILine(t, out.Stdout, "  2026-04-04,sub check")
}

func testAXIPathWithSpaces(t *testing.T) {
	f := contract.NewFixture(t, contract.WithSpacePath())
	f.WriteFile(".bench/learnings.md", "## 2026-05-05 — spaced  [open]\n")
	f.WriteFile("decisions/s.md", `## #1: q?

Type: Grill

### Answer
— (open)
`)

	learnings := f.Bench("learnings")
	learnings.RequireExit(0)
	requireAXILine(t, learnings.Stdout, "  2026-05-05,spaced")

	maps := f.Bench("maps")
	maps.RequireExit(0)
	requireAXILine(t, maps.Stdout, "  s,1,Grill,open")
}

func testAXIMapsOverMatchAnchoringFence(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("decisions/o.md", `# Over-match map

## #1: genuine?

Type: Grill

### Answer
— (open)

## #2: prose mentions?

Type: Grill

### Answer
Decided: a mid-line GRILL DEFERRED mention is not an unresolved banner.

## #3: fenced example?

Type: Grill

### Answer
Decided: the placeholder looks like this:

`+"```"+`
— (open)
`+"```"+`

so authors recognize it.
`)

	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXIFirstLine(t, out.Stdout, "maps[1]{map,ticket,type,state}:")
	requireAXILine(t, out.Stdout, "  o,1,Grill,open")
	requireNoAXILineMatching(t, out.Stdout, `^  o,2,`)
	requireNoAXILineMatching(t, out.Stdout, `^  o,3,`)
}

func testAXIMapsCRLFStripping(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("decisions/c.md", "## #1: q?\r\nType: Grill\r\n### Answer\r\n— (open)\r\n")

	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  c,1,Grill,open")
	if strings.Contains(out.Stdout, "\r") {
		t.Fatalf("CRLF leaked carriage returns into maps output:\n%s", out.Stdout)
	}
}

func testAXIMapsNoTypeTicket(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("decisions/n.md", `## #1: typeless?

### Answer
— (open)
`)

	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  n,1,unknown,open")
}

func testAXILearningsASCIISeparatorTitle(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile(".bench/learnings.md", "## 2026-01-01 - ascii title  [open]\n")

	out := f.Bench("learnings")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  2026-01-01,ascii title")
}

func testAXIMapsHandoffCloseReadiness(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("decisions/hm.md", `# HM
## #1: q?
Type: Grill
### Answer
Decided: yes.
`)
	f.WriteFile("decisions/hf.md", `# HF
## #1: q?
Type: Grill
### Answer
Decided: yes.

## Handoff
1. Module boundaries. n/a — single unit.
2. Contracts. n/a — no CLI surface.
`)
	f.WriteFile("decisions/ho.md", `# HO
## #1: q?
Type: Grill
### Answer
— (open)
`)
	f.WriteFile("decisions/hx.md", `# HX
## #1: q?
Type: Grill
### Answer
Decided: yes. A handoff section looks like:

`+"```"+`
## Handoff
1. Module boundaries.
`+"```"+`
`)
	f.WriteFile("decisions/hp.md", `# HP
## #1: q?
Type: Grill
### Answer
Decided: yes.

## Handoff
1. Module boundaries.
— (open)
`)
	// A README documents the directory rather than claiming membership in it, so it is
	// not a decision map and never reaches the parser: it earns no row, adds nothing to
	// the unresolved tally, and cannot flip the command's exit code. Every other file
	// here is close-ready or carries only a handoff row, so the exit is 0.
	f.WriteFile("decisions/README.md", "# Decisions index\nNotes about this folder — not a map.\n")

	out := f.Bench("maps")

	out.RequireExit(0)
	requireAXILine(t, out.Stdout, "  hm,handoff,handoff,missing")
	requireAXILine(t, out.Stdout, "  hx,handoff,handoff,missing")
	requireAXILine(t, out.Stdout, "  hp,handoff,handoff,open")
	requireNoAXILineMatching(t, out.Stdout, `^  hf,handoff`)
	requireNoAXILineMatching(t, out.Stdout, `^  ho,handoff`)
	requireAXILine(t, out.Stdout, "  ho,1,Grill,open")
	requireNoAXILineMatching(t, out.Stdout, `^  hp,1,`)
	requireNoAXILineMatching(t, out.Stdout, `^  README,`)
}

func testAXIMapsCountAdapter(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("decisions/m.md", "## #1: a?\nType: Grill\n### Answer\n— (open)\n")
	f.WriteFile("decisions/scope.md", "### Answer\n— (deferred)\n")
	f.WriteFile("decisions/.hidden.md", "## #1: h?\nType: Grill\n### Answer\n— (open)\n")

	count := f.Bench("maps", "--count")
	count.RequireExit(0)
	if got := strings.TrimSpace(count.Stdout); got != "2" {
		t.Fatalf("maps --count = %q, want 2\nstdout:\n%s\nstderr:\n%s", got, count.Stdout, count.Stderr)
	}

	rows := f.Bench("maps")
	rows.RequireExit(0)
	requireAXILine(t, rows.Stdout, "  m,1,Grill,open")
	requireNoAXILineMatching(t, rows.Stdout, `^  scope,`)
	requireNoAXILineMatching(t, rows.Stdout, `^  \.?hidden,`)
}
