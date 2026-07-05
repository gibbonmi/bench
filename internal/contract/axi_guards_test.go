package contract

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestAXIGuardsContracts(t *testing.T) {
	t.Parallel()
	skipIfSubjectBenchMissing(t)
	runParallel(t, "AXI guards aggregation contract", testAXIGuardsAggregation)
	runParallel(t, "AXI guards --brief contract", testAXIGuardsBrief)
	runParallel(t, "AXI guards usage/subdirectory contract", testAXIGuardsUsageSubdirectory)
	runParallel(t, "AXI guards path-with-spaces contract", testAXIGuardsPathWithSpaces)
	runParallel(t, "AXI guards --describe timeout-bound contract", testAXIGuardsDescribeTimeoutBound)
	runParallel(t, "AXI guards unmanaged-pre-push safety contract", testAXIGuardsUnmanagedPrePushSafety)
	runParallel(t, "AXI block-dangerous-git core-unreachable manifest contract", testAXIBlockDangerousGitCoreUnreachableManifest)
	runParallel(t, "AXI block-dangerous-git linked-worktree classification contract", testAXIBlockDangerousGitLinkedWorktreeClassification)
	runParallel(t, "session-start guard-brief injection contract", testSessionStartGuardBriefInjection)
	runParallel(t, "session-start never-blocks-outside-repo contract", testSessionStartNeverBlocksOutsideRepo)
}

func testAXIGuardsAggregation(t *testing.T) {
	noteContractFailure(t, "AXI guards aggregation contract failed")
	f := linkedGuardsFixture(t)

	out := f.Bench("guards")
	out.RequireExit(0)
	requireGuardsFirstLine(t, out.Stdout, "guards[4]{guard,boundary,denies}:")
	for _, guard := range []string{"block-dangerous-git", "check-agent-line", "stop", "pre-push"} {
		requireGuardsLineMatching(t, out.Stdout, "^  "+regexp.QuoteMeta(guard)+",")
	}
	requireNoGuardsLineMatching(t, out.Stdout, `^  session-start,`)

	prepush := gitPrePushPath(t, f)
	manifest := f.Run("bash", prepush, "--describe")
	manifest.RequireExit(0)
	for _, key := range []string{"name", "boundary", "denies", "why"} {
		requireGuardsLineMatching(t, manifest.Stdout, "^"+key+": ")
	}

	f.WriteExecutable(".bench/hooks/extra.sh", "#!/usr/bin/env bash\ncat >/dev/null\nexit 0\n")
	withExtra := f.Bench("guards")
	withExtra.RequireExit(0)
	requireGuardsLine(t, withExtra.Stdout, `  extra,"",no manifest`)

	if err := os.Remove(prepush); err != nil {
		t.Fatalf("remove generated pre-push: %v", err)
	}
	withoutPrePush := f.Bench("guards")
	withoutPrePush.RequireExit(0)
	requireGuardsLine(t, withoutPrePush.Stdout, `  pre-push,"",not installed`)
}

func testAXIGuardsBrief(t *testing.T) {
	f := linkedGuardsFixture(t)

	out := f.Bench("guards", "--brief")

	out.RequireExit(0)
	requireGuardsStringEqual(t, fmt.Sprint(strings.Count(out.Stdout, "full manifests: bench guards")), "1", "brief footer count")
	requireGuardsStringEqual(t, fmt.Sprint(nonEmptyGuardsLineCount(out.Stdout)), "5", "brief line count")
	out.RequireContains(out.Stdout, "block-dangerous-git: destructive git")
}

func testAXIGuardsUsageSubdirectory(t *testing.T) {
	f := linkedGuardsFixture(t)

	usage := f.Bench("guards", "bogusarg")
	usage.RequireExit(2)
	usage.RequireContains(strings.ToLower(usage.Stdout), "usage")

	subdir := filepath.Join(f.Root, "sub", "deeper")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatalf("create subdir: %v", err)
	}
	subFixture := Fixture{t: t, Root: subdir, Env: f.Env}
	fromSubdir := subFixture.Bench("guards")
	fromSubdir.RequireExit(0)
	if !strings.HasPrefix(firstGuardsLine(fromSubdir.Stdout), "guards[") {
		t.Fatalf("guards from subdirectory lost root resolution:\n%s", fromSubdir.Stdout)
	}
}

func testAXIGuardsPathWithSpaces(t *testing.T) {
	f := linkedGuardsFixture(t, WithSpacePath())

	out := f.Bench("guards")

	out.RequireExit(0)
	if !strings.HasPrefix(firstGuardsLine(out.Stdout), "guards[") {
		t.Fatalf("space-path guards failed:\n%s", out.Stdout)
	}
}

func testAXIGuardsDescribeTimeoutBound(t *testing.T) {
	f := NewFixture(t)
	f.WriteExecutable(".bench/hooks/slow.sh", "#!/usr/bin/env bash\nif [ \"${1:-}\" = \"--describe\" ]; then sleep 30; fi\nexit 0\n")

	start := time.Now()
	out := f.Bench("guards")
	elapsed := time.Since(start)

	out.RequireExit(0)
	if elapsed >= 10*time.Second {
		t.Fatalf("guards did not bound a slow --describe (took %v)", elapsed)
	}
	requireGuardsLine(t, out.Stdout, `  slow,"",no manifest (timed out)`)
}

func testAXIGuardsUnmanagedPrePushSafety(t *testing.T) {
	f := NewFixture(t)
	sentinel := filepath.Join(t.TempDir(), "ran-foreign-prepush")
	hooks := gitHooksPath(t, f)
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatalf("create hooks dir: %v", err)
	}
	prepush := filepath.Join(hooks, "pre-push")
	script := fmt.Sprintf("#!/usr/bin/env bash\ntouch %q\nexit 1\n", sentinel)
	if err := os.WriteFile(prepush, []byte(script), 0o755); err != nil {
		t.Fatalf("write foreign pre-push: %v", err)
	}

	out := f.Bench("guards")

	out.RequireExit(0)
	requireGuardsLine(t, out.Stdout, `  pre-push,"",unmanaged (no manifest)`)
	if _, err := os.Stat(sentinel); err == nil {
		t.Fatal("bench guards executed a foreign pre-push")
	}
}

func testAXIBlockDangerousGitCoreUnreachableManifest(t *testing.T) {
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh")
	if _, err := os.Stat(hook); err != nil {
		t.Skipf("block-dangerous-git hook unavailable: %v", err)
	}
	f := NewFixture(t)

	out := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", hook, "--describe")

	out.RequireExit(0)
	out.RequireContains(out.Stdout, "manifest unavailable (analyzer missing)")
}

func testAXIBlockDangerousGitLinkedWorktreeClassification(t *testing.T) {
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh")
	if _, err := os.Stat(hook); err != nil {
		t.Skipf("block-dangerous-git hook unavailable: %v", err)
	}
	f := NewFixture(t)
	f.Env["GIT_AUTHOR_NAME"] = "Bench"
	f.Env["GIT_AUTHOR_EMAIL"] = "bench@local"
	f.Env["GIT_COMMITTER_NAME"] = "Bench"
	f.Env["GIT_COMMITTER_EMAIL"] = "bench@local"
	copyGuardsExecutable(t, filepath.Join(SubjectRoot(t), "bin", "bench.sh"), filepath.Join(f.Root, "bin", "bench.sh"))
	f.CommitAll("init")
	copyGuardsExecutable(t, filepath.Join(SubjectRoot(t), "dist", "bench"), filepath.Join(f.Root, "dist", "bench"))
	linked := filepath.Join(t.TempDir(), "linked")
	f.Git("worktree", "add", "-q", "--detach", linked, "HEAD")
	wt := Fixture{t: t, Root: linked, Env: f.Env}

	readonly := wt.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", "-c", `printf '{"tool_input":{"command":"git status"}}' | bash "$1"`, "sh", hook)
	readonly.RequireExit(0)
	destructive := wt.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", "-c", `printf '{"tool_input":{"command":"git reset --hard"}}' | bash "$1"`, "sh", hook)
	destructive.RequireExit(2)
}

func testSessionStartGuardBriefInjection(t *testing.T) {
	f := linkedGuardsFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	hook := filepath.Join(f.Root, ".bench", "hooks", "session-start.sh")

	out := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", hook)

	out.RequireExit(0)
	combined := out.Stdout + out.Stderr
	out.RequireContains(combined, "full manifests: bench guards")
	requireGuardsLineMatching(t, combined, `^bench CLI: .*\.bench/bin/bench\.sh \(bench not on PATH; invoke by path`)
}

func testSessionStartNeverBlocksOutsideRepo(t *testing.T) {
	f := NewFixture(t, WithNoRepo())
	hook := filepath.Join(SubjectRoot(t), ".bench", "hooks", "session-start.sh")

	out := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", hook)

	out.RequireExit(0)
	if out.Stdout != "" || out.Stderr != "" {
		t.Fatalf("session-start printed outside a repo\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
}

func linkedGuardsFixture(t *testing.T, opts ...FixtureOption) Fixture {
	t.Helper()
	f := NewFixture(t, opts...)
	f.Bench("link").RequireExit(0)
	return f
}

func gitHooksPath(t *testing.T, f Fixture) string {
	t.Helper()
	hooks := strings.TrimSpace(f.Git("rev-parse", "--git-path", "hooks").Stdout)
	if filepath.IsAbs(hooks) {
		return hooks
	}
	return filepath.Join(f.Root, hooks)
}

func gitPrePushPath(t *testing.T, f Fixture) string {
	t.Helper()
	return filepath.Join(gitHooksPath(t, f), "pre-push")
}

func copyGuardsExecutable(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read executable %s: %v", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatalf("create executable dir %s: %v", filepath.Dir(dst), err)
	}
	if err := os.WriteFile(dst, data, 0o755); err != nil {
		t.Fatalf("write executable %s: %v", dst, err)
	}
}

func firstGuardsLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}

func requireGuardsFirstLine(t *testing.T, output, want string) {
	t.Helper()
	if got := firstGuardsLine(output); got != want {
		t.Fatalf("first line = %q, want %q\noutput:\n%s", got, want, output)
	}
}

func requireGuardsLine(t *testing.T, output, want string) {
	t.Helper()
	for _, line := range strings.Split(output, "\n") {
		if line == want {
			return
		}
	}
	t.Fatalf("missing line %q\noutput:\n%s", want, output)
}

func requireGuardsLineMatching(t *testing.T, output, pattern string) {
	t.Helper()
	re := regexp.MustCompile(pattern)
	for _, line := range strings.Split(output, "\n") {
		if re.MatchString(line) {
			return
		}
	}
	t.Fatalf("missing line matching %q\noutput:\n%s", pattern, output)
}

func requireNoGuardsLineMatching(t *testing.T, output, pattern string) {
	t.Helper()
	re := regexp.MustCompile(pattern)
	for _, line := range strings.Split(output, "\n") {
		if re.MatchString(line) {
			t.Fatalf("unexpected line matching %q: %s\noutput:\n%s", pattern, line, output)
		}
	}
}

func nonEmptyGuardsLineCount(s string) int {
	count := 0
	for _, line := range strings.Split(s, "\n") {
		if line != "" {
			count++
		}
	}
	return count
}

func requireGuardsStringEqual(t *testing.T, got, want, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %q, want %q", msg, got, want)
	}
}
