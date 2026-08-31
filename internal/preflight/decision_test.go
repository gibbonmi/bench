package preflight

import (
	"testing"

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
		WritesPathExists: map[string]bool{"internal/example/foo.go": true},
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
	if naSeen != 6 {
		t.Fatalf("fixture invalid: build mode with no tickets/ gave %d not-applicable rows, want 6", naSeen)
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
