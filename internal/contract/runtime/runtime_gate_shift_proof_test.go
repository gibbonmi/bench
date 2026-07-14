package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/intent"
)

type shiftProofRun struct {
	f       contract.Fixture
	home    string
	pool    runtimePoolWorktrees
	base    string
	variant string
}

func proveShiftResult(t *testing.T, variant string) {
	run := newShiftProofRun(t, variant)
	var held *os.File
	if variant == "lock" {
		held = holdShiftGateLock(t, run)
		defer held.Close()
	}
	if variant == "persistence" {
		cache := filepath.Join(shiftProofGitDir(t, run), "bench-last-gate")
		if err := os.Mkdir(cache, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if variant == "cancellation" {
		proveCancelledShift(t, run)
		return
	}

	probe := run.f.BenchEnv(run.env(), "shift", "R15-"+variant)
	if probe.ExitCode != 1 {
		t.Fatalf("%s shift exit = %d, want 1\n%s%s", variant, probe.ExitCode, probe.Stdout, probe.Stderr)
	}
	want := map[string]string{
		"red":         "gate failed — snapshotting iteration 1",
		"lock":        "gate execution already in progress",
		"persistence": "gate pending persistence failed",
		"drift":       "gate subject changed during execution",
	}[variant]
	if !strings.Contains(probe.Stdout+probe.Stderr, want) {
		t.Fatalf("%s shift result missing %q:\n%s%s", variant, want, probe.Stdout, probe.Stderr)
	}
	requirePreservedShift(t, run, probe.Stdout)
}

func newShiftProofRun(t *testing.T, variant string) shiftProofRun {
	t.Helper()
	gate := "#!/usr/bin/env bash\nexit 23\n"
	if variant == "drift" {
		gate = `#!/usr/bin/env bash
printf drift >> tracked.txt
git status --porcelain=v1 --untracked-files=all > "$BENCH_TEST_STATE/status"
git diff --binary HEAD -- > "$BENCH_TEST_STATE/diff"
exit 0
`
	}
	if variant == "cancellation" {
		gate = `#!/usr/bin/env bash
echo $$ > "$BENCH_TEST_STATE/gate-pgid"
: > "$BENCH_TEST_STATE/gate-started"
while :; do sleep .05; done
`
	}
	f := shiftFixture(t, gate)
	f.WriteExecutable("agent", `#!/usr/bin/env bash
mkdir -p "$BENCH_TEST_STATE"
printf 'charged\n' > charged.txt
cp charged.txt "$BENCH_TEST_STATE/charged"
cp .bench-objective "$BENCH_TEST_STATE/objective"
cp .bench-notes.md "$BENCH_TEST_STATE/notes"
cp "$(git rev-parse --git-path bench-lease)" "$BENCH_TEST_STATE/lease"
cp "$(git rev-parse --path-format=absolute --git-common-dir)/bench-intent.json" "$BENCH_TEST_STATE/intent"
git status --porcelain=v1 --untracked-files=all > "$BENCH_TEST_STATE/status"
git diff --binary HEAD -- > "$BENCH_TEST_STATE/diff"
`)
	if variant == "drift" {
		f.WriteFile("tracked.txt", "base\n")
		f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":["BENCH_TEST_PROMPTS","BENCH_TEST_STATE"],"paths":["tracked.txt"],"tools":[]}`+"\n")
	}
	f.CommitAll("R15 shift fixture")
	home := t.TempDir()
	return shiftProofRun{f: f, home: home, pool: addRuntimePoolWorktrees(t, f, home), base: headSha(f), variant: variant}
}

func (r shiftProofRun) env() map[string]string {
	env := map[string]string{
		"BENCH_AGENT":      filepath.Join(r.f.Root, "agent"),
		"BENCH_HOME":       r.home,
		"BENCH_MAX_ITERS":  "1",
		"BENCH_TEST_STATE": filepath.Join(r.home, "preservation-state"),
	}
	return env
}

func shiftProofGitDir(t *testing.T, run shiftProofRun) string {
	t.Helper()
	return strings.TrimSpace(contract.RunAt(t, run.f, run.pool.Warm, nil, "git", "rev-parse", "--absolute-git-dir").Stdout)
}

func holdShiftGateLock(t *testing.T, run shiftProofRun) *os.File {
	t.Helper()
	lock := filepath.Join(shiftProofGitDir(t, run), "bench-gate.lock")
	held, err := os.OpenFile(lock, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	acquireTestGateLock(t, held)
	return held
}

func proveCancelledShift(t *testing.T, run shiftProofRun) {
	t.Helper()
	state := run.env()["BENCH_TEST_STATE"]
	if err := os.MkdirAll(state, 0o700); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	cmd := exec.Command("bash", benchPath(t), "shift", "R15-cancellation")
	cmd.Dir = run.f.Root
	cmd.Env = surfaceEnv(run.f, run.env())
	cmd.Stdout, cmd.Stderr = &out, &out
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	started := filepath.Join(state, "gate-started")
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(started); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			t.Fatalf("shift gate did not start: %s", out.String())
		}
		time.Sleep(10 * time.Millisecond)
	}
	var gatePGID int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(mustReadRuntime(t, filepath.Join(state, "gate-pgid")))), "%d", &gatePGID); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		_ = syscall.Kill(-gatePGID, syscall.SIGKILL)
		t.Fatalf("cancelled shift did not exit: %s", out.String())
	}
	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if !ok || status.ExitStatus() != 130 {
		t.Fatalf("cancelled shift status = %v, want exit 130: %s", cmd.ProcessState.Sys(), out.String())
	}
	deadline = time.Now().Add(3 * time.Second)
	for syscall.Kill(-gatePGID, 0) == nil && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if syscall.Kill(-gatePGID, 0) == nil {
		t.Fatalf("cancelled gate process group %d survived", gatePGID)
	}
	requirePreservedShift(t, run, out.String())
}

func requirePreservedShift(t *testing.T, run shiftProofRun, output string) {
	t.Helper()
	worktree := shiftWorktree(t, output)
	if worktree != run.pool.Warm {
		t.Fatalf("shift acquired %q, want seeded warm worktree %q", worktree, run.pool.Warm)
	}
	branch := shiftBranchFromStart(t, output)
	state := run.env()["BENCH_TEST_STATE"]

	// FT79 replaces the old "retain the charged worktree verbatim" special case with the
	// uniform rule: snapshot the dirty tree to refs/bench/recovery/<branch> and release
	// the pool worktree normally. The preserved evidence now lives in the recovery ref's
	// tree, not the (now-clean, released) physical worktree.
	ref := "refs/bench/recovery/" + branch
	run.f.Git("show-ref", "--verify", ref)
	tree := run.f.Git("ls-tree", "-r", "--name-only", ref).Stdout
	if !strings.Contains(tree, "charged.txt") {
		t.Fatalf("recovery snapshot did not preserve charged.txt:\n%s", tree)
	}
	for _, scratch := range []string{".bench-objective", ".bench-notes.md"} {
		if strings.Contains(tree, scratch) {
			t.Fatalf("recovery snapshot rode scratch %s into the tree:\n%s", scratch, tree)
		}
	}
	if got := run.f.Git("show", ref+":charged.txt").Stdout; got != string(mustReadRuntime(t, filepath.Join(state, "charged"))) {
		t.Fatalf("recovery snapshot charged.txt = %q, want %q", got, string(mustReadRuntime(t, filepath.Join(state, "charged"))))
	}
	requireNoWorktreeLease(t, worktree)
	if got := strings.TrimSpace(runGitAt(t, worktree, "status", "--porcelain")); got != "" {
		t.Fatalf("released pool worktree is not clean:\n%s", got)
	}

	intentPath, err := intent.Address(run.f.Root)
	if err != nil {
		t.Fatal(err)
	}
	// FT79 now upserts the shift's final outcome and recovery pointer onto its intent
	// entry at every exit path (including this preserving one), so the ledger is no
	// longer byte-identical to the pre-outcome snapshot the agent captured mid-run —
	// only those two fields legitimately drift, to the values this preserving failure
	// resolves to.
	wantOutcome := "failed"
	if run.variant == "cancellation" {
		wantOutcome = "interrupted"
	}
	requirePreservedIntentLedger(t, intentPath, filepath.Join(state, "intent"), wantOutcome, "ref:"+ref)
	entries, err := intent.Snapshot(run.f.Root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Kind == intent.KindShift && entry.Objective == "R15-"+run.variant && entry.Worktree == worktree && entry.Branch == branch {
			return
		}
	}
	t.Fatalf("failed shift lost correlated intent for branch=%s worktree=%s: %#v", branch, worktree, entries)
}

func requirePreservedFile(t *testing.T, gotPath, snapshotPath string) {
	t.Helper()
	gotInfo, gotErr := os.Lstat(gotPath)
	wantInfo, wantErr := os.Lstat(snapshotPath)
	if gotErr != nil || wantErr != nil {
		t.Fatalf("preserved file existence %s: got=%v snapshot=%v", gotPath, gotErr, wantErr)
	}
	if !gotInfo.Mode().IsRegular() || !wantInfo.Mode().IsRegular() || gotInfo.Mode() != wantInfo.Mode() {
		t.Fatalf("preserved file mode %s: got=%v snapshot=%v, want matching regular files", gotPath, gotInfo.Mode(), wantInfo.Mode())
	}
	got, gotErr := os.ReadFile(gotPath)
	want, wantErr := os.ReadFile(snapshotPath)
	if gotErr != nil || wantErr != nil || !bytes.Equal(got, want) {
		t.Fatalf("preserved file bytes %s: got=%q/%v snapshot=%q/%v", gotPath, got, gotErr, want, wantErr)
	}
}

// requireNoWorktreeLease asserts the pool worktree's bench-lease file is gone — the
// signal that preserveAndRecover released it back to the pool rather than retaining it.
func requireNoWorktreeLease(t *testing.T, worktree string) {
	t.Helper()
	lease := strings.TrimSpace(runGitAt(t, worktree, "rev-parse", "--git-path", "bench-lease"))
	if !filepath.IsAbs(lease) {
		lease = filepath.Join(worktree, lease)
	}
	if _, err := os.Stat(lease); err == nil {
		t.Fatalf("worktree lease %s was not released", lease)
	}
}

// requirePreservedIntentLedger compares the current intent ledger to the pre-outcome
// snapshot the agent captured mid-run, tolerating exactly two legitimate drifts: each
// entry's Outcome and Recovery fields, which FT79's finish()/checkpoint() now set at
// every shift exit path — including this preserving one. Every other field, and the
// entry count and order, must match byte-for-byte-equivalent (via a decoded structural
// compare, since the Outcome/Recovery fields shift the encoded byte length).
func requirePreservedIntentLedger(t *testing.T, gotPath, snapshotPath, wantOutcome, wantRecovery string) {
	t.Helper()
	gotInfo, gotErr := os.Lstat(gotPath)
	wantInfo, wantErr := os.Lstat(snapshotPath)
	if gotErr != nil || wantErr != nil {
		t.Fatalf("preserved intent ledger existence %s: got=%v snapshot=%v", gotPath, gotErr, wantErr)
	}
	if !gotInfo.Mode().IsRegular() || !wantInfo.Mode().IsRegular() {
		t.Fatalf("preserved intent ledger mode %s: got=%v snapshot=%v, want matching regular files", gotPath, gotInfo.Mode(), wantInfo.Mode())
	}
	gotBytes, err := os.ReadFile(gotPath)
	if err != nil {
		t.Fatal(err)
	}
	wantBytes, err := os.ReadFile(snapshotPath)
	if err != nil {
		t.Fatal(err)
	}
	var got, want intent.Ledger
	if err := json.Unmarshal(gotBytes, &got); err != nil {
		t.Fatalf("decode intent ledger %s: %v", gotPath, err)
	}
	if err := json.Unmarshal(wantBytes, &want); err != nil {
		t.Fatalf("decode intent ledger snapshot %s: %v", snapshotPath, err)
	}
	if len(got.Entries) != len(want.Entries) {
		t.Fatalf("preserved intent ledger entry count = %d, want %d", len(got.Entries), len(want.Entries))
	}
	for i := range got.Entries {
		g, w := got.Entries[i], want.Entries[i]
		if g.Outcome != wantOutcome {
			t.Fatalf("preserved intent entry[%d] outcome = %q, want %q", i, g.Outcome, wantOutcome)
		}
		if g.Recovery != wantRecovery {
			t.Fatalf("preserved intent entry[%d] recovery = %q, want %q", i, g.Recovery, wantRecovery)
		}
		g.Outcome, g.Recovery = "", ""
		if g != w {
			t.Fatalf("preserved intent entry[%d] drifted beyond outcome/recovery:\ngot=%#v\nwant=%#v", i, g, w)
		}
	}
}

func shiftBranchFromStart(t *testing.T, output string) string {
	t.Helper()
	const prefix = "▶ shift on "
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			branch, _, ok := strings.Cut(strings.TrimPrefix(line, prefix), " — objective:")
			if ok && branch != "" {
				return branch
			}
		}
	}
	t.Fatalf("shift start did not name branch:\n%s", output)
	return ""
}
