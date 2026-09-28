package env

import (
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestKitTestRunDisablesVerdictWindows(t *testing.T) {
	run, err := OpenKitTestRun(kitRunBase(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := run.Close(); err != nil {
			t.Error(err)
		}
	})
	if got := kitEnvValue(run.Entries(), bounds.UnboundedWaitsEnv); got != "1" {
		t.Fatalf("kit verdict-window switch = %q, want 1", got)
	}
}
