//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/compatibility/compatibilitytest"
)

// The managed assets that the lifecycle journeys lose and expect back. lifecycleDoc
// precedes lifecycleReference in the link plan.
const (
	lifecycleDoc       = ".bench/BENCH.md"
	lifecycleReference = ".bench/BENCH-reference.md"
	lifecycleManifest  = ".bench/link-manifest.tsv"
	lifecycleHook      = ".git/hooks/pre-push"
	lifecycleProse     = "# Team rules\n\nKeep this rule.\n"
)

// runLifecycleEntry runs one selected-binary entry in repo. The overrides follow the
// selected binary and kit, so a later BENCH_KIT wins. TMPDIR keeps the recovery directory
// of a failed or interrupted publication inside the test.
func runLifecycleEntry(t *testing.T, repo string, overrides []string, args ...string) processResult {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	env := append([]string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit, "TMPDIR=" + t.TempDir()}, overrides...)
	return owner.runAt(repo, env, owner.selected.path, args...)
}

// linkedLifecycleRepo is a fresh repository, apart from the owner's shared ones, that the
// selected binary has linked over project AGENTS.md prose.
func linkedLifecycleRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if result := owner.runAt(repo, nil, "git", "init", "-q", "-b", "main"); result.code != 0 {
		t.Fatalf("git init = %#v", result)
	}
	commitmenttest.Write(t, repo, "AGENTS.md", lifecycleProse)
	if result := owner.runSelected(repo, "link", "copy"); result.code != 0 {
		t.Fatalf("link fixture = %#v", result)
	}
	return repo
}

func removeLifecycleFiles(t *testing.T, repo string, rels ...string) {
	t.Helper()
	for _, rel := range rels {
		if err := os.Remove(filepath.Join(repo, filepath.FromSlash(rel))); err != nil {
			t.Fatal(err)
		}
	}
}

// requireLifecycleState fails unless each rel holds its kit bytes when present is true, and
// is absent otherwise.
func requireLifecycleState(t *testing.T, repo string, present bool, rels ...string) {
	t.Helper()
	for _, rel := range rels {
		data, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(rel)))
		if !present {
			if !os.IsNotExist(err) {
				t.Fatalf("%s: %v, want absent", rel, err)
			}
			continue
		}
		if want := readSpillFile(t, filepath.Join(owner.kit, filepath.FromSlash(rel))); err != nil || string(data) != want {
			t.Fatalf("%s = %q, %v; want the kit bytes", rel, data, err)
		}
	}
}

func requireLifecycleProse(t *testing.T, repo string) {
	t.Helper()
	if got := readSpillFile(t, filepath.Join(repo, "AGENTS.md")); got != lifecycleProse {
		t.Fatalf("AGENTS.md = %q, want the project prose", got)
	}
}

// loseLifecycleAssets removes two owned assets and pins the manifest to an older kit, so a
// relink publishes both assets and then a changed manifest. It returns the clean manifest.
func loseLifecycleAssets(t *testing.T, repo string) string {
	t.Helper()
	manifest := readSpillFile(t, filepath.Join(repo, lifecycleManifest))
	_, rows, _ := strings.Cut(manifest, "\n")
	commitmenttest.Write(t, repo, lifecycleManifest, "#kit\t0.0.0\n"+rows)
	removeLifecycleFiles(t, repo, lifecycleDoc, lifecycleReference)
	return manifest
}

// trimLifecycleUnlink keeps only the manifest rows named in rels, so unlink removes those
// rows, publishes the stripped AGENTS.md, and then removes the managed hook.
func trimLifecycleUnlink(t *testing.T, repo string, rels ...string) {
	t.Helper()
	path := filepath.Join(repo, lifecycleManifest)
	kept := []string{}
	for i, line := range strings.SplitAfter(readSpillFile(t, path), "\n") {
		if rel, _, _ := strings.Cut(line, "\t"); i == 0 || slices.Contains(rels, rel) {
			kept = append(kept, line)
		}
	}
	commitmenttest.Write(t, repo, lifecycleManifest, strings.Join(kept, ""))
}

// TestLifecycleEntryJourney drives each applying lifecycle entry through the selected binary
// and observes the effect that its lifecycle decision selects. Link installs absent managed
// assets. Confirmed Setup seeds an absent gate input without a manifest row. An applying
// Upgrade and a compatibility repair each restore an owned asset that the repository lost.
// Unlink removes the clean managed assets with the manifest and restores the project prose.
func TestLifecycleEntryJourney(t *testing.T) {
	repo := linkedLifecycleRepo(t)
	requireLifecycleState(t, repo, true, lifecycleDoc, lifecycleReference)

	seed := ".bench/gate-inputs.json"
	if result := owner.runSelected(repo, "setup", "--yes"); result.code != 3 {
		t.Fatalf("setup --yes = %#v, want exit 3 with the gate stub", result)
	}
	if _, err := os.Stat(filepath.Join(repo, seed)); err != nil || strings.Contains(readSpillFile(t, filepath.Join(repo, lifecycleManifest)), seed+"\t") {
		t.Fatalf("seed %s: %v, want it present without a manifest row", seed, err)
	}

	loseLifecycleAssets(t, repo)
	if result := owner.runSelected(repo, "upgrade"); result.code != 0 {
		t.Fatalf("upgrade = %#v, want exit 0", result)
	}
	requireLifecycleState(t, repo, true, lifecycleDoc, lifecycleReference)

	fixture := newCompatibilitySession(t)
	removeLifecycleFiles(t, fixture.repo, lifecycleDoc)
	if result := runLifecycleEntry(t, fixture.repo, fixture.overrides, "doctor", "--compat", "codex-cli", "--fix"); result.code != 1 || result.stderr != "" {
		t.Fatalf("compatibility repair = %#v, want exit 1 with no error", result)
	}
	requireLifecycleState(t, fixture.repo, true, lifecycleDoc)

	if result := owner.runSelected(repo, "unlink"); result.code != 0 {
		t.Fatalf("unlink = %#v, want exit 0", result)
	}
	requireLifecycleState(t, repo, false, lifecycleDoc, lifecycleReference, lifecycleManifest)
	requireLifecycleProse(t, repo)
}

// lifecycleKit writes a kit tree that holds one file at rel, or a directory when the file
// is empty, and returns the BENCH_KIT override that selects it.
func lifecycleKit(t *testing.T, rel, content string) []string {
	t.Helper()
	kit := t.TempDir()
	if content == "" {
		if err := os.MkdirAll(filepath.Join(kit, filepath.FromSlash(rel)), 0o755); err != nil {
			t.Fatal(err)
		}
	} else {
		commitmenttest.Write(t, kit, rel, content)
	}
	return []string{"BENCH_KIT=" + kit}
}

// TestLifecyclePublicationCompatibility pins the failure and recovery results of lifecycle
// publication through the selected binary. Each failure exits 1 and names its cause. A
// link failure and an unlink promotion failure leave the repository unchanged. An unlink
// report failure follows its publication, so a clean asset is gone and the manifest stays.
// After an interruption at a publication boundary, a fresh process converges the
// repository to the clean link or the clean unlink.
func TestLifecyclePublicationCompatibility(t *testing.T) {
	control := ".agents/commands/control\x01.md"
	link, unlink := []string{"link", "copy"}, []string{"unlink"}
	for _, tc := range []struct {
		name, cause string
		args        []string
		changed     bool
		arrange     func(t *testing.T, repo string) []string
	}{
		{name: "missing source", cause: lifecycleDoc, args: link, arrange: func(t *testing.T, _ string) []string {
			return lifecycleKit(t, lifecycleDoc, "")
		}},
		{name: "stage failure", cause: "permission denied", args: link, arrange: func(t *testing.T, repo string) []string {
			// The hook stages beside its destination, so a read-only hooks directory fails
			// that stage after the lease is held.
			hooks := filepath.Dir(filepath.Join(repo, lifecycleHook))
			removeLifecycleFiles(t, repo, lifecycleHook)
			if err := os.Chmod(hooks, 0o555); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.Chmod(hooks, 0o755) })
			if probe, err := os.CreateTemp(hooks, "lifecycle-stage-probe-"); err == nil {
				_ = probe.Close()
				_ = os.Remove(probe.Name())
				capability.Capability(t, capability.Privilege, "read-only hooks directory remains writable by this user")
			}
			return nil
		}},
		{name: "report failure", cause: "control", args: link, arrange: func(t *testing.T, repo string) []string {
			commitmenttest.Write(t, repo, control, "project command\n")
			return lifecycleKit(t, control, "kit command\n")
		}},
		{name: "promotion failure", cause: "fault 1", args: link, arrange: func(t *testing.T, repo string) []string {
			loseLifecycleAssets(t, repo)
			return []string{"BENCH_LINK_FAULT=1"}
		}},
		{name: "final manifest failure", cause: "fault last", args: link, arrange: func(t *testing.T, repo string) []string {
			loseLifecycleAssets(t, repo)
			return []string{"BENCH_LINK_FAULT=last"}
		}},
		{name: "unlink promotion failure", cause: "fault last", args: unlink, arrange: func(t *testing.T, repo string) []string {
			// Without the hook, the stripped AGENTS.md is the last change, so the fault
			// meets a replacement.
			trimLifecycleUnlink(t, repo, lifecycleDoc)
			removeLifecycleFiles(t, repo, lifecycleHook)
			return []string{"BENCH_LINK_FAULT=last"}
		}},
		{name: "unlink report failure", cause: "control", args: unlink, changed: true, arrange: func(t *testing.T, repo string) []string {
			manifest := filepath.Join(repo, lifecycleManifest)
			commitmenttest.Write(t, repo, lifecycleManifest, readSpillFile(t, manifest)+"../control\x01\trefused\n")
			return nil
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := linkedLifecycleRepo(t)
			overrides := tc.arrange(t, repo)
			before := compatibilitytest.SnapshotHome(t, repo)
			manifest := readSpillFile(t, filepath.Join(repo, lifecycleManifest))
			if result := runLifecycleEntry(t, repo, overrides, tc.args...); result.code != 1 || !strings.Contains(result.stdout+result.stderr, tc.cause) {
				t.Fatalf("%s = %#v, want exit 1 naming %q", tc.name, result, tc.cause)
			}
			if !tc.changed {
				if after := compatibilitytest.SnapshotHome(t, repo); !reflect.DeepEqual(after, before) {
					t.Fatalf("failed publication changed the repository:\n%+v\nwant:\n%+v", after, before)
				}
				return
			}
			requireLifecycleState(t, repo, false, lifecycleDoc)
			if got := readSpillFile(t, filepath.Join(repo, lifecycleManifest)); got != manifest {
				t.Fatalf("manifest after the report failure = %q, want the kept %q", got, manifest)
			}
		})
	}
	// A relink of the lost assets publishes lifecycleDoc, lifecycleReference, CLAUDE.md,
	// the hook, and then the manifest. Unlink of a trimmed repository removes its kept rows,
	// publishes AGENTS.md, and then removes the hook. The interruption seam stops only a
	// rename, so an unlink interruption stops at AGENTS.md with the hook removal still
	// pending, and the final manifest removal of Unlink has no interruption case.
	for _, tc := range []struct {
		name, interrupt               string
		args, rows, published, absent []string
	}{
		{name: "link pre-publication", interrupt: "interrupt:1", args: link, absent: []string{lifecycleDoc, lifecycleReference}},
		{name: "link middle promotion", interrupt: "interrupt:2", args: link, published: []string{lifecycleDoc}, absent: []string{lifecycleReference}},
		{name: "link final manifest", interrupt: "interrupt:5", args: link, published: []string{lifecycleDoc, lifecycleReference}},
		{name: "unlink pre-publication", interrupt: "interrupt:1", args: unlink},
		{name: "unlink middle promotion", interrupt: "interrupt:2", args: unlink, rows: []string{lifecycleDoc}, absent: []string{lifecycleDoc}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := linkedLifecycleRepo(t)
			clean := ""
			if tc.args[0] == "unlink" {
				trimLifecycleUnlink(t, repo, tc.rows...)
			} else {
				clean = loseLifecycleAssets(t, repo)
			}
			manifest := readSpillFile(t, filepath.Join(repo, lifecycleManifest))
			if result := runLifecycleEntry(t, repo, []string{"BENCH_LINK_FAULT=" + tc.interrupt}, tc.args...); result.code != 97 {
				t.Fatalf("interrupted %s = %#v, want exit 97", tc.name, result)
			}
			requireLifecycleState(t, repo, true, tc.published...)
			requireLifecycleState(t, repo, false, tc.absent...)
			if got := readSpillFile(t, filepath.Join(repo, lifecycleManifest)); got != manifest {
				t.Fatalf("manifest after the interruption = %q, want the unpublished %q", got, manifest)
			}
			if clean == "" {
				if _, err := os.Lstat(filepath.Join(repo, lifecycleHook)); err != nil {
					t.Fatalf("hook after the interruption: %v, want its removal pending", err)
				}
			}
			if result := runLifecycleEntry(t, repo, nil, tc.args...); result.code != 0 {
				t.Fatalf("fresh-process %s = %#v, want exit 0", tc.name, result)
			}
			if clean == "" {
				requireLifecycleState(t, repo, false, append(tc.rows, lifecycleManifest, lifecycleHook)...)
				requireLifecycleProse(t, repo)
				return
			}
			requireLifecycleState(t, repo, true, lifecycleDoc, lifecycleReference)
			if got := readSpillFile(t, filepath.Join(repo, lifecycleManifest)); got != clean {
				t.Fatalf("manifest after recovery = %q, want the clean %q", got, clean)
			}
		})
	}
}
