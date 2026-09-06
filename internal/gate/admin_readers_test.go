package gate

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/gittest"
)

// A reader that cannot resolve the checkout administration directory reports the
// composed-green probe false, instead of propagating the failure as a panic or a wrong
// true. A non-repository root never reaches this reader, so the fixture is a real one.
func TestComposedGreenIsFalseWhenTheReaderFails(t *testing.T) {
	root := outcomeFixture(t)
	logPath := filepath.Join(t.TempDir(), "git.log")
	gittest.StubGit(t, root, "fail-git-dir", logPath)

	if got := composedGreenAtKit(root, root); got {
		t.Fatalf("composedGreenAtKit = %v, want false when the reader fails", got)
	}
}

// A reader that cannot resolve the checkout administration directory refuses the fast
// lane with its own error text, instead of continuing past a gitdir it cannot trust.
func TestRunLaneRefusesWhenTheReaderFails(t *testing.T) {
	root := outcomeFixture(t)
	logPath := filepath.Join(t.TempDir(), "git.log")
	gittest.StubGit(t, root, "fail-git-dir", logPath)

	_, err := runLane(context.Background(), LaneRequest{
		Root:   root,
		Checks: []Phase{{Name: "noop", Argv: []string{"true"}}},
		Stdout: io.Discard,
		Stderr: io.Discard,
	})

	if err == nil || err.Error() != "gate: git directory unavailable" {
		t.Fatalf("runLane err = %v, want %q", err, "gate: git directory unavailable")
	}
}

// A reader that cannot resolve the checkout administration directory reports the
// verdict inspection unavailable, instead of a stale or ready state read off no cache.
func TestInspectSubjectReportsUnavailableWhenTheReaderFails(t *testing.T) {
	root := outcomeFixture(t)
	logPath := filepath.Join(t.TempDir(), "git.log")
	gittest.StubGit(t, root, "fail-git-dir", logPath)

	gi := inspectSubjectAt(root, subject{}, time.Now().UTC())

	if gi.State != Unavailable || gi.Reason != "git directory unavailable" {
		t.Fatalf("inspection = %#v, want state Unavailable and reason %q", gi, "git directory unavailable")
	}
}

// A reader that cannot resolve the checkout administration directory refuses the run
// transaction with an operational result, instead of starting the gate against a gitdir
// it cannot trust.
func TestRunTransactionRefusesWhenTheReaderFails(t *testing.T) {
	root := outcomeFixture(t)
	logPath := filepath.Join(t.TempDir(), "git.log")
	gittest.StubGit(t, root, "fail-git-dir", logPath)
	plan := subject{
		Tree:       strings.Repeat("a", 40),
		Resolution: Resolution{Kind: GateSh},
		Closed:     true,
	}
	var stdout, stderr bytes.Buffer

	result := executeSubjectWithRunBinary(context.Background(), root, root, &stdout, &stderr, nil, forceRun, plannedEvaluation{plan: plan}, nil, "")

	if result.ActionExit != 1 || result.GateExit != 0 {
		t.Fatalf("result = %#v, stderr=%q, want ActionExit 1 and GateExit 0", result, stderr.String())
	}
	if !strings.Contains(stderr.String(), "git directory unavailable") {
		t.Fatalf("stderr = %q, want it to hold %q", stderr.String(), "git directory unavailable")
	}
}
