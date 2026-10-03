package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestReauthorizeCommandGrammarKeepsFlagValuesOutOfPath(t *testing.T) {
	t.Parallel()
	flags := []string{"--assignment", "--request", "--base", "--source-tip"}
	for _, flag := range flags {
		t.Run(flag, func(t *testing.T) {
			f := reauthorizeFixture(t)
			before := reauthorizeEvidence(t, f.root, f.creation.Path)
			values := map[string]string{
				"--assignment": f.creation.Assignment.ID,
				"--request":    "replacement-request",
				"--base":       f.base,
				"--source-tip": f.tip,
			}
			values[flag] = f.creation.Path
			args := make([]string, 0, len(flags)*2)
			args = append(args, flag, values[flag])
			for _, other := range flags {
				if other == flag {
					continue
				}
				args = append(args, other, values[other])
			}
			if r := runVerb(t, verbReauthorize, f.call(args...)); r.exit != 2 {
				t.Fatalf("path-only-as-%s exit = %d, want 2; stdout=%q stderr=%q", flag, r.exit, r.stdout, r.stderr)
			}
			if got := reauthorizeEvidence(t, f.root, f.creation.Path); got != before {
				t.Fatalf("path-only-as-%s changed retained state\nbefore=%q\nafter=%q", flag, before, got)
			}
		})
	}

	f := reauthorizeFixture(t)
	args := []string{"--assignment", f.creation.Assignment.ID, "--request", "control-token", "--base", f.base, "--source-tip", f.tip, "--", f.creation.Path}
	if r := runVerb(t, verbReauthorize, f.call(args...)); r.exit != 0 {
		t.Fatalf("-- path control exit = %d; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
}

func TestReauthorizeCommandRequiredFlagsKeepDeclaredHelp(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"--request", "r", "--base", "b", "--source-tip", "s", "path"},
		{"--assignment", "a", "--base", "b", "--source-tip", "s", "path"},
		{"--assignment", "a", "--request", "r", "--source-tip", "s", "path"},
		{"--assignment", "a", "--request", "r", "--base", "b", "path"},
	} {
		if r := runVerb(t, verbReauthorize, verbCall{home: Home(), args: args}); r.exit != 2 || r.stdout != "" || r.stderr != reauthorizeGrammar.Help+"\n" {
			t.Fatalf("reauthorize %q = (%d, %q, %q), want (2, empty, %q)", args, r.exit, r.stdout, r.stderr, reauthorizeGrammar.Help+"\n")
		}
	}
}

func TestReauthorizeCommandEscapesControlBearingBase(t *testing.T) {
	t.Parallel()
	f := reauthorizeFixture(t)
	args := []string{"--assignment", f.creation.Assignment.ID, "--request", "replacement-request", "--base", "not-a-commit\nforged-output", "--source-tip", f.tip, f.creation.Path}
	r := runVerb(t, verbReauthorize, f.call(args...))
	if r.exit != 1 {
		t.Fatalf("control-bearing base exit = %d, want 1; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	if strings.Count(r.stderr, "\n") != 1 || !strings.Contains(r.stderr, `\n`) {
		t.Fatalf("control-bearing base forged terminal output: %q", r.stderr)
	}
}

func TestReauthorizeCommandNamesRecordedStartWhenNotAncestor(t *testing.T) {
	t.Parallel()
	f := reauthorizeFixture(t)
	commitInWorktree(t, f.root, "later.txt", "later\n", "later")
	nonAncestorBase := gitOutput(t, f.root, "rev-parse", "HEAD")
	want := "bench worktree reauthorize: review base is not an ancestor of source tip; wanted=" + f.creation.Assignment.Start + "\n"
	if r := runVerb(t, verbReauthorize, f.call("--assignment", f.creation.Assignment.ID, "--request", "replacement", "--base", nonAncestorBase, "--source-tip", f.tip, f.creation.Path)); r.exit != 1 || r.stdout != "" || r.stderr != want {
		t.Fatalf("ancestry refusal = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
}

func TestReauthorizeCommandRefusesRecordedStartOutsideSourceHistory(t *testing.T) {
	t.Parallel()
	f := reauthorizeFixture(t)
	commitInWorktree(t, f.root, "later.txt", "later\n", "later")
	a := f.creation.Assignment
	a.Start = gitOutput(t, f.root, "rev-parse", "HEAD")
	if err := intent.PutAssignment(f.root, a); err != nil {
		t.Fatal(err)
	}
	gitRun(t, f.root, "worktree", "unlock", f.creation.Path)
	gitRun(t, f.root, "worktree", "lock", "--reason", lockReason(a), f.creation.Path)
	want := "bench worktree reauthorize: recorded start is not an ancestor of source tip; wanted=" + a.Start + "\n"
	if r := runVerb(t, verbReauthorize, f.call("--assignment", a.ID, "--request", "replacement", "--base", f.base, "--source-tip", f.tip, f.creation.Path)); r.exit != 1 || r.stdout != "" || r.stderr != want {
		t.Fatalf("recorded-start refusal = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
}

func TestReauthorizeCommandProvesExactIdentityAndChangesOnlyRequest(t *testing.T) {
	t.Parallel()
	f := reauthorizeFixture(t)
	before := reauthorizeEvidence(t, f.root, f.creation.Path)
	args := []string{"--assignment", f.creation.Assignment.ID, "--request", "replacement-request", "--base", f.base, "--source-tip", f.tip, f.creation.Path}
	unknown := append([]string(nil), args...)
	unknown[1] = strings.Repeat("f", 32)
	if r := runVerb(t, verbReauthorize, f.call(unknown...)); r.exit != 1 {
		t.Fatalf("unknown assignment exit = %d, want 1; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	if got := reauthorizeEvidence(t, f.root, f.creation.Path); got != before {
		t.Fatalf("unknown assignment changed retained state\nbefore=%q\nafter=%q", before, got)
	}
	r := runVerb(t, verbReauthorize, f.call(args...))
	if r.exit != 0 {
		t.Fatalf("reauthorize exit = %d; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	if want := "reauthorized{assignment=" + f.creation.Assignment.ID + ",recorded_start=" + f.creation.Assignment.Start + ",approved_base=" + f.base + ",source_tip=" + f.tip + ",state=active}\n"; r.stdout != want {
		t.Fatalf("reauthorize stdout = %q, want %q", r.stdout, want)
	}
	if r.stderr != "" {
		t.Fatalf("reauthorize stderr = %q, want empty", r.stderr)
	}
	afterSuccess := reauthorizeEvidence(t, f.root, f.creation.Path)
	expectedBefore := afterSuccess
	expectedBefore.Request = before.Request
	expectedBefore.Lock = before.Lock
	if expectedBefore != before {
		t.Fatalf("reauthorize changed more than request\nbefore=%q\nafter=%q", before, afterSuccess)
	}
	assignment, err := assignmentByID(f.root, f.creation.Assignment.ID)
	if err != nil {
		t.Fatal(err)
	}
	if assignment.Request != intent.RequestDigest("replacement-request") {
		t.Fatalf("assignment request = %q, want replacement digest", assignment.Request)
	}
	expectedAssignment := assignment
	expectedAssignment.Request = f.creation.Assignment.Request
	if !reflect.DeepEqual(expectedAssignment, f.creation.Assignment) {
		t.Fatalf("ledger record changed beyond request: before=%#v after=%#v", f.creation.Assignment, assignment)
	}

	gitRun(t, f.creation.Path, "checkout", "--detach", f.tip)
	beforeDetached := reauthorizeEvidence(t, f.root, f.creation.Path)
	if r := runVerb(t, verbReauthorize, f.call(args...)); r.exit != 1 {
		t.Fatalf("detached assignment exit = %d, want 1; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	if got := reauthorizeEvidence(t, f.root, f.creation.Path); got != beforeDetached {
		t.Fatalf("detached refusal changed retained state\nbefore=%q\nafter=%q", beforeDetached, got)
	}
}

func TestReauthorizeCommandRollsBackLockRefreshAndCASLoss(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		install func(*testing.T, reauthorizeSet, joins) joins
		// stderr is the exact refusal when a case pins one. A denied administration
		// directory also fails the relock, so the retained state alone cannot show that
		// the refusal came from the unlock.
		stderr string
	}{
		{
			// The unlock removes the lock file from the administration directory, so a
			// directory that denies writes fails the real unlock. Root bypasses that mode.
			name:   "unlock failure",
			stderr: "bench worktree reauthorize: refresh ownership lock\n",
			install: func(t *testing.T, f reauthorizeSet, j joins) joins {
				if os.Geteuid() == 0 {
					capability.Capability(t, capability.Privilege, "root bypasses directory permissions; cannot deny writes to fail the unlock")
				}
				admin, err := git.AdminDir(f.creation.Path)
				mustNoError(t, err)
				t.Cleanup(func() { _ = os.Chmod(admin, 0o700) })
				mustNoError(t, os.Chmod(admin, 0o500))
				return j
			},
		},
		{
			name: "relock failure",
			install: func(_ *testing.T, f reauthorizeSet, j joins) joins {
				old := j.reauthorizeLock
				next := f.creation.Assignment
				next.Request = intent.RequestDigest("replacement-request")
				j.reauthorizeLock = func(root, path, reason string) error {
					if reason == lockReason(next) {
						return errors.New("injected relock failure")
					}
					return old(root, path, reason)
				}
				return j
			},
		},
		{
			name: "expected-old loss",
			install: func(_ *testing.T, _ reauthorizeSet, j joins) joins {
				j.reauthorizeBeforeCAS = func(a *intent.Assignment) { a.Request = intent.RequestDigest("concurrent-winner") }
				return j
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			f := reauthorizeFixture(t)
			before := reauthorizeEvidence(t, f.root, f.creation.Path)
			j := testCase.install(t, f, defaultJoins())
			args := []string{"--assignment", f.creation.Assignment.ID, "--request", "replacement-request", "--base", f.base, "--source-tip", f.tip, f.creation.Path}
			if r := runVerb(t, verbReauthorize, f.callWith(j, args...)); r.exit != 1 || (testCase.stderr != "" && r.stderr != testCase.stderr) {
				t.Fatalf("%s exit = %d, want 1; stdout=%q stderr=%q, want stderr %q", testCase.name, r.exit, r.stdout, r.stderr, testCase.stderr)
			}
			if got := reauthorizeEvidence(t, f.root, f.creation.Path); got != before {
				t.Fatalf("%s changed authority or worktree state\nbefore=%q\nafter=%q", testCase.name, before, got)
			}
		})
	}
}

type reauthorizeState struct {
	Request string
	Tree    string
	Index   string
	Status  string
	Refs    string
	Lock    string
}

func reauthorizeEvidence(t *testing.T, root, path string) reauthorizeState {
	t.Helper()
	assignments, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(assignments) != 1 {
		t.Fatalf("assignment count = %d, want 1", len(assignments))
	}
	indexPath := mustAdminPath(t, path, "index")
	index, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	tracked, err := os.ReadFile(filepath.Join(path, "reviewed.txt"))
	if err != nil {
		t.Fatal(err)
	}
	return reauthorizeState{
		Request: assignments[0].Request,
		Tree:    string(tracked),
		Index:   string(index),
		Status:  gitOutput(t, path, "status", "--porcelain=v1"),
		Refs:    gitOutput(t, root, "show-ref", "--head"),
		Lock:    gitOutput(t, root, "worktree", "list", "--porcelain"),
	}
}

// reauthorizeFixture returns the repository's one registration with the reviewed base and
// tip. The home is explicit, so the fixture binds no process environment.
func reauthorizeFixture(t *testing.T) reauthorizeSet {
	t.Helper()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	if err := os.WriteFile(filepath.Join(root, "reviewed.txt"), []byte("reviewed source\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "reviewed.txt")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "reviewed source")
	creation := mustCreate(t, root, home, "lost-request", "reauthorize fixture")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	if creation.Assignment.Start != tip {
		t.Fatalf("fixture start = %s, tip = %s, want equal", creation.Assignment.Start, tip)
	}
	return reauthorizeSet{ownedAssignment: ownedAssignment{repoHome: repoHome{root: root, home: home}, creation: creation}, base: base, tip: tip}
}
