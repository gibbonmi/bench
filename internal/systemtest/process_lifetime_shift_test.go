//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestProcessLifetimeCleanProtectedRecovery is PL85. The first shift CLI dies while its
// registered adapter group lives on, and its worktree holds only scratch. A second shift CLI
// runs the real recovery owner before its own acquire. The recovery must retain the worktree
// with its scratch, because an empty dirty set is no proof that the adapter group cannot
// change the worktree.
func TestProcessLifetimeCleanProtectedRecovery(t *testing.T) {
	t.Parallel()
	fixture := scaffoldRecordedPublicationRepo(t)
	commitmenttest.Register(t, fixture.root, "shift-protected")
	commitmenttest.Admit(t, fixture.root, "shift-protected", "specs/record/spec.md")
	marks := t.TempDir()
	mark, adapter := filepath.Join(marks, "adapter.pid"), filepath.Join(marks, "agent")
	// The adapter changes no worktree file, and it leaves the CLI's streams, so only its
	// group keeps the worktree in use after the CLI dies. It names its worktree, because the
	// killed CLI publishes no bounded response.
	script := "#!/bin/sh\ncat >/dev/null\npwd -P > '" + mark + ".dir'\necho $$ > '" + mark + ".partial' && mv '" + mark + ".partial' '" + mark + "'\nexec sleep 300 </dev/null >/dev/null 2>&1\n"
	if err := os.WriteFile(adapter, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	first, stdout, stderr := systemStartSelected(t, fixture.root, append(fixture.environment(fixture.home),
		"BENCH_AGENT="+adapter, "BENCH_MAX_ITERS=1", "TMPDIR="+t.TempDir()), "shift", "--outcome", "delivery", "protected scratch")
	processtest.WaitPID(t, mark)
	if err := first.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = exitWithin(t, first, 0)
	dir, err := os.ReadFile(mark + ".dir")
	if err != nil {
		t.Fatal(err)
	}
	worktree := strings.TrimSpace(string(dir))
	scratch := filepath.Join(worktree, ".bench-objective")
	if _, err := os.Stat(scratch); worktree == "" || err != nil {
		t.Fatalf("the killed shift left worktree %q with scratch %v; stdout %q stderr %q", worktree, err, stdout.String(), stderr.String())
	}

	second := systemSelected(t, fixture.root, append(fixture.environment(fixture.home),
		"BENCH_AGENT=true", "BENCH_MAX_ITERS=1", "TMPDIR="+t.TempDir()), "shift", "--outcome", "delivery", "recovering shift")
	if !strings.Contains(second.stdout, "bench shift recovery: recovered 1") {
		t.Fatalf("second shift = (%d, %q, %q), want the recovery of the killed shift", second.code, second.stdout, second.stderr)
	}
	if _, err := os.Stat(scratch); err != nil {
		t.Fatalf("the recovery removed the protected scratch: %v", err)
	}
	listing := systemGitOutput(t, fixture.root, "worktree", "list", "--porcelain")
	at := strings.Index(listing, "worktree "+worktree+"\n")
	if at < 0 {
		t.Fatalf("git lists no worktree %s after the recovery: %q", worktree, listing)
	}
	if block, _, _ := strings.Cut(listing[at:], "\n\n"); !strings.Contains(block, "\nlocked") {
		t.Fatalf("worktree %s after the recovery = %q, want it retained and locked", worktree, block)
	}
}

// TestProcessLifetimeShiftDifferential is PL102. The pass evidence records the adapter's own
// exit, and the shift result keeps its taxonomy: a zero adapter with no change is a no-op,
// and a nonzero one fails the shift. A surviving descendant gets no normal-exit signal. The
// wall deadline and a pulled line keep their checkpoint exits, stop the adapter group, and
// keep the bounded response with its shift_result block.
func TestProcessLifetimeShiftDifferential(t *testing.T) {
	t.Parallel()
	args := []string{"shift", "--outcome", "delivery", "shape objective"}
	setup := func(t *testing.T, adapter string, extra ...string) (recordedPublicationRepo, []string) {
		fixture := scaffoldRecordedPublicationRepo(t)
		commitmenttest.Register(t, fixture.root, "shift-shape")
		commitmenttest.Admit(t, fixture.root, "shift-shape", "specs/record/spec.md")
		env := append(fixture.environment(fixture.home), "BENCH_AGENT="+adapter, "BENCH_MAX_ITERS=1", "TMPDIR="+t.TempDir())
		return fixture, append(env, extra...)
	}
	shift := func(t *testing.T, adapter string, extra ...string) (recordedPublicationRepo, processResult) {
		fixture, env := setup(t, adapter, extra...)
		return fixture, systemSelected(t, fixture.root, env, args...)
	}
	// A held adapter publishes its group and waits, so only a checkpoint exit ends the shift.
	hold := func(t *testing.T) (adapter, mark string) {
		dir := t.TempDir()
		mark = filepath.Join(dir, "group")
		return callerTool(t, dir, "agent", "cat >/dev/null\n"+publishGroup("mv", mark)+"exec sleep 300\n"), mark
	}
	// A checkpoint exit keeps its code and the bounded response that names its outcome.
	checkpointExit := func(t *testing.T, path string, code int, stdout, stderr string, want int, outcome string) {
		t.Helper()
		if code != want || !strings.Contains(stdout, "shift_result[1]") || !strings.Contains(stdout, outcome) {
			t.Fatalf("shift %s = (%d, %q, %q), want exit %d with the shift_result outcome %q", path, code, stdout, stderr, want, outcome)
		}
	}
	for _, shape := range preservationShapes {
		t.Run(shape.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			mark := filepath.Join(dir, "group")
			fixture, result := shift(t, callerTool(t, dir, "agent", "cat >/dev/null\n"+shapeScript(shape, "sleep", "mv", mark)))
			pgid := groupOf(t, mark)
			if shape.descendant {
				requireGroup(t, pgid, protection.GroupPresent, "shift adapter")
			}
			// The record reader keeps string values only, so the integer exit is read raw.
			files, _ := filepath.Glob(filepath.Join(fixture.home, "otel", "*", "traces.jsonl"))
			raw := []byte{}
			if len(files) == 1 {
				raw, _ = os.ReadFile(files[0])
			}
			exit := `{"key":"` + otelrecord.AttrAdapterExit + `","value":{"intValue":"` + strconv.Itoa(shape.exit) + `"}}`
			ran := `{"key":"` + otelrecord.AttrAdapterResult + `","value":{"stringValue":"` + otelrecord.AdapterExited + `"}}`
			if !strings.Contains(string(raw), exit) || !strings.Contains(string(raw), ran) {
				t.Fatalf("shift record %s = %q, want the pass evidence %s and %s", files, raw, ran, exit)
			}
			// A green gate with no change is a no-op, unless the adapter exited nonzero.
			code, outcome := 4, "\n  no-op,"
			if shape.exit != 0 {
				code, outcome = 1, "\n  failed,"
			}
			if result.code != code || !strings.Contains(result.stdout, "shift_result[1]") || !strings.Contains(result.stdout, outcome) {
				t.Fatalf("shift over %s = (%d, %q, %q), want exit %d and outcome %q", shape.name, result.code, result.stdout, result.stderr, code, outcome)
			}
		})
	}
	t.Run("wall deadline", func(t *testing.T) {
		t.Parallel()
		adapter, mark := hold(t)
		_, result := shift(t, adapter, "BENCH_MAX_WALL=1s")
		pgid := groupOf(t, mark)
		checkpointExit(t, "past its wall deadline", result.code, result.stdout, result.stderr, 3, "\n  incomplete,")
		requireGroup(t, pgid, protection.GroupAbsent, "shift adapter past its deadline")
	})
	t.Run("interrupt", func(t *testing.T) {
		t.Parallel()
		adapter, mark := hold(t)
		fixture, env := setup(t, adapter)
		cmd, stdout, stderr := systemStartSelected(t, fixture.root, env, args...)
		pgid := groupOf(t, mark)
		if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
			t.Fatal(err)
		}
		code := exitWithin(t, cmd, 0)
		checkpointExit(t, "after a pulled line", code, stdout.String(), stderr.String(), 130, "\n  interrupted,")
		requireGroup(t, pgid, protection.GroupAbsent, "shift adapter after a pulled line")
	})
}
