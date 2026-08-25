package worktree

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestReauthorizeCommandGrammarKeepsFlagValuesOutOfPath(t *testing.T) {
	flags := []string{"--assignment", "--request", "--base", "--source-tip"}
	for _, flag := range flags {
		t.Run(flag, func(t *testing.T) {
			root, creation, base, tip := reauthorizeFixture(t)
			before := reauthorizeEvidence(t, root, creation.Path)
			values := map[string]string{
				"--assignment": creation.Assignment.ID,
				"--request":    "replacement-request",
				"--base":       base,
				"--source-tip": tip,
			}
			values[flag] = creation.Path
			args := make([]string, 0, len(flags)*2)
			args = append(args, flag, values[flag])
			for _, other := range flags {
				if other == flag {
					continue
				}
				args = append(args, other, values[other])
			}
			var stdout, stderr bytes.Buffer
			if code := ReauthorizeCommand(root, Home(), args, &stdout, &stderr); code != 2 {
				t.Fatalf("path-only-as-%s exit = %d, want 2; stdout=%q stderr=%q", flag, code, stdout.String(), stderr.String())
			}
			if got := reauthorizeEvidence(t, root, creation.Path); got != before {
				t.Fatalf("path-only-as-%s changed retained state\nbefore=%q\nafter=%q", flag, before, got)
			}
		})
	}

	root, creation, base, tip := reauthorizeFixture(t)
	var stdout, stderr bytes.Buffer
	args := []string{"--assignment", creation.Assignment.ID, "--request", "control-token", "--base", base, "--source-tip", tip, "--", creation.Path}
	if code := ReauthorizeCommand(root, Home(), args, &stdout, &stderr); code != 0 {
		t.Fatalf("-- path control exit = %d; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
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
		var stdout, stderr bytes.Buffer
		if code := ReauthorizeCommand("", Home(), args, &stdout, &stderr); code != 2 || stdout.Len() != 0 || stderr.String() != reauthorizeGrammar.Help+"\n" {
			t.Fatalf("ReauthorizeCommand(%q) = (%d, %q, %q), want (2, empty, %q)", args, code, stdout.String(), stderr.String(), reauthorizeGrammar.Help+"\n")
		}
	}
}

func TestReauthorizeCommandEscapesControlBearingBase(t *testing.T) {
	root, creation, _, tip := reauthorizeFixture(t)
	var stdout, stderr bytes.Buffer
	args := []string{"--assignment", creation.Assignment.ID, "--request", "replacement-request", "--base", "not-a-commit\nforged-output", "--source-tip", tip, creation.Path}
	if code := ReauthorizeCommand(root, Home(), args, &stdout, &stderr); code != 1 {
		t.Fatalf("control-bearing base exit = %d, want 1; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if strings.Count(stderr.String(), "\n") != 1 || !strings.Contains(stderr.String(), `\n`) {
		t.Fatalf("control-bearing base forged terminal output: %q", stderr.String())
	}
}

func TestReauthorizeCommandNamesRecordedStartWhenNotAncestor(t *testing.T) {
	root, creation, _, tip := reauthorizeFixture(t)
	commitInWorktree(t, root, "later.txt", "later\n", "later")
	nonAncestorBase := gitOutput(t, root, "rev-parse", "HEAD")
	var stdout, stderr bytes.Buffer
	want := "bench worktree reauthorize: review base is not an ancestor of source tip; wanted=" + creation.Assignment.Start + "\n"
	if code := ReauthorizeCommand(root, Home(), []string{"--assignment", creation.Assignment.ID, "--request", "replacement", "--base", nonAncestorBase, "--source-tip", tip, creation.Path}, &stdout, &stderr); code != 1 || stdout.Len() != 0 || stderr.String() != want {
		t.Fatalf("ancestry refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
}

func TestReauthorizeCommandRefusesRecordedStartOutsideSourceHistory(t *testing.T) {
	root, creation, base, tip := reauthorizeFixture(t)
	commitInWorktree(t, root, "later.txt", "later\n", "later")
	a := creation.Assignment
	a.Start = gitOutput(t, root, "rev-parse", "HEAD")
	if err := intent.PutAssignment(root, a); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "worktree", "unlock", creation.Path)
	gitRun(t, root, "worktree", "lock", "--reason", lockReason(a), creation.Path)
	var stdout, stderr bytes.Buffer
	want := "bench worktree reauthorize: recorded start is not an ancestor of source tip; wanted=" + a.Start + "\n"
	if code := ReauthorizeCommand(root, Home(), []string{"--assignment", a.ID, "--request", "replacement", "--base", base, "--source-tip", tip, creation.Path}, &stdout, &stderr); code != 1 || stdout.Len() != 0 || stderr.String() != want {
		t.Fatalf("recorded-start refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
}

func TestReauthorizeCommandProvesExactIdentityAndChangesOnlyRequest(t *testing.T) {
	root, creation, base, tip := reauthorizeFixture(t)
	before := reauthorizeEvidence(t, root, creation.Path)
	args := []string{"--assignment", creation.Assignment.ID, "--request", "replacement-request", "--base", base, "--source-tip", tip, creation.Path}
	var stdout, stderr bytes.Buffer
	unknown := append([]string(nil), args...)
	unknown[1] = strings.Repeat("f", 32)
	if code := ReauthorizeCommand(root, Home(), unknown, &stdout, &stderr); code != 1 {
		t.Fatalf("unknown assignment exit = %d, want 1; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if got := reauthorizeEvidence(t, root, creation.Path); got != before {
		t.Fatalf("unknown assignment changed retained state\nbefore=%q\nafter=%q", before, got)
	}
	stdout.Reset()
	stderr.Reset()
	if code := ReauthorizeCommand(root, Home(), args, &stdout, &stderr); code != 0 {
		t.Fatalf("reauthorize exit = %d; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if want := "reauthorized{assignment=" + creation.Assignment.ID + ",recorded_start=" + creation.Assignment.Start + ",approved_base=" + base + ",source_tip=" + tip + ",state=active}\n"; stdout.String() != want {
		t.Fatalf("reauthorize stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("reauthorize stderr = %q, want empty", stderr.String())
	}
	afterSuccess := reauthorizeEvidence(t, root, creation.Path)
	expectedBefore := afterSuccess
	expectedBefore.Request = before.Request
	expectedBefore.Lock = before.Lock
	if expectedBefore != before {
		t.Fatalf("reauthorize changed more than request\nbefore=%q\nafter=%q", before, afterSuccess)
	}
	assignment, err := assignmentByID(root, creation.Assignment.ID)
	if err != nil {
		t.Fatal(err)
	}
	if assignment.Request != intent.RequestDigest("replacement-request") {
		t.Fatalf("assignment request = %q, want replacement digest", assignment.Request)
	}
	expectedAssignment := assignment
	expectedAssignment.Request = creation.Assignment.Request
	if !reflect.DeepEqual(expectedAssignment, creation.Assignment) {
		t.Fatalf("ledger record changed beyond request: before=%#v after=%#v", creation.Assignment, assignment)
	}

	gitRun(t, creation.Path, "checkout", "--detach", tip)
	beforeDetached := reauthorizeEvidence(t, root, creation.Path)
	stdout.Reset()
	stderr.Reset()
	if code := ReauthorizeCommand(root, Home(), args, &stdout, &stderr); code != 1 {
		t.Fatalf("detached assignment exit = %d, want 1; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if got := reauthorizeEvidence(t, root, creation.Path); got != beforeDetached {
		t.Fatalf("detached refusal changed retained state\nbefore=%q\nafter=%q", beforeDetached, got)
	}
}

func TestReauthorizeCommandRollsBackLockRefreshAndCASLoss(t *testing.T) {
	cases := []struct {
		name    string
		install func(*testing.T, Creation)
	}{
		{
			name: "unlock failure",
			install: func(t *testing.T, _ Creation) {
				old := reauthorizeUnlock
				reauthorizeUnlock = func(string, string) error { return errors.New("injected unlock failure") }
				t.Cleanup(func() { reauthorizeUnlock = old })
			},
		},
		{
			name: "relock failure",
			install: func(t *testing.T, creation Creation) {
				old := reauthorizeLock
				next := creation.Assignment
				next.Request = intent.RequestDigest("replacement-request")
				reauthorizeLock = func(root, path, reason string) error {
					if reason == lockReason(next) {
						return errors.New("injected relock failure")
					}
					return old(root, path, reason)
				}
				t.Cleanup(func() { reauthorizeLock = old })
			},
		},
		{
			name: "expected-old loss",
			install: func(t *testing.T, _ Creation) {
				old := reauthorizeBeforeCAS
				reauthorizeBeforeCAS = func(a *intent.Assignment) { a.Request = intent.RequestDigest("concurrent-winner") }
				t.Cleanup(func() { reauthorizeBeforeCAS = old })
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			root, creation, base, tip := reauthorizeFixture(t)
			before := reauthorizeEvidence(t, root, creation.Path)
			testCase.install(t, creation)
			var stdout, stderr bytes.Buffer
			args := []string{"--assignment", creation.Assignment.ID, "--request", "replacement-request", "--base", base, "--source-tip", tip, creation.Path}
			if code := ReauthorizeCommand(root, Home(), args, &stdout, &stderr); code != 1 {
				t.Fatalf("%s exit = %d, want 1; stdout=%q stderr=%q", testCase.name, code, stdout.String(), stderr.String())
			}
			if got := reauthorizeEvidence(t, root, creation.Path); got != before {
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
	indexPath := gitOutput(t, path, "rev-parse", "--git-path", "index")
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

func reauthorizeFixture(t *testing.T) (string, Creation, string, string) {
	t.Helper()
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
	base := gitOutput(t, root, "rev-parse", "HEAD")
	if err := os.WriteFile(filepath.Join(root, "reviewed.txt"), []byte("reviewed source\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "reviewed.txt")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "reviewed source")
	creation := mustCreate(t, root, "lost-request", "reauthorize fixture")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	if creation.Assignment.Start != tip {
		t.Fatalf("fixture start = %s, tip = %s, want equal", creation.Assignment.Start, tip)
	}
	return root, creation, base, tip
}
