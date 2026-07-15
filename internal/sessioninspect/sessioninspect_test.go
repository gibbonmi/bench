package sessioninspect

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

func TestInspectDeadlineWarnsAndReturnsZero(t *testing.T) {
	original := phases
	t.Cleanup(func() { phases = original })
	phases = []phase{func(ctx context.Context, _, _ io.Writer, _ string) int {
		<-ctx.Done()
		return 1
	}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	var out bytes.Buffer
	if code := Inspect(ctx, &out, t.TempDir()); code != 0 {
		t.Fatalf("Inspect exit = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "warning: bench session-inspect: deadline exceeded; session inspection stopped") {
		t.Fatalf("Inspect warning missing:\n%s", out.String())
	}
}

func TestCommandInstallsTenSecondDeadline(t *testing.T) {
	original := runInspect
	t.Cleanup(func() { runInspect = original })
	runInspect = func(ctx context.Context, _ io.Writer, _ string) int {
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("Command handed Inspect a context without a deadline")
		}
		remaining := time.Until(deadline)
		if remaining < 9*time.Second || remaining > 10*time.Second {
			t.Fatalf("deadline remaining = %v, want concrete 10s timeout", remaining)
		}
		return 0
	}
	if code := Command(nil, io.Discard, io.Discard); code != 0 {
		t.Fatalf("Command exit = %d, want 0", code)
	}
}
