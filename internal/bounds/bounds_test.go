package bounds

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestProductionPolicyValues(t *testing.T) {
	if ProviderTimeout != 10*time.Second || GitRefreshTimeout != 30*time.Second || GuardScanTimeout != 5*time.Second || GateTimeout != 45*time.Minute {
		t.Fatalf("duration policy changed: provider=%s refresh=%s guard=%s gate=%s", ProviderTimeout, GitRefreshTimeout, GuardScanTimeout, GateTimeout)
	}
	if ModelReadLimit != 5<<20 || OutlineFileLimit != 2<<20 || OutlineRowLimit != 200 {
		t.Fatalf("read/output policy changed: model=%d outline_file=%d outline_rows=%d", ModelReadLimit, OutlineFileLimit, OutlineRowLimit)
	}
	if IterationMin != 1 || IterationMax != 100 || MainIterationsDefault != 12 || RefactorIterationsDefault != 4 || MaxWall != 24*time.Hour {
		t.Fatalf("shift policy changed: range=[%d,%d] defaults=%d/%d max_wall=%s", IterationMin, IterationMax, MainIterationsDefault, RefactorIterationsDefault, MaxWall)
	}
}

func TestRunClassifiesProcessOutcomes(t *testing.T) {
	tests := []struct {
		name   string
		parent func() (context.Context, context.CancelFunc)
		limit  time.Duration
		cmd    func(context.Context) *exec.Cmd
		want   ProcessStatus
		exit   int
	}{
		{name: "complete", parent: liveContext, limit: time.Second, cmd: func(ctx context.Context) *exec.Cmd { return exec.CommandContext(ctx, "sh", "-c", "printf ok") }, want: ProcessComplete},
		{name: "timeout", parent: liveContext, limit: 20 * time.Millisecond, cmd: func(ctx context.Context) *exec.Cmd { return exec.CommandContext(ctx, "sh", "-c", "sleep 5") }, want: ProcessTimeout},
		{name: "parent cancellation", parent: canceledContext, limit: time.Second, cmd: func(ctx context.Context) *exec.Cmd { return exec.CommandContext(ctx, "sh", "-c", "sleep 5") }, want: ProcessCanceled},
		{name: "nonzero exit", parent: liveContext, limit: time.Second, cmd: func(ctx context.Context) *exec.Cmd {
			return exec.CommandContext(ctx, "sh", "-c", "printf bad >&2; exit 23")
		}, want: ProcessExit, exit: 23},
		{name: "start failure", parent: liveContext, limit: time.Second, cmd: func(ctx context.Context) *exec.Cmd {
			return exec.CommandContext(ctx, "/definitely/missing/bench-command")
		}, want: ProcessStart},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.parent()
			defer cancel()
			got := Run(ctx, tt.limit, tt.cmd(ctx))
			if got.Status != tt.want || got.Exit != tt.exit {
				t.Fatalf("Run status/exit = %q/%d, want %q/%d (err=%v output=%q)", got.Status, got.Exit, tt.want, tt.exit, got.Err, got.Output)
			}
		})
	}
}

func TestReadClassifiesExactLimitAndLimitPlusOne(t *testing.T) {
	for _, tt := range []struct {
		name string
		body string
		want ReadStatus
	}{
		{name: "exact", body: "12345", want: ReadComplete},
		{name: "plus one", body: "123456", want: ReadOversized},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Read(strings.NewReader(tt.body), 5)
			if got.Status != tt.want {
				t.Fatalf("Read status = %q, want %q", got.Status, tt.want)
			}
			if got.Status == ReadComplete && string(got.Data) != tt.body {
				t.Fatalf("Read data = %q, want %q", got.Data, tt.body)
			}
		})
	}
	failing := Read(errorReader{}, 5)
	if failing.Status != ReadFailed || failing.Err == nil {
		t.Fatalf("failing reader = %#v, want failed with error", failing)
	}
}

func liveContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
func canceledContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx, cancel
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }
