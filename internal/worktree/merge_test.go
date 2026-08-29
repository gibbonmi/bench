// The worktree merge verb's command-seam tests: composition, publication, and the
// advertised operand grammar, driven in process against real fixture repositories.
package worktree

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/intent"
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

// --- refusals before publication ---

// commitCapture commits one capture-directory file at the mode the caller names, so a
// fixture repository tracks the capture file this repository git-ignores.
func commitCapture(t *testing.T, dir, name, body string, mode os.FileMode, message string) {
	t.Helper()
	mustMkdirAll(t, filepath.Dir(filepath.Join(dir, name)), 0o755)
	mustWrite(t, filepath.Join(dir, name), []byte(body), mode)
	gitRun(t, dir, "add", name)
	gitRun(t, dir, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", message)
}

// requireMergeRefusal pins the refusal surface every pre-publication refusal shares: exit
// 1, the `refused{` record on stdout alone, and the fragments the row names.
func requireMergeRefusal(t *testing.T, code int, stdout, stderr string, fragments ...string) {
	t.Helper()
	if code != 1 || !strings.HasPrefix(stdout, "refused{") || stderr != "" {
		t.Fatalf("refusal = (%d, %q, %q), want exit 1 with refused{ on stdout alone", code, stdout, stderr)
	}
	for _, fragment := range fragments {
		if !strings.Contains(stdout, fragment) {
			t.Fatalf("refusal = %q, want it to carry %q", stdout, fragment)
		}
	}
}

// requireMergeUnchanged pins what a refusal leaves behind: the branch tip, the checkout
// HEAD, a clean checkout, and no lane record.
func requireMergeUnchanged(t *testing.T, root, path, branch, previous, tally string) {
	t.Helper()
	if tip := gitOutput(t, root, "rev-parse", branch); tip != previous {
		t.Fatalf("branch tip = %s, want it unchanged at %s", tip, previous)
	}
	if head := gitOutput(t, path, "rev-parse", "HEAD"); head != previous {
		t.Fatalf("checkout HEAD = %s, want it unchanged at %s", head, previous)
	}
	if _, err := os.Stat(tally); !os.IsNotExist(err) {
		t.Fatalf("the refusal ran the lane: %v", err)
	}
}

// WM9: a conflicting non-capture path refuses with the path table and changes nothing. A
// verb that writes conflict markers leaves the worktree mid-merge.
func TestMergeRefusesAConflictingNonCapturePath(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration")
	target := created[0]
	commitInWorktree(t, target.Path, "tracked.txt", "target edit\n", "target edit")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	incoming := commitOnDefault(t, root, "tracked.txt", "incoming edit\n")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "composition conflict: textual", "paths_total=1", "tracked.txt")
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
	if status := gitOutput(t, target.Path, "status", "--porcelain=v1", "--untracked-files=all"); status != "" {
		t.Fatalf("checkout status = %q, want no conflict residue", status)
	}
}

// WM10: a conflicted `capture/learnings.md` composes as the union and the verb discloses
// the side it took. A verb that bypasses the rule table refuses what the landing settles.
func TestMergeSettlesAConflictedCaptureJournalAsTheUnion(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	commitCapture(t, target.Path, "capture/learnings.md", "target entry\n", 0o644, "target journal")
	commitCapture(t, root, "capture/learnings.md", "incoming entry\n", 0o644, "incoming journal")
	incoming := gitOutput(t, root, "rev-parse", "HEAD")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "merge composition{resolved=capture/learnings.md:union}") {
		t.Fatalf("stderr = %q, want the union disclosure", stderr)
	}
	settled, err := os.ReadFile(filepath.Join(target.Path, "capture/learnings.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range []string{"target entry", "incoming entry"} {
		if !strings.Contains(string(settled), entry) {
			t.Fatalf("settled journal = %q, want it to keep %q", settled, entry)
		}
	}
}

// WM11: a dirty target refuses before composition and names the path. A reconcile over a
// dirty checkout discards the edit.
func TestMergeRefusesADirtyTarget(t *testing.T) {
	t.Parallel()
	for name, dirty := range map[string]func(t *testing.T, path string){
		"untracked": func(t *testing.T, path string) {
			mustWrite(t, filepath.Join(path, "untracked.txt"), []byte("scratch\n"), 0o644)
		},
		"modified": func(t *testing.T, path string) {
			mustWrite(t, filepath.Join(path, "tracked.txt"), []byte("edited\n"), 0o644)
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			j, root, home, tally, created := mergeFixture(t, "integration")
			target := created[0]
			previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
			incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
			dirty(t, target.Path)

			code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
			want := "untracked.txt"
			if name == "modified" {
				want = "tracked.txt"
			}
			requireMergeRefusal(t, code, stdout, stderr, "merge target checkout is not clean", want)
			if tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch); tip != previous {
				t.Fatalf("branch tip = %s, want it unchanged at %s", tip, previous)
			}
			if _, err := os.Stat(tally); !os.IsNotExist(err) {
				t.Fatalf("the refusal ran the lane: %v", err)
			}
		})
	}
}

// WM12: a sibling contributes its committed branch tip alone, so a dirty sibling names
// `bench commit` at the sibling and a detached sibling names its assignment branch ref. A
// fold of the branch tip would silently omit the uncommitted delegate work.
func TestMergeRefusesADirtyOrDetachedSibling(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration", "delegate")
	target, sibling := created[0], created[1]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	mustWrite(t, filepath.Join(sibling.Path, "sibling.txt"), []byte("uncommitted\n"), 0o644)

	code, stdout, stderr := runMerge(t, j, root, home, "--from", sibling.Assignment.Label, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "sibling checkout is not clean",
		"next=bench worktree exec "+sibling.Assignment.ID+" -- bench commit", "sibling.txt")
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)

	gitRun(t, sibling.Path, "checkout", "-q", "--", "sibling.txt")
	gitRun(t, sibling.Path, "checkout", "-q", "--detach", "HEAD")
	code, stdout, stderr = runMerge(t, j, root, home, "--from", sibling.Assignment.Label, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, sibling.Assignment.Branch)
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
}

// WM13: a `--from` both lookups answer is ambiguous, and the refusal names the assignment
// id and the full commit. A first-match resolver merges whichever lookup runs first.
func TestMergeRefusesAnAmbiguousFrom(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration", "delegate")
	target, sibling := created[0], created[1]
	commitInWorktree(t, sibling.Path, "sibling.txt", "sibling\n", "sibling work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	gitRun(t, root, "branch", sibling.Assignment.Label, "HEAD")
	commit := gitOutput(t, root, "rev-parse", "HEAD")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", sibling.Assignment.Label, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "--from names both an assignment and a commit",
		"observed="+sibling.Assignment.Label, "wanted="+sibling.Assignment.ID+" or "+commit)
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
}

// WM14: a `--from` that resolves to the target's own assignment refuses. A self-merge
// prints `kind=current` and hides the operand error.
func TestMergeRefusesAFromThatNamesTheTarget(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration")
	target := created[0]
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", target.Assignment.Label, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "--from resolves to the target itself", "observed="+target.Assignment.ID)
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
}

// WM15: a `--from` that names nothing refuses naming the value, rather than reporting a
// silent `kind=current` on a typo.
func TestMergeRefusesAnUnknownFrom(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration")
	target := created[0]
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", "no-such-thing", target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "--from names no assignment and no commit", "observed=no-such-thing")
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
}

// WM16: the identity bundle is the verb's whole authority, so each failed component names
// itself and a checkout off its assignment branch names the branch ref. A verb that skips
// the bundle writes to a foreign tree.
func TestMergeRefusesAFailedTargetIdentityComponent(t *testing.T) {
	t.Parallel()
	for _, component := range []string{componentOwnerMarker, componentLock, componentAssignmentState} {
		t.Run(component, func(t *testing.T) {
			t.Parallel()
			j, root, home, tally, created := mergeFixture(t, "integration")
			target := created[0]
			previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
			incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
			fixture := identityComponentFixtureFor(t, component)
			fixture.mutate(t, root, target)

			code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
			requireMergeRefusal(t, code, stdout, stderr, fixture.want(target, "", ""))
			requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
		})
	}
	t.Run("off-branch", func(t *testing.T) {
		t.Parallel()
		j, root, home, tally, created := mergeFixture(t, "integration")
		target := created[0]
		previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
		incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
		gitRun(t, target.Path, "checkout", "-q", "--detach", "HEAD")

		code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
		requireMergeRefusal(t, code, stdout, stderr, target.Assignment.Branch)
		requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
	})
}

// WM17: the compare-and-swap old value is the branch tip, so a branch tip that is not the
// checkout HEAD refuses naming both commits. The pair cannot diverge inside one
// repository, so the check reads a root whose view of the branch is another commit —
// exactly the race the check guards.
func TestMergeRefusesABranchTipThatIsNotTheCheckoutHead(t *testing.T) {
	t.Parallel()
	_, _, _, _, created := mergeFixture(t, "integration")
	target := created[0]
	head := gitOutput(t, target.Path, "rev-parse", "HEAD")
	other := newWorktreeRepo(t)
	commitOnDefault(t, other, "elsewhere.txt", "elsewhere\n")
	gitRun(t, other, "update-ref", target.Assignment.Branch, gitOutput(t, other, "rev-parse", "HEAD"))
	elsewhere := gitOutput(t, other, "rev-parse", target.Assignment.Branch)
	if elsewhere == head {
		t.Fatalf("fixture premise failed: both roots report %s", head)
	}

	_, err := mergeTargetTip(other, target.Assignment)
	var typed refusalError
	if !errors.As(err, &typed) {
		t.Fatalf("mergeTargetTip error = %v, want a refusal", err)
	}
	if typed.detail != "merge target branch tip is not the checkout HEAD" || typed.observed != head || typed.wanted != elsewhere {
		t.Fatalf("refusal = %+v, want both commits named", typed.refusal)
	}
}

// WM31: a conflicting path that holds a control byte renders escaped, so no raw control
// byte splits the refusal record.
func TestMergeRefusalEscapesAControlBytePath(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	name := "board\x1bfile.md"
	commitInWorktree(t, target.Path, name, "target bytes\n", "target control-byte path")
	incoming := commitOnDefault(t, root, name, "incoming bytes\n")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "composition conflict: textual", `"board\\u001bfile.md"`)
	if strings.ContainsRune(stdout, '\x1b') {
		t.Fatalf("refusal = %q, want no raw control byte", stdout)
	}
}

// WM33: two sides of an add/add capture path that disagree on the file mode leave no
// settled mode, so the rule table cannot settle it and the refusal names the path.
func TestMergeRefusesACaptureAddAddWithDisagreeingModes(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration")
	target := created[0]
	commitCapture(t, target.Path, "capture/learnings.md", "target entry\n", 0o644, "target journal")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	commitCapture(t, root, "capture/learnings.md", "incoming entry\n", 0o755, "incoming journal")
	incoming := gitOutput(t, root, "rev-parse", "HEAD")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "composition conflict: mode", "capture/learnings.md")
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
}

// WM34: only the default branch's history and a sibling tip authorize the lane, so an
// off-branch commit refuses naming the commit and no lane executes an unowned tree.
func TestMergeRefusesAnOffBranchFromCommit(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration")
	target := created[0]
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	gitRun(t, root, "checkout", "-q", "-b", "unowned")
	offBranch := commitOnDefault(t, root, "off.txt", "off branch\n")
	gitRun(t, root, "checkout", "-q", "main")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", offBranch, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr,
		"--from is outside the default branch's history and is no sibling tip", "observed="+offBranch)
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
}

// --- the publication boundary ---

// mergeLaneOf replaces the fixture's lane with the checks one row needs, so a row
// controls the outcome the boundary reacts to without a second fixture.
func mergeLaneOf(checks ...gate.Phase) func(string) ([]gate.Phase, string, error) {
	return func(string) ([]gate.Phase, string, error) { return checks, "", nil }
}

// requireMergeLaneRefusal pins the surface a lane fail leaves: exit 1, the lane's own
// fail line naming the check, and the refusal record after it.
func requireMergeLaneRefusal(t *testing.T, code int, stdout, check string) {
	t.Helper()
	if code != 1 {
		t.Fatalf("merge exit = %d, want 1; stdout=%q", code, stdout)
	}
	fail := "lane{outcome=fail,check=" + check + "}"
	if !strings.Contains(stdout, fail) {
		t.Fatalf("stdout = %q, want %q", stdout, fail)
	}
	if index := strings.Index(stdout, "refused{"); index < 0 || index < strings.Index(stdout, fail) {
		t.Fatalf("stdout = %q, want the refusal record after the lane's fail line", stdout)
	}
}

// WM19: a failing lane check names itself at exit 1, and the tip and the checkout stay
// unchanged. A verb that publishes on a lane fail ships an ungraded tree.
func TestMergeRefusesAFailingLaneCheck(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
	j.mergeLane = mergeLaneOf(gate.Phase{Name: "unit", Argv: []string{"sh", "-c", "echo the lane says no; exit 1"}})

	code, stdout, _ := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	requireMergeLaneRefusal(t, code, stdout, "unit")
	if tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch); tip != previous {
		t.Fatalf("branch tip = %s, want it unchanged at %s", tip, previous)
	}
	if status := gitOutput(t, target.Path, "status", "--porcelain=v1", "--untracked-files=all"); status != "" {
		t.Fatalf("checkout status = %q, want the checkout untouched", status)
	}
	if _, err := os.Stat(filepath.Join(target.Path, "incoming.txt")); !os.IsNotExist(err) {
		t.Fatalf("incoming.txt reached the checkout of a refused merge: %v", err)
	}
}

// WM23: a fast-forward grades the incoming tree too, so a commit whose tree fails the
// lane refuses and the tip stays. A lane skipped on fast-forward publishes an ungraded
// tree.
func TestMergeRefusesAFastForwardTheLaneFails(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "delegate")
	target := created[0]
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
	// The check reads the graded tree's own checkout, so it passes on the previous tip
	// and fails on the incoming one.
	j.mergeLane = mergeLaneOf(gate.Phase{Name: "unit", Argv: []string{"sh", "-c", "test ! -f incoming.txt"}})

	code, stdout, _ := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	requireMergeLaneRefusal(t, code, stdout, "unit")
	if tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch); tip != previous {
		t.Fatalf("branch tip = %s, want it unchanged at %s", tip, previous)
	}
}

// WM22: a checkout edited between the pre-check and the ref update refuses before the ref
// moves, and the edit survives. A verb that trusts the first read resets over a fresh
// edit.
func TestMergeRefusesACheckoutEditedDuringTheLane(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
	fresh := filepath.Join(target.Path, "fresh.txt")
	j.mergeLane = mergeLaneOf(gate.Phase{Name: "unit", Argv: []string{"sh", "-c",
		"printf 'fresh\\n' > " + sanitize.ShellQuote(fresh)}})

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	// The lane ran and passed, so its own line precedes the refusal record on stdout.
	if code != 1 || stderr != "" {
		t.Fatalf("merge = (%d, %q), want exit 1 with stderr empty; stdout=%q", code, stderr, stdout)
	}
	if !strings.Contains(stdout, "refused{detail=merge target checkout changed}") {
		t.Fatalf("stdout = %q, want the fingerprint recheck's refusal", stdout)
	}
	if tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch); tip != previous {
		t.Fatalf("branch tip = %s, want it unchanged at %s", tip, previous)
	}
	if body, err := os.ReadFile(fresh); err != nil || string(body) != "fresh\n" {
		t.Fatalf("fresh.txt = %q, %v; want the edit to survive the refusal", body, err)
	}
}

// WM21: a reconcile failure after the ref update exits 3 and names the repair at the
// published commit. A refusal-shaped exit hides that the ref moved.
func TestMergeExitsThreeWhenTheReconcileFails(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
	j.mergeReconcile = func(string, string) error { return errors.New("reset refused") }

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	if code != 3 {
		t.Fatalf("merge exit = %d, want 3; stdout=%q stderr=%q", code, stdout, stderr)
	}
	tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch)
	if tip == previous {
		t.Fatal("the exit-3 record reports a publication the ref never took")
	}
	record := mergedRecord(t, stdout)
	if !strings.Contains(record, ",tip="+tip+",") {
		t.Fatalf("record = %q, want it to name the published tip %s", record, tip)
	}
	if !strings.Contains(record, "next=git -C ") || !strings.Contains(record, target.Path) {
		t.Fatalf("record = %q, want the repair at the target checkout %s", record, target.Path)
	}
	if !strings.HasSuffix(record, " reset --merge "+tip+"}") {
		t.Fatalf("record = %q, want the repair to name the published commit %s", record, tip)
	}
	// A linked worktree's HEAD follows the branch ref whether or not the checkout was
	// reconciled, so the unreconciled state reads through the status and the tree.
	if status := gitOutput(t, target.Path, "status", "--porcelain=v1", "--untracked-files=all"); status == "" {
		t.Fatal("checkout status is clean, so the fixture never left the reconcile undone")
	}
	if _, err := os.Stat(filepath.Join(target.Path, "incoming.txt")); !os.IsNotExist(err) {
		t.Fatalf("incoming.txt is in the checkout, so the reconcile ran: %v", err)
	}
}

// WM24: the lane's prose placeholder is exactly the Markdown that differs between the
// previous tip's tree and the composed tree. An empty placeholder grades no incoming
// prose, and a whole-tree placeholder grades unchanged files.
func TestMergeResolvesTheProsePlaceholderToTheIncomingMarkdown(t *testing.T) {
	t.Parallel()
	j, root, home, _, created := mergeFixture(t, "integration")
	target := created[0]
	// README.md rides the base commit unchanged; the target side changes one Markdown
	// file the composed tree shares with the previous tip; the incoming side changes one.
	commitInWorktree(t, target.Path, "target-only.md", "target prose\n", "target prose")
	incoming := commitOnDefault(t, root, "incoming-only.md", "incoming prose\n")
	argv := filepath.Join(t.TempDir(), "prose-argv")
	j.mergeLane = mergeLaneOf(gate.Phase{Name: "prose", Argv: []string{"sh", "-c",
		`for path in "$@"; do printf '%s\n' "$path"; done > ` + sanitize.ShellQuote(argv),
		"prose", gate.LaneNamedMarkdownToken}})

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	recorded, err := os.ReadFile(argv)
	if err != nil {
		t.Fatalf("the prose check recorded no argv: %v", err)
	}
	if got := string(recorded); got != "incoming-only.md\n" {
		t.Fatalf("prose argv = %q, want the incoming side's Markdown alone", got)
	}
}

// --- the review repairs ---

// WM35: an ambiguous `--from` prefix refuses naming every matching assignment id, whether
// or not a branch of the same spelling exists. A resolver that swallows the ambiguity
// merges whichever object happens to carry the name.
func TestMergeRefusesAnAmbiguousFromPrefix(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "shared-prefix-one", "shared-prefix-two", "integration")
	target := created[2]
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")

	for _, branched := range []bool{false, true} {
		if branched {
			gitRun(t, root, "branch", "shared-p", "HEAD")
		}
		code, stdout, stderr := runMerge(t, j, root, home, "--from", "shared-p", target.Assignment.ID)
		requireMergeRefusal(t, code, stdout, stderr, "target is ambiguous",
			created[0].Assignment.ID, created[1].Assignment.ID)
		requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
	}
}

// WM37: the sibling lookup runs over active assignments only, so a retired assignment's
// label names no sibling and falls through to the commit lookup. A lookup over every
// state refuses a legitimate default-branch commit by the assignment's state.
func TestMergeResolvesARetiredSiblingLabelThroughTheCommitLookup(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration", "retired")
	target, retired := created[0], created[1].Assignment
	commitInWorktree(t, target.Path, "target.txt", "target\n", "target work")
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	retired.State = intent.StateComplete
	if err := intent.PutAssignment(root, retired); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := runMerge(t, j, root, home, "--from", retired.Label, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "--from names no assignment and no commit", "observed="+retired.Label)
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)

	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
	gitRun(t, root, "branch", retired.Label, incoming)
	code, stdout, stderr = runMerge(t, j, root, home, "--from", retired.Label, target.Assignment.ID)
	if code != 0 {
		t.Fatalf("merge exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	tip := gitOutput(t, root, "rev-parse", target.Assignment.Branch)
	if second := gitOutput(t, root, "rev-parse", tip+"^2"); second != incoming {
		t.Fatalf("second parent = %s, want the commit the branch names %s", second, incoming)
	}
}

// WM38: an unresolved default branch is a failed query, not a classification, so the
// commit lookup refuses naming the query. A fold into `owned=false` reports the commit as
// outside a history the verb never read.
func TestMergeRefusesAnUnresolvedDefaultBranch(t *testing.T) {
	t.Parallel()
	j, root, home, tally, created := mergeFixture(t, "integration")
	target := created[0]
	previous := gitOutput(t, target.Path, "rev-parse", "HEAD")
	incoming := commitOnDefault(t, root, "incoming.txt", "incoming\n")
	// Two local branches remain and none of them is the `main` candidate, so no default
	// branch resolves and the ancestry the lookup needs has no subject.
	gitRun(t, root, "branch", "-m", "main", "trunk")

	code, stdout, stderr := runMerge(t, j, root, home, "--from", incoming, target.Assignment.ID)
	requireMergeRefusal(t, code, stdout, stderr, "default branch is unresolved")
	requireMergeUnchanged(t, root, target.Path, target.Assignment.Branch, previous, tally)
}

// S1: the landing destination and the merge target read one checkout-clean predicate, so
// an unreadable status names the failed read at both. The destination's own derivation
// reported an unreadable status as a dirty destination and named no path.
func TestLandingDestinationNamesAnUnreadableStatusThroughTheSharedPredicate(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	// A corrupt index fails `git status` while every ref read the destination proof runs
	// before it still answers.
	mustWrite(t, filepath.Join(root, ".git", "index"), []byte("not an index\n"), 0o644)

	_, _, _, _, err := landingDestination(defaultJoins(), root)
	if err == nil || !strings.Contains(err.Error(), "checkout status is unreadable") {
		t.Fatalf("landing destination error = %v, want the shared predicate's unreadable-status refusal", err)
	}
}
