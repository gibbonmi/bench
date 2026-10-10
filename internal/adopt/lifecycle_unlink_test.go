package adopt

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

// moveManifestRow renames the manifest row of from to the key to and keeps its fingerprint.
func moveManifestRow(t *testing.T, root, from, to string) {
	t.Helper()
	editPreviewManifest(t, root, func(rows map[string]string) {
		hash := rows[from]
		delete(rows, from)
		rows[to] = hash
	})
}

// TestLifecycleUnlinkPartitions drives real Unlink through each partition that leaves no
// residual. A clean row loses its asset. A row whose asset is already absent stays out of
// the report, and an emptied directory that only that row names is swept. A clean row that
// names its asset through sub/../ resolves to that asset and removes it. Each run exits 0,
// removes the manifest, and removes the AGENTS.md that link created with no project prose.
func TestLifecycleUnlinkPartitions(t *testing.T) {
	for _, tc := range []struct {
		name, row         string
		absent, dedicated bool
	}{
		{name: "clean", row: lifecycleDoc},
		{name: "absent", row: lifecycleDoc, absent: true},
		{name: "absent in an emptied dedicated parent", row: "dedicated/asset", absent: true, dedicated: true},
		{name: "normalized", row: "sub/../" + lifecycleDoc},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := linkedCopy(t, goldenPlain)
			path := filepath.Join(root, tc.row)
			moveManifestRow(t, root, lifecycleDoc, tc.row)
			if tc.absent {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
			}
			out := runUnlink(t, nil, 0)
			requireAbsent(t, path)
			requireAbsent(t, previewManifestPath(root))
			requireAbsent(t, filepath.Join(root, "AGENTS.md"))
			if tc.dedicated {
				requireAbsent(t, filepath.Dir(path))
			}
			if mentioned := strings.Contains(out, tc.row+"\n"); mentioned == tc.absent {
				t.Fatalf("unlink output names %s: %t, want %t\n%s", tc.row, mentioned, !tc.absent, out)
			}
		})
	}
}

// TestLifecycleUnlinkRefusals drives real Unlink over each residual of a linked
// repository. A modified asset keeps its bytes. A row that escapes the repository, a row
// under a symlink parent, and a row that names a FIFO or a directory keep their targets. An
// AGENTS.md that is a FIFO or holds a malformed managed block keeps its state without a
// manifest row. Each run exits 3, reports the residual with its reason, and keeps the
// manifest bytes, while a clean row still leaves.
func TestLifecycleUnlinkRefusals(t *testing.T) {
	plantFile := func(t *testing.T, path string) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		writeFixtureFile(t, path, "kit asset\n", 0o644)
	}
	outside := "../" + filepath.Base(t.TempDir()) + "/asset with [glob]*"
	for _, tc := range []struct {
		name, rel, reason string
		bespoke           bool
		arrange           func(t *testing.T, root, path string)
	}{
		{name: "AGENTS.md special file", rel: "AGENTS.md", reason: reasonRefused, bespoke: true, arrange: func(t *testing.T, _, path string) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "AGENTS.md malformed block", rel: "AGENTS.md", reason: reasonRefused, bespoke: true, arrange: func(t *testing.T, _, path string) {
			writeFixtureFile(t, path, readFile(t, path)+benchStartMarker+"\n", 0o644)
		}},
		{name: "modified", rel: lifecycleDoc, reason: reasonModified, arrange: func(t *testing.T, _, path string) {
			writeFixtureFile(t, path, "project edit\n", 0o644)
		}},
		{name: "escape", rel: outside, reason: reasonRefused, arrange: func(t *testing.T, _, path string) { plantFile(t, path) }},
		{name: "symlink parent", rel: "linked/asset", reason: reasonRefused, arrange: func(t *testing.T, root, path string) {
			plantSymlinkParent(t, root)
			plantFile(t, path)
		}},
		{name: "special file", rel: "fifo asset", reason: reasonRefused, arrange: func(t *testing.T, _, path string) {
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "directory", rel: "directory asset", reason: reasonRefused, arrange: func(t *testing.T, _, path string) {
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := linkedCopy(t, goldenPlain)
			path := filepath.Join(root, tc.rel)
			tc.arrange(t, root, path)
			if tc.rel != lifecycleDoc && !tc.bespoke {
				// A refused row records the fingerprint of a target that has one, so only the
				// refusal keeps that target.
				recorded, err := fingerprintPath(path)
				if err != nil {
					recorded = "no fingerprint"
				}
				editPreviewManifest(t, root, func(rows map[string]string) { rows[tc.rel] = recorded })
			}
			before := observeEntry(t, root, tc.rel)
			manifest := readFile(t, previewManifestPath(root))
			out := runUnlink(t, nil, 3)
			residual, err := renderVerdicts("residuals", []lifecycleVerdict{{tc.rel, tc.reason}})
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out, residual) {
				t.Fatalf("unlink output = %q, want the residual block %q", out, residual)
			}
			requireKept(t, root, tc.rel, before)
			if got := readFile(t, previewManifestPath(root)); got != manifest {
				t.Fatalf("manifest after a partial unlink = %q, want the kept bytes %q", got, manifest)
			}
			requireAbsent(t, filepath.Join(root, lifecycleExec))
		})
	}
}

// TestLifecycleDryRunSnapshot drives real Unlink --dry-run over a linked repository that
// holds clean rows, a modified asset, an absent asset, a refused row, project AGENTS.md
// prose around the managed block, and the managed hook. The dry run names a clean row in
// its plan and leaves every byte, mode, symlink target, and directory entry of the
// repository unchanged, the manifest included.
func TestLifecycleDryRunSnapshot(t *testing.T) {
	root := linkedCopy(t, goldenProse)
	writeFixtureFile(t, filepath.Join(root, lifecycleDoc), "project edit\n", 0o644)
	if err := os.Remove(filepath.Join(root, lifecycleExec)); err != nil {
		t.Fatal(err)
	}
	plantSymlinkParent(t, root)
	editPreviewManifest(t, root, func(rows map[string]string) { rows["linked/asset"] = "refused" })
	var out string
	requireUnchanged(t, root, func() { out = runUnlink(t, []string{"--dry-run"}, 3) })
	if clean := ".bench/BENCH-reference.md"; !strings.Contains(out, clean+"\n") {
		t.Fatalf("dry-run output = %q, want the clean row %s in its plan", out, clean)
	}
}
