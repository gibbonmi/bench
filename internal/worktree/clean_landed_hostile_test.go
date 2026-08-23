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
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, "home space *"))
	clean := mustCreate(t, root, "landed-hostile-clean", "clean")
	dirty := mustCreate(t, root, "landed-hostile-dirty", "dirty")
	landAssignment(t, root, clean, "clean.txt")
	landAssignment(t, root, dirty, "dirty.txt")
	mustWrite(t, filepath.Join(dirty.Path, "dirty.txt"), []byte("changed\n"), 0o644)

	stdout, stderr, code := runCleanLanded(t, root, "--landed")
	if code != 0 || stderr != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, clean.Path+",remove,") || !strings.Contains(stdout, "bench worktree clean '"+dirty.Path+"'") || !strings.Contains(stdout, "bench worktree clean --landed --apply ") {
		t.Fatalf("output=%q, want safe row and pasteable help", stdout)
	}
	_, applyErr, applyCode := runCleanLanded(t, root, "--landed", "--apply", landedRowFingerprint(t, stdout))
	if applyCode != 0 || applyErr != "" {
		t.Fatalf("apply exit=%d stderr=%q", applyCode, applyErr)
	}
	if _, err := os.Lstat(clean.Path); !os.IsNotExist(err) {
		t.Fatalf("safe hostile path was not removed: %v", err)
	}
	if _, err := os.Lstat(dirty.Path); err != nil {
		t.Fatalf("dirty hostile path disappeared: %v", err)
	}
}

func TestCleanLandedControlBytePathRetained(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, "home\x1bunsafe"))
	creation := mustCreate(t, root, "landed-control", "control")
	landAssignment(t, root, creation, "control.txt")

	stdout, stderr, code := runCleanLanded(t, root, "--landed")
	pointer := "bench worktree exec " + creation.Assignment.ID + " -- bench worktree clean ."
	if code != 0 || stderr != "" || strings.ContainsRune(stdout, '\x1b') || !strings.Contains(stdout, "sha256:") || !strings.Contains(stdout, ",retain,") || !strings.Contains(stdout, "unsafe control bytes") || strings.Count(stdout, pointer) != 1 {
		t.Fatalf("plan exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	fingerprint := landedRowFingerprint(t, stdout)
	applied, applyErr, applyCode := runCleanLanded(t, root, "--landed", "--apply", fingerprint)
	if applyCode != 0 || applyErr != "" || strings.ContainsRune(applied, '\x1b') {
		t.Fatalf("apply exit=%d stdout=%q stderr=%q", applyCode, applied, applyErr)
	}
	if _, err := os.Lstat(creation.Path); err != nil {
		t.Fatalf("retained control-byte path disappeared: %v", err)
	}
}

func TestCleanLandedTabPathRendersOneRow(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(root, "home\tpath"))
	creation := mustCreate(t, root, "landed-tab", "tab")
	landAssignment(t, root, creation, "tab.txt")

	stdout, stderr, code := runCleanLanded(t, root, "--landed")
	if code != 0 || stderr != "" || !strings.HasPrefix(stdout, "worktree_cleanup[1]") || strings.ContainsRune(stdout, '\t') || !strings.Contains(stdout, `\t`) {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestCleanLandedSpecialPathsRetainedWithoutOpening(t *testing.T) {
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
			previousPlanner := planLandedExplicitWithOptions
			plannerCalls := 0
			planLandedExplicitWithOptions = func(root, path string, options CleanupOptions) (CleanupPlan, error) {
				plannerCalls++
				return previousPlanner(root, path, options)
			}
			t.Cleanup(func() { planLandedExplicitWithOptions = previousPlanner })
			root := newWorktreeRepo(t)
			bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
			creation := mustCreate(t, root, "landed-special-"+strings.ReplaceAll(tc.name, " ", "-"), tc.name)
			landAssignment(t, root, creation, "special.txt")
			if err := os.RemoveAll(creation.Path); err != nil {
				t.Fatal(err)
			}
			tc.make(t, creation.Path)

			stdout, stderr, code := runCleanLanded(t, root, "--landed")
			if code != 0 || stderr != "" || !strings.Contains(stdout, ",retain,") || !strings.Contains(stdout, "assignment path shape is "+tc.shape) {
				t.Fatalf("plan exit=%d stdout=%q stderr=%q", code, stdout, stderr)
			}
			fingerprint := landedRowFingerprint(t, stdout)
			_, applyErr, applyCode := runCleanLanded(t, root, "--landed", "--apply", fingerprint)
			if applyCode != 0 || applyErr != "" {
				t.Fatalf("apply exit=%d stderr=%q", applyCode, applyErr)
			}
			if plannerCalls != 0 {
				t.Fatalf("special path reached explicit planner %d time(s)", plannerCalls)
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
			bindEnv(t, "BENCH_HOME", filepath.Join(root, ".bench-home"))
			creation := mustCreate(t, root, "landed-metadata-"+tc.name, tc.name)
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
			chdir(t, root)

			if listing, code := ListCommand(nil); code != 0 || !strings.Contains(listing, creation.Assignment.ID) {
				t.Fatalf("ListCommand = (%d, %q), want complete assignment row", code, listing)
			}
			assertNoTargetGitCalls(t, log, "ListCommand")
			var stdout, stderr strings.Builder
			if code := ResumeCleanCommand(nil, &stdout, &stderr); code != 0 {
				t.Fatalf("ResumeCleanCommand = (%d, %q, %q), want completion", code, stdout.String(), stderr.String())
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
