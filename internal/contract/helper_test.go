package contract

import (
	"path/filepath"
	"strings"
	"testing"
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

func TestBenchRunsKitWrapperFromFixture(t *testing.T) {
	t.Parallel()
	f := NewFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nprintf 'fixture gate cwd=%s\\n' \"$PWD\"\nexit 0\n")
	f.CommitAll("init")

	probe := f.Bench("gate")

	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "fixture gate cwd="+f.Root)
}
