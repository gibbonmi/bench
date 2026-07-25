package runtime

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/handoff"
)

// handoffWithState builds a fixture whose handoff is present but untracked. Untracked is
// deliberate: the file is one dirty path before the command runs and one after, so the
// derived facts do not move underneath a second invocation.
func handoffWithState(t *testing.T, body string) contract.Fixture {
	t.Helper()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	f.WriteFile(handoffFile, "# Session handoff\n\nstale header line\n\n## State\n\n"+body+"\n\n## Next command\n\n`stale command`\n\n## Shape\n\nstale shape text\n")
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
	written := f.ReadFile(handoffFile)
	if !strings.Contains(written, "## State\n\n"+body+"\n\n") {
		t.Fatalf("State section was not preserved byte-for-byte:\n%s", written)
	}
	contract.RequireNotContains(t, written, "stale header line")
}

func TestHandoffRegeneratesShape(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffWithState(t, "carried judgment")

	f.Bench("handoff").RequireExit(0)
	written := f.ReadFile(handoffFile)
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
	scaffolded := f.ReadFile(handoffFile)
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
	contract.RequireNotContains(t, written.ReadFile(handoffFile), "cold-start pin")
}

func TestHandoffScaffoldsMissing(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	if f.Exists(handoffFile) {
		t.Fatal("fixture already carries a handoff")
	}

	out := f.Bench("handoff")
	out.RequireExit(0)
	written := f.ReadFile(handoffFile)
	contract.RequireContains(t, written, "## State\n\n## Next command")
}

func TestHandoffRefusesUnparseable(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	before := "# Session handoff\n\nhand-edited, and the section heading is gone\n\n## Next command\n\n`bench status`\n"
	f.WriteFile(handoffFile, before)

	out := f.Bench("handoff")
	out.RequireExit(1)
	contract.RequireContains(t, out.Stdout, "## State")
	contract.RequireContains(t, out.Stdout, "error: ")
	if after := f.ReadFile(handoffFile); after != before {
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
			f.WriteFile(handoffFile, before)

			out := f.Bench("handoff")
			out.RequireExit(1)
			contract.RequireContains(t, out.Stdout, "## State")
			if after := f.ReadFile(handoffFile); after != before {
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
	firstFile := f.ReadFile(handoffFile)

	second := f.Bench("handoff")
	second.RequireExit(0)
	if second.Stdout != first.Stdout {
		t.Fatalf("stdout changed between runs\nfirst:\n%s\nsecond:\n%s", first.Stdout, second.Stdout)
	}
	if secondFile := f.ReadFile(handoffFile); secondFile != firstFile {
		t.Fatalf("file changed between runs\nfirst:\n%s\nsecond:\n%s", firstFile, secondFile)
	}
}
