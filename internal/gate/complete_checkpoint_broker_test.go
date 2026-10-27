package gate

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// BrokerExec is one call of the broker exec seam: the broker, its arguments after the program
// name, and the working directory at the call.
type BrokerExec struct {
	Broker, Dir string
	Args        []string
}

// StubBrokerExec replaces the broker exec seam until t ends with one that records each call
// and returns as a successful replacement, with no process replacement. It answers the
// recorded calls.
func StubBrokerExec(t *testing.T) *[]BrokerExec {
	t.Helper()
	var calls []BrokerExec
	previous := execBroker
	t.Cleanup(func() { execBroker = previous })
	execBroker = func(broker string, args, _ []string) error {
		dir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		calls = append(calls, BrokerExec{Broker: broker, Dir: dir, Args: args})
		return nil
	}
	return &calls
}

// setProcessArgs makes args the process argument vector until t ends.
func setProcessArgs(t *testing.T, args ...string) {
	t.Helper()
	previous := os.Args
	t.Cleanup(func() { os.Args = previous })
	os.Args = args
}

// InstallRunning is commitmenttest.BindRunning under the name that the external gate test
// package reads.
var InstallRunning = commitmenttest.BindRunning

// DeclareBuildInputs prepares a checkpoint fixture whose repository declares Bench build
// inputs. The fixture ignores the wrapper directory, so a broker manifest there leaves the
// checkout clean.
func DeclareBuildInputs(f *recordtest.Fixture) {
	ignoreWrapperDirectory(f)
	f.Write(benchfreshness.BuildInputsManifest, "cmd\n")
}

// ignoreWrapperDirectory ignores the wrapper directory of the fixture's checkout, which holds
// the broker manifest.
func ignoreWrapperDirectory(f *recordtest.Fixture) {
	ignored, err := os.ReadFile(filepath.Join(f.Root, ".gitignore"))
	if err != nil {
		f.T.Fatal(err)
	}
	f.Write(".gitignore", string(ignored)+"/bin/\n")
}

// brokerSource is a completed source whose checkout ignores the wrapper directory. With
// declare, its repository declares Bench build inputs.
func brokerSource(t *testing.T, declare bool) *recordtest.Fixture {
	t.Helper()
	prepare := ignoreWrapperDirectory
	if declare {
		prepare = DeclareBuildInputs
	}
	f := attachedCheckpointFixture(t, recordtest.Attach, prepare)
	retainCompletion(f)
	return f
}

// bindBroker writes a broker that is not the running executable and binds it in the broker
// manifest of root. It answers the bound broker path.
func bindBroker(t *testing.T, root string) string {
	t.Helper()
	broker := filepath.Join(t.TempDir(), "broker")
	if err := os.WriteFile(broker, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	commitmenttest.BindManifest(t, root, broker)
	resolved, err := filepath.EvalSymlinks(broker)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

// requireOwnProcess requires an exit 0 after one oracle run and no broker exec.
func requireOwnProcess(t *testing.T, f *recordtest.Fixture, calls *[]BrokerExec, code int, out string) {
	t.Helper()
	if code != 0 || len(*calls) != 0 || outcomeRuns(t, f.Root) != 1 {
		t.Fatalf("checkpoint = (%d, %q) with broker execs %v; want exit 0 after one oracle run in its own process", code, out, *calls)
	}
}

// requireBrokerUnavailable requires the checkpoint-broker-unavailable refusal with no
// successful broker exec and no oracle run.
func requireBrokerUnavailable(t *testing.T, f *recordtest.Fixture, calls *[]BrokerExec, code int, out string) {
	t.Helper()
	requireBrokerRoute(t, code, out)
	if len(*calls) != 0 {
		t.Fatalf("checkpoint refused after broker execs %v, want none", *calls)
	}
	assertNoOracleRun(t, f)
}

// requireBrokerRoute requires a refusal whose route repeats the landing's broker repair.
func requireBrokerRoute(t *testing.T, code int, out string) {
	t.Helper()
	if next, printed := routetest.Next(out); code != 1 || !printed || !strings.Contains(next, "republish the promotion broker") {
		t.Fatalf("checkpoint = (%d, %q), want the checkpoint-broker-unavailable refusal", code, out)
	}
}

// GR59: a complete checkpoint whose running executable differs from the broker execs the
// broker with the same arguments in the same working directory, and runs no gate itself.
func TestCompleteCheckpointExecsTheBroker(t *testing.T) {
	f := brokerSource(t, true)
	broker := bindBroker(t, f.Root)
	calls := StubBrokerExec(t)
	setProcessArgs(t, "bench", "gate", "--checkpoint", "specs/example/spec.md", "--complete")
	dir := t.TempDir()
	t.Chdir(dir)
	code, out := completeCheckpoint(t, f)
	want := BrokerExec{Broker: broker, Dir: dir, Args: []string{"gate", "--checkpoint", "specs/example/spec.md", "--complete"}}
	if code != 0 || len(*calls) != 1 || (*calls)[0].Broker != want.Broker || (*calls)[0].Dir != want.Dir || !slices.Equal((*calls)[0].Args, want.Args) {
		t.Fatalf("checkpoint = (%d, %q) with broker execs %v; want one exec %v", code, out, *calls, want)
	}
	assertNoOracleRun(t, f)
}

// GR60: a complete checkpoint whose running executable is the broker runs in its own process.
func TestCompleteCheckpointRunsUnderTheBroker(t *testing.T) {
	f := brokerSource(t, true)
	commitmenttest.BindRunning(t, f.Root)
	calls := StubBrokerExec(t)
	code, out := completeCheckpoint(t, f)
	requireOwnProcess(t, f, calls, code, out)
}

// GR61: a manifest that does not read refuses the complete checkpoint.
func TestCompleteCheckpointRefusesAnUnreadableManifest(t *testing.T) {
	f := brokerSource(t, true)
	calls := StubBrokerExec(t)
	code, out := completeCheckpoint(t, f)
	requireBrokerUnavailable(t, f, calls, code, out)
}

// GR62: a broker whose bytes changed after its manifest bound them refuses the complete
// checkpoint.
func TestCompleteCheckpointRefusesABrokerDigestMismatch(t *testing.T) {
	f := brokerSource(t, true)
	broker := bindBroker(t, f.Root)
	if err := os.WriteFile(broker, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	calls := StubBrokerExec(t)
	code, out := completeCheckpoint(t, f)
	requireBrokerUnavailable(t, f, calls, code, out)
}

// GR63: a repository that declares no build inputs runs the complete checkpoint under its
// running executable, and a manifest that binds another broker does not move it.
func TestCompleteCheckpointKeepsTheExecutableWithoutBuildInputs(t *testing.T) {
	f := brokerSource(t, false)
	bindBroker(t, f.Root)
	calls := StubBrokerExec(t)
	code, out := completeCheckpoint(t, f)
	requireOwnProcess(t, f, calls, code, out)
}

// GR63: a checkout that declares no build inputs keeps its running executable even when its
// primary checkout does not resolve.
func TestCompleteCheckpointWithoutBuildInputsNeedsNoPrimary(t *testing.T) {
	calls := StubBrokerExec(t)
	var out bytes.Buffer
	checkpoint := Checkpoint{Spec: recordtest.Spec, Complete: true}
	if _, done := runUnderBroker(context.Background(), t.TempDir(), checkpoint, &out, reuseFreshGreen); done || len(*calls) != 0 {
		t.Fatalf("broker step = (done %t, %q) with broker execs %v; want the run kept in its own process", done, out.String(), *calls)
	}
}

// A running executable whose digest does not read refuses the complete checkpoint, so the
// broker that it would exec cannot exec itself again.
func TestCompleteCheckpointRefusesAnUndigestedExecutable(t *testing.T) {
	f := brokerSource(t, true)
	bindBroker(t, f.Root)
	setRunningExecutable(t, "")
	calls := StubBrokerExec(t)
	code, out := completeCheckpoint(t, f)
	requireBrokerUnavailable(t, f, calls, code, out)
}

// GR64: a chunk checkpoint runs under its running executable when the broker differs.
func TestChunkCheckpointKeepsTheRunningExecutable(t *testing.T) {
	f := brokerSource(t, true)
	bindBroker(t, f.Root)
	calls := StubBrokerExec(t)
	code, out := runCheckpoint(t, f)
	requireOwnProcess(t, f, calls, code, out)
}

// Each other cause that leaves the broker unusable refuses the complete checkpoint: a checkout
// that declares build inputs and whose primary checkout does not resolve, whose manifest the
// run must never look up beside the working directory, and an exec that fails.
func TestCompleteCheckpointRefusesAnUnusableBroker(t *testing.T) {
	t.Run("unresolved primary", func(t *testing.T) {
		dir := t.TempDir()
		commitmenttest.BindRunning(t, dir)
		commitmenttest.Write(t, dir, benchfreshness.BuildInputsManifest, "cmd\n")
		t.Chdir(dir)
		var out bytes.Buffer
		code := RunCommand([]string{dir, "--checkpoint", recordtest.Spec, "--complete"}, &out, &out)
		requireBrokerRoute(t, code, out.String())
	})
	t.Run("failed exec", func(t *testing.T) {
		f := brokerSource(t, true)
		bindBroker(t, f.Root)
		previous := execBroker
		t.Cleanup(func() { execBroker = previous })
		execBroker = func(string, []string, []string) error { return errors.New("exec format error") }
		code, out := completeCheckpoint(t, f)
		requireBrokerRoute(t, code, out)
		assertNoOracleRun(t, f)
	})
}
