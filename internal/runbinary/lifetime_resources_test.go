package runbinary

import (
	"context"
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

// binaryProtocol admits the selection kind alone, as the composition root's binary row does.
func binaryProtocol() *protection.Protocol {
	return &protection.Protocol{Validators: protection.Validators{protection.KindBinary: ValidateProtection}, ReadLimit: bounds.ControlRecordLimit}
}

// TestBinaryClose is PL36. A selection user exits and leaves a live descendant in its
// registered group, so Close must keep the executable and report the retention.
func TestBinaryClose(t *testing.T) {
	t.Parallel()
	factory := Factory{
		TempRoot:   t.TempDir(),
		Protection: binaryProtocol(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
	selection, err := factory.Own(context.Background(), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if resources := selection.Resources(); len(resources) != 1 || ValidateProtection(resources[0]) != nil {
		t.Fatalf("selection resources = %+v, want one valid binary generation", resources)
	}
	pidFile := filepath.Join(t.TempDir(), "descendant")
	user := exec.Command("sh", "-c", `sleep 300 </dev/null >/dev/null 2>&1 & echo $! > "$1.partial" && mv "$1.partial" "$1"`, "sh", pidFile)
	policy := subprocess.Policy{Normal: subprocess.Preserve, Final: bounds.ProcessWindow(bounds.FixedWindow(bounds.ProcessFinalWindow)), Poll: bounds.ProcessPoll}
	lifetime := subprocess.NewLifetime(user, policy, subprocess.Streams{}, subprocess.Registration{Protocol: factory.Protection, Resources: selection.Resources()})
	if err := lifetime.Start(); err != nil {
		t.Fatal(err)
	}
	outcome, _ := lifetime.Wait()
	processtest.WaitPID(t, pidFile)
	if outcome.Process.Code != 0 || outcome.Protection.State != protection.Protected {
		t.Fatalf("user outcome = %+v, want exit 0 with its live group protected", outcome)
	}
	closeErr := selection.Close()
	if closeErr == nil || !strings.Contains(closeErr.Error(), string(protection.Protected)) {
		t.Fatalf("Close with a live registered group = %v, want a protected retention", closeErr)
	}
	if info, err := os.Stat(selection.Path); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("retained executable = %v, %v; want the selected file", info, err)
	}
}
