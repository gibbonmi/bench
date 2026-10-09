package adopt

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/toon"
)

// previewCounts is one upgrade count-table row.
type previewCounts struct{ added, changed, removed int }

// previewTable renders the count table that a real Upgrade prints before any guard.
func previewTable(t *testing.T, from, to string, want previewCounts) string {
	t.Helper()
	block, err := toon.TableTyped(upgradeTable, upgradeColumns, [][]any{{from, to, want.added, want.changed, want.removed}})
	if err != nil {
		t.Fatal(err)
	}
	return block
}

func runPreviewUpgrade(t *testing.T, args []string, version string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := Upgrade(args, &stdout, &stderr, version)
	return code, stdout.String(), stderr.String()
}

func previewManifestPath(root string) string {
	return filepath.Join(root, ".bench", "link-manifest.tsv")
}

// editPreviewManifest rewrites the linked manifest rows through edit and keeps its version.
func editPreviewManifest(t *testing.T, root string, edit func(map[string]string)) {
	t.Helper()
	m, err := ReadManifest(previewManifestPath(root))
	if err != nil {
		t.Fatal(err)
	}
	edit(m.hashes)
	if err := writeManifest(previewManifestPath(root), m.KitVersion, m.Rows()); err != nil {
		t.Fatal(err)
	}
}

// previewEntry returns the first real link plan entry of kind.
func previewEntry(t *testing.T, kind string) string {
	t.Helper()
	plan, err := buildLinkPlan(gate.KitDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range plan {
		if entry.kind == kind {
			return entry.rel
		}
	}
	t.Fatalf("link plan has no %s entry", kind)
	return ""
}

func TestLifecycleUpgradePreviewCompatibility(t *testing.T) {
	setHash := func(kind, hash string) func(*testing.T, string) {
		return func(t *testing.T, root string) {
			rel := previewEntry(t, kind)
			editPreviewManifest(t, root, func(rows map[string]string) { rows[rel] = hash })
		}
	}
	dropRow := func(kind string) func(*testing.T, string) {
		return func(t *testing.T, root string) {
			rel := previewEntry(t, kind)
			editPreviewManifest(t, root, func(rows map[string]string) { delete(rows, rel) })
		}
	}
	for _, tc := range []struct {
		name    string
		version string
		arrange func(*testing.T, string)
		want    previewCounts
	}{
		{name: "converged", version: "1.0.1", arrange: func(*testing.T, string) {}},
		{name: "owned file drift", version: "1.0.1", arrange: setHash("file", "stale"), want: previewCounts{changed: 1}},
		{name: "owned inline drift", version: "1.0.1", arrange: setHash("inline", "stale")},
		{name: "owned adapter drift", version: "1.0.1", arrange: setHash("adapter", "stale")},
		{name: "unowned file", version: "1.0.1", arrange: dropRow("file"), want: previewCounts{added: 1}},
		{name: "unowned inline", version: "1.0.1", arrange: dropRow("inline"), want: previewCounts{added: 1}},
		{name: "unowned adapter", version: "1.0.1", arrange: dropRow("adapter"), want: previewCounts{added: 1}},
		{name: "dropped row", version: "1.0.1", arrange: func(t *testing.T, root string) {
			editPreviewManifest(t, root, func(rows map[string]string) { rows[".agents/commands/retired.md"] = "old" })
		}, want: previewCounts{removed: 1}},
		{name: "preserved CLAUDE.md", version: "1.0.1", arrange: func(t *testing.T, root string) {
			editPreviewManifest(t, root, func(rows map[string]string) {
				if rows["CLAUDE.md"] == "" {
					t.Fatal("Link recorded no CLAUDE.md row")
				}
				rows["CLAUDE.md"] = "project edit"
			})
		}},
		{name: "hook refresh", version: "1.0.1", arrange: removePreviewHook, want: previewCounts{changed: 1}},
		{name: "hook refresh at equal version", version: "1.0.0", arrange: removePreviewHook},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := consumerRepo(t)
			linkConsumer(t, root, 0)
			tc.arrange(t, root)
			before := readFile(t, previewManifestPath(root))
			code, stdout, stderr := runPreviewUpgrade(t, []string{"--check"}, tc.version)
			if want := previewTable(t, "1.0.0", tc.version, tc.want); code != 0 || stdout != want || stderr != "" {
				t.Fatalf("Upgrade --check = (%d, %q, %q), want (0, %q, \"\")", code, stdout, stderr, want)
			}
			if got := readFile(t, previewManifestPath(root)); got != before {
				t.Fatalf("Upgrade --check rewrote the manifest:\n%s", got)
			}
		})
	}
	t.Run("plan kinds", testPreviewPlanKinds)
}

func removePreviewHook(t *testing.T, root string) {
	t.Helper()
	if err := os.Remove(filepath.Join(root, ".git", "hooks", "pre-push")); err != nil {
		t.Fatal(err)
	}
}

// testPreviewPlanKinds drives upgradePlanCounts for every plan kind and an unreadable
// file source. Only Setup produces inline-exec and seed entries, so a real Upgrade
// cannot reach them.
func testPreviewPlanKinds(t *testing.T) {
	src := filepath.Join(t.TempDir(), "asset")
	writeFixtureFile(t, src, "kit\n", 0o644)
	current, err := fingerprintPath(src)
	if err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(t.TempDir(), "missing")
	for _, tc := range []struct {
		name  string
		entry planEntry
		rows  map[string]string
		want  previewCounts
	}{
		{name: "file current", entry: planEntry{src: src, rel: "a", kind: "file"}, rows: map[string]string{"a": current}},
		{name: "file stale", entry: planEntry{src: src, rel: "a", kind: "file"}, rows: map[string]string{"a": "stale"}, want: previewCounts{changed: 1}},
		{name: "file unowned", entry: planEntry{src: src, rel: "a", kind: "file"}, want: previewCounts{added: 1}},
		{name: "file unreadable", entry: planEntry{src: missing, rel: "a", kind: "file"}, rows: map[string]string{"a": "stale"}, want: previewCounts{changed: 1}},
		{name: "file unreadable unowned", entry: planEntry{src: missing, rel: "a", kind: "file"}, want: previewCounts{added: 1}},
		{name: "inline stale", entry: planEntry{rel: "a", kind: "inline", content: "x\n"}, rows: map[string]string{"a": "stale"}},
		{name: "inline unowned", entry: planEntry{rel: "a", kind: "inline", content: "x\n"}, want: previewCounts{added: 1}},
		{name: "inline-exec stale", entry: planEntry{rel: "a", kind: "inline-exec", content: "x\n"}, rows: map[string]string{"a": "stale"}},
		{name: "inline-exec unowned", entry: planEntry{rel: "a", kind: "inline-exec", content: "x\n"}, want: previewCounts{added: 1}},
		{name: "adapter stale", entry: planEntry{src: src, rel: "a", kind: "adapter"}, rows: map[string]string{"a": "stale"}},
		{name: "adapter unowned", entry: planEntry{src: src, rel: "a", kind: "adapter"}, want: previewCounts{added: 1}},
		{name: "seed stale", entry: planEntry{rel: "a", kind: "seed", content: "x\n"}, rows: map[string]string{"a": "stale"}},
		{name: "seed unowned", entry: planEntry{rel: "a", kind: "seed", content: "x\n"}, want: previewCounts{added: 1}},
		{name: "dropped row", entry: planEntry{rel: "a", kind: "inline"}, rows: map[string]string{"a": "h", "gone": "h"}, want: previewCounts{removed: 1}},
		{name: "CLAUDE.md row", entry: planEntry{rel: "a", kind: "inline"}, rows: map[string]string{"a": "h", "CLAUDE.md": "h"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			hashes := map[string]string{}
			for rel, hash := range tc.rows {
				hashes[rel] = hash
			}
			added, changed, removed := upgradePlanCounts(Manifest{KitVersion: "1.0.0", hashes: hashes}, []planEntry{tc.entry})
			if got := (previewCounts{added, changed, removed}); got != tc.want {
				t.Fatalf("upgradePlanCounts = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestLifecycleUpgradeVersionGuards(t *testing.T) {
	const gitignore = ".bench/dist/.gitignore"
	for _, tc := range []struct {
		name, version, wantKit, wantErr string
		args                            []string
		code                            int
		applied                         bool
	}{
		{name: "equal", version: "1.0.0", wantKit: "1.0.0"},
		{name: "newer", version: "1.0.1", wantKit: "1.0.1", applied: true},
		{name: "older", version: "0.9.0", wantKit: "1.0.0", code: 1, wantErr: "installed kit is older than the linked kit"},
		{name: "older check", version: "0.9.0", args: []string{"--check"}, wantKit: "1.0.0", code: 1, wantErr: "installed kit is older than the linked kit"},
		{name: "forced older", version: "0.9.0", args: []string{"--force"}, wantKit: "0.9.0", applied: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := consumerRepo(t)
			linkConsumer(t, root, 0)
			if err := os.Remove(filepath.Join(root, gitignore)); err != nil {
				t.Fatal(err)
			}
			code, stdout, stderr := runPreviewUpgrade(t, tc.args, tc.version)
			if code != tc.code || !strings.HasPrefix(stdout, previewTable(t, "1.0.0", tc.version, previewCounts{})) || !strings.Contains(stderr, tc.wantErr) {
				t.Fatalf("Upgrade %v at %s = (%d, %q, %q), want exit %d and %q", tc.args, tc.version, code, stdout, stderr, tc.code, tc.wantErr)
			}
			m, err := ReadManifest(previewManifestPath(root))
			if err != nil || m.KitVersion != tc.wantKit {
				t.Fatalf("manifest kit = %q, %v; want %q", m.KitVersion, err, tc.wantKit)
			}
			if _, err := os.Lstat(filepath.Join(root, gitignore)); (err == nil) != tc.applied {
				t.Fatalf("%s present = %t, want %t", gitignore, err == nil, tc.applied)
			}
		})
	}
	for _, tc := range []struct{ name, manifest, wantErr string }{
		{name: "absent manifest", wantErr: "repository is not Bench-linked"},
		{name: "empty manifest", manifest: "", wantErr: "link manifest is empty"},
		{name: "missing header", manifest: "a\tdeadbeef\n", wantErr: "link manifest carries no #kit version header"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := consumerRepo(t)
			if tc.name != "absent manifest" {
				if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
					t.Fatal(err)
				}
				writeFixtureFile(t, previewManifestPath(root), tc.manifest, 0o644)
			}
			code, stdout, stderr := runPreviewUpgrade(t, nil, "1.0.1")
			if code != 1 || stdout != "" || !strings.Contains(stderr, tc.wantErr) {
				t.Fatalf("Upgrade = (%d, %q, %q), want exit 1 and %q", code, stdout, stderr, tc.wantErr)
			}
		})
	}
}
