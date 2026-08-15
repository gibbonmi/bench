package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestGateRunGreenReusesRecordedVerdict(t *testing.T) {
	root := outcomeFixture(t)
	var firstOut, firstErr bytes.Buffer
	first := Execute(context.Background(), root, &firstOut, &firstErr)
	if first.GateExit != 0 || first.ActionExit != 0 {
		t.Fatalf("first result = %#v, stderr=%q", first, firstErr.String())
	}
	inspection := Inspect(root)
	if inspection.State != Ready || inspection.Status != "green" || !inspection.ReusableGreen {
		t.Fatalf("inspection = %#v, want reusable ready green", inspection)
	}
	record := outcomeRecord(t, root)

	var secondOut, secondErr bytes.Buffer
	second := Execute(context.Background(), root, &secondOut, &secondErr)
	if second.GateExit != 0 || second.ActionExit != 0 || !second.Inspection.ReusableGreen {
		t.Fatalf("second result = %#v, stderr=%q", second, secondErr.String())
	}
	if got := secondOut.String(); got != "gate: green (fresh verdict reused for this tree)\n" {
		t.Fatalf("second stdout = %q", got)
	}
	if got := outcomeRuns(t, root); got != 1 {
		t.Fatalf("runs = %d, want 1", got)
	}
	if got := outcomeRecord(t, root); !bytes.Equal(got, record) {
		t.Fatalf("recorded verdict changed on reuse\nbefore: %s\nafter: %s", record, got)
	}
}

func TestGateRunRedInvalidatesPriorGreen(t *testing.T) {
	redRoot := outcomeFixture(t)
	outcomeWrite(t, redRoot, ".gate-red", "\n", 0o644)
	redResult := Execute(context.Background(), redRoot, &bytes.Buffer{}, &bytes.Buffer{})
	if redResult.GateExit != 7 || redResult.ActionExit != 7 {
		t.Fatalf("red result = %#v, want gate and action exit 7", redResult)
	}

	root := outcomeFixture(t)
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("green result = %#v", result)
	}
	outcomeWrite(t, root, ".gate-red", "\n", 0o644)
	var redOut, redErr bytes.Buffer
	if got := RunCommand([]string{"--fresh", root}, &redOut, &redErr); got != 7 {
		t.Fatalf("red exit = %d, stderr=%q", got, redErr.String())
	}
	inspection := Inspect(root)
	if inspection.State != Ready || inspection.Status != "red" || inspection.ReusableGreen {
		t.Fatalf("red inspection = %#v", inspection)
	}
	if err := os.Remove(filepath.Join(root, ".gate-red")); err != nil {
		t.Fatal(err)
	}
	var rerunOut, rerunErr bytes.Buffer
	if result := Execute(context.Background(), root, &rerunOut, &rerunErr); result.ActionExit != 0 {
		t.Fatalf("following result = %#v, stderr=%q", result, rerunErr.String())
	}
	if got := outcomeRuns(t, root); got != 3 {
		t.Fatalf("runs after green, red, following green = %d, want 3", got)
	}
}

func TestGateRunCommandFreshRerunsGreen(t *testing.T) {
	root := outcomeFixture(t)
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("green result = %#v", result)
	}
	var stdout, stderr bytes.Buffer
	if got := RunCommand([]string{"--fresh", root}, &stdout, &stderr); got != 0 {
		t.Fatalf("fresh exit = %d, stderr=%q", got, stderr.String())
	}
	if got := outcomeRuns(t, root); got != 2 {
		t.Fatalf("runs after --fresh = %d, want 2", got)
	}
}

func TestGateRunScriptWitnessesPendingRecord(t *testing.T) {
	root := outcomeFixture(t)
	var stdout, stderr bytes.Buffer
	if result := Execute(context.Background(), root, &stdout, &stderr); result.ActionExit != 0 {
		t.Fatalf("result = %#v, stderr=%q", result, stderr.String())
	}
	var copied struct {
		State    State  `json:"state"`
		Tree     string `json:"tree"`
		OwnerPID int    `json:"owner_pid"`
	}
	if err := json.Unmarshal(outcomeRead(t, filepath.Join(root, ".gate-record-during")), &copied); err != nil {
		t.Fatalf("parse script-witnessed record: %v", err)
	}
	inspection := Inspect(root)
	if copied.State != Pending || copied.OwnerPID != os.Getpid() || copied.Tree != inspection.CurrentTree {
		t.Fatalf("script-witnessed record = %#v, inspection = %#v", copied, inspection)
	}
}

func TestGateRunRefusesMovedSubject(t *testing.T) {
	root := outcomeFixture(t)
	outcomeWrite(t, root, ".gate-drift", "\n", 0o644)
	var stdout, stderr bytes.Buffer
	result := Execute(context.Background(), root, &stdout, &stderr)
	if result.GateExit != 0 || result.ActionExit != 1 {
		t.Fatalf("result = %#v, stderr=%q", result, stderr.String())
	}
	if !strings.Contains(stderr.String(), "gate subject changed during execution") {
		t.Fatalf("stderr = %q", stderr.String())
	}
	if result.Inspection.ReusableGreen || Inspect(root).ReusableGreen {
		t.Fatalf("moved subject retained reusable green: result=%#v inspect=%#v", result.Inspection, Inspect(root))
	}
}

func TestGateRunWithoutGateExitsThree(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	outcomeWrite(t, root, "tracked.txt", "tracked\n", 0o644)
	outcomeCommit(t, root, "fixture")
	t.Setenv("BENCH_GATE", "")
	var stdout, stderr bytes.Buffer
	result := Execute(context.Background(), root, &stdout, &stderr)
	if result.GateExit != 3 || result.ActionExit != 3 {
		t.Fatalf("result = %#v, stderr=%q", result, stderr.String())
	}
	if !strings.Contains(stderr.String(), "no gate found") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestGateRunReloadsDeadOwnerPendingRecord(t *testing.T) {
	root := outcomeFixture(t)
	before := Inspect(root)
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")
	pending := fmt.Sprintf(`{"schema":1,"state":"pending","tree":%q,"oracle":%q,"started_at":%q,"owner_pid":999999}`+"\n", before.CurrentTree, strings.Repeat("0", 64), time.Now().UTC().Add(-time.Minute).Truncate(time.Second).Format(time.RFC3339))
	if err := os.WriteFile(filepath.Join(gitdir, "bench-last-gate"), []byte(pending), 0o600); err != nil {
		t.Fatal(err)
	}
	if inspection := Inspect(root); inspection.State != Pending {
		t.Fatalf("dead-owner inspection = %#v, want pending", inspection)
	}
	var stdout, stderr bytes.Buffer
	result := Execute(context.Background(), root, &stdout, &stderr)
	if result.ActionExit != 0 {
		t.Fatalf("result = %#v, stderr=%q", result, stderr.String())
	}
	inspection := Inspect(root)
	if inspection.State != Ready || inspection.Status != "green" || !inspection.ReusableGreen {
		t.Fatalf("after reload inspection = %#v", inspection)
	}
	if got := outcomeRuns(t, root); got != 1 {
		t.Fatalf("runs = %d, want 1", got)
	}
}

func outcomeFixture(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	outcomeWrite(t, root, ".gitignore", ".gate-*\n", 0o644)
	outcomeWrite(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n", 0o644)
	outcomeWrite(t, root, ".bench/gate.sh", `#!/bin/sh
set -eu
count=0
if [ -f .gate-run-count ]; then count=$(cat .gate-run-count); fi
printf '%s' "$((count + 1))" > .gate-run-count
gitdir=$(git rev-parse --absolute-git-dir)
cp "$gitdir/bench-last-gate" .gate-record-during
if [ -e .gate-red ]; then exit 7; fi
if [ -e .gate-drift ]; then printf 'moved\n' >> tracked.txt; fi
`+"\n", 0o755)
	outcomeWrite(t, root, "tracked.txt", "tracked\n", 0o644)
	outcomeCommit(t, root, "fixture")
	return root
}

func outcomeCommit(t *testing.T, root, message string) {
	t.Helper()
	outcomeGit(t, root, "add", ".")
	outcomeGit(t, root, "commit", "-qm", message)
}

func outcomeGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	output, err := benchgit.Output(append([]string{"-C", root}, args...)...)
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(output)
}

func outcomeWrite(t *testing.T, root, path, value string, mode os.FileMode) {
	t.Helper()
	full := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(value), mode); err != nil {
		t.Fatal(err)
	}
}

func outcomeRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func outcomeRecord(t *testing.T, root string) []byte {
	t.Helper()
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")
	return outcomeRead(t, filepath.Join(gitdir, "bench-last-gate"))
}

func outcomeRuns(t *testing.T, root string) int {
	t.Helper()
	runs, err := strconv.Atoi(strings.TrimSpace(string(outcomeRead(t, filepath.Join(root, ".gate-run-count")))))
	if err != nil {
		t.Fatal(err)
	}
	return runs
}
