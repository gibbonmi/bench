package diff

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestCommandRefusesEachSnapshotIdentityDimensionAfterOneRetry(t *testing.T) {
	tests := []struct {
		name      string
		dimension string
		prepare   func(*testing.T, string, string, string) func(int)
	}{
		{
			name: "HEAD", dimension: "HEAD",
			prepare: func(t *testing.T, _, _, _ string) func(int) {
				return func(int) { runGit(t, "commit", "-q", "--allow-empty", "-m", "move HEAD") }
			},
		},
		{
			name: "default tip", dimension: "default tip",
			prepare: func(t *testing.T, _, c1, c3 string) func(int) {
				return func(call int) {
					if call%2 == 1 {
						runGit(t, "update-ref", "refs/heads/main", c1)
					} else {
						runGit(t, "update-ref", "refs/heads/main", c3)
					}
				}
			},
		},
		{
			name: "recorded base", dimension: "recorded base",
			prepare: func(t *testing.T, c0, c1, _ string) func(int) {
				runGit(t, "config", "branch.feature.benchBase", c0)
				return func(call int) {
					if call%2 == 1 {
						runGit(t, "config", "branch.feature.benchBase", c1)
					} else {
						runGit(t, "config", "branch.feature.benchBase", c0)
					}
				}
			},
		},
		{
			name: "raw index", dimension: "index",
			prepare: func(t *testing.T, _, _, _ string) func(int) {
				mustWriteFile(t, "index.txt", "index\n")
				runGit(t, "add", "index.txt")
				return func(call int) {
					flag := "+x"
					if call%2 == 0 {
						flag = "-x"
					}
					runGit(t, "update-index", "--chmod="+flag, "index.txt")
				}
			},
		},
		{
			name: "raw porcelain", dimension: "porcelain",
			prepare: func(t *testing.T, _, _, _ string) func(int) {
				mustWriteFile(t, "untracked-a.txt", "new\n")
				return func(call int) {
					from, to := "untracked-a.txt", "untracked-b.txt"
					if call%2 == 0 {
						from, to = to, from
					}
					if err := os.Rename(from, to); err != nil {
						t.Fatal(err)
					}
				}
			},
		},
		{
			name: "dirty content", dimension: "dirty content",
			prepare: func(t *testing.T, _, _, _ string) func(int) {
				mustWriteFile(t, "f.txt", "dirty 0\n")
				return func(call int) { mustWriteFile(t, "f.txt", "dirty "+string(rune('0'+call))+"\n") }
			},
		},
		{
			name: "dirty mode", dimension: "dirty mode",
			prepare: func(t *testing.T, _, _, _ string) func(int) {
				mustWriteFile(t, "f.txt", "dirty\n")
				return func(call int) {
					mode := os.FileMode(0o755)
					if call%2 == 0 {
						mode = 0o644
					}
					if err := os.Chmod("f.txt", mode); err != nil {
						t.Fatal(err)
					}
				}
			},
		},
		{
			name: "symlink target", dimension: "dirty content",
			prepare: func(t *testing.T, _, _, _ string) func(int) {
				mustWriteFile(t, "target-a", "a\n")
				if err := os.Symlink("target-a", "link"); err != nil {
					t.Fatal(err)
				}
				return func(call int) {
					if err := os.Remove("link"); err != nil {
						t.Fatal(err)
					}
					target := "target-a"
					if call%2 == 1 {
						target = "target-b"
						mustWriteFile(t, target, "b\n")
					}
					if err := os.Symlink(target, "link"); err != nil {
						t.Fatal(err)
					}
				}
			},
		},
		{
			name: "non-regular stat", dimension: "dirty content",
			prepare: func(t *testing.T, _, _, _ string) func(int) {
				if err := syscall.Mkfifo("pipe", 0o600); err != nil {
					t.Fatal(err)
				}
				return func(call int) {
					stamp := time.Unix(int64(call), 0)
					if err := os.Chtimes("pipe", stamp, stamp); err != nil {
						t.Fatal(err)
					}
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, c0, c1, c3 := seedDivergedRepo(t)
			mutate := tc.prepare(t, c0, c1, c3)
			assertRepeatedDrift(t, tc.dimension, mutate)
		})
	}

	t.Run("dirty gitlink HEAD", func(t *testing.T) {
		_, _, _, _ = seedDivergedRepo(t)
		subject := t.TempDir()
		runGitAt(t, subject, "init", "-q", "-b", "main")
		runGitAt(t, subject, "config", "user.email", "t@example.com")
		runGitAt(t, subject, "config", "user.name", "t")
		mustWriteFileAt(t, subject, "sub.txt", "zero\n")
		runGitAt(t, subject, "add", "sub.txt")
		runGitAt(t, subject, "commit", "-q", "-m", "sub zero")
		zero := runGitAt(t, subject, "rev-parse", "HEAD")
		mustWriteFileAt(t, subject, "sub.txt", "one\n")
		runGitAt(t, subject, "commit", "-qam", "sub one")
		one := runGitAt(t, subject, "rev-parse", "HEAD")
		mustWriteFileAt(t, subject, "sub.txt", "two\n")
		runGitAt(t, subject, "commit", "-qam", "sub two")
		runGit(t, "-c", "protocol.file.allow=always", "submodule", "add", "-q", subject, "sub")
		runGit(t, "commit", "-q", "-m", "add submodule")
		runGitAt(t, "sub", "checkout", "-q", zero)
		assertRepeatedDrift(t, "dirty gitlink HEAD", func(call int) {
			target := one
			if call%2 == 0 {
				target = zero
			}
			runGitAt(t, "sub", "checkout", "-q", target)
		})
	})
}

func TestCommandResolvesRecordedBaseOncePerAttempt(t *testing.T) {
	root, c0, _, _ := seedDivergedRepo(t)
	runGit(t, "config", "branch.feature.benchBase", c0)
	mustWriteFile(t, "f.txt", "dirty\n")
	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	logPath := filepath.Join(bin, "git.log")
	wrapper := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> " + shellQuote(logPath) + "\nexec " + shellQuote(realGit) + " \"$@\"\n"
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte(wrapper), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	previous := snapshotAfterRead
	defer func() { snapshotAfterRead = previous }()
	snapshotAfterRead = func() { mustWriteFile(t, filepath.Join(root, "f.txt"), "move\n") }
	_, _ = Command(nil)
	raw, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(raw), "config branch.feature.benchBase\n"); got != 4 {
		t.Fatalf("recorded-base reads = %d, want one resolution plus one after-read capture per each of two attempts; trace:\n%s", got, raw)
	}
	if got := strings.Count(string(raw), "symbolic-ref --quiet --short refs/remotes/origin/HEAD\n"); got != 2 {
		t.Fatalf("default resolutions = %d, want one fresh resolution per attempt; trace:\n%s", got, raw)
	}
	if got := strings.Count(string(raw), "merge-base --is-ancestor"); got != 2 {
		t.Fatalf("recorded-base semantic resolutions = %d, want one fresh resolution per attempt; trace:\n%s", got, raw)
	}
}

// TestIndexIdentityRefusesUnresolvedAnswer is GR26: indexIdentity must call the
// named reader and surface its typed failure rather than swallow an
// unresolved answer. This test sets PATH and so stays serial.
func TestIndexIdentityRefusesUnresolvedAnswer(t *testing.T) {
	root, _, _, _ := seedDivergedRepo(t)
	gittest.StubGit(t, root, "fail-git-path", filepath.Join(t.TempDir(), "argv"))
	_, err := indexIdentity(root)
	var resolutionErr *git.ResolutionError
	if !errors.As(err, &resolutionErr) {
		t.Fatalf("indexIdentity(%s) error = %v, want a *git.ResolutionError", root, err)
	}
}

func assertRepeatedDrift(t *testing.T, dimension string, mutate func(int)) {
	t.Helper()
	previous := snapshotAfterRead
	defer func() { snapshotAfterRead = previous }()
	calls := 0
	snapshotAfterRead = func() {
		calls++
		mutate(calls)
	}
	out, code := Command(nil)
	if code != 1 || calls != 2 {
		t.Fatalf("Command drift = (exit %d, mutation calls %d), want (1, 2); output:\n%s", code, calls, out)
	}
	for _, want := range []string{
		"error: snapshot drift",
		"the " + dimension + " changed while reading",
		"help[1]{cmd,why}:\n  bench diff,retry after the repository stopped moving\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("drift response missing %q:\n%s", want, out)
		}
	}
	for _, forbidden := range []string{"revision[", "aggregate[", "files[", "checkout[", "whitespace["} {
		if strings.Contains(out, forbidden) {
			t.Errorf("drift response leaked a partial snapshot block %q:\n%s", forbidden, out)
		}
	}
}

func runGitAt(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git -C %s %v: %v\n%s", root, args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func mustWriteFileAt(t *testing.T, root, path, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, path), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
