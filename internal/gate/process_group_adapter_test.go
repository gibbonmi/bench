package gate

import (
	"context"
	"os/exec"
	"testing"
)

func TestProcessGroupAdapterRunsOneControlledCommand(t *testing.T) {
	result := runProcessGroupCommand(context.Background(), exec.Command("sh", "-c", "exit 0"))
	if result.Code != 0 || result.StartErr != nil || result.Cancelled {
		t.Fatalf("runProcessGroupCommand() = %#v, want a completed green process", result)
	}
}
