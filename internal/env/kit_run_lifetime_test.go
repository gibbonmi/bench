package env

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestKitRunClose is PL37. A child registers on the vector that Entries hands it, writes into
// the private home, and leaves a live descendant in its group. Close must keep that file.
func TestKitRunClose(t *testing.T) {
	p := &protection.Protocol{Validators: protection.Validators{protection.KindKitTest: ValidateProtection}, ReadLimit: bounds.ControlRecordLimit}
	// The run under test owns the only vector, so an outer run's resources stay out of it.
	run, err := openKitTestRun(protection.Environ(kitRunBase(t), nil), p)
	if err != nil {
		t.Fatal(err)
	}
	vector, err := protection.ParseVector(kitEnvValue(run.Entries(), protection.Env))
	if err != nil {
		t.Fatal(err)
	}
	pidFile := filepath.Join(t.TempDir(), "descendant")
	user := exec.Command("sh", "-c", `echo held >"$HOME/owned" || exit 1; sleep 300 </dev/null >/dev/null 2>&1 & echo $! >"$1.partial" && mv "$1.partial" "$1"`, "sh", pidFile)
	user.Env = append(os.Environ(), run.Entries()...)
	policy := subprocess.Policy{Normal: subprocess.Preserve, Final: bounds.ProcessWindow(bounds.FixedWindow(bounds.ProcessFinalWindow)), Poll: bounds.ProcessPoll}
	lifetime := subprocess.NewLifetime(user, policy, subprocess.Streams{}, subprocess.Registration{Protocol: p, Resources: vector})
	if err := lifetime.Start(); err != nil {
		t.Fatal(err)
	}
	outcome, _ := lifetime.Wait()
	processtest.WaitPID(t, pidFile)
	if outcome.Process.Code != 0 || outcome.Protection.State != protection.Protected {
		t.Fatalf("user outcome = %+v, want exit 0 with its live group protected", outcome)
	}
	closeErr := run.Close()
	if closeErr == nil || !strings.Contains(closeErr.Error(), string(protection.Protected)) {
		t.Fatalf("Close with a live registered group = %v, want a protected retention", closeErr)
	}
	if data, err := os.ReadFile(filepath.Join(kitEnvValue(run.Entries(), "HOME"), "owned")); err != nil || string(data) != "held\n" {
		t.Fatalf("child-owned file after Close = %q, %v; want it retained", data, err)
	}
}
