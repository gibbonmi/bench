package worktree

import (
	"bytes"
	"path/filepath"
	"testing"
)

// A path-addressed clean is destructive. When its operand names nothing this repo can
// act on, the caller must be able to tell that from the exit code alone: an agent that
// scripts cleanup reads a zero as "planned" and moves on, so a silent no-op there loses
// the work the operator meant to preserve.
func TestCleanCommandRefusesAnOperandItCannotResolve(t *testing.T) {
	root, _ := newOwnedAssignment(t, "operand")
	for _, tc := range []struct {
		name, target, detail string
	}{
		{"unregistered", filepath.Join(t.TempDir(), "absent"), "target is not registered"},
		{"tilde-prefixed", "~/.bench/worktrees/absent", "target is not registered"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(root)
			var stdout, stderr bytes.Buffer
			code := CleanCommand([]string{tc.target}, &stdout, &stderr)
			if code == 0 || !bytes.Contains(stdout.Bytes(), []byte(tc.detail)) {
				t.Fatalf("%s clean code=%d stdout=%q", tc.name, code, stdout.String())
			}
		})
	}
}

// A resolved operand that this repo declines to remove is a verdict, not a bad operand:
// the plan is the answer and the exit code stays zero.
func TestCleanCommandReportsAResolvedRetainVerdictAsSuccess(t *testing.T) {
	root, creation := newPendingAssignment(t, "verdict")
	mustWrite(t, filepath.Join(creation.Path, "residual.log"), []byte("x\n"), 0o644)
	mustWrite(t, filepath.Join(creation.Path, ".gitignore"), []byte("residual.log\n"), 0o644)
	gitRun(t, creation.Path, "add", ".gitignore")
	gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore residual")
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	code := CleanCommand([]string{creation.Path}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("resolved retain code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}
