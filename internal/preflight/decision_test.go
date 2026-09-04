package preflight

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/tickets"
)

// baseFacts is a conformant Facts value: every check should answer green over it
// unmodified. Each test case starts from a copy and breaks exactly one field.
func baseFacts() Facts {
	return Facts{
		Mode:                  "review",
		SpecPath:              "specs/example/spec.md",
		DefaultBranchResolved: true,
		DefaultBranchCurrent:  true,
		ReviewBaseResolved:    true,
		ChangedPaths:          []string{"internal/example/foo.go"},
		FenceEntries:          []string{"internal/example"},
		DeclaredRowIDs:        []string{"PF1", "PF2"},
		SpecTag:               "PF",
		Tickets: []tickets.Ticket{{
			Name:   "one.md",
			Writes: []string{"internal/example/foo.go"},
			Covers: []string{"PF1", "PF2"},
		}},
		WritesPathExists:  map[string]bool{"internal/example/foo.go": true},
		BinarySealPresent: true,
	}
}

// coveringTickets is one ticket citing exactly tokens, with no `Writes:` entry
// of its own. A citation-shaped test says nothing about ownership.
func coveringTickets(tokens ...string) []tickets.Ticket {
	return []tickets.Ticket{{Name: "one.md", Covers: tokens}}
}

func checkRow(v Verdict, name string) (CheckResult, bool) {
	for _, c := range v.Checks {
		if c.Check == name {
			return c, true
		}
	}
	return CheckResult{}, false
}

// TestDecideAllGreen is the tracer: every check present by name, all
// green, and a second Decide over the same Facts is byte-identical. This
// is the rerun guarantee the verdict core promises.
func TestDecideAllGreen(t *testing.T) {
	f := baseFacts()
	first := Decide(f)
	second := Decide(f)

	wantChecks := []string{
		"base-current", "paths-authorized",
		"tickets-parse", "blockers-resolve", "writes-resolve",
		"fixture-closure", "registry-closure", "kit-pin",
		"rows-owned", "rows-membership", "diff-nonempty",
	}
	if len(first.Checks) != len(wantChecks) {
		t.Fatalf("Checks = %#v, want %d rows", first.Checks, len(wantChecks))
	}
	for i, name := range wantChecks {
		if first.Checks[i].Check != name {
			t.Errorf("Checks[%d].Check = %q, want %q", i, first.Checks[i].Check, name)
		}
		if first.Checks[i].Verdict != verdictGreen {
			t.Errorf("Checks[%d] (%s) verdict = %q, want green: %+v", i, name, first.Checks[i].Verdict, first.Checks[i])
		}
	}
	if first.Red {
		t.Errorf("Verdict.Red = true, want false for an all-green Facts value")
	}
	if !decisionsEqual(first, second) {
		t.Errorf("Decide is not deterministic: first=%#v second=%#v", first, second)
	}

	// WF35: a remedy belongs to a red the verdict can answer, so no green row
	// carries one. The build-mode pass below adds the not-applicable rows, which
	// answer nothing and so carry none either.
	for _, c := range first.Checks {
		if c.Next != "" {
			t.Errorf("green %s row carries Next = %q, want empty", c.Check, c.Next)
		}
	}
	build := baseFacts()
	build.Mode = modeBuild
	build.TicketsDirExists = false
	naSeen := 0
	for _, c := range Decide(build).Checks {
		if c.Verdict == verdictNA {
			naSeen++
		}
		if c.Next != "" {
			t.Errorf("%s row (%s) carries Next = %q, want empty", c.Check, c.Verdict, c.Next)
		}
	}
	// Every row that reads a parsed ticket, plus diff-nonempty, is
	// not-applicable in a build with no tickets/ directory at all.
	if naSeen != 9 {
		t.Fatalf("fixture invalid: build mode with no tickets/ gave %d not-applicable rows, want 9", naSeen)
	}
}

func decisionsEqual(a, b Verdict) bool {
	if a.Red != b.Red || len(a.Checks) != len(b.Checks) {
		return false
	}
	for i := range a.Checks {
		if a.Checks[i] != b.Checks[i] {
			return false
		}
	}
	return true
}

func TestDecideBaseCurrent(t *testing.T) {
	f := baseFacts()
	f.DefaultBranchResolved = false
	v := Decide(f)
	c, ok := checkRow(v, "base-current")
	if !ok || c.Verdict != verdictRed {
		t.Fatalf("base-current with unresolved default branch = %+v, want red", c)
	}
	if !v.Red {
		t.Error("Verdict.Red = false, want true")
	}

	f = baseFacts()
	f.DefaultBranchCurrent = false
	f.DefaultBranch = "main"
	f.AssignmentTarget = "0123456789abcdef0123456789abcdef"
	c, _ = checkRow(Decide(f), "base-current")
	if c.Verdict != verdictRed {
		t.Fatalf("base-current with stale default branch = %+v, want red", c)
	}
	// WF33: the stale-base red is the one red that answers itself. The remedy
	// names the gathered branch and the assignment the operator is standing in.
	if c.Next != "bench worktree merge --from main 0123456789abcdef0123456789abcdef" {
		t.Errorf("stale base Next = %q, want the merge remedy naming the assignment id", c.Next)
	}

	// WF34: with no assignment owning this root, the id slot renders the
	// placeholder rather than an empty word that would print a broken command.
	f.AssignmentTarget = ""
	c, _ = checkRow(Decide(f), "base-current")
	if c.Next != "bench worktree merge --from main <target>" {
		t.Errorf("stale base Next without an assignment = %q, want the placeholder target", c.Next)
	}
}

// TestDecideOtherRedsCarryNoNext is WF36. One red answers itself; every other
// red states its detail and stops. A remedy copied onto all of them would name
// a merge for an unresolved branch or for a fence miss no merge repairs.
func TestDecideOtherRedsCarryNoNext(t *testing.T) {
	cases := []struct {
		name   string
		check  string
		detail string
		mutate func(*Facts)
	}{
		{"unresolved default branch", "base-current", "default branch does not resolve", func(f *Facts) {
			f.DefaultBranchResolved = false
		}},
		{"unresolved source base", "base-current", "source base does not resolve: --base is not an ancestor", func(f *Facts) {
			f.ExplicitSourceRange = true
			f.ReviewBaseResolved = false
			f.ReviewBaseHint = "--base is not an ancestor"
		}},
		{"unresolved source tip", "tip-current", "source tip does not resolve, so --source-tip cafe cannot be verified", func(f *Facts) {
			f.PinnedSourceTip = "cafe"
			f.SourceTip = ""
		}},
		{"pinned tip mismatch", "tip-current", "--source-tip cafe is not the derived source tip beef", func(f *Facts) {
			f.PinnedSourceTip = "cafe"
			f.SourceTip = "beef"
		}},
		{"paths unresolved review base", "paths-authorized", "review base does not resolve: no resolvable default branch", func(f *Facts) {
			f.ReviewBaseResolved = false
			f.ReviewBaseHint = "no resolvable default branch"
		}},
		{"unfenced path", "paths-authorized", "not authorized by any ownership fence: unfenced/path.go", func(f *Facts) {
			f.ChangedPaths = []string{"unfenced/path.go"}
		}},
		{"uncited row", "rows-owned", "declared row(s) cited by no ticket file: PF3", func(f *Facts) {
			f.DeclaredRowIDs = []string{"PF1", "PF2", "PF3"}
		}},
		{"phantom token", "rows-membership", "ticket token(s) under this spec's tag name no declared row: PF99", func(f *Facts) {
			f.Tickets = coveringTickets("PF1", "PF2", "PF99")
		}},
		{"diff unresolved review base", "diff-nonempty", "review base does not resolve: no resolvable default branch", func(f *Facts) {
			f.ReviewBaseResolved = false
			f.ReviewBaseHint = "no resolvable default branch"
		}},
		{"empty diff", "diff-nonempty", "no changed files since the resolved review base", func(f *Facts) {
			f.ChangedPaths = nil
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := baseFacts()
			f.DefaultBranch = "main"
			f.AssignmentTarget = "0123456789abcdef0123456789abcdef"
			tc.mutate(&f)
			c, ok := checkRow(Decide(f), tc.check)
			if !ok || c.Verdict != verdictRed || c.Detail != tc.detail {
				t.Fatalf("%s row = %+v, want the red detailed %q", tc.check, c, tc.detail)
			}
			if c.Next != "" {
				t.Errorf("%s red carries Next = %q, want empty", tc.check, c.Next)
			}
		})
	}
}

func TestDecideExplicitSourceRangeReplacesDefaultBranchAncestry(t *testing.T) {
	// A retained source may predate destination-only phase commits. Its frozen
	// base-to-tip range, not destination ancestry, is the applicable validity fact.
	f := baseFacts()
	f.DefaultBranchCurrent = false
	f.ExplicitSourceRange = true
	if c, _ := checkRow(Decide(f), "base-current"); c.Verdict != verdictGreen {
		t.Fatalf("explicit source range with stale destination = %+v, want green", c)
	}

	// Mutation control: omitting --base restores the default-branch predicate.
	f.ExplicitSourceRange = false
	if c, _ := checkRow(Decide(f), "base-current"); c.Verdict != verdictRed {
		t.Fatalf("bare stale destination = %+v, want red", c)
	}

	f.ExplicitSourceRange = true
	f.ReviewBaseResolved = false
	f.ReviewBaseHint = "--base is not an ancestor"
	if c, _ := checkRow(Decide(f), "base-current"); c.Verdict != verdictRed || c.Detail != "source base does not resolve: --base is not an ancestor" {
		t.Fatalf("invalid explicit source range = %+v, want source-validity red", c)
	}
}

func TestDecidePathsAuthorized(t *testing.T) {
	f := baseFacts()
	f.ChangedPaths = []string{"unfenced/path.go"}
	c, _ := checkRow(Decide(f), "paths-authorized")
	if c.Verdict != verdictRed || c.Detail == "" {
		t.Fatalf("paths-authorized with an out-of-fence path = %+v, want red naming the path", c)
	}

	// A path equal to a fence entry, or under a fence prefix with a "/"
	// separator, stays green. A same-string-prefix path that lacks the
	// separator does not.
	f = baseFacts()
	f.FenceEntries = []string{"internal/git"}
	f.ChangedPaths = []string{"internal/git", "internal/git/sub.go"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictGreen {
		t.Errorf("paths-authorized with exact and prefix matches = %+v, want green", c)
	}
	f.ChangedPaths = []string{"internal/git2"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictRed {
		t.Errorf("paths-authorized: internal/git2 must not match fence internal/git = %+v, want red", c)
	}

	// An empty changed set is green.
	f = baseFacts()
	f.ChangedPaths = nil
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictGreen {
		t.Errorf("paths-authorized with no changed paths = %+v, want green", c)
	}

	f = baseFacts()
	f.ReviewBaseResolved = false
	f.ReviewBaseHint = "no resolvable default branch"
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictRed {
		t.Errorf("paths-authorized with an unresolved review base = %+v, want red", c)
	}
}

func TestDecideRowsOwned(t *testing.T) {
	f := baseFacts()
	f.DeclaredRowIDs = []string{"PF1", "PF2", "PF3"}
	f.Tickets = coveringTickets("PF1", "PF2")
	c, _ := checkRow(Decide(f), "rows-owned")
	if c.Verdict != verdictRed || c.Detail == "" {
		t.Fatalf("rows-owned with one uncited declared row = %+v, want red naming PF3", c)
	}
}

func TestDecideRowsMembership(t *testing.T) {
	// A token under the spec's own tag naming no declared row is red.
	f := baseFacts()
	f.Tickets = coveringTickets("PF1", "PF2", "PF99")
	c, _ := checkRow(Decide(f), "rows-membership")
	if c.Verdict != verdictRed || c.Detail == "" {
		t.Fatalf("rows-membership with a phantom same-tag token = %+v, want red naming PF99", c)
	}

	// A foreign-tag token is ignored, not flagged.
	f = baseFacts()
	f.Tickets = coveringTickets("PF1", "PF2", "FT93")
	c, _ = checkRow(Decide(f), "rows-membership")
	if c.Verdict != verdictGreen {
		t.Fatalf("rows-membership with a foreign-tag token = %+v, want green (ignored)", c)
	}
}

// TestWritesResolveNamesAbsentPath is TG12. A typo path charges a delegate
// against nothing, so the row names the offending entry and the ticket that
// declared it.
func TestWritesResolveNamesAbsentPath(t *testing.T) {
	f := baseFacts()
	f.Tickets = []tickets.Ticket{{
		Name:   "one.md",
		Writes: []string{"internal/example/foo.go", "internal/example/tpyo.go"},
		Covers: []string{"PF1", "PF2"},
	}}
	f.WritesPathExists = map[string]bool{
		"internal/example/foo.go":  true,
		"internal/example/tpyo.go": false,
	}

	c, ok := checkRow(Decide(f), "writes-resolve")
	want := "Writes: entry names no tree path and carries no (new) marker: one.md: internal/example/tpyo.go"
	if !ok || c.Verdict != verdictRed || c.Detail != want {
		t.Fatalf("writes-resolve = %+v, want the red detailed %q", c, want)
	}
}

// TestWritesResolveAcceptsNewMarker is TG13. A blocker ticket may land the
// file first, so a declared new file is green while the tree still lacks it.
// An unconditional exists check would red every new-file ticket.
func TestWritesResolveAcceptsNewMarker(t *testing.T) {
	f := baseFacts()
	f.Tickets = []tickets.Ticket{{
		Name:   "one.md",
		Writes: []string{"internal/example/next.go (new)"},
		Covers: []string{"PF1", "PF2"},
	}}
	f.WritesPathExists = map[string]bool{"internal/example/next.go (new)": false}
	if c, _ := checkRow(Decide(f), "writes-resolve"); c.Verdict != verdictGreen {
		t.Fatalf("writes-resolve over a (new) entry the tree lacks = %+v, want green", c)
	}

	// Mutation control: strip the marker and the same absent path reds.
	f.Tickets[0].Writes = []string{"internal/example/next.go"}
	f.WritesPathExists = map[string]bool{"internal/example/next.go": false}
	if c, _ := checkRow(Decide(f), "writes-resolve"); c.Verdict != verdictRed {
		t.Fatalf("writes-resolve over the same path without the marker = %+v, want red", c)
	}
}

// TestCoversDrivesRowOwnership is TG21 and TG22. The parsed Covers: field is
// the one citation source, so an uncited declared row and a phantom token each
// red their own row by name.
func TestCoversDrivesRowOwnership(t *testing.T) {
	// TG21: a declared row no ticket cites reds rows-owned naming the row.
	f := baseFacts()
	f.DeclaredRowIDs = []string{"PF1", "PF2", "PF3"}
	c, _ := checkRow(Decide(f), "rows-owned")
	want := "declared row(s) cited by no ticket file: PF3"
	if c.Verdict != verdictRed || c.Detail != want {
		t.Fatalf("rows-owned = %+v, want the red detailed %q", c, want)
	}

	// TG22: a Covers: token naming no declared row reds rows-membership naming
	// the token.
	f = baseFacts()
	f.Tickets = coveringTickets("PF1", "PF2", "PF99")
	c, _ = checkRow(Decide(f), "rows-membership")
	want = "ticket token(s) under this spec's tag name no declared row: PF99"
	if c.Verdict != verdictRed || c.Detail != want {
		t.Fatalf("rows-membership = %+v, want the red detailed %q", c, want)
	}

	// Two tickets may share one row: rows-owned needs one owner, not exactly one.
	f = baseFacts()
	f.Tickets = []tickets.Ticket{
		{Name: "one.md", Covers: []string{"PF1"}},
		{Name: "two.md", Covers: []string{"PF1", "PF2"}},
	}
	if c, _ := checkRow(Decide(f), "rows-owned"); c.Verdict != verdictGreen {
		t.Errorf("rows-owned with a shared row = %+v, want green", c)
	}
}

// TestProseRowIDIsNotOwnership is TG20 at the preflight seam. A row ID that
// appears only in a ticket's prose lands in no parsed Covers: set, so it is
// not evidence. The deleted bare-token scrape read it as ownership.
func TestProseRowIDIsNotOwnership(t *testing.T) {
	f := baseFacts()
	f.Tickets = []tickets.Ticket{{
		Name:   "one.md",
		Title:  "Cover PF1 and mention PF2 in passing",
		Covers: []string{"PF1"},
	}}
	c, _ := checkRow(Decide(f), "rows-owned")
	want := "declared row(s) cited by no ticket file: PF2"
	if c.Verdict != verdictRed || c.Detail != want {
		t.Fatalf("rows-owned with PF2 named only in prose = %+v, want the red detailed %q", c, want)
	}

	// A phantom token in prose is equally invisible to rows-membership.
	f.DeclaredRowIDs = []string{"PF1"}
	f.Tickets[0].Title = "Cover PF1; PF99 is only prose"
	if c, _ := checkRow(Decide(f), "rows-membership"); c.Verdict != verdictGreen {
		t.Errorf("rows-membership with PF99 named only in prose = %+v, want green", c)
	}
}

// TestTicketGrammarRowsCarryOwnDetails is the attribution contract: each
// grammar row reds on its own fact with its own detail, so one red never
// hides another.
func TestTicketGrammarRowsCarryOwnDetails(t *testing.T) {
	f := baseFacts()
	f.TicketDiagnostics = []string{"one.md: missing Writes"}
	f.BlockerCycles = []string{"cycle edge one.md -> two.md"}
	f.Tickets = []tickets.Ticket{{Name: "one.md", Writes: []string{"gone.go"}, Covers: []string{"PF1", "PF2"}}}
	f.WritesPathExists = map[string]bool{"gone.go": false}

	for _, tc := range []struct{ check, detail string }{
		{"tickets-parse", "ticket grammar fault(s): one.md: missing Writes"},
		{"blockers-resolve", "blocker cycle(s): cycle edge one.md -> two.md"},
		{"writes-resolve", "Writes: entry names no tree path and carries no (new) marker: one.md: gone.go"},
	} {
		c, ok := checkRow(Decide(f), tc.check)
		if !ok || c.Verdict != verdictRed || c.Detail != tc.detail {
			t.Errorf("%s row = %+v, want the red detailed %q", tc.check, c, tc.detail)
		}
	}
}

func TestDecideDiffNonempty(t *testing.T) {
	f := baseFacts()
	f.ChangedPaths = nil
	c, _ := checkRow(Decide(f), "diff-nonempty")
	if c.Verdict != verdictRed {
		t.Fatalf("diff-nonempty with an empty changed set = %+v, want red", c)
	}

	f = baseFacts()
	f.ReviewBaseResolved = false
	c, _ = checkRow(Decide(f), "diff-nonempty")
	if c.Verdict != verdictRed {
		t.Fatalf("diff-nonempty with an unresolved review base = %+v, want red", c)
	}
}

// TestDecidePathsAuthorizedImplicitSpecFolder covers LS7-LS11. The active
// spec's own folder authorizes changed paths without a declared fence
// entry, in every mode, at a path-segment boundary. It works beside the
// declared entries rather than in place of them.
func TestDecidePathsAuthorizedImplicitSpecFolder(t *testing.T) {
	// LS7: build preflight authorizes a path under the active spec's folder with no
	// self-fence entry.
	f := baseFacts()
	f.Mode = "build"
	f.ChangedPaths = []string{"specs/example/spec.md"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictGreen {
		t.Errorf("LS7: build mode, path under the spec's own folder = %+v, want green", c)
	}

	// LS8: review preflight authorizes that same path the same way. A
	// mode-keyed implicit entry leaves one of these two rows red.
	f = baseFacts()
	f.Mode = "review"
	f.ChangedPaths = []string{"specs/example/spec.md"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictGreen {
		t.Errorf("LS8: review mode, path under the spec's own folder = %+v, want green", c)
	}

	// LS9: a path under a different spec's folder stays unauthorized.
	f = baseFacts()
	f.ChangedPaths = []string{"specs/other/spec.md"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictRed {
		t.Errorf("LS9: path under a foreign spec's folder = %+v, want red", c)
	}

	// LS10: a sibling folder whose name merely extends the slug stays
	// unauthorized. The boundary is a path segment, not a string prefix.
	f = baseFacts()
	f.ChangedPaths = []string{"specs/example-two/notes.md"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictRed {
		t.Errorf("LS10: sibling folder extending the slug = %+v, want red", c)
	}

	// LS11: declared entries keep their exact semantics. A path authorized
	// only by a declared fence stays green, and one outside every fence and
	// the spec folder is red.
	f = baseFacts()
	f.ChangedPaths = []string{"internal/example/foo.go"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictGreen {
		t.Errorf("LS11: declared-fence-only path = %+v, want green", c)
	}
	f.ChangedPaths = []string{"unfenced/path.go"}
	if c, _ := checkRow(Decide(f), "paths-authorized"); c.Verdict != verdictRed {
		t.Errorf("LS11: path outside every fence and the spec folder = %+v, want red", c)
	}
}

// TestFixtureClosureNamesUnnamedFixture covers TG14. A ticket that writes a
// fixture-pinned path and does not name the pinning fixture directory reds with
// both the path and the fixture named.
func TestFixtureClosureNamesUnnamedFixture(t *testing.T) {
	f := baseFacts()
	f.WritesFixturePins = map[string][]string{
		"internal/example/foo.go": {"tests/canary/example-family/pinning-fixture"},
	}
	c, ok := checkRow(Decide(f), "fixture-closure")
	const want = "Writes: entry names a fixture-pinned path without naming the fixture: " +
		"one.md: internal/example/foo.go is pinned by tests/canary/example-family/pinning-fixture"
	if !ok || c.Verdict != verdictRed || c.Detail != want {
		t.Fatalf("fixture-closure row = %+v, want the red detailed %q", c, want)
	}

	// Naming the fixture directory closes the row.
	f.Tickets[0].Writes = append(f.Tickets[0].Writes, "tests/canary/example-family/pinning-fixture")
	if c, _ := checkRow(Decide(f), "fixture-closure"); c.Verdict != verdictGreen {
		t.Errorf("fixture-closure with the fixture named = %+v, want green", c)
	}
}

// TestRegistryClosureNamesOmittedRegistry covers TG16. A ticket that writes a
// bound package and omits a bound file reds with the omitted file named.
func TestRegistryClosureNamesOmittedRegistry(t *testing.T) {
	f := baseFacts()
	f.WritesBoundFiles = map[string][]string{
		"internal/example/foo.go": {"cmd/bench/command_registry.go", "cmd/bench/main_test.go"},
	}
	f.Tickets[0].Writes = []string{"internal/example/foo.go", "cmd/bench/command_registry.go"}
	c, ok := checkRow(Decide(f), "registry-closure")
	const want = "Writes: entry names a bound package without naming every bound file: " +
		"one.md: internal/example/foo.go requires cmd/bench/main_test.go"
	if !ok || c.Verdict != verdictRed || c.Detail != want {
		t.Fatalf("registry-closure row = %+v, want the red detailed %q", c, want)
	}

	f.Tickets[0].Writes = append(f.Tickets[0].Writes, "cmd/bench/main_test.go")
	if c, _ := checkRow(Decide(f), "registry-closure"); c.Verdict != verdictGreen {
		t.Errorf("registry-closure with every bound file named = %+v, want green", c)
	}
}

// TestKitPinRequiresBenchKit covers TG19 and its green counterpart. A written
// system-tagged test file reds unless the ticket body states BENCH_KIT; a test
// file with no system tag stays green either way.
func TestKitPinRequiresBenchKit(t *testing.T) {
	f := baseFacts()
	f.Tickets[0].Writes = []string{"internal/example/sys_test.go"}
	f.WritesPathExists = map[string]bool{"internal/example/sys_test.go": true}
	f.WritesSystemTagged = map[string]bool{"internal/example/sys_test.go": true}
	c, ok := checkRow(Decide(f), "kit-pin")
	const want = "ticket writes a system-tagged test file without stating BENCH_KIT: " +
		"one.md: internal/example/sys_test.go"
	if !ok || c.Verdict != verdictRed || c.Detail != want {
		t.Fatalf("kit-pin row = %+v, want the red detailed %q", c, want)
	}

	f.TicketPinsKit = map[string]bool{"one.md": true}
	if c, _ := checkRow(Decide(f), "kit-pin"); c.Verdict != verdictGreen {
		t.Errorf("kit-pin with BENCH_KIT stated = %+v, want green", c)
	}

	untagged := baseFacts()
	untagged.Tickets[0].Writes = []string{"internal/example/plain_test.go"}
	untagged.WritesPathExists = map[string]bool{"internal/example/plain_test.go": true}
	if c, _ := checkRow(Decide(untagged), "kit-pin"); c.Verdict != verdictGreen {
		t.Errorf("kit-pin over a test file with no system tag = %+v, want green", c)
	}
}

// TestBinarySealRedsAMismatchedDestination covers BF26. A build whose
// published dist/bench fails its seal reds binary-seal, carries the
// verifier's refusal whole, and makes the whole verdict red, so a --full run
// names the rebuild before it opens the transaction.
func TestBinarySealRedsAMismatchedDestination(t *testing.T) {
	f := baseFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = true
	f.BinarySealPresent = true
	const refusal = `bench binary "/root/dist/bench" is untrusted: source digest mismatch; ` +
		`rebuild with cd '/root' && bash scripts/go-build.sh '/root' '/root/dist/bench'`
	f.BinarySealRefusal = refusal

	v := Decide(f)
	c, ok := checkRow(v, "binary-seal")
	if !ok {
		t.Fatalf("Checks = %#v, want a binary-seal row", v.Checks)
	}
	if c.Verdict != verdictRed || c.Detail != refusal {
		t.Errorf("binary-seal row = %+v, want the red detailed %q", c, refusal)
	}
	if !v.Red {
		t.Errorf("Verdict.Red = false, want true when the destination seal mismatches")
	}
	if at := rowIndex(v, "binary-seal"); at != rowIndex(v, "kit-pin")+1 {
		t.Errorf("binary-seal sits at index %d, want directly after kit-pin at %d", at, rowIndex(v, "kit-pin"))
	}
}

// TestBinarySealNotApplicableWithoutADestinationBinary covers BF27. A root
// that publishes no dist/bench reports the row not applicable, keeps the
// verdict green, and leaves every other row exactly as it was.
func TestBinarySealNotApplicableWithoutADestinationBinary(t *testing.T) {
	f := baseFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = true
	f.BinarySealPresent = false

	v := Decide(f)
	c, ok := checkRow(v, "binary-seal")
	if !ok {
		t.Fatalf("Checks = %#v, want a binary-seal row", v.Checks)
	}
	if c.Verdict != verdictNA || c.Detail != "" || c.Next != "" {
		t.Errorf("binary-seal row = %+v, want not-applicable with no detail and no remedy", c)
	}
	if v.Red {
		t.Errorf("Verdict.Red = true, want false when no destination binary is published")
	}
	for _, other := range v.Checks {
		if other.Check == "binary-seal" || other.Check == "diff-nonempty" {
			continue
		}
		if other.Verdict != verdictGreen {
			t.Errorf("%s row = %+v, want green and unchanged by an absent binary", other.Check, other)
		}
	}
}

// TestBinarySealIsBuildModeOnly pins the mode gate. Story 26 names
// `bench preflight build`, so a review preflight renders no binary-seal row
// at all, even over a root whose published binary fails its seal.
func TestBinarySealIsBuildModeOnly(t *testing.T) {
	f := baseFacts()
	f.BinarySealPresent = true
	f.BinarySealRefusal = "bench binary is untrusted: source digest mismatch"

	v := Decide(f)
	if _, ok := checkRow(v, "binary-seal"); ok {
		t.Errorf("review Checks = %#v, want no binary-seal row", v.Checks)
	}
	if v.Red {
		t.Errorf("Verdict.Red = true, want false: a review preflight does not grade the binary")
	}
}

// sealFixtureRoot is the smallest rebuildable root the seal primitives
// accept: a module whose ./cmd/bench resolves, plus the build-inputs
// manifest the source digest reads.
func sealFixtureRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range map[string]string{
		"go.mod":                  "module example.com/sealfixture\n\ngo 1.25\n",
		"cmd/bench/main.go":       "package main\n\nfunc main() {}\n",
		"scripts/go-build.sh":     "#!/usr/bin/env bash\n",
		"scripts/go-build.inputs": "build_script=scripts/go-build.sh\n",
	} {
		mustWriteFile(t, filepath.Join(root, filepath.FromSlash(name)), body)
	}
	return root
}

// publishSealedBinary seals a dist/bench pair for root through the same
// publisher a build runs, so the fixture's green case is a real seal rather
// than a hand-written sidecar.
func publishSealedBinary(t *testing.T, root string) string {
	t.Helper()
	staged := filepath.Join(root, "staged-bench")
	mustWriteFile(t, staged, "Bench executable")
	if err := os.Chmod(staged, 0o755); err != nil {
		t.Fatal(err)
	}
	executable := freshness.PublishedExecutable(root)
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := freshness.Publish(root, staged, executable, "1.2.3"); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	return executable
}

// TestBinarySealFactsGradeARootOnDisk covers BF26 and BF27 at the gatherer
// seam. Both row tests set the facts by hand, so the path spelling and the
// verifier call answer to nothing. This drives binarySealFacts against real
// trees instead. A dangling dist/bench is the state the two derivations
// disagreed on: doctor lstats and reds it, so the build preflight grades it
// too rather than reading a broken link as no binary at all.
func TestBinarySealFactsGradeARootOnDisk(t *testing.T) {
	for _, tc := range []struct {
		name    string
		place   func(*testing.T, string)
		present bool
		refused bool
	}{
		{
			name:  "no binary published",
			place: func(*testing.T, string) {},
		},
		{
			name:    "sealed pair",
			place:   func(t *testing.T, root string) { publishSealedBinary(t, root) },
			present: true,
		},
		{
			name: "sources changed after the seal",
			place: func(t *testing.T, root string) {
				publishSealedBinary(t, root)
				mustWriteFile(t, filepath.Join(root, "cmd", "bench", "main.go"), "package main\n\nfunc main() { _ = 1 }\n")
			},
			present: true,
			refused: true,
		},
		{
			name: "dangling symbolic link",
			place: func(t *testing.T, root string) {
				executable := freshness.PublishedExecutable(root)
				if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(filepath.Join(root, "absent-bench"), executable); err != nil {
					t.Fatal(err)
				}
			},
			present: true,
			refused: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := sealFixtureRoot(t)
			tc.place(t, root)

			present, refusal := binarySealFacts(root)
			if present != tc.present {
				t.Errorf("binarySealFacts present = %v, want %v (refusal %q)", present, tc.present, refusal)
			}
			if refused := refusal != ""; refused != tc.refused {
				t.Errorf("binarySealFacts refusal = %q, want refused = %v", refusal, tc.refused)
			}
		})
	}
}

// TestDanglingDestinationBinaryRedsBinarySealInBuildMode carries the
// gathered facts of a broken dist/bench link through to the verdict. The
// row an operator reads is the whole point: a --full build that trusts a
// link to nothing runs an executable nobody can name.
func TestDanglingDestinationBinaryRedsBinarySealInBuildMode(t *testing.T) {
	root := sealFixtureRoot(t)
	executable := freshness.PublishedExecutable(root)
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "absent-bench"), executable); err != nil {
		t.Fatal(err)
	}

	f := baseFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = true
	f.BinarySealPresent, f.BinarySealRefusal = binarySealFacts(root)

	v := Decide(f)
	c, ok := checkRow(v, "binary-seal")
	if !ok {
		t.Fatalf("Checks = %#v, want a binary-seal row", v.Checks)
	}
	if c.Verdict != verdictRed {
		t.Errorf("binary-seal row = %+v, want red over a dangling dist/bench", c)
	}
	if !v.Red {
		t.Errorf("Verdict.Red = false, want true when dist/bench links to nothing")
	}
}

// rowIndex is the position of one named row, or -1. Order is part of the
// verdict's contract, so a test that grades placement reads it here.
func rowIndex(v Verdict, name string) int {
	for i, c := range v.Checks {
		if c.Check == name {
			return i
		}
	}
	return -1
}

// TestSixRowsNotApplicableWithoutTickets covers TG39. In build mode with no
// tickets/ directory, the six grammar rows render not-applicable, in order.
func TestSixRowsNotApplicableWithoutTickets(t *testing.T) {
	f := baseFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = false
	v := Decide(f)

	want := []string{
		"tickets-parse", "blockers-resolve", "writes-resolve",
		"fixture-closure", "registry-closure", "kit-pin",
	}
	at := -1
	for i, c := range v.Checks {
		if c.Check == want[0] {
			at = i
			break
		}
	}
	if at < 0 || at+len(want) > len(v.Checks) {
		t.Fatalf("Checks = %#v, want the six grammar rows", v.Checks)
	}
	for i, name := range want {
		row := v.Checks[at+i]
		if row.Check != name || row.Verdict != verdictNA || row.Detail != "" {
			t.Errorf("grammar row %d = %+v, want %s not-applicable with no detail", i, row, name)
		}
	}
	if v.Red {
		t.Errorf("Verdict.Red = true, want false when every grammar row is not-applicable")
	}
}
