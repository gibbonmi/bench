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

// The listing takes its slug from the operand the retire parsed, so a "--" on either side of
// the operand still lists the candidate.
func TestRetireListsCandidatesForTerminatedOperand(t *testing.T) {
	for _, operands := range [][]string{{"--", retireSlug}, {retireSlug, "--"}} {
		t.Run(strings.Join(operands, " "), func(t *testing.T) {
			repo, _ := retireListingRepo(t, "retire-host")
			build := plantAssignment(t, repo, "r77", "alpha-build")
			out, code := dispatch(append([]string{"spec", "retire"}, operands...)...)
			if got := candidateLines(out); code != 0 || !slices.Equal(got, []string{candidateLine(build)}) {
				t.Errorf("retire %q = (%d, %q), want the candidate %q", operands, code, out, candidateLine(build))
			}
		})
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

// RI50, RI100: two unique unrecorded refs print the count and the unclaimed plan command, and
// faulted rows add a faulted suffix without counting as unique.
func TestRetireCountsUniqueRefs(t *testing.T) {
	cases := []struct {
		name  string
		plant func(t *testing.T, repo string)
		want  string
	}{
		{"two unique refs", func(t *testing.T, repo string) {
			plantUniqueShiftBranch(t, repo, "one")
			plantUniqueShiftBranch(t, repo, "two")
		}, "unique refs: 2 — bench worktree clean --discard-branch --unclaimed"},
		{"a symref to main and a blob tip", func(t *testing.T, repo string) {
			runAXIGit(t, "-C", repo, "symbolic-ref", "refs/heads/bench/assign/orphan/symref", "refs/heads/main")
			content := filepath.Join(t.TempDir(), "blob.txt")
			writeAXIFixture(t, content, "blob\n")
			blob := strings.TrimSpace(runAXIGit(t, "-C", repo, "hash-object", "-w", "--", content))
			gitDir := strings.TrimSpace(runAXIGit(t, "-C", repo, "rev-parse", "--absolute-git-dir"))
			writeAXIFixture(t, filepath.Join(gitDir, "refs", "heads", "bench", "assign", "orphan", "blob"), blob+"\n")
		}, "unique refs: 0, 2 faulted — bench worktree clean --discard-branch --unclaimed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo, _ := retireListingRepo(t, "retire-host")
			tc.plant(t, repo)
			out, code := dispatch("spec", "retire", retireSlug)
			var counts []string
			for _, line := range strings.Split(out, "\n") {
				if strings.HasPrefix(line, "unique refs:") {
					counts = append(counts, line)
				}
			}
			if code != 0 || !slices.Equal(counts, []string{tc.want}) {
				t.Errorf("retire = (%d, %q), want the one count line %q", code, out, tc.want)
			}
		})
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

// RI80, RI97: help and a refusal print no listing, even with a candidate on the ledger. The
// retire reflects its operand into a refusal, so an operand that carries a next: line holds a
// real marker, and only the exit code keeps that refusal free of the listing.
func TestRetireHelpAndRefusalPrintNoListing(t *testing.T) {
	repo, host := retireListingRepo(t, "alpha-host")
	plantAssignment(t, repo, "r80", "alpha-build")
	writeAXIFixture(t, filepath.Join(host.Path, "specs", "beta", "spec.md"), "# Beta\n\nStatus: staged\n")
	cases := []struct {
		argv []string
		code int
	}{
		{[]string{"spec", "retire", "--help"}, 0},
		{[]string{"spec", "retire", "beta"}, 1},
		{[]string{"spec", "--help"}, 0},
		{[]string{"spec", "retire", "gamma\nnext: alpha"}, 1},
	}
	for _, tc := range cases {
		out, code := dispatch(tc.argv...)
		if code != tc.code || strings.Contains(out, "superseded candidate:") || strings.Contains(out, "unique refs:") {
			t.Errorf("%q = (%d, %q), want exit %d and no listing", tc.argv, code, out, tc.code)
		}
	}
}
