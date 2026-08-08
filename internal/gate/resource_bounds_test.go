package gate

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestExecuteDeadlineRecordsDistinctTimeout(t *testing.T) {
	root := gateTestRepo(t, "#!/bin/sh\nsleep 30 &\necho $! > .git/timeout-child\nwait\n", "")
	old := gateTimeout
	gateTimeout = 50 * time.Millisecond
	t.Cleanup(func() { gateTimeout = old })
	parent, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	var stderr bytes.Buffer
	got := Execute(parent, root, io.Discard, &stderr)
	inspection := Inspect(root)
	if got.ActionExit != 124 || got.GateExit != 124 || inspection.State != Ready || inspection.Status != "timeout" || inspection.ReusableGreen {
		t.Fatalf("result=%+v inspection=%+v stderr=%q", got, inspection, stderr.String())
	}
	if !strings.Contains(stderr.String(), "gate: timeout") {
		t.Fatalf("timeout message missing: %q", stderr.String())
	}
	child := waitForPIDFile(t, filepath.Join(root, ".git", "timeout-child"))
	t.Cleanup(func() { _ = syscall.Kill(child, syscall.SIGKILL) })
	waitForProcessExit(t, child)
}

// TestGateRunDeadlineTermGraceThenKill grades the kill cascade on the deadline path.
// A deadline fires with no operator present, so the gate's own dying words are the
// only account of what the run was stuck on; a teardown that opens with SIGKILL takes
// that account with it. The stubborn half pins the opposite edge — the grace is a
// bounded courtesy, not a licence for a gate that ignores TERM to outlive its deadline.
func TestGateRunDeadlineTermGraceThenKill(t *testing.T) {
	const trapMarker = "gate-term-trap-ran"
	for _, tc := range []struct {
		name    string
		script  string
		reports bool
	}{
		{
			name:    "trapping gate speaks before it dies",
			script:  "#!/usr/bin/env bash\ntrap 'echo " + trapMarker + "; exit 143' TERM\nsleep 30 &\nwait\n",
			reports: true,
		},
		{
			name:   "gate ignoring TERM is killed after the grace",
			script: "#!/usr/bin/env bash\ntrap '' TERM\nwhile :; do sleep 0.2; done\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := gateTestRepo(t, tc.script, "")
			old := gateTimeout
			gateTimeout = 50 * time.Millisecond
			t.Cleanup(func() { gateTimeout = old })

			var stdout, stderr bytes.Buffer
			done := make(chan Result, 1)
			ctx := withProcessGroupCancelGrace(context.Background(), fastProcessGroupCancelGrace)
			go func() { done <- Execute(ctx, root, &stdout, &stderr) }()
			var got Result
			select {
			case got = <-done:
			case <-time.After(processGroupCancelGrace + 30*time.Second):
				t.Fatal("deadline path never returned; the grace did not escalate to SIGKILL")
			}
			if got.ActionExit != 124 || got.GateExit != 124 {
				t.Fatalf("deadline result = %+v, want exit 124; stdout=%q stderr=%q", got, stdout.String(), stderr.String())
			}
			if tc.reports && !strings.Contains(stdout.String(), trapMarker) {
				t.Fatalf("gate's TERM trap never ran, so nothing it had to say survived the teardown; stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
		})
	}
}

func TestExecuteHealthyGateJustBelowDeadlineRecordsOrdinaryGreen(t *testing.T) {
	root := gateTestRepo(t, "#!/bin/sh\nsleep 0.6\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
	old := gateTimeout
	gateTimeout = time.Second
	t.Cleanup(func() { gateTimeout = old })
	var stderr bytes.Buffer
	got := Execute(context.Background(), root, io.Discard, &stderr)
	inspection := Inspect(root)
	if got.ActionExit != 0 || got.GateExit != 0 || inspection.State != Ready || inspection.Status != "green" || !inspection.ReusableGreen {
		t.Fatalf("near-boundary result=%+v inspection=%+v stderr=%q", got, inspection, stderr.String())
	}
}

// TestManifestEntryLimitConstant pins the production entry-limit value. The runtime
// boundary proof runs at a lowered limit (BENCH_GATE_ENTRY_LIMIT) for speed, so this is
// the independent expectation that a change to the real ceiling turns red — without it,
// the cheap proof could pass while the shipped limit silently drifted.
func TestManifestEntryLimitConstant(t *testing.T) {
	t.Parallel()
	if defaultManifestEntryLimit != 100000 {
		t.Fatalf("defaultManifestEntryLimit = %d, want 100000 (the shipped gate identity ceiling)", defaultManifestEntryLimit)
	}
}

// TestManifestEntryLimitOverrideIsTightenOnly pins the fail-safe: BENCH_GATE_ENTRY_LIMIT
// may only lower the ceiling. A value at or above the default, or a malformed one, is
// ignored — so the override can never raise the limit and can never enable a false green.
func TestManifestEntryLimitOverrideIsTightenOnly(t *testing.T) {
	cases := []struct {
		value string
		want  int
	}{
		{"", defaultManifestEntryLimit},
		{"10", 10},
		{"0", 0},
		{"99999", 99999},
		{"100000", defaultManifestEntryLimit},   // equal — ignored, cannot raise
		{"100001", defaultManifestEntryLimit},   // above — ignored
		{"-1", defaultManifestEntryLimit},       // negative — ignored
		{"nonsense", defaultManifestEntryLimit}, // malformed — ignored
	}
	for _, tc := range cases {
		t.Setenv("BENCH_GATE_ENTRY_LIMIT", tc.value)
		if tc.value == "" {
			os.Unsetenv("BENCH_GATE_ENTRY_LIMIT")
		}
		if got := manifestEntryLimit(); got != tc.want {
			t.Errorf("manifestEntryLimit() with %q = %d, want %d", tc.value, got, tc.want)
		}
	}
}
