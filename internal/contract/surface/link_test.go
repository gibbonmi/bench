package surface

import (
	"crypto/sha256"
	"fmt"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLinkContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench init does not scaffold .bench/learnings.md (self-learning journal)", testInitScaffoldsLearnings)
	contract.RunParallel(t, "a second bench init clobbered an existing .bench/gate.sh", testInitExistingGateIdempotence)
	contract.RunParallel(t, "bench link safe fresh/relink contract failed", testLinkSafeFreshRelink)
	contract.RunParallel(t, "bench link dist/.gitignore contract failed", testLinkWritesDistGitignore)
	contract.RunParallel(t, "bench link existing AGENTS.md contract failed", testLinkExistingAgents)
	contract.RunParallel(t, "bench link conflict contract failed", testLinkConflictWithoutManifest)
	contract.RunParallel(t, "bench link modified-managed contract failed", testLinkModifiedManaged)
	contract.RunParallel(t, "bench linked hooks local-CLI contract failed", testLinkedHooksLocalCLI)
	contract.RunParallel(t, "bench link metachar kit-path contract failed", testLinkMetacharKitPath)
	contract.RunParallel(t, "bench link worktree contract failed", testLinkWorktree)
	contract.RunParallel(t, "bench link hooksPath contract failed", testLinkHooksPath)
	contract.RunParallel(t, "bench link default-branch resolution contract failed", testLinkDefaultBranchResolution)
	contract.RunParallel(t, "bench link hooksPath conflict contract failed", testLinkHooksPathConflict)
	contract.RunParallel(t, "managed pre-push gate pin contract failed", testManagedPrePushGatePinning)
}

func testInitScaffoldsLearnings(t *testing.T) {
	f := contract.NewFixture(t)

	f.Bench("init").RequireExit(0)

	if !f.Exists(".bench/learnings.md") {
		t.Fatal("bench init does not scaffold .bench/learnings.md")
	}
	requireFixtureFileContains(t, f, ".bench/learnings.md", "/bench-what-next", "scaffolded journal header does not name /bench-what-next as the journal exit")
	requireFixtureFileNotContains(t, f, ".bench/learnings.md", "/bench-integrate-learnings", "scaffolded journal header still names the retired learnings-integration phase")
}

func testInitExistingGateIdempotence(t *testing.T) {
	f := contract.NewFixture(t)

	f.Bench("init").RequireExit(0)
	appendFile(t, filepath.Join(f.Root, ".bench", "gate.sh"), "# configured by hand\n")
	f.Bench("init").RequireExit(0)

	requireFixtureFileContains(t, f, ".bench/gate.sh", "# configured by hand", "a second bench init clobbered an existing .bench/gate.sh")
}

func testLinkSafeFreshRelink(t *testing.T) {
	f := contract.NewFixture(t)

	linkOK(t, f)

	requireLinkFile(t, f, "AGENTS.md")
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "fresh link did not create exactly one managed start marker")
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:end -->", 1, "fresh link did not create exactly one managed end marker")
	requireLinkFile(t, f, ".bench/BENCH.md")
	requireLinkFile(t, f, ".bench/BENCH-reference.md")
	requireFixtureFileNotContains(t, f, "CLAUDE.md", "@.bench/BENCH-reference.md", "fresh link @-imported the reference file; it must stay on-demand")
	requireExecutable(t, filepath.Join(f.Root, ".bench", "bin", "bench.sh"), "fresh link did not install local hook CLI .bench/bin/bench.sh")
	for _, rel := range []string{".bench/bin/bench-link.sh", ".bench/bin/bench-init.sh", ".bench/bin/bench-doctor.sh", ".bench/bin/bench-worktree.sh"} {
		requireLinkNotExists(t, f, rel, "fresh link still installed deleted local hook CLI "+filepath.Base(rel))
	}
	requireLinkFile(t, f, ".bench/link-manifest.tsv")
	requireFixtureFileContains(t, f, ".bench/link-manifest.tsv", "#kit\t", "fresh link did not stamp the kit version in the manifest")
	requireFixtureFileContains(t, f, ".bench/link-manifest.tsv", ".bench/BENCH.md\t", "fresh link manifest lost managed file rows after adding the stamp")
	requireLinkFile(t, f, ".agents/commands/bench-implement-spec.md")
	requireLinkFile(t, f, ".agents/skills/bench-craft-seams/SKILL.md")
	requireLinkFile(t, f, ".agents/skills/bench-implement-spec/SKILL.md")
	requireLinkFile(t, f, ".agents/skills/bench-implement-spec/agents/openai.yaml")
	requireLinkFile(t, f, ".claude/README.md")
	requireFixtureFileContains(t, f, ".claude/README.md", ".agents/", "Claude adapter README does not explain .agents")
	requireFixtureFileContains(t, f, ".claude/README.md", ".bench/hooks/", "Claude adapter README does not explain shared hooks")
	requireLinkFile(t, f, ".claude/commands/bench-implement-spec.md")
	requireLinkFile(t, f, ".claude/skills/bench-craft-seams/SKILL.md")
	requireLinkNotExists(t, f, ".claude/skills/bench-implement-spec", "fresh link installed a Codex-only phase adapter into .claude/skills (duplicates the /bench-implement-spec menu entry)")
	requireLinkFile(t, f, ".codex/hooks.json")
	requireLinkFile(t, f, ".bench/hooks/block-dangerous-git.sh")
	requireExecutable(t, filepath.Join(f.Root, ".bench", "adapters", "claude"), "fresh link did not install executable reference adapters")
	requireLinkFile(t, f, ".bench/lib/resolve-bench.sh")
	requireLinkFile(t, f, ".bench/hooks/session-start.sh")
	requireFixtureFileContains(t, f, ".claude/settings.json", "SessionStart", "fresh link .claude/settings.json has no SessionStart wiring")
	requireExecutable(t, filepath.Join(f.Root, ".git", "hooks", "pre-push"), "fresh link did not install git pre-push hook")
	requireFixtureFileContains(t, f, "CLAUDE.md", "@.bench/BENCH.md", "fresh link CLAUDE.md does not import .bench/BENCH.md")
	requireNotSymlink(t, filepath.Join(f.Root, ".agents", "commands", "bench-implement-spec.md"), "default link mode symlinked portable commands")

	beforeManifest := fileSum(t, filepath.Join(f.Root, ".bench", "link-manifest.tsv"))
	for _, sub := range []string{"link", "init", "doctor"} {
		probe := f.Run(filepath.Join(".bench", "bin", "bench.sh"), sub)
		if probe.ExitCode == 0 {
			t.Fatalf("planted .bench/bin/bench.sh %s succeeded instead of refusing", sub)
		}
		probe.RequireContains(probe.Stderr, "real Bench kit")
	}
	requireLinkEqual(t, fileSum(t, filepath.Join(f.Root, ".bench", "link-manifest.tsv")), beforeManifest, "planted adoption command mutated the link manifest")
	requireLinkNotExists(t, f, ".bench/gate.sh", "planted init scaffolded a gate despite refusing")

	linkOK(t, f)
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "relink duplicated managed Bench block")
	if got := manifestRowCount(t, f, "CLAUDE.md"); got != 1 {
		t.Fatalf("relink of an already link-owned CLAUDE.md left %d manifest rows, want 1", got)
	}

	f.WriteFile("CLAUDE.md", "# Bench\n\nCanonical agreement in AGENTS.md.\n\n@AGENTS.md\n")
	linkOK(t, f)
	requireFixtureFileContains(t, f, "CLAUDE.md", "@.bench/BENCH.md", "relink did not retrofit the legacy bench-generated CLAUDE.md")
	if got := manifestRowCount(t, f, "CLAUDE.md"); got != 1 {
		t.Fatalf("relink of a legacy-shaped CLAUDE.md left %d manifest rows, want 1", got)
	}

	f.WriteFile("CLAUDE.md", "# Custom\n\nproject-owned claude config\n")
	linkOK(t, f)
	requireFixtureFileContains(t, f, "CLAUDE.md", "project-owned claude config", "relink rewrote a project-owned CLAUDE.md")
	requireFixtureFileNotContains(t, f, "CLAUDE.md", "@.bench/BENCH.md", "relink injected an import into a project-owned CLAUDE.md")
	if got := manifestRowCount(t, f, "CLAUDE.md"); got != 0 {
		t.Fatalf("relink recorded a project-owned CLAUDE.md in the manifest (%d rows), want 0", got)
	}
}

func testLinkWritesDistGitignore(t *testing.T) {
	f := contract.NewFixture(t)

	linkOK(t, f)

	requireLinkFile(t, f, ".bench/dist/.gitignore")
	requireFixtureFileContains(t, f, ".bench/dist/.gitignore", "bench", "dist gitignore does not ignore the copied binary")
	requireFixtureFileNotContains(t, f, ".bench/dist/.gitignore", ".gitignore", "dist gitignore ignores itself; it must travel with the repo")
	requireFixtureFileContains(t, f, ".bench/link-manifest.tsv", ".bench/dist/.gitignore\t", "dist gitignore has no manifest row")

	before := f.ReadFile(".bench/dist/.gitignore")
	beforeRows := manifestRowCount(t, f, ".bench/dist/.gitignore")

	linkOK(t, f)

	if got := f.ReadFile(".bench/dist/.gitignore"); got != before {
		t.Fatalf("relink rewrote .bench/dist/.gitignore non-identically:\nbefore:\n%q\nafter:\n%q", before, got)
	}
	afterRows := manifestRowCount(t, f, ".bench/dist/.gitignore")
	if beforeRows != 1 || afterRows != 1 {
		t.Fatalf("dist gitignore manifest rows: first link %d, relink %d, want 1 and 1", beforeRows, afterRows)
	}
}

func manifestRowCount(t *testing.T, f contract.Fixture, rel string) int {
	t.Helper()
	n := 0
	for _, line := range strings.Split(f.ReadFile(".bench/link-manifest.tsv"), "\n") {
		if strings.HasPrefix(line, rel+"\t") {
			n++
		}
	}
	return n
}

func testLinkExistingAgents(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("AGENTS.md", "PROJECT RULES\n")

	linkOK(t, f)

	requireFixtureFileContains(t, f, "AGENTS.md", "PROJECT RULES", "existing AGENTS.md content was clobbered")
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "existing AGENTS.md did not get exactly one managed block")
}

func testLinkConflictWithoutManifest(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile(".agents/commands/bench-implement-spec.md", "project command\n")

	probe := f.Bench("link")

	if probe.ExitCode == 0 {
		t.Fatal("link succeeded despite a project-owned command conflict")
	}
	probe.RequireContains(strings.ToLower(probe.Stderr+probe.Stdout), "conflict")
	requireFixtureFileContains(t, f, ".agents/commands/bench-implement-spec.md", "project command", "conflicting project command was overwritten")
	requireLinkNotExists(t, f, ".bench/link-manifest.tsv", "conflicting link wrote a manifest despite failing")
}

func testLinkModifiedManaged(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	appendFile(t, filepath.Join(f.Root, ".agents", "commands", "bench-implement-spec.md"), "\nlocal edit\n")

	probe := f.Bench("link")

	if probe.ExitCode == 0 {
		t.Fatal("relink overwrote a locally modified managed file")
	}
	probe.RequireContains(strings.ToLower(probe.Stderr+probe.Stdout), "modified")
}

func testLinkedHooksLocalCLI(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	f.WriteExecutable(".bench/gate.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("linked")

	f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, filepath.Join(".bench", "hooks", "session-start.sh")).RequireExit(0)
	f.RunEnv(map[string]string{"PATH": "/usr/bin:/bin"}, "bash", "-c", "printf '{}\n' | BENCH_SHIFT=1 .bench/hooks/stop.sh >/dev/null 2>&1").RequireExit(0)
}

func testLinkMetacharKitPath(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	kitcopy := filepath.Join(f.Root, "kit[1]")
	for _, dir := range []string{
		filepath.Join(kitcopy, ".bench"),
		filepath.Join(kitcopy, "dist"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	root := contract.KitRoot(t)
	copyPaths(t, kitcopy,
		filepath.Join(root, "bin"),
		filepath.Join(root, ".agents"),
		filepath.Join(root, ".claude"),
		filepath.Join(root, ".codex"),
	)
	copyFileTo(t, filepath.Join(root, "dist", "bench"), filepath.Join(kitcopy, "dist", "bench"))
	copyFileTo(t, filepath.Join(root, "AGENTS.md"), filepath.Join(kitcopy, "AGENTS.md"))
	copyFileTo(t, filepath.Join(root, ".bench", "BENCH.md"), filepath.Join(kitcopy, ".bench", "BENCH.md"))
	copyFileTo(t, filepath.Join(root, ".bench", "BENCH-reference.md"), filepath.Join(kitcopy, ".bench", "BENCH-reference.md"))
	copyPaths(t, filepath.Join(kitcopy, ".bench"),
		filepath.Join(root, ".bench", "hooks"),
		filepath.Join(root, ".bench", "adapters"),
		filepath.Join(root, ".bench", "lib"),
	)

	repo := filepath.Join(f.Root, "repo")
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	r := linkFixtureAt(t, repo, f.Env)
	r.Git("init", "-q")
	probe := r.BenchEnv(map[string]string{"BENCH_KIT": kitcopy}, "link")

	probe.RequireExit(0)
	requireExecutable(t, filepath.Join(r.Root, ".bench", "bin", "bench.sh"), "metachar kit path scattered installed files")
	requireLinkFile(t, r, ".agents/commands/bench-implement-spec.md")
}

func testLinkWorktree(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	mainRepo := filepath.Join(f.Root, "main-repo")
	f.Run("git", "init", "-q", mainRepo).RequireExit(0)
	main := linkFixtureAt(t, mainRepo, f.Env)
	main.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "--allow-empty", "-m", "init")
	wt := filepath.Join(f.Root, "wt")
	main.Git("worktree", "add", "-q", wt, "-b", "side", "HEAD")
	worktree := linkFixtureAt(t, wt, f.Env)

	linkOK(t, worktree)

	probe := worktree.Run("sh", "-c", `hooks="$(git rev-parse --git-path hooks)" && test -x "$hooks/pre-push"`)
	if probe.ExitCode != 0 {
		t.Fatalf("worktree link did not install pre-push in the effective hooks dir\nstdout:\n%s\nstderr:\n%s", probe.Stdout, probe.Stderr)
	}
}

func testLinkHooksPath(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("config", "core.hooksPath", ".husky")

	linkOK(t, f)

	requireExecutable(t, filepath.Join(f.Root, ".husky", "pre-push"), "pre-push not installed into configured hooksPath")
	requireFixtureFileContains(t, f, ".husky/pre-push", "bench:managed-pre-push", "hooksPath pre-push is not bench-managed")
}

func testLinkDefaultBranchResolution(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	f.Run("git", "init", "-q", "--bare", "-b", "master", "remote.git").RequireExit(0)
	f.Run("git", "init", "-q", "-b", "master", "repo").RequireExit(0)
	repo := linkFixtureAt(t, filepath.Join(f.Root, "repo"), f.Env)
	repo.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "--allow-empty", "-m", "init")
	repo.Git("remote", "add", "origin", filepath.Join(f.Root, "remote.git"))
	repo.Git("push", "-q", "origin", "master")
	repo.Git("fetch", "-q", "origin")
	repo.GitAllow("symbolic-ref", "-d", "refs/remotes/origin/HEAD")

	linkOK(t, repo)

	hooks := strings.TrimSpace(repo.Git("rev-parse", "--git-path", "hooks").Stdout)
	requireFixtureFileContains(t, repo, filepath.ToSlash(filepath.Join(hooks, "pre-push")), "protected=\"master\"", "pre-push bakes the wrong branch fallback when origin/HEAD is unset")
}

func testLinkHooksPathConflict(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("config", "core.hooksPath", ".husky")
	f.WriteExecutable(".husky/pre-push", "#!/bin/sh\nexit 0\n")

	probe := f.Bench("link")

	if probe.ExitCode == 0 {
		t.Fatal("link succeeded over a non-managed pre-push in hooksPath")
	}
	probe.RequireContains(strings.ToLower(probe.Stderr+probe.Stdout), "conflict")
	requireFixtureFileContains(t, f, ".husky/pre-push", "exit 0", "hooksPath conflict overwrote the project hook")
}

func linkOK(t *testing.T, f contract.Fixture) contract.Probe {
	t.Helper()
	probe := f.Bench("link")
	probe.RequireExit(0)
	return probe
}

func linkFixtureAt(t testing.TB, root string, env map[string]string) contract.Fixture {
	t.Helper()
	return contract.NewFixtureAt(t, root, env)
}

func requireLinkFile(t *testing.T, f contract.Fixture, rel string) {
	t.Helper()
	if !f.Exists(rel) {
		t.Fatalf("missing %s", rel)
	}
}

func requireLinkNotExists(t *testing.T, f contract.Fixture, rel, msg string) {
	t.Helper()
	if _, err := os.Lstat(filepath.Join(f.Root, filepath.FromSlash(rel))); err == nil {
		t.Fatal(msg)
	}
}

func requireExecutable(t *testing.T, path, msg string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
	if info.Mode()&0o111 == 0 {
		t.Fatal(msg)
	}
}

func requireNotSymlink(t *testing.T, path, msg string) {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatal(msg)
	}
}

func requireFixtureFileContains(t *testing.T, f contract.Fixture, rel, needle, msg string) {
	t.Helper()
	data := f.ReadFile(rel)
	if !strings.Contains(data, needle) {
		t.Fatalf("%s: missing %q in %s\n%s", msg, needle, rel, data)
	}
}

func requireFixtureFileNotContains(t *testing.T, f contract.Fixture, rel, needle, msg string) {
	t.Helper()
	data := f.ReadFile(rel)
	if strings.Contains(data, needle) {
		t.Fatalf("%s: unexpected %q in %s\n%s", msg, needle, rel, data)
	}
}

func requireLiteralCount(t *testing.T, f contract.Fixture, rel, needle string, want int, msg string) {
	t.Helper()
	got := strings.Count(f.ReadFile(rel), needle)
	if got != want {
		t.Fatalf("%s: got %d, want %d", msg, got, want)
	}
}

func requireLinkEqual(t *testing.T, got, want, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %q, want %q", msg, got, want)
	}
}

func appendFile(t *testing.T, path, contents string) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open append %s: %v", path, err)
	}
	defer file.Close()
	if _, err := file.WriteString(contents); err != nil {
		t.Fatalf("append %s: %v", path, err)
	}
}

func fileSum(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func copyPaths(t *testing.T, destParent string, sources ...string) {
	t.Helper()
	args := append([]string{"-R"}, sources...)
	args = append(args, destParent)
	out, err := exec.Command("cp", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("copy paths to %s: %v\n%s", destParent, err, out)
	}
}

func copyFileTo(t *testing.T, src, dest string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(dest), err)
	}
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	info, err := os.Stat(src)
	if err != nil {
		t.Fatalf("stat %s: %v", src, err)
	}
	if err := os.WriteFile(dest, data, info.Mode().Perm()); err != nil {
		t.Fatalf("write %s: %v", dest, err)
	}
}
