package adopt

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
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

// plantSymlinkParent makes root/linked a symlink to the directory root/target.
func plantSymlinkParent(t *testing.T, root string) {
	t.Helper()
	if err := os.Mkdir(filepath.Join(root, "target"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("target", filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
}

// requireUnchanged runs run and fails unless every entry under root keeps its state.
func requireUnchanged(t *testing.T, root string, run func()) {
	t.Helper()
	before := compatibilitytest.SnapshotHome(t, root)
	run()
	if after := compatibilitytest.SnapshotHome(t, root); !reflect.DeepEqual(after, before) {
		t.Fatalf("refused run changed the repository:\n%+v\nwant:\n%+v", after, before)
	}
}

// TestLifecycleHardRefusalOrder drives the shared execution caller with plans whose write
// needs a symlink parent. The refused entry follows an accepted one. The refusal wins over
// the later entries, and nothing publishes. A seed refuses whether its leaf exists or not.
func TestLifecycleHardRefusalOrder(t *testing.T) {
	accepted := planEntry{rel: "accepted.txt", kind: "inline", content: "accepted\n"}
	asset := planEntry{rel: "linked/asset", kind: "inline", content: "asset\n"}
	missing := planEntry{src: filepath.Join(t.TempDir(), "absent-kit-source"), rel: "missing.txt", kind: "file"}
	seed := func(rel string) planEntry { return planEntry{rel: rel, kind: "seed", content: "seed\n"} }
	for _, tc := range []struct {
		name    string
		refused planEntry
		later   []planEntry
	}{
		{name: "write", refused: asset},
		{name: "competing refusals", refused: asset, later: []planEntry{missing}},
		{name: "present seed", refused: seed("linked/seed")},
		{name: "absent seed", refused: seed("linked/absent-seed")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			runAdoptGit(t, root, "init", "-q")
			writeFixtureFile(t, filepath.Join(root, "AGENTS.md"), "# Project rules\n", 0o644)
			writeFixtureFile(t, filepath.Join(root, "CLAUDE.md"), "Project notes.\n", 0o644)
			plantSymlinkParent(t, root)
			writeFixtureFile(t, filepath.Join(root, "target", "seed"), "project seed\n", 0o644)
			plan := append([]planEntry{accepted, tc.refused}, tc.later...)
			var stdout, stderr bytes.Buffer
			var code int
			var changed bool
			requireUnchanged(t, root, func() {
				code, changed = transactionalLink(root, t.TempDir(), "copy", "test", plan, &stdout, &stderr)
			})
			if want := "conflict: " + tc.refused.rel + " has a symlink parent directory\n"; code != 1 || changed || stderr.String() != want {
				t.Fatalf("transactionalLink = (%d, %t), stderr %q; want (1, false), stderr %q", code, changed, stderr.String(), want)
			}
		})
	}
}

// TestLifecycleRepairCompatibility drives real compatibility repair over an owned asset
// whose bytes are clean but whose mode a project changed. Strict repair keeps the asset
// and reports it, where ordinary relink would restore the mode.
func TestLifecycleRepairCompatibility(t *testing.T) {
	root := linkedCopy(t, goldenPlain)
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

// TestLifecycleClaudeOwnership drives real relink over each CLAUDE.md partition, starting
// from a linked repository. Project content, including an empty file, keeps its bytes and
// takes no import and no row. Owned content that the project edited keeps its bytes and
// row and is reported. An absent, legacy, or current form takes the current form, and a
// dangling link, owned or not, reads as absent.
func TestLifecycleClaudeOwnership(t *testing.T) {
	write := func(content string) func(*testing.T, string, string) {
		return func(t *testing.T, _, path string) { writeFixtureFile(t, path, content, 0o644) }
	}
	project := func(content string) func(*testing.T, string, string) {
		return func(t *testing.T, root, path string) {
			write(content)(t, root, path)
			editPreviewManifest(t, root, func(rows map[string]string) { delete(rows, "CLAUDE.md") })
		}
	}
	dangling := func(owned bool) func(*testing.T, string, string) {
		return func(t *testing.T, root, path string) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("missing-target", path); err != nil {
				t.Fatal(err)
			}
			if !owned {
				editPreviewManifest(t, root, func(rows map[string]string) { delete(rows, "CLAUDE.md") })
			}
		}
	}
	for _, tc := range []struct {
		name, reason string
		installed    bool
		arrange      func(t *testing.T, root, path string)
	}{
		{name: "project-authored", arrange: project("Project notes.\n")},
		{name: "empty", arrange: project("")},
		{name: "modified owned", reason: "modified-managed", arrange: write("Project notes.\n")},
		{name: "absent", installed: true, arrange: func(t *testing.T, _, path string) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "legacy form", installed: true, arrange: write(legacyClaudeMD())},
		{name: "current form", installed: true, arrange: func(*testing.T, string, string) {}},
		{name: "unowned dangling link", installed: true, arrange: dangling(false)},
		{name: "owned dangling link", installed: true, arrange: dangling(true)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := linkedCopy(t, goldenPlain)
			tc.arrange(t, root, filepath.Join(root, "CLAUDE.md"))
			var before entryState
			if !tc.installed {
				before = observeEntry(t, root, "CLAUDE.md")
			}
			requireReported(t, runLink(nil), "CLAUDE.md", tc.reason)
			if tc.installed {
				requireClaudeInstalled(t, root)
			} else {
				requireKept(t, root, "CLAUDE.md", before)
			}
		})
	}
}

// TestLifecycleSpecialInstructionFiles drives real Link with a writerless FIFO at an
// instruction file. Link must finish without opening it, report it project-owned, keep
// its file kind, and drop an ownership row that an earlier link wrote.
func TestLifecycleSpecialInstructionFiles(t *testing.T) {
	for _, tc := range []struct {
		name, rel string
		owned     bool
	}{
		{name: "AGENTS.md", rel: "AGENTS.md"},
		{name: "CLAUDE.md", rel: "CLAUDE.md"},
		{name: "owned CLAUDE.md", rel: "CLAUDE.md", owned: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var root string
			if tc.owned {
				root = linkedCopy(t, goldenPlain)
			} else {
				root = consumerRepo(t)
			}
			path := filepath.Join(root, tc.rel)
			if tc.owned {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			}
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				t.Fatal(err)
			}
			requireReported(t, runLinkBounded(t), tc.rel, "project-owned")
			if info, err := os.Lstat(path); err != nil || info.Mode()&os.ModeNamedPipe == 0 {
				t.Fatalf("%s after Link: %v, want the FIFO kept", tc.rel, err)
			}
			if row := manifestHash(root, tc.rel); row != "" {
				t.Fatalf("%s manifest row = %q, want no row", tc.rel, row)
			}
		})
	}
}

// hookRepo points core.hooksPath of root at its own hooks directory when diverted is set,
// and returns root with the effective pre-push path. A copy of a diverted golden repository
// takes the setting again, because the copied setting names the golden directory.
func hookRepo(t *testing.T, root string, diverted bool) (string, string) {
	t.Helper()
	hooks := filepath.Join(root, ".git", "hooks")
	if diverted {
		hooks = filepath.Join(root, "hooks")
		if err := os.MkdirAll(hooks, 0o755); err != nil {
			t.Fatal(err)
		}
		runAdoptGit(t, root, "config", "core.hooksPath", hooks)
	}
	return root, filepath.Join(hooks, "pre-push")
}

// TestLifecycleHookCompatibility drives real Link over each effective pre-push state. An
// absent or stale managed hook takes the current managed hook, also where core.hooksPath
// diverts the hooks directory. A foreign hook, or any hook in a diverted directory that
// Bench does not manage, refuses the whole link and leaves the repository unchanged.
// Real Unlink then meets each state in a linked repository. It removes and reports a hook
// that carries the managed marker. It keeps a foreign hook without a refusal and leaves an
// absent hook absent, and it reports neither.
func TestLifecycleHookCompatibility(t *testing.T) {
	foreign := func(t *testing.T, path string) { writeFixtureFile(t, path, "#!/bin/sh\nexit 0\n", 0o755) }
	goldens := map[bool]*golden{false: goldenPlain, true: goldenDiverted}
	for _, tc := range []struct {
		name     string
		diverted bool
		refused  bool
		absent   bool
		arrange  func(t *testing.T, path string)
	}{
		{name: "absent", absent: true, arrange: func(*testing.T, string) {}},
		{name: "stale managed", arrange: func(t *testing.T, path string) {
			writeFixtureFile(t, path, strings.TrimSuffix(renderPrePushBranch("baked"), "\n"), 0o755)
		}},
		{name: "diverted absent", diverted: true, absent: true, arrange: func(*testing.T, string) {}},
		{name: "foreign", refused: true, arrange: foreign},
		{name: "diverted foreign", diverted: true, refused: true, arrange: foreign},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, path := hookRepo(t, consumerRepo(t), tc.diverted)
			tc.arrange(t, path)
			if !tc.refused {
				requireReported(t, runLink(nil), "pre-push", "")
				if got := InspectPrePush(root); got.State != PrePushManaged || got.Currency != PrePushCurrent || got.Path != path {
					t.Fatalf("hook after Link = %#v, want the current managed hook at %s", got, path)
				}
				return
			}
			var run linkRun
			requireUnchanged(t, root, func() { run = runLink(nil) })
			if want := foreignHookConflict(path); run.code != 1 || run.stderr != want {
				t.Fatalf("Link = %d with stderr %q, want 1 with %q", run.code, run.stderr, want)
			}
		})
		t.Run("unlink "+tc.name, func(t *testing.T) {
			_, path := hookRepo(t, linkedCopy(t, goldens[tc.diverted]), tc.diverted)
			tc.arrange(t, path)
			var before map[string]compatibilitytest.HomeEntry
			if tc.absent {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			} else {
				before = compatibilitytest.SnapshotHome(t, path)
			}
			out := runUnlink(t, nil, 0)
			if reported, want := strings.Contains(out, "pre-push"), !tc.refused && !tc.absent; reported != want {
				t.Fatalf("Unlink reports the hook: %t, want %t\n%s", reported, want, out)
			}
			if !tc.refused {
				requireAbsent(t, path)
			} else if after := compatibilitytest.SnapshotHome(t, path); !reflect.DeepEqual(after, before) {
				t.Fatalf("foreign hook after Unlink = %+v, want the kept %+v", after, before)
			}
		})
	}
}
