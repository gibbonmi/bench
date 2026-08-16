package worktree

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"
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

// [CP1][CP2] `bench worktree path` prints a portable `~`-prefixed path and the help rows
// steer the operator straight from it into `bench worktree clean`, so the two have to
// compose. Plan and apply must also agree on one canonical target: a fingerprint taken
// against the unexpanded path and applied against the expanded one would remove a
// checkout the operator never saw named.
func TestCleanCommandAcceptsThePortablePathThatPathPrints(t *testing.T) {
	root, creation := newOwnedAssignment(t, "portable")
	// The fixture worktree lives under BENCH_HOME inside root, so pointing HOME at root
	// is what makes `path` print the portable form this test is about.
	t.Setenv("HOME", root)
	t.Chdir(root)
	var printed, stderr bytes.Buffer
	if code := PathCommand(root, []string{creation.Assignment.Label}, &printed, &stderr); code != 0 {
		t.Fatalf("path exited %d: %s", code, stderr.String())
	}
	portable := strings.TrimSpace(printed.String())
	if !strings.HasPrefix(portable, "~/") {
		t.Fatalf("path printed %q, want the portable home form", portable)
	}
	var planned bytes.Buffer
	if code := CleanCommand([]string{portable}, &planned, &stderr); code != 0 {
		t.Fatalf("clean %q exited %d: %s", portable, code, planned.String())
	}
	if !strings.Contains(planned.String(), ",remove,") {
		t.Fatalf("clean %q planned no removal: %s", portable, planned.String())
	}
	fingerprint := regexp.MustCompile(`[0-9a-f]{64}`).FindString(planned.String())
	if fingerprint == "" {
		t.Fatalf("plan carried no fingerprint: %s", planned.String())
	}
	var applied bytes.Buffer
	if code := CleanCommand([]string{portable, "--apply", fingerprint}, &applied, &stderr); code != 0 {
		t.Fatalf("apply against the portable path exited %d: %s", code, applied.String())
	}
	if !strings.Contains(applied.String(), ",removed,") {
		t.Fatalf("apply did not remove: %s", applied.String())
	}
}

// [CP3] An unsupported home form is a bad operand in its own right, not a path that
// merely happens to be unregistered once it has been canonicalized against the repo root.
func TestCleanCommandRefusesAnUnsupportedHomeTarget(t *testing.T) {
	root, _ := newOwnedAssignment(t, "homeform")
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	code := CleanCommand([]string{"~someone/else"}, &stdout, &stderr)
	if code == 0 || !bytes.Contains(stdout.Bytes(), []byte("unsupported home target")) {
		t.Fatalf("unsupported home target code=%d stdout=%q", code, stdout.String())
	}
}
