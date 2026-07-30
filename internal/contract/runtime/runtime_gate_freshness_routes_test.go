package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestFreshnessRefusalReachesEveryGateRoute(t *testing.T) {
	t.Run("direct kit gate", func(t *testing.T) {
		f := freshnessRouteFixture(t)
		probe := contract.RunAt(t, f, f.Root, nil, "bash", filepath.Join(f.Root, ".bench", "gate.sh"))
		requireFreshnessRouteRefusal(t, "direct kit gate", probe.Stdout+probe.Stderr)
		if probe.ExitCode != 1 {
			t.Fatalf("direct kit gate exit = %d, want 1", probe.ExitCode)
		}
	})

	t.Run("linked repository by-path gate from nested cwd", func(t *testing.T) {
		f := freshnessRouteFixture(t)
		nested := filepath.Join(f.Root, "nested [*] path")
		contract.Mkdir(t, nested)
		probe := contract.RunAt(t, f, nested, nil, "bash", filepath.Join(f.Root, ".bench", "gate.sh"))
		requireFreshnessRouteRefusal(t, "linked repository by-path gate", probe.Stdout+probe.Stderr)
		if probe.ExitCode != 1 {
			t.Fatalf("linked repository by-path gate exit = %d, want 1", probe.ExitCode)
		}
	})

	t.Run("armed Stop hook", func(t *testing.T) {
		f := freshnessRouteFixture(t)
		probe := contract.RunAtWithInput(t, f, f.Root, map[string]string{"BENCH_SHIFT": "1"}, "{}\n", "bash", filepath.Join(f.Root, ".bench", "hooks", "stop.sh"))
		requireFreshnessRouteRefusal(t, "armed Stop hook", probe.Stdout+probe.Stderr)
		if probe.ExitCode != 2 {
			t.Fatalf("armed Stop hook exit = %d, want 2", probe.ExitCode)
		}
	})

	t.Run("shift through configured Codex adapter", func(t *testing.T) {
		f := freshnessRouteFixture(t)
		home := t.TempDir()
		addRuntimePoolWorktrees(t, f, home)
		probe := f.BenchEnv(map[string]string{
			"BENCH_AGENT":     filepath.Join(f.Root, ".bench", "adapters", "codex"),
			"BENCH_HOME":      home,
			"BENCH_MAX_ITERS": "1",
			"PATH":            filepath.Join(f.Root, "tools") + string(os.PathListSeparator) + os.Getenv("PATH"),
		}, "shift", "FT131 freshness route")
		output := probe.Stdout + probe.Stderr
		if !strings.Contains(output, "ft131 adapter invoked") {
			t.Fatalf("configured adapter did not run:\n%s", output)
		}
		requireFreshnessRouteRefusal(t, "shift through configured adapter", output)
		if probe.ExitCode != 1 {
			t.Fatalf("shift exit = %d, want 1\n%s", probe.ExitCode, output)
		}
	})
}

func TestGateReplacementRunsTheRebuiltPhaseTableOnce(t *testing.T) {
	root := contract.SubjectRoot(t)
	f := contract.NewExecFixtureAt(t, t.TempDir())
	clone := filepath.Join(f.Root, "replacement [*] source")
	f.Run("git", "clone", "-q", "--no-hardlinks", root, clone).RequireExit(0)
	copyRuntimeFile(t, filepath.Join(root, ".bench", "gate.sh"), filepath.Join(clone, ".bench", "gate.sh"), 0o755)

	build := contract.RunAt(t, f, clone, nil, "bash", filepath.Join(clone, "scripts", "go-build.sh"), clone, filepath.Join(clone, "dist", "bench"))
	if build.ExitCode != 0 {
		t.Fatalf("initial prescribed build exit = %d:\n%s%s", build.ExitCode, build.Stdout, build.Stderr)
	}
	phasePath := filepath.Join(clone, "internal", "gate", "phases.go")
	phaseSource, err := os.ReadFile(phasePath)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(phaseSource), `"shellcheck"`, `"fresh-shellcheck"`, 1)
	updated = strings.Replace(updated, `shellcheckArgv(kit)`, `[]string{"bash", "-c", "printf 'fresh-shellcheck\\n'"}`, 1)
	if updated == string(phaseSource) || !strings.Contains(updated, `"fresh-shellcheck"`) || !strings.Contains(updated, `printf 'fresh-shellcheck\\n'`) {
		t.Fatal("phase table did not contain the shellcheck phase to replace")
	}
	if err := os.WriteFile(phasePath, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}

	stale := contract.RunAt(t, f, clone, map[string]string{"BENCH_CANARY_INNER": "1", "BENCH_CANARY_PHASE": "fresh-shellcheck"}, "bash", filepath.Join(clone, ".bench", "gate.sh"))
	requireFreshnessRouteRefusal(t, "stale replacement table", stale.Stdout+stale.Stderr)
	if stale.ExitCode != 1 {
		t.Fatalf("stale gate exit = %d, want 1", stale.ExitCode)
	}

	rebuild := contract.RunAt(t, f, clone, nil, "bash", filepath.Join(clone, "scripts", "go-build.sh"), clone, filepath.Join(clone, "dist", "bench"))
	if rebuild.ExitCode != 0 {
		t.Fatalf("replacement prescribed build exit = %d:\n%s%s", rebuild.ExitCode, rebuild.Stdout, rebuild.Stderr)
	}
	fresh := contract.RunAt(t, f, clone, map[string]string{"BENCH_CANARY_INNER": "1", "BENCH_CANARY_PHASE": "fresh-shellcheck"}, "bash", filepath.Join(clone, ".bench", "gate.sh"))
	if fresh.ExitCode != 0 {
		t.Fatalf("rebuilt gate exit = %d:\n%s%s", fresh.ExitCode, fresh.Stdout, fresh.Stderr)
	}
	output := fresh.Stdout + fresh.Stderr
	if strings.Count(output, "fresh-shellcheck\n") != 1 || strings.Contains(output, "shellcheck: skipped") {
		t.Fatalf("rebuilt phase output did not resolve the replacement table exactly once:\n%s", output)
	}
}

func freshnessRouteFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := copiedCLIHookFixture(t, true)
	root := contract.SubjectRoot(t)
	for _, rel := range []string{
		".bench/gate.sh",
		".bench/hooks/stop.sh",
		".bench/lib/resolve-bench.sh",
		".bench/adapters/codex",
	} {
		copyRuntimeFile(t, filepath.Join(root, rel), filepath.Join(f.Root, rel), 0o755)
	}
	f.WriteExecutable("tools/codex", "#!/usr/bin/env bash\nprintf 'ft131 adapter invoked\\n' >&2\n")
	f.CommitAll("FT131 freshness route fixture")
	return f
}

func requireFreshnessRouteRefusal(t *testing.T, route, output string) {
	t.Helper()
	if !strings.Contains(output, "bench binary ") || strings.Count(output, "rebuild with bash scripts/go-build.sh ") != 1 {
		t.Fatalf("%s did not report the stable freshness rebuild action:\n%s", route, output)
	}
	if strings.Contains(output, "phase ") || strings.Contains(output, "old phase") {
		t.Fatalf("%s resolved or scheduled phases before freshness refusal:\n%s", route, output)
	}
}
