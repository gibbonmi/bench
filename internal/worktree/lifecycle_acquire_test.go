package worktree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestAcquireCreatesPrivatePoolAndLease(t *testing.T) {
	oldUmask := syscall.Umask(0)
	t.Cleanup(func() { syscall.Umask(oldUmask) })
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	wt, err := Acquire(root, "", "")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	t.Cleanup(func() { Release(wt) })
	assertMode := func(path string, want os.FileMode) {
		t.Helper()
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %s: %v", path, err)
		}
		if got := info.Mode().Perm(); got != want {
			t.Errorf("mode %s = %04o, want %04o", path, got, want)
		}
	}
	assertMode(Pool(root), 0o700)
	lease, err := LeaseFile(wt)
	if err != nil {
		t.Fatalf("LeaseFile: %v", err)
	}
	assertMode(lease, 0o600)
}

func TestAcquireTightensExistingPool(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	pool := Pool(root)
	if err := os.MkdirAll(pool, 0o777); err != nil {
		t.Fatalf("mkdir loose pool: %v", err)
	}
	if err := os.Chmod(pool, 0o777); err != nil {
		t.Fatalf("chmod loose pool: %v", err)
	}
	wt, err := Acquire(root, "", "")
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	t.Cleanup(func() { Release(wt) })
	info, err := os.Stat(pool)
	if err != nil {
		t.Fatalf("stat pool: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Fatalf("pool mode after Acquire = %04o, want 0700", got)
	}
}

// TestAcquireWithUnresolvableDefaultAddsAtHead covers the empty-remote-ref end of the
// pool-minting fallback. With no default branch to start from, the first add is already
// the HEAD one. So the mint still succeeds rather than spending its attempt twice.
func TestAcquireWithUnresolvableDefaultAddsAtHead(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	gitRun(t, root, "branch", "-M", "master")
	gitRun(t, root, "branch", "feature")
	if ref := git.RemoteDefaultRef(root); ref != "" {
		t.Fatalf("fixture default resolved to %q, want no resolvable default", ref)
	}

	wt, err := Acquire(root, "", "")

	if err != nil {
		t.Fatalf("Acquire with an unresolvable default: %v", err)
	}
	t.Cleanup(func() { Release(wt) })
	if head := gitOutput(t, wt, "rev-parse", "HEAD"); head != gitOutput(t, root, "rev-parse", "HEAD") {
		t.Fatalf("pool worktree HEAD = %q, want the repository HEAD", head)
	}
}

func TestAcquireContinuesWhenPoolTightenFails(t *testing.T) {
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	root := newWorktreeRepo(t)
	pool := Pool(root)
	old := chmodPool
	called := false
	chmodPool = func(path string, mode os.FileMode) error {
		if path == pool {
			called = true
			return os.ErrPermission
		}
		return os.Chmod(path, mode)
	}
	t.Cleanup(func() { chmodPool = old })
	wt, err := Acquire(root, "", "")
	if err != nil {
		t.Fatalf("Acquire after pool chmod failure: %v", err)
	}
	t.Cleanup(func() { Release(wt) })
	if !called {
		t.Fatal("Acquire did not attempt to tighten the pool")
	}
}

func TestRecoveryPreservesEveryGitVisibleLayerWithoutMovingBranchOrIndex(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*testing.T, string)
		wantLayers []string
		check      func(*testing.T, string, map[string]string)
	}{
		{"staged-only", func(t *testing.T, path string) {
			mustWrite(t, filepath.Join(path, "README.md"), []byte("staged\n"), 0o644)
			gitRun(t, path, "add", "README.md")
		}, []string{"staged", "working"}, func(t *testing.T, root string, layers map[string]string) {
			requireBlob(t, root, layers["staged"], "README.md", "staged\n")
		}},
		{"unstaged-only", func(t *testing.T, path string) {
			mustWrite(t, filepath.Join(path, "README.md"), []byte("working\n"), 0o644)
		}, []string{"working"}, func(t *testing.T, root string, layers map[string]string) {
			requireBlob(t, root, layers["working"], "README.md", "working\n")
		}},
		{"divergent-staged-working", func(t *testing.T, path string) {
			mustWrite(t, filepath.Join(path, "README.md"), []byte("staged\n"), 0o644)
			gitRun(t, path, "add", "README.md")
			mustWrite(t, filepath.Join(path, "README.md"), []byte("working\n"), 0o644)
		}, []string{"staged", "working"}, func(t *testing.T, root string, layers map[string]string) {
			if layers["staged"] == layers["working"] {
				t.Fatal("divergent staged and working layers collapsed")
			}
			requireBlob(t, root, layers["staged"], "README.md", "staged\n")
			requireBlob(t, root, layers["working"], "README.md", "working\n")
		}},
		{"deleted", func(t *testing.T, path string) { mustRemove(t, filepath.Join(path, "README.md")) }, []string{"working"}, func(t *testing.T, root string, layers map[string]string) {
			requireMissingBlob(t, root, layers["working"], "README.md")
		}},
		{"untracked", func(t *testing.T, path string) { mustWrite(t, filepath.Join(path, "new.txt"), []byte("new\n"), 0o644) }, []string{"working"}, func(t *testing.T, root string, layers map[string]string) {
			requireBlob(t, root, layers["working"], "new.txt", "new\n")
		}},
		{"renamed", func(t *testing.T, path string) { gitRun(t, path, "mv", "README.md", "renamed.md") }, []string{"staged", "working"}, func(t *testing.T, root string, layers map[string]string) {
			requireBlob(t, root, layers["working"], "renamed.md", "initial\n")
			requireMissingBlob(t, root, layers["working"], "README.md")
		}},
		{"symlink", func(t *testing.T, path string) {
			mustRemove(t, filepath.Join(path, "README.md"))
			if err := os.Symlink("target.txt", filepath.Join(path, "README.md")); err != nil {
				t.Fatal(err)
			}
		}, []string{"working"}, func(t *testing.T, root string, layers map[string]string) {
			entry := gitOutput(t, root, "ls-tree", layers["working"], "README.md")
			if !strings.HasPrefix(entry, "120000 blob ") {
				t.Fatalf("symlink tree entry = %q", entry)
			}
		}},
		{"executable-mode", func(t *testing.T, path string) { mustChmod(t, filepath.Join(path, "README.md"), 0o755) }, []string{"working"}, func(t *testing.T, root string, layers map[string]string) {
			entry := gitOutput(t, root, "ls-tree", layers["working"], "README.md")
			if !strings.HasPrefix(entry, "100755 blob ") {
				t.Fatalf("executable tree entry = %q", entry)
			}
		}},
		{"conflicted", setupConflict, []string{"base", "ours", "theirs", "working"}, func(t *testing.T, root string, layers map[string]string) {
			requireBlob(t, root, layers["base"], "README.md", "base\n")
			requireBlob(t, root, layers["ours"], "README.md", "ours\n")
			requireBlob(t, root, layers["theirs"], "README.md", "theirs\n")
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root, creation := newOwnedAssignment(t, "recovery-"+tc.name)
			tc.setup(t, creation.Path)
			markPending(t, root, creation.Assignment)
			indexPath := gitOutput(t, creation.Path, "rev-parse", "--path-format=absolute", "--git-path", "index")
			indexBefore, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatal(err)
			}
			branchBefore := gitOutput(t, root, "rev-parse", creation.Assignment.Branch)
			// Preservation belongs to the explicit path-addressed clean. The automatic planner
			// an unattended resume or release drives retains a checkout it could only remove
			// by preserving first. So the layer capture runs through the surface that still
			// reaches it.
			restore := cleanupTransactionBoundary
			defer func() { cleanupTransactionBoundary = restore }()
			plan, err := PlanExplicit(root, creation.Path)
			if err != nil {
				t.Fatal(err)
			}
			stop := errors.New("stop after recovery metadata")
			cleanupTransactionBoundary = failLifecycleStep(StepRecoveryMetadata, stop)
			_, err = ApplyExplicit(root, creation.Path, plan.Fingerprint)
			if !errors.Is(err, stop) {
				t.Fatalf("ApplyExplicit error = %v, want recovery-metadata fault", err)
			}
			assignments, readErr := intent.Assignments(root)
			if readErr != nil || len(assignments) != 1 || len(assignments[0].Recovery) != 1 {
				t.Fatalf("recovery metadata = %#v, %v", assignments, readErr)
			}
			if exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", assignments[0].Recovery[0].Ref).Run() == nil {
				t.Fatal("recovery ref exists before its metadata checkpoint replays")
			}
			stop = errors.New("stop after durable recovery ref")
			replayPlan, err := PlanExplicit(root, creation.Path)
			if err != nil {
				t.Fatal(err)
			}
			cleanupTransactionBoundary = failLifecycleStep(StepRecoveryRef, stop)
			_, err = ApplyExplicit(root, creation.Path, replayPlan.Fingerprint)
			if !errors.Is(err, stop) {
				t.Fatalf("ApplyExplicit replay error = %v, want recovery-ref fault", err)
			}
			cleanupTransactionBoundary = restore
			indexAfter, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(indexBefore, indexAfter) {
				t.Fatal("recovery changed the real index")
			}
			if after := gitOutput(t, root, "rev-parse", creation.Assignment.Branch); after != branchBefore {
				t.Fatalf("assignment branch moved from %s to %s", branchBefore, after)
			}
			assignments, err = intent.Assignments(root)
			if err != nil || len(assignments) != 1 || len(assignments[0].Recovery) != 1 {
				t.Fatalf("recovery assignment = %#v, %v", assignments, err)
			}
			recovery := assignments[0].Recovery[0]
			if got := gitOutput(t, root, "rev-parse", recovery.Ref); got != recovery.Root {
				t.Fatalf("recovery ref = %s, want root %s", got, recovery.Root)
			}
			manifestRaw := gitOutput(t, root, "show", recovery.Root+":manifest.json")
			var manifest struct {
				Layers map[string]string `json:"layers"`
			}
			if err := json.Unmarshal([]byte(manifestRaw), &manifest); err != nil {
				t.Fatalf("manifest: %v\n%s", err, manifestRaw)
			}
			for _, layer := range tc.wantLayers {
				if manifest.Layers[layer] == "" {
					t.Fatalf("manifest lacks %s layer: %#v", layer, manifest.Layers)
				}
			}
			for _, payload := range recovery.Payloads {
				gitRun(t, root, "merge-base", "--is-ancestor", payload, recovery.Root)
			}
			tc.check(t, root, manifest.Layers)
			registration := gitOutput(t, root, "worktree", "list", "--porcelain")
			if !strings.Contains(registration, "worktree "+creation.Path) || !strings.Contains(registration, "locked "+lockReason(assignments[0])) {
				t.Fatal("recovery-ref fault did not leave a locked attributable checkout")
			}
		})
	}
}

func TestExplicitApplyRevalidatesSafetyEvidence(t *testing.T) {
	type fixture struct {
		root     string
		creation Creation
	}
	newFixture := func(t *testing.T) fixture {
		t.Helper()
		root := newWorktreeRepo(t)
		t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
		creation, err := Create(root, "drift-"+strings.ReplaceAll(t.Name(), "/", "-"), "drift", nil)
		if err != nil {
			t.Fatal(err)
		}
		return fixture{root: root, creation: creation}
	}
	tests := []struct {
		name   string
		setup  func(*testing.T, fixture)
		mutate func(*testing.T, fixture)
	}{
		{name: "default ref oid", mutate: func(t *testing.T, f fixture) {
			mustWrite(t, filepath.Join(f.root, "default-drift.txt"), []byte("drift\n"), 0o644)
			gitRun(t, f.root, "add", "default-drift.txt")
			gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "default drift")
		}},
		{name: "marker bytes", mutate: func(t *testing.T, f fixture) {
			path, _ := markerPath(f.creation.Path)
			body, _ := os.ReadFile(path)
			mustWrite(t, path, append(body, ' '), 0o600)
		}},
		{name: "marker mode", mutate: func(t *testing.T, f fixture) {
			path, _ := markerPath(f.creation.Path)
			if err := os.Chmod(path, 0o640); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "assignment state", mutate: func(t *testing.T, f fixture) {
			a := f.creation.Assignment
			a.State = intent.StateCleanupPending
			if err := intent.PutAssignment(f.root, a); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "assignment start", mutate: func(t *testing.T, f fixture) {
			a := f.creation.Assignment
			a.Start = gitOutput(t, f.root, "commit-tree", "HEAD^{tree}", "-p", "HEAD", "-m", "alternate start")
			if err := intent.PutAssignment(f.root, a); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "head and branch ref", mutate: func(t *testing.T, f fixture) {
			gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "--allow-empty", "-qm", "head drift")
		}},
		{name: "index", mutate: func(t *testing.T, f fixture) {
			gitRun(t, f.creation.Path, "update-index", "--chmod=+x", "README.md")
		}},
		{name: "lease", mutate: func(t *testing.T, f fixture) {
			lease, _ := LeaseFile(f.creation.Path)
			mustWrite(t, lease, []byte("live lease\n"), 0o600)
		}},
		{name: "full lock reason", mutate: func(t *testing.T, f fixture) {
			gitRun(t, f.root, "worktree", "unlock", f.creation.Path)
			gitRun(t, f.root, "worktree", "lock", "--reason", "changed full lock reason", f.creation.Path)
		}},
		{name: "nested state", setup: func(t *testing.T, f fixture) {
			nested := filepath.Join(f.creation.Path, "nested")
			if err := os.MkdirAll(nested, 0o755); err != nil {
				t.Fatal(err)
			}
			gitRun(t, nested, "init", "-q", "-b", "main")
			mustWrite(t, filepath.Join(nested, "n.txt"), []byte("base\n"), 0o644)
			gitRun(t, nested, "add", "n.txt")
			gitRun(t, nested, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "base")
		}, mutate: func(t *testing.T, f fixture) {
			mustWrite(t, filepath.Join(f.creation.Path, "nested", "n.txt"), []byte("dirty\n"), 0o644)
		}},
		{name: "ignored inventory", setup: func(t *testing.T, f fixture) {
			mustWrite(t, filepath.Join(f.root, ".git", "info", "exclude"), []byte("ignored.txt\n"), 0o644)
		}, mutate: func(t *testing.T, f fixture) {
			mustWrite(t, filepath.Join(f.creation.Path, "ignored.txt"), []byte("metadata only\n"), 0o644)
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			if tc.setup != nil {
				tc.setup(t, f)
			}
			plan, err := PlanExplicit(f.root, f.creation.Path)
			if err != nil || plan.Fingerprint == "" {
				t.Fatalf("plan before drift = %#v, %v", plan, err)
			}
			tc.mutate(t, f)
			before := explicitDurableState(t, f.root, f.creation.Path)
			current, err := ApplyExplicit(f.root, f.creation.Path, plan.Fingerprint)
			if !errors.Is(err, errStaleFingerprint) || current.Fingerprint == plan.Fingerprint {
				t.Fatalf("apply old fingerprint = %#v, %v; want current non-mutating plan", current, err)
			}
			if after := explicitDurableState(t, f.root, f.creation.Path); after != before {
				t.Fatalf("stale apply mutated state\nbefore=%q\nafter=%q", before, after)
			}
		})
	}
}
func explicitDurableState(t *testing.T, root, target string) string {
	t.Helper()
	ledger, _ := intent.Address(root)
	marker, _ := markerPath(target)
	lease, _ := LeaseFile(target)
	var state bytes.Buffer
	for _, path := range []string{ledger, marker, lease} {
		body, err := os.ReadFile(path)
		fmt.Fprintf(&state, "%s:%v:%x\n", path, err, body)
	}
	state.WriteString(gitOutput(t, root, "worktree", "list", "--porcelain"))
	state.WriteString(gitOutput(t, root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/recovery/"))
	return state.String()
}
func setupConflict(t *testing.T, path string) {
	t.Helper()
	base := hashBlob(t, path, "base\n")
	ours := hashBlob(t, path, "ours\n")
	theirs := hashBlob(t, path, "theirs\n")
	mustWrite(t, filepath.Join(path, "README.md"), []byte("ours\n"), 0o644)
	input := "100644 " + base + " 1\tREADME.md\n" +
		"100644 " + ours + " 2\tREADME.md\n" +
		"100644 " + theirs + " 3\tREADME.md\n"
	cmd := exec.Command("git", "-C", path, "update-index", "--index-info")
	cmd.Stdin = strings.NewReader(input)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("create conflict index: %v\n%s", err, out)
	}
}
func hashBlob(t *testing.T, root, contents string) string {
	t.Helper()
	cmd := exec.Command("git", "-C", root, "hash-object", "-w", "--stdin")
	cmd.Stdin = strings.NewReader(contents)
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}
func requireBlob(t *testing.T, root, commit, path, want string) {
	t.Helper()
	if got := gitOutput(t, root, "show", commit+":"+path); got+"\n" != want {
		t.Fatalf("%s:%s = %q, want %q", commit, path, got+"\n", want)
	}
}
func requireMissingBlob(t *testing.T, root, commit, path string) {
	t.Helper()
	if exec.Command("git", "-C", root, "cat-file", "-e", commit+":"+path).Run() == nil {
		t.Fatalf("%s unexpectedly contains %s", commit, path)
	}
}
