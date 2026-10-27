package gate

// The delta checkpoint: after a whole green at the reviewed tip, a chunk checkpoint whose
// change since that tip is test-only runs the delta set. A green delta writes only the delta
// record, which no landing and no whole checkpoint reads. A later chunk checkpoint of the same
// base and subject reuses it.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// deltaScheduleEnv carries the package operands of the delta set, as a JSON list, to the one
// gate-phases process that runs it.
const deltaScheduleEnv = "BENCH_GATE_DELTA"

// conformancePhaseName names the delta set's run of the root-conformance entry.
const conformancePhaseName = "conformance"

// RootConformanceRunPattern is the run pattern that selects the root-conformance entry alone.
const RootConformanceRunPattern = "^" + registry.RootConformanceTest + "$"

// deltaReusedLine announces a reused delta green, so no reader takes it for a whole green.
const deltaReusedLine = "gate: green (delta green reused for this subject)"

// deltaRun is an eligible delta: the subject whose environment names the delta set, the set
// as that environment carries it, and the base tree.
type deltaRun struct {
	plan      subject
	set, base string
}

// writeDeltaGreen installs the green of run, which the executable at the content address
// runBinary graded at recordedAt, as the delta record of gitdir through replacer. The verdict
// cache and the evidence store stay untouched, so no landing and no whole checkpoint reads it.
func writeDeltaGreen(replacer durablefile.Replacer, gitdir string, run deltaRun, runBinary string, recordedAt time.Time) error {
	rec := storeRecord(run.plan.Tree)
	rec.RecordedAt, rec.RunBinary, rec.Delta, rec.BaseTree = recordTime(recordedAt), runBinary, run.set, run.base
	return durableReplaceAt(replacer, gitdir, benchgit.DeltaCacheFile, rec)
}

// deltaPlan answers the delta of a chunk checkpoint, whose plan carries the environment entry
// that names the delta set, and whether the delta is eligible. Only a chunk checkpoint carries
// its HEAD source, so every other run finds no delta. The delta needs
// every condition: the run may reuse, its whole schedule is the kit table, the change from the
// base to the subject is test-only, no changed test file needs a whole-only phase, and the
// base holds a whole green with no freshness window.
func deltaPlan(ctx context.Context, root, runtimeRoot string, plan subject, mode runMode) (deltaRun, bool) {
	checkpoint, _ := ctx.Value(checkpointKey{}).(Checkpoint)
	source, _ := ctx.Value(deltaSourceKey{}).(chunkSource)
	if mode != reuseFreshGreen || !narrowsKitTable(runtimeRoot, root, plan.Resolution) {
		return deltaRun{}, false
	}
	base, err := deltaBase(root, checkpoint, source)
	if err != nil {
		return deltaRun{}, false
	}
	files, err := reviewrecord.TestOnlyFiles(root, checkpoint.Spec, base, plan.Tree)
	if err != nil {
		return deltaRun{}, false
	}
	packages, eligible := deltaPackages(root, files)
	if !eligible || !inspectProspective(context.Background(), root, base, time.Now(), false).ReusableGreen {
		return deltaRun{}, false
	}
	// A list of strings always encodes.
	set, _ := json.Marshal(packages)
	plan.Env = mergeEnv(plan.Env, []string{deltaScheduleEnv + "=" + string(set)})
	return deltaRun{plan: plan, set: string(set), base: base}, true
}

// narrowsKitTable reports whether the whole schedule of the run is the kit table of a Go
// module, the one table that the delta set narrows. The gate routes through the phase table,
// and root, the baseline whose manifest selects a chunk checkpoint's table, has no manifest.
// A manifest that the loader refuses also keeps the whole gate.
func narrowsKitTable(runtimeRoot, root string, resolution Resolution) bool {
	_, declared, err := readManifest(manifestPath(root))
	return phaseTableGate(runtimeRoot, resolution) && GoModuleRoot(runtimeRoot) && err == nil && !declared
}

// deltaBase answers the base of the delta: the source digest of the reviewed tip of the chunk
// in the HEAD record, which the derivation finds against the HEAD commit of the subject.
func deltaBase(root string, checkpoint Checkpoint, source chunkSource) (string, error) {
	chunk, err := reviewrecord.ChunkEntry(root, source.tree, checkpoint.Spec, checkpoint.Chunk)
	if err != nil {
		return "", err
	}
	reviewed, err := reviewrecord.ReviewedAt(root, checkpoint.Spec, chunk, source.tip)
	if err != nil {
		return "", err
	}
	return reviewed.Digest, nil
}

// deltaPackages answers the package operands of the delta set: the directory of each changed
// test file, once, in file order. A test file builds with the package of its directory. It
// answers false for a test file in the system suite or in a race-table package, whose tests
// run only in the whole schedule, and for a delta with no test file.
func deltaPackages(root string, files []string) ([]string, bool) {
	wholeOnly := map[string]bool{}
	operands, _ := SystemSuite(root)
	for _, operand := range operands {
		if dir, ok := strings.CutPrefix(operand, "./"); ok {
			wholeOnly[dir] = true
		}
	}
	for _, test := range raceTests {
		wholeOnly[strings.TrimPrefix(test.PackagePath, "./")] = true
	}
	var packages []string
	for _, file := range files {
		dir := path.Dir(file)
		if wholeOnly[dir] {
			return nil, false
		}
		if operand := "./" + dir; !slices.Contains(packages, operand) {
			packages = append(packages, operand)
		}
	}
	return packages, len(packages) != 0
}

// runDelta answers steps 2 and 3 of the run decision for an eligible delta. It reuses a delta
// green of the same base and subject, and otherwise runs the oracle over the delta set. A
// green run writes the delta record of gitdir through replacer, and no pending record, no
// evidence, and no verdict cache. A red delta stays red with no whole run after it.
func runDelta(ctx, runCtx context.Context, runtimeRoot, storageRoot, gitdir string, declaredPaths []string, delta deltaRun, evaluation executionEvaluation, local []protection.Descriptor, selection *runbinary.Selection, release func() error, replacer durablefile.Replacer, stdout, stderr io.Writer) Result {
	if reusableDelta(gitdir, delta) {
		logGateEvent(ctx, gateLogRecord{Event: "gate.reused", Root: storageRoot})
		fmt.Fprintln(stdout, deltaReusedLine)
		return Result{Inspection: inspectSubjectAt(storageRoot, delta.plan, time.Now().UTC()), Schedule: ScheduleDelta}
	}
	return runReduced(ctx, runCtx, runtimeRoot, storageRoot, declaredPaths, delta.plan, ScheduleDelta, evaluation, local, selection, release, stdout, stderr, func(runBinary string) error {
		return writeDeltaGreen(replacer, gitdir, delta, runBinary, time.Now())
	})
}

// runReduced runs the oracle over the reduced schedule that plan's environment names. A green
// run passes the digest of the run binary to record, which writes the one record of the
// schedule. The run writes no pending record, no evidence, and no verdict cache, and a red
// run records nothing.
func runReduced(ctx, runCtx context.Context, runtimeRoot, storageRoot string, declaredPaths []string, plan subject, schedule Schedule, evaluation executionEvaluation, local []protection.Descriptor, selection *runbinary.Selection, release func() error, stdout, stderr io.Writer, record func(runBinary string) error) (result Result) {
	defer func() { result.Schedule = schedule }()
	logGateEvent(ctx, gateLogRecord{Event: "oracle.start", Root: runtimeRoot})
	run := runCaptured(runCtx, runtimeRoot, plan, stdout, stderr, local...)
	rc := run.Code
	logGateEvent(ctx, gateLogRecord{Event: "oracle.finish", Root: runtimeRoot, Exit: &rc})
	if code, _, interrupted := interruptedExit(ctx, runCtx, rc, stderr); interrupted {
		return Result{GateExit: code, ActionExit: code, Inspection: inspectAt(storageRoot, time.Now().UTC())}
	}
	// An eligible reduced schedule runs the phase table, so the run selected its executable.
	// The release in settleRun can remove those bytes, so the digest reads first.
	runBinary, err := fileDigest(selection.Path)
	after, refusal, refused := settleRun(evaluation, runtimeRoot, storageRoot, declaredPaths, plan, run, release, stderr)
	if refused {
		return refusal
	}
	if rc == 0 {
		if err == nil {
			err = record(runBinary)
		}
		if err != nil {
			return persistenceFailure(storageRoot, rc, stderr, "gate "+schedule.String()+" persistence failed", err)
		}
	}
	return Result{GateExit: rc, ActionExit: rc, Inspection: inspectSubjectAt(storageRoot, after, time.Now().UTC())}
}

// reusableDelta reports whether the delta record of gitdir is a green of the base and the
// subject of delta. Only the delta class carries a base tree, so no other class matches.
func reusableDelta(gitdir string, delta deltaRun) bool {
	loaded := loadVerdict(filepath.Join(gitdir, benchgit.DeltaCacheFile), time.Now().UTC())
	return loaded.state == Ready && loaded.record.BaseTree == delta.base && loaded.record.Tree == delta.plan.Tree
}

// interruptedExit classifies an oracle run that a cancel of the gate or the gate timeout ended.
// A cancel answers the child's code. A timeout prints its line and answers 124. A whole run and
// a reduced run both read it, so they end alike. interrupted is false when neither ended the run.
func interruptedExit(ctx, runCtx context.Context, rc int, stderr io.Writer) (code int, timedOut, interrupted bool) {
	if ctx.Err() != nil {
		return rc, false, true
	}
	if errors.Is(context.Cause(runCtx), errGateTimeout) {
		fmt.Fprintln(stderr, "gate: timeout")
		return 124, true, true
	}
	return rc, false, false
}

// deltaScheduled reports whether this process runs the delta set.
func deltaScheduled() bool { return os.Getenv(deltaScheduleEnv) != "" }

// scheduledTable answers the reader set or the delta set when this process runs it, and the
// phase table otherwise. A malformed schedule, or an empty delta schedule, refuses.
func scheduledTable(root, kit string) ([]Phase, error) {
	if readerScheduled() {
		return readerScheduledTable(root, kit)
	}
	if !deltaScheduled() {
		return phaseTable(root, kit)
	}
	var packages []string
	if err := json.Unmarshal([]byte(os.Getenv(deltaScheduleEnv)), &packages); err != nil || len(packages) == 0 {
		return nil, fmt.Errorf("gate: malformed delta schedule in %s", deltaScheduleEnv)
	}
	return deltaPhases(root, kit, packages), nil
}

// deltaPhases is the delta set: gofmt and vet over the whole root, the test phase over
// packages, and the root-conformance entry alone where the kit runs it.
func deltaPhases(root, kit string, packages []string) []Phase {
	phases := append(sourcePhases(root), testPhase(root, kit, packages...))
	if RootConformanceRuns(root, kit) {
		conformance := testPhase(root, kit, "./"+conformancePackagePath, "-run", RootConformanceRunPattern)
		conformance.Name = conformancePhaseName
		phases = append(phases, conformance)
	}
	return phases
}
