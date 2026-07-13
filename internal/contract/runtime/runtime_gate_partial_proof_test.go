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
	"github.com/gibbonmi/bench/internal/gate"
)

var r17BridgeBinary string

func buildR17Bridge(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "gate.test")
	cmd := exec.Command("go", "test", "-c", "-o", path, "./internal/gate")
	cmd.Dir = contract.SubjectRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("compile R17 private bridge: %v\n%s", err, out)
	}
	return path
}

func proveR17Result(t *testing.T, variant string) {
	switch {
	case variant == "drift":
		proveR17Drift(t)
	case variant == "cancellation":
		proveR17Cancellation(t)
	case strings.HasPrefix(variant, "invalid-"):
		proveR17Invalid(t, strings.TrimPrefix(variant, "invalid-"))
	case strings.HasPrefix(variant, "fault-"):
		proveR17Fault(t, strings.TrimPrefix(variant, "fault-"))
	default:
		t.Fatalf("unknown R17 proof %q", variant)
	}
}

func proveR17Fault(t *testing.T, op string) {
	t.Helper()
	cmd := exec.Command(r17BridgeBinary, "-test.run=^TestR17PrivateFaultBridge$", "-test.v")
	cmd.Env = append(os.Environ(), "FT78_R17_FAULT="+op)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("private fault %s: %v\n%s", op, err, out)
	}
	durable, attempts := "interrupted-pending", 2
	if op == "lock-open" || op == "lock-acquisition" {
		durable = "ready-green"
		attempts = 1
	} else if op == "post-run-subject-rebuild" {
		attempts = 1
	}
	want := fmt.Sprintf("R17-TUPLE op=%s calls=2 gate=0 action=1 returned_reusable=false durable=%s attempts=%d,%d temps=0", op, durable, attempts, attempts)
	if !strings.Contains(string(out), want) {
		t.Fatalf("private fault tuple = %q, want literal %q", out, want)
	}
}

func proveR17Drift(t *testing.T) {
	t.Helper()
	f := proofFixture(t)
	f.WriteFile("tracked.txt", "base\n")
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":["tracked.txt"],"tools":[]}`)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\necho run >> .git/ft78-runs\nprintf drift >> tracked.txt\n")
	f.CommitAll("drift fixture")
	for call := 1; call <= 2; call++ {
		f.Bench("gate").RequireExit(1)
		got := gate.Inspect(f.Root)
		if got.State != gate.Pending || got.PendingStatus != "interrupted-pending" || got.ReusableGreen {
			t.Fatalf("drift call %d = %+v, want interrupted pending", call, got)
		}
	}
	assertRuns(t, f, 2)
	assertNoR17Temps(t, f)
}

func proveR17Invalid(t *testing.T, kind string) {
	t.Helper()
	f := proofFixture(t)
	f.WriteFile(".bench/gate-inputs.json", localManifest)
	f.Bench("gate").RequireExit(0)
	cache := filepath.Join(gitDir(t, f), "bench-last-gate")
	for call := 1; call <= 2; call++ {
		plantR17Invalid(t, cache, kind)
		if got := gate.Inspect(f.Root); got.State != gate.Invalid || got.ReusableGreen {
			t.Fatalf("%s call %d planted inspection = %+v", kind, call, got)
		}
		if kind == "directory" {
			f.Bench("gate").RequireExit(1)
			if got := gate.Inspect(f.Root); got.State != gate.Invalid || got.ReusableGreen {
				t.Fatalf("%s call %d failed recovery = %+v", kind, call, got)
			}
			continue
		}
		f.Bench("gate").RequireExit(0)
		if got := gate.Inspect(f.Root); got.State != gate.Ready || got.Status != "green" || !got.ReusableGreen {
			t.Fatalf("%s call %d recovery = %+v", kind, call, got)
		}
	}
	wantRuns := 3
	if kind == "directory" {
		wantRuns = 1
	}
	assertRuns(t, f, wantRuns)
	assertNoR17Temps(t, f)
}

func plantR17Invalid(t *testing.T, cache, kind string) {
	t.Helper()
	if kind == "wrong-mode" {
		if err := os.Chmod(cache, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	if err := os.RemoveAll(cache); err != nil {
		t.Fatal(err)
	}
	switch kind {
	case "zero-byte":
		contract.WriteFileAbs(t, cache, "")
	case "malformed":
		contract.WriteFileAbs(t, cache, "{\n")
	case "legacy":
		contract.WriteFileAbs(t, cache, "green old-tree\n")
	case "oversized":
		contract.WriteFileAbs(t, cache, strings.Repeat("x", 16_385))
	case "symlink":
		target := cache + "-r17-target"
		contract.WriteFileAbs(t, target, "old green\n")
		if err := os.Symlink(target, cache); err != nil {
			t.Fatal(err)
		}
	case "directory":
		if err := os.Mkdir(cache, 0o700); err != nil {
			t.Fatal(err)
		}
	default:
		t.Fatalf("unknown invalid record %q", kind)
	}
}

func proveR17Cancellation(t *testing.T) {
	t.Helper()
	f := proofFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\necho run >> .git/ft78-runs\necho $$ > .git/r17-pgid\ntouch .git/r17-started\nwhile :; do sleep .05; done\n")
	for call := 1; call <= 2; call++ {
		cmd := exec.Command("bash", benchPath(t), "gate")
		cmd.Dir, cmd.Env, cmd.SysProcAttr = f.Root, surfaceEnv(f, nil), &syscall.SysProcAttr{Setpgid: true}
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		waitR17File(t, filepath.Join(gitDir(t, f), "r17-started"))
		var gatePGID int
		if _, err := fmt.Sscanf(strings.TrimSpace(string(mustReadRuntime(t, filepath.Join(gitDir(t, f), "r17-pgid")))), "%d", &gatePGID); err != nil {
			t.Fatal(err)
		}
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		_ = syscall.Kill(-gatePGID, syscall.SIGKILL)
		if err := cmd.Wait(); err == nil {
			t.Fatalf("cancellation call %d exited zero: %s", call, out.String())
		}
		contract.Remove(t, filepath.Join(gitDir(t, f), "r17-started"))
		if got := gate.Inspect(f.Root); got.State != gate.Pending || got.PendingStatus != "interrupted-pending" || got.ReusableGreen {
			t.Fatalf("cancellation call %d = %+v", call, got)
		}
	}
	assertRuns(t, f, 2)
	assertNoR17Temps(t, f)
}

func waitR17File(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("marker %s not written", path)
}

func assertNoR17Temps(t *testing.T, f contract.Fixture) {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(gitDir(t, f), ".bench-last-gate-*"))
	if err != nil || len(paths) != 0 {
		t.Fatalf("temporary gate evidence = %v/%v", paths, err)
	}
}
