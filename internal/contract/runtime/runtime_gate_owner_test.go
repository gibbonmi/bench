package runtime

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/gate"
)

func testRuntimeGateOwnerAndSignalCleanup(t *testing.T) {
	for _, test := range []struct {
		name   string
		signal syscall.Signal
	}{
		{name: "SIGINT", signal: syscall.SIGINT},
		{name: "SIGTERM", signal: syscall.SIGTERM},
		{name: "SIGHUP", signal: syscall.SIGHUP},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			testRuntimeGateSignalCleanup(t, test.signal)
		})
	}
}

func testRuntimeGateSignalCleanup(t *testing.T, signal syscall.Signal) {
	f := contract.NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", `#!/usr/bin/env bash
echo $$ > .git/ft88-gate-pgid
touch .git/ft88-gate-started
while :; do sleep .05; done
`)
	f.WriteFile(".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`+"\n")
	f.CommitAll("blocking gate")
	gitdir := gitDir(t, f)
	ownerPath := filepath.Join(gitdir, "bench-gate-owner")
	var output bytes.Buffer
	cmd := exec.Command("bash", benchPath(t), "gate")
	cmd.Dir, cmd.Env, cmd.Stdout, cmd.Stderr = f.Root, surfaceEnv(f, nil), &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	deadline := time.Now().Add(5 * time.Second)
	var ownerPID int
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(ownerPath)
		if err == nil {
			fields := strings.Fields(string(data))
			if len(fields) == 2 {
				ownerPID, err = strconv.Atoi(fields[0])
				if err == nil && syscall.Kill(ownerPID, 0) == nil {
					break
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	if ownerPID == 0 {
		_ = cmd.Process.Kill()
		t.Fatalf("gate owner record did not name a live gate-run process: %s", output.String())
	}
	blocked := f.Bench("gate")
	if blocked.ExitCode == 0 || !strings.Contains(blocked.Stderr, "gate execution already in progress") || !strings.Contains(blocked.Stderr, "gate owner: pid "+strconv.Itoa(ownerPID)+" (alive)") {
		t.Fatalf("blocked gate diagnostic = exit %d stderr %q, want owner pid %d alive", blocked.ExitCode, blocked.Stderr, ownerPID)
	}
	deadPID := 99999999
	if err := os.WriteFile(ownerPath, []byte(strconv.Itoa(deadPID)+" 2026-07-23T12:00:00Z"), 0o600); err != nil {
		t.Fatal(err)
	}
	dead := f.Bench("gate")
	if dead.ExitCode == 0 || !strings.Contains(dead.Stderr, "gate execution already in progress") || !strings.Contains(dead.Stderr, "gate owner: pid "+strconv.Itoa(deadPID)+" (not alive)") {
		t.Fatalf("dead owner diagnostic = exit %d stderr %q, want owner pid %d not alive", dead.ExitCode, dead.Stderr, deadPID)
	}
	requireBareRefusal := func(label string, probe contract.Probe) {
		t.Helper()
		if probe.ExitCode == 0 || strings.TrimSpace(probe.Stderr) != "gate execution already in progress" {
			t.Fatalf("%s owner refusal = exit %d stderr %q", label, probe.ExitCode, probe.Stderr)
		}
	}
	if err := os.Remove(ownerPath); err != nil {
		t.Fatal(err)
	}
	requireBareRefusal("absent", f.Bench("gate"))
	if err := os.WriteFile(ownerPath, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	requireBareRefusal("empty", f.Bench("gate"))
	if err := os.WriteFile(ownerPath, []byte("unparseable\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	requireBareRefusal("unparseable", f.Bench("gate"))
	pgid, err := strconv.Atoi(strings.TrimSpace(string(mustReadRuntime(t, filepath.Join(gitdir, "ft88-gate-pgid")))))
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Process.Signal(signal); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
		t.Fatalf("gate-run did not exit after %s", signal)
	}
	for deadline = time.Now().Add(3 * time.Second); syscall.Kill(-pgid, 0) == nil && time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
	}
	if syscall.Kill(-pgid, 0) == nil {
		t.Fatalf("gate script process group %d survived %s", pgid, signal)
	}
	if got := gate.Inspect(f.Root); got.State != gate.Pending || got.Status == "red" {
		t.Fatalf("%s gate inspection = %+v, want pending and not red", signal, got)
	}
	if _, err := os.Stat(ownerPath); !os.IsNotExist(err) {
		t.Fatalf("owner record remains after signal: %v", err)
	}

	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("green gate")
	f.Bench("gate").RequireExit(0)
	if _, err := os.Stat(ownerPath); !os.IsNotExist(err) {
		t.Fatalf("owner record remains after green exit: %v", err)
	}
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 7\n")
	f.CommitAll("red gate")
	f.Bench("gate").RequireExit(7)
	if _, err := os.Stat(ownerPath); !os.IsNotExist(err) {
		t.Fatalf("owner record remains after red exit: %v", err)
	}
}
