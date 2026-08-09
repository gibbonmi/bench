package surface

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestUpgradeContracts pins `bench upgrade` as the consumer's supported route onto a
// newer kit (FT85 story 4): it reports a from/to plan, writes nothing under --check,
// applies the existing transactional relink otherwise, and refuses a downgrade by
// name. Every case drives the shipped wrapper against a throwaway linked repo, so the
// assertions are what a consumer sees rather than the shape of any internal plan.
func TestUpgradeContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench upgrade check-then-apply contract failed", testUpgradeChecksThenApplies)
	contract.RunParallel(t, "bench upgrade hook refresh plan contract failed", testUpgradePlansHookRefresh)
	contract.RunParallel(t, "bench upgrade live hook render plan contract failed", testUpgradePlansLiveHookRender)
	contract.RunParallel(t, "bench upgrade configured hooksPath refresh plan contract failed", testUpgradePlansAbsentConfiguredHookRefresh)
	contract.RunParallel(t, "bench upgrade configured hooksPath error plan contract failed", testUpgradeRejectsConfiguredHookPathErrors)
	contract.RunParallel(t, "bench upgrade special and dangling hook planning contract failed", testUpgradeRefusesUnwritableHooks)
	contract.RunParallel(t, "bench upgrade unlinked-repo contract failed", testUpgradeRefusesWithoutManifest)
	contract.RunParallel(t, "bench upgrade same-version no-op contract failed", testUpgradeSameVersionIsNoOp)
	contract.RunParallel(t, "bench upgrade downgrade refusal contract failed", testUpgradeRefusesDowngrade)
	contract.RunParallel(t, "bench upgrade invocation-path contract failed", testUpgradeResolvesThroughLauncherAndSymlink)
	contract.RunParallel(t, "bench upgrade subdirectory contract failed", testUpgradeFromSubdirectory)
	contract.RunParallel(t, "bench upgrade prerelease-to-release contract failed", testUpgradePrereleaseToRelease)
	contract.RunParallel(t, "bench upgrade argument contract failed", testUpgradeArgumentEdges)
	contract.RunParallel(t, "bench upgrade unusable-manifest contract failed", testUpgradeUnusableManifestStates)
}

// TestUpgradePayloadReconcileContracts pins FT85 story 5: upgrading across the release
// that withheld the kit-only surfaces removes the assets a pre-exclusion manifest still
// owns, from the tree and from the manifest, and a consumer's own edit to a managed
// asset survives with the partial-result posture the relink contract already defines.
func TestUpgradePayloadReconcileContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench upgrade kit-only withdrawal contract failed", testUpgradeRemovesPreExclusionKitOnlyAssets)
	contract.RunParallel(t, "bench upgrade modified-asset posture contract failed", testUpgradeKeepsModifiedManagedAsset)
}

const preExclusionKitOnlyRel = ".agents/skills/bench-craft-synthesis/SKILL.md"

// upgradePlanHeader is the TOON schema `bench upgrade` prints. Pinned as one literal so
// every case asserts the same block contract a consumer parses.
const upgradePlanHeader = "upgrade[1]{from,to,added,changed,removed}:"

// linkedRepoAtOlderVersion links f, then rewrites only the manifest's #kit header to
// older, which is exactly the state a repo linked by a previous release is in. It
// returns the installed kit version the header carried before the rewrite.
func linkedRepoAtOlderVersion(t *testing.T, f contract.Fixture, older string) string {
	t.Helper()
	linkOK(t, f)
	installed := manifestKitVersion(t, f)
	repinManifestKitVersion(t, f, older)
	return installed
}

// repinManifestKitVersion rewrites only the manifest's #kit header, leaving every owned
// row intact — the state a repo linked by another release is in.
func repinManifestKitVersion(t *testing.T, f contract.Fixture, pinned string) {
	t.Helper()
	current := manifestKitVersion(t, f)
	f.WriteFile(".bench/link-manifest.tsv", strings.Replace(f.ReadFile(".bench/link-manifest.tsv"), "#kit\t"+current+"\n", "#kit\t"+pinned+"\n", 1))
}

func manifestKitVersion(t *testing.T, f contract.Fixture) string {
	t.Helper()
	for _, line := range strings.Split(f.ReadFile(".bench/link-manifest.tsv"), "\n") {
		if strings.HasPrefix(line, "#kit\t") {
			return strings.TrimPrefix(line, "#kit\t")
		}
	}
	t.Fatal("link manifest carries no #kit version header")
	return ""
}

// requireUpgradePlanRow asserts the printed plan names this from/to pair. The three
// count columns are deliberately not pinned to exact numbers: they move with the kit's
// asset inventory, and the behavior under test is the version plan, not the census.
func requireUpgradePlanRow(t *testing.T, probe contract.Probe, from, to string) {
	t.Helper()
	probe.RequireContains(probe.Stdout, upgradePlanHeader)
	probe.RequireContains(probe.Stdout, "\n  "+from+","+to+",")
}

func upgradePlanChangedCount(t *testing.T, probe contract.Probe) int {
	t.Helper()
	for _, line := range strings.Split(probe.Stdout, "\n") {
		if !strings.HasPrefix(line, "  ") {
			continue
		}
		cells := strings.Split(strings.TrimSpace(line), ",")
		if len(cells) != 5 {
			continue
		}
		changed, err := strconv.Atoi(cells[3])
		if err != nil {
			t.Fatalf("upgrade plan changed count is not numeric: %q", line)
		}
		return changed
	}
	t.Fatalf("upgrade plan carries no data row:\n%s", probe.Stdout)
	return 0
}

func upgradeSubject(t *testing.T, f contract.Fixture, args ...string) contract.Probe {
	t.Helper()
	return f.BenchWrapperEnv(map[string]string{"BENCH_KIT": contract.SubjectRoot(t)}, args...)
}

func testUpgradePlansHookRefresh(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "0.0.1")

	current := upgradeSubject(t, f, "upgrade", "--check")
	current.RequireExit(0)
	currentCount := upgradePlanChangedCount(t, current)

	hook := prePushPath(t, f)
	managedHook := contract.ReadFileAbs(t, hook)
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	absent := upgradeSubject(t, f, "upgrade", "--check")
	absent.RequireExit(0)
	if got := upgradePlanChangedCount(t, absent); got != currentCount+1 {
		t.Fatalf("upgrade plan changed count: got %d, want %d with an absent hook refresh", got, currentCount+1)
	}
	contract.WriteFileAbs(t, hook, "#!/bin/sh\nexit 0\n")
	// A non-managed current file is refused rather than promised as a refresh.
	foreign := upgradeSubject(t, f, "upgrade", "--check")
	foreign.RequireExit(0)
	if got := upgradePlanChangedCount(t, foreign); got != currentCount {
		t.Fatalf("upgrade plan changed count: got %d, want %d with a refused hook", got, currentCount)
	}
	contract.WriteFileAbs(t, hook, managedHook)
	staleHook := contract.ReadFileAbs(t, hook) + "\n# stale hook\n"
	contract.WriteFileAbs(t, hook, staleHook)
	before := fixtureState(t, f)
	stale := upgradeSubject(t, f, "upgrade", "--check")
	stale.RequireExit(0)
	if got := upgradePlanChangedCount(t, stale); got != currentCount+1 {
		t.Fatalf("upgrade plan changed count: got %d, want %d after a stale hook refresh", got, currentCount+1)
	}
	if after := fixtureState(t, f); after != before {
		t.Fatalf("bench upgrade --check refreshed the hook\nbefore:\n%s\nafter:\n%s", before, after)
	}

	apply := upgradeSubject(t, f, "upgrade")
	apply.RequireExit(0)
	requireLinkEqual(t, manifestKitVersion(t, f), installed, "upgrading did not restamp the manifest kit version")
	contract.WriteFileAbs(t, hook, staleHook)
	before = fixtureState(t, f)
	equal := upgradeSubject(t, f, "upgrade", "--check")
	equal.RequireExit(0)
	if got := upgradePlanChangedCount(t, equal); got != 0 {
		t.Fatalf("equal-version upgrade plan counted stale hook refresh: got %d, want 0", got)
	}
	if after := fixtureState(t, f); after != before {
		t.Fatalf("equal-version bench upgrade --check refreshed the hook\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func testUpgradePlansLiveHookRender(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	f.Run("git", "init", "-q", "--bare", "-b", "master", "remote.git").RequireExit(0)
	f.Run("git", "init", "-q", "-b", "master", "repo").RequireExit(0)
	repo := linkFixtureAt(t, filepath.Join(f.Root, "repo"), f.Env)
	repo.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "--allow-empty", "-m", "master")
	repo.Git("remote", "add", "origin", filepath.Join(f.Root, "remote.git"))
	repo.Git("push", "-q", "origin", "master")
	repo.Git("remote", "set-head", "origin", "--auto")

	installed := linkedRepoAtOlderVersion(t, repo, "0.0.1")
	baseline := upgradeSubject(t, repo, "upgrade", "--check")
	baseline.RequireExit(0)
	count := upgradePlanChangedCount(t, baseline)

	repo.Git("checkout", "-q", "-b", "main")
	repo.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "--allow-empty", "-m", "main")
	repo.Git("push", "-q", "origin", "main")
	f.Run("git", "-C", filepath.Join(f.Root, "remote.git"), "symbolic-ref", "HEAD", "refs/heads/main").RequireExit(0)

	check := upgradeSubject(t, repo, "upgrade", "--check")
	check.RequireExit(0)
	requireUpgradePlanRow(t, check, "0.0.1", installed)
	if got := upgradePlanChangedCount(t, check); got != count+1 {
		t.Fatalf("upgrade plan changed count: got %d, want %d when the remote default branch changes", got, count+1)
	}
	requireLinkEqual(t, manifestKitVersion(t, repo), "0.0.1", "bench upgrade --check restamped the manifest kit version")
	requireLinkEqual(t, strings.TrimSpace(repo.Git("symbolic-ref", "--short", "refs/remotes/origin/HEAD").Stdout), "origin/master", "local origin/HEAD moved while the remote default branch changed")
}

func testUpgradeRefusesUnwritableHooks(t *testing.T) {
	f := contract.NewFixture(t)
	linkedRepoAtOlderVersion(t, f, "0.0.1")
	hook := prePushPath(t, f)
	baseline := upgradeSubject(t, f, "upgrade", "--check")
	baseline.RequireExit(0)
	count := upgradePlanChangedCount(t, baseline)

	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	f.WriteFifo(filepath.ToSlash(strings.TrimPrefix(hook, f.Root+string(os.PathSeparator))))
	fifo := f.BenchWrapperDeadlined("upgrade", "--check")
	if fifo.TimedOut {
		t.Fatal("bench upgrade --check blocked on a writerless FIFO at the effective pre-push path")
	}
	fifo.RequireExit(0)
	if got := upgradePlanChangedCount(t, fifo); got != count {
		t.Fatalf("upgrade plan changed count: got %d, want %d with a refused FIFO hook", got, count)
	}
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing-pre-push", hook); err != nil {
		t.Fatal(err)
	}
	check := upgradeSubject(t, f, "upgrade", "--check")
	check.RequireExit(0)
	if got := upgradePlanChangedCount(t, check); got != count {
		t.Fatalf("upgrade plan changed count: got %d, want %d with a dangling hook", got, count)
	}
	apply := upgradeSubject(t, f, "upgrade")
	apply.RequireExit(1)
	apply.RequireContains(apply.Stderr, "conflict")
}

func testUpgradePlansAbsentConfiguredHookRefresh(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("config", "core.hooksPath", ".husky")
	linkedRepoAtOlderVersion(t, f, "0.0.1")
	hook := filepath.Join(f.Root, ".husky", "pre-push")
	baseline := upgradeSubject(t, f, "upgrade", "--check")
	baseline.RequireExit(0)
	count := upgradePlanChangedCount(t, baseline)
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}

	check := upgradeSubject(t, f, "upgrade", "--check")
	check.RequireExit(0)
	if got := upgradePlanChangedCount(t, check); got != count+1 {
		t.Fatalf("upgrade plan changed count: got %d, want %d with an absent configured hook refresh", got, count+1)
	}
	apply := upgradeSubject(t, f, "upgrade")
	apply.RequireExit(0)
	requireExecutable(t, hook, "upgrade did not restore the absent configured pre-push hook")
}

func testUpgradeRejectsConfiguredHookPathErrors(t *testing.T) {
	f := contract.NewFixture(t)
	f.Git("config", "core.hooksPath", ".husky")
	linkedRepoAtOlderVersion(t, f, "0.0.1")

	baseline := upgradeSubject(t, f, "upgrade", "--check")
	baseline.RequireExit(0)
	count := upgradePlanChangedCount(t, baseline)

	hook := filepath.Join(f.Root, ".husky", "pre-push")
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Dir(hook)); err != nil {
		t.Fatal(err)
	}
	contract.WriteFileAbs(t, filepath.Dir(hook), "not a hooks directory\n")

	check := upgradeSubject(t, f, "upgrade", "--check")
	check.RequireExit(0)
	if got := upgradePlanChangedCount(t, check); got != count {
		t.Fatalf("upgrade plan changed count: got %d, want %d with an unreadable configured pre-push path", got, count)
	}
}

func testUpgradeChecksThenApplies(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "0.0.1")

	before := fixtureState(t, f)
	check := f.BenchWrapper("upgrade", "--check")
	check.RequireExit(0)
	requireUpgradePlanRow(t, check, "0.0.1", installed)
	if after := fixtureState(t, f); after != before {
		t.Fatalf("bench upgrade --check wrote to the repository\nbefore:\n%s\nafter:\n%s", before, after)
	}

	apply := f.BenchWrapper("upgrade")
	apply.RequireExit(0)
	requireUpgradePlanRow(t, apply, "0.0.1", installed)
	requireLinkEqual(t, manifestKitVersion(t, f), installed, "applying upgrade did not restamp the manifest kit version")
	requireLinkFile(t, f, ".agents/commands/bench-implement-spec.md")
}

// testUpgradeRefusesWithoutManifest covers the absent-vs-present-but-empty edge: no
// manifest is "this repo was never linked" and names `bench link`; an empty manifest is
// a distinct malformed state, because a repo that has a manifest file but no version
// header needs a different remedy than one that has none.
func testUpgradeRefusesWithoutManifest(t *testing.T) {
	f := contract.NewFixture(t)

	absent := f.BenchWrapper("upgrade")
	absent.RequireExit(1)
	absent.RequireContains(absent.Stderr, "bench link")

	f.WriteFile(".bench/link-manifest.tsv", "")
	empty := f.BenchWrapper("upgrade")
	empty.RequireExit(1)
	empty.RequireContains(empty.Stderr, "empty")
	empty.RequireNotContains(empty.Stderr, "run 'bench link' first")
}

func testUpgradeSameVersionIsNoOp(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "0.0.1")
	f.BenchWrapper("upgrade").RequireExit(0)

	before := fixtureState(t, f)
	again := f.BenchWrapper("upgrade")
	again.RequireExit(0)
	requireUpgradePlanRow(t, again, installed, installed)
	if after := fixtureState(t, f); after != before {
		t.Fatalf("a second bench upgrade at the same version rewrote the repository\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func testUpgradeRefusesDowngrade(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "99.0.0")

	before := fixtureState(t, f)
	refused := f.BenchWrapper("upgrade")
	refused.RequireExit(1)
	refused.RequireContains(refused.Stderr, "--force")
	if after := fixtureState(t, f); after != before {
		t.Fatalf("a refused downgrade touched the repository\nbefore:\n%s\nafter:\n%s", before, after)
	}

	forced := f.BenchWrapper("upgrade", "--force")
	forced.RequireExit(0)
	requireUpgradePlanRow(t, forced, "99.0.0", installed)
	requireLinkEqual(t, manifestKitVersion(t, f), installed, "forced downgrade did not restamp the manifest kit version")
}

// testUpgradeResolvesThroughLauncherAndSymlink covers the by-path invocation edge: a
// symlink to the kit wrapper and the linked repo's own .bench/bin launcher both reach
// the same implementation and print the same plan. Without a kit to relink from, the
// repo-local launcher fails closed exactly as the other adoption routes do rather than
// planning against an empty asset tree.
func testUpgradeResolvesThroughLauncherAndSymlink(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "0.0.1")
	kit := contract.SubjectRoot(t)

	direct := f.BenchWrapper("upgrade", "--check")
	direct.RequireExit(0)

	symlink := filepath.Join(f.Root, "bench-symlink.sh")
	if err := os.Symlink(filepath.Join(kit, "bin", "bench.sh"), symlink); err != nil {
		t.Fatal(err)
	}
	viaSymlink := f.Run("bash", symlink, "upgrade", "--check")
	viaSymlink.RequireExit(0)
	requireLinkEqual(t, viaSymlink.Stdout, direct.Stdout, "upgrade through a symlinked wrapper resolved a different implementation")

	bare := f.Run(filepath.Join(".bench", "bin", "bench.sh"), "upgrade", "--check")
	if bare.ExitCode == 0 {
		t.Fatal("the repo-local launcher planned an upgrade with no kit asset tree to relink from")
	}
	bare.RequireContains(bare.Stderr, "real Bench kit")

	viaLauncher := f.RunEnv(map[string]string{"BENCH_KIT": kit}, filepath.Join(".bench", "bin", "bench.sh"), "upgrade", "--check")
	viaLauncher.RequireExit(0)
	requireUpgradePlanRow(t, viaLauncher, "0.0.1", installed)
	requireLinkEqual(t, viaLauncher.Stdout, direct.Stdout, "upgrade through the repo-local launcher resolved a different implementation")
}

func testUpgradeFromSubdirectory(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "0.0.1")
	deep := filepath.Join(f.Root, "src", "nested")
	contract.Mkdir(t, deep)

	probe := contract.RunAt(t, f, deep, nil, "bash", filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh"), "upgrade")

	probe.RequireExit(0)
	requireUpgradePlanRow(t, probe, "0.0.1", installed)
	requireLinkEqual(t, manifestKitVersion(t, f), installed, "upgrade from a subdirectory did not resolve the repo root")
}

// testUpgradePrereleaseToRelease covers the release a consumer is likeliest to be
// pinned behind: the prerelease of the version they now have installed. Ordering by
// release components alone made this pair compare equal, so upgrade printed a plan,
// exited 0, and left the manifest stamped at the prerelease — a silent no-op that looked
// like success. It must relink and restamp instead.
func testUpgradePrereleaseToRelease(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	installed := manifestKitVersion(t, f)
	prerelease := installed + "-rc1"
	repinManifestKitVersion(t, f, prerelease)

	probe := f.BenchWrapper("upgrade")

	probe.RequireExit(0)
	requireUpgradePlanRow(t, probe, prerelease, installed)
	requireLinkEqual(t, manifestKitVersion(t, f), installed, "upgrading from a prerelease of the installed version left the manifest stamped at the prerelease")
}

// testUpgradeArgumentEdges pins the two argument shapes a consumer reaches by accident:
// an unrecognized flag is a usage error rather than a silently ignored word, and
// --check with --force still writes nothing, because --force widens what an applying run
// will do and never turns a dry run into one.
func testUpgradeArgumentEdges(t *testing.T) {
	f := contract.NewFixture(t)
	linkedRepoAtOlderVersion(t, f, "99.0.0")

	unknown := f.BenchWrapper("upgrade", "--dry-run")
	unknown.RequireExit(2)
	unknown.RequireContains(unknown.Stderr, "usage: bench upgrade")

	before := fixtureState(t, f)
	checkForce := f.BenchWrapper("upgrade", "--check", "--force")
	checkForce.RequireExit(0)
	if after := fixtureState(t, f); after != before {
		t.Fatalf("bench upgrade --check --force wrote to the repository\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// testUpgradeUnusableManifestStates covers the two manifests a hand edit produces. A
// version string nothing can parse must not read as a downgrade — refusing there would
// strand the repo with no route forward — so upgrade applies and restamps. A manifest
// that cannot be read at all fails closed by name instead of being treated as unlinked.
func testUpgradeUnusableManifestStates(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "not-a-version")

	applied := f.BenchWrapper("upgrade")
	applied.RequireExit(0)
	requireUpgradePlanRow(t, applied, "not-a-version", installed)
	requireLinkEqual(t, manifestKitVersion(t, f), installed, "upgrading from an unparseable pinned version did not restamp the manifest")

	if os.Geteuid() == 0 {
		return
	}
	manifest := filepath.Join(f.Root, ".bench", "link-manifest.tsv")
	if err := os.Chmod(manifest, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(manifest, 0o644)
	unreadable := f.BenchWrapper("upgrade")
	unreadable.RequireExit(1)
	unreadable.RequireContains(unreadable.Stderr, "unreadable")
	unreadable.RequireNotContains(unreadable.Stderr, "run 'bench link' first")
}

// testUpgradeRemovesPreExclusionKitOnlyAssets is the FT85 story 5 red signal. The
// fixture reconstructs what a repo linked before the exclusion looks like — the
// kit-only asset present on disk and owned by the manifest — and the manifest's last
// row deliberately lacks a trailing newline, the hand-edited-file edge. Upgrading must
// withdraw the asset from both the tree and the manifest.
func testUpgradeRemovesPreExclusionKitOnlyAssets(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "0.0.1")

	kitOnly := contract.ReadFileAbs(t, filepath.Join(contract.SubjectRoot(t), filepath.FromSlash(preExclusionKitOnlyRel)))
	f.WriteFile(preExclusionKitOnlyRel, kitOnly)
	f.WriteFile(".bench/link-manifest.tsv", f.ReadFile(".bench/link-manifest.tsv")+
		preExclusionKitOnlyRel+"\t"+fmt.Sprintf("%x", sha256.Sum256([]byte(kitOnly))))

	probe := f.BenchWrapper("upgrade")

	probe.RequireExit(0)
	requireUpgradePlanRow(t, probe, "0.0.1", installed)
	requireLinkNotExists(t, f, preExclusionKitOnlyRel, "upgrade left a pre-exclusion kit-only asset in the tree")
	requireFixtureFileNotContains(t, f, ".bench/link-manifest.tsv", preExclusionKitOnlyRel+"\t", "upgrade left a pre-exclusion kit-only asset manifest-owned")
}

func testUpgradeKeepsModifiedManagedAsset(t *testing.T) {
	f := contract.NewFixture(t)
	linkedRepoAtOlderVersion(t, f, "0.0.1")
	f.WriteFile(".agents/commands/bench-implement-spec.md", "consumer's own edit\n")

	probe := f.BenchWrapper("upgrade")

	probe.RequireExit(3)
	probe.RequireContains(probe.Stdout, "  .agents/commands/bench-implement-spec.md,modified-managed")
	requireLinkEqual(t, f.ReadFile(".agents/commands/bench-implement-spec.md"), "consumer's own edit\n", "upgrade overwrote a modified consumer-owned asset")
}
