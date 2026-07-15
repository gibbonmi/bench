package axi

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestAXIGuardsContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "AXI guards aggregation contract", testAXIGuardsAggregation)
	contract.RunParallel(t, "AXI guards --brief contract", testAXIGuardsBrief)
	contract.RunParallel(t, "AXI guards usage/subdirectory contract", testAXIGuardsUsageSubdirectory)
	contract.RunParallel(t, "AXI guards path-with-spaces contract", testAXIGuardsPathWithSpaces)
	contract.RunParallel(t, "AXI guards unmanaged-pre-push safety contract", testAXIGuardsUnmanagedPrePushSafety)
	contract.RunParallel(t, "AXI guards static non-execution sentinel contract", testAXIGuardsStaticNonExecution)
	contract.RunParallel(t, "AXI block-dangerous-git linked-worktree classification contract", testAXIBlockDangerousGitLinkedWorktreeClassification)
	contract.RunParallel(t, "session-start guard-brief injection contract", testSessionStartGuardBriefInjection)
	contract.RunParallel(t, "session-start resume failure warning contract", testSessionStartSurfacesResumeFailure)
	contract.RunParallel(t, "session-start never-blocks-outside-repo contract", testSessionStartNeverBlocksOutsideRepo)
	contract.RunParallel(t, "AXI maps/guards help contract", testAXIMapsGuardsHelp)
}

func testAXIGuardsAggregation(t *testing.T) {
	contract.NoteContractFailure(t, "AXI guards aggregation contract failed")
	f := linkedGuardsFixture(t)

	out := f.Bench("guards")
	out.RequireExit(0)
	requireGuardsFirstLine(t, out.Stdout, "guards[4]{guard,boundary,denies,wired}:")
	for _, guard := range []string{"block-dangerous-git", "check-agent-line", "stop", "pre-push"} {
		requireGuardsLineMatching(t, out.Stdout, "^  "+regexp.QuoteMeta(guard)+",")
	}
	requireNoGuardsLineMatching(t, out.Stdout, `^  session-start,`)

	// The wired cell names exactly the harness configs that reference each script.
	// The kit's own wiring is asymmetric: check-agent-line is wired only in
	// .claude/settings.json (its Agent matcher), so its cell is the bare "claude";
	// block-dangerous-git is wired in both configs, so its cell is the comma-joined,
	// TOON-quoted "claude,codex".
	requireGuardsLineMatching(t, out.Stdout, `^  check-agent-line,.*,claude$`)
	requireGuardsLineMatching(t, out.Stdout, `^  block-dangerous-git,.*,"claude,codex"$`)

	prepush := gitPrePushPath(t, f)
	// An orphan hook script referenced by neither config renders the definitive
	// `none`, never a blank cell.
	f.WriteExecutable(".bench/hooks/extra.sh", "#!/usr/bin/env bash\ncat >/dev/null\nexit 0\n")
	withExtra := f.Bench("guards")
	withExtra.RequireExit(0)
	requireGuardsLine(t, withExtra.Stdout, `  extra,"",no manifest,none`)

	if err := os.Remove(prepush); err != nil {
		t.Fatalf("remove generated pre-push: %v", err)
	}
	withoutPrePush := f.Bench("guards")
	withoutPrePush.RequireExit(0)
	requireGuardsLine(t, withoutPrePush.Stdout, `  pre-push,"",not installed,git`)
}

func testAXIGuardsBrief(t *testing.T) {
	f := linkedGuardsFixture(t)

	out := f.Bench("guards", "--brief")

	out.RequireExit(0)
	requireGuardsStringEqual(t, fmt.Sprint(strings.Count(out.Stdout, "full manifests: bench guards")), "1", "brief footer count")
	requireGuardsStringEqual(t, fmt.Sprint(nonEmptyGuardsLineCount(out.Stdout)), "5", "brief line count")
	out.RequireContains(out.Stdout, "block-dangerous-git: destructive git")
	// --brief carries the wired harnesses per guard so the SessionStart injection
	// stays honest about which configs can actually fire each hook.
	requireGuardsLineMatching(t, out.Stdout, `^check-agent-line: .*\[wired: claude\]$`)
	requireGuardsLineMatching(t, out.Stdout, `^block-dangerous-git: .*\[wired: claude,codex\]$`)
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
	subFixture := contract.NewFixtureAt(t, subdir, f.Env)
	fromSubdir := subFixture.Bench("guards")
	fromSubdir.RequireExit(0)
	if !strings.HasPrefix(firstGuardsLine(fromSubdir.Stdout), "guards[") {
		t.Fatalf("guards from subdirectory lost root resolution:\n%s", fromSubdir.Stdout)
	}
}

func testAXIGuardsPathWithSpaces(t *testing.T) {
	f := linkedGuardsFixture(t, contract.WithSpacePath())

	out := f.Bench("guards")

	out.RequireExit(0)
	if !strings.HasPrefix(firstGuardsLine(out.Stdout), "guards[") {
		t.Fatalf("space-path guards failed:\n%s", out.Stdout)
	}
}

func testAXIGuardsUnmanagedPrePushSafety(t *testing.T) {
	f := contract.NewFixture(t)
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
	requireGuardsLine(t, out.Stdout, `  pre-push,"",unmanaged (no manifest),git`)
	if _, err := os.Stat(sentinel); err == nil {
		t.Fatal("bench guards executed a foreign pre-push")
	}
}

func testAXIGuardsStaticNonExecution(t *testing.T) {
	f := contract.NewFixture(t)
	type sentinel struct {
		name   string
		header string
		path   string
	}
	fixtures := []sentinel{
		{"full", "# name: full\n# boundary: test\n# denies: a guarded action\n# why: sentinel\n", filepath.Join(t.TempDir(), "full-executed")},
		{"incomplete", "# name: incomplete\n# boundary: test\n# denies: a guarded action\n", filepath.Join(t.TempDir(), "incomplete-executed")},
		{"absent", "", filepath.Join(t.TempDir(), "absent-executed")},
		{"informational", "# name: informational\n# boundary: test\n# denies: nothing (informational)\n# why: sentinel\n", filepath.Join(t.TempDir(), "informational-executed")},
	}
	for _, fixture := range fixtures {
		body := fmt.Sprintf("#!/usr/bin/env bash\n%stouch %q\nexit 0\n", fixture.header, fixture.path)
		f.WriteExecutable(filepath.Join(".bench", "hooks", fixture.name+".sh"), body)
	}

	guards := f.Bench("guards")
	guards.RequireExit(0)
	requireGuardsLine(t, guards.Stdout, `  full,test,a guarded action,none`)
	requireGuardsLine(t, guards.Stdout, `  incomplete,"",no manifest,none`)
	requireGuardsLine(t, guards.Stdout, `  absent,"",no manifest,none`)
	requireNoGuardsLineMatching(t, guards.Stdout, `^  informational,`)
	f.Bench("session-inspect").RequireExit(0)
	for _, fixture := range fixtures {
		if _, err := os.Stat(fixture.path); !os.IsNotExist(err) {
			t.Fatalf("%s header fixture executed; evidence stat err=%v", fixture.name, err)
		}
	}
}

func testAXIBlockDangerousGitLinkedWorktreeClassification(t *testing.T) {
	hook := filepath.Join(contract.SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh")
	if _, err := os.Stat(hook); err != nil {
		t.Skipf("block-dangerous-git hook unavailable: %v", err)
	}
	f := contract.NewFixture(t)
	f.Env["GIT_AUTHOR_NAME"] = "Bench"
	f.Env["GIT_AUTHOR_EMAIL"] = "bench@local"
	f.Env["GIT_COMMITTER_NAME"] = "Bench"
	f.Env["GIT_COMMITTER_EMAIL"] = "bench@local"
	copyGuardsExecutable(t, filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), filepath.Join(f.Root, "bin", "bench.sh"))
	f.CommitAll("init")
	copyGuardsExecutable(t, filepath.Join(contract.SubjectRoot(t), "dist", "bench"), filepath.Join(f.Root, "dist", "bench"))
	linked := filepath.Join(t.TempDir(), "linked")
	f.Git("worktree", "add", "-q", "--detach", linked, "HEAD")
	wt := contract.NewFixtureAt(t, linked, f.Env)

	readonly := wt.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", "-c", `printf '{"tool_input":{"command":"git status"}}' | bash "$1"`, "sh", hook)
	readonly.RequireExit(0)
	destructive := wt.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", "-c", `printf '{"tool_input":{"command":"git reset --hard"}}' | bash "$1"`, "sh", hook)
	destructive.RequireExit(2)
}

func testSessionStartGuardBriefInjection(t *testing.T) {
	contract.NoteContractFailure(t, "session-start conservative cleanup contract failed")
	f := linkedGuardsFixture(t)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("session start fixture")
	hook := filepath.Join(f.Root, ".bench", "hooks", "session-start.sh")
	ordinary := filepath.Join(t.TempDir(), "ordinary foreign worktree")
	detached := filepath.Join(t.TempDir(), "detached foreign worktree")
	f.Git("worktree", "add", "-q", "-b", "foreign-session", ordinary, "HEAD")
	f.Git("worktree", "add", "-q", "--detach", detached, "HEAD")
	contract.WriteFileAbs(t, filepath.Join(detached, "unique.txt"), "unique\n")
	contract.RunAt(t, f, detached, nil, "git", "add", "unique.txt").RequireExit(0)
	contract.RunAt(t, f, detached, nil, "git", "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "unique detached").RequireExit(0)

	out := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", hook)

	out.RequireExit(0)
	combined := out.Stdout + out.Stderr
	out.RequireContains(combined, "full manifests: bench guards")
	out.RequireContains(combined, "retained foreign=2")
	resumeAt := strings.Index(combined, "bench resume:")
	statusAt := strings.Index(combined, "out-of-pool worktree")
	guardsAt := strings.Index(combined, "full manifests: bench guards")
	if resumeAt < 0 || statusAt < 0 || guardsAt < 0 || !(resumeAt < statusAt && statusAt < guardsAt) {
		t.Fatalf("session-start output order resume=%d status=%d guards=%d:\n%s", resumeAt, statusAt, guardsAt, combined)
	}
	for _, path := range []string{ordinary, detached} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("session-start removed foreign worktree %q: %v", path, err)
		}
	}
	requireGuardsLineMatching(t, combined, `^bench CLI: .*\.bench/bin/bench\.sh \(bench not on PATH; invoke by path`)
}

func testSessionStartNeverBlocksOutsideRepo(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	hook := filepath.Join(contract.SubjectRoot(t), ".bench", "hooks", "session-start.sh")

	out := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", hook)

	out.RequireExit(0)
	if out.Stdout != "" || out.Stderr != "" {
		t.Fatalf("session-start printed outside a repo\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
}

func testSessionStartSurfacesResumeFailure(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteExecutable(".bench/bin/bench.sh", `#!/usr/bin/env bash
case "${1:-}" in
  session-inspect)
    printf 'injected resume failure\n' >&2
    printf 'warning: bench session-start: resume-clean failed; inspect retained worktree state\n' >&2
    printf 'bench: clean — nothing pending\n'
    printf 'full manifests: bench guards\n'
    ;;
esac
`)
	hook := filepath.Join(contract.SubjectRoot(t), ".bench", "hooks", "session-start.sh")

	out := f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", hook)

	out.RequireExit(0)
	out.RequireContains(out.Stderr, "injected resume failure")
	out.RequireContains(out.Stderr, "warning: bench session-start: resume-clean failed; inspect retained worktree state")
	combined := out.Stdout + out.Stderr
	statusAt := strings.Index(combined, "bench: clean")
	guardsAt := strings.Index(combined, "full manifests: bench guards")
	if statusAt < 0 || guardsAt < 0 || statusAt >= guardsAt {
		t.Fatalf("session-start failure output order status=%d guards=%d:\n%s", statusAt, guardsAt, combined)
	}
}

func linkedGuardsFixture(t *testing.T, opts ...contract.FixtureOption) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t, opts...)
	f.Bench("link").RequireExit(0)
	return f
}

func gitHooksPath(t *testing.T, f contract.Fixture) string {
	t.Helper()
	hooks := strings.TrimSpace(f.Git("rev-parse", "--git-path", "hooks").Stdout)
	if filepath.IsAbs(hooks) {
		return hooks
	}
	return filepath.Join(f.Root, hooks)
}

func gitPrePushPath(t *testing.T, f contract.Fixture) string {
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

// testAXIMapsGuardsHelp pins row 8: the real, agent-invocable flags are documented
// in their own command's help, not findable only in Go source.
func testAXIMapsGuardsHelp(t *testing.T) {
	contract.NoteContractFailure(t, "AXI maps/guards help contract failed")
	f := contract.NewFixture(t)

	maps := f.Bench("maps", "-h")
	maps.RequireExit(0)
	requireContainsFold(t, maps.Stdout, "--count")

	guards := f.Bench("guards", "-h")
	guards.RequireExit(0)
	requireContainsFold(t, guards.Stdout, "--brief")
}
