package surface

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

// This intentionally stays outside TestLinkContracts' parallel fan-out: a FIFO must
// prove it cannot block the CLI without scheduler contention masking that deadline.
func TestLinkLifecycleHostileContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	testLinkRejectsHostileDroppedRows(t)
	testLinkRejectsFIFOInAllowlistedKitTree(t)
	testLinkRejectsSymlinkInAllowlistedKitTree(t)
}

func TestLinkSymlinkLifecycleContracts(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	testLinkConvergesManagedSymlink(t)
	testLinkRejectsDriftedManagedSymlink(t)
	testLinkRejectsNewEntryUnderManagedSymlink(t)
	testUpgradeRefreshesHookThroughManagedSymlink(t)
	testLinkManagedSymlinkIsIdempotent(t)
	testLinkLeavesCleanEntryInPlace(t)
	testLinkConvergesAdapterDirectorySymlink(t)
	testLinkRejectsDriftedAdapterDirectorySymlink(t)
}

// adapterFileRel is the .claude mirror of managedFileRel — the adapter entry whose staged
// symlink names that canonical file.
var adapterFileRel = strings.Replace(managedFileRel, ".agents/", ".claude/", 1)

// testLinkConvergesAdapterDirectorySymlink pins the shape a repo takes when it satisfies
// the whole adapter mirror itself: .claude/commands is one directory symlink to the
// canonical directory, so every adapter destination resolves to exactly the file its
// staged link names. That is converged, and a link that read it as needing a write would
// hit the symlink-parent refusal on the first adapter entry.
func testLinkConvergesAdapterDirectorySymlink(t *testing.T) {
	f := adapterDirectorySymlinkFixture(t, "../.agents/commands")
	before := fixtureStateExceptManifest(t, f)

	probe := f.Bench("link")

	probe.RequireExit(0)
	if after := fixtureStateExceptManifest(t, f); after != before {
		t.Fatalf("link wrote to a converged destination under an adapter directory symlink\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// testLinkRejectsDriftedAdapterDirectorySymlink keeps the hard refusal intact for the
// case content-awareness must not swallow: the adapter directory symlink points at a
// stale mirror, so once the canonical file changes the destination no longer resolves to
// what the staged link would name and the entry genuinely needs a write.
func testLinkRejectsDriftedAdapterDirectorySymlink(t *testing.T) {
	f := adapterDirectorySymlinkFixture(t, "../adapter-mirror")
	f.WriteFile(managedFileRel, "canonical content, locally changed\n")
	before := fixtureState(t, f)

	probe := f.Bench("link")

	probe.RequireExit(1)
	probe.RequireContains(probe.Stderr, adapterFileRel+" has a symlink parent directory")
	if after := fixtureState(t, f); after != before {
		t.Fatalf("drifted adapter refusal still promoted part of the transaction\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// adapterDirectorySymlinkFixture replaces a linked repo's .claude/commands directory with
// a symlink to target. The "../adapter-mirror" target is seeded with a copy of the
// canonical directory so a later edit to a canonical file makes exactly one adapter
// destination drift; any other target is used as given.
func adapterDirectorySymlinkFixture(t *testing.T, target string) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	linkOK(t, f)
	if target == "../adapter-mirror" {
		copyPaths(t, f.Root, filepath.Join(f.Root, ".agents", "commands"))
		if err := os.Rename(filepath.Join(f.Root, "commands"), filepath.Join(f.Root, "adapter-mirror")); err != nil {
			t.Fatal(err)
		}
	}
	dir := filepath.Join(f.Root, ".claude", "commands")
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, dir); err != nil {
		t.Fatal(err)
	}
	return f
}

// fixtureStateExceptManifest elides the link manifest from fixtureState. A converged entry
// records its destination's own fingerprint, so an adapter destination reached through a
// directory symlink restamps its row from the link it no longer is to the bytes it now
// resolves to — a manifest change with no destination write behind it.
func fixtureStateExceptManifest(t *testing.T, f contract.Fixture) string {
	t.Helper()
	kept := []string{}
	for _, line := range strings.Split(fixtureState(t, f), "\n") {
		if !strings.HasPrefix(line, ".bench/link-manifest.tsv ") {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

func testLinkConvergesManagedSymlink(t *testing.T) {
	f := convergedManagedSymlinkFixture(t)

	f.Bench("link").RequireExit(0)
}

// testLinkRejectsDriftedManagedSymlink pins the symlink-parent abort ahead of the soft
// conflict classification: a drifted entry needs a write, so the deliberately symlinked
// directory beneath it is a hard refusal for the whole transaction rather than one row
// of an exit-3 conflicts report that still promotes everything else.
func testLinkRejectsDriftedManagedSymlink(t *testing.T) {
	f := convergedManagedSymlinkFixture(t)
	f.WriteFile(".agents/commands/bench-implement-spec.md", "consumer drift\n")
	before := fixtureState(t, f)

	probe := f.Bench("link")

	probe.RequireExit(1)
	probe.RequireContains(probe.Stderr, ".agents/commands/bench-implement-spec.md has a symlink parent directory")
	if after := fixtureState(t, f); after != before {
		t.Fatalf("symlink-parent refusal still promoted part of the transaction\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func testLinkRejectsNewEntryUnderManagedSymlink(t *testing.T) {
	f := contract.NewFixture(t)
	kit := payloadTestKit(t, f, "symlink-new-entry")
	f.BenchEnv(map[string]string{"BENCH_KIT": kit}, "link").RequireExit(0)
	makeManagedAgentsSymlink(t, f)
	contract.WriteFileAbs(t, filepath.Join(kit, ".agents", "commands", "new-managed-command.md"), "new managed command\n")

	probe := f.BenchEnv(map[string]string{"BENCH_KIT": kit}, "link")

	probe.RequireExit(1)
	probe.RequireContains(probe.Stderr, "new-managed-command.md has a symlink parent directory")
}

func testUpgradeRefreshesHookThroughManagedSymlink(t *testing.T) {
	f := contract.NewFixture(t)
	linkedRepoAtOlderVersion(t, f, "0.0.1")
	makeManagedAgentsSymlink(t, f)
	hook := prePushPath(t, f)
	current := contract.ReadFileAbs(t, hook)
	contract.WriteFileAbs(t, hook, current+"\n# stale hook\n")

	f.Bench("upgrade").RequireExit(0)

	requireLinkEqual(t, contract.ReadFileAbs(t, hook), current, "upgrade through a converged managed symlink did not restore current hook bytes")
}

func testLinkManagedSymlinkIsIdempotent(t *testing.T) {
	f := convergedManagedSymlinkFixture(t)
	f.Bench("link").RequireExit(0)
	before := fixtureState(t, f)

	f.Bench("link").RequireExit(0)

	if after := fixtureState(t, f); after != before {
		t.Fatalf("second link through a converged managed symlink changed tree or manifest\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// testLinkLeavesCleanEntryInPlace pins the skip as a genuine one rather than a rewrite
// of identical bytes. Comparing content would pass either way, so the inode is the
// observable that separates them: promotion renames a staged replacement over the
// destination and changes it, while a skipped entry keeps the file it already had.
func testLinkLeavesCleanEntryInPlace(t *testing.T) {
	plain := contract.NewFixture(t)
	linkOK(t, plain)
	requireCleanEntryInodeStable(t, plain, "clean entry under an ordinary parent")
	requireCleanEntryInodeStable(t, convergedManagedSymlinkFixture(t), "clean entry under a symlink parent")
}

func requireCleanEntryInodeStable(t *testing.T, f contract.Fixture, context string) {
	t.Helper()
	before := fixtureInode(t, f, managedFileRel)

	f.Bench("link").RequireExit(0)

	if after := fixtureInode(t, f, managedFileRel); after != before {
		t.Fatalf("%s: link replaced %s (inode %d -> %d) instead of skipping it", context, managedFileRel, before, after)
	}
}

func fixtureInode(t *testing.T, f contract.Fixture, rel string) uint64 {
	t.Helper()
	info, err := os.Stat(filepath.Join(f.Root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		t.Fatalf("inode unavailable for %s on this platform", rel)
	}
	return uint64(stat.Ino)
}

func convergedManagedSymlinkFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	linkOK(t, f)
	makeManagedAgentsSymlink(t, f)
	return f
}

func makeManagedAgentsSymlink(t *testing.T, f contract.Fixture) {
	t.Helper()
	managed := filepath.Join(f.Root, "managed-agents")
	if err := os.Rename(filepath.Join(f.Root, ".agents"), managed); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("managed-agents", filepath.Join(f.Root, ".agents")); err != nil {
		t.Fatal(err)
	}
}

// TestLinkPayloadAllowlistContracts pins the consumer-payload allowlist as the single
// source of what a fresh link writes: kit-only surfaces reach no destination (FT85
// story 1), a stray file under an allowlisted tree's kit-only subtree is not linked
// even though the tree itself is walked (story 2), and a space-bearing path inside an
// allowlisted tree survives link and unlink intact (story 2 edge inventory row).
func TestLinkPayloadAllowlistContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench link kit-only exclusion contract failed", testLinkExcludesKitOnlySurfaces)
	contract.RunParallel(t, "bench link stray kit-only-tree asset contract failed", testLinkExcludesStrayAssetUnderKitOnlyTree)
	contract.RunParallel(t, "bench link space-bearing allowlisted path contract failed", testLinkSpaceBearingAllowlistedTreePath)
}

// testLinkExcludesKitOnlySurfaces is the FT85 story 1 red signal: a fresh link must
// write no bench-assess, bench-update-kit, or craft-synthesis asset to any destination,
// including the .claude/ adapter mirrors. Before the allowlist's audience filter
// existed, buildLinkPlan copied both .agents trees wholesale, so every one of these
// paths would be present.
func testLinkExcludesKitOnlySurfaces(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)

	for _, rel := range []string{
		".agents/commands/bench-assess.md",
		".agents/commands/bench-update-kit.md",
		".agents/skills/bench-assess/SKILL.md",
		".agents/skills/bench-update-kit/SKILL.md",
		".agents/skills/bench-craft-synthesis/SKILL.md",
		".claude/commands/bench-assess.md",
		".claude/commands/bench-update-kit.md",
		".claude/skills/bench-craft-synthesis/SKILL.md",
	} {
		requireLinkNotExists(t, f, rel, "fresh link wrote kit-only surface "+rel)
	}
}

// testLinkExcludesStrayAssetUnderKitOnlyTree is the FT85 story 2 red signal: the link
// plan's consumer entries equal the allowlist's consumer rows, so a file added under a
// kit-only subtree of an otherwise-walked allowlisted tree (.agents/skills) is not
// written even though the tree itself is linked. Before the plan was allowlist-driven,
// the tree walk had no exclusion at all, so this stray file would appear.
func testLinkExcludesStrayAssetUnderKitOnlyTree(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	kit := payloadTestKit(t, f, "kit-stray")
	contract.WriteFileAbs(t, filepath.Join(kit, ".agents", "skills", "bench-assess", "stray-new-file.md"), "not part of any released bench-assess skill\n")

	repo := filepath.Join(f.Root, "repo")
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	r := linkFixtureAt(t, repo, f.Env)
	r.Git("init", "-q")
	r.BenchEnv(map[string]string{"BENCH_KIT": kit}, "link").RequireExit(0)

	requireLinkNotExists(t, r, ".agents/skills/bench-assess/stray-new-file.md", "link wrote a stray file added under a kit-only subtree of an allowlisted tree")
}

// testLinkSpaceBearingAllowlistedTreePath is the FT85 story 2 edge-inventory row: an
// allowlisted tree containing a space-bearing path links and unlinks intact.
func testLinkSpaceBearingAllowlistedTreePath(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	kit := payloadTestKit(t, f, "kit-space")
	spaceRel := filepath.Join(".agents", "skills", "bench-craft-seams", "extra notes.md")
	contract.WriteFileAbs(t, filepath.Join(kit, spaceRel), "space-bearing path payload\n")

	repo := filepath.Join(f.Root, "repo")
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	r := linkFixtureAt(t, repo, f.Env)
	r.Git("init", "-q")
	r.BenchEnv(map[string]string{"BENCH_KIT": kit}, "link").RequireExit(0)
	requireLinkFile(t, r, filepath.ToSlash(spaceRel))

	unlinkOK(t, r)
	requireLinkNotExists(t, r, filepath.ToSlash(spaceRel), "unlink left a space-bearing allowlisted path behind")
}

// testLinkRejectsFIFOInAllowlistedKitTree is the FT85 story 2 edge-inventory row for
// special files in script-discovery paths: a non-regular file inside an allowlisted
// tree on the kit side is refused by name before it is read, rather than blocking the
// plan builder or silently vanishing from the linked tree.
func testLinkRejectsFIFOInAllowlistedKitTree(t *testing.T) {
	requireLinkRefusesSpecialFile(t, "kit-fifo", filepath.Join(".agents", "skills", "bench-craft-seams", "blocked.fifo"), "blocked.fifo",
		func(path string) error { return syscall.Mkfifo(path, 0o644) })
}

// testLinkRejectsSymlinkInAllowlistedKitTree covers the same refusal for the special
// file a real kit is far likelier to grow — a symbolic link — and does it under
// .agents/commands, so the walk's refusal is pinned for an allowlisted tree other than
// .agents/skills. The link points at a regular file that link would otherwise be happy
// to copy, so what is refused is the indirection itself, not an unreadable target.
func testLinkRejectsSymlinkInAllowlistedKitTree(t *testing.T) {
	requireLinkRefusesSpecialFile(t, "kit-symlink", filepath.Join(".agents", "commands", "linked-command.md"), "symbolic link",
		func(path string) error {
			return os.Symlink(filepath.Join(filepath.Dir(path), "bench-implement-spec.md"), path)
		})
}

// requireLinkRefusesSpecialFile plants one special file inside an allowlisted kit tree
// and asserts link refuses by name without blocking and without leaving a manifest
// behind. The deadline is the point for the FIFO case, so every case carries it: a plan
// builder that opens what it walks would hang here rather than fail.
func requireLinkRefusesSpecialFile(t *testing.T, kitName, relPath, wantStderr string, create func(path string) error) {
	t.Helper()
	f := contract.NewFixture(t, contract.WithNoRepo())
	kit := payloadTestKit(t, f, kitName)
	if err := create(filepath.Join(kit, relPath)); err != nil {
		t.Fatal(err)
	}

	repo := filepath.Join(f.Root, "repo")
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	r := linkFixtureAt(t, repo, f.Env)
	r.Git("init", "-q")
	probe := contract.RunAtWithTimeout(t, r, r.Root, map[string]string{"BENCH_KIT": kit}, time.Second, "bash", filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), "link")
	if probe.TimedOut {
		t.Fatalf("link blocked on %s inside an allowlisted kit tree", relPath)
	}
	probe.RequireExit(1)
	probe.RequireContains(probe.Stderr, wantStderr)
	requireLinkNotExists(t, r, ".bench/link-manifest.tsv", "link refusing a hostile kit tree still wrote a manifest")
}

// payloadTestKit builds a full working kit copy (every source the allowlist can name)
// so a test can mutate one file before linking against it via BENCH_KIT.
func payloadTestKit(t *testing.T, f contract.Fixture, name string) string {
	t.Helper()
	root := contract.KitRoot(t)
	kit := filepath.Join(f.Root, name)
	contract.Mkdir(t, filepath.Join(kit, ".bench"))
	contract.Mkdir(t, filepath.Join(kit, "dist"))
	copyPaths(t, kit, filepath.Join(root, "bin"), filepath.Join(root, ".agents"), filepath.Join(root, ".claude"), filepath.Join(root, ".codex"))
	copyFileTo(t, filepath.Join(root, ".bench", "BENCH.md"), filepath.Join(kit, ".bench", "BENCH.md"))
	copyFileTo(t, filepath.Join(root, ".bench", "BENCH-reference.md"), filepath.Join(kit, ".bench", "BENCH-reference.md"))
	copyFileTo(t, filepath.Join(root, "dist", "bench"), filepath.Join(kit, "dist", "bench"))
	copyPaths(t, filepath.Join(kit, ".bench"), filepath.Join(root, ".bench", "hooks"), filepath.Join(root, ".bench", "adapters"), filepath.Join(root, ".bench", "lib"))
	return kit
}

func TestLinkMarkerFenceContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench link malformed marker contract failed", testLinkMalformedMarker)
	contract.RunParallel(t, "bench link fenced-marker contract failed", testLinkFencedMarker)
	contract.RunParallel(t, "bench link unclosed-fence contract failed", testLinkUnclosedFence)
}

func TestLinkPrePushLifecycleContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench link leaves installation open when origin head probing fails", testLinkOriginHeadProbeFailureOpen)
	contract.RunParallel(t, "bench link refuses a dangling pre-push hook", testLinkDanglingPrePushRefusal)
	contract.RunParallel(t, "bench link refuses a marker-less pre-push hook", testLinkMarkerlessPrePushRefusal)
}

func testLinkOriginHeadProbeFailureOpen(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("remote", "add", "origin", filepath.Join(f.Root, "missing-remote.git"))

	linkOK(t, f)
	hooks := strings.TrimSpace(f.Git("rev-parse", "--git-path", "hooks").Stdout)
	requireExecutable(t, filepath.Join(f.Root, hooks, "pre-push"), "link failed open when origin head probe failed")
}

func testLinkDanglingPrePushRefusal(t *testing.T) {
	f := contract.NewFixture(t)
	hooks := strings.TrimSpace(f.Git("rev-parse", "--git-path", "hooks").Stdout)
	path := filepath.Join(f.Root, hooks, "pre-push")
	if err := os.Symlink("missing-pre-push", path); err != nil {
		t.Fatal(err)
	}

	probe := f.Bench("link")
	if probe.ExitCode == 0 {
		t.Fatal("link succeeded over a dangling pre-push hook")
	}
	probe.RequireContains(strings.ToLower(probe.Stderr+probe.Stdout), "conflict")
	target, err := os.Readlink(path)
	if err != nil || target != "missing-pre-push" {
		t.Fatalf("dangling pre-push after refusal = %q, %v", target, err)
	}
}

func testLinkMarkerlessPrePushRefusal(t *testing.T) {
	f := contract.NewFixture(t)
	hooks := strings.TrimSpace(f.Git("rev-parse", "--git-path", "hooks").Stdout)
	rel := filepath.ToSlash(filepath.Join(hooks, "pre-push"))
	contents := "#!/bin/sh\nexit 0\n"
	f.WriteExecutable(rel, contents)

	probe := f.Bench("link")
	if probe.ExitCode == 0 {
		t.Fatal("link succeeded over a marker-less pre-push hook")
	}
	probe.RequireContains(strings.ToLower(probe.Stderr+probe.Stdout), "conflict")
	requireFixtureFileContains(t, f, rel, contents, "marker-less pre-push changed after refusal")
}

func testLinkMalformedMarker(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("AGENTS.md", "PROJECT BEFORE\n<!-- bench:end -->\nPROJECT MIDDLE\n<!-- bench:start -->\nPROJECT AFTER\n")

	probe := f.Bench("link")

	if probe.ExitCode == 0 {
		t.Fatal("link succeeded despite reversed Bench managed block markers")
	}
	probe.RequireContains(strings.ToLower(probe.Stderr+probe.Stdout), "malformed")
	requireFixtureFileContains(t, f, "AGENTS.md", "PROJECT AFTER", "malformed marker failure still rewrote project-owned text")
}

func testLinkFencedMarker(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("AGENTS.md", "# Project rules\n\nHow Bench marks its block:\n\n```\n<!-- bench:start -->\nmanaged content example\n<!-- bench:end -->\n```\n\nKEEP-ME project text.\n")

	linkOK(t, f)

	requireFixtureFileContains(t, f, "AGENTS.md", "KEEP-ME project text.", "fenced-marker link lost project text")
	requireFixtureFileContains(t, f, "AGENTS.md", "managed content example", "fenced example content was rewritten")
	linkOK(t, f)
	requireFixtureFileContains(t, f, "AGENTS.md", "managed content example", "relink consumed the fenced example")
	requireLiteralCount(t, f, "AGENTS.md", "## Bench", 1, "fenced markers caused duplicate managed blocks")
}

func testLinkUnclosedFence(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("AGENTS.md", "# Project rules\n\nBroken docs with an unclosed fence:\n\n```\n<!-- bench:start -->\n<!-- bench:end -->\n\nKEEP-ME text after the unclosed fence.\n")

	probe := f.Bench("link")

	if probe.ExitCode == 0 {
		t.Fatal("link succeeded despite an unclosed fence around Bench markers")
	}
	probe.RequireContains(strings.ToLower(probe.Stderr+probe.Stdout), "fence")
	requireFixtureFileContains(t, f, "AGENTS.md", "KEEP-ME text after the unclosed fence.", "unclosed-fence failure rewrote project text")
	requireLiteralCount(t, f, "AGENTS.md", "## Bench", 0, "unclosed-fence link still installed a managed block")
}

func testLinkRollsBackOnFault(t *testing.T) {
	for _, k := range []string{"1", "last"} {
		t.Run("promotion-"+k, func(t *testing.T) {
			f := contract.NewFixture(t)
			if k == "last" {
				f.Git("config", "core.hooksPath", ".husky")
			}
			before := fixtureState(t, f)
			f.BenchEnv(map[string]string{"BENCH_LINK_FAULT": k}, "link").RequireExit(1)
			if after := fixtureState(t, f); after != before {
				t.Fatalf("fault %s left repository residue\nbefore:\n%s\nafter:\n%s", k, before, after)
			}
			if k == "last" {
				requireLinkNotExists(t, f, ".husky", "last promotion fault left the newly-created hooks directory")
			}
		})
	}
}

func testLinkRelinkRollsBackOnFault(t *testing.T) {
	root := filepath.Join(t.TempDir(), "repo [1]")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	f := contract.NewFixtureAt(t, root, contract.IsolatedEnv(t, root))
	f.Git("init", "-q")
	kitA, kitB := lifecycleKits(t, f)
	f.BenchEnv(map[string]string{"BENCH_KIT": kitA}, "link").RequireExit(0)
	contract.WriteFileAbs(t, filepath.Join(kitB, ".bench", "BENCH.md"), "replacement from kit B\n")
	before := fixtureState(t, f)
	f.BenchEnv(map[string]string{"BENCH_KIT": kitB, "BENCH_LINK_FAULT": "last"}, "link").RequireExit(1)
	if after := fixtureState(t, f); after != before {
		t.Fatalf("relink fault left repository residue\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func fixtureState(t *testing.T, f contract.Fixture) string {
	t.Helper()
	var rows []string
	err := filepath.Walk(f.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(f.Root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			rows = append(rows, fmt.Sprintf("%s mode=%#o symlink:%s", rel, info.Mode(), target))
		} else if info.Mode().IsRegular() {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rows = append(rows, fmt.Sprintf("%s mode=%#o bytes=%x", rel, info.Mode(), b))
		} else if info.IsDir() {
			rows = append(rows, fmt.Sprintf("%s mode=%#o dir", rel, info.Mode()))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(rows)
	return strings.Join(rows, "\n") + "\nstatus=" + f.Git("status", "--porcelain=v1").Stdout
}

func testLinkReconcilesKitVersions(t *testing.T) {
	f := contract.NewFixture(t)
	kitA, kitB := lifecycleKits(t, f)
	f.BenchEnv(map[string]string{"BENCH_KIT": kitA}, "link").RequireExit(0)
	f.WriteFile(".bench/link-manifest.tsv", strings.Replace(f.ReadFile(".bench/link-manifest.tsv"), "#kit\t", "#kit\tsentinel-a-", 1))
	f.BenchEnv(map[string]string{"BENCH_KIT": kitB}, "link").RequireExit(0)
	requireLinkNotExists(t, f, ".agents/commands/lifecycle-x.md", "clean dropped asset survived relink")
	requireLinkFile(t, f, ".agents/commands/lifecycle-y.md")
	requireFixtureFileNotContains(t, f, ".bench/link-manifest.tsv", "lifecycle-x.md\t", "dropped asset stayed manifest-owned")
	requireFixtureFileContains(t, f, ".bench/link-manifest.tsv", "lifecycle-y.md\t", "added asset is missing from manifest")
	requireFixtureFileNotContains(t, f, ".bench/link-manifest.tsv", "sentinel-a-", "relink did not restamp kit version")
}

func testLinkKeepsModifiedDroppedAsset(t *testing.T) {
	f := contract.NewFixture(t)
	kitA, kitB := lifecycleKits(t, f)
	f.BenchEnv(map[string]string{"BENCH_KIT": kitA}, "link").RequireExit(0)
	oldRow := manifestRow(t, f, ".agents/commands/lifecycle-x.md")
	f.WriteFile(".agents/commands/lifecycle-x.md", "asset x, locally changed\n")
	probe := f.BenchEnv(map[string]string{"BENCH_KIT": kitB}, "link")
	probe.RequireExit(3)
	requireLinkEqual(t, f.ReadFile(".agents/commands/lifecycle-x.md"), "asset x, locally changed\n", "modified dropped asset bytes changed")
	if got := manifestRow(t, f, ".agents/commands/lifecycle-x.md"); got != oldRow {
		t.Fatalf("modified dropped asset manifest row = %q, want old row %q", got, oldRow)
	}
	probe.RequireContains(probe.Stdout, "conflicts[1]{path,reason}:")
	probe.RequireContains(probe.Stdout, "  .agents/commands/lifecycle-x.md,kept-modified-removed")
}

func testLinkRejectsHostileDroppedRows(t *testing.T) {
	for _, hostile := range []struct {
		name string
		seed func(*testing.T, contract.Fixture)
	}{
		{"traversal", func(t *testing.T, f contract.Fixture) {
			contract.WriteFileAbs(t, filepath.Join(filepath.Dir(f.Root), "relink-outside-sentinel"), "do not touch\n")
			appendFile(t, filepath.Join(f.Root, ".bench", "link-manifest.tsv"), "../relink-outside-sentinel\tdeadbeef\n")
		}},
		{"fifo", func(t *testing.T, f contract.Fixture) {
			path := filepath.Join(f.Root, ".agents", "commands", "lifecycle-x.md")
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := syscall.Mkfifo(path, 0o644); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(hostile.name, func(t *testing.T) {
			f := contract.NewFixture(t)
			kitA, kitB := hostileLifecycleKits(t, f)
			f.BenchEnv(map[string]string{"BENCH_KIT": kitA}, "link").RequireExit(0)
			before := f.ReadFile(".bench/link-manifest.tsv")
			hostile.seed(t, f)
			probe := contract.RunAtWithTimeout(t, f, f.Root, map[string]string{"BENCH_KIT": kitB}, time.Second, "bash", filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), "link")
			if probe.TimedOut {
				t.Fatal("relink blocked on hostile dropped target")
			}
			probe.RequireExit(1)
			requireLinkNotExists(t, f, ".agents/commands/lifecycle-y.md", "hard reconcile failure wrote planned asset")
			if hostile.name == "traversal" && !strings.Contains(contract.ReadFileAbs(t, filepath.Join(filepath.Dir(f.Root), "relink-outside-sentinel")), "do not touch") {
				t.Fatal("relink touched outside sentinel")
			}
			if hostile.name == "fifo" {
				if info, err := os.Lstat(filepath.Join(f.Root, ".agents", "commands", "lifecycle-x.md")); err != nil || info.Mode()&os.ModeNamedPipe == 0 {
					t.Fatal("relink deleted or replaced FIFO")
				}
			}
			if got := f.ReadFile(".bench/link-manifest.tsv"); got != before && hostile.name == "fifo" {
				t.Fatal("FIFO hard failure rewrote manifest")
			}
		})
	}
}

func testLinkLifecycleMatrix(t *testing.T) {
	f := contract.NewFixture(t)
	for cycle := 1; cycle <= 2; cycle++ {
		linkOK(t, f)
		unlinkOK(t, f)
	}
	requireLinkNotExists(t, f, ".bench/link-manifest.tsv", "second link/unlink cycle left manifest residue")
	requireLinkNotExists(t, f, managedFileRel, "second link/unlink cycle left managed asset residue")
	relink := contract.NewFixture(t)
	linkOK(t, relink)
	linkOK(t, relink)
	if manifestRowCount(t, relink, managedFileRel) != 1 {
		t.Fatal("second relink duplicated manifest ownership")
	}
}

// lifecycleSharedRel is the managed asset the two lifecycle kits both ship with
// different bytes — the shape a real release takes when it edits an asset a consumer
// already owns and has never touched.
const (
	lifecycleSharedRel = ".agents/commands/lifecycle-shared.md"
	lifecycleSharedA   = "shared asset, kit A\n"
	lifecycleSharedB   = "shared asset, kit B, revised\n"
)

// TestLinkCleanSkipPropagationContracts pins that an untouched managed asset tracks the
// kit it is linked against rather than the bytes it was first written with. Skipping a
// destination because it matches the recorded manifest hash makes every release a
// content no-op for exactly the files a consumer never touched, and leaves the manifest
// asserting ownership of bytes no kit ships.
func TestLinkCleanSkipPropagationContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench link clean-entry propagation contract failed", testLinkPropagatesChangedKitBytes)
	contract.RunParallel(t, "bench upgrade clean-entry propagation contract failed", testUpgradePropagatesChangedKitBytes)
}

func testLinkPropagatesChangedKitBytes(t *testing.T) {
	f := contract.NewFixture(t)
	kitA, kitB := lifecycleKits(t, f)
	f.BenchEnv(map[string]string{"BENCH_KIT": kitA}, "link").RequireExit(0)
	requireLinkEqual(t, f.ReadFile(lifecycleSharedRel), lifecycleSharedA, "link did not install the kit A shared asset")

	f.BenchEnv(map[string]string{"BENCH_KIT": kitB}, "link").RequireExit(0)

	requireLinkEqual(t, f.ReadFile(lifecycleSharedRel), lifecycleSharedB, "relink left an untouched clean asset at the previous kit's bytes")
	requireManifestHash(t, f, lifecycleSharedRel, lifecycleSharedB, "relink kept the previous kit's manifest hash for a rewritten asset")
}

func testUpgradePropagatesChangedKitBytes(t *testing.T) {
	f := contract.NewFixture(t)
	kitA, kitB := lifecycleKits(t, f)
	f.BenchEnv(map[string]string{"BENCH_KIT": kitA}, "link").RequireExit(0)
	repinManifestKitVersion(t, f, "0.0.1")

	f.BenchEnv(map[string]string{"BENCH_KIT": kitB}, "upgrade").RequireExit(0)

	requireLinkEqual(t, f.ReadFile(lifecycleSharedRel), lifecycleSharedB, "upgrade left an untouched clean asset at the previous kit's bytes")
	requireManifestHash(t, f, lifecycleSharedRel, lifecycleSharedB, "upgrade kept the previous kit's manifest hash for a rewritten asset")
}

// requireManifestHash asserts the manifest owns rel at exactly the hash content
// fingerprints to, which is what separates a propagated rewrite from a retained row.
func requireManifestHash(t *testing.T, f contract.Fixture, rel, content, msg string) {
	t.Helper()
	want := rel + "\t" + fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	if got := manifestRow(t, f, rel); got != want {
		t.Fatalf("%s: manifest row = %q, want %q", msg, got, want)
	}
}

func lifecycleKits(t *testing.T, f contract.Fixture) (string, string) {
	t.Helper()
	root := contract.KitRoot(t)
	makeKit := func(name string) string {
		kit := filepath.Join(f.Root, name)
		contract.Mkdir(t, filepath.Join(kit, ".bench"))
		contract.Mkdir(t, filepath.Join(kit, "dist"))
		copyPaths(t, kit, filepath.Join(root, "bin"), filepath.Join(root, ".agents"), filepath.Join(root, ".claude"), filepath.Join(root, ".codex"))
		copyFileTo(t, filepath.Join(root, ".bench", "BENCH.md"), filepath.Join(kit, ".bench", "BENCH.md"))
		copyFileTo(t, filepath.Join(root, ".bench", "BENCH-reference.md"), filepath.Join(kit, ".bench", "BENCH-reference.md"))
		copyFileTo(t, filepath.Join(root, "dist", "bench"), filepath.Join(kit, "dist", "bench"))
		copyPaths(t, filepath.Join(kit, ".bench"), filepath.Join(root, ".bench", "hooks"), filepath.Join(root, ".bench", "adapters"), filepath.Join(root, ".bench", "lib"))
		return kit
	}
	kitA, kitB := makeKit("kit-a"), makeKit("kit-b")
	contract.WriteFileAbs(t, filepath.Join(kitA, ".agents", "commands", "lifecycle-x.md"), "asset x\n")
	contract.WriteFileAbs(t, filepath.Join(kitB, ".agents", "commands", "lifecycle-y.md"), "asset y\n")
	contract.WriteFileAbs(t, filepath.Join(kitA, filepath.FromSlash(lifecycleSharedRel)), lifecycleSharedA)
	contract.WriteFileAbs(t, filepath.Join(kitB, filepath.FromSlash(lifecycleSharedRel)), lifecycleSharedB)
	return kitA, kitB
}

// hostileLifecycleKits contains only the plan's fixed files and the two command
// inputs needed to make X a dropped managed target and Y a new one. This keeps the
// FIFO deadline about classification, not unrelated full-kit staging work.
func hostileLifecycleKits(t *testing.T, f contract.Fixture) (string, string) {
	t.Helper()
	root := contract.KitRoot(t)
	makeKit := func(name string) string {
		kit := filepath.Join(f.Root, name)
		for _, rel := range []string{
			".bench/BENCH.md", ".bench/BENCH-reference.md", ".claude/README.md",
			".claude/settings.json", ".codex/hooks.json",
		} {
			copyFileTo(t, filepath.Join(root, filepath.FromSlash(rel)), filepath.Join(kit, filepath.FromSlash(rel)))
		}
		copyFileTo(t, filepath.Join(root, "dist", "bench"), filepath.Join(kit, "dist", "bench"))
		return kit
	}
	kitA, kitB := makeKit("hostile-kit-a"), makeKit("hostile-kit-b")
	contract.WriteFileAbs(t, filepath.Join(kitA, ".agents", "commands", "lifecycle-x.md"), "asset x\n")
	contract.WriteFileAbs(t, filepath.Join(kitB, ".agents", "commands", "lifecycle-y.md"), "asset y\n")
	return kitA, kitB
}

func manifestRow(t *testing.T, f contract.Fixture, rel string) string {
	t.Helper()
	for _, line := range strings.Split(f.ReadFile(".bench/link-manifest.tsv"), "\n") {
		if strings.HasPrefix(line, rel+"\t") {
			return line
		}
	}
	t.Fatalf("missing manifest row for %s", rel)
	return ""
}
