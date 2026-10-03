package worktree

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

func TestCleanLandedQuotesSpaceAndGlobPaths(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, "home space *")
	clean := mustCreate(t, root, home, "landed-hostile-clean", "clean")
	dirty := mustCreate(t, root, home, "landed-hostile-dirty", "dirty")
	landAssignment(t, root, clean, "clean.txt")
	landAssignment(t, root, dirty, "dirty.txt")
	mustWrite(t, filepath.Join(dirty.Path, "dirty.txt"), []byte("changed\n"), 0o644)

	plan := runVerb(t, verbClean, repoHome{root, home}.call("--landed"))
	if plan.exit != 0 || plan.stderr != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	if !strings.Contains(plan.stdout, clean.Path+",remove,") || !strings.Contains(plan.stdout, "bench worktree clean '"+dirty.Path+"'") || !strings.Contains(plan.stdout, "bench worktree clean --landed --apply ") {
		t.Fatalf("output=%q, want safe row and pasteable help", plan.stdout)
	}
	applied := runVerb(t, verbClean, repoHome{root, home}.call("--landed", "--apply", plan.mustFingerprint(t)))
	if applied.exit != 0 || applied.stderr != "" {
		t.Fatalf("apply exit=%d stderr=%q", applied.exit, applied.stderr)
	}
	if _, err := os.Lstat(clean.Path); !os.IsNotExist(err) {
		t.Fatalf("safe hostile path was not removed: %v", err)
	}
	if _, err := os.Lstat(dirty.Path); err != nil {
		t.Fatalf("dirty hostile path disappeared: %v", err)
	}
}

func TestCleanLandedControlBytePathRetained(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, "home\x1bunsafe")
	creation := mustCreate(t, root, home, "landed-control", "control")
	landAssignment(t, root, creation, "control.txt")

	plan := runVerb(t, verbClean, repoHome{root, home}.call("--landed"))
	pointer := "bench worktree exec " + creation.Assignment.ID + " -- bench worktree clean ."
	if plan.exit != 0 || plan.stderr != "" || strings.ContainsRune(plan.stdout, '\x1b') || !strings.Contains(plan.stdout, "sha256:") || !strings.Contains(plan.stdout, ",retain,") || !strings.Contains(plan.stdout, "unsafe control bytes") || strings.Count(plan.stdout, pointer) != 1 {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
	applied := runVerb(t, verbClean, repoHome{root, home}.call("--landed", "--apply", plan.mustFingerprint(t)))
	if applied.exit != 0 || applied.stderr != "" || strings.ContainsRune(applied.stdout, '\x1b') {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q", applied.exit, applied.stdout, applied.stderr)
	}
	if _, err := os.Lstat(creation.Path); err != nil {
		t.Fatalf("retained control-byte path disappeared: %v", err)
	}
}

func TestCleanLandedTabPathRendersOneRow(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, "home\tpath")
	creation := mustCreate(t, root, home, "landed-tab", "tab")
	landAssignment(t, root, creation, "tab.txt")

	plan := runVerb(t, verbClean, repoHome{root, home}.call("--landed"))
	if plan.exit != 0 || plan.stderr != "" || !strings.HasPrefix(plan.stdout, cleanupTable+"[1]") || strings.ContainsRune(plan.stdout, '\t') || !strings.Contains(plan.stdout, `\t`) {
		t.Fatalf("exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
	}
}

func TestCleanLandedSpecialPathsRetainedWithoutOpening(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, shape string
		make        func(*testing.T, string)
	}{
		{name: "fifo", shape: "non-directory", make: func(t *testing.T, path string) {
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
			}
		}},
		{name: "dangling symlink", shape: "dangling-symlink", make: func(t *testing.T, path string) {
			if err := os.Symlink(path+"-missing", path); err != nil {
				capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
			}
		}},
		{name: "socket", shape: "non-directory", make: func(t *testing.T, path string) {
			listener, err := net.Listen("unix", path)
			if err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("unix sockets unavailable: %v", err))
			}
			t.Cleanup(func() { _ = listener.Close() })
		}},
		{name: "device", shape: "non-directory", make: func(t *testing.T, path string) {
			if info, err := os.Lstat("/dev/null"); err != nil || info.Mode()&os.ModeDevice == 0 {
				capability.Capability(t, capability.Fifo, "no /dev/null device node")
			}
			if err := os.Symlink("/dev/null", path); err != nil {
				capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			home := filepath.Join(root, ".bench-home")
			creation := mustCreate(t, root, home, "landed-special-"+strings.ReplaceAll(tc.name, " ", "-"), tc.name)
			landAssignment(t, root, creation, "special.txt")
			if err := os.RemoveAll(creation.Path); err != nil {
				t.Fatal(err)
			}
			tc.make(t, creation.Path)

			plan := runVerb(t, verbClean, repoHome{root, home}.call("--landed"))
			if plan.exit != 0 || plan.stderr != "" || !strings.Contains(plan.stdout, ",retain,") || !strings.Contains(plan.stdout, "assignment path shape is "+tc.shape) {
				t.Fatalf("plan exit=%d stdout=%q stderr=%q", plan.exit, plan.stdout, plan.stderr)
			}
			applied := runVerb(t, verbClean, repoHome{root, home}.call("--landed", "--apply", plan.mustFingerprint(t)))
			if applied.exit != 0 || applied.stderr != "" {
				t.Fatalf("apply exit=%d stderr=%q", applied.exit, applied.stderr)
			}
			if _, err := os.Lstat(creation.Path); err != nil {
				t.Fatalf("special path disappeared: %v", err)
			}
		})
	}
	markProof(t, "landing/journey/hostile-residue")
}

func TestLandedConsumersRejectSpecialGitMetadataBeforePlanning(t *testing.T) {
	for _, tc := range []struct {
		name string
		make func(*testing.T, string)
	}{
		{name: "fifo", make: func(t *testing.T, path string) {
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
			}
		}},
		{name: "socket", make: func(t *testing.T, path string) {
			listener, err := net.Listen("unix", path)
			if err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("unix sockets unavailable: %v", err))
			}
			t.Cleanup(func() { _ = listener.Close() })
		}},
		{name: "symlink", make: func(t *testing.T, path string) {
			if err := os.Symlink(path+"-target", path); err != nil {
				capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newWorktreeRepo(t)
			home := filepath.Join(root, ".bench-home")
			creation := mustCreate(t, root, home, "landed-metadata-"+tc.name, tc.name)
			landAssignment(t, root, creation, "metadata.txt")
			metadata := filepath.Join(creation.Path, ".git")
			if err := os.Remove(metadata); err != nil {
				t.Fatal(err)
			}
			tc.make(t, metadata)

			realGit, err := exec.LookPath("git")
			if err != nil {
				t.Fatal(err)
			}
			wrapper := t.TempDir()
			log := filepath.Join(wrapper, "target-git.log")
			wrapperPath := filepath.Join(wrapper, "git")
			wrapperSource := "#!/bin/sh\nprevious=\nfor argument in \"$@\"; do\n  if [ \"$previous\" = -C ] && [ \"$argument\" = \"$BENCH_TEST_FORBIDDEN_GIT_C\" ]; then\n    printf '%s\\n' \"$@\" >> \"$BENCH_TEST_TARGET_GIT_LOG\"\n    exit 97\n  fi\n  previous=$argument\ndone\nexec \"$BENCH_TEST_REAL_GIT\" \"$@\"\n"
			if err := os.WriteFile(wrapperPath, []byte(wrapperSource), 0o755); err != nil {
				t.Fatal(err)
			}
			bindEnv(t, "BENCH_TEST_FORBIDDEN_GIT_C", creation.Path)
			bindEnv(t, "BENCH_TEST_TARGET_GIT_LOG", log)
			bindEnv(t, "BENCH_TEST_REAL_GIT", realGit)
			bindEnv(t, "PATH", wrapper+string(os.PathListSeparator)+os.Getenv("PATH"))

			if listing := runVerb(t, verbList, repoHome{root, home}.call()); listing.exit != 0 || !strings.Contains(listing.stdout, creation.Assignment.ID) {
				t.Fatalf("ListCommand = (%d, %q), want complete assignment row", listing.exit, listing.stdout)
			}
			assertNoTargetGitCalls(t, log, "ListCommand")
			if resumed := runVerb(t, verbResumeClean, repoHome{root, home}.call()); resumed.exit != 0 {
				t.Fatalf("ResumeCleanCommand = (%d, %q, %q), want completion", resumed.exit, resumed.stdout, resumed.stderr)
			}
			assertNoTargetGitCalls(t, log, "ResumeCleanCommand")
		})
	}
}

func assertNoTargetGitCalls(t *testing.T, log, consumer string) {
	t.Helper()
	if calls, err := os.ReadFile(log); err == nil && len(calls) != 0 {
		t.Fatalf("%s reached target git: %q", consumer, calls)
	} else if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
}
