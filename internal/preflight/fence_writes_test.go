package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/tickets"
)

const fenceWritesLead = "spec fence and ticket Writes: union differ: "

// fencedFacts is baseFacts with the given fence and one ticket that writes writes.
func fencedFacts(fence []string, writes ...string) Facts {
	f := baseFacts()
	f.FenceEntries = fence
	f.Tickets = []tickets.Ticket{{Name: "one.md", Writes: writes, Covers: []string{"PF1", "PF2"}}}
	return f
}

func assertFenceWrites(t *testing.T, f Facts, want string) {
	t.Helper()
	c, ok := checkRow(Decide(f), "fence-writes")
	if want == "" {
		if !ok || c.Verdict != verdictGreen {
			t.Fatalf("fence-writes row = %+v, want green", c)
		}
		return
	}
	if !ok || c.Verdict != verdictRed || c.Detail != fenceWritesLead+want {
		t.Fatalf("fence-writes row = %+v, want the red detailed %q", c, fenceWritesLead+want)
	}
}

// TestFenceWritesRedsDifferentSets covers SC14. Each containment direction reds on
// its own.
func TestFenceWritesRedsDifferentSets(t *testing.T) {
	assertFenceWrites(t, fencedFacts([]string{"a/"}, "a", "b.go"), "Writes only: b.go")
	assertFenceWrites(t, fencedFacts([]string{"a/", "b.go"}, "b.go"), "fence only: a")
}

// TestFenceWritesNamesBothSides covers SC15. The detail names each side apart, in
// fence-then-Writes order, and sorts each list.
func TestFenceWritesNamesBothSides(t *testing.T) {
	assertFenceWrites(t, fencedFacts([]string{"a/"}, "b.go"), "fence only: a; Writes only: b.go")
	assertFenceWrites(t, fencedFacts([]string{"z", "a/"}, "y.go", "b.go"), "fence only: a, z; Writes only: b.go, y.go")
}

// TestFenceWritesGreenWhenExact covers SC16.
func TestFenceWritesGreenWhenExact(t *testing.T) {
	assertFenceWrites(t, fencedFacts([]string{"a/", "b.go"}, "b.go", "a"), "")
}

// TestFenceWritesIgnoresReviewPickup covers SC17. The pickup the review record path
// owner names for specs/example/spec.md is on neither side.
func TestFenceWritesIgnoresReviewPickup(t *testing.T) {
	assertFenceWrites(t, fencedFacts([]string{"a/", "reviews/example.md"}, "a"), "")
	assertFenceWrites(t, fencedFacts([]string{"a/"}, "a", "reviews/example.md (new)"), "")
}

// TestFenceWritesNormalizesSpelling covers SC18, a spaced path too.
func TestFenceWritesNormalizesSpelling(t *testing.T) {
	assertFenceWrites(t, fencedFacts([]string{"internal/x/", "internal/spaced dir/"}, "internal/x (new)", "internal/spaced dir (new)"), "")
}

// TestFenceWritesIgnoresImplicitAuthority covers SC19. The spec folder and capture
// authorize a path with no fence line, so neither side counts it.
func TestFenceWritesIgnoresImplicitAuthority(t *testing.T) {
	assertFenceWrites(t, fencedFacts([]string{"a/"}, "a", "specs/example/spec.md", "capture/learnings.md"), "")
	assertFenceWrites(t, fencedFacts([]string{"a/", "specs/example/", "capture/"}, "a"), "")
}

// TestCommandBuildFenceWritesReportsProseToken covers SC20. A backticked token in a
// prose sentence of the fence section parses as a fence entry, and no ticket writes it.
func TestCommandBuildFenceWritesReportsProseToken(t *testing.T) {
	slug := "example"
	preflighttest.StartRepo(t)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", preflighttest.SpecBody(slug, "Edit the shared rules only in `.bench/BENCH.md`."))
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", preflighttest.TicketDoc("One", "PF1", "PF2"))
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c0")
	preflighttest.RunGit(t, "checkout", "-q", "-b", "feature")
	preflighttest.MustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	preflighttest.RunGit(t, "add", ".")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")

	out, code := Command([]string{"build", slug})
	want := "  fence-writes,red,\"" + fenceWritesLead + "fence only: .bench/BENCH.md\",\"\"\n"
	if code != 1 || !strings.Contains(out, want) {
		t.Fatalf("build = (%d), want the row %q:\n%s", code, want, out)
	}
}

// TestProposeWritesToleratesFenceWrites covers SC21. The proposal reports the fence
// state the red requires, while a charge still refuses on the same red.
func TestProposeWritesToleratesFenceWrites(t *testing.T) {
	root, slug := seedProposalWith(t, "internal/example/foo.go", nil, "- `unowned/`")
	out, code := Command(proposalArgs(t, root, slug))
	if code != 0 || !strings.Contains(out, "writes_proposal[") {
		t.Fatalf("proposal with a fence-writes red = (%d):\n%s", code, out)
	}
	charge := []string{"build", slug, "--charge", "--ticket", "one.md", "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")}
	if chargeOut, chargeCode := Command(charge); chargeCode != 1 || !strings.Contains(chargeOut, "fence-writes") {
		t.Fatalf("charge with a fence-writes red = (%d):\n%s", chargeCode, chargeOut)
	}
}
