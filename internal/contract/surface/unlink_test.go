package surface

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestUnlinkContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench unlink removed-tree contract failed", testUnlinkRemovesManagedTree)
	contract.RunParallel(t, "bench unlink kept-modified contract failed", testUnlinkKeepsModified)
	contract.RunParallel(t, "bench unlink managed/foreign pre-push contract failed", testUnlinkPrePushHook)
	contract.RunParallel(t, "bench unlink AGENTS.md fence contract failed", testUnlinkStripsAgentsFence)
	contract.RunParallel(t, "bench unlink user-artifact contract failed", testUnlinkSparesUserArtifacts)
	contract.RunParallel(t, "bench unlink manifest-last contract failed", testUnlinkManifestConditional)
	contract.RunParallel(t, "bench unlink traversal-refusal contract failed", testUnlinkRefusesTraversal)
	contract.RunParallel(t, "bench unlink --dry-run no-writes contract failed", testUnlinkDryRunNoWrites)
	contract.RunParallel(t, "bench unlink report contract failed", testUnlinkReportsBothRuns)
	contract.RunParallel(t, "bench unlink absent-manifest contract failed", testUnlinkAbsentManifest)
	contract.RunParallel(t, "bench unlink by-path CLI route contract failed", testUnlinkByPathRoute)
	contract.RunParallel(t, "bench unlink link-created CLAUDE.md contract failed", testUnlinkRemovesLinkCreatedClaudeMD)
	contract.RunParallel(t, "bench unlink pre-existing CLAUDE.md contract failed", testUnlinkSparesExistingClaudeMD)
}

const managedFileRel = ".agents/commands/bench-implement-spec.md"

// Story 1: link then unlink leaves only user content; managed files and now-empty managed
// directories are gone.
func testUnlinkRemovesManagedTree(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	requireLinkFile(t, f, managedFileRel)

	unlinkOK(t, f)

	requireLinkNotExists(t, f, managedFileRel, "unlink left a managed file on disk")
	requireLinkNotExists(t, f, ".agents/commands", "unlink left the emptied managed directory behind")
	requireLinkNotExists(t, f, ".agents", "unlink left an emptied managed directory tree behind")
}

// Story 2: a managed file edited after link is kept and reported, not deleted.
func testUnlinkKeepsModified(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	appendFile(t, filepath.Join(f.Root, filepath.FromSlash(managedFileRel)), "\nlocal edit\n")

	probe := f.Bench("unlink")
	probe.RequireExit(0)

	requireFixtureFileContains(t, f, managedFileRel, "local edit", "unlink deleted a locally modified managed file")
	probe.RequireContains(probe.Stdout, "modified")
	probe.RequireContains(probe.Stdout, managedFileRel)
}

// Story 3: the managed pre-push hook is removed; a foreign one is left.
func testUnlinkPrePushHook(t *testing.T) {
	managed := contract.NewFixture(t)
	linkOK(t, managed)
	managedHook := prePushPath(t, managed)
	if !managed.Exists(relTo(t, managed, managedHook)) {
		t.Fatal("link did not install a managed pre-push hook to remove")
	}
	unlinkOK(t, managed)
	if managed.Exists(relTo(t, managed, managedHook)) {
		t.Fatal("unlink left the bench-managed pre-push hook in place")
	}

	foreign := contract.NewFixture(t)
	linkOK(t, foreign)
	foreignHook := prePushPath(t, foreign)
	contract.WriteExecutableAbs(t, foreignHook, "#!/bin/sh\n# a hook I installed myself\nexit 0\n")
	unlinkOK(t, foreign)
	if !foreign.Exists(relTo(t, foreign, foreignHook)) {
		t.Fatal("unlink clobbered a foreign (non-managed) pre-push hook")
	}
	if got := contract.ReadFileAbs(t, foreignHook); !strings.Contains(got, "a hook I installed myself") {
		t.Fatalf("unlink altered a foreign pre-push hook:\n%s", got)
	}
}

// Story 4: the AGENTS.md fenced block is stripped while user prose survives.
func testUnlinkStripsAgentsFence(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("AGENTS.md", "PROJECT RULES\n")
	linkOK(t, f)
	requireFixtureFileContains(t, f, "AGENTS.md", "<!-- bench:start -->", "link did not install a managed AGENTS.md block")

	unlinkOK(t, f)

	requireFixtureFileContains(t, f, "AGENTS.md", "PROJECT RULES", "unlink lost the user's AGENTS.md prose")
	requireFixtureFileNotContains(t, f, "AGENTS.md", "bench:start", "unlink left the managed AGENTS.md marker behind")
}

// Story 5: user artifacts are never removed.
func testUnlinkSparesUserArtifacts(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	artifacts := map[string]string{
		"ROADMAP.md":          "# Roadmap\n",
		"IDEAS.md":            "- idea\n",
		"CONTEXT.md":          "# Context\n",
		".bench/learnings.md": "- learning\n",
		".bench/gate.sh":      "#!/usr/bin/env bash\n# drifted by hand\nexit 0\n",
		"specs/mine.md":       "# my spec\n",
	}
	for rel, body := range artifacts {
		f.WriteFile(rel, body)
	}

	unlinkOK(t, f)

	for rel := range artifacts {
		requireLinkFile(t, f, rel)
	}
	requireFixtureFileContains(t, f, ".bench/gate.sh", "drifted by hand", "unlink overwrote a drifted gate.sh")
}

// Story 6: the manifest is removed only when nothing was refused.
func testUnlinkManifestConditional(t *testing.T) {
	clean := contract.NewFixture(t)
	linkOK(t, clean)
	unlinkOK(t, clean)
	requireLinkNotExists(t, clean, ".bench/link-manifest.tsv", "clean unlink left the manifest behind")

	refused := contract.NewFixture(t)
	linkOK(t, refused)
	appendFile(t, filepath.Join(refused.Root, filepath.FromSlash(managedFileRel)), "\nlocal edit\n")
	unlinkOK(t, refused)
	requireLinkFile(t, refused, ".bench/link-manifest.tsv")
}

// Story 7: a traversal or root-escaping manifest row is refused, not deleted.
func testUnlinkRefusesTraversal(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	outside := filepath.Join(filepath.Dir(f.Root), "outside-sentinel")
	contract.WriteFileAbs(t, outside, "do not delete me\n")
	appendFile(t, filepath.Join(f.Root, ".bench", "link-manifest.tsv"), "../outside-sentinel\tdeadbeef\n")

	probe := f.Bench("unlink")
	probe.RequireExit(0)

	if got := contract.ReadFileAbs(t, outside); !strings.Contains(got, "do not delete me") {
		t.Fatalf("unlink followed a traversal row and touched a file outside the repo:\n%s", got)
	}
	probe.RequireContains(probe.Stdout, "refused")
	requireLinkFile(t, f, ".bench/link-manifest.tsv")
}

// Story 8: --dry-run prints the plan and changes nothing.
func testUnlinkDryRunNoWrites(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)
	f.CommitAll("linked")

	probe := f.Bench("unlink", "--dry-run")
	probe.RequireExit(0)
	probe.RequireContains(probe.Stdout, "would remove")

	status := f.Git("status", "--porcelain").Stdout
	if strings.TrimSpace(status) != "" {
		t.Fatalf("dry-run mutated the working tree:\n%s", status)
	}
	requireLinkFile(t, f, managedFileRel)
	requireLinkFile(t, f, ".bench/link-manifest.tsv")
}

// Story 9: the report lists removed and kept-modified paths on both runs.
func testUnlinkReportsBothRuns(t *testing.T) {
	real := contract.NewFixture(t)
	linkOK(t, real)
	appendFile(t, filepath.Join(real.Root, filepath.FromSlash(managedFileRel)), "\nlocal edit\n")
	realProbe := real.Bench("unlink")
	realProbe.RequireExit(0)
	realProbe.RequireContains(realProbe.Stdout, "removed: .bench/BENCH.md")
	realProbe.RequireContains(realProbe.Stdout, managedFileRel)
	realProbe.RequireContains(realProbe.Stdout, "modified")

	dry := contract.NewFixture(t)
	linkOK(t, dry)
	appendFile(t, filepath.Join(dry.Root, filepath.FromSlash(managedFileRel)), "\nlocal edit\n")
	dryProbe := dry.Bench("unlink", "--dry-run")
	dryProbe.RequireExit(0)
	dryProbe.RequireContains(dryProbe.Stdout, "would remove: .bench/BENCH.md")
	dryProbe.RequireContains(dryProbe.Stdout, managedFileRel)
	dryProbe.RequireContains(dryProbe.Stdout, "modified")
	requireFixtureFileContains(t, dry, managedFileRel, "local edit", "dry-run report path removed the kept file")
}

// Story 10: absent or unreadable manifest exits 1; a rowless manifest exits 0; a second
// unlink after a full run also exits 1.
func testUnlinkAbsentManifest(t *testing.T) {
	absent := contract.NewFixture(t)
	noManifest := absent.Bench("unlink")
	if noManifest.ExitCode != 1 {
		t.Fatalf("unlink on a repo with no manifest exit = %d, want 1\nstdout:\n%s\nstderr:\n%s", noManifest.ExitCode, noManifest.Stdout, noManifest.Stderr)
	}
	noManifest.RequireContains(noManifest.Stderr, "no link manifest")

	rowless := contract.NewFixture(t)
	rowless.WriteFile(".bench/link-manifest.tsv", "#kit\t0.0.0\n")
	rowlessProbe := rowless.Bench("unlink")
	rowlessProbe.RequireExit(0)
	rowlessProbe.RequireContains(rowlessProbe.Stdout, "nothing to remove")

	rerun := contract.NewFixture(t)
	linkOK(t, rerun)
	unlinkOK(t, rerun)
	second := rerun.Bench("unlink")
	if second.ExitCode != 1 {
		t.Fatalf("second unlink after a full run exit = %d, want 1\nstderr:\n%s", second.ExitCode, second.Stderr)
	}
}

// Story 11: unlink resolves to one implementation from the linked-repo by-path CLI, not the
// adoption-route refusal.
func testUnlinkByPathRoute(t *testing.T) {
	f := contract.NewFixture(t)
	linkOK(t, f)

	probe := f.Run(filepath.Join(".bench", "bin", "bench.sh"), "unlink", "--dry-run")
	probe.RequireExit(0)
	probe.RequireNotContains(probe.Stdout+probe.Stderr, "real Bench kit")
	probe.RequireContains(probe.Stdout, "would remove")
}

// Story 11: a CLAUDE.md link owns — created fresh, or retrofitted from the retired
// legacy bench shape — is removed by unlink.
func testUnlinkRemovesLinkCreatedClaudeMD(t *testing.T) {
	for name, seed := range map[string]func(f contract.Fixture){
		"no prior file": func(contract.Fixture) {},
		"legacy bench shape": func(f contract.Fixture) {
			f.WriteFile("CLAUDE.md", "# Bench\n\nCanonical agreement in AGENTS.md.\n\n@AGENTS.md\n")
		},
	} {
		t.Run(name, func(t *testing.T) {
			f := contract.NewFixture(t)
			seed(f)
			linkOK(t, f)
			requireLinkFile(t, f, "CLAUDE.md")
			requireFixtureFileContains(t, f, ".bench/link-manifest.tsv", "CLAUDE.md\t", "link did not record the CLAUDE.md it owns in the manifest")

			unlinkOK(t, f)

			requireLinkNotExists(t, f, "CLAUDE.md", "unlink left a link-owned CLAUDE.md on disk")
		})
	}
}

// Story 11: a CLAUDE.md that predates link — including one present but empty — is never
// recorded in the link manifest and survives unlink with identical bytes.
func testUnlinkSparesExistingClaudeMD(t *testing.T) {
	for name, body := range map[string]string{
		"user content": "# Custom\n\nproject-owned claude config\n",
		"empty file":   "",
	} {
		t.Run(name, func(t *testing.T) {
			f := contract.NewFixture(t)
			f.WriteFile("CLAUDE.md", body)

			linkOK(t, f)
			if got := f.ReadFile("CLAUDE.md"); got != body {
				t.Fatalf("link altered a pre-existing user CLAUDE.md:\nbefore:\n%q\nafter:\n%q", body, got)
			}
			requireFixtureFileNotContains(t, f, ".bench/link-manifest.tsv", "CLAUDE.md\t", "link recorded a pre-existing user CLAUDE.md in the manifest")

			unlinkOK(t, f)
			if got := f.ReadFile("CLAUDE.md"); got != body {
				t.Fatalf("unlink altered a pre-existing user CLAUDE.md:\nbefore:\n%q\nafter:\n%q", body, got)
			}
		})
	}
}

func unlinkOK(t *testing.T, f contract.Fixture) contract.Probe {
	t.Helper()
	probe := f.Bench("unlink")
	probe.RequireExit(0)
	return probe
}

func relTo(t *testing.T, f contract.Fixture, abs string) string {
	t.Helper()
	rel, err := filepath.Rel(f.Root, abs)
	if err != nil {
		t.Fatalf("relativize %s against %s: %v", abs, f.Root, err)
	}
	return rel
}
