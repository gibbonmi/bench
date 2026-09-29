package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/worktree"
)

// useCensusFamily registers one bounded family whose `print` leaf prints lines numbered
// lines and exits exit, and one bounded command `plain` with no leaves that does the same.
func useCensusFamily(t *testing.T, lines, exit int) {
	t.Helper()
	old := commandRegistry
	t.Cleanup(func() { commandRegistry = old })
	printLines := func(c Command, args []string) int {
		linesHandler(stdoutOf, lines)(c, args)
		return exit
	}
	commandRegistry = []commandDefinition{
		{
			Name: "family", Inventory: publicInventory(), Bound: boundResponse, LeafUsage: func() string { return "usage\n" },
			Leaves: []commandLeaf{{Name: "print", Root: rootNone, Bound: boundResponse, Run: func(c Command, _ string, args []string) int {
				return printLines(c, args)
			}}},
		},
		{Name: "plain", Inventory: publicInventory(), Bound: boundResponse, Run: printLines},
	}
}

// dispatch runs argv through the production dispatcher and answers its response.
func dispatch(argv ...string) (string, int) {
	var out bytes.Buffer
	code := Command{Stdout: &out, Stderr: &out}.Run(argv)
	return withoutTreeRow(out.String()), code
}

// censusAssignment creates one Bench assignment of a new repository under a private Bench
// home. It answers the home, the repository, and the creation.
func censusAssignment(t *testing.T, request string) (string, string, worktree.Creation) {
	t.Helper()
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	repo := newAXIEnvelopeRepo(t)
	creation, err := worktree.Create(repo, request, request, nil)
	if err != nil {
		t.Fatal(err)
	}
	return home, repo, creation
}

// BO57 at the dispatcher: a bounded leaf run inside an assignment worktree appends one
// record under that assignment with the verb head and the complete output's size.
func TestDispatcherRecordsOutputInAssignment(t *testing.T) {
	home, repo, creation := censusAssignment(t, "census-output-leaf")
	t.Chdir(creation.Path)
	useCensusFamily(t, 11, 0)
	if out, code := dispatch("family", "print"); code != 0 {
		t.Fatalf("family print = (%d, %q)", code, out)
	}
	want := fmt.Sprintf("bench family print=1/%d", len(numberedLines(11)))
	if got := census.OutputBreakdown(home, repo, creation.Assignment.ID); got != want {
		t.Fatalf("output breakdown = %q, want %q", got, want)
	}
}

// A bounded command with no leaves records under the head `bench <name>`.
func TestDispatcherRecordsHeadOfLeaflessVerb(t *testing.T) {
	home, repo, creation := censusAssignment(t, "census-output-plain")
	t.Chdir(creation.Path)
	useCensusFamily(t, 2, 0)
	if out, code := dispatch("plain"); code != 0 {
		t.Fatalf("plain = (%d, %q)", code, out)
	}
	want := fmt.Sprintf("bench plain=1/%d", len(numberedLines(2)))
	if got := census.OutputBreakdown(home, repo, creation.Assignment.ID); got != want {
		t.Fatalf("output breakdown = %q, want %q", got, want)
	}
}

// BO58: `bench worktree exec <target>` from the primary checkout records its response
// under the target assignment, with the exec head.
func TestExecOutputRecordUsesTarget(t *testing.T) {
	home, repo, creation := censusAssignment(t, "census-output-exec")
	t.Chdir(repo)
	if out, code := dispatch("worktree", "exec", creation.Assignment.Label, "--", "echo", "child"); code != 0 {
		t.Fatalf("exec = (%d, %q)", code, out)
	}
	want := fmt.Sprintf("bench worktree exec=1/%d", len("child\n"))
	if got := census.OutputBreakdown(home, repo, creation.Assignment.ID); got != want {
		t.Fatalf("output breakdown under the target = %q, want %q", got, want)
	}
}

// A verb that retires its own working tree writes no record, because the ledger no longer
// holds that assignment as active when the verb returns. The row reads the retirement from
// the removed tree rather than from the exit code, because a release run from inside its
// own tree removes the tree and then exits 1 on a later lookup under the removed root.
func TestReleaseOfOwnWorktreeWritesNoRecord(t *testing.T) {
	request := "census-output-release"
	home, repo, creation := censusAssignment(t, request)
	t.Chdir(creation.Path)
	out, code := dispatch("worktree", "release", "--request", request, creation.Path)
	if _, err := os.Stat(creation.Path); !os.IsNotExist(err) {
		t.Fatalf("release = (%d, %q), want the tree removed: %v", code, out, err)
	}
	if got := census.OutputBreakdown(home, repo, creation.Assignment.ID); got != "" {
		t.Fatalf("output breakdown of the released assignment = %q, want no record", got)
	}
}

// A retiring verb run from a live worktree records under that worktree's assignment.
func TestReleaseFromLiveWorktreeRecordsUnderIt(t *testing.T) {
	request := "census-output-released"
	home, repo, released := censusAssignment(t, request)
	live, err := worktree.Create(repo, "census-output-live", "census-output-live", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(live.Path)
	out, code := dispatch("worktree", "release", "--request", request, released.Path)
	if code != 0 {
		t.Fatalf("release = (%d, %q)", code, out)
	}
	want := fmt.Sprintf("bench worktree release=1/%d", len(out))
	if got := census.OutputBreakdown(home, repo, live.Assignment.ID); got != want {
		t.Fatalf("output breakdown under the live worktree = %q, want %q", got, want)
	}
}

// A reported assignment that the ledger no longer holds as active gets no record, so an
// exec whose child retired its target leaves no file behind the retirement's drop.
func TestRecordOutputSkipsReleasedAssignment(t *testing.T) {
	request := "census-output-gone"
	home, repo, creation := censusAssignment(t, request)
	t.Chdir(repo)
	if out, code := dispatch("worktree", "release", "--request", request, creation.Path); code != 0 {
		t.Fatalf("release = (%d, %q)", code, out)
	}
	size := responsebound.Size{Lines: 1, Bytes: int64(len("child\n"))}
	if err := recordOutput(size, "bench worktree exec", creation.Assignment.ID, repo); err != nil {
		t.Fatal(err)
	}
	if got := census.OutputBreakdown(home, repo, creation.Assignment.ID); got != "" {
		t.Fatalf("output breakdown of the released assignment = %q, want no record", got)
	}
}

// A reported assignment that the ledger still holds, but not as active, gets no record: the
// active check, and not only the ledger lookup, gates the record.
func TestRecordOutputSkipsInactiveAssignment(t *testing.T) {
	home, repo, creation := censusAssignment(t, "census-output-pending")
	pending := creation.Assignment
	pending.State = intent.StateCleanupPending
	if err := intent.PutAssignment(repo, pending); err != nil {
		t.Fatal(err)
	}
	size := responsebound.Size{Lines: 1, Bytes: int64(len("child\n"))}
	if err := recordOutput(size, "bench worktree exec", pending.ID, repo); err != nil {
		t.Fatal(err)
	}
	if got := census.OutputBreakdown(home, repo, pending.ID); got != "" {
		t.Fatalf("output breakdown of the cleanup-pending assignment = %q, want no record", got)
	}
}

// BO62: a symlink at the census directory refuses the record, and the verb's response and
// exit code stay the command's own. The same run without the symlink records, so the
// refusal is the record path failing and not a path the dispatcher never reached.
func TestOutputRecordFailureKeepsVerdict(t *testing.T) {
	home, repo, creation := censusAssignment(t, "census-output-symlink")
	t.Chdir(creation.Path)
	useCensusFamily(t, 3, 3)
	elsewhere := t.TempDir()
	dir := census.Dir(home, repo)
	if err := os.MkdirAll(filepath.Dir(dir), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(elsewhere, dir); err != nil {
		t.Fatal(err)
	}
	if out, code := dispatch("family", "print"); code != 3 || out != numberedLines(3) {
		t.Fatalf("family print with a symlinked census = (%d, %q), want (3, %q)", code, out, numberedLines(3))
	}
	if entries, err := os.ReadDir(elsewhere); err != nil || len(entries) != 0 {
		t.Fatalf("the record followed the symlink: %v (%v)", entries, err)
	}
	if err := os.Remove(dir); err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("bench family print=1/%d", len(numberedLines(3)))
	if out, code := dispatch("family", "print"); code != 3 || census.OutputBreakdown(home, repo, creation.Assignment.ID) != want {
		t.Fatalf("family print without the symlink = (%d, %q), want exit 3 and the record %q", code, out, want)
	}
}
