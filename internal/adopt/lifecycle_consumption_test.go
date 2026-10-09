package adopt

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/compatibility/compatibilitytest"
)

// The ordinary managed assets that the real Link fixtures change: a regular kit file, and
// an executable kit file for the mode fixtures. lifecycleWithdrawn is a path that the
// manifest owns and the current kit no longer plans.
const (
	lifecycleDoc       = ".bench/BENCH.md"
	lifecycleExec      = ".bench/bin/bench.sh"
	lifecycleWithdrawn = ".bench/withdrawn.md"
)

// linkRun is the exit and output of one real Link.
type linkRun struct {
	code           int
	stdout, stderr string
}

func runLink(args []string) linkRun {
	var stdout, stderr bytes.Buffer
	code := Link(args, &stdout, &stderr, "1.0.0")
	return linkRun{code, stdout.String(), stderr.String()}
}

// runLinkBounded runs real Link off the test goroutine behind a deadline, because opening
// a writerless FIFO never returns.
func runLinkBounded(t *testing.T) linkRun {
	t.Helper()
	done := make(chan linkRun, 1)
	go func() { done <- runLink(nil) }()
	select {
	case run := <-done:
		return run
	case <-time.After(bounds.TestDeadline(bounds.TestDeadlineFloor)):
		t.Fatal("Link blocked: it opened a special file before it checked the file kind")
	}
	return linkRun{}
}

// requireReported fails unless run exited 0, or, when reason is set, exited 3 with the
// reason row for rel.
func requireReported(t *testing.T, run linkRun, rel, reason string) {
	t.Helper()
	code := 0
	if reason != "" {
		code = 3
	}
	if run.code != code || reason != "" && !strings.Contains(run.stdout, rel+","+reason+"\n") {
		t.Fatalf("Link = %d, want %d with the %q row for %s\nstdout:\n%s\nstderr:\n%s", run.code, code, reason, rel, run.stdout, run.stderr)
	}
}

// plantWithdrawn writes a withdrawn asset at path and returns the fingerprint that a
// manifest row records for it.
func plantWithdrawn(t *testing.T, path string) string {
	t.Helper()
	writeFixtureFile(t, path, "withdrawn kit asset\n", 0o644)
	recorded, err := fingerprintPath(path)
	if err != nil {
		t.Fatal(err)
	}
	return recorded
}

// requireWithdrawn fails unless the withdrawn asset and its manifest row are both gone.
func requireWithdrawn(t *testing.T, root string) {
	t.Helper()
	requireAbsent(t, filepath.Join(root, lifecycleWithdrawn))
	if row := manifestHash(root, lifecycleWithdrawn); row != "" {
		t.Fatalf("%s manifest row = %q, want no row", lifecycleWithdrawn, row)
	}
}

// requireClaudeInstalled fails unless CLAUDE.md holds the current managed form and the
// manifest records its fingerprint.
func requireClaudeInstalled(t *testing.T, root string) {
	t.Helper()
	want := benchClaudeMD()
	if got, row := readFile(t, filepath.Join(root, "CLAUDE.md")), manifestHash(root, "CLAUDE.md"); got != want || row != hashBytes([]byte(want)) {
		t.Fatalf("CLAUDE.md = %q with row %q, want the current form and its fingerprint", got, row)
	}
}

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
// staged link names, so it stays a regular file. A consumer root whose name holds a space
// and glob characters takes the absent asset at its literal path.
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
		name, rel, reason, dir string
		args                   []string
		linked, installed      bool
		arrange                func(t *testing.T, root, path string)
	}{
		{name: "absent", rel: lifecycleDoc, installed: true, arrange: func(*testing.T, string, string) {}},
		{name: "absent under a spaced glob root", rel: lifecycleDoc, dir: "dir with [glob]*", installed: true, arrange: func(*testing.T, string, string) {}},
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
			if tc.dir != "" {
				root = filepath.Join(root, tc.dir)
				if err := os.Mkdir(root, 0o755); err != nil {
					t.Fatal(err)
				}
				runAdoptGit(t, root, "init", "-q")
				t.Chdir(root)
			}
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
			requireReported(t, runLink(tc.args), tc.rel, tc.reason)
			if tc.installed {
				requireInstalled(t, root, tc.rel)
			} else {
				requireKept(t, root, tc.rel, before)
			}
		})
	}
}

// TestLifecycleOwnerConsumption drives each applying entry that consumes the complete
// decision. Each linking entry starts from a manifest that pins the linked kit and owns
// only a clean withdrawn asset, next to project AGENTS.md prose and no CLAUDE.md. Each
// linking entry must install an absent ordinary asset and the current CLAUDE.md form,
// withdraw the clean asset with its row, and add the managed block after the project
// prose. Unlink starts from a linked repository with one modified asset. It must remove a
// clean asset, keep the modified one with its row, and leave the project prose alone.
func TestLifecycleOwnerConsumption(t *testing.T) {
	const prose = "# Team rules\n\nKeep this rule.\n"
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
			if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
				t.Fatal(err)
			}
			writeFixtureFile(t, filepath.Join(root, "AGENTS.md"), prose, 0o644)
			withdrawn := manifestRow{lifecycleWithdrawn, plantWithdrawn(t, filepath.Join(root, lifecycleWithdrawn))}
			if err := writeManifest(previewManifestPath(root), "1.0.0", []manifestRow{withdrawn}); err != nil {
				t.Fatal(err)
			}
			entry.run(t, root)
			requireInstalled(t, root, lifecycleDoc)
			requireWithdrawn(t, root)
			requireClaudeInstalled(t, root)
			if got := readFile(t, filepath.Join(root, "AGENTS.md")); !strings.HasPrefix(got, prose) || !strings.Contains(got, BenchAgentsBlock()) {
				t.Fatalf("AGENTS.md = %q, want the project prose followed by the managed block", got)
			}
		})
	}
	t.Run("unlink", func(t *testing.T) {
		root := consumerRepo(t)
		writeFixtureFile(t, filepath.Join(root, "AGENTS.md"), prose, 0o644)
		linkConsumer(t, root, 0)
		writeFixtureFile(t, filepath.Join(root, lifecycleDoc), "project edit\n", 0o644)
		before := observeEntry(t, root, lifecycleDoc)
		runUnlink(t, nil, 3)
		requireAbsent(t, filepath.Join(root, lifecycleExec))
		requireKept(t, root, lifecycleDoc, before)
		if got := readFile(t, filepath.Join(root, "AGENTS.md")); got != prose {
			t.Fatalf("AGENTS.md after unlink = %q, want the project prose %q", got, prose)
		}
	})
}

// TestLifecycleDroppedRows drives real relink over a manifest row that the kit no longer
// plans. A clean withdrawn asset leaves with its row. A modified one keeps its bytes and
// row and is reported. An absent one loses its row without a report. A row under a
// symlink parent, a row that escapes the repository, and a special file each refuse the
// relink with the exact refusal and keep the planted file and its row.
func TestLifecycleDroppedRows(t *testing.T) {
	keep := func(*testing.T, string) {}
	for _, tc := range []struct {
		name, rel, reason, refusal string
		symlinkParent              bool
		arrange                    func(t *testing.T, path string)
	}{
		{name: "clean", arrange: keep},
		{name: "modified", reason: "kept-modified-removed", arrange: func(t *testing.T, path string) {
			writeFixtureFile(t, path, "project edit\n", 0o644)
		}},
		{name: "absent", arrange: func(t *testing.T, path string) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "symlink parent", rel: "linked/withdrawn.md", symlinkParent: true,
			refusal: "refused manifest path: linked/withdrawn.md\n", arrange: keep},
		{name: "escape", rel: "../withdrawn.md", refusal: "refused manifest path: ../withdrawn.md\n", arrange: keep},
		{name: "special file", refusal: "refused non-regular dropped asset: " + lifecycleWithdrawn + "\n", arrange: func(t *testing.T, path string) {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := consumerRepo(t)
			linkConsumer(t, root, 0)
			rel := tc.rel
			if rel == "" {
				rel = lifecycleWithdrawn
			}
			if tc.symlinkParent {
				plantSymlinkParent(t, root)
			}
			path := filepath.Join(root, rel)
			recorded := plantWithdrawn(t, path)
			editPreviewManifest(t, root, func(rows map[string]string) { rows[rel] = recorded })
			tc.arrange(t, path)
			var before entryState
			if tc.reason != "" || tc.refusal != "" {
				before = observeEntry(t, root, rel)
			}
			run := runLinkBounded(t)
			if tc.refusal != "" {
				if run.code != 1 || run.stderr != tc.refusal {
					t.Fatalf("Link = %d with stderr %q, want 1 with %q", run.code, run.stderr, tc.refusal)
				}
			} else {
				requireReported(t, run, rel, tc.reason)
			}
			if tc.reason != "" || tc.refusal != "" {
				requireKept(t, root, rel, before)
				return
			}
			requireWithdrawn(t, root)
		})
	}
}
