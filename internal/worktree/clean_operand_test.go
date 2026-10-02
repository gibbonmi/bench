package worktree

import (
	"path/filepath"
	"strings"
	"testing"
)

// A path-addressed clean is destructive. When its operand names nothing this repo can
// act on, the caller must be able to tell that from the exit code alone. An agent that
// scripts cleanup reads a zero as "planned" and moves on. A silent no-op there loses
// the work the operator meant to preserve.
func TestCleanCommandRefusesAnOperandItCannotResolve(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "operand")
	for _, tc := range []struct {
		name, target, detail string
	}{
		{"unregistered", filepath.Join(t.TempDir(), "absent"), "target is not registered"},
		{"tilde-prefixed", "~/.bench/worktrees/absent", "target is not registered"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := runVerb(t, verbClean, f.call(tc.target))
			if result.exit == 0 || !strings.Contains(result.stdout, tc.detail) {
				t.Fatalf("%s clean code=%d stdout=%q", tc.name, result.exit, result.stdout)
			}
		})
	}
}

// A resolved operand that this repo declines to remove is a verdict, not a bad operand:
// the plan is the answer and the exit code stays zero.
func TestCleanCommandReportsAResolvedRetainVerdictAsSuccess(t *testing.T) {
	t.Parallel()
	f := newPendingAssignment(t, "verdict")
	mustWrite(t, filepath.Join(f.creation.Path, "residual.log"), []byte("x\n"), 0o644)
	mustWrite(t, filepath.Join(f.creation.Path, ".gitignore"), []byte("residual.log\n"), 0o644)
	gitRun(t, f.creation.Path, "add", ".gitignore")
	gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "ignore residual")
	result := runVerb(t, verbClean, f.call(f.creation.Path))
	if result.exit != 0 {
		t.Fatalf("resolved retain code=%d stdout=%q stderr=%q", result.exit, result.stdout, result.stderr)
	}
}

// [CP1][CP2] `bench worktree path` prints the resolved absolute path, and the help
// rows steer the operator straight from it into `bench worktree clean`. So the two
// have to compose, and the printed form has to work verbatim when quoted — which the
// `~` form never does. Plan and apply must also agree on one canonical target.
func TestCleanCommandAcceptsTheAbsolutePathThatPathPrints(t *testing.T) {
	f := newOwnedAssignment(t, "portable")
	bindEnv(t, "HOME", f.root)
	printed := runVerb(t, verbPath, f.call(f.creation.Assignment.Label))
	if printed.exit != 0 {
		t.Fatalf("path exited %d: %s", printed.exit, printed.stderr)
	}
	portable := strings.TrimSpace(printed.stdout)
	if !filepath.IsAbs(portable) {
		t.Fatalf("path printed %q, want a resolved absolute path", portable)
	}
	planned := runVerb(t, verbClean, f.call(portable))
	if planned.exit != 0 {
		t.Fatalf("clean %q exited %d: %s", portable, planned.exit, planned.stdout)
	}
	if !strings.Contains(planned.stdout, ",remove,") {
		t.Fatalf("clean %q planned no removal: %s", portable, planned.stdout)
	}
	applied := runVerb(t, verbClean, f.call(portable, "--apply", planned.mustFingerprint(t)))
	if applied.exit != 0 {
		t.Fatalf("apply against the portable path exited %d: %s", applied.exit, applied.stdout)
	}
	if !strings.Contains(applied.stdout, ",removed,") {
		t.Fatalf("apply did not remove: %s", applied.stdout)
	}
}

// [CP3] An unsupported home form is a bad operand in its own right. It is not a path
// that merely happens to be unregistered once it has been canonicalized against the
// repo root.
func TestCleanCommandRefusesAnUnsupportedHomeTarget(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "homeform")
	result := runVerb(t, verbClean, f.call("~someone/else"))
	if result.exit == 0 || !strings.Contains(result.stdout, "unsupported home target") {
		t.Fatalf("unsupported home target code=%d stdout=%q", result.exit, result.stdout)
	}
}
