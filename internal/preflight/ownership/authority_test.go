package ownership_test

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// breakdownSeam is the PF2 seam cell of the breakdown spec. No ticket names it, so only the
// coverage parser's row can supply it.
const breakdownSeam = "breakdown seam y"

// breakdownRepo seeds the conformant build repository whose main branch carries the
// breakdown spec and tickets as its whole tickets folder, and whose feature branch edits
// internal/example/foo.go. A nil map removes the folder, and an empty map leaves it present
// and empty. It answers the breakdown and the ordinary build arguments over that tip.
func breakdownRepo(t *testing.T, tickets map[string]string) (breakdown, build []string) {
	t.Helper()
	root, _ := preflighttest.SeedConformant(t)
	preflighttest.RunGit(t, "checkout", "-q", "main")
	preflighttest.MustWriteFile(t, specPath, strings.Replace(preflighttest.SpecBody(slug), "| PF2 | 1 | does y | cli seam |", "| PF2 | 1 | does y | "+breakdownSeam+" |", 1))
	dir := "specs/" + slug + "/tickets"
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	for name, body := range tickets {
		preflighttest.MustWriteFile(t, dir+"/"+name, body)
	}
	commitmenttest.SeedAdmission(t, root, specPath)
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", "breakdown tickets")
	preflighttest.RunGit(t, "checkout", "-q", "-B", "feature")
	preflighttest.MustWriteFile(t, "internal/example/foo.go", "package example\n\nconst c1 = 1\n")
	preflighttest.RunGit(t, "add", "internal/example/foo.go")
	preflighttest.RunGit(t, "commit", "-q", "-m", "c1")
	if tickets != nil && len(tickets) == 0 {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	build = []string{"build", slug, "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")}
	return append([]string{"build", slug, "--breakdown"}, build[2:]...), build
}

// mustBreakdown runs the actual breakdown and fails the test unless it answers.
func mustBreakdown(t *testing.T, args []string) string {
	t.Helper()
	out, code := preflight.Command(args)
	if code != 0 {
		t.Fatalf("breakdown = (%d):\n%s", code, out)
	}
	return out
}

// uncitedRows is the row cell of each uncited breakdown row, in response order.
func uncitedRows(t *testing.T, out string) []string {
	t.Helper()
	var rows []string
	for _, row := range table(t, out, "uncited") {
		rows = append(rows, cell(row, "row"))
	}
	return rows
}

// TestBreakdownRows drives the read-only breakdown and the ordinary build over one tip each.
func TestBreakdownRows(t *testing.T) {
	t.Run("OC61 OC62 uncited row", func(t *testing.T) {
		breakdown, build := breakdownRepo(t, map[string]string{"one.md": preflighttest.TicketDoc("One", "PF1") + "\nPF2 waits for a later ticket.\n"})
		out := mustBreakdown(t, breakdown)
		row := one(t, out, "uncited", map[string]string{"row": "PF2"})
		if got := uncitedRows(t, out); !slices.Equal(got, []string{"PF2"}) || cell(row, "story") != "1" || cell(row, "seam") != breakdownSeam {
			t.Fatalf("uncited = %v (%v), want only PF2 with story 1 and seam %q:\n%s", got, row, breakdownSeam, out)
		}
		verdict, code := preflight.Command(build)
		if code != 1 || !strings.Contains(verdict, "rows-owned,red") || !strings.Contains(verdict, "cited by no ticket file: PF2") {
			t.Fatalf("ordinary build = (%d), want the rows-owned red for PF2:\n%s", code, verdict)
		}
		one(t, out, "checks", map[string]string{"check": "rows-owned", "verdict": "red"})
	})
	// The closure tree's ticket cites every row, so its breakdown shows the closure and the
	// premise diagnostics alone: a fixture pin and the unverified kit registry source.
	t.Run("closure and premise diagnostics", func(t *testing.T) {
		tr := closureTree()
		tr.OwnInventory, tr.Tail = true, fenced(completeEmpty)
		tr.Seed(t)
		out := mustBreakdown(t, []string{"build", slug, "--breakdown", "--base", preflighttest.RunGit(t, "rev-parse", "main"), "--source-tip", preflighttest.RunGit(t, "rev-parse", "HEAD")})
		one(t, out, "closure", map[string]string{"ticket": "one.md", "path": "internal/example/pinned.go", "requires": "tests/canary/fam/pin-fixture"})
		if row := one(t, out, "premises", map[string]string{"id": "kit-registry"}); cell(row, "state") != "unknown" {
			t.Fatalf("kit-registry premise = %v, want unknown:\n%s", row, out)
		}
	})
	// The uncited PF2 case also holds the tolerated rows-owned red, which precedes
	// rows-membership in check order, so the refusal must still name rows-membership.
	t.Run("OC63 foreign or phantom token", func(t *testing.T) {
		for _, test := range []struct {
			covers []string
			check  string
		}{
			{[]string{"PF1", "PF2", "PF9"}, "rows-membership"},
			{[]string{"PF1", "PF9"}, "rows-membership"},
			{[]string{"PF1", "PF2", "FT93"}, "tickets-parse"},
		} {
			breakdown, _ := breakdownRepo(t, map[string]string{"one.md": preflighttest.TicketDoc("One", test.covers...)})
			out, code := preflight.Command(breakdown)
			if code != 1 || !strings.Contains(out, test.check) || strings.Contains(out, "uncited[") {
				t.Errorf("breakdown covering %v = (%d), want the %s refusal and no uncited table:\n%s", test.covers, code, test.check, out)
			}
		}
	})
	t.Run("OC64 present empty directory", func(t *testing.T) {
		breakdown, build := breakdownRepo(t, map[string]string{})
		out := mustBreakdown(t, breakdown)
		if got := uncitedRows(t, out); !slices.Equal(got, []string{"PF1", "PF2"}) {
			t.Fatalf("uncited = %v, want every row:\n%s", got, out)
		}
		one(t, out, "checks", map[string]string{"check": "rows-owned", "verdict": "red"})
		if verdict, _ := preflight.Command(build); !strings.Contains(verdict, "rows-owned,red") {
			t.Fatalf("ordinary build over an empty directory lacks the rows-owned red:\n%s", verdict)
		}
	})
	t.Run("OC65 absent directory", func(t *testing.T) {
		breakdown, build := breakdownRepo(t, nil)
		out := mustBreakdown(t, breakdown)
		if got := uncitedRows(t, out); !slices.Equal(got, []string{"PF1", "PF2"}) {
			t.Fatalf("uncited = %v, want every row:\n%s", got, out)
		}
		if rows := matching(t, out, "checks", map[string]string{"check": "rows-owned"}); len(rows) != 0 {
			t.Fatalf("absent directory breakdown reports rows-owned %v, want it not applicable", rows)
		}
		if verdict, _ := preflight.Command(build); strings.Contains(verdict, "rows-owned") {
			t.Fatalf("ordinary build over an absent directory grades rows-owned:\n%s", verdict)
		}
	})
}

// TestBreakdownGrammar grades the breakdown selector through the actual command grammar.
func TestBreakdownGrammar(t *testing.T) {
	t.Run("OC66 incompatible selector", func(t *testing.T) {
		pinned := []string{"build", slug, "--breakdown", "--base", "base", "--source-tip", "tip"}
		for _, test := range []struct {
			args []string
			want string
		}{
			{append(slices.Clone(pinned), "--charge", "--ticket", "one.md"), "cannot be combined"},
			{append(slices.Clone(pinned), "--plan-only"), "cannot be combined"},
			{append(slices.Clone(pinned), "--propose-writes", "--ticket", "one.md"), "cannot be combined"},
			{append(slices.Clone(pinned), "--verify"), "cannot be combined"},
			{append(slices.Clone(pinned), "--ticket", "one.md"), "--ticket"},
			{[]string{"review", slug, "--breakdown"}, "--breakdown requires build"},
		} {
			if out, code := preflight.Command(test.args); code != 2 || !strings.Contains(out, test.want) {
				t.Errorf("%v = (%d), want exit 2 naming %q:\n%s", test.args, code, test.want, out)
			}
		}
	})
	// The root help inventory owns the form's exact text, so this witness finds the one
	// registered root help row that the breakdown selector names.
	t.Run("OC67 help derives the selector", func(t *testing.T) {
		var forms []evidencecmd.HelpRow
		for _, row := range evidencecmd.HelpRows() {
			if slices.Contains(strings.Fields(row.Suffix), "--breakdown") {
				forms = append(forms, row)
			}
		}
		if len(forms) != 1 || forms[0].Description == "" {
			t.Fatalf("root help rows %v, want one described --breakdown form", evidencecmd.HelpRows())
		}
		if out, code := preflight.Command([]string{"--help"}); code != 0 || !strings.Contains(out, "bench preflight"+forms[0].Suffix+"\n") {
			t.Errorf("preflight help = (%d), want the form %q:\n%s", code, forms[0].Suffix, out)
		}
	})
}

// TestExplicitAuthority grades that a proposal grants no authority: the spec fence and the
// ticket Writes: line must both name a path before a charge publishes.
func TestExplicitAuthority(t *testing.T) {
	pinned := []string{"internal/example/pinned.go"}
	t.Run("OC69 outside the spec", func(t *testing.T) {
		root := ownershiptest.ProposalRepo{Writes: pinned, Files: coveredTree().Files}.Seed(t)
		args := proposalArgs(t, root)
		before := authorityState(t, root)
		out, code := preflight.Command(args)
		row := one(t, out, "writes_proposal", map[string]string{"path": "tests/canary/fam/pin-fixture"})
		if code != 0 || cell(row, "fence") != "spec fence expansion required" || authorityState(t, root) != before {
			t.Fatalf("proposal = (%d) %v, want the expansion and unchanged authority:\n%s", code, row, out)
		}
		args[2] = "--charge"
		assertChargeRefused(t, root, args, "fixture-closure")
	})
	t.Run("OC70 wide spec without ticket Writes", func(t *testing.T) {
		root := ownershiptest.ProposalRepo{Writes: pinned, Others: map[string][]string{"two.md": {"tests/canary/fam/pin-fixture"}}, Files: coveredTree().Files}.Seed(t)
		args := proposalArgs(t, root)
		out, _ := preflight.Command(args)
		if row := one(t, out, "writes_proposal", map[string]string{"path": "tests/canary/fam/pin-fixture"}); cell(row, "fence") != "covered" {
			t.Fatalf("proposal row = %v, want a fence-covered path that the ticket still lacks", row)
		}
		args[2] = "--charge"
		assertChargeRefused(t, root, args, "fixture-closure")
	})
	t.Run("OC86 ticket path outside the spec fence", func(t *testing.T) {
		tr := coveredTree()
		tr.Fence = []string{}
		for _, write := range tr.Writes {
			tr.Fence = append(tr.Fence, "- `"+write+"`")
		}
		tr.Writes = append(tr.Writes, "internal/other/extra.go")
		tr.Files["internal/other/extra.go"] = "package other\n"
		root := tr.Seed(t)
		args := proposalArgs(t, root)
		args[2] = "--charge"
		assertChargeRefused(t, root, args, "fence-writes", "Writes only: internal/other/extra.go")
	})
}
