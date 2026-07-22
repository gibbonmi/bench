package surface

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// setup_rerun_test.go covers FT76 story 8: a re-run on an already-linked repository
// converges and reports. An unchanged tree reports already converged with an
// unchanged content hash; a managed asset modified or removed between runs
// reconciles through FT84's relink semantics on the next run.

func TestSetupRerunContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench setup re-run on an unchanged tree reports converged with an identical hash", testSetupRerunUnchangedTree)
	contract.RunParallel(t, "bench setup re-run recreates a removed managed asset", testSetupRerunRecreatesRemovedAsset)
	contract.RunParallel(t, "bench setup re-run preserves a hand-modified managed asset as a conflict", testSetupRerunPreservesModifiedAsset)
	contract.RunParallel(t, "bench setup --plan on a re-run previews without writing", testSetupRerunPlanDoesNotWrite)
	contract.RunParallel(t, "bench setup re-run with the manifest deleted reclassifies managed assets as conflicts", testSetupRerunWithoutManifestReclassifies)
}

func testSetupRerunUnchangedTree(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	first := f.Bench("setup", "--yes")
	first.RequireExit(0)
	// A fresh converged tree wrote something - it must never claim "already converged".
	if strings.Contains(strings.ToLower(first.Stdout), "already converged") {
		t.Fatalf("first run over a fresh tree reported already converged:\n%s", first.Stdout)
	}
	before := hashManagedTree(t, f.Root)

	probe := f.Bench("setup", "--yes")

	probe.RequireExit(0)
	// Strengthened over a bare "converged" substring (satisfied by either wording):
	// a true no-op re-run must use the distinct "already converged" phrasing.
	probe.RequireContains(strings.ToLower(probe.Stdout), "already converged")
	after := hashManagedTree(t, f.Root)
	if before != after {
		t.Fatalf("re-run over an unchanged tree changed content hash: %s -> %s", before, after)
	}
}

// testSetupRerunWithoutManifestReclassifies pins the FT84 lifecycle's actual behavior
// when .bench/link-manifest.tsv is deleted but the assets it tracked are still intact:
// with no manifest row to prove an asset is bench's, every managed asset the plan
// loop inspects (AGENTS.md/CLAUDE.md are reconciled through their own separate path
// and are unaffected) reads as project-owned and is preserved as a conflict rather
// than silently reclaimed. This pins the observed lifecycle behavior; it does not
// change production semantics.
func testSetupRerunWithoutManifestReclassifies(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	f.Bench("setup", "--yes").RequireExit(0)
	before := f.ReadFile(".bench/gate.sh")
	contract.Remove(t, filepath.Join(f.Root, ".bench", "link-manifest.tsv"))

	probe := f.Bench("setup", "--yes")

	probe.RequireExit(3)
	probe.RequireContains(strings.ToLower(probe.Stdout+probe.Stderr), "conflict")
	probe.RequireContains(probe.Stdout, ".bench/gate.sh")
	probe.RequireContains(probe.Stdout, "project-owned")
	if got := f.ReadFile(".bench/gate.sh"); got != before {
		t.Fatalf(".bench/gate.sh was not preserved byte-identical once the manifest could no longer prove ownership:\n%s", got)
	}
	// AGENTS.md/CLAUDE.md reconcile through their own marker/import-line path, not the
	// manifest-gated plan loop, so they still converge even with no manifest to read.
	requireLiteralCount(t, f, "AGENTS.md", "<!-- bench:start -->", 1, "AGENTS.md diverged when re-run without a prior manifest")
}

func testSetupRerunRecreatesRemovedAsset(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	f.Bench("setup", "--yes").RequireExit(0)
	contract.Remove(t, filepath.Join(f.Root, ".bench", "BENCH.md"))

	probe := f.Bench("setup", "--yes")

	probe.RequireExit(0)
	requireLinkFile(t, f, ".bench/BENCH.md")
}

func testSetupRerunPreservesModifiedAsset(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	f.Bench("setup", "--yes").RequireExit(0)
	const handEdit = "# hand-edited by the project, not bench\n"
	f.WriteFile(".bench/BENCH.md", handEdit)

	probe := f.Bench("setup", "--yes")

	probe.RequireExit(3)
	probe.RequireContains(strings.ToLower(probe.Stdout+probe.Stderr), "conflict")
	// The verdict vocabulary, not just the presence of a conflict: a genuinely
	// tampered manifest-tracked asset reads "modified-managed", distinct from a
	// never-managed foreign file's "project-owned".
	probe.RequireContains(probe.Stdout, "modified-managed")
	if got := f.ReadFile(".bench/BENCH.md"); got != handEdit {
		t.Fatalf("re-run overwrote a hand-modified managed asset:\n%s", got)
	}
	// The rest of the plan still converges alongside the preserved conflict.
	requireLinkFile(t, f, ".bench/gate.sh")
}

func testSetupRerunPlanDoesNotWrite(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("go.mod", "module example.com/fixture\n\ngo 1.22\n")
	f.Bench("setup", "--yes").RequireExit(0)
	before := hashManagedTree(t, f.Root)

	probe := f.Bench("setup", "--plan")

	probe.RequireExit(0)
	after := hashManagedTree(t, f.Root)
	if before != after {
		t.Fatalf("--plan on a re-run wrote to the tree: %s -> %s", before, after)
	}
}

// hashManagedTree hashes the content of every regular file bench setup manages,
// independent of git or mtimes, so a re-run's idempotency can be pinned by content
// alone (transactionalLink always restages and repromotes every accepted entry, so
// mtimes/inodes churn even when content does not).
func hashManagedTree(t *testing.T, root string) string {
	t.Helper()
	var paths []string
	for _, dir := range []string{".bench", ".agents", ".claude", ".codex", "AGENTS.md", "CLAUDE.md", "projects"} {
		full := filepath.Join(root, dir)
		_ = filepath.WalkDir(full, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil || !info.Mode().IsRegular() {
				return nil
			}
			rel, _ := filepath.Rel(root, path)
			paths = append(paths, filepath.ToSlash(rel))
			return nil
		})
	}
	sort.Strings(paths)
	h := sha256.New()
	for _, rel := range paths {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		h.Write([]byte(rel + "\x00"))
		h.Write(data)
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}
