package gate

// The delta checkpoint: after a whole green at the reviewed tip, a chunk checkpoint whose
// change since that tip is test-only runs the delta set. The oracle runs over the delta set
// alone and records nothing, so no landing and no later checkpoint finds a delta green.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// deltaScheduleEnv carries the package operands of the delta set, as a JSON list, to the one
// gate-phases process that runs it.
const deltaScheduleEnv = "BENCH_GATE_DELTA"

// conformancePhaseName names the delta set's run of the root-conformance entry.
const conformancePhaseName = "conformance"

// RootConformanceRunPattern is the run pattern that selects the root-conformance entry alone.
const RootConformanceRunPattern = "^" + registry.RootConformanceTest + "$"

// deltaPlan answers plan with the environment entry that names the delta set of a chunk
// checkpoint, and whether the delta is eligible. Only a chunk checkpoint that accepts its
// committed source carries that source, so every other run finds no delta. The delta needs
// every condition: the run may reuse, its whole schedule is the kit table, the change from the
// base to the subject is test-only, no changed test file needs a whole-only phase, and the
// base holds a whole green with no freshness window.
func deltaPlan(ctx context.Context, root, runtimeRoot string, plan subject, mode runMode) (subject, bool) {
	checkpoint, _ := ctx.Value(checkpointKey{}).(Checkpoint)
	source, _ := ctx.Value(chunkSourceKey{}).(chunkSource)
	if mode != reuseFreshGreen || !narrowsKitTable(runtimeRoot, root, plan.Resolution) {
		return plan, false
	}
	base, err := deltaBase(root, checkpoint, source)
	if err != nil {
		return plan, false
	}
	files, err := reviewrecord.TestOnlyFiles(root, checkpoint.Spec, base, plan.Tree)
	if err != nil {
		return plan, false
	}
	packages, eligible := deltaPackages(root, files)
	if !eligible || !inspectProspective(context.Background(), root, base, time.Now(), false).ReusableGreen {
		return plan, false
	}
	// A list of strings always encodes.
	set, _ := json.Marshal(packages)
	plan.Env = mergeEnv(plan.Env, []string{deltaScheduleEnv + "=" + string(set)})
	return plan, true
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

// runDelta runs the oracle over the delta set. It writes no pending record, no evidence, and
// no verdict cache, and a red delta stays red with no whole run after it.
func runDelta(ctx, runCtx context.Context, runtimeRoot, storageRoot string, declaredPaths []string, plan subject, evaluation executionEvaluation, local []protection.Descriptor, release func() error, stdout, stderr io.Writer) Result {
	logGateEvent(ctx, gateLogRecord{Event: "oracle.start", Root: runtimeRoot})
	run := runCaptured(runCtx, runtimeRoot, plan, stdout, stderr, local...)
	rc := run.Code
	logGateEvent(ctx, gateLogRecord{Event: "oracle.finish", Root: runtimeRoot, Exit: &rc})
	if code, _, interrupted := interruptedExit(ctx, runCtx, rc, stderr); interrupted {
		return Result{GateExit: code, ActionExit: code, Inspection: inspectAt(storageRoot, time.Now().UTC())}
	}
	after, refusal, refused := settleRun(evaluation, runtimeRoot, storageRoot, declaredPaths, plan, run, release, stderr)
	if refused {
		return refusal
	}
	return Result{GateExit: rc, ActionExit: rc, Inspection: inspectSubjectAt(storageRoot, after, time.Now().UTC())}
}

// interruptedExit classifies an oracle run that a cancel of the gate or the gate timeout ended.
// A cancel answers the child's code. A timeout prints its line and answers 124. A whole run and
// a delta run both read it, so they end alike. interrupted is false when neither ended the run.
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

// scheduledTable answers the delta set when this process runs it, and the phase table
// otherwise. A malformed or empty schedule refuses.
func scheduledTable(root, kit string) ([]Phase, error) {
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
