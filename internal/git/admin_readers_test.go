package git

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gittest"
)

// independentRevParse runs an independent `git rev-parse` for comparison against a
// reader's answer. It reuses runGit, the package's one exec.Command site, and
// fails the test on a nonzero exit.
func independentRevParse(t *testing.T, root string, args ...string) string {
	t.Helper()
	return strings.TrimRight(runGit(t, root, append([]string{"rev-parse"}, args...)...), "\n")
}

// bareRepo initializes a bare repository and returns its root. It reuses runGit
// rather than a second repository constructor, because the ordinary-build census
// allows the package exactly one.
func bareRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init", "-q", "--bare")
	return root
}

// linkedWorktree adds a detached linked worktree off root and returns its path.
func linkedWorktree(t *testing.T, root string) string {
	t.Helper()
	linked := filepath.Join(t.TempDir(), "linked")
	runGit(t, root, "worktree", "add", "-q", "--detach", linked, "HEAD")
	return linked
}

func TestAdminDirMatchesIndependentRevParse(t *testing.T) {
	primary := newRepo(t)
	linked := linkedWorktree(t, primary)
	bare := bareRepo(t)
	for name, root := range map[string]string{
		"primary checkout": primary,
		"linked worktree":  linked,
		"bare repository":  bare,
	} {
		t.Run(name, func(t *testing.T) {
			want := independentRevParse(t, root, "--path-format=absolute", "--git-dir")
			got, err := AdminDir(root)
			if err != nil || got != want {
				t.Fatalf("AdminDir(%s) = %q, %v, want %q, nil", root, got, err, want)
			}
		})
	}
}

func TestAdminPathMatchesIndependentRevParse(t *testing.T) {
	primary := newRepo(t)
	linked := linkedWorktree(t, primary)
	bare := bareRepo(t)
	roots := map[string]string{
		"primary checkout": primary,
		"linked worktree":  linked,
		"bare repository":  bare,
	}
	if aliased, ok := symlinkedParentAlias(t, primary); ok {
		roots["primary checkout via a symlinked parent"] = aliased
	}
	for name, root := range roots {
		t.Run(name, func(t *testing.T) {
			want := independentRevParse(t, root, "--path-format=absolute", "--git-path", "index")
			got, err := AdminPath(root, "index")
			if err != nil || got != want {
				t.Fatalf("AdminPath(%s, index) = %q, %v, want %q, nil", root, got, err, want)
			}
		})
	}
}

// symlinkedParentAlias creates a symlink to repo's parent directory and returns repo
// addressed through that alias, so a reader that keeps the alias instead of resolving
// it reds against the independent, canonicalizing rev-parse run. (Coverage row GR2.)
func symlinkedParentAlias(t *testing.T, repo string) (string, bool) {
	t.Helper()
	parent := filepath.Dir(repo)
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(parent, alias); err != nil {
		capability.Capability(t, capability.Symlink, "symlinks unavailable: "+err.Error())
		return "", false
	}
	return filepath.Join(alias, filepath.Base(repo)), true
}

func TestBareRepositoryReadersAgree(t *testing.T) {
	root := bareRepo(t)
	dir, err := AdminDir(root)
	if err != nil {
		t.Fatalf("AdminDir: %v", err)
	}
	common, err := CommonDir(root)
	if err != nil {
		t.Fatalf("CommonDir: %v", err)
	}
	if dir != common {
		t.Fatalf("AdminDir=%q, CommonDir=%q, want equal over a bare repository", dir, common)
	}
	primary, err := IsPrimaryCheckout(root)
	if err != nil || !primary {
		t.Fatalf("IsPrimaryCheckout(bare) = %v, %v, want true, nil", primary, err)
	}
}

func TestReadersRefuseBadAdministrationDirectories(t *testing.T) {
	for _, tc := range []struct {
		name, mode string
		reader     func(root string) (string, error)
		answer     func(root string) string
		setup      func(t *testing.T, root, answer string)
		want       []string
	}{
		{"AdminDir missing", "bad-git-dir", AdminDir, func(root string) string { return filepath.Join(root, "missing-admin") }, nil, []string{"missing path"}},
		{"AdminDir empty", "empty-git-dir", AdminDir, func(string) string { return "" }, nil, []string{"empty path"}},
		{"AdminDir symlink", "symlink-git-dir", AdminDir, func(root string) string { return filepath.Join(root, "symlink-admin") }, symlinkAt, []string{"symlink"}},
		{"AdminDir non-directory", "file-git-dir", AdminDir, func(root string) string { return filepath.Join(root, "file-admin") }, regularFileAt, []string{"non-directory"}},
		{"CommonDir missing", "bad-rev-parse", CommonDir, func(root string) string { return filepath.Join(root, "missing-common") }, nil, []string{"missing path"}},
		{"CommonDir empty", "empty-rev-parse", CommonDir, func(string) string { return "" }, nil, []string{"empty path"}},
		{"CommonDir symlink", "symlink-rev-parse", CommonDir, func(root string) string { return filepath.Join(root, "symlink-common") }, symlinkAt, []string{"symlink"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			gittest.StubGit(t, root, tc.mode, filepath.Join(t.TempDir(), "argv"))
			if tc.setup != nil {
				tc.setup(t, root, tc.answer(root))
			}
			_, err := tc.reader(root)
			var resolution *ResolutionError
			if !errors.As(err, &resolution) || resolution.Action != investigateGitFailureAction {
				t.Fatalf("%s: err = %v, want a typed refusal", tc.name, err)
			}
			for _, fragment := range tc.want {
				if !strings.Contains(err.Error(), fragment) {
					t.Fatalf("%s: err = %v, want %q", tc.name, err, fragment)
				}
			}
		})
	}
}

// symlinkAt creates a symlink at answer, standing in for a hostile symlinked
// administration directory or common directory. It needs a real directory
// target — root's own .git — because the stub's answer path is fixed and the
// only free variable is what the test plants there.
func symlinkAt(t *testing.T, root, answer string) {
	t.Helper()
	if answer == "" {
		t.Fatal("symlinkAt: empty answer path")
	}
	if err := os.Symlink(filepath.Join(root, ".git"), answer); err != nil {
		capability.Capability(t, capability.Symlink, "symlinks unavailable: "+err.Error())
	}
}

// regularFileAt writes a regular file at answer, standing in for a hostile
// non-directory administration directory answer.
func regularFileAt(t *testing.T, root, answer string) {
	t.Helper()
	if err := os.WriteFile(answer, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestAdminPathJoinsRelativeAnswerOntoRoot proves the file reader joins git's
// relative answer onto the absolute root, because git ran with -C root.
// (Coverage row GR5.)
func TestAdminPathJoinsRelativeAnswerOntoRoot(t *testing.T) {
	root := newRepo(t)
	gittest.StubGit(t, root, "relative-git-path", filepath.Join(t.TempDir(), "argv"))
	absolute, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(absolute, ".git", "index")
	got, err := AdminPath(root, "index")
	if err != nil || got != want {
		t.Fatalf("AdminPath(%s, index) = %q, %v, want %q, nil", root, got, err, want)
	}

	t.Run("root reaches the repository through a symlink followed by dot-dot", func(t *testing.T) {
		base := t.TempDir()
		physical := filepath.Join(base, "physical")
		child := filepath.Join(physical, "child")
		if err := os.MkdirAll(child, 0o755); err != nil {
			t.Fatal(err)
		}
		jump := filepath.Join(base, "jump")
		if err := os.Symlink(child, jump); err != nil {
			capability.Capability(t, capability.Symlink, "symlinks unavailable: "+err.Error())
		}
		// A lexical filepath.Join would clean "jump/.." away before AdminPath ever
		// sees it, so the root is built by concatenation to keep the symlinked
		// component AdminPath must resolve physically.
		root := jump + string(filepath.Separator) + ".."
		gittest.StubGit(t, root, "relative-git-path", filepath.Join(t.TempDir(), "argv"))
		want := filepath.Join(physical, ".git", "index")
		got, err := AdminPath(root, "index")
		if err != nil || got != want {
			t.Fatalf("AdminPath(%s, index) = %q, %v, want %q, nil", root, got, err, want)
		}
	})
}

// TestAdminPathKeepsSymlinkPath proves the file reader answers the symlink's
// own path, not its target, so an Lstat at a file site still sees the symlink.
// (Coverage row GR40.)
func TestAdminPathKeepsSymlinkPath(t *testing.T) {
	root := newRepo(t)
	target := filepath.Join(t.TempDir(), "outside-lease")
	if err := os.WriteFile(target, []byte("lease"), 0o600); err != nil {
		t.Fatal(err)
	}
	lease := filepath.Join(root, ".git", BenchLeaseFilename)
	if err := os.Symlink(target, lease); err != nil {
		capability.Capability(t, capability.Symlink, "symlinks unavailable: "+err.Error())
	}
	got, err := AdminPath(root, BenchLeaseFilename)
	if err != nil {
		t.Fatalf("AdminPath(%s, %s) = %v, want the symlink's own path", root, BenchLeaseFilename, err)
	}
	info, err := os.Lstat(got)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("AdminPath returned %q, whose Lstat mode is %v, want a symlink", got, info.Mode())
	}
}

// TestRelativeGitPathStubPassesThroughToRealGit proves the pass-through the
// relative-git-path stub mode owes every sibling ticket: only the file query
// it targets is stubbed, and every other invocation reaches the real git the
// stub located before the test replaced PATH.
func TestRelativeGitPathStubPassesThroughToRealGit(t *testing.T) {
	root := newRepo(t)
	gittest.StubGit(t, root, "relative-git-path", filepath.Join(t.TempDir(), "argv"))
	got := strings.TrimRight(runGit(t, root, "rev-parse", "--show-toplevel"), "\n")
	want, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("show-toplevel through stub = %q, want %q", got, want)
	}
}

func TestReadersTimeOutUnderTheWorktreeListBound(t *testing.T) {
	restore := SetWorktreeListTimeoutForTest(50 * time.Millisecond)
	t.Cleanup(restore)

	for _, tc := range []struct {
		name, mode string
		call       func(root string) error
	}{
		{"AdminDir", "block-git-dir", func(root string) error { _, err := AdminDir(root); return err }},
		{"AdminPath", "block-git-path", func(root string) error { _, err := AdminPath(root, "index"); return err }},
		{"CommonDir", "block-rev-parse", func(root string) error { _, err := CommonDir(root); return err }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			gittest.StubGit(t, root, tc.mode, filepath.Join(t.TempDir(), "argv"))
			done := make(chan error, 1)
			go func() { done <- tc.call(root) }()
			select {
			case err := <-done:
				var typed *ResolutionError
				if !errors.As(err, &typed) || !strings.Contains(err.Error(), "timed out") {
					t.Fatalf("%s timeout = %v, want a typed timeout", tc.name, err)
				}
			case <-time.After(time.Second):
				t.Fatalf("%s did not return within one second", tc.name)
			}
		})
	}
}

func TestReadersTypeStartFailures(t *testing.T) {
	root := newRepo(t)
	t.Setenv("PATH", t.TempDir())
	for _, tc := range []struct {
		name string
		call func() error
	}{
		{"AdminDir", func() error { _, err := AdminDir(root); return err }},
		{"AdminPath", func() error { _, err := AdminPath(root, "index"); return err }},
		{"CommonDir", func() error { _, err := CommonDir(root); return err }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			var typed *ResolutionError
			if !errors.As(err, &typed) || !strings.Contains(err.Error(), "rev-parse") || !strings.Contains(err.Error(), "executable file not found") {
				t.Fatalf("%s start failure = %v, want rev-parse + executable file not found", tc.name, err)
			}
		})
	}
}

func TestResolutionErrorNamesItsSubject(t *testing.T) {
	for _, tc := range []struct {
		name, subject string
	}{
		{"common directory", subjectCommonDir},
		{"admin directory", subjectAdminDir},
		{"admin path", subjectAdminPath},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := &ResolutionError{Err: errors.New("boom"), Action: investigateGitFailureAction, Subject: tc.subject}
			if !strings.HasPrefix(err.Error(), tc.subject) {
				t.Fatalf("Error() = %q, want prefix %q", err.Error(), tc.subject)
			}
		})
	}
}
