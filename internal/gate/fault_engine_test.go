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
	"sort"
	"testing"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

const r21CompletenessFailure = "FT78 proof ledger completeness contract failed"

type r21ProofCase struct {
	id     string
	driver func(*testing.T)
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
	"R21/lock-open",
	"R21/lock-acquisition",
	"R21/temporary-create",
	"R21/mode-establishment",
	"R21/write",
	"R21/file-sync",
	"R21/file-close",
	"R21/atomic-rename",
	"R21/directory-open",
	"R21/directory-sync",
	"R21/directory-close",
	"R21/post-run-subject-rebuild",
	"R21/byte-bound-16384",
	"R21/byte-bound-16385",
	"R21/future-clock",
}

type faultEngine struct {
	productionGateEngine
	now    time.Time
	failOp string
	failed bool
	trace  []string
}

func (e *faultEngine) Now() time.Time { return e.now }

func (e *faultEngine) operation(name string) error {
	e.trace = append(e.trace, name)
	if name == e.failOp && !e.failed {
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

func TestR21DeterministicFaultProofRegistryCompleteness(t *testing.T) {
	got := make([]string, 0, len(r21ProofRegistry))
	seen := map[string]bool{}
	for _, proof := range r21ProofRegistry {
		if proof.id == "" || proof.driver == nil || seen[proof.id] {
			t.Fatalf("%s: invalid or duplicate registration %q", r21CompletenessFailure, proof.id)
		}
		seen[proof.id] = true
		got = append(got, proof.id)
	}
	want := append([]string(nil), r21ExpectedProofIDs...)
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: got IDs %v, want %v", r21CompletenessFailure, got, want)
	}
}

func TestR21DeterministicFaultEngine(t *testing.T) {
	for _, proof := range r21ProofRegistry {
		proof := proof
		t.Run(proof.id, proof.driver)
	}
}

func r21FaultCase(id, failOp string, wantTrace []string, durablePending bool) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
		now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
		engine := &faultEngine{now: now, failOp: failOp}
		got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
		if !reflect.DeepEqual(engine.trace, wantTrace) {
			t.Fatalf("trace = %v, want %v", engine.trace, wantTrace)
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
		if !durablePending {
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
		got := inspectAt(root, time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC))
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
	got := inspectAt(root, now)
	if got.State != Invalid || got.Reason != "invalid cache record" || got.ReusableGreen {
		t.Fatalf("future-clock inspection = %+v, want invalid non-reusable record", got)
	}
	data, err := os.ReadFile(filepath.Join(gitdir, benchgit.GateCacheFile))
	if err != nil || !bytes.Contains(data, []byte(fmt.Sprintf(`"recorded_at":%q`, rec.RecordedAt))) {
		t.Fatalf("future-clock durable bytes = %q/error %v", data, err)
	}
}
