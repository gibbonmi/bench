package runtime

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/handoff"
	"github.com/gibbonmi/bench/internal/status"
)

// handoffWithState builds a fixture whose handoff is present but untracked, so the file is
// one dirty path before the command runs and one after.
//
// That is a real limit on what the idempotence row can prove, not a convenience. The
// command's write is itself a change to the tree it reports on, so on a *tracked* clean
// repo the second run legitimately sees a dirty path the first did not, and the board grows
// a git signal the first run's derivation never had. The pin block's own dirty clause
// excludes the handoff for that reason (see landedState) and is stable in both
// configurations; the board-derived Next command is not, and making it stable needs an
// exclusion `bench status` does not offer. Until it does, this row pins the sinks against
// timestamps and reordering, and the tracked-tree difference is recorded as out of scope in
// the spec rather than left to look like a passing guarantee.
func handoffWithState(t *testing.T, body string) contract.Fixture {
	t.Helper()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	f.WriteFile(status.HandoffFile, "# Session handoff\n\nstale header line\n\n## State\n\n"+body+"\n\n## Next command\n\n`stale command`\n\n## Shape\n\nstale shape text\n")
	return f
}

func TestHandoffPreservesState(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	// Distinctive prose with interior blank lines, indentation, and markdown a reflowing
	// writer would normalize.
	body := "- **The canary deadlocks** rather than failing.\n\n      indented block\n\n- second point with a trailing double space  "
	f := handoffWithState(t, body)

	f.Bench("handoff").RequireExit(0)
	written := f.ReadFile(status.HandoffFile)
	if !strings.Contains(written, "## State\n\n"+body+"\n\n") {
		t.Fatalf("State section was not preserved byte-for-byte:\n%s", written)
	}
	contract.RequireNotContains(t, written, "stale header line")
}

// TestHandoffPreservesStateHeadings covers the section the reviewer owns carrying level-two
// headings of its own. The splitter ends the body at a generated heading, never at any
// heading, because a reviewer who organizes their own notes must not have everything below
// the first one silently dropped on a zero exit.
func TestHandoffPreservesStateHeadings(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	body := "judgment above the heading\n\n## Notes\n\njudgment below it\n\n## Open questions\n\nand below that"
	f := handoffWithState(t, body)

	f.Bench("handoff").RequireExit(0)
	written := f.ReadFile(status.HandoffFile)
	if !strings.Contains(written, "## State\n\n"+body+"\n\n") {
		t.Fatalf("a heading inside State truncated the section:\n%s", written)
	}
	contract.RequireContains(t, written, "judgment below it")
	contract.RequireContains(t, written, "and below that")
}

// TestHandoffEmptyDocument covers the present-but-empty file, which is a distinct state
// from an absent one: absent is scaffolded, empty is refused. Collapsing empty into the
// scaffold branch would overwrite a file whose State section a reader cannot confirm was
// ever there.
func TestHandoffEmptyDocument(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.WriteFile(status.HandoffFile, "")
	f.CommitAll("empty handoff")

	out := f.Bench("handoff")
	out.RequireExit(1)
	if written := f.ReadFile(status.HandoffFile); written != "" {
		t.Fatalf("refusal wrote to the file:\n%s", written)
	}
}

func TestHandoffRegeneratesShape(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffWithState(t, "carried judgment")

	f.Bench("handoff").RequireExit(0)
	written := f.ReadFile(status.HandoffFile)
	contract.RequireContains(t, written, handoff.ShapeSection)
	contract.RequireNotContains(t, written, "stale shape text")
	if n := strings.Count(written, "## Shape"); n != 1 {
		t.Fatalf("want exactly one ## Shape heading, got %d:\n%s", n, written)
	}
	if n := strings.Count(written, "## Next command"); n != 1 {
		t.Fatalf("want exactly one ## Next command heading, got %d:\n%s", n, written)
	}
}

func TestHandoffSkeletonCarriesConventions(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")

	f.Bench("handoff").RequireExit(0)
	scaffolded := f.ReadFile(status.HandoffFile)
	// The guidance is unexported, so these expectations are authored independently of it.
	// That independence is what makes an implementation scaffolding bare structure — the
	// headings with nothing telling the first session what the document is for — red here.
	for _, phrase := range []string{"cold-start pin", "The State section is yours", "regenerated from the git tree"} {
		contract.RequireContains(t, scaffolded, phrase)
	}

	// The guidance is addressed to a session facing an empty State, so it goes once that
	// session has written one.
	written := handoffWithState(t, "carried judgment")
	written.Bench("handoff").RequireExit(0)
	contract.RequireNotContains(t, written.ReadFile(status.HandoffFile), "cold-start pin")
}

func TestHandoffScaffoldsMissing(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	if f.Exists(status.HandoffFile) {
		t.Fatal("fixture already carries a handoff")
	}

	out := f.Bench("handoff")
	out.RequireExit(0)
	written := f.ReadFile(status.HandoffFile)
	contract.RequireContains(t, written, "## State\n\n## Next command")
}

func TestHandoffRefusesUnparseable(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	before := "# Session handoff\n\nhand-edited, and the section heading is gone\n\n## Next command\n\n`bench status`\n"
	f.WriteFile(status.HandoffFile, before)

	out := f.Bench("handoff")
	out.RequireExit(1)
	contract.RequireContains(t, out.Stdout, "## State")
	contract.RequireContains(t, out.Stdout, "error: ")
	if after := f.ReadFile(status.HandoffFile); after != before {
		t.Fatalf("refusal rewrote the file:\n%s", after)
	}
}

func TestHandoffRefusesAmbiguousState(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	cases := map[string]string{
		// A fence-blind splitter finds exactly one heading here and writes over the
		// document; the fence means the section is not really present at all.
		"the only heading is fenced": "# Session handoff\n\n```\n## State\n\nan example, not a section\n```\n",
		"two headings":               "# Session handoff\n\n## State\n\nfirst\n\n## State\n\nsecond\n",
	}
	for name, before := range cases {
		t.Run(name, func(t *testing.T) {
			f := handoffFixtureOnMain(t)
			f.WriteFile("tracked.txt", "base\n")
			f.CommitAll("base")
			f.WriteFile(status.HandoffFile, before)

			out := f.Bench("handoff")
			out.RequireExit(1)
			contract.RequireContains(t, out.Stdout, "## State")
			if after := f.ReadFile(status.HandoffFile); after != before {
				t.Fatalf("refusal rewrote the file:\n%s", after)
			}
		})
	}
}

func TestHandoffIdempotent(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffWithState(t, "judgment that must survive a second run")

	first := f.Bench("handoff")
	first.RequireExit(0)
	firstFile := f.ReadFile(status.HandoffFile)

	second := f.Bench("handoff")
	second.RequireExit(0)
	if second.Stdout != first.Stdout {
		t.Fatalf("stdout changed between runs\nfirst:\n%s\nsecond:\n%s", first.Stdout, second.Stdout)
	}
	if secondFile := f.ReadFile(status.HandoffFile); secondFile != firstFile {
		t.Fatalf("file changed between runs\nfirst:\n%s\nsecond:\n%s", firstFile, secondFile)
	}
}
