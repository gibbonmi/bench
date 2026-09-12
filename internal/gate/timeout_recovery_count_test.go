package gate

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGateRunTimeoutInvalidatesOldEvidence(t *testing.T) {
	for _, test := range []struct {
		name    string
		fixture func(*testing.T) string
	}{
		{name: "normal fixture", fixture: failureOutcomeFixture},
		{name: "delayed counter fixture", fixture: delayedCounterOutcomeFixture},
	} {
		t.Run(test.name, func(t *testing.T) {
			previousTimeout := gateTimeout
			t.Cleanup(func() { gateTimeout = previousTimeout })
			root := test.fixture(t)
			if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
				t.Fatalf("green result = %#v", result)
			}

			gateTimeout = 100 * time.Millisecond
			outcomeWrite(t, root, ".gate-sleep", "\n", 0o644)
			var stdout, stderr bytes.Buffer
			if got := RunCommand([]string{"--fresh", root}, &stdout, &stderr); got != 124 {
				t.Fatalf("timeout exit = %d, stderr=%q", got, stderr.String())
			}
			if !strings.Contains(stderr.String(), "gate: timeout") {
				t.Fatalf("stderr = %q", stderr.String())
			}
			if inspection := Inspect(root); inspection.State != Ready || inspection.Status != "timeout" || inspection.ReusableGreen {
				t.Fatalf("timeout inspection = %#v", inspection)
			}

			runsAfterTimeout := outcomeRuns(t, root)
			if err := os.Remove(filepath.Join(root, ".gate-sleep")); err != nil {
				t.Fatal(err)
			}
			gateTimeout = previousTimeout
			if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
				t.Fatalf("ordinary run after timeout = %#v", result)
			}
			if got := outcomeRuns(t, root); got != runsAfterTimeout+1 {
				t.Fatalf("runs after recovery = %d, want %d", got, runsAfterTimeout+1)
			}
		})
	}
}

func delayedCounterOutcomeFixture(t *testing.T) string {
	t.Helper()
	root := outcomeFixture(t)
	path := filepath.Join(root, ".bench/gate.sh")
	script := string(outcomeRead(t, path))
	delayed := strings.Replace(script, "count=0\n", "if [ -e .gate-sleep ]; then sleep 5; fi\ncount=0\n", 1)
	if delayed == script {
		t.Fatal("gate fixture has no counter initialization")
	}
	outcomeWrite(t, root, ".bench/gate.sh", delayed, 0o755)
	outcomeCommit(t, root, "delay counter")
	return root
}
