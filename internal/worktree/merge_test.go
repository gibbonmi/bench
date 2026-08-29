// The worktree merge verb's command-seam tests: composition, publication, and the
// advertised operand grammar, driven in process against real fixture repositories.
package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// mergeFixture is one repository plus one owned assignment per label, and the seam set
// whose declared fast lane appends to a tally file. The tally is the lane record the
// idempotence row reads: a lane that ran leaves a byte behind, one that did not leaves
// no file at all. The lane arrives through the seam rather than the process environment,
// so every row below stays parallel-eligible.
func mergeFixture(t *testing.T, labels ...string) (j joins, root, home, tally string, created []Creation) {
	t.Helper()
	home = filepath.Join(t.TempDir(), "bench-home")
	tally = filepath.Join(t.TempDir(), "lane-tally")
	root = newWorktreeRepo(t)
	for _, label := range labels {
		created = append(created, mustCreate(t, root, home, "merge-"+label, label))
	}
	j = defaultJoins()
	j.mergeLane = func(string) ([]gate.Phase, string, error) {
		return []gate.Phase{{Name: "unit", Argv: []string{"sh", "-c", "printf g >> " + sanitize.ShellQuote(tally)}}}, "", nil
	}
	return j, root, home, tally, created
}

// commitOnDefault advances the default branch, so the commit it returns is one the
// bootstrap authority accepts as `--from`.
func commitOnDefault(t *testing.T, root, name, body string) string {
	t.Helper()
	mustWrite(t, filepath.Join(root, name), []byte(body), 0o644)
	gitRun(t, root, "add", name)
	gitRun(t, root, "commit", "-q", "-m", "advance "+name)
	return gitOutput(t, root, "rev-parse", "HEAD")
}

func runMerge(t *testing.T, j joins, root, home string, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := mergeWith(j, root, home, args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// mergedRecord is the single stdout line the verb's record grammar owns. The lane writes
// its own line to the same stream, so the reader selects rather than compares the whole.
func mergedRecord(t *testing.T, stdout string) string {
	t.Helper()
	var found []string
	for _, line := range strings.Split(strings.TrimSuffix(stdout, "\n"), "\n") {
		if strings.HasPrefix(line, "merged{") {
			found = append(found, line)
		}
	}
	if len(found) != 1 {
		t.Fatalf("merged records = %d in %q, want exactly one", len(found), stdout)
	}
	return found[0]
}

// WM1: a diverged target publishes the merge-tree of the pair, and its checkout follows
// the ref. A verb that moves only the ref leaves the checkout at the previous tip.
func TestMergePublishesTheMergeTreeAndMovesTheCheckout(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch)
	if tip == previous {
		t.Fatal("the branch tip did not move")
	}
	want := gitOutput(t, root, "merge-tree", "--write-tree", previous, incoming)
	if got := gitOutput(t, root, "rev-parse", tip+"^{tree}"); got != want {
		t.Fatalf("published tree = %s, want the merge-tree of the pair %s", got, want)
	}
	if head := gitOutput(t, target.Path, "rev-parse", "HEAD"); head != tip {
		t.Fatalf("checkout HEAD = %s, want the published commit %s", head, tip)
	}
	// HEAD alone follows the branch ref a linked worktree is attached to, so it reads the
	// same whether the checkout moved or not. The status is what separates them.
	if status := gitOutput(t, target.Path, "status", "--porcelain=v1", "--untracked-files=all"); status != "" {
		t.Fatalf("checkout status = %q, want the checkout reconciled to the published commit", status)
	}
	if _, err := os.Stat(filepath.Join(target.Path, "incoming.txt")); err != nil {
		t.Fatalf("incoming.txt is absent from the checkout: %v", err)
	}
}

// WM2: a target behind the incoming commit moves by fast-forward. A merge-always
// implementation mints a redundant merge commit.
func TestMergeFastForwardsATargetBehindTheIncomingCommit(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "delegate")
	target := created[0]
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if record := mergedRecord(t, stdout); !strings.Contains(record, "kind=fast-forward") {
		t.Fatalf("record = %q, want kind=fast-forward", record)
	}
	if tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch); tip != incoming {
		t.Fatalf("branch tip = %s, want the incoming commit %s with no new object", tip, incoming)
	}
}

// WM3: a target that already contains the incoming commit is idempotent. A verb that
// reruns the lane or mints an empty merge changes state on a re-run.
func TestMergeReportsCurrentAndChangesNothing(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration")
	target := created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	contained := gitOutput(t, root, "rev-parse", "HEAD")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", contained, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if record := mergedRecord(t, stdout); !strings.Contains(record, "kind=current") {
		t.Fatalf("record = %q, want kind=current", record)
	}
	if tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch); tip != previous {
		t.Fatalf("branch tip = %s, want it unchanged at %s", tip, previous)
	}
	if head := gitOutput(t, target.Path, "rev-parse", "HEAD"); head != previous {
		t.Fatalf("checkout HEAD = %s, want it unchanged at %s", head, previous)
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("the lane ran on a current target: %v", err)
	}
}

// WM5: the verb derives the subject, so no `-m` exists and the log reads one way.
func TestMergeDerivesThePublishedSubject(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")

	if code, stdout, stderr := runMerge(t, j, root, home, "--from", "main", target.Assignment.ID); code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch)
	want := "merge: compose main " + incoming[:8] + " into " + target.Assignment.Label
	if got := gitOutput(t, root, "log", "-1", "--format=%s", tip); got != want {
		t.Fatalf("subject = %q, want %q", got, want)
	}
}

// WM6: `--from <sibling label>` folds the sibling's committed branch tip. A from-resolver
// that takes only commits refuses every sibling.
func TestMergeFoldsASiblingBranchTip(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration", "delegate")
	target, sibling := created[0], created[1]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	siblingTip := gitOutput(t, sibling.Path, "rev-parse", "HEAD")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", sibling.Assignment.Label, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch)
	if first := gitOutput(t, root, "rev-parse", tip+"^1"); first != previous {
		t.Fatalf("first parent = %s, want the previous tip %s", first, previous)
	}
	if second := gitOutput(t, root, "rev-parse", tip+"^2"); second != siblingTip {
		t.Fatalf("second parent = %s, want the sibling's branch tip %s", second, siblingTip)
	}
}

// WM7: the record names every identity the next phase reads, so no caller runs `git log`
// to find the new tip.
func TestMergeRecordNamesEveryIdentity(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch)
	want := "merged{worktree=" + target.Assignment.ID + ",from=" + incoming +
		",kind=merge,previous_tip=" + previous + ",tip=" + tip +
		",tree=" + gitOutput(t, root, "rev-parse", tip+"^{tree}") + "}"
	if got := mergedRecord(t, stdout); got != want {
		t.Fatalf("record = %q, want %q", got, want)
	}
}

// WM8: the target operand takes every address `exec` takes, and an ambiguous prefix
// refuses naming both ids. A path-only operand breaks the `exec`-style address.
func TestMergeTargetOperandTakesEveryAddress(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "shared-prefix-one", "shared-prefix-two")
	target := created[0]
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")

	for _, operand := range []string{
		target.Assignment.Label,
		target.Assignment.ID,
		target.Assignment.ID[:10],
		target.Path,
	} {
		code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, operand)
		if code != 0 {
			t.Fatalf("merge %q exit = %d, want 0; stdout=%q stderr=%q", operand, code, stdout, stderr)
		}
		if record := mergedRecord(t, stdout); !strings.Contains(record, "worktree="+target.Assignment.ID+",") {
			t.Fatalf("merge %q record = %q, want the target's assignment id", operand, record)
		}
	}
	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, "shared-p")
	if code != 1 || stderr != "" {
		t.Fatalf("ambiguous prefix = (%d, %q), want a refusal on stdout alone", code, stderr)
	}
	for _, id := range []string{created[0].Assignment.ID, created[1].Assignment.ID} {
		if !strings.Contains(stdout, id) {
			t.Fatalf("ambiguity refusal = %q, want it to name %s", stdout, id)
		}
	}
}

// WM18: the declared lane grades the composed tree before the ref moves. A verb that
// skips the lane publishes a broken build.
func TestMergeRunsTheDeclaredLaneOnTheComposedTree(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration")
	target := created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "lane{outcome=pass,checks=unit}") {
		t.Fatalf("stdout = %q, want the lane's pass line", stdout)
	}
	if recorded, err := os.ReadFile(tally); err != nil || string(recorded) != "g" {
		t.Fatalf("lane tally = %q, %v; want exactly one lane run", recorded, err)
	}
}

// WM32: a missing, empty, or over-supplied operand is invalid usage, not a merge. A
// defaulted `--from` would merge the default branch silently.
func TestMergeRefusesInvalidUsage(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	for name, args := range map[string][]string{
		"missing --from":     {target.Assignment.ID},
		"empty --from":       {"--from", "", target.Assignment.ID},
		"second positional":  {"--from", "main", target.Assignment.ID, target.Assignment.ID},
		"missing positional": {"--from", "main"},
	} {
		t.Run(name, func(t *testing.T) {
			code, stdout, stderr := runMerge(t, j, root, home, args...)
			if code != 2 || stdout != "" || stderr == "" {
				t.Fatalf("%s = (%d, %q, %q), want exit 2 with usage on stderr alone", name, code, stdout, stderr)
			}
		})
	}
}
