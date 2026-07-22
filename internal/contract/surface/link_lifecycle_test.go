package surface

import (
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
}

func TestLinkMarkerFenceContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench link malformed marker contract failed", testLinkMalformedMarker)
	contract.RunParallel(t, "bench link fenced-marker contract failed", testLinkFencedMarker)
	contract.RunParallel(t, "bench link unclosed-fence contract failed", testLinkUnclosedFence)
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
