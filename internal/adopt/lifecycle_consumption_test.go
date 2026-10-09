package adopt

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/compatibility/compatibilitytest"
)

// The ordinary managed assets that the real Link fixtures change: a regular kit file, and
// an executable kit file for the mode fixtures.
const (
	lifecycleDoc  = ".bench/BENCH.md"
	lifecycleExec = ".bench/bin/bench.sh"
)

// entryState is what a preserved destination keeps: its bytes or link target and mode,
// its file identity, and its manifest row.
type entryState struct {
	snapshot map[string]compatibilitytest.HomeEntry
	info     os.FileInfo
	row      string
}

func observeEntry(t *testing.T, root, rel string) entryState {
	t.Helper()
	path := filepath.Join(root, rel)
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	return entryState{compatibilitytest.SnapshotHome(t, path), info, manifestHash(root, rel)}
}

// requireKept fails unless rel keeps the state that observeEntry recorded before the run.
func requireKept(t *testing.T, root, rel string, before entryState) {
	t.Helper()
	after := observeEntry(t, root, rel)
	if !reflect.DeepEqual(after.snapshot, before.snapshot) || !os.SameFile(after.info, before.info) || after.row != before.row {
		t.Fatalf("%s = %+v with row %q, want the preserved %+v with row %q", rel, after.snapshot, after.row, before.snapshot, before.row)
	}
}

// requireInstalled fails unless rel holds the bytes and mode of its kit source, and the
// manifest records the fingerprint of that source.
func requireInstalled(t *testing.T, root, rel string) {
	t.Helper()
	entry := linkPlanEntry(t, rel, func(e planEntry) bool { return e.rel == rel })
	want, err := fingerprintPath(entry.src)
	if err != nil {
		t.Fatal(err)
	}
	if got := compatibilitytest.SnapshotHome(t, filepath.Join(root, rel)); !reflect.DeepEqual(got, compatibilitytest.SnapshotHome(t, entry.src)) {
		t.Fatalf("%s = %+v, want the kit source bytes and mode", rel, got)
	}
	if row := manifestHash(root, rel); row != want {
		t.Fatalf("%s manifest row = %q, want the kit fingerprint %q", rel, row, want)
	}
}

// TestLifecycleLinkPartitions drives real Link through each ordinary desired partition. An
// installed path takes the kit bytes, mode, and fingerprint. Every other path keeps its
// bytes, identity, and manifest row, and a reported path names its conflict reason. A
// copy relinked in symlink mode is converged because it resolves to the bytes that the
// staged link names, so it stays a regular file.
func TestLifecycleLinkPartitions(t *testing.T) {
	ownOld := func(t *testing.T, root, path string) {
		writeFixtureFile(t, path, "old kit\n", 0o644)
		old, err := fingerprintPath(path)
		if err != nil {
			t.Fatal(err)
		}
		editPreviewManifest(t, root, func(rows map[string]string) { rows[lifecycleDoc] = old })
	}
	for _, tc := range []struct {
		name, rel, reason string
		args              []string
		linked, installed bool
		arrange           func(t *testing.T, root, path string)
	}{
		{name: "absent", rel: lifecycleDoc, installed: true, arrange: func(*testing.T, string, string) {}},
		{name: "clean old", rel: lifecycleDoc, linked: true, installed: true, arrange: ownOld},
		{name: "converged", rel: lifecycleDoc, linked: true, arrange: func(*testing.T, string, string) {}},
		{name: "converged in symlink mode", rel: lifecycleDoc, args: []string{"symlink"}, linked: true, arrange: func(*testing.T, string, string) {}},
		{name: "mode", rel: lifecycleExec, linked: true, installed: true, arrange: func(t *testing.T, _, path string) {
			if err := os.Chmod(path, 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "modified", rel: lifecycleDoc, linked: true, reason: "modified-managed", arrange: func(t *testing.T, _, path string) {
			writeFixtureFile(t, path, "project edit\n", 0o644)
		}},
		{name: "equal-byte collision", rel: lifecycleDoc, reason: "project-owned", arrange: func(t *testing.T, _, path string) {
			src := linkPlanEntry(t, lifecycleDoc, func(e planEntry) bool { return e.rel == lifecycleDoc }).src
			writeFixtureFile(t, path, readFile(t, src), 0o644)
		}},
		{name: "dangling leaf collision", rel: lifecycleDoc, reason: "project-owned", arrange: func(t *testing.T, _, path string) {
			if err := os.Symlink("missing-target", path); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := consumerRepo(t)
			path := filepath.Join(root, tc.rel)
			if tc.linked {
				linkConsumer(t, root, 0)
			} else if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			tc.arrange(t, root, path)
			var before entryState
			if !tc.installed {
				before = observeEntry(t, root, tc.rel)
			}
			code := 0
			if tc.reason != "" {
				code = 3
			}
			var stdout, stderr bytes.Buffer
			if got := Link(tc.args, &stdout, &stderr, "1.0.0"); got != code || tc.reason != "" && !strings.Contains(stdout.String(), tc.rel+","+tc.reason+"\n") {
				t.Fatalf("Link %v = %d, want %d with the %q row for %s\nstdout:\n%s\nstderr:\n%s", tc.args, got, code, tc.reason, tc.rel, &stdout, &stderr)
			}
			if tc.installed {
				requireInstalled(t, root, tc.rel)
			} else {
				requireKept(t, root, tc.rel, before)
			}
		})
	}
}

// TestLifecycleOwnerConsumption drives each applying entry that consumes the desired
// decision. Each entry starts from a manifest that pins the linked kit and owns no row,
// and it must install an absent ordinary asset that the owner adds.
func TestLifecycleOwnerConsumption(t *testing.T) {
	for _, entry := range []struct {
		name string
		run  func(t *testing.T, root string)
	}{
		{name: "link", run: func(t *testing.T, root string) { linkConsumer(t, root, 0) }},
		{name: "confirmed setup", run: func(t *testing.T, _ string) { runSetupYes(t) }},
		{name: "applying upgrade", run: func(t *testing.T, _ string) {
			if code, stdout, stderr := runPreviewUpgrade(t, nil, "1.0.1"); code != 0 {
				t.Fatalf("Upgrade = (%d, %q, %q), want exit 0", code, stdout, stderr)
			}
		}},
		{name: "compatibility repair", run: func(t *testing.T, _ string) {
			if code, stdout, stderr := runCompatibilityRepair(t); code != 1 || stderr != "" {
				t.Fatalf("compatibility repair = (%d, %q, %q), want exit 1 with no error", code, stdout, stderr)
			}
		}},
	} {
		t.Run(entry.name, func(t *testing.T) {
			root := consumerRepo(t)
			if err := writeManifest(previewManifestPath(root), "1.0.0", nil); err != nil {
				t.Fatal(err)
			}
			entry.run(t, root)
			requireInstalled(t, root, lifecycleDoc)
		})
	}
}
