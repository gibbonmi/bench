package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

const anchorClosureLead = "Writes: entry names an anchored guidance path without naming every anchor registry file: "

// anchoredFacts is baseFacts with one ticket that writes each guidance path, and an
// anchor map that names files for those paths.
func anchoredFacts(anchored map[string][]string, writes ...string) Facts {
	f := baseFacts()
	f.Tickets[0].Writes = append(f.Tickets[0].Writes, writes...)
	f.WritesAnchorFiles = anchored
	return f
}

func assertAnchorClosure(t *testing.T, f Facts, want string) {
	t.Helper()
	c, ok := checkRow(Decide(f), "anchor-closure")
	if want == "" {
		if !ok || c.Verdict != verdictGreen {
			t.Fatalf("anchor-closure row = %+v, want green", c)
		}
		return
	}
	if !ok || c.Verdict != verdictRed || c.Detail != anchorClosureLead+want {
		t.Fatalf("anchor-closure row = %+v, want the red detailed %q", c, anchorClosureLead+want)
	}
}

// TestAnchorClosureRedsUnnamedRegistry covers SC1. A ticket that writes an anchored
// guidance path and omits the registry file that names it reds, a spaced path too.
func TestAnchorClosureRedsUnnamedRegistry(t *testing.T) {
	assertAnchorClosure(t, anchoredFacts(map[string][]string{
		".agents/x/SKILL.md": {"internal/anchors/registry_a.go"},
	}, ".agents/x/SKILL.md"), "one.md: .agents/x/SKILL.md is anchored by internal/anchors/registry_a.go")
	assertAnchorClosure(t, anchoredFacts(map[string][]string{
		".agents/spaced dir/SKILL.md": {"internal/anchors/registry_a.go"},
	}, ".agents/spaced dir/SKILL.md"), "one.md: .agents/spaced dir/SKILL.md is anchored by internal/anchors/registry_a.go")
}

// TestAnchorClosureNamesEachMissingFile covers SC2. The detail lists each missing file
// once, and it leaves out a file the ticket already names.
func TestAnchorClosureNamesEachMissingFile(t *testing.T) {
	f := anchoredFacts(map[string][]string{
		".agents/x/SKILL.md": {"internal/anchors/registry_a.go", "internal/anchors/registry_a_test.go", "internal/anchors/registry_b.go"},
	}, ".agents/x/SKILL.md", "internal/anchors/registry_b.go")
	assertAnchorClosure(t, f, "one.md: .agents/x/SKILL.md is anchored by internal/anchors/registry_a.go, "+
		"one.md: .agents/x/SKILL.md is anchored by internal/anchors/registry_a_test.go")
}

// TestAnchorClosureGreenWhenNamed covers SC7. Each named registry file, or the
// anchor registry directory, closes the row.
func TestAnchorClosureGreenWhenNamed(t *testing.T) {
	anchored := map[string][]string{".agents/x/SKILL.md": {"internal/anchors/registry_a.go", "internal/anchors/registry_a_test.go"}}
	assertAnchorClosure(t, anchoredFacts(anchored, ".agents/x/SKILL.md", "internal/anchors/registry_a.go", "internal/anchors/registry_a_test.go (new)"), "")
	assertAnchorClosure(t, anchoredFacts(anchored, ".agents/x/SKILL.md", "internal/anchors"), "")
}

// seedAnchoredTicket plants a repository whose own anchor registry file names the
// guidance path .agents/x/SKILL.md, and whose one ticket writes writes.
func seedAnchoredTicket(t *testing.T, writes string) (root, slug string) {
	t.Helper()
	slug = "example"
	root = preflighttest.StartRepo(t)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug, "- `.agents/x/`", "- `internal/anchors/`"))
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md",
		strings.Replace(preflighttest.TicketDoc("One", "PF1", "PF2"), "Writes: specs", "Writes: "+writes, 1))
	preflighttest.MustWriteFile(t, ".agents/x/SKILL.md", "# X\n")
	preflighttest.MustWriteFile(t, "internal/anchors/extra.go", "package anchors\n\nvar extra = []string{`.agents/x/SKILL.md`}\n")
	for _, source := range []string{chargesource.DelegateSkill, chargesource.DelegateProcedure, chargesource.BuildPhase} {
		preflighttest.MustWriteFile(t, source, "# Source\n")
	}
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	preflighttest.MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	preflighttest.RunGit(t, "add", "internal/"+slug+"/foo.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}

func assertAnchorClosureRow(t *testing.T, out, detail string) {
	t.Helper()
	want := "  anchor-closure,red,\"" + anchorClosureLead + detail + "\",\"\"\n"
	if !strings.Contains(out, want) {
		t.Fatalf("output missing the anchor-closure row %q:\n%s", want, out)
	}
}

// TestCommandBuildAnchorClosureReadsTree covers SC12. The row reads the graded tree's
// own anchor registry, so a file no binary knows reds the row.
func TestCommandBuildAnchorClosureReadsTree(t *testing.T) {
	_, slug := seedAnchoredTicket(t, ".agents/x/SKILL.md")
	out, code := Command([]string{"build", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	assertAnchorClosureRow(t, out, "one.md: .agents/x/SKILL.md is anchored by internal/anchors/extra.go")
}

// TestAnchorClosureCoversDirectoryEntry covers SC3. A directory entry takes the
// closure of the anchored file under it.
func TestAnchorClosureCoversDirectoryEntry(t *testing.T) {
	_, slug := seedAnchoredTicket(t, ".agents/x")
	out, _ := Command([]string{"build", slug})
	assertAnchorClosureRow(t, out, "one.md: .agents/x is anchored by internal/anchors/extra.go")
}

// TestProposeWritesListsAnchorClosure covers SC9. The proposal lists the missing
// anchor file with its source and fence state, and the anchor-closure red does not
// refuse it.
func TestProposeWritesListsAnchorClosure(t *testing.T) {
	root, slug := seedAnchoredTicket(t, ".agents/x/SKILL.md")
	out, code := Command(proposalArgs(t, root, slug))
	if code != 0 || !strings.Contains(out, "writes_proposal[1]{path,source,fence}") {
		t.Fatalf("proposal = (%d):\n%s", code, out)
	}
	assertProposalRow(t, out, "internal/anchors/extra.go", "anchor .agents/x/SKILL.md", "covered")
}
