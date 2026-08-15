package gate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
)

type runBinaryOwner func(context.Context, string) (*runbinary.Selection, error)

func productionRunBinaryOwner() runBinaryOwner {
	if _, inherited := os.LookupEnv(runbinary.Env); inherited {
		return runbinary.ReuseOrOwn
	}
	if os.Getenv("BENCH_WRAPPER") != "" {
		return runbinary.Own
	}
	return nil
}

func executeSubjectWithRunBinary(ctx context.Context, runtimeRoot, storageRoot string, stdout, stderr io.Writer, engine gateEngine, arm postAcquireContextArm, mode runMode, evaluation executionEvaluation, owner runBinaryOwner) Result {
	plan, err := evaluation.acceptPre()
	if err != nil {
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate subject unavailable")
	}
	decision := Decide(DecisionInput{Subject: plan.Tree, Resolution: plan.Resolution})
	if plan.Resolution.Kind == None {
		fmt.Fprintln(stderr, "no gate found: add an executable .bench/gate.sh or set BENCH_GATE")
		return Result{GateExit: 3, ActionExit: 3, Inspection: inspectAt(storageRoot, engine.Now())}
	}
	if !decision.Accepted {
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate decision refused: "+decision.Refusal)
	}
	logGateEvent(ctx, gateLogRecord{Event: "subject.accepted", Root: storageRoot, Mode: mode.String(), Detail: plan.Tree})
	if mode == reuseFreshGreen {
		if reuse := reusableEvidence(storageRoot, plan, engine.Now()); reuse.ReusableGreen {
			logGateEvent(ctx, gateLogRecord{Event: "gate.reused", Root: storageRoot})
			return reusedGreenResult(stdout, reuse)
		}
	}
	gitdir, err := engine.GitDir(storageRoot)
	if err != nil {
		return operationalWithEngine(engine, storageRoot, 0, stderr, "git directory unavailable")
	}
	lock, err := engine.OpenLock(filepath.Join(gitdir, "bench-gate.lock"))
	if err != nil {
		persistInterruptedIfGreen(engine, storageRoot, gitdir, plan)
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate lock unavailable")
	}
	defer lock.Close()
	if err := engine.Acquire(lock); err != nil {
		persistInterruptedIfGreen(engine, storageRoot, gitdir, plan)
		fmt.Fprintln(stderr, "gate execution already in progress")
		writeOwnerDiagnostic(stderr, filepath.Join(gitdir, "bench-gate-owner"))
		inspection := inspectAt(storageRoot, engine.Now())
		inspection.ReusableGreen = false
		return Result{ActionExit: 1, Inspection: inspection}
	}
	defer engine.Unlock(lock)
	logGateEvent(ctx, gateLogRecord{Event: "gate.locked", Root: storageRoot, Path: filepath.Join(gitdir, "bench-gate.lock")})
	if arm != nil {
		var stop func()
		ctx, stop = arm(ctx)
		defer stop()
	}
	ownerPath := filepath.Join(gitdir, "bench-gate-owner")
	defer func() { _ = engine.Remove(ownerPath) }()
	if err := engine.WriteFile(ownerPath, ownerRecord(engine.Now()), 0o600); err != nil {
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate owner persistence failed")
	}
	underLock, err := evaluation.validatePre()
	if err != nil || !sameSubject(plan, underLock) {
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate subject changed before execution")
	}
	// A reusable green is answered from the record without touching it: re-recording the
	// verdict would push RecordedAt forward on every read and make the freshness window
	// unbounded. The check sits ahead of the pending replace so a reuse returns with nothing
	// written — no pending record to leave behind, no verdict to restore.
	if mode == reuseFreshGreen {
		if reuse := reusableEvidence(storageRoot, plan, engine.Now()); reuse.ReusableGreen {
			logGateEvent(ctx, gateLogRecord{Event: "gate.reused", Root: storageRoot})
			return reusedGreenResult(stdout, reuse)
		}
	}
	runCtx, cancelRun := bounds.ContextCause(ctx, gateTimeout, errGateTimeout)
	defer cancelRun()
	var selection *runbinary.Selection
	if owner != nil && phaseTableGate(runtimeRoot, plan.Resolution) {
		source := runBinarySource(runtimeRoot, storageRoot, plan)
		logGateEvent(ctx, gateLogRecord{Event: "binary.select.start", Root: source})
		selection, err = owner(runCtx, source)
		if err != nil {
			logGateEvent(ctx, gateLogRecord{Event: "binary.select.finish", Root: source, Detail: err.Error()})
			if errors.Is(context.Cause(runCtx), errGateTimeout) {
				fmt.Fprintln(stderr, "gate: timeout")
				return Result{GateExit: 124, ActionExit: 124, Inspection: inspectAt(storageRoot, engine.Now())}
			}
			if ctx.Err() != nil {
				return Result{GateExit: 130, ActionExit: 130, Inspection: inspectAt(storageRoot, engine.Now())}
			}
			return operationalWithEngine(engine, storageRoot, 0, stderr, "gate Bench executable unavailable")
		}
		logGateEvent(ctx, gateLogRecord{Event: "binary.select.finish", Root: source, Path: selection.Path})
		defer func() {
			err := selection.Close()
			record := gateLogRecord{Event: "binary.cleanup", Path: selection.Path}
			if err != nil {
				record.Detail = err.Error()
			}
			logGateEvent(ctx, record)
		}()
		plan.Env = runbinary.WithEnv(plan.Env, selection.Path)
		plan.Env = mergeEnv(plan.Env, []string{"BENCH_KIT=" + selection.SourceRoot})
	}
	plan.Env = withGateRunLogEnv(ctx, plan.Env)
	// Persist the interrupted posture before the oracle starts. The branch-native gate
	// always grades the complete subject; there is no component partition to retain.
	pending := interruptedRecord(plan, engine.Now())
	if err := durableReplaceWithEngine(engine, gitdir, pending); err != nil {
		_ = durableReplaceWithEngine(engine, gitdir, pending)
		return operationalWithEngine(engine, storageRoot, 0, stderr, "gate pending persistence failed")
	}
	logGateEvent(ctx, gateLogRecord{Event: "oracle.start", Root: runtimeRoot})
	rc := runCaptured(runCtx, runtimeRoot, plan, stdout, stderr)
	logGateEvent(ctx, gateLogRecord{Event: "oracle.finish", Root: runtimeRoot, Exit: &rc})
	if ctx.Err() != nil {
		return Result{GateExit: rc, ActionExit: rc, Inspection: inspectAt(storageRoot, engine.Now())}
	}
	if errors.Is(context.Cause(runCtx), errGateTimeout) {
		fmt.Fprintln(stderr, "gate: timeout")
		if err := invalidateEvidence(storageRoot, plan); err != nil {
			return operationalWithEngine(engine, storageRoot, 124, stderr, "gate evidence invalidation failed")
		}
		ready := verdictRecord{Schema: 1, State: Ready, Status: "timeout", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: engine.Now().UTC().Truncate(time.Second).Format(time.RFC3339)}
		if err := durableReplaceWithEngine(engine, gitdir, ready); err != nil {
			_ = durableReplaceWithEngine(engine, gitdir, pending)
			return operationalWithEngine(engine, storageRoot, 124, stderr, "gate timeout persistence failed")
		}
		return Result{GateExit: 124, ActionExit: 124, Inspection: inspectAt(storageRoot, engine.Now())}
	}
	after, err := evaluation.capturePost()
	if err != nil || !sameSubject(plan, after) {
		fmt.Fprintln(stderr, "gate subject changed during execution")
		inspection := inspectAt(storageRoot, engine.Now())
		if err == nil {
			inspection = inspectSubjectAt(storageRoot, after, engine.Now())
		}
		return Result{GateExit: rc, ActionExit: 1, Inspection: inspection}
	}
	status := "red"
	if rc == 0 {
		status = "green"
	}
	recordedAt := engine.Now()
	ready := verdictRecord{Schema: 1, State: Ready, Status: status, Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: recordedAt.UTC().Truncate(time.Second).Format(time.RFC3339)}
	if status == "green" {
		if err := retainGreen(storageRoot, plan, recordedAt); err != nil {
			fmt.Fprintln(stderr, "gate evidence persistence failed")
			return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, engine.Now())}
		}
	} else {
		if err := invalidateEvidence(storageRoot, plan); err != nil {
			return operationalWithEngine(engine, storageRoot, rc, stderr, "gate evidence invalidation failed")
		}
	}
	if err := durableReplaceWithEngine(engine, gitdir, ready); err != nil {
		_ = durableReplaceWithEngine(engine, gitdir, pending)
		fmt.Fprintln(stderr, "gate final persistence failed")
		return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, engine.Now())}
	}
	return Result{GateExit: rc, ActionExit: rc, Inspection: inspectSubjectAt(storageRoot, after, engine.Now())}
}

func persistInterruptedIfGreen(engine gateEngine, root, gitdir string, plan subject) {
	if !inspectAt(root, engine.Now()).ReusableGreen {
		return
	}
	pending := interruptedRecord(plan, engine.Now())
	if err := durableReplaceWithEngine(engine, gitdir, pending); err != nil {
		_ = durableReplaceWithEngine(engine, gitdir, pending)
	}
}

func retainGreen(root string, plan subject, recordedAt time.Time) error {
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		return err
	}
	dir := filepath.Join(gitdir, "bench-gate-evidence")
	if err := ensureEvidenceDir(gitdir, dir); err != nil {
		return err
	}
	record := verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: recordedAt.UTC().Truncate(time.Second).Format(time.RFC3339)}
	return durableReplaceAt(dir, evidenceName(plan), record)
}

func ensureEvidenceDir(parent, dir string) error {
	if err := os.Mkdir(dir, 0o700); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() || info.Mode().Perm() != 0o700 {
		return errors.New("invalid evidence directory")
	}
	parentFile, err := os.Open(parent)
	if err != nil {
		return err
	}
	defer parentFile.Close()
	return parentFile.Sync()
}

func invalidateEvidence(root string, plan subject) error {
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		return err
	}
	dir := filepath.Join(gitdir, "bench-gate-evidence")
	if err := os.Remove(evidencePath(gitdir, plan)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	file, err := os.Open(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}

func ownerRecord(now time.Time) []byte {
	return []byte(strconv.Itoa(os.Getpid()) + " " + now.UTC().Truncate(time.Second).Format(time.RFC3339) + "\n")
}

func writeOwnerDiagnostic(stderr io.Writer, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	fields := strings.Fields(string(data))
	if len(fields) != 2 {
		return
	}
	pid, err := strconv.Atoi(fields[0])
	if err != nil || pid <= 0 {
		return
	}
	if _, err := time.Parse(time.RFC3339, fields[1]); err != nil {
		return
	}
	liveness := "alive"
	if err := syscall.Kill(pid, 0); err != nil && err != syscall.EPERM {
		liveness = "not alive"
	}
	fmt.Fprintf(stderr, "gate owner: pid %d (%s)\n", pid, liveness)
}

func interruptedRecord(plan subject, now time.Time) verdictRecord {
	return verdictRecord{Schema: 1, State: Pending, Tree: plan.Tree, Oracle: plan.Oracle, StartedAt: now.UTC().Truncate(time.Second).Format(time.RFC3339), OwnerPID: os.Getpid()}
}

func sameSubject(a, b subject) bool {
	return a.Tree == b.Tree && a.Oracle == b.Oracle && a.Resolution == b.Resolution && a.Closed == b.Closed && a.Reason == b.Reason
}

func runCaptured(ctx context.Context, root string, s subject, stdout, stderr io.Writer) int {
	return runResolved(ctx, root, s.Resolution, s.Env, controlSafeWriter{stdout}, controlSafeWriter{stderr}, true).Code
}

func lockHeld(gitdir string) (bool, error) {
	path := filepath.Join(gitdir, "bench-gate.lock")
	executionLockOwners.Lock()
	defer executionLockOwners.Unlock()
	if executionLockOwners.paths[path] {
		return true, nil
	}
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer f.Close()
	lock := recordLock(syscall.F_RDLCK)
	if err := syscall.FcntlFlock(f.Fd(), syscall.F_GETLK, &lock); err != nil {
		return false, err
	}
	return lock.Type != syscall.F_UNLCK, nil
}
