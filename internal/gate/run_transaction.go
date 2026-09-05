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
	"sync"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/runbinary"
)

type runBinaryOwner func(context.Context, string) (*runbinary.Selection, error)

var executionLockOwners = struct {
	sync.Mutex
	paths map[string]bool
}{paths: map[string]bool{}}

func productionRunBinaryOwner() runBinaryOwner {
	if _, inherited := os.LookupEnv(runbinary.Env); inherited {
		return runbinary.ReuseOrOwn
	}
	if os.Getenv("BENCH_WRAPPER") != "" {
		return runbinary.Own
	}
	return nil
}

// baseline, when set, is the landing destination whose phase schedule grades the
// prospective tree. It reaches the oracle as one environment entry, so the candidate's
// own phase manifest never selects what runs against it.
func executeSubjectWithRunBinary(ctx context.Context, runtimeRoot, storageRoot string, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode, evaluation executionEvaluation, owner runBinaryOwner, baseline string) Result {
	plan, err := evaluation.acceptPre()
	if err != nil {
		return operational(storageRoot, 0, stderr, fmt.Sprintf("gate subject unavailable: %v", err))
	}
	var declaredPaths []string
	if m, _, reason := loadManifest(runtimeRoot); reason == "" {
		declaredPaths = append(declaredPaths, m.Paths...)
	}
	decision := Decide(DecisionInput{Subject: plan.Tree, Resolution: plan.Resolution})
	if plan.Resolution.Kind == None {
		fmt.Fprintln(stderr, "no gate found: add an executable .bench/gate.sh or set BENCH_GATE")
		return Result{GateExit: 3, ActionExit: 3, Inspection: inspectAt(storageRoot, time.Now().UTC())}
	}
	if !decision.Accepted {
		return operational(storageRoot, 0, stderr, "gate decision refused: "+decision.Refusal)
	}
	logGateEvent(ctx, gateLogRecord{Event: "subject.accepted", Root: storageRoot, Mode: mode.String(), Detail: plan.Tree})
	if mode == reuseFreshGreen {
		if reuse := reusableEvidence(storageRoot, plan, time.Now().UTC()); reuse.ReusableGreen {
			logGateEvent(ctx, gateLogRecord{Event: "gate.reused", Root: storageRoot})
			return reusedGreenResult(stdout, reuse)
		}
	}
	gitdir, err := benchgit.AdminDir(storageRoot)
	if err != nil {
		return operational(storageRoot, 0, stderr, "git directory unavailable")
	}
	lock, err := os.OpenFile(filepath.Join(gitdir, "bench-gate.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		persistInterruptedIfGreen(storageRoot, gitdir, plan)
		return operational(storageRoot, 0, stderr, "gate lock unavailable")
	}
	defer lock.Close()
	if err := acquireExecutionLock(lock); err != nil {
		persistInterruptedIfGreen(storageRoot, gitdir, plan)
		fmt.Fprintln(stderr, "gate execution already in progress")
		writeOwnerDiagnostic(stderr, filepath.Join(gitdir, "bench-gate-owner"))
		inspection := inspectAt(storageRoot, time.Now().UTC())
		inspection.ReusableGreen = false
		return Result{ActionExit: 1, Inspection: inspection}
	}
	defer unlockExecutionLock(lock)
	logGateEvent(ctx, gateLogRecord{Event: "gate.locked", Root: storageRoot, Path: filepath.Join(gitdir, "bench-gate.lock")})
	if arm != nil {
		var stop func()
		ctx, stop = arm(ctx)
		defer stop()
	}
	ownerPath := filepath.Join(gitdir, "bench-gate-owner")
	defer func() { _ = os.Remove(ownerPath) }()
	if err := os.WriteFile(ownerPath, ownerRecord(time.Now().UTC()), 0o600); err != nil {
		return operational(storageRoot, 0, stderr, "gate owner persistence failed")
	}
	underLock, err := evaluation.validatePre()
	if err != nil || !sameSubject(plan, underLock) {
		return operational(storageRoot, 0, stderr, "gate subject changed before execution")
	}
	// A reusable green is answered from the record without touching it. Re-recording the
	// verdict would push RecordedAt forward on every read and make the freshness window
	// unbounded. The check sits ahead of the pending replace, so a reuse returns with
	// nothing written: no pending record to leave behind, no verdict to restore.
	if mode == reuseFreshGreen {
		if reuse := reusableEvidence(storageRoot, plan, time.Now().UTC()); reuse.ReusableGreen {
			logGateEvent(ctx, gateLogRecord{Event: "gate.reused", Root: storageRoot})
			return reusedGreenResult(stdout, reuse)
		}
	}
	runCtx, cancelRun := bounds.ContextCause(ctx, gateTimeout, errGateTimeout)
	defer cancelRun()
	// The cache lock spans the whole run: it opens before the run-binary build, the first
	// Go child of the run, and the deferred release closes it after every teardown above.
	// A clean that arrives between those two points refuses instead of removing an archive
	// a phase is reading. A lock this run cannot take refuses the run: a gate that graded
	// the tree unlocked would compile beside a clean removing the archives it is writing.
	// A closure that names no home is the one exception. It derives no Bench cache, so its
	// children write outside every directory a clean can reach, and there is no hold to
	// take. Release answers a nil holder, so that run needs no second branch here.
	holder, err := gocache.Hold(plan.Env)
	if err != nil {
		logGateEvent(ctx, gateLogRecord{Event: "cache.lock.unavailable", Root: storageRoot, Detail: err.Error()})
		if gocache.Declared(plan.Env) {
			return operational(storageRoot, 0, stderr, gocache.Refusal(plan.Env, err))
		}
	}
	defer holder.Release()
	var selection *runbinary.Selection
	if owner != nil && phaseTableGate(runtimeRoot, plan.Resolution) {
		source := runBinarySource(runtimeRoot, storageRoot, plan)
		logGateEvent(ctx, gateLogRecord{Event: "binary.select.start", Root: source})
		selection, err = owner(runCtx, source)
		if err != nil {
			logGateEvent(ctx, gateLogRecord{Event: "binary.select.finish", Root: source, Detail: err.Error()})
			if errors.Is(context.Cause(runCtx), errGateTimeout) {
				fmt.Fprintln(stderr, "gate: timeout")
				return Result{GateExit: 124, ActionExit: 124, Inspection: inspectAt(storageRoot, time.Now().UTC())}
			}
			if ctx.Err() != nil {
				return Result{GateExit: 130, ActionExit: 130, Inspection: inspectAt(storageRoot, time.Now().UTC())}
			}
			return operational(storageRoot, 0, stderr, "gate Bench executable unavailable")
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
	if baseline != "" {
		plan.Env = mergeEnv(plan.Env, []string{baselinePolicyEnv + "=" + baseline})
	}
	plan.Env = withGateRunLogEnv(ctx, plan.Env)
	// Persist the interrupted posture before the oracle starts. The branch-native gate
	// always grades the complete subject; there is no component partition to retain.
	pending := interruptedRecord(plan, time.Now().UTC())
	if err := durableReplace(gitdir, pending); err != nil {
		_ = durableReplace(gitdir, pending)
		return operational(storageRoot, 0, stderr, "gate pending persistence failed")
	}
	logGateEvent(ctx, gateLogRecord{Event: "oracle.start", Root: runtimeRoot})
	rc := runCaptured(runCtx, runtimeRoot, plan, stdout, stderr)
	logGateEvent(ctx, gateLogRecord{Event: "oracle.finish", Root: runtimeRoot, Exit: &rc})
	if ctx.Err() != nil {
		return Result{GateExit: rc, ActionExit: rc, Inspection: inspectAt(storageRoot, time.Now().UTC())}
	}
	if errors.Is(context.Cause(runCtx), errGateTimeout) {
		fmt.Fprintln(stderr, "gate: timeout")
		if err := invalidateEvidence(storageRoot, plan); err != nil {
			return operational(storageRoot, 124, stderr, "gate evidence invalidation failed")
		}
		ready := verdictRecord{Schema: 1, State: Ready, Status: "timeout", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)}
		if err := durableReplace(gitdir, ready); err != nil {
			_ = durableReplace(gitdir, pending)
			return operational(storageRoot, 124, stderr, "gate timeout persistence failed")
		}
		return Result{GateExit: 124, ActionExit: 124, Inspection: inspectAt(storageRoot, time.Now().UTC())}
	}
	after, err := evaluation.capturePost()
	if err != nil || !sameSubject(plan, after) {
		fmt.Fprintln(stderr, subjectChangeDiagnostic(runtimeRoot, declaredPaths, plan, after, err))
		inspection := inspectAt(storageRoot, time.Now().UTC())
		if err == nil {
			inspection = inspectSubjectAt(storageRoot, after, time.Now().UTC())
		}
		return Result{GateExit: rc, ActionExit: 1, Inspection: inspection}
	}
	status := "red"
	if rc == 0 {
		status = "green"
	}
	recordedAt := time.Now().UTC()
	ready := verdictRecord{Schema: 1, State: Ready, Status: status, Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: recordedAt.UTC().Truncate(time.Second).Format(time.RFC3339)}
	if status == "green" {
		if err := retainGreen(storageRoot, plan, recordedAt); err != nil {
			fmt.Fprintln(stderr, "gate evidence persistence failed")
			return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, time.Now().UTC())}
		}
	} else {
		if err := invalidateEvidence(storageRoot, plan); err != nil {
			return operational(storageRoot, rc, stderr, "gate evidence invalidation failed")
		}
	}
	if err := durableReplace(gitdir, ready); err != nil {
		_ = durableReplace(gitdir, pending)
		fmt.Fprintln(stderr, "gate final persistence failed")
		return Result{GateExit: rc, ActionExit: 1, Inspection: inspectAt(storageRoot, time.Now().UTC())}
	}
	return Result{GateExit: rc, ActionExit: rc, Inspection: inspectSubjectAt(storageRoot, after, time.Now().UTC())}
}

func persistInterruptedIfGreen(root, gitdir string, plan subject) {
	if !inspectAt(root, time.Now().UTC()).ReusableGreen {
		return
	}
	pending := interruptedRecord(plan, time.Now().UTC())
	if err := durableReplace(gitdir, pending); err != nil {
		_ = durableReplace(gitdir, pending)
	}
}

func acquireExecutionLock(f *os.File) error {
	executionLockOwners.Lock()
	defer executionLockOwners.Unlock()
	if executionLockOwners.paths[f.Name()] {
		return syscall.EAGAIN
	}
	lock := gocache.RecordLock(syscall.F_WRLCK)
	if err := syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &lock); err != nil {
		return err
	}
	executionLockOwners.paths[f.Name()] = true
	return nil
}

func unlockExecutionLock(f *os.File) error {
	executionLockOwners.Lock()
	defer executionLockOwners.Unlock()
	lock := gocache.RecordLock(syscall.F_UNLCK)
	err := syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &lock)
	delete(executionLockOwners.paths, f.Name())
	return err
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

func subjectChangeDiagnostic(root string, declaredPaths []string, before, after subject, err error) string {
	if err == nil && before.Tree != after.Tree {
		if changed, ok := benchgit.ChangedPathsBetweenTrees(root, before.Tree, after.Tree); ok {
			for _, path := range declaredPaths {
				for _, moved := range changed {
					if moved == path || strings.HasPrefix(moved, path+"/") {
						return "gate subject changed during execution: declared input " + path
					}
				}
			}
		}
		return "gate subject changed during execution: tree " + after.Tree
	}
	return "gate subject changed during execution"
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
	lock := gocache.RecordLock(syscall.F_RDLCK)
	if err := syscall.FcntlFlock(f.Fd(), syscall.F_GETLK, &lock); err != nil {
		return false, err
	}
	return lock.Type != syscall.F_UNLCK, nil
}
