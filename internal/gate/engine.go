package gate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/runbinary"
)

type gateEngine interface {
	Now() time.Time
	BuildSubject(string) (subject, error)
	PostRunSubject(string) (subject, error)
	GitDir(string) (string, error)
	OpenLock(string) (gateFile, error)
	Acquire(gateFile) error
	Unlock(gateFile) error
	CreateTemp(string, string) (gateFile, error)
	Rename(string, string) error
	OpenDir(string) (gateFile, error)
	WriteFile(string, []byte, os.FileMode) error
	Remove(string) error
}

type productionGateEngine struct{}

var executionLockOwners = struct {
	sync.Mutex
	paths map[string]bool
}{paths: map[string]bool{}}

func recordLock(typ int16) syscall.Flock_t {
	return syscall.Flock_t{Type: typ, Whence: int16(io.SeekStart), Start: 0, Len: 0}
}

func (productionGateEngine) Now() time.Time                              { return time.Now().UTC() }
func (productionGateEngine) BuildSubject(root string) (subject, error)   { return buildSubject(root) }
func (productionGateEngine) PostRunSubject(root string) (subject, error) { return buildSubject(root) }
func (productionGateEngine) GitDir(root string) (string, error) {
	return benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
}
func (productionGateEngine) OpenLock(path string) (gateFile, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
}
func (productionGateEngine) Acquire(f gateFile) error {
	executionLockOwners.Lock()
	defer executionLockOwners.Unlock()
	if executionLockOwners.paths[f.Name()] {
		return syscall.EAGAIN
	}
	lock := recordLock(syscall.F_WRLCK)
	if err := syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &lock); err != nil {
		return err
	}
	executionLockOwners.paths[f.Name()] = true
	return nil
}
func (productionGateEngine) Unlock(f gateFile) error {
	executionLockOwners.Lock()
	defer executionLockOwners.Unlock()
	lock := recordLock(syscall.F_UNLCK)
	err := syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, &lock)
	delete(executionLockOwners.paths, f.Name())
	return err
}
func (productionGateEngine) CreateTemp(dir, pattern string) (gateFile, error) {
	return os.CreateTemp(dir, pattern)
}
func (productionGateEngine) Rename(oldpath, newpath string) error  { return os.Rename(oldpath, newpath) }
func (productionGateEngine) OpenDir(path string) (gateFile, error) { return os.Open(path) }
func (productionGateEngine) WriteFile(path string, data []byte, mode os.FileMode) error {
	return os.WriteFile(path, data, mode)
}
func (productionGateEngine) Remove(path string) error { return os.Remove(path) }

// EvidenceInspection reports whether retained gate evidence can authorize one tree.
type EvidenceInspection struct {
	Tree          string
	Oracle        string
	RecordedAt    time.Time
	ReusableGreen bool
	Reason        string
}

// InspectTree reports retained exact evidence for an unpublished Git tree.
func InspectTree(root, tree string) EvidenceInspection {
	return inspectProspective(root, tree, time.Now())
}

// ExecuteTree runs or reuses the gate for an unpublished Git tree.
func ExecuteTree(ctx context.Context, root, tree string, stdout, stderr io.Writer) Result {
	ctx, finishLog := beginGateRunLog(ctx, root, stderr, "prospective")
	result := executeTreeWithOwner(ctx, root, tree, stdout, stderr, runbinary.ReuseOrOwn)
	finishLog(result)
	return result
}

func executeTreeWithOwner(ctx context.Context, root, tree string, stdout, stderr io.Writer, owner runBinaryOwner) Result {
	checkout, cleanup, err := prospectiveCheckout(root, tree)
	if err != nil {
		fmt.Fprintln(stderr, "prospective gate subject unavailable")
		return Result{ActionExit: 1}
	}
	defer cleanup()
	evaluation := newProspectiveTreeEvaluation(checkout, root, tree)
	return executeSubjectWithRunBinary(ctx, checkout, root, stdout, stderr, productionGateEngine{}, nil, reuseFreshGreen, evaluation, owner)
}

// ValidateProjectGreen reports whether branch's tip and marker have retained exact green evidence.
func ValidateProjectGreen(root, branch string) EvidenceInspection {
	current, err := benchgit.Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil || current != branch {
		return EvidenceInspection{Reason: "working branch changed"}
	}
	tip, err := benchgit.Output("-C", root, "rev-parse", "HEAD")
	if err != nil {
		return EvidenceInspection{Reason: "working tip unavailable"}
	}
	marker, err := benchgit.Output("-C", root, "rev-parse", "--verify", "refs/bench/green/"+branch)
	if err != nil || marker != tip {
		return EvidenceInspection{Reason: "project-green marker changed"}
	}
	plan, err := buildSubject(root)
	if err != nil {
		return EvidenceInspection{Reason: "subject unavailable"}
	}
	return inspectEvidence(root, plan, time.Now())
}

func inspectProspective(root, tree string, now time.Time) EvidenceInspection {
	checkout, cleanup, err := prospectiveCheckout(root, tree)
	if err != nil {
		return EvidenceInspection{Reason: "subject unavailable"}
	}
	defer cleanup()
	plan, err := buildProspectiveSubjectFor(checkout, root)
	if err != nil || plan.Tree != tree {
		return EvidenceInspection{Reason: "subject unavailable"}
	}
	return inspectEvidence(root, plan, now)
}

func inspectEvidence(root string, plan subject, now time.Time) EvidenceInspection {
	return inspectEvidenceWindowed(root, plan, now, true)
}

func inspectEvidenceWindowed(root string, plan subject, now time.Time, expires bool) EvidenceInspection {
	inspection := EvidenceInspection{Tree: plan.Tree, Oracle: plan.Oracle, Reason: plan.Reason}
	if !plan.Closed {
		return inspection
	}
	gitdir, err := commonGitDir(root)
	if err != nil {
		inspection.Reason = "evidence unavailable"
		return inspection
	}
	loaded := loadVerdict(evidencePath(gitdir, plan), now)
	if loaded.state != Ready || loaded.record.Status != "green" {
		if loaded.state == Ready {
			inspection.Reason = "recorded " + loaded.record.Status
		} else if loaded.state != Absent {
			inspection.Reason = "evidence unavailable"
		} else {
			inspection.Reason = "evidence absent"
		}
		return inspection
	}
	if loaded.record.Tree != plan.Tree || loaded.record.Oracle != plan.Oracle {
		inspection.Reason = "evidence changed"
		return inspection
	}
	recorded, _ := time.Parse(time.RFC3339, loaded.record.RecordedAt)
	inspection.RecordedAt = recorded
	if expires && now.Sub(recorded) >= freshness {
		inspection.Reason = "verdict expired"
		return inspection
	}
	inspection.ReusableGreen, inspection.Reason = true, ""
	return inspection
}

func retainGreen(root string, plan subject, recordedAt time.Time) error {
	gitdir, err := commonGitDir(root)
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
	gitdir, err := commonGitDir(root)
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

func evidencePath(gitdir string, plan subject) string {
	return filepath.Join(gitdir, "bench-gate-evidence", evidenceName(plan))
}

func evidenceName(plan subject) string {
	h := sha256.New()
	frame(h, plan.Tree)
	frame(h, plan.Oracle)
	return hex.EncodeToString(h.Sum(nil))
}

func commonGitDir(root string) (string, error) {
	return benchgit.Output("-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir")
}

func prospectiveCheckout(root, tree string) (string, func(), error) {
	if !treeHashRE.MatchString(tree) {
		return "", nil, errors.New("invalid tree")
	}
	path, err := os.MkdirTemp("", "bench-gate-subject-")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() {
		_ = exec.Command("git", "-C", root, "worktree", "remove", "--force", path).Run()
		_ = os.RemoveAll(path)
	}
	if output, err := exec.Command("git", "-C", root, "worktree", "add", "--quiet", "--detach", path, "HEAD").CombinedOutput(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("create prospective checkout: %s", strings.TrimSpace(string(output)))
	}
	if output, err := exec.Command("git", "-C", path, "read-tree", "--reset", "-u", tree).CombinedOutput(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("materialize prospective tree: %s", strings.TrimSpace(string(output)))
	}
	return path, cleanup, nil
}

func durableReplaceAt(dir, name string, rec verdictRecord) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return durableReplaceRecordAt(dir, name, data)
}

// durableReplaceRecordAt installs data as the store entry named name: written to a private
// temporary beside it, synced, renamed over whatever was there, and the directory synced so
// the rename itself survives a crash. Every record class in the store is published through
// here, so a reader never observes a half-written or world-readable record whatever class
// wrote it.
func durableReplaceRecordAt(dir, name string, data []byte) error {
	data = append(data, '\n')
	tmp, err := os.CreateTemp(dir, ".bench-gate-evidence-")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = tmp.Close(); _ = os.Remove(tmpName) }()
	if err := tmp.Chmod(0o600); err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, filepath.Join(dir, name)); err != nil {
		return err
	}
	file, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}

// reusableEvidence projects the verdict state from the plan the caller already accepted:
// both call sites hold the current subject, so rebuilding one here would be a second
// independent capture standing beside the generation that authorized the plan.
func reusableEvidence(root string, plan subject, now time.Time) Inspection {
	projection := inspectSubjectAt(root, plan, now)
	if projection.State == Pending || projection.State == Invalid || projection.State == Unavailable {
		return Inspection{}
	}
	evidence := inspectEvidence(root, plan, now)
	if !evidence.ReusableGreen {
		return Inspection{}
	}
	return Inspection{State: Ready, Status: "green", CachedTree: plan.Tree, CurrentTree: plan.Tree, RecordedAt: evidence.RecordedAt, ReusableGreen: true}
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

func runResolved(ctx context.Context, root string, res Resolution, env []string, stdout, stderr io.Writer, processGroup bool) processGroupResult {
	cmd := res.command(root)
	if cmd == nil {
		return processGroupResult{Code: 3}
	}
	cmd.Dir, cmd.Stdout, cmd.Stderr, cmd.Env = root, stdout, stderr, append([]string(nil), env...)
	if processGroup {
		return runProcessGroupCommand(ctx, cmd)
	}
	if err := cmd.Run(); err != nil {
		if cmd.ProcessState != nil {
			if code := cmd.ProcessState.ExitCode(); code > 0 {
				return processGroupResult{Code: code}
			}
		}
		return processGroupResult{Code: 1, StartErr: err}
	}
	return processGroupResult{}
}

type controlSafeWriter struct{ io.Writer }

func (w controlSafeWriter) Write(p []byte) (int, error) {
	safe := make([]byte, 0, len(p))
	for _, b := range p {
		if (b >= 0x20 && b != 0x7f) || b == '\n' || b == '\r' || b == '\t' {
			safe = append(safe, b)
		}
	}
	if _, err := w.Writer.Write(safe); err != nil {
		return 0, err
	}
	return len(p), nil
}
