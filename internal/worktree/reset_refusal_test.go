package worktree

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func runReset(t *testing.T, root, home string, args ...string) (int, string, string) {
	t.Helper()
	return runResetWith(t, defaultJoins(), root, home, args...)
}

func runResetWith(t *testing.T, j joins, root, home string, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := resetWith(j, root, home, args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestResetRefusesNeitherMode(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".git", intent.Filename), []byte("{"), 0o600)
	code, out, errout := runReset(t, root, t.TempDir(), "target")
	requireTest(t, code == 2 && strings.Contains(out+errout, "bench worktree reset"),
		"missing mode = %d %s %s", code, out, errout)
	requireTest(t, !strings.Contains(out+errout, "ledger"), "grammar read the ledger: %s %s", out, errout)
}

func TestResetRefusesAControlByteCheckpoint(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".git", intent.Filename), []byte("{"), 0o600)
	code, out, _ := runReset(t, root, t.TempDir(), "--to", "bad\x1bvalue", "target")
	requireTest(t, code == 1 && strings.Contains(out, "--to contains control characters"), "control checkpoint = %d %s", code, out)
}

func requireResetRefusal(t *testing.T, root, home, to, target, detail string) string {
	t.Helper()
	code, out, errout := runReset(t, root, home, "--to", to, target)
	requireTest(t, code == 1 && strings.Contains(out, detail), "refusal wants %q: %d %s %s", detail, code, out, errout)
	return out
}

func TestResetRefusesANonCommitCheckpoint(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-noncommit")
	requireResetRefusal(t, root, home, "not-a-commit", creation.Assignment.ID, "checkpoint is not a commit")
}

func TestResetRefusesACheckpointOutsideTheHistory(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-range")
	commitInWorktree(t, root, "default-only", "new\n", "default ahead")
	checkpoint := gitOutput(t, root, "rev-parse", "HEAD")
	out := requireResetRefusal(t, root, home, checkpoint, creation.Assignment.ID, "checkpoint is outside the assignment history")
	requireTest(t, strings.Contains(out, "wanted="+creation.Assignment.Start+".."+creation.Assignment.Start), "missing history range: %s", out)
}

func TestResetRefusesAShiftBranchCheckpoint(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-shift-checkpoint")
	gitRun(t, creation.Path, "switch", "-c", "bench/shift-reset-checkpoint")
	commitInWorktree(t, creation.Path, "shift-only", "new\n", "shift ahead")
	requireResetRefusal(t, root, home, gitOutput(t, creation.Path, "rev-parse", "HEAD"),
		creation.Assignment.ID, "checkpoint is outside the assignment history")
}

func TestResetRefusesThePrimaryCheckout(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-primary")
	requireResetRefusal(t, root, home, creation.Assignment.Start, root, "target is unassigned")
}

func TestResetRefusesARetiredAssignment(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-retired")
	assignment := creation.Assignment
	assignment.State = intent.StateComplete
	mustNoError(t, intent.PutAssignment(root, assignment))
	requireResetRefusal(t, root, home, assignment.Start, assignment.Label, "assignment "+assignment.ID+" is not active")
}

func TestResetRefusesAConflictedIndex(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-conflict")
	setupConflict(t, creation.Path)
	out := requireResetRefusal(t, root, home, creation.Assignment.Start, creation.Assignment.ID, "checkout is conflicted")
	requireTest(t, strings.Contains(out, "next=bench worktree clean "+creation.Assignment.ID), "missing cleanup route: %s", out)
}

func TestResetRefusesADirtyNestedRepository(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedSubmoduleAssignment(t, "reset-submodule")
	mustWrite(t, filepath.Join(creation.Path, "sub", "sub.txt"), []byte("dirty\n"), 0o644)
	requireResetRefusal(t, root, home, creation.Assignment.Start, creation.Assignment.ID, "nested repository is dirty")
}

func resetEmbedded(t *testing.T, path string) string {
	t.Helper()
	nested := filepath.Join(path, "nested")
	mustMkdirAll(t, nested, 0o755)
	gitRun(t, nested, "init", "-q", "-b", "main")
	commitInWorktree(t, nested, "nested.txt", "base\n", "nested")
	return nested
}

func TestResetRefusesAnEmbeddedRepository(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-embedded")
	resetEmbedded(t, creation.Path)
	requireResetRefusal(t, root, home, creation.Assignment.Start, creation.Assignment.ID, "embedded repository is retained")
}

func TestResetRefusesADirtyEmbeddedRepository(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-dirty-embedded")
	nested := resetEmbedded(t, creation.Path)
	mustWrite(t, filepath.Join(nested, "nested.txt"), []byte("dirty\n"), 0o644)
	requireResetRefusal(t, root, home, creation.Assignment.Start, creation.Assignment.ID, "embedded repository is retained")
}

func TestResetRefusesAnUnknownNestedState(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-unknown-nested")
	nested := filepath.Join(creation.Path, "broken")
	mustMkdirAll(t, nested, 0o755)
	mustWrite(t, filepath.Join(nested, ".git"), []byte("gitdir: missing\n"), 0o644)
	requireResetRefusal(t, root, home, creation.Assignment.Start, creation.Assignment.ID, "nested repository state is unknown")
}

func TestResetRefusesALiveLease(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-lease")
	lease, err := LeaseFile(creation.Path)
	mustNoError(t, err)
	mustWrite(t, lease, []byte(fmt.Sprintf("%d 2026-09-11T00:00:00Z\n", os.Getppid())), 0o600)
	requireResetRefusal(t, root, home, creation.Assignment.Start, creation.Assignment.ID, "assignment has a live lease")
}

func TestResetRefusesAnAmbiguousCheckpoint(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-ambiguous")
	tree := gitOutput(t, root, "rev-parse", "HEAD^{tree}")
	seen := map[string]string{}
	for nonce := 0; nonce <= 65536; nonce++ {
		body := fmt.Sprintf("tree %s\nparent %s\nauthor bench <bench@local> 1 +0000\ncommitter bench <bench@local> 1 +0000\n\n%d\n",
			tree, creation.Assignment.Start, nonce)
		digest := fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("commit %d%c%s", len(body), byte(0), body))))
		prefix := digest[:4]
		if previous, found := seen[prefix]; found {
			for _, commit := range []string{previous, body} {
				_, err := gitInput(root, nil, []byte(commit), "hash-object", "-t", "commit", "-w", "--stdin")
				mustNoError(t, err)
			}
			requireTest(t, len(strings.Fields(gitOutput(t, root, "rev-parse", "--disambiguate="+prefix))) >= 2,
				"fixture did not create an ambiguous prefix")
			requireResetRefusal(t, root, home, prefix, creation.Assignment.ID, "checkpoint is not a commit")
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
			root, creation, home := newOwnedAssignment(t, "reset-hidden"+flag)
			gitRun(t, creation.Path, "update-index", flag, "README.md")
			mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("hidden edit\n"), 0o644)
			gitRun(t, creation.Path, "switch", "--detach", "HEAD")
			out := requireResetRefusal(t, root, home, creation.Assignment.Start, creation.Assignment.ID, "index carries hidden flags")
			requireTest(t, strings.Contains(out, "refusal_paths[1]{path}:\n  README.md\n"), "hidden flag paths = %s", out)
			body, err := os.ReadFile(filepath.Join(creation.Path, "README.md"))
			mustNoError(t, err)
			requireTest(t, string(body) == "hidden edit\n", "hidden edit changed: %q", body)
		})
	}
}

func TestResetApplyRefusesAFingerprintForANonePlan(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-apply")
	code, out, _ := runReset(t, root, home, "--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", "fingerprint")
	requireTest(t, code == 1 && strings.Contains(out, "reset plan is stale") && strings.Contains(out, "wanted=none") && !strings.Contains(out, "reset_plan"),
		"none-plan apply = %d %s", code, out)
}
