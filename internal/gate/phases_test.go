package gate

// The engine's whole composition, driven at its own entry point with fixture phases a
// manifest declares. The exec path demands a sealed run binary no unit test holds, so a
// constructed selection is what makes the exit codes and the run log observable here.

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// fixturePhaseRoot builds a graded root whose manifest declares exactly the phases the
// JSON names. A manifest is authoritative, so these phases reach the runner with no Go
// module, no toolchain, and no built kit in the fixture.
//
// The baseline selector is cleared for the run. It names the tree whose manifest selects
// the schedule, and an inherited one would point the loader at a root this test never
// wrote.
func fixturePhaseRoot(t *testing.T, manifest string) string {
	t.Helper()
	t.Setenv(baselinePolicyEnv, "")
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(canary.PhaseManifestPath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// fixtureSelection is the run binary the phase composer merges into every phase's
// environment. No fixture phase reads it, so the pair only has to be non-empty.
func fixtureSelection(root string) *runbinary.Selection {
	return &runbinary.Selection{Path: filepath.Join(root, "bench"), SourceRoot: root}
}

func runFixturePhases(ctx context.Context, t *testing.T, root string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := phasesCommandAtKitWithSelection(ctx, root, root, fixtureSelection(root), &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// BG13: a red fixture phase settles the run at exit 1, and its line reaches the table as
// the phase's one row.
func TestFixturePhaseRedRunExitsOne(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"red","argv":["sh","-c","echo a canned finding; exit 1"]}]}`)

	code, stdout, stderr := runFixturePhases(context.Background(), t, root)

	if code != 1 {
		t.Fatalf("red run exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if rows := rowsForPhase(t, stdout, "red"); len(rows) != 1 || rows[0] != "a canned finding" {
		t.Errorf("red rows = %q, want the phase's one line", rows)
	}
	if !strings.HasSuffix(stdout, "gate: red\n") {
		t.Errorf("red stdout = %q, want it to close on the verdict", stdout)
	}
}

// BG34: a green fixture phase settles the run at exit 0, and stdout is the phase table,
// the skip count, and the verdict.
func TestFixturePhaseGreenRunExitsZero(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"green","argv":["true"]}]}`)

	code, stdout, stderr := runFixturePhases(context.Background(), t, root)

	if code != 0 {
		t.Fatalf("green run exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.HasPrefix(stdout, "phases[1]{phase,verdict,elapsed_ms}:\n") {
		t.Errorf("green stdout = %q, want it to lead with the phase table", stdout)
	}
	if !strings.Contains(stdout, "\n"+skipRowPrefix+": 0 (") {
		t.Errorf("green stdout = %q, want the skip-count line", stdout)
	}
	if !strings.HasSuffix(stdout, "gate: green\n") {
		t.Errorf("green stdout = %q, want it to close on the verdict", stdout)
	}
}

// BG36: a cancelled run exits 130 and publishes no verdict. An interrupt grades nothing,
// so a table here would report phases that never got to answer.
func TestFixturePhaseCancelledRunExitsOneHundredThirty(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"slow","argv":["sh","-c","sleep 5"]}]}`)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	time.AfterFunc(100*time.Millisecond, cancel)

	code, stdout, stderr := runFixturePhases(ctx, t, root)

	if code != 130 {
		t.Fatalf("cancelled run exit = %d, want 130; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if stdout != "" {
		t.Errorf("cancelled stdout = %q, want no verdict", stdout)
	}
}

// BG23: the green table's `elapsed_ms` cell is the same measurement the phase's
// `phase.finish` record carries. One reading of the clock feeds both, so a second timer
// anywhere makes the table and the progress log disagree about one phase.
func TestPhaseElapsedCellEqualsItsFinishRecord(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[{"name":"timed","argv":["sh","-c","sleep 0.05"]}]}`)
	stubGateLogPathIgnored(t)
	var stdout, stderr bytes.Buffer

	ctx, finish := beginGateRunLog(context.Background(), root, &stderr, "dev")
	if strings.Contains(stderr.String(), "progress logging unavailable") {
		t.Fatalf("the run opened no record, so this asserts nothing: %q", stderr.String())
	}
	code := phasesCommandAtKitWithSelection(ctx, root, root, fixtureSelection(root), &stdout, &stderr)
	finish(Result{GateExit: code, ActionExit: code})

	if code != 0 {
		t.Fatalf("timed run exit = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	recorded := phaseFinishElapsed(t, root, "timed")
	if recorded <= 0 {
		t.Fatalf("recorded elapsed_ms = %d, want the phase's own wall time", recorded)
	}
	if cell := phaseTableRow(t, stdout.String(), "timed")[2]; cell != strconv.FormatInt(recorded, 10) {
		t.Errorf("elapsed_ms cell = %q, want the finish record's %d", cell, recorded)
	}
}

// BG24 at the seam the row names, the in-package phases command: the one production line
// that hands the engine's buffer the run's stream file is driven end to end here, and the
// file is then read off disk. A genuinely killed process is out of reach for a unit test,
// so this proves the wiring rather than the kill: the phases the manifest declares reach
// `.logs/gate-<run>.out`, in arrival order, each line under the phase that wrote it. That
// a line lands when it arrives rather than when the phase settles, which is what makes a
// killed run keep anything, is pinned beside it in run_stream_test.go.
func TestFixturePhaseLinesReachTheRunsStreamFile(t *testing.T) {
	root := fixturePhaseRoot(t, `{"phases":[`+
		`{"name":"first","argv":["sh","-c","echo first said one; echo first said two"]},`+
		`{"name":"second","argv":["sh","-c","echo second said one"]}]}`)
	stubGateLogPathIgnored(t)
	var stdout, stderr bytes.Buffer

	ctx, finish := beginGateRunLog(context.Background(), root, &stderr, "dev")
	stream := gateRunStreamFile(ctx)
	if stream == nil {
		t.Fatalf("the run opened no stream file, so this asserts nothing: %q", stderr.String())
	}
	code := phasesCommandAtKitWithSelection(ctx, root, root, fixtureSelection(root), &stdout, &stderr)
	finish(Result{GateExit: code, ActionExit: code})

	if code != 0 {
		t.Fatalf("fixture run exit = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	data, err := os.ReadFile(stream.Name())
	if err != nil {
		t.Fatal(err)
	}
	want := "[first] first said one\n[first] first said two\n[second] second said one\n"
	if string(data) != want {
		t.Errorf("%s =\n%q\nwant\n%q", stream.Name(), data, want)
	}
}

// phaseTableRow answers one phase's row of the green table as its three cells. The green
// table's cells are a filtered name and two engine-owned values, so a plain split is
// exact here.
func phaseTableRow(t *testing.T, stdout, phase string) []string {
	t.Helper()
	for _, line := range strings.Split(stdout, "\n") {
		if !strings.HasPrefix(line, "  ") {
			continue
		}
		cells := strings.Split(strings.TrimPrefix(line, "  "), ",")
		if len(cells) == 3 && cells[0] == phase {
			return cells
		}
	}
	t.Fatalf("phase %s has no row in\n%s", phase, stdout)
	return nil
}

// phaseFinishElapsed answers the elapsed_ms of one phase's finish record in the run log
// this run wrote. The run token is the writer's own, so the record is found by the name
// shape rather than by a token the test would have to predict.
func phaseFinishElapsed(t *testing.T, root, phase string) int64 {
	t.Helper()
	records, err := filepath.Glob(filepath.Join(gateLogDir(root), gateLogName("*", gateLogRecordSuffix)))
	if err != nil || len(records) != 1 {
		t.Fatalf("progress logs = %q, %v; want exactly the one this run wrote", records, err)
	}
	data, err := os.ReadFile(records[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
		var record gateLogRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("progress record %q: %v", line, err)
		}
		if record.Event == "phase.finish" && record.Phase == phase {
			return record.ElapsedMS
		}
	}
	t.Fatalf("phase %s has no finish record in\n%s", phase, data)
	return 0
}

// WF18: the gate's system phase reads the one system-suite producer. The named check
// `bench test --check system` reads the same producer, so neither surface can hold an
// argv list the other does not.
func TestBenchkitPhasesSystemPhaseReadsTheProducer(t *testing.T) {
	kit, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	operands, suiteEnv := SystemSuite(kit)
	var system Phase
	found := false
	for _, phase := range BenchkitPhases(kit, kit) {
		if phase.Name == SystemPhaseName {
			system, found = phase, true
		}
	}
	if !found {
		t.Fatal("BenchkitPhases holds no system phase")
	}
	tail := system.Argv[len(system.Argv)-len(operands):]
	if !reflect.DeepEqual(tail, operands) {
		t.Fatalf("system argv tail = %#v, want the producer's operands %#v", tail, operands)
	}
	if !reflect.DeepEqual(system.Env, suiteEnv) {
		t.Fatalf("system env = %#v, want the producer's environment %#v", system.Env, suiteEnv)
	}
	if want := SystemRootEnv + "=" + kit; !slices.Contains(system.Env, want) {
		t.Fatalf("system env = %#v, want %q", system.Env, want)
	}
}
