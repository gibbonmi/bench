package adopt

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/compatibility"
	"github.com/gibbonmi/bench/internal/compatibility/compatibilitytest"
)

// lifecycleSeed is the Setup seed that the seed fixtures arrange.
const lifecycleSeed = ".bench/gate-inputs.json"

// runCompatibilityRepair runs a real compatibility repair in the working repository. The
// repair keeps its recovery records in a private Bench home.
func runCompatibilityRepair(t *testing.T) (int, string, string) {
	t.Helper()
	t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "bench-state"))
	var stdout, stderr bytes.Buffer
	code := compatibilityDoctor([]string{string(compatibility.CodexCLI), "--fix"}, &stdout, &stderr, func(compatibility.Interface) compatibility.Input {
		return compatibilityFixture(compatibility.CodexCLI)
	})
	return code, stdout.String(), stderr.String()
}

// TestLifecycleSeedOwnership drives confirmed Setup through each seed partition. An absent
// seed takes the scaffold bytes without a manifest row. A present seed keeps its bytes,
// identity, and manifest row, including an ownership row that an earlier release wrote.
func TestLifecycleSeedOwnership(t *testing.T) {
	for _, tc := range []struct {
		name      string
		installed bool
		arrange   func(t *testing.T, root, path string)
	}{
		{name: "absent", installed: true, arrange: func(*testing.T, string, string) {}},
		{name: "empty", arrange: func(t *testing.T, _, path string) { writeFixtureFile(t, path, "", 0o644) }},
		{name: "edited", arrange: func(t *testing.T, _, path string) { writeFixtureFile(t, path, "edited\n", 0o644) }},
		{name: "dangling leaf", arrange: func(t *testing.T, _, path string) {
			if err := os.Symlink("missing-target", path); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "old row", arrange: func(t *testing.T, root, path string) {
			runSetupYes(t)
			writeFixtureFile(t, path, "edited\n", 0o644)
			editPreviewManifest(t, root, func(rows map[string]string) { rows[lifecycleSeed] = "old-row" })
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := setupPromptTestRepo(t)
			path := filepath.Join(root, lifecycleSeed)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			tc.arrange(t, root, path)
			var before entryState
			if !tc.installed {
				before = observeEntry(t, root, lifecycleSeed)
			}
			runSetupYes(t)
			if !tc.installed {
				requireKept(t, root, lifecycleSeed, before)
				return
			}
			if got, row := readFile(t, path), manifestHash(root, lifecycleSeed); got != scaffoldGateInputs() || row != "" {
				t.Fatalf("%s = %q with row %q, want the scaffold bytes and no row", lifecycleSeed, got, row)
			}
		})
	}
}

// TestLifecycleHardRefusalOrder drives the shared execution caller with plans whose write
// needs a symlink parent. The refusal wins over the later entries, and nothing publishes.
func TestLifecycleHardRefusalOrder(t *testing.T) {
	accepted := planEntry{rel: "accepted.txt", kind: "inline", content: "accepted\n"}
	asset := planEntry{rel: "linked/asset", kind: "inline", content: "asset\n"}
	missing := planEntry{src: filepath.Join(t.TempDir(), "absent-kit-source"), rel: "missing.txt", kind: "file"}
	for _, tc := range []struct {
		name, stderr string
		plan         []planEntry
	}{
		{name: "write", stderr: "conflict: linked/asset has a symlink parent directory\n", plan: []planEntry{accepted, asset}},
		{name: "competing refusals", stderr: "conflict: linked/asset has a symlink parent directory\n", plan: []planEntry{accepted, asset, missing}},
		{name: "present seed", stderr: "conflict: linked/seed has a symlink parent directory\n",
			plan: []planEntry{accepted, {rel: "linked/seed", kind: "seed", content: "seed\n"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			runAdoptGit(t, root, "init", "-q")
			writeFixtureFile(t, filepath.Join(root, "AGENTS.md"), "# Project rules\n", 0o644)
			writeFixtureFile(t, filepath.Join(root, "CLAUDE.md"), "Project notes.\n", 0o644)
			if err := os.Mkdir(filepath.Join(root, "target"), 0o755); err != nil {
				t.Fatal(err)
			}
			writeFixtureFile(t, filepath.Join(root, "target", "seed"), "project seed\n", 0o644)
			if err := os.Symlink("target", filepath.Join(root, "linked")); err != nil {
				t.Fatal(err)
			}
			before := compatibilitytest.SnapshotHome(t, root)
			var stdout, stderr bytes.Buffer
			code, changed := transactionalLink(root, t.TempDir(), "copy", "test", tc.plan, &stdout, &stderr)
			if code != 1 || changed || stderr.String() != tc.stderr {
				t.Fatalf("transactionalLink = (%d, %t), stderr %q; want (1, false), stderr %q", code, changed, stderr.String(), tc.stderr)
			}
			if after := compatibilitytest.SnapshotHome(t, root); !reflect.DeepEqual(after, before) {
				t.Fatalf("refused plan changed the repository:\n%+v\nwant:\n%+v", after, before)
			}
		})
	}
}

// TestLifecycleRepairCompatibility drives real compatibility repair over an owned asset
// whose bytes are clean but whose mode a project changed. Strict repair keeps the asset
// and reports it, where ordinary relink would restore the mode.
func TestLifecycleRepairCompatibility(t *testing.T) {
	root := consumerRepo(t)
	linkConsumer(t, root, 0)
	if err := os.Chmod(filepath.Join(root, lifecycleDoc), 0o600); err != nil {
		t.Fatal(err)
	}
	before := observeEntry(t, root, lifecycleDoc)
	code, stdout, stderr := runCompatibilityRepair(t)
	if code != 1 || !strings.Contains(stdout, lifecycleDoc+",modified-managed\n") || !strings.Contains(stdout, ",incomplete,") {
		t.Fatalf("compatibility repair = (%d, %q, %q), want exit 1 with the modified-managed row and an incomplete repair", code, stdout, stderr)
	}
	requireKept(t, root, lifecycleDoc, before)
}
