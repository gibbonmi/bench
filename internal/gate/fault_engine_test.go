package gate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

const r21CompletenessFailure = "FT78 proof ledger completeness contract failed"

func story3Fixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/bin/sh\necho run >> .git/ft78-runs\nexit 0\n")
	f.WriteFile(".gitignore", ".bench/gate.sh\n.bench/gate-inputs.json\nft78-*\ninputs/\ntools/\n")
	f.CommitAll("base")
	return f
}

func mustSubject(t *testing.T, root string) subject {
	t.Helper()
	s, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

type story3Verdict struct {
	Oracle string `json:"oracle"`
}

func story3GitDir(f contract.Fixture) string { return filepath.Join(f.Root, ".git") }

func story3ReadVerdict(t *testing.T, f contract.Fixture) story3Verdict {
	t.Helper()
	var rec story3Verdict
	if err := json.Unmarshal([]byte(contract.ReadFileAbs(t, filepath.Join(story3GitDir(f), "bench-last-gate"))), &rec); err != nil {
		t.Fatal(err)
	}
	return rec
}

func writeCache(t *testing.T, path, body string, mode os.FileMode) {
	t.Helper()
	if err := os.RemoveAll(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), mode); err != nil {
		t.Fatal(err)
	}
}

type r21ProofCase struct {
	id     string
	driver func(*testing.T)
}

var ft78Story3Proofs = []r21ProofCase{
	codecProof("R5/cache-absent", codecAbsent), codecProof("R5/ready-red", codecReadyRed),
	codecProof("R5/subject-unavailable", codecUnavailable), codecProof("R5/zero-byte", codecZeroByte),
	codecProof("R5/ready-green-no-final-newline", codecNoFinalNewlineGreen), codecProof("R5/ready-red-no-final-newline", codecNoFinalNewlineRed),
	codecProof("R5/trailing-json", codecTrailing), codecProof("R5/duplicate-field", codecDuplicate),
	codecProof("R5/unknown-field", codecUnknown), codecProof("R5/wrong-field-type", codecWrongType),
	codecProof("R5/wrong-schema", codecWrongSchema), codecProof("R5/wrong-state-enum", codecWrongState),
	codecProof("R5/wrong-status-enum", codecWrongStatus), codecProof("R5/wrong-hash", codecWrongHash),
	codecProof("R5/wrong-time", codecWrongTime), codecProof("R5/legacy-record", codecLegacy),
	codecProof("R5/truncated-record", codecTruncated), codecProof("R5/byte-bound-16384", codec16384),
	codecProof("R5/byte-bound-16385", codec16385), codecProof("R5/cache-symlink", codecSymlink),
	codecProof("R5/cache-directory", codecDirectory), codecProof("R5/cache-unreadable", codecUnreadable),
	freshnessProof("R6/after-ten-minutes", freshnessAfter), freshnessProof("R6/exact-ten-minutes", freshnessExact),
	freshnessProof("R6/future-record", freshnessFuture), freshnessProof("R6/malformed-time", freshnessMalformed),
	freshnessProof("R6/fingerprint-mismatch", freshnessFingerprint), freshnessProof("R6/freshness-policy-mismatch", freshnessPolicy),
	secretProof("R7/sentinel-command", secretCommand), secretProof("R7/sentinel-environment-name", secretEnvironmentName),
	secretProof("R7/sentinel-environment-value", secretEnvironmentValue), secretProof("R7/sentinel-manifest-path", secretManifestPath),
	secretProof("R7/sentinel-input-content", secretInputContent), secretProof("R7/sentinel-tool-output", secretToolOutput),
	secretProof("R7/sentinel-gate-output", secretGateOutput), secretProof("R7/unsafe-control-bytes", secretControlBytes),
	hostileProof("R8/repository-path-spaces", hostileRepositorySpaces), hostileProof("R8/repository-path-glob", hostileRepositoryGlob),
	hostileProof("R8/declared-path-spaces", hostileDeclaredSpaces), hostileProof("R8/declared-path-glob", hostileDeclaredGlob),
	hostileProof("R8/manifest-no-final-newline", hostileManifestNoNewline), hostileProof("R8/symlink-chain", hostileSymlinkChain),
	hostileProof("R8/external-symlink-target", hostileExternalTarget), hostileProof("R8/missing-global-bench", hostileMissingGlobalBench),
	hostileProof("R8/executable-mode-change", hostileExecutableMode), hostileProof("R8/control-byte-safe-output", hostileControlOutput),
}

var ft78Story3ExpectedIDs = []string{
	"R5/cache-absent", "R5/ready-red", "R5/subject-unavailable", "R5/zero-byte", "R5/ready-green-no-final-newline", "R5/ready-red-no-final-newline", "R5/trailing-json", "R5/duplicate-field", "R5/unknown-field", "R5/wrong-field-type", "R5/wrong-schema", "R5/wrong-state-enum", "R5/wrong-status-enum", "R5/wrong-hash", "R5/wrong-time", "R5/legacy-record", "R5/truncated-record", "R5/byte-bound-16384", "R5/byte-bound-16385", "R5/cache-symlink", "R5/cache-directory", "R5/cache-unreadable",
	"R6/after-ten-minutes", "R6/exact-ten-minutes", "R6/future-record", "R6/malformed-time", "R6/fingerprint-mismatch", "R6/freshness-policy-mismatch",
	"R7/sentinel-command", "R7/sentinel-environment-name", "R7/sentinel-environment-value", "R7/sentinel-manifest-path", "R7/sentinel-input-content", "R7/sentinel-tool-output", "R7/sentinel-gate-output", "R7/unsafe-control-bytes",
	"R8/repository-path-spaces", "R8/repository-path-glob", "R8/declared-path-spaces", "R8/declared-path-glob", "R8/manifest-no-final-newline", "R8/symlink-chain", "R8/external-symlink-target", "R8/missing-global-bench", "R8/executable-mode-change", "R8/control-byte-safe-output",
}

func TestFT78Story3ProofLedgerCompleteness(t *testing.T) {
	contract.NoteContractFailure(t, r21CompletenessFailure)
	seen := map[string]int{}
	for _, proof := range ft78Story3Proofs {
		seen[proof.id]++
		if proof.driver == nil {
			t.Fatalf("%s: nil real driver", proof.id)
		}
	}
	if len(seen) != len(ft78Story3ExpectedIDs) {
		t.Fatalf("registered IDs = %d, want %d", len(seen), len(ft78Story3ExpectedIDs))
	}
	for _, id := range ft78Story3ExpectedIDs {
		if seen[id] != 1 {
			t.Fatalf("%s registrations = %d, want 1", id, seen[id])
		}
	}
}

func TestFT78Story3ProofLedger(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	for _, proof := range ft78Story3Proofs {
		proof := proof
		t.Run(proof.id, proof.driver)
	}
}

var r21ProofRegistry = []r21ProofCase{
	r21FaultCase("R21/lock-open", "lock-open", []string{"lock-open"}, false),
	r21FaultCase("R21/lock-acquisition", "lock-acquisition", []string{"lock-open", "lock-acquisition"}, false),
	r21FaultCase("R21/temporary-create", "temporary-create", []string{"lock-open", "lock-acquisition", "temporary-create"}, false),
	r21FaultCase("R21/mode-establishment", "mode-establishment", []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "file-close"}, false),
	r21FaultCase("R21/write", "write", []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "write", "file-close"}, false),
	r21FaultCase("R21/file-sync", "file-sync", []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "write", "file-sync", "file-close"}, false),
	r21FaultCase("R21/file-close", "file-close", []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "write", "file-sync", "file-close"}, false),
	r21FaultCase("R21/atomic-rename", "atomic-rename", []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "write", "file-sync", "file-close", "atomic-rename"}, false),
	r21FaultCase("R21/directory-open", "directory-open", []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "write", "file-sync", "file-close", "atomic-rename", "directory-open"}, true),
	r21FaultCase("R21/directory-sync", "directory-sync", []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "write", "file-sync", "file-close", "atomic-rename", "directory-open", "directory-sync", "directory-close"}, true),
	r21FaultCase("R21/directory-close", "directory-close", []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "write", "file-sync", "file-close", "atomic-rename", "directory-open", "directory-sync", "directory-close"}, true),
	r21FaultCase("R21/post-run-subject-rebuild", "post-run-subject-rebuild", []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "write", "file-sync", "file-close", "atomic-rename", "directory-open", "directory-sync", "directory-close", "post-run-subject-rebuild"}, true),
	{id: "R21/byte-bound-16384", driver: r21ByteBoundDriver(16_384, "invalid cache framing")},
	{id: "R21/byte-bound-16385", driver: r21ByteBoundDriver(16_385, "invalid cache record")},
	{id: "R21/future-clock", driver: r21FutureClockDriver},
}

var r21ExpectedProofIDs = []string{
	"R21/lock-open", "R21/lock-acquisition", "R21/temporary-create", "R21/mode-establishment", "R21/write",
	"R21/file-sync", "R21/file-close", "R21/atomic-rename", "R21/directory-open", "R21/directory-sync",
	"R21/directory-close", "R21/post-run-subject-rebuild", "R21/byte-bound-16384", "R21/byte-bound-16385", "R21/future-clock",
}

type faultEngine struct {
	productionGateEngine
	now      time.Time
	failOp   string
	failAt   int
	opCounts map[string]int
	failed   bool
	trace    []string
}

func (e *faultEngine) Now() time.Time { return e.now }

func (e *faultEngine) operation(name string) error {
	e.trace = append(e.trace, name)
	if e.opCounts == nil {
		e.opCounts = map[string]int{}
	}
	e.opCounts[name]++
	failAt := e.failAt
	if failAt == 0 {
		failAt = 1
	}
	if name == e.failOp && e.opCounts[name] == failAt && !e.failed {
		e.failed = true
		return errors.New("scripted " + name + " failure")
	}
	return nil
}

func (e *faultEngine) OpenLock(path string) (gateFile, error) {
	if err := e.operation("lock-open"); err != nil {
		return nil, err
	}
	return e.productionGateEngine.OpenLock(path)
}

func (e *faultEngine) Acquire(f gateFile) error {
	if err := e.operation("lock-acquisition"); err != nil {
		return err
	}
	return e.productionGateEngine.Acquire(f)
}

func (e *faultEngine) CreateTemp(dir, pattern string) (gateFile, error) {
	if err := e.operation("temporary-create"); err != nil {
		return nil, err
	}
	f, err := e.productionGateEngine.CreateTemp(dir, pattern)
	if err != nil {
		return nil, err
	}
	return &faultFile{gateFile: f, engine: e, kind: "file"}, nil
}

func (e *faultEngine) Rename(oldpath, newpath string) error {
	if err := e.operation("atomic-rename"); err != nil {
		return err
	}
	return e.productionGateEngine.Rename(oldpath, newpath)
}

func (e *faultEngine) OpenDir(path string) (gateFile, error) {
	if err := e.operation("directory-open"); err != nil {
		return nil, err
	}
	f, err := e.productionGateEngine.OpenDir(path)
	if err != nil {
		return nil, err
	}
	return &faultFile{gateFile: f, engine: e, kind: "directory"}, nil
}

func (e *faultEngine) WriteFile(path string, data []byte, mode os.FileMode) error {
	if err := e.operation("owner-write"); err != nil {
		return err
	}
	return e.productionGateEngine.WriteFile(path, data, mode)
}

func (e *faultEngine) PostRunSubject(root string) (subject, error) {
	if err := e.operation("post-run-subject-rebuild"); err != nil {
		return subject{}, err
	}
	return e.productionGateEngine.PostRunSubject(root)
}

type faultFile struct {
	gateFile
	engine    *faultEngine
	kind      string
	closeSeen bool
}

func (f *faultFile) Chmod(mode os.FileMode) error {
	if err := f.engine.operation("mode-establishment"); err != nil {
		return err
	}
	return f.gateFile.Chmod(mode)
}

func (f *faultFile) Write(data []byte) (int, error) {
	if err := f.engine.operation("write"); err != nil {
		return 0, err
	}
	return f.gateFile.Write(data)
}

func (f *faultFile) Sync() error {
	op := "file-sync"
	if f.kind == "directory" {
		op = "directory-sync"
	}
	if err := f.engine.operation(op); err != nil {
		return err
	}
	return f.gateFile.Sync()
}

func (f *faultFile) Close() error {
	if f.closeSeen {
		return f.gateFile.Close()
	}
	f.closeSeen = true
	op := "file-close"
	if f.kind == "directory" {
		op = "directory-close"
	}
	if err := f.engine.operation(op); err != nil {
		return err
	}
	return f.gateFile.Close()
}

func r21FaultCase(id, failOp string, wantTrace []string, durablePending bool) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
		now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
		engine := &faultEngine{now: now, failOp: failOp}
		got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
		expectedTrace, expectPending := wantTrace, durablePending
		if failOp != "lock-open" && failOp != "lock-acquisition" {
			expectedTrace = append(append([]string{}, wantTrace[:2]...), append([]string{"owner-write"}, wantTrace[2:]...)...)
		}
		if failOp != "lock-open" && failOp != "lock-acquisition" && failOp != "post-run-subject-rebuild" {
			expectedTrace = append(expectedTrace, pendingTrace[3:]...)
			expectPending = true
		}
		if !reflect.DeepEqual(engine.trace, expectedTrace) {
			t.Fatalf("trace = %v, want %v", engine.trace, expectedTrace)
		}
		if got.GateExit != 0 || got.ActionExit != 1 {
			t.Fatalf("exits = gate %d/action %d, want 0/1", got.GateExit, got.ActionExit)
		}
		gitdir, err := benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
		if err != nil {
			t.Fatal(err)
		}
		cache := filepath.Join(gitdir, benchgit.GateCacheFile)
		data, readErr := os.ReadFile(cache)
		if !expectPending {
			if !errors.Is(readErr, os.ErrNotExist) || got.Inspection.State != Absent {
				t.Fatalf("durable floor = bytes %q/error %v/inspection %+v, want absent", data, readErr, got.Inspection)
			}
			return
		}
		plan, err := buildSubject(root)
		if err != nil {
			t.Fatal(err)
		}
		wantRecord, err := json.Marshal(verdictRecord{Schema: 1, State: Pending, Tree: plan.Tree, Oracle: plan.Oracle, StartedAt: now.Format(time.RFC3339), OwnerPID: os.Getpid()})
		if err != nil {
			t.Fatal(err)
		}
		wantRecord = append(wantRecord, '\n')
		if readErr != nil || !bytes.Equal(data, wantRecord) {
			t.Fatalf("durable bytes = %q/error %v, want %q", data, readErr, wantRecord)
		}
		if got.Inspection.State != Pending || got.Inspection.PendingStatus != "locked-pending" || got.Inspection.CacheBytes != len(wantRecord) {
			t.Fatalf("inspection = %+v, want locked pending with %d bytes", got.Inspection, len(wantRecord))
		}
	}}
}

func r21ByteBoundDriver(size int, wantReason string) func(*testing.T) {
	return func(t *testing.T) {
		root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
		gitdir, err := benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
		if err != nil {
			t.Fatal(err)
		}
		data := bytes.Repeat([]byte{'x'}, size)
		data[0] = '{'
		if err := os.WriteFile(filepath.Join(gitdir, benchgit.GateCacheFile), data, 0o600); err != nil {
			t.Fatal(err)
		}
		got := inspectWithEngine(root, &faultEngine{now: time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)})
		if got.State != Invalid || got.Reason != wantReason || got.CacheBytes != size {
			t.Fatalf("inspection = %+v, want invalid/%q/%d bytes", got, wantReason, size)
		}
		persisted, err := os.ReadFile(filepath.Join(gitdir, benchgit.GateCacheFile))
		if err != nil || !bytes.Equal(persisted, data) {
			t.Fatalf("inspection mutated durable bytes: equal=%v err=%v", bytes.Equal(persisted, data), err)
		}
	}
}

func r21FutureClockDriver(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	plan, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	gitdir, err := benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	rec := verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: now.Add(time.Second).Format(time.RFC3339)}
	if err := durableReplace(gitdir, rec); err != nil {
		t.Fatal(err)
	}
	got := inspectWithEngine(root, &faultEngine{now: now})
	if got.State != Invalid || got.Reason != "invalid cache record" || got.ReusableGreen {
		t.Fatalf("future-clock inspection = %+v, want invalid non-reusable record", got)
	}
	data, err := os.ReadFile(filepath.Join(gitdir, benchgit.GateCacheFile))
	if err != nil || !bytes.Contains(data, []byte(fmt.Sprintf(`"recorded_at":%q`, rec.RecordedAt))) {
		t.Fatalf("future-clock durable bytes = %q/error %v", data, err)
	}
}
