package surface

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
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
	contract.RunParallel(t, "bench upgrade unlinked-repo contract failed", testUpgradeRefusesWithoutManifest)
	contract.RunParallel(t, "bench upgrade same-version no-op contract failed", testUpgradeSameVersionIsNoOp)
	contract.RunParallel(t, "bench upgrade downgrade refusal contract failed", testUpgradeRefusesDowngrade)
	contract.RunParallel(t, "bench upgrade invocation-path contract failed", testUpgradeResolvesThroughLauncherAndSymlink)
	contract.RunParallel(t, "bench upgrade subdirectory contract failed", testUpgradeFromSubdirectory)
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
	f.WriteFile(".bench/link-manifest.tsv", strings.Replace(f.ReadFile(".bench/link-manifest.tsv"), "#kit\t"+installed+"\n", "#kit\t"+older+"\n", 1))
	return installed
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

func testUpgradeChecksThenApplies(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "0.0.1")

	before := fixtureState(t, f)
	check := f.Bench("upgrade", "--check")
	check.RequireExit(0)
	requireUpgradePlanRow(t, check, "0.0.1", installed)
	if after := fixtureState(t, f); after != before {
		t.Fatalf("bench upgrade --check wrote to the repository\nbefore:\n%s\nafter:\n%s", before, after)
	}

	apply := f.Bench("upgrade")
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

	absent := f.Bench("upgrade")
	absent.RequireExit(1)
	absent.RequireContains(absent.Stderr, "bench link")

	f.WriteFile(".bench/link-manifest.tsv", "")
	empty := f.Bench("upgrade")
	empty.RequireExit(1)
	empty.RequireContains(empty.Stderr, "empty")
	empty.RequireNotContains(empty.Stderr, "run 'bench link' first")
}

func testUpgradeSameVersionIsNoOp(t *testing.T) {
	f := contract.NewFixture(t)
	installed := linkedRepoAtOlderVersion(t, f, "0.0.1")
	f.Bench("upgrade").RequireExit(0)

	before := fixtureState(t, f)
	again := f.Bench("upgrade")
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
	refused := f.Bench("upgrade")
	refused.RequireExit(1)
	refused.RequireContains(refused.Stderr, "--force")
	if after := fixtureState(t, f); after != before {
		t.Fatalf("a refused downgrade touched the repository\nbefore:\n%s\nafter:\n%s", before, after)
	}

	forced := f.Bench("upgrade", "--force")
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

	direct := f.Bench("upgrade", "--check")
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

	probe := f.Bench("upgrade")

	probe.RequireExit(0)
	requireUpgradePlanRow(t, probe, "0.0.1", installed)
	requireLinkNotExists(t, f, preExclusionKitOnlyRel, "upgrade left a pre-exclusion kit-only asset in the tree")
	requireFixtureFileNotContains(t, f, ".bench/link-manifest.tsv", preExclusionKitOnlyRel+"\t", "upgrade left a pre-exclusion kit-only asset manifest-owned")
}

func testUpgradeKeepsModifiedManagedAsset(t *testing.T) {
	f := contract.NewFixture(t)
	linkedRepoAtOlderVersion(t, f, "0.0.1")
	f.WriteFile(".agents/commands/bench-implement-spec.md", "consumer's own edit\n")

	probe := f.Bench("upgrade")

	probe.RequireExit(3)
	probe.RequireContains(probe.Stdout, "  .agents/commands/bench-implement-spec.md,modified-managed")
	requireLinkEqual(t, f.ReadFile(".agents/commands/bench-implement-spec.md"), "consumer's own edit\n", "upgrade overwrote a modified consumer-owned asset")
}
