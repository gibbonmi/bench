package gate

// The delta checkpoint. The delta gate records the schedule of each run, so a test tells a
// run of the delta set from a whole run by that record and by the evidence store.

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// TestDeltaCheckpointRunsTheDeltaSet covers RC58. A red planted outside the delta set reds
// the whole schedule, and the delta set stays green over it with exactly its own phases. A
// root that is not the kit runs no root conformance, and a malformed schedule refuses.
func TestDeltaCheckpointRunsTheDeltaSet(t *testing.T) {
	t.Setenv(baselinePolicyEnv, "")
	root := gittest.Repo(t)
	for path, body := range map[string]string{
		"go.mod":                  "module fixture\n\ngo 1.21\n",
		"changed/changed_test.go": goTest("changed", "TestChanged", ""),
		"other/other_test.go":     goTest("other", "TestPlanted", `t.Fatal("planted red")`),
		conformancePackagePath + "/entry_test.go": goTest("conformance", registry.RootConformanceTest, ""),
	} {
		outcomeWrite(t, root, path, body, 0o644)
	}
	outcomeWrite(t, root, "bench", "#!/bin/sh\nexit 0\n", 0o755)
	planted := func(row string) bool { return strings.Contains(row, "TestPlanted") }
	if code, out, diagnostic := runSchedule(t, root, root, ""); code == 0 || !slices.ContainsFunc(rowsForPhase(t, out, canary.PhaseTest), planted) {
		t.Fatalf("whole schedule = %d %s %s; want the planted red", code, out, diagnostic)
	}
	for _, tc := range []struct {
		kit  string
		want []string
	}{
		{root, []string{canary.PhaseGofmt, canary.PhaseVet, canary.PhaseTest, conformancePhaseName}},
		{t.TempDir(), []string{canary.PhaseGofmt, canary.PhaseVet, canary.PhaseTest}},
	} {
		code, out, diagnostic := runSchedule(t, root, tc.kit, `["./changed"]`)
		if names := tablePhases(out); code != 0 || !slices.Equal(names, tc.want) || !strings.HasSuffix(out, "\ngate: green (delta checkpoint)\n") {
			t.Fatalf("delta set under kit %s = %d %s %s; want the phases %q and the delta verdict", tc.kit, code, out, diagnostic, tc.want)
		}
	}
	if code, out, _ := runSchedule(t, root, root, "./changed"); code == 0 || out != "" {
		t.Fatalf("malformed delta schedule = %d %s; want a refusal before any phase", code, out)
	}
}

// TestDeltaCheckpointPackagesComeFromTheTestGapList covers RC59: the packages are the
// directories of the changed test files, once each, a file name with a space included.
func TestDeltaCheckpointPackagesComeFromTheTestGapList(t *testing.T) {
	f := deltaFixture(t, graded(t, func(f *recordtest.Fixture) {
		modifyTest(f)
		f.Write("spaced/a b_test.go", "package spaced\n")
		f.Write("spaced/c_test.go", "package spaced\n")
	}))
	requireDelta(t, f, 0, "./.", "./spaced")
}

// TestDeltaCheckpointNeedsAWholeBaseGreen covers RC61.
func TestDeltaCheckpointNeedsAWholeBaseGreen(t *testing.T) {
	requireWhole(t, deltaFixture(t, modifyTest))
}

// TestDeltaCheckpointReadsTheBaseWithoutExpiry covers RC62. Witness: with no changed test
// file, an expired green of the subject itself takes no delta, so the whole gate runs.
func TestDeltaCheckpointReadsTheBaseWithoutExpiry(t *testing.T) {
	f := deltaFixture(t, graded(t, modifyTest))
	expireEvidence(t, f.Root)
	requireDelta(t, f, 0, "./.")
	f = frozenDeltaFixture(t)
	requirePrepared(t, f)
	expireEvidence(t, f.Root)
	requireWhole(t, f)
}

// TestDeltaCheckpointRunsWholeForANonTestDelta covers RC63. A comment-only source commit
// keeps the evidence pass, and the delta from the reviewed tip holds that source file.
func TestDeltaCheckpointRunsWholeForANonTestDelta(t *testing.T) {
	f := deltaFixture(t, graded(t, modifyTest))
	f.Write(repairSource, "package fixture\n\n// value is the repaired constant.\nconst value = 1\n")
	f.Commit("comment the source")
	requireWhole(t, f)
}

// TestDeltaCheckpointAcceptsAnAssignmentOnlyPlanCommit covers RC64.
func TestDeltaCheckpointAcceptsAnAssignmentOnlyPlanCommit(t *testing.T) {
	f := frozenDeltaFixture(t)
	repairChunk(f, graded(t, modifyTest))
	dispatchRepair(f)
	retainWrittenConfirmation(t, f)
	requireDelta(t, f, 0, "./.")
}

// TestDeltaCheckpointRunsWholeForTheSystemSuite covers RC65.
func TestDeltaCheckpointRunsWholeForTheSystemSuite(t *testing.T) {
	operands, _ := SystemSuite(t.TempDir())
	suite := operands[slices.IndexFunc(operands, func(operand string) bool { return strings.HasPrefix(operand, "./") })]
	requireWhole(t, deltaFixture(t, graded(t, addTestIn(suite))))
}

// TestDeltaCheckpointRunsWholeForARacePackage covers RC66.
func TestDeltaCheckpointRunsWholeForARacePackage(t *testing.T) {
	requireWhole(t, deltaFixture(t, graded(t, addTestIn(raceTests[0].PackagePath))))
}

// TestDeltaGreenAuthorizesNoLanding covers RC67: the landing's inspection of the subject
// finds no green after a green delta.
func TestDeltaGreenAuthorizesNoLanding(t *testing.T) {
	f := deltaFixture(t, graded(t, modifyTest))
	requireDelta(t, f, 0, "./.")
	if got := InspectTree(f.Root, strippedSubject(t, f)); got.ReusableGreen {
		t.Fatalf("landing inspection after a delta green = %+v; want no reusable green", got)
	}
}

// TestDeltaCheckpointFreshRunsWhole covers RC71.
func TestDeltaCheckpointFreshRunsWhole(t *testing.T) {
	requireWhole(t, deltaFixture(t, graded(t, modifyTest)), flagFresh)
}

// TestDeltaCheckpointOpensTheChunkGap covers RC72: the checkout's assignment binds the build,
// and a green delta of chunk 1 opens the gap before chunk 2.
func TestDeltaCheckpointOpensTheChunkGap(t *testing.T) {
	t.Setenv("BENCH_HOME", t.TempDir())
	f := deltaFixture(t, graded(t, modifyTest), func(f *recordtest.Fixture) { commitmenttest.SeedAdmission(t, f.Root, recordtest.Spec) })
	commitmenttest.Register(t, f.Root, "build")
	commitmenttest.Admit(t, f.Root, "build", recordtest.Spec)
	requireDelta(t, f, 0, "./.")
	gaps, claim := commitmenttest.ClaimedGaps(t, f.Root)
	if want := []intent.ChunkGap{{Milestone: claim.Milestone, Outcome: claim.Outcome, Chunk: "1", Open: true}}; !reflect.DeepEqual(gaps, want) {
		t.Fatalf("gap records after a delta green = %+v, want %+v", gaps, want)
	}
}

// TestDeltaCheckpointRedStaysRed covers RC73: a red delta reds the checkpoint, records no
// green, and runs no whole gate after it.
func TestDeltaCheckpointRedStaysRed(t *testing.T) {
	f := deltaFixture(t, graded(t, func(f *recordtest.Fixture) { f.Write(plantedRed, "package fixture\n") }))
	requireDelta(t, f, 7, "./.")
	if got := InspectTree(f.Root, strippedSubject(t, f)); got.ReusableGreen {
		t.Fatalf("evidence after a red delta = %+v; want no reusable green", got)
	}
}

// TestDeltaCheckpointNarrowsOnlyThePhaseTable is the witness of the schedule rule: a gate with
// no phase-table hand-off, and a declared phase manifest, each keep the whole gate.
func TestDeltaCheckpointNarrowsOnlyThePhaseTable(t *testing.T) {
	for name, prepare := range map[string]func(*recordtest.Fixture){
		"no hand-off": func(f *recordtest.Fixture) {
			script := string(outcomeRead(t, filepath.Join(f.Root, ".bench", "gate.sh")))
			outcomeWrite(t, f.Root, ".bench/gate.sh", strings.Replace(script, phaseTableHandoff, "", 1), 0o755)
		},
		"manifest": func(f *recordtest.Fixture) {
			f.Write(canary.PhaseManifestPath, `{"phases":[{"name":"x","argv":["true"]}]}`)
		},
	} {
		t.Run(name, func(t *testing.T) { requireWhole(t, deltaFixture(t, graded(t, modifyTest), prepare)) })
	}
}

// TestPhaseChildrenDropTheDeltaSchedule is the witness that no phase child of a delta set runs
// a delta set of its own.
func TestPhaseChildrenDropTheDeltaSchedule(t *testing.T) {
	t.Setenv(deltaScheduleEnv, `["./."]`)
	env, err := gateEnv()
	if err != nil || slices.ContainsFunc(env, func(entry string) bool { return strings.HasPrefix(entry, deltaScheduleEnv+"=") }) {
		t.Fatalf("phase child environment = %q, %v; want no delta schedule", env, err)
	}
}

// scheduleWitness is the file in the Git common directory where the delta gate records the
// schedule of each run: the delta set that the run names, or wholeSchedule.
const scheduleWitness outcomeWitness = ".gate-schedules"

const (
	wholeSchedule = "whole"
	// plantedRed is the test file whose presence reds the delta gate.
	plantedRed = "planted_red_test.go"
	// phaseTableHandoff is the gate line that routes a gate through the phase table.
	phaseTableHandoff = "if false; then \"$bench\" gate-phases; fi\n"
)

// deltaGate routes the outcome gate of f through the phase table and makes it record the
// schedule of each run. The planted red file reds the gate, and go.mod makes the tree a module.
func deltaGate(t *testing.T) func(*recordtest.Fixture) {
	return func(f *recordtest.Fixture) {
		script := string(outcomeRead(t, filepath.Join(f.Root, ".bench", "gate.sh")))
		record := "printf '%s\\n' \"${" + deltaScheduleEnv + ":-" + wholeSchedule + "}\" >> \"$common/" + string(scheduleWitness) + "\"\n"
		outcomeWrite(t, f.Root, ".bench/gate.sh", script+record+phaseTableHandoff+"if [ -e "+plantedRed+" ]; then exit 7; fi\n", 0o755)
		f.Write("go.mod", "module fixture\n\ngo 1.21\n")
	}
}

// inheritFixtureBinary makes the phase-table route inherit a fixture executable whose
// freshness check passes, so no run builds Bench.
func inheritFixtureBinary(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	outcomeWrite(t, dir, "bench", "#!/bin/sh\nexit 0\n", 0o755)
	t.Setenv("BENCH_KIT", "")
	t.Setenv(runbinary.Env, filepath.Join(dir, "bench"))
	previous := prospectiveRunBinary
	prospectiveRunBinary = runbinary.Factory{Verify: func(string, string) error { return nil }}
	t.Cleanup(func() { prospectiveRunBinary = previous })
}

// frozenDeltaFixture is the frozen and reviewed chunk 1 of a two-chunk delegated plan, behind
// the delta gate. prepare runs before the freeze.
func frozenDeltaFixture(t *testing.T, prepare ...func(*recordtest.Fixture)) *recordtest.Fixture {
	t.Helper()
	inheritFixtureBinary(t)
	twoChunks := func(tb testing.TB, root string, _ int) *recordtest.Fixture {
		return recordtest.AttachDelegated(tb, root, 2)
	}
	return attachedCheckpointFixture(t, twoChunks, append([]func(*recordtest.Fixture){seedRepair, deltaGate(t)}, prepare...)...)
}

// deltaFixture is the frozen delta fixture after the confirmed repair repair.
func deltaFixture(t *testing.T, repair func(*recordtest.Fixture), prepare ...func(*recordtest.Fixture)) *recordtest.Fixture {
	t.Helper()
	f := frozenDeltaFixture(t, prepare...)
	confirmRepair(f, repairChunk(f, repair))
	return f
}

// graded runs the prepared run of the frozen chunk, so its reviewed tip holds a whole green,
// and then applies change.
func graded(t *testing.T, change func(*recordtest.Fixture)) func(*recordtest.Fixture) {
	return func(f *recordtest.Fixture) {
		requirePrepared(t, f)
		change(f)
	}
}

// addTestIn answers the repair that adds one test file in the package at operand.
func addTestIn(operand string) func(*recordtest.Fixture) {
	return func(f *recordtest.Fixture) {
		f.Write(strings.TrimPrefix(operand, "./")+"/delta_test.go", "package delta\n")
	}
}

// requireDelta requires that the chunk checkpoint of f exits code after one run of the delta
// set of packages, and that the run adds no evidence.
func requireDelta(t *testing.T, f *recordtest.Fixture, code int, packages ...string) {
	t.Helper()
	before := evidenceEntries(t, f.Root)
	got, out := runCheckpoint(t, f)
	if after := evidenceEntries(t, f.Root); got != code || after != before {
		t.Fatalf("delta checkpoint = %d %s with %d evidence entries after %d; want exit %d and no new entry", got, out, after, before, code)
	}
	set, err := json.Marshal(packages)
	if err != nil {
		t.Fatal(err)
	}
	if schedules := gateSchedules(t, f); !slices.Equal(schedules, []string{wholeSchedule, string(set)}) {
		t.Fatalf("gate schedules = %q; want the prepared whole run and then the delta set %s", schedules, set)
	}
}

// requireWhole requires that the chunk checkpoint of f, with flags, exits 0 after one more
// whole run, which leaves a reusable green of its subject.
func requireWhole(t *testing.T, f *recordtest.Fixture, flags ...string) {
	t.Helper()
	before := gateSchedules(t, f)
	code, out := runCheckpointWith(t, f, flags...)
	if schedules := gateSchedules(t, f); code != 0 || !slices.Equal(schedules, append(before, wholeSchedule)) {
		t.Fatalf("checkpoint = %d %s after the schedules %q; want one more whole run", code, out, schedules)
	}
	if got := InspectTree(f.Root, strippedSubject(t, f)); !got.ReusableGreen {
		t.Fatalf("evidence after a whole checkpoint = %+v; want a reusable green", got)
	}
}

// gateSchedules answers the schedule of each run of the delta gate of f, in run order.
func gateSchedules(t *testing.T, f *recordtest.Fixture) []string {
	t.Helper()
	data, err := os.ReadFile(scheduleWitness.path(t, f.Root))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
}

// evidenceEntries counts the records in the evidence store of root.
func evidenceEntries(t *testing.T, root string) int {
	t.Helper()
	common, err := benchgit.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(filepath.Dir(evidencePath(common, subject{})))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return len(entries)
}

// strippedSubject is the HEAD tree of f without its review record.
func strippedSubject(t *testing.T, f *recordtest.Fixture) string {
	t.Helper()
	tree, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	return tree
}

// runSchedule runs the phase command of root under kit, with the delta set that schedule
// names when it is not empty. It answers the exit, stdout, and stderr.
func runSchedule(t *testing.T, root, kit, schedule string) (int, string, string) {
	t.Helper()
	if schedule != "" {
		t.Setenv(deltaScheduleEnv, schedule)
	}
	var stdout, stderr bytes.Buffer
	code := phasesCommandAtKitWithSelection(context.Background(), root, kit, fixtureSelection(root), &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// tablePhases answers the phase cells of the green table in stdout, in table order.
func tablePhases(stdout string) []string {
	var names []string
	for _, cells := range greenTableRows(stdout) {
		names = append(names, cells[0])
	}
	return names
}

// goTest is a Go test file of package pkg that declares the test name with body.
func goTest(pkg, name, body string) string {
	return "package " + pkg + "\n\nimport \"testing\"\n\nfunc " + name + "(t *testing.T) {" + body + "}\n"
}
