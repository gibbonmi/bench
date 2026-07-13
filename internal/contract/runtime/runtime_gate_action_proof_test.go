package runtime

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	gatepkg "github.com/gibbonmi/bench/internal/gate"
)

const ft78LedgerFailure = "FT78 proof ledger completeness contract failed"

type actionProof struct {
	id     string
	driver func(*testing.T)
}

func proof(id string, driver func(*testing.T)) actionProof { return actionProof{id, driver} }
func commitProof(id, variant string) actionProof {
	return proof(id, func(t *testing.T) { proveCommitResult(t, variant) })
}

var ft78Story5Proofs = []actionProof{
	commitProof("R14/exact-green-reuse", "reuse"), commitProof("R14/stale-tree-rerun", "stale"),
	commitProof("R14/ordinary-red", "red"), commitProof("R14/oracle-mismatch", "oracle-mismatch"),
	commitProof("R14/absent-inspection", "absent"), commitProof("R14/ready-red-inspection", "ready-red"),
	commitProof("R14/stale-green-inspection", "stale-green"), commitProof("R14/open-subject-green-inspection", "open-green"),
	commitProof("R14/locked-pending-inspection", "locked-pending"), commitProof("R14/interrupted-pending-inspection", "interrupted"),
	commitProof("R14/invalid-inspection", "invalid"), commitProof("R14/unavailable-inspection", "unavailable"),
	commitProof("R14/lock-open-result", "lock-open"), commitProof("R14/lock-acquire-result", "lock-acquire"),
	commitProof("R14/pending-persistence-result", "pending-persistence"), commitProof("R14/final-persistence-result", "final-persistence"),
	commitProof("R14/subject-build-result", "subject-build"), commitProof("R14/subject-recheck-drift-result", "drift"),
	commitProof("R14/cancellation-result", "cancellation"), commitProof("R14/start-failure-result", "start-failure"),
	commitProof("R14/no-gate-result", "no-gate"),
}

var ft78Story5ExpectedIDs = []string{
	"R14/exact-green-reuse", "R14/stale-tree-rerun", "R14/ordinary-red", "R14/oracle-mismatch", "R14/absent-inspection", "R14/ready-red-inspection", "R14/stale-green-inspection", "R14/open-subject-green-inspection", "R14/locked-pending-inspection", "R14/interrupted-pending-inspection", "R14/invalid-inspection", "R14/unavailable-inspection", "R14/lock-open-result", "R14/lock-acquire-result", "R14/pending-persistence-result", "R14/final-persistence-result", "R14/subject-build-result", "R14/subject-recheck-drift-result", "R14/cancellation-result", "R14/start-failure-result", "R14/no-gate-result",
}

func TestFT78Story5ProofLedgerCompleteness(t *testing.T) {
	contract.NoteContractFailure(t, ft78LedgerFailure)
	seen := map[string]int{}
	for _, c := range ft78Story5Proofs {
		seen[c.id]++
		if c.driver == nil {
			t.Fatalf("%s: nil real driver", c.id)
		}
	}
	if len(seen) != len(ft78Story5ExpectedIDs) {
		t.Fatalf("%s: registered IDs = %d, want %d", ft78LedgerFailure, len(seen), len(ft78Story5ExpectedIDs))
	}
	for _, id := range ft78Story5ExpectedIDs {
		if seen[id] != 1 {
			t.Fatalf("%s: %s registrations = %d, want 1", ft78LedgerFailure, id, seen[id])
		}
	}
}

func TestFT78Story5ProofLedger(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	for _, c := range ft78Story5Proofs {
		c := c
		t.Run(c.id, c.driver)
	}
}

type actionSnapshot struct{ head, index, status, spec string }

func snapshotAction(t *testing.T, f contract.Fixture) actionSnapshot {
	t.Helper()
	s := actionSnapshot{head: headSha(f), index: strings.TrimSpace(f.Git("write-tree").Stdout), status: f.Git("status", "--porcelain=v1").Stdout}
	if f.Exists("specs/proof.md") {
		s.spec = f.ReadFile("specs/proof.md")
	}
	return s
}

func requireActionUnchanged(t *testing.T, f contract.Fixture, before actionSnapshot) {
	t.Helper()
	after := snapshotAction(t, f)
	if after.head != before.head || after.index != before.index || after.spec != before.spec || !strings.HasPrefix(f.ReadFile("work.txt"), "charged") {
		t.Fatalf("failed action changed HEAD/index/status/spec/path\nbefore=%+v\nafter=%+v", before, after)
	}
}

func proveCommitResult(t *testing.T, variant string) {
	if variant == "reuse" {
		testCommitFreshVerdictReused(t)
		return
	}
	if variant == "stale" {
		testCommitStaleVerdictRerunsGate(t)
		return
	}
	gateBody := "#!/usr/bin/env bash\nexit 23\n"
	manifest := `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}` + "\n"
	switch variant {
	case "oracle-mismatch":
		gateBody = ""
	case "start-failure":
		gateBody = "#!/definitely/missing\n"
	case "final-persistence":
		gateBody = "#!/usr/bin/env bash\nrm -f .git/bench-last-gate; mkdir .git/bench-last-gate; exit 0\n"
	case "drift":
		gateBody = "#!/usr/bin/env bash\nprintf drift >> tracked.txt\nexit 0\n"
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":["tracked.txt"],"tools":[]}` + "\n"
	case "open-green":
		manifest = `{"schema":1,"closure":"remote","environment":[],"paths":[],"tools":[]}` + "\n"
		gateBody = "#!/usr/bin/env bash\n[ ! -f .git/story5-red ]\n"
	case "stale-green":
		gateBody = "#!/usr/bin/env bash\n[ ! -f .git/story5-red ]\n"
	case "locked-pending", "interrupted":
		gateBody = "#!/usr/bin/env bash\nif [ ! -f .git/story5-owner-once ]; then\n  touch .git/story5-owner-once\n  echo $$ > .git/story5-gate-pgid\n  touch .git/story5-owner-started\n  while :; do sleep .05; done\nfi\nexit 23\n"
	case "no-gate":
		gateBody = ""
	}
	f := story5ActionFixture(t, gateBody, manifest)
	if variant == "open-green" {
		f.Bench("gate").RequireExit(0)
		contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "story5-red"), "red\n")
	}
	f.WriteFile("work.txt", "charged\n")
	env := map[string]string{}
	if variant == "oracle-mismatch" {
		f.WriteExecutable("gate-green", "#!/usr/bin/env bash\nexit 0\n")
		f.WriteExecutable("gate-red", "#!/usr/bin/env bash\nexit 23\n")
		f.CommitAll("oracle commands")
		f.WriteFile("work.txt", "charged\n")
		env["BENCH_GATE"] = filepath.Join(f.Root, "gate-green")
		f.BenchEnv(env, "gate").RequireExit(0)
		env["BENCH_GATE"] = filepath.Join(f.Root, "gate-red")
	}
	cache, lockPath := filepath.Join(gitDir(t, f), "bench-last-gate"), filepath.Join(gitDir(t, f), "bench-gate.lock")
	if variant == "unavailable" || variant == "subject-build" {
		env["PATH"] = story5FailingGitPath(t) + string(os.PathListSeparator) + os.Getenv("PATH")
	}
	var held *os.File
	var owner *story5GateOwner
	switch variant {
	case "ready-red":
		f.Bench("gate").RequireExit(23)
	case "stale-green":
		f.Bench("gate").RequireExit(0)
		contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "story5-red"), "red\n")
		f.WriteFile("work.txt", "charged-again\n")
	case "invalid":
		contract.WriteFileAbs(t, cache, "invalid\n")
	case "locked-pending", "interrupted":
		owner = startStory5GateOwner(t, f)
		if variant == "interrupted" {
			owner.stop(t)
			owner = nil
		} else {
			defer owner.stop(t)
		}
	case "lock-open":
		if err := os.Mkdir(lockPath, 0o700); err != nil {
			t.Fatal(err)
		}
	case "lock-acquire":
		var err error
		held, err = os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		if err = syscall.Flock(int(held.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
			t.Fatal(err)
		}
		defer held.Close()
	case "pending-persistence":
		if err := os.Mkdir(cache, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	wantState := map[string]gatepkg.State{"absent": gatepkg.Absent, "ready-red": gatepkg.Ready, "stale-green": gatepkg.Ready, "open-green": gatepkg.Ready, "locked-pending": gatepkg.Pending, "interrupted": gatepkg.Pending, "invalid": gatepkg.Invalid, "unavailable": gatepkg.Unavailable}[variant]
	if variant == "unavailable" {
		old := os.Getenv("PATH")
		_ = os.Setenv("PATH", env["PATH"])
		got := gatepkg.Inspect(f.Root)
		_ = os.Setenv("PATH", old)
		if got.State != gatepkg.Unavailable {
			t.Fatalf("unavailable pre-action inspection = %+v", got)
		}
	}
	if wantState != "" && variant != "unavailable" {
		got := gatepkg.Inspect(f.Root)
		if got.State != wantState {
			t.Fatalf("%s pre-action inspection = %+v, want state %s", variant, got, wantState)
		}
		if variant == "locked-pending" && got.PendingStatus != "locked-pending" {
			t.Fatalf("locked pending projection = %+v", got)
		}
		if variant == "interrupted" && got.PendingStatus != "interrupted-pending" {
			t.Fatalf("interrupted pending projection = %+v", got)
		}
		if variant == "ready-red" && (got.Status != "red" || got.ReusableGreen) {
			t.Fatalf("ready-red inspection = %+v", got)
		}
		if (variant == "stale-green" || variant == "open-green") && (got.Status != "green" || got.ReusableGreen) {
			t.Fatalf("%s inspection = %+v, want non-reusable green", variant, got)
		}
		if variant == "invalid" && got.Reason == "" {
			t.Fatalf("invalid inspection omitted reason: %+v", got)
		}
	}
	before := snapshotAction(t, f)
	if variant == "cancellation" {
		proveCancelledCommit(t, f, before)
		return
	}
	p := f.BenchEnv(env, "commit", "-m", "must refuse", "--spec", "proof", "work.txt")
	if p.ExitCode == 0 {
		t.Fatalf("%s commit authorized without reusable green", variant)
	}
	want := map[string]string{"lock-open": "gate lock unavailable", "lock-acquire": "gate execution already in progress", "pending-persistence": "gate pending persistence failed", "final-persistence": "gate final persistence failed", "drift": "gate subject changed during execution", "no-gate": "no gate found", "subject-build": "gate subject unavailable", "start-failure": "gate is red"}[variant]
	if want != "" && !strings.Contains(p.Stdout+p.Stderr, want) {
		t.Fatalf("%s result missing %q:\n%s%s", variant, want, p.Stdout, p.Stderr)
	}
	if variant == "start-failure" {
		got := gatepkg.Inspect(f.Root)
		if got.State != gatepkg.Ready || got.Status != "red" {
			t.Fatalf("start-failure result = %+v, want ready red", got)
		}
	}
	if variant == "no-gate" {
		if got := gatepkg.Inspect(f.Root); got.State != gatepkg.Absent {
			t.Fatalf("no-gate inspection = %+v, want absent", got)
		}
	}
	requireActionUnchanged(t, f, before)
}

func story5FailingGitPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	real, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	contract.WriteExecutableAbs(t, filepath.Join(dir, "git"), "#!/usr/bin/env bash\n[ -n \"${GIT_INDEX_FILE:-}\" ] && exit 1\nexec "+real+" \"$@\"\n")
	return dir
}

func proveCancelledCommit(t *testing.T, f contract.Fixture, before actionSnapshot) {
	t.Helper()
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\necho $$ > .git/story5-cancel-pgid\ntouch .git/story5-cancel-started\nwhile :; do sleep .05; done\n")
	f.Git("add", ".bench/gate.sh")
	f.Git("commit", "-q", "-m", "blocking gate")
	before = snapshotAction(t, f)
	var out bytes.Buffer
	cmd := exec.Command("bash", benchPath(t), "commit", "-m", "cancel", "--spec", "proof", "work.txt")
	cmd.Dir = f.Root
	cmd.Env = surfaceEnv(f, nil)
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(gitDir(t, f), "story5-cancel-started")
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			t.Fatal("commit gate did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	var gatePGID int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(mustReadRuntime(t, filepath.Join(gitDir(t, f), "story5-cancel-pgid")))), "%d", &gatePGID); err != nil {
		t.Fatal(err)
	}
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
	_ = syscall.Kill(-gatePGID, syscall.SIGKILL)
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	var waitErr error
	select {
	case waitErr = <-done:
	case <-time.After(3 * time.Second):
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		t.Fatalf("timed out reaping cancelled commit and gate pgid=%d", gatePGID)
	}
	if waitErr == nil {
		t.Fatalf("cancelled commit exited zero: %s", out.String())
	}
	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if !ok || !status.Signaled() || status.Signal() != syscall.SIGINT {
		t.Fatalf("cancelled commit status = %v, want SIGINT", cmd.ProcessState.Sys())
	}
	deadline = time.Now().Add(3 * time.Second)
	for syscall.Kill(-gatePGID, 0) == nil && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if syscall.Kill(-gatePGID, 0) == nil {
		t.Fatalf("cancelled gate process group %d survived", gatePGID)
	}
	if got := gatepkg.Inspect(f.Root); got.State != gatepkg.Pending {
		t.Fatalf("cancelled commit inspection = %+v, want pending", got)
	}
	requireActionUnchanged(t, f, before)
}

type story5GateOwner struct {
	cmd      *exec.Cmd
	gatePGID int
}

func startStory5GateOwner(t *testing.T, f contract.Fixture) *story5GateOwner {
	t.Helper()
	cmd := exec.Command("bash", benchPath(t), "gate")
	cmd.Dir = f.Root
	cmd.Env = surfaceEnv(f, nil)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	marker := filepath.Join(gitDir(t, f), "story5-owner-started")
	for time.Now().Before(deadline) {
		if _, err := os.Stat(marker); err == nil {
			var pgid int
			if _, err := fmt.Sscanf(strings.TrimSpace(string(mustReadRuntime(t, filepath.Join(gitDir(t, f), "story5-gate-pgid")))), "%d", &pgid); err != nil {
				t.Fatal(err)
			}
			return &story5GateOwner{cmd: cmd, gatePGID: pgid}
		}
		time.Sleep(10 * time.Millisecond)
	}
	_ = cmd.Process.Kill()
	t.Fatal("gate owner did not reach pending state")
	return nil
}

func (o *story5GateOwner) stop(t *testing.T) {
	t.Helper()
	_ = syscall.Kill(-o.gatePGID, syscall.SIGKILL)
	_ = syscall.Kill(-o.cmd.Process.Pid, syscall.SIGKILL)
	done := make(chan error, 1)
	go func() { done <- o.cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatalf("timed out reaping Story5 gate owner pgid=%d", o.gatePGID)
	}
}

func mustReadRuntime(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func story5ActionFixture(t *testing.T, gateBody, manifest string) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	f.Git("config", "user.email", "bench@local")
	f.Git("config", "user.name", "bench")
	if gateBody != "" {
		f.WriteExecutable(".bench/gate.sh", gateBody)
		f.WriteFile(".bench/gate-inputs.json", manifest)
	}
	f.WriteFile("tracked.txt", "base\n")
	f.WriteFile("specs/proof.md", "# proof\nStatus: staged\n")
	f.CommitAll("base")
	return f
}
