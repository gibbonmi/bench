package contract

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestFixtureInitializesGitRepoByDefault(t *testing.T) {
	t.Parallel()
	f := NewFixture(t)

	probe := f.Run("git", "rev-parse", "--show-toplevel")

	probe.RequireExit(0)
	if got := strings.TrimSpace(probe.Stdout); got != f.Root {
		t.Fatalf("git root = %q, want %q", got, f.Root)
	}
}

func TestFixtureOptions(t *testing.T) {
	t.Parallel()
	noRepo := NewFixture(t, WithNoRepo())
	noRepo.GitAllow("rev-parse", "--show-toplevel").RequireExit(128)

	spacePath := NewFixture(t, WithSpacePath())
	if !strings.Contains(spacePath.Root, "space dir") {
		t.Fatalf("space-path fixture root = %q, want a path containing space dir", spacePath.Root)
	}
	spacePath.Git("rev-parse", "--show-toplevel").RequireExit(0)
}

func TestRunCapturesOutputExitAndEnvironment(t *testing.T) {
	t.Setenv("HOME", filepath.Join(t.TempDir(), "ambient-home"))
	t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "ambient-bench-home"))
	f := NewFixture(t)

	probe := f.RunEnv(map[string]string{"CONTRACT_VALUE": "ok"}, "sh", "-c", `
printf 'stdout:%s:%s\n' "$CONTRACT_VALUE" "$HOME"
printf 'stderr:%s\n' "$BENCH_HOME" >&2
exit 7
`)

	probe.RequireExit(7)
	probe.RequireContains(probe.Stdout, "stdout:ok:"+filepath.Join(f.Root, ".bench-contract-env", "home"))
	probe.RequireContains(probe.Stderr, "stderr:"+filepath.Join(f.Root, ".bench-contract-env", "bench-home"))
	probe.RequireNotContains(probe.Stdout+probe.Stderr, "ambient-")
}

func TestFixtureRunnerReapsSpawnedProcessGroup(t *testing.T) {
	f := NewFixture(t)
	pids := filepath.Join(f.Root, "runner-pids")

	f.Run("sh", "-c", "sleep 30 >/dev/null 2>&1 & printf '%s %s\\n' \"$$\" \"$!\" > \"$1\"", "sh", pids).RequireExit(0)

	data, err := os.ReadFile(pids)
	if err != nil {
		t.Fatal(err)
	}
	fields := strings.Fields(string(data))
	if len(fields) != 2 {
		t.Fatalf("runner PID record = %q, want group and child PID", data)
	}
	group, err := strconv.Atoi(fields[0])
	if err != nil {
		t.Fatalf("parse group PID %q: %v", fields[0], err)
	}
	child, err := strconv.Atoi(fields[1])
	if err != nil {
		t.Fatalf("parse child PID %q: %v", fields[1], err)
	}

	RequireProcessGroupDrained(t, group, 2*time.Second, "runner returned before reaping what it spawned", child)
}

func TestBenchRunsKitWrapperFromFixture(t *testing.T) {
	t.Parallel()
	f := NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nprintf 'fixture gate cwd=%s\\n' \"$PWD\"\nexit 0\n")
	f.CommitAll("init")

	probe := f.Bench("gate")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "fixture gate cwd="+f.Root)
}
