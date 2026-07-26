package axi

import (
	"github.com/gibbonmi/bench/internal/contract"
	"path/filepath"
	"strings"
	"testing"
)

// axiCoverageStories and axiCoverageHeader are the fixture bodies every coverage
// --check contract test builds a spec from, so the declared story set (1, 2, 3)
// and the canonical header are asserted from one source rather than re-typed per test.
// axiCoverageGappedStories declares 1, 2, 4 — a set with a hole below its maximum, so a
// row referencing story 3 is a number a max-only check would let through but exact-set
// membership must reject.
const (
	axiCoverageStories       = "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n3. As g, I want h, so i.\n"
	axiCoverageGappedStories = "## User stories\n1. As a, I want b, so c.\n2. As d, I want e, so f.\n4. As g, I want h, so i.\n"
	axiCoverageHeader        = "| story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|"
)

func TestAXICoverageContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI coverage extraction contract", testAXICoverageExtraction)
	contract.RunParallel(t, "AXI coverage state/error contract", testAXICoverageStateError)
	contract.RunParallel(t, "AXI coverage --check validation contract", testAXICoverageCheckValidation)
	contract.RunParallel(t, "AXI coverage root-anchored slug contract", testAXICoverageSlugRootAnchored)
}

// testAXICoverageSlugRootAnchored pins the slug fallback to the repository root. It
// compares the subdirectory run against the root run rather than asserting the
// subdirectory run succeeds: an implementation that resolved both against the process
// cwd would still pass a bare success assertion when the root run happens to work.
func testAXICoverageSlugRootAnchored(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/anchored.md", "# anchored\n\n## User stories\n1. As a, I want b, so c.\n\n### Acceptance coverage map\n| story | behavior | seam | red signal | why it catches the failure |\n|---|---|---|---|---|\n| 1 | b | s | r | w |\n")
	subdir := filepath.Join(f.Root, "sub", "deeper")
	contract.Mkdir(t, subdir)

	root := f.Bench("coverage", "anchored")
	nested := runBenchInDir(t, f, subdir, "coverage", "anchored")
	explicit := f.Bench("coverage", "specs/anchored.md")

	if nested.ExitCode != root.ExitCode {
		t.Fatalf("exit code from subdirectory = %d, from root = %d\nsubdir stdout:\n%s\nroot stdout:\n%s", nested.ExitCode, root.ExitCode, nested.Stdout, root.Stdout)
	}
	if nested.Stdout != root.Stdout {
		t.Fatalf("stdout differs by cwd\nsubdir:\n%s\nroot:\n%s", nested.Stdout, root.Stdout)
	}
	// The slug and explicit-path forms name the same file, so they must render it the
	// same way. An absolute rendering would satisfy the cwd comparison above while still
	// leaking the machine's layout and differing from the path form.
	if explicit.Stdout != root.Stdout {
		t.Fatalf("slug form differs from explicit-path form\nslug:\n%s\npath:\n%s", root.Stdout, explicit.Stdout)
	}
	root.RequireExit(0)
	requireOutputLine(t, root, "spec: specs/anchored.md")
	requireOutputLine(t, root, "state: mapped")
}

func testAXICoverageExtraction(t *testing.T) {
	contract.NoteContractFailure(t, "AXI coverage extraction contract failed")
	f := contract.NewFixture(t)
	spec := `# t

## User stories
1. As a, I want b, so c.
2. As d, I want e, so f.
3. As g, I want h, so i.

## Testing decisions

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 2–3 | does x \| y | cli seam | cmd fails, loudly | catches z |
| edge of 1 | edge case | gate | already covered | catches w |

## Out of scope
`
	f.WriteFile("specs/t.md", spec)

	out := f.Bench("coverage", "specs/t.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "spec: specs/t.md")
	requireOutputLine(t, out, "state: mapped")
	requireOutputLine(t, out, "rows[2]{story,seam,red_signal}:")
	requireOutputLine(t, out, `  2–3,cli seam,"cmd fails, loudly"`)
	requireOutputLine(t, out, "  edge of 1,gate,already covered")

	f.WriteFile("specs/t.md", strings.ReplaceAll(spec, "\n", "\r\n"))
	out = f.Bench("coverage", "specs/t.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "state: mapped")
	requireOutputLine(t, out, "rows[2]{story,seam,red_signal}:")
}

func testAXICoverageStateError(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("specs/h.md", "# h\n\n<!-- coverage-map: historical -->\n\n### Acceptance coverage map\n| a |\n")

	out := f.Bench("coverage", "specs/h.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "state: historical")
	requireOutputLine(t, out, "rows[0]{story,seam,red_signal}:")

	f.WriteFile("specs/n.md", "# n\nno map here\n")
	out = f.Bench("coverage", "specs/n.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "state: no-map")

	out = f.Bench("coverage")

	out.RequireExit(2)
	out.RequireContains(strings.ToLower(out.Stdout), "usage")

	out = f.Bench("coverage", "--bogus", "specs/n.md")

	out.RequireExit(2)

	out = f.Bench("coverage", "specs/absent.md")

	out.RequireExit(1)
	requireOutputPrefix(t, out, "error: ")
}

func testAXICoverageCheckValidation(t *testing.T) {
	f := contract.NewFixture(t)
	stories := axiCoverageStories
	header := axiCoverageHeader
	f.WriteFile("specs/v.md", "# v\n\n"+stories+"\n### Acceptance coverage map\n"+header+"\n| 1, 2–3 | b | s | r | w |\n")

	out := f.Bench("coverage", "--check", "specs/v.md")

	// row 2: a valid map's --check gets a definitive pass line, not silence —
	// stdout was empty before this story, indistinguishable from a check that
	// silently produced nothing.
	out.RequireExit(0)
	requireOutputLine(t, out, "ok: coverage map valid — 1 row(s)")

	// A historical spec's pass is also a definitive line, not silence — see
	// TestAXICoverageHistoricalPasses for the pinned assertion.
	f.WriteFile("specs/h.md", "# h\n<!-- coverage-map: historical -->\n### Acceptance coverage map\n|bad|\n")
	out = f.Bench("coverage", "--check", "specs/h.md")
	out.RequireExit(0)
	requireOutputLine(t, out, "ok: coverage map historical — validation skipped")

	// An unmapped spec with no historical marker now fails --check outright —
	// see TestAXICoverageNoMapFails / TestAXICoverageHistoricalPasses for the pair
	// that pins this rule and its opt-out.
	f.WriteFile("specs/n.md", "# n\nprose only\n")
	out = f.Bench("coverage", "--check", "specs/n.md")
	out.RequireExit(1)
	out.RequireContains(out.Stdout, "coverage map missing")

	cases := []struct {
		name string
		path string
		body string
		want string
	}{
		{
			name: "missing canonical header",
			path: "specs/b1.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n",
			want: "coverage map missing the canonical header",
		},
		{
			name: "no data rows",
			path: "specs/b2.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n",
			want: "coverage map has no data rows",
		},
		{
			name: "wrong cell count",
			path: "specs/b3.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| 1 | b | s | r |\n",
			want: "coverage map row 1 has 4 cells (want 5)",
		},
		{
			name: "empty cell",
			path: "specs/b4.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| 1 | b |  | r | w |\n",
			want: "coverage map row 1 has an empty 'seam' cell",
		},
		{
			name: "story reference beyond numbering",
			path: "specs/b5.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| 9 | b | s | r | w |\n",
			want: "references story 9, which the spec does not declare (has: 1, 2, 3)",
		},
		{
			name: "unrecognized story reference",
			path: "specs/b6.md",
			body: "# b\n\n" + stories + "\n### Acceptance coverage map\n" + header + "\n| x | b | s | r | w |\n",
			want: "has an unrecognized story reference 'x'",
		},
		{
			// A spec with a map but no story section declares nothing, so the message
			// must name that state outright — the set-listing phrasing degenerates into
			// a dangling "(has: )" here and reads like a truncated line.
			name: "story reference with no declared stories",
			path: "specs/b7.md",
			body: "# b\n\n### Acceptance coverage map\n" + header + "\n| 1 | b | s | r | w |\n",
			want: "references story 1, but the spec declares no stories",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			f := contract.NewFixture(t)
			f.WriteFile(c.path, c.body)

			out := f.Bench("coverage", "--check", c.path)

			out.RequireExit(1)
			out.RequireContains(out.Stdout, c.want)
		})
	}

	// row 3: a violation line renders through the one canonical `toon.Errorf`
	// shape (`error: <kind> — <hint>`), not a hand-rolled line of its own —
	// pinned byte-for-byte rather than by substring.
	f.WriteFile("specs/b1.md", "# b\n\n"+stories+"\n### Acceptance coverage map\n| a | b |\n|---|---|\n| 1 | x |\n")
	out = f.Bench("coverage", "--check", "specs/b1.md")
	out.RequireExit(1)
	requireOutputLine(t, out, "error: specs/b1.md coverage map missing the canonical header — fix the map or mark it <!-- coverage-map: historical -->")
}

// TestAXICoverageNoMapFails pins story 14: a spec with no acceptance coverage map and
// no historical marker fails --check. The current implementation treats "nothing to
// validate" as a pass, which this exit-1 assertion rejects.
func TestAXICoverageNoMapFails(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("specs/n.md", "# n\nprose only, no map here\n")

	out := f.Bench("coverage", "--check", "specs/n.md")

	out.RequireExit(1)
	out.RequireContains(out.Stdout, "coverage map missing")
}

// TestAXICoverageHistoricalPasses pairs with TestAXICoverageNoMapFails: a spec that
// carries the historical marker but still has no map passes --check. The pass is a
// definitive one-line result, not silence — a silent empty state is exactly the
// defect class this spec exists to eliminate, and the historical opt-out is new code
// that would otherwise reintroduce it.
func TestAXICoverageHistoricalPasses(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("specs/h.md", "# h\n\n<!-- coverage-map: historical -->\n\nno map here either\n")

	out := f.Bench("coverage", "--check", "specs/h.md")

	out.RequireExit(0)
	requireOutputLine(t, out, "ok: coverage map historical — validation skipped")
}

// TestAXICoverageStoryMembership pins story 15: story 0, a reference outside the
// declared story set, and a reversed range each fail --check with their own message.
// The non-member case uses axiCoverageGappedStories (1, 2, 4) and references story 3 —
// a number below the maximum but inside the gap — so the fixture is discriminating: a
// max-only validator (the prior implementation) would let it through, while only exact
// set membership rejects it. The range case reuses the same gapped set to reach the
// number a range spans without naming: `2-4` has declared endpoints, so only a check
// that walks the numbers between them rejects it.
func TestAXICoverageStoryMembership(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	cases := []struct {
		name, path, stories, storyCell, want string
	}{
		{"story zero", "specs/m0.md", axiCoverageStories, "0", "references story 0, which is not a valid story number"},
		{"non-member story", "specs/m1.md", axiCoverageGappedStories, "3", "references story 3, which the spec does not declare (has: 1, 2, 4)"},
		{"reversed range", "specs/m2.md", axiCoverageStories, "3-1", "has a story range with end before start '3-1'"},
		{"non-member inside a range", "specs/m3.md", axiCoverageGappedStories, "2-4", "references story 3, which the spec does not declare (has: 1, 2, 4)"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			f := contract.NewFixture(t)
			f.WriteFile(c.path, "# b\n\n"+c.stories+"\n### Acceptance coverage map\n"+axiCoverageHeader+
				"\n| "+c.storyCell+" | b | s | r | w |\n")

			out := f.Bench("coverage", "--check", c.path)

			out.RequireExit(1)
			out.RequireContains(out.Stdout, c.want)
		})
	}
}

// TestAXICoverageValidStoryRefs pairs with TestAXICoverageStoryMembership: a valid
// comma list and a valid forward range still pass, forbidding a validator so strict it
// rejects legitimate multi-story rows.
func TestAXICoverageValidStoryRefs(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t)
	f.WriteFile("specs/valid.md", "# v\n\n"+axiCoverageStories+"\n### Acceptance coverage map\n"+axiCoverageHeader+
		"\n| 1, 2-3 | b | s | r | w |\n")

	out := f.Bench("coverage", "--check", "specs/valid.md")

	out.RequireExit(0)
	out.RequireContains(out.Stdout, "ok: coverage map valid")
}
