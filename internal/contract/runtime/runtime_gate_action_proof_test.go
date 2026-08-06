package runtime

import (
	"bytes"
	"context"
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

type actionPathSnapshot struct {
	exists bool
	mode   os.FileMode
	bytes  string
}

type actionSnapshot struct {
	head, index, status string
	work, spec          actionPathSnapshot
}

func snapshotActionPath(t *testing.T, path string) actionPathSnapshot {
	t.Helper()
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return actionPathSnapshot{}
	}
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return actionPathSnapshot{exists: true, mode: info.Mode(), bytes: string(data)}
}

func snapshotAction(t *testing.T, f contract.Fixture) actionSnapshot {
	t.Helper()
	return actionSnapshot{
		head: headSha(f), index: strings.TrimSpace(f.Git("write-tree").Stdout),
		status: f.Git("status", "--porcelain=v1", "--untracked-files=all", "--", "work.txt", "specs/proof/spec.md").Stdout,
		work:   snapshotActionPath(t, filepath.Join(f.Root, "work.txt")),
		spec:   snapshotActionPath(t, filepath.Join(f.Root, "specs", "proof", "spec.md")),
	}
}

func requireActionUnchanged(t *testing.T, f contract.Fixture, before actionSnapshot) {
	t.Helper()
	after := snapshotAction(t, f)
	if after != before {
		t.Fatalf("failed action changed HEAD/index/status/spec/path\nbefore=%+v\nafter=%+v", before, after)
	}
}

func story5CoordinatingGate(body string) string {
	return "#!/usr/bin/env bash\n" + commonGitDirGateBody(body)
}

func proveStory5GateUsesCommonDirFromProspectiveCheckout(t *testing.T, f contract.Fixture) {
	t.Helper()
	var output bytes.Buffer
	tree := strings.TrimSpace(f.Git("rev-parse", "HEAD^{tree}").Stdout)
	if got := gatepkg.ExecuteTree(t.Context(), f.Root, tree, &output, &output); got.ActionExit != 23 {
		t.Fatalf("prospective interrupted gate exit = %d, want 23\n%s", got.ActionExit, output.String())
	}
}

func proveStory5CancellationFromProspectiveCheckout(t *testing.T, f contract.Fixture) {
	t.Helper()
	gitdir := gitDir(t, f)
	marker := filepath.Join(gitdir, "story5-cancel-started")
	pidPath := filepath.Join(gitdir, "story5-cancel-pgid")
	contract.Remove(t, marker)
	contract.Remove(t, pidPath)
	defer contract.Remove(t, marker)
	defer contract.Remove(t, pidPath)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var output bytes.Buffer
	var got gatepkg.Result
	done := make(chan struct{})
	tree := strings.TrimSpace(f.Git("rev-parse", "HEAD^{tree}").Stdout)
	go func() {
		got = gatepkg.ExecuteTree(ctx, f.Root, tree, &output, &output)
		close(done)
	}()
	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		select {
		case <-done:
			t.Fatalf("prospective cancellation gate exited before its common-directory marker: exit=%d\n%s", got.ActionExit, output.String())
		default:
		}
		if time.Now().After(deadline) {
			cancel()
			select {
			case <-done:
				t.Fatalf("prospective cancellation gate missed its common-directory marker: exit=%d\n%s", got.ActionExit, output.String())
			case <-time.After(3 * time.Second):
				t.Fatal("prospective cancellation gate missed its common-directory marker and could not be reaped")
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	var gatePGID int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(mustReadRuntime(t, pidPath))), "%d", &gatePGID); err != nil {
		t.Fatal(err)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out cancelling prospective gate")
	}
	if got.ActionExit != 130 {
		t.Fatalf("prospective cancellation gate exit = %d, want 130\n%s", got.ActionExit, output.String())
	}
	contract.RequireProcessGroupDrained(t, gatePGID, 3*time.Second, "prospective cancellation gate survived")
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
		gateBody = story5CoordinatingGate("rm -f \"$gitdir/bench-last-gate\"; mkdir \"$gitdir/bench-last-gate\"; exit 0\n")
	case "drift":
		gateBody = "#!/usr/bin/env bash\nprintf drift >> tracked.txt\nexit 0\n"
		manifest = `{"schema":1,"closure":"local","environment":[],"paths":["tracked.txt"],"tools":[]}` + "\n"
	case "open-green":
		manifest = `{"schema":1,"closure":"remote","environment":[],"paths":[],"tools":[]}` + "\n"
		gateBody = story5CoordinatingGate("[ ! -f \"$gitdir/story5-red\" ]\n")
	case "stale-green":
		gateBody = story5CoordinatingGate("[ ! -f \"$gitdir/story5-red\" ]\n")
	case "locked-pending", "interrupted":
		gateBody = story5CoordinatingGate("if [ ! -f \"$gitdir/story5-owner-once\" ]; then\n  touch \"$gitdir/story5-owner-once\"\n  echo $$ > \"$gitdir/story5-gate-pgid\"\n  touch \"$gitdir/story5-owner-started\"\n  while :; do sleep .05; done\nfi\nexit 23\n")
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
		acquireTestGateLock(t, held)
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
	result := p.Stdout + p.Stderr
	want := map[string]string{"lock-open": "gate lock unavailable", "lock-acquire": "gate execution already in progress", "pending-persistence": "gate pending persistence failed", "final-persistence": "gate final persistence failed", "drift": "gate subject changed during execution", "no-gate": "no gate found"}[variant]
	if want != "" && !strings.Contains(p.Stdout+p.Stderr, want) {
		t.Fatalf("%s result missing %q:\n%s%s", variant, want, p.Stdout, p.Stderr)
	}
	if variant == "subject-build" && !strings.Contains(result, "gate subject unavailable") && !strings.Contains(result, "exit status") {
		t.Fatalf("subject-build result did not surface the subject or composition failure:\n%s", result)
	}
	if variant == "start-failure" && !strings.Contains(result, "gate is red") && !strings.Contains(result, "prospective authorization refused") {
		t.Fatalf("start-failure result did not identify the gate or prospective authorization stage:\n%s%s", p.Stdout, p.Stderr)
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
	if variant == "interrupted" {
		proveStory5GateUsesCommonDirFromProspectiveCheckout(t, f)
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
	f.WriteExecutable(".bench/gate.sh", story5CoordinatingGate("echo $$ > \"$gitdir/story5-cancel-pgid\"\ntouch \"$gitdir/story5-cancel-started\"\nwhile :; do sleep .05; done\n"))
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
	proveStory5CancellationFromProspectiveCheckout(t, f)
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
	f.WriteFile("specs/proof/spec.md", "# proof\nStatus: staged\n")
	f.CommitAll("base")
	return f
}
