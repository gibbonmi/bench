package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/worktree"
)

// retireSlug is the fixture spec the retire listing tests retire. No fixture label or request
// token outside a planted candidate contains it.
const retireSlug = "alpha"

// retireListingRepo creates a repository and a host assignment whose worktree commits a
// merged-implemented spec, and it moves into that worktree, where retire runs.
func retireListingRepo(t *testing.T, hostRequest string) (string, worktree.Creation) {
	t.Helper()
	_, repo, host := censusAssignment(t, hostRequest)
	writeAXIFixture(t, filepath.Join(host.Path, "specs", retireSlug, "spec.md"), "# Alpha\n\nStatus: implemented\n\nRoadmap: FT7\n")
	runAXIGit(t, "-C", host.Path, "add", ".")
	runAXIGit(t, "-C", host.Path, "commit", "-q", "-m", "implement alpha")
	t.Chdir(host.Path)
	return repo, host
}

// plantAssignment records one live assignment with its own worktree and answers its record.
func plantAssignment(t *testing.T, repo, request, label string) intent.Assignment {
	t.Helper()
	creation, err := worktree.Create(repo, request, label, nil)
	if err != nil {
		t.Fatal(err)
	}
	return creation.Assignment
}

// plantUniqueShiftBranch leaves one unrecorded shift branch with a commit of its own and no
// checkout.
func plantUniqueShiftBranch(t *testing.T, repo, name string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), name)
	runAXIGit(t, "-C", repo, "worktree", "add", "-q", "-b", strings.TrimPrefix(intent.ShiftBranchPrefix(), "refs/heads/")+name, dir, "main")
	writeAXIFixture(t, filepath.Join(dir, name+".txt"), name+"\n")
	runAXIGit(t, "-C", dir, "add", ".")
	runAXIGit(t, "-C", dir, "commit", "-q", "-m", name)
	runAXIGit(t, "-C", repo, "worktree", "remove", "--force", dir)
}

func candidateLine(a intent.Assignment) string {
	return "superseded candidate: " + a.ID + " " + a.Label + " — bench worktree clean --discard-branch --target " + a.ID
}

func candidateLines(out string) []string {
	var lines []string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "superseded candidate:") {
			lines = append(lines, line)
		}
	}
	return lines
}

func branchTips(t *testing.T, repo string) string {
	t.Helper()
	return runAXIGit(t, "-C", repo, "for-each-ref", "refs/heads", "--format=%(refname) %(objectname)")
}

const zeroUniqueLine = "unique refs: 0 — bench worktree clean --discard-branch --unclaimed"

// RI48, RI49, RI51, RI79: a candidate prints with its exact --target command, the calling
// worktree's own assignment prints too, the listing sits between the retired: lines and the
// next: line, the retire exits 0, and every branch keeps its tip.
func TestRetireListsSupersededCandidatesBeforeNext(t *testing.T) {
	repo, host := retireListingRepo(t, "alpha-host")
	build := plantAssignment(t, repo, "r48", "alpha-build")
	before := branchTips(t, repo)
	out, code := dispatch("spec", "retire", retireSlug)
	if code != 0 {
		t.Fatalf("retire = (%d, %q)", code, out)
	}
	got := candidateLines(out)
	for _, want := range []string{candidateLine(build), candidateLine(host.Assignment)} {
		if !slices.Contains(got, want) {
			t.Errorf("candidate lines = %q, want %q", got, want)
		}
	}
	if len(got) != 2 {
		t.Errorf("candidate lines = %q, want exactly two", got)
	}
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	want := append([]string{"retired: specs/" + retireSlug}, got...)
	want = append(want, zeroUniqueLine)
	if len(lines) != len(want)+1 || !slices.Equal(lines[:len(want)], want) || !strings.HasPrefix(lines[len(want)], "next: ") {
		t.Errorf("retire lines = %q, want %q then the next: line", lines, want)
	}
	if after := branchTips(t, repo); after != before {
		t.Errorf("branch tips changed:\nbefore %s\nafter  %s", before, after)
	}
}

// RI77: a path operand derives the same slug as the bare slug.
func TestRetireListsCandidatesForPathOperand(t *testing.T) {
	repo, _ := retireListingRepo(t, "retire-host")
	build := plantAssignment(t, repo, "r77", "alpha-build")
	out, code := dispatch("spec", "retire", filepath.Join("specs", retireSlug, "spec.md"))
	if got := candidateLines(out); code != 0 || !slices.Equal(got, []string{candidateLine(build)}) {
		t.Errorf("retire of the path = (%d, %q), want the candidate %q", code, out, candidateLine(build))
	}
}

// RI70: the request token matches when the label does not.
func TestRetireListsCandidateByRequestToken(t *testing.T) {
	repo, _ := retireListingRepo(t, "retire-host")
	other := plantAssignment(t, repo, "alpha-run", "other")
	out, code := dispatch("spec", "retire", retireSlug)
	if got := candidateLines(out); code != 0 || !slices.Equal(got, []string{candidateLine(other)}) {
		t.Errorf("retire = (%d, %q), want the candidate %q", code, out, candidateLine(other))
	}
}

// RI71: a complete record is retired work and prints no candidate line.
func TestRetireSkipsCompleteRecord(t *testing.T) {
	repo, _ := retireListingRepo(t, "retire-host")
	done := plantAssignment(t, repo, "r71", "alpha-build")
	done.State = intent.StateComplete
	if err := intent.PutAssignment(repo, done); err != nil {
		t.Fatal(err)
	}
	out, code := dispatch("spec", "retire", retireSlug)
	if got := candidateLines(out); code != 0 || len(got) != 0 {
		t.Errorf("retire = (%d, %q), want no candidate line", code, out)
	}
}

// RI50: two unique unrecorded refs print the count and the unclaimed plan command.
func TestRetireCountsUniqueRefs(t *testing.T) {
	repo, _ := retireListingRepo(t, "retire-host")
	plantUniqueShiftBranch(t, repo, "one")
	plantUniqueShiftBranch(t, repo, "two")
	out, code := dispatch("spec", "retire", retireSlug)
	if want := "\nunique refs: 2 — bench worktree clean --discard-branch --unclaimed\n"; code != 0 || !strings.Contains(out, want) {
		t.Errorf("retire = (%d, %q), want the line %q", code, out, want)
	}
}

// RI52: no match still prints the count line, at zero, and no candidate line.
func TestRetireWithNoCandidatePrintsZeroCount(t *testing.T) {
	retireListingRepo(t, "retire-host")
	out, code := dispatch("spec", "retire", retireSlug)
	if code != 0 || len(candidateLines(out)) != 0 || !strings.Contains(out, "\n"+zeroUniqueLine+"\n") {
		t.Errorf("retire = (%d, %q), want %q and no candidate line", code, out, zeroUniqueLine)
	}
}

// RI69: an unreadable ledger prints the unavailable count line and keeps exit 0.
func TestRetireWithUnreadableLedgerKeepsExitZero(t *testing.T) {
	repo, _ := retireListingRepo(t, "alpha-host")
	address, err := intent.Address(repo)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(address, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	out, code := dispatch("spec", "retire", retireSlug)
	if code != 0 || !strings.Contains(out, "\nunique refs: unavailable — read intent ledger: ") || len(candidateLines(out)) != 0 {
		t.Errorf("retire with an unreadable ledger = (%d, %q), want exit 0 and the unavailable count line", code, out)
	}
}

// RI80: help and a refusal print no listing, even with a candidate on the ledger.
func TestRetireHelpAndRefusalPrintNoListing(t *testing.T) {
	repo, host := retireListingRepo(t, "alpha-host")
	plantAssignment(t, repo, "r80", "alpha-build")
	writeAXIFixture(t, filepath.Join(host.Path, "specs", "beta", "spec.md"), "# Beta\n\nStatus: staged\n")
	for _, argv := range [][]string{{"spec", "retire", "--help"}, {"spec", "retire", "beta"}, {"spec", "--help"}} {
		out, _ := dispatch(argv...)
		if strings.Contains(out, "superseded candidate:") || strings.Contains(out, "unique refs:") {
			t.Errorf("%q printed the listing: %q", argv, out)
		}
	}
}
