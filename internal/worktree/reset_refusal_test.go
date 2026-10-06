package worktree

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestResetRefusesNeitherMode(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".git", intent.Filename), []byte("{"), 0o600)
	result := runVerb(t, verbReset, verbCall{root: root, home: t.TempDir(), args: []string{"target"}})
	requireTest(t, result.exit == 2 && strings.Contains(result.stdout+result.stderr, "bench worktree reset"),
		"missing mode = %d %s %s", result.exit, result.stdout, result.stderr)
	requireTest(t, !strings.Contains(result.stdout+result.stderr, "ledger"), "grammar read the ledger: %s %s", result.stdout, result.stderr)
}

func TestResetRefusesAControlByteCheckpoint(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".git", intent.Filename), []byte("{"), 0o600)
	result := runVerb(t, verbReset, verbCall{root: root, home: t.TempDir(), args: []string{"--to", "bad\x1bvalue", "target"}})
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "--to contains control characters"), "control checkpoint = %d %s", result.exit, result.stdout)
}

func requireResetRefusal(t *testing.T, root, home, to, target, detail string) string {
	t.Helper()
	result := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", to, target}})
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, detail), "refusal wants %q: %d %s %s", detail, result.exit, result.stdout, result.stderr)
	return result.stdout
}

func TestResetRefusesANonCommitCheckpoint(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-noncommit")
	requireResetRefusal(t, f.root, f.home, "not-a-commit", f.creation.Assignment.ID, "checkpoint is not a commit")
}

func TestResetRefusesACheckpointOutsideTheHistory(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-range")
	commitInWorktree(t, f.root, "default-only", "new\n", "default ahead")
	checkpoint := gitOutput(t, f.root, "rev-parse", "HEAD")
	out := requireResetRefusal(t, f.root, f.home, checkpoint, f.creation.Assignment.ID, "checkpoint is outside the assignment history")
	requireTest(t, strings.Contains(out, "wanted="+f.creation.Assignment.Start+".."+f.creation.Assignment.Start), "missing history range: %s", out)
}

func TestResetRefusesAShiftBranchCheckpoint(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-shift-checkpoint")
	gitRun(t, f.creation.Path, "switch", "-c", "bench/shift-reset-checkpoint")
	commitInWorktree(t, f.creation.Path, "shift-only", "new\n", "shift ahead")
	requireResetRefusal(t, f.root, f.home, gitOutput(t, f.creation.Path, "rev-parse", "HEAD"),
		f.creation.Assignment.ID, "checkpoint is outside the assignment history")
}

func TestResetRefusesThePrimaryCheckout(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-primary")
	requireResetRefusal(t, f.root, f.home, f.creation.Assignment.Start, f.root, "target is unassigned")
}

func TestResetRefusesARetiredAssignment(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-retired")
	assignment := f.creation.Assignment
	assignment.State = intent.StateComplete
	mustNoError(t, intent.PutAssignment(f.root, assignment))
	requireResetRefusal(t, f.root, f.home, assignment.Start, assignment.Label, "assignment "+assignment.ID+" is not active")
}

func TestResetRefusesAConflictedIndex(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-conflict")
	setupConflict(t, f.creation.Path)
	out := requireResetRefusal(t, f.root, f.home, f.creation.Assignment.Start, f.creation.Assignment.ID, "checkout is conflicted")
	requireTest(t, strings.Contains(out, "next=bench worktree clean "+f.creation.Assignment.ID), "missing cleanup route: %s", out)
}

func TestResetRefusesADirtyNestedRepository(t *testing.T) {
	t.Parallel()
	f := newOwnedSubmoduleAssignment(t, "reset-submodule")
	mustWrite(t, filepath.Join(f.creation.Path, "sub", "sub.txt"), []byte("dirty\n"), 0o644)
	requireResetRefusal(t, f.root, f.home, f.creation.Assignment.Start, f.creation.Assignment.ID, "nested repository is dirty")
}

func TestResetRefusesAnEmbeddedRepository(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-embedded")
	plantNestedRepository(t, f.creation.Path)
	requireResetRefusal(t, f.root, f.home, f.creation.Assignment.Start, f.creation.Assignment.ID, "embedded repository is retained")
}

func TestResetRefusesADirtyEmbeddedRepository(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-dirty-embedded")
	nested := plantNestedRepository(t, f.creation.Path)
	mustWrite(t, filepath.Join(nested, nestedRepositoryFile), []byte("dirty\n"), 0o644)
	requireResetRefusal(t, f.root, f.home, f.creation.Assignment.Start, f.creation.Assignment.ID, "embedded repository is retained")
}

func TestResetRefusesAnUnknownNestedState(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-unknown-nested")
	nested := filepath.Join(f.creation.Path, "broken")
	mustMkdirAll(t, nested, 0o755)
	mustWrite(t, filepath.Join(nested, ".git"), []byte("gitdir: missing\n"), 0o644)
	requireResetRefusal(t, f.root, f.home, f.creation.Assignment.Start, f.creation.Assignment.ID, "nested repository state is unknown")
}

func TestResetRefusesALiveLease(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-lease")
	lease, err := LeaseFile(f.creation.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(fmt.Sprintf("%d 2026-09-11T00:00:00Z\n", os.Getppid())), 0o600)
	requireResetRefusal(t, f.root, f.home, f.creation.Assignment.Start, f.creation.Assignment.ID, "assignment has a live lease")
}

func TestResetRefusesAnAmbiguousCheckpoint(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-ambiguous")
	tree := gitOutput(t, f.root, "rev-parse", "HEAD^{tree}")
	seen := map[string]string{}
	for nonce := 0; nonce <= 65536; nonce++ {
		body := fmt.Sprintf("tree %s\nparent %s\nauthor bench <bench@local> 1 +0000\ncommitter bench <bench@local> 1 +0000\n\n%d\n",
			tree, f.creation.Assignment.Start, nonce)
		digest := fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("commit %d%c%s", len(body), byte(0), body))))
		prefix := digest[:4]
		if previous, found := seen[prefix]; found {
			for _, commit := range []string{previous, body} {
				_, err := gitInput(f.root, nil, []byte(commit), "hash-object", "-t", "commit", "-w", "--stdin")
				mustNoError(t, err)
			}
			requireTest(t, len(strings.Fields(gitOutput(t, f.root, "rev-parse", "--disambiguate="+prefix))) >= 2,
				"fixture did not create an ambiguous prefix")
			requireResetRefusal(t, f.root, f.home, prefix, f.creation.Assignment.ID, "checkpoint is not a commit")
			return
		}
		seen[prefix] = body
	}
	t.Fatal("no commit prefix collision")
}

func TestResetRefusesHiddenIndexFlags(t *testing.T) {
	t.Parallel()
	for _, flag := range []string{"--assume-unchanged", "--skip-worktree"} {
		t.Run(flag, func(t *testing.T) {
			t.Parallel()
			f := newOwnedAssignment(t, "reset-hidden"+flag)
			gitRun(t, f.creation.Path, "update-index", flag, "README.md")
			mustWrite(t, filepath.Join(f.creation.Path, "README.md"), []byte("hidden edit\n"), 0o644)
			gitRun(t, f.creation.Path, "switch", "--detach", "HEAD")
			out := requireResetRefusal(t, f.root, f.home, f.creation.Assignment.Start, f.creation.Assignment.ID, "index carries hidden flags")
			requireTest(t, strings.Contains(out, refusalPathsTable+"[1]{path}:\n  README.md\n"), "hidden flag paths = %s", out)
			body, err := os.ReadFile(filepath.Join(f.creation.Path, "README.md"))
			mustNoError(t, err)
			requireTest(t, string(body) == "hidden edit\n", "hidden edit changed: %q", body)
		})
	}
}

func TestResetApplyRefusesAFingerprintForANonePlan(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "reset-apply")
	result := runVerb(t, verbReset, f.call("--to", f.creation.Assignment.Start, f.creation.Assignment.ID, "--apply", "fingerprint"))
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset plan is stale") && strings.Contains(result.stdout, "wanted=none") && !strings.Contains(result.stdout, "reset_plan"),
		"none-plan apply = %d %s", result.exit, result.stdout)
}
