package gate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

var ft78Story4Proofs = []r21ProofCase{
	{id: "R9/pending-mode-and-order", driver: r9Order},
	r9Fault("R9/temporary-create-failure", "temporary-create"), r9Fault("R9/mode-establishment-failure", "mode-establishment"),
	r9Fault("R9/write-failure", "write"), r9Fault("R9/file-sync-failure", "file-sync"),
	r9Fault("R9/file-close-failure", "file-close"), r9Fault("R9/atomic-rename-failure", "atomic-rename"),
	r9Fault("R9/directory-open-failure", "directory-open"), r9Fault("R9/directory-sync-failure", "directory-sync"),
	r9Fault("R9/directory-close-failure", "directory-close"),
	r10Control("R10/durable-green", 0), r10Control("R10/durable-red-exit-23", 23),
	r10Fault("R10/green-write-failure", 0, "write"), r10Fault("R10/red-write-failure", 23, "write"),
	r10Fault("R10/green-file-sync-failure", 0, "file-sync"), r10Fault("R10/red-file-sync-failure", 23, "file-sync"),
	r10Fault("R10/green-atomic-rename-failure", 0, "atomic-rename"), r10Fault("R10/red-atomic-rename-failure", 23, "atomic-rename"),
	r10Fault("R10/green-directory-sync-failure", 0, "directory-sync"), r10Fault("R10/red-directory-sync-failure", 23, "directory-sync"),
	r10Recheck("R10/green-subject-recheck-failure", 0), r10Recheck("R10/red-subject-recheck-failure", 23),
	r11Drift("R11/command-drift", "command"), r11Drift("R11/manifest-drift", "manifest"),
	r11Drift("R11/environment-drift", "environment"), r11Drift("R11/path-drift", "path"),
	r11Drift("R11/tool-drift", "tool"), r11Drift("R11/launcher-drift", "launcher"),
	r11Drift("R11/auto-kind-drift", "auto-kind"),
	{id: "R11/cancellation-interruption", driver: r11Cancellation},
	{id: "R11/cancellation-kills-process-group", driver: r11Cancellation},
	r12Contention("R12/standalone-gate-blocked", "gate"), r12Contention("R12/commit-blocked", "commit"),
	r12Contention("R12/shift-blocked", "shift"), r12Contention("R12/armed-stop-blocked", "stop"),
	{id: "R12/separate-git-directory-concurrent", driver: r12SeparateGitDir},
	r13Pending("R13/live-pid-young", os.Getpid(), time.Second), r13Pending("R13/live-pid-old", os.Getpid(), 24*time.Hour),
	r13Pending("R13/dead-pid-young", 99999999, time.Second), r13Pending("R13/dead-pid-old", 99999999, 24*time.Hour),
	{id: "R13/lock-free-pending-interrupted", driver: r13LockFree},
	{id: "R13/killed-owner-interrupted", driver: r13KilledOwner},
	{id: "R13/next-run-replaces-and-finishes", driver: r13Recovery},
}

var ft78Story4ExpectedIDs = []string{
	"R9/pending-mode-and-order", "R9/temporary-create-failure", "R9/mode-establishment-failure", "R9/write-failure", "R9/file-sync-failure", "R9/file-close-failure", "R9/atomic-rename-failure", "R9/directory-open-failure", "R9/directory-sync-failure", "R9/directory-close-failure",
	"R10/durable-green", "R10/durable-red-exit-23", "R10/green-write-failure", "R10/red-write-failure", "R10/green-file-sync-failure", "R10/red-file-sync-failure", "R10/green-atomic-rename-failure", "R10/red-atomic-rename-failure", "R10/green-directory-sync-failure", "R10/red-directory-sync-failure", "R10/green-subject-recheck-failure", "R10/red-subject-recheck-failure",
	"R11/command-drift", "R11/manifest-drift", "R11/environment-drift", "R11/path-drift", "R11/tool-drift", "R11/launcher-drift", "R11/auto-kind-drift", "R11/cancellation-interruption", "R11/cancellation-kills-process-group",
	"R12/standalone-gate-blocked", "R12/commit-blocked", "R12/shift-blocked", "R12/armed-stop-blocked", "R12/separate-git-directory-concurrent",
	"R13/live-pid-young", "R13/live-pid-old", "R13/dead-pid-young", "R13/dead-pid-old", "R13/lock-free-pending-interrupted", "R13/killed-owner-interrupted", "R13/next-run-replaces-and-finishes",
}

func TestFT78Story4ProofLedgerCompleteness(t *testing.T) {
	contract.NoteContractFailure(t, r21CompletenessFailure)
	seen := map[string]int{}
	for _, proof := range ft78Story4Proofs {
		seen[proof.id]++
		if proof.driver == nil {
			t.Fatalf("%s: nil real driver", proof.id)
		}
	}
	if len(seen) != len(ft78Story4ExpectedIDs) {
		t.Fatalf("registered IDs = %d, want %d", len(seen), len(ft78Story4ExpectedIDs))
	}
	for _, id := range ft78Story4ExpectedIDs {
		if seen[id] != 1 {
			t.Fatalf("%s: %s registrations = %d, want 1", r21CompletenessFailure, id, seen[id])
		}
	}
}

func TestFT78Story4ProofLedger(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	for _, proof := range ft78Story4Proofs {
		proof := proof
		t.Run(proof.id, proof.driver)
	}
}

var pendingTrace = []string{"lock-open", "lock-acquisition", "temporary-create", "mode-establishment", "write", "file-sync", "file-close", "atomic-rename", "directory-open", "directory-sync", "directory-close"}

func story4Repo(t *testing.T, exit int) string {
	return gateTestRepo(t, fmt.Sprintf("#!/usr/bin/env bash\ntouch .git/gate-marker\nexit %d\n", exit), `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
}

func cachePath(t *testing.T, root string) string {
	t.Helper()
	dir, err := benchgit.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(dir, benchgit.GateCacheFile)
}

func requirePending(t *testing.T, root string, got Result, gateExit int, marker bool, started time.Time, plan subject) {
	t.Helper()
	_, markerErr := os.Stat(filepath.Join(root, ".git", "gate-marker"))
	if (markerErr == nil) != marker {
		t.Fatalf("gate marker present = %v, want %v", markerErr == nil, marker)
	}
	wantBytes := []byte(fmt.Sprintf(`{"schema":1,"state":"pending","tree":%q,"oracle":%q,"started_at":%q,"owner_pid":%d}`+"\n", plan.Tree, plan.Oracle, started.UTC().Truncate(time.Second).Format(time.RFC3339), os.Getpid()))
	path := cachePath(t, root)
	data, readErr := os.ReadFile(path)
	info, statErr := os.Lstat(path)
	if readErr != nil || statErr != nil || !bytes.Equal(data, wantBytes) || !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 || info.Size() != int64(len(wantBytes)) {
		t.Fatalf("durable pending = %q/read=%v/info=%v/stat=%v, want exact %q regular 0600", data, readErr, info, statErr, wantBytes)
	}
	current := mustSubject(t, root)
	wantInspection := Inspection{State: Pending, PendingStatus: "locked-pending", CachedTree: plan.Tree, CurrentTree: current.Tree, Reason: current.Reason, CacheBytes: len(wantBytes)}
	if got.GateExit != gateExit || got.ActionExit != 1 || !reflect.DeepEqual(got.Inspection, wantInspection) {
		t.Fatalf("result = %+v, want gate %d/action 1/inspection %+v", got, gateExit, wantInspection)
	}
}

func r9Order(t *testing.T) {
	root := story4Repo(t, 0)
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	plan := mustSubject(t, root)
	engine := &faultEngine{now: now, failOp: "post-run-subject-rebuild"}
	got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
	want := append(append([]string{}, pendingTrace...), "post-run-subject-rebuild")
	if !reflect.DeepEqual(engine.trace, want) {
		t.Fatalf("trace = %v, want %v", engine.trace, want)
	}
	requirePending(t, root, got, 0, true, now, plan)
}

func r9Fault(id, op string) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := story4Repo(t, 0)
		now := time.Now().UTC()
		engine := &faultEngine{now: now, failOp: op}
		got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
		want := append(append([]string{}, pendingTrace[:2]...), failedPersistenceTrace(op)...)
		want = append(want, pendingTrace[2:]...)
		if !reflect.DeepEqual(engine.trace, want) {
			t.Fatalf("trace = %v, want %v", engine.trace, want)
		}
		if got.ActionExit != 1 || got.GateExit != 0 {
			t.Fatalf("result = %+v, want 0/1", got)
		}
		if _, err := os.Stat(filepath.Join(root, ".git", "gate-marker")); !os.IsNotExist(err) {
			t.Fatalf("pre-run fault ran gate: %v", err)
		}
		data, err := os.ReadFile(cachePath(t, root))
		if err != nil || got.Inspection.State != Pending || !bytes.Contains(data, []byte(`"state":"pending"`)) {
			t.Fatalf("durable floor = %q/%v/%+v, want pending", data, err, got.Inspection)
		}
	}}
}

func failedPersistenceTrace(op string) []string {
	persistence := pendingTrace[2:]
	idx := 0
	for i, name := range persistence {
		if name == op {
			idx = i
			break
		}
	}
	want := append([]string{}, persistence[:idx+1]...)
	if op == "mode-establishment" || op == "write" || op == "file-sync" {
		want = append(want, "file-close")
	}
	if op == "directory-sync" {
		want = append(want, "directory-close")
	}
	return want
}

func r10Control(id string, exit int) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := story4Repo(t, exit)
		got := Execute(context.Background(), root, io.Discard, io.Discard)
		wantStatus := "green"
		if exit != 0 {
			wantStatus = "red"
		}
		if got.GateExit != exit || got.ActionExit != exit || got.Inspection.State != Ready || got.Inspection.Status != wantStatus {
			t.Fatalf("result = %+v, want %d/%d ready %s", got, exit, exit, wantStatus)
		}
		data, err := os.ReadFile(cachePath(t, root))
		if err != nil || !bytes.Contains(data, []byte(`"status":"`+wantStatus+`"`)) {
			t.Fatalf("durable bytes = %q/%v", data, err)
		}
	}}
}

func r10Fault(id string, exit int, op string) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := story4Repo(t, exit)
		now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
		plan := mustSubject(t, root)
		engine := &faultEngine{now: now, failOp: op, failAt: 2}
		got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
		want := append(append([]string{}, pendingTrace...), "post-run-subject-rebuild")
		want = append(want, failedPersistenceTrace(op)...)
		want = append(want, pendingTrace[2:]...)
		if !reflect.DeepEqual(engine.trace, want) {
			t.Fatalf("trace = %v, want %v", engine.trace, want)
		}
		requirePending(t, root, got, exit, true, now, plan)
	}}
}

type recheckFaultEngine struct {
	faultEngine
	failRecheck bool
}

func (e *recheckFaultEngine) PostRunSubject(root string) (subject, error) {
	e.trace = append(e.trace, "post-run-subject-rebuild")
	if e.failRecheck {
		return subject{}, fmt.Errorf("scripted recheck")
	}
	return buildSubject(root)
}
func r10Recheck(id string, exit int) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := story4Repo(t, exit)
		now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
		plan := mustSubject(t, root)
		engine := &recheckFaultEngine{faultEngine: faultEngine{now: now}, failRecheck: true}
		got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine)
		want := append(append([]string{}, pendingTrace...), "post-run-subject-rebuild")
		if !reflect.DeepEqual(engine.trace, want) {
			t.Fatalf("trace = %v, want %v", engine.trace, want)
		}
		requirePending(t, root, got, exit, true, now, plan)
	}}
}

func r11Drift(id, kind string) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := story4Repo(t, 0)
		manifest := filepath.Join(root, ".bench", "gate-inputs.json")
		const blockingGate = "#!/usr/bin/env bash\ngitdir=\"$(git rev-parse --absolute-git-dir)\"\ntouch \"$gitdir/gate-marker\"\nwhile [ ! -f \"$gitdir/release-gate\" ]; do sleep .01; done\nexit 0\n"
		var mutate func()
		switch kind {
		case "command":
			mutate = func() {
				_ = os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte(blockingGate+"# drift\n"), 0o755)
			}
		case "manifest":
			mutate = func() {
				_ = os.WriteFile(manifest, []byte(`{"schema":1,"closure":"remote","environment":[],"paths":[],"tools":[]}`), 0o644)
			}
		case "environment":
			t.Setenv("FT78_DRIFT", "before")
			_ = os.WriteFile(manifest, []byte(`{"schema":1,"closure":"local","environment":["FT78_DRIFT"],"paths":[],"tools":[]}`), 0o644)
			mutate = func() { _ = os.Setenv("FT78_DRIFT", "after") }
		case "path":
			_ = os.MkdirAll(filepath.Join(root, "inputs"), 0o755)
			_ = os.WriteFile(filepath.Join(root, "inputs", "x"), []byte("before"), 0o644)
			_ = os.WriteFile(manifest, []byte(`{"schema":1,"closure":"local","environment":[],"paths":["inputs/x"],"tools":[]}`), 0o644)
			mutate = func() { _ = os.WriteFile(filepath.Join(root, "inputs", "x"), []byte("after"), 0o644) }
		case "tool":
			tool := filepath.Join(root, "ft78-tool")
			_ = os.WriteFile(tool, []byte("#!/bin/sh\nexit 0\n"), 0o755)
			t.Setenv("PATH", root+string(os.PathListSeparator)+os.Getenv("PATH"))
			_ = os.WriteFile(manifest, []byte(`{"schema":1,"closure":"local","environment":[],"paths":[],"tools":["ft78-tool"]}`), 0o644)
			mutate = func() { _ = os.WriteFile(tool, []byte("#!/bin/sh\nexit 0\n# drift\n"), 0o755) }
		case "launcher":
			mutate = func() { _ = os.Chmod(filepath.Join(root, ".bench", "gate.sh"), 0o700) }
		case "auto-kind":
			_ = os.Remove(filepath.Join(root, ".bench", "gate.sh"))
			_ = os.Remove(manifest)
			_ = os.WriteFile(filepath.Join(root, "package.json"), []byte("{}\n"), 0o644)
			_ = os.WriteFile(filepath.Join(root, "npm"), []byte(blockingGate), 0o755)
			t.Setenv("PATH", root+string(os.PathListSeparator)+os.Getenv("PATH"))
			mutate = func() { _ = os.WriteFile(filepath.Join(root, "pnpm-lock.yaml"), []byte("lock\n"), 0o644) }
		}
		if kind != "auto-kind" {
			_ = os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte(blockingGate), 0o755)
		}
		plan := mustSubject(t, root)
		gitdir := filepath.Dir(cachePath(t, root))
		done := make(chan Result, 1)
		go func() { done <- Execute(context.Background(), root, io.Discard, io.Discard) }()
		defer func() { _ = os.WriteFile(filepath.Join(gitdir, "release-gate"), nil, 0o600) }()
		waitFile(t, filepath.Join(gitdir, "gate-marker"))
		pendingBytes := mustRead(t, cachePath(t, root))
		var pending verdictRecord
		if err := strictJSON(pendingBytes, &pending); err != nil {
			t.Fatal(err)
		}
		started, err := time.Parse(time.RFC3339, pending.StartedAt)
		if err != nil {
			t.Fatal(err)
		}
		mutate()
		if err := os.WriteFile(filepath.Join(gitdir, "release-gate"), nil, 0o600); err != nil {
			t.Fatal(err)
		}
		select {
		case got := <-done:
			requirePending(t, root, got, 0, true, started, plan)
		case <-time.After(5 * time.Second):
			t.Fatal("drifted gate did not return")
		}
	}}
}

func r11Cancellation(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\nsleep 30 & echo $! > .git/child-pid; wait\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan Result, 1)
	go func() { done <- Execute(ctx, root, io.Discard, io.Discard) }()
	waitFile(t, filepath.Join(root, ".git", "child-pid"))
	cancel()
	got := <-done
	if got.GateExit != 130 || got.ActionExit != 130 || got.Inspection.State != Pending {
		t.Fatalf("cancellation result = %+v, want 130/130 pending", got)
	}
	var pid int
	_, _ = fmt.Sscanf(strings.TrimSpace(string(mustRead(t, filepath.Join(root, ".git", "child-pid")))), "%d", &pid)
	if syscall.Kill(pid, 0) == nil {
		t.Fatalf("gate child %d survived cancellation", pid)
	}
}

func r12SeparateGitDir(t *testing.T) {
	root1 := story4Repo(t, 0)
	root2 := story4Repo(t, 0)
	lock, err := os.OpenFile(filepath.Join(root1, ".git", "bench-gate.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	held := recordLock(syscall.F_WRLCK)
	if err := syscall.FcntlFlock(lock.Fd(), syscall.F_SETLK, &held); err != nil {
		t.Fatal(err)
	}
	got := Execute(context.Background(), root2, io.Discard, io.Discard)
	if got.ActionExit != 0 || got.Inspection.State != Ready {
		t.Fatalf("separate Git dir result = %+v", got)
	}
}

func r13Pending(id string, pid int, age time.Duration) r21ProofCase {
	return r21ProofCase{id: id, driver: func(t *testing.T) {
		root := story4Repo(t, 0)
		plan, _ := buildSubject(root)
		rec := verdictRecord{Schema: 1, State: Pending, Tree: plan.Tree, Oracle: plan.Oracle, StartedAt: time.Now().UTC().Add(-age).Truncate(time.Second).Format(time.RFC3339), OwnerPID: pid}
		if err := durableReplace(filepath.Dir(cachePath(t, root)), rec); err != nil {
			t.Fatal(err)
		}
		got := Inspect(root)
		if got.State != Pending || got.PendingStatus != "interrupted-pending" || got.ReusableGreen {
			t.Fatalf("PID/age inspection = %+v", got)
		}
	}}
}
func r13LockFree(t *testing.T) { r13Pending("", os.Getpid(), time.Second).driver(t) }

func r13KilledOwner(t *testing.T) {
	root := gateTestRepo(t, "#!/usr/bin/env bash\necho $$ > .git/gate-pid\ntouch .git/live-owner\nsleep 30\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	binary := filepath.Join(contract.SubjectRoot(t), "dist", "bench")
	cmd := exec.Command(binary, "gate-run", root)
	cmd.Dir = root
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	waitFile(t, filepath.Join(root, ".git", "live-owner"))
	var gatePID int
	_, _ = fmt.Sscanf(strings.TrimSpace(string(mustRead(t, filepath.Join(root, ".git", "gate-pid")))), "%d", &gatePID)
	defer func() { _ = syscall.Kill(-gatePID, syscall.SIGKILL) }()
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	_ = cmd.Wait()
	got := Inspect(root)
	if got.State != Pending || got.PendingStatus != "interrupted-pending" {
		t.Fatalf("killed-owner inspection = %+v", got)
	}
	_ = syscall.Kill(-gatePID, syscall.SIGKILL)
	waitForProcessExit(t, gatePID)
}
func r13Recovery(t *testing.T) {
	root := story4Repo(t, 0)
	plan, _ := buildSubject(root)
	rec := verdictRecord{Schema: 1, State: Pending, Tree: plan.Tree, Oracle: plan.Oracle, StartedAt: time.Now().UTC().Add(-time.Hour).Truncate(time.Second).Format(time.RFC3339), OwnerPID: 99999999}
	if err := durableReplace(filepath.Dir(cachePath(t, root)), rec); err != nil {
		t.Fatal(err)
	}
	got := Execute(context.Background(), root, io.Discard, io.Discard)
	if got.ActionExit != 0 || got.Inspection.State != Ready || !got.Inspection.ReusableGreen {
		t.Fatalf("recovery result = %+v", got)
	}
}

func waitFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", path)
}
func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
