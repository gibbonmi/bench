package worktree

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
)

// The fixtures here exercise git.LandedInDefault at the seam every cleanup path asks:
// its verdict is the only thing standing between a branch and deletion, so these tests
// pin both the squash-landing it must prove and every ambiguity it must refuse.

// onNewBranch cuts name from main, runs build inside it, and returns to main.
func onNewBranch(t *testing.T, root, name string, build func()) {
	t.Helper()
	gitRun(t, root, "switch", "-q", "-c", name, "main")
	build()
	gitRun(t, root, "switch", "-q", "main")
}

func commitAll(t *testing.T, root, message string) {
	t.Helper()
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-qm", message)
}

// squashLand composes every commit on branch beyond main into one commit on main —
// the landing shape ancestry, merge detection, and git cherry all miss.
func squashLand(t *testing.T, root, branch string) {
	t.Helper()
	gitRun(t, root, "cherry-pick", "--no-commit", "main.."+branch)
	gitRun(t, root, "commit", "-qm", "squash "+branch)
}

func verdict(t *testing.T, root, branch string) (bool, bool) {
	t.Helper()
	landed, byContent, err := git.LandedInDefault(root, branch, "main")
	mustNoError(t, err)
	return landed, byContent
}

// [RL1]
func TestLandedInDefaultProvesSquashLanding(t *testing.T) {
	t.Run("multi-form squash is proven landed", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "feature", func() {
			mustWrite(t, filepath.Join(root, "feature.txt"), []byte("feature\n"), 0o644)
			mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("modified\n"), 0o644)
			mustWrite(t, filepath.Join(root, "bin.dat"), []byte{0x00, 0x01, 0xff, 0xfe, 0x00}, 0o644)
			mustWrite(t, filepath.Join(root, "crlf.txt"), []byte("one\r\ntwo\r\n"), 0o644)
			gitRun(t, root, "mv", "README.md", "docs.md")
			commitAll(t, root, "content forms")
			if err := os.Symlink("tracked.txt", filepath.Join(root, "link")); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(filepath.Join(root, "feature.txt"), 0o755); err != nil {
				t.Fatal(err)
			}
			commitAll(t, root, "symlink and mode")
		})
		squashLand(t, root, "feature")
		landed, byContent := verdict(t, root, "feature")
		requireTest(t, landed, "squash-landed branch reported not landed")
		requireTest(t, byContent, "squash proof did not report as a content proof")
	})
	t.Run("ancestry still proves", func(t *testing.T) {
		root := newWorktreeRepo(t)
		gitRun(t, root, "branch", "behind", "main")
		mustWrite(t, filepath.Join(root, "ahead.txt"), []byte("ahead\n"), 0o644)
		commitAll(t, root, "advance main")
		landed, byContent := verdict(t, root, "behind")
		requireTest(t, landed && !byContent, "ancestor branch verdict = %t/%t, want landed by ancestry", landed, byContent)
	})
	t.Run("patch equivalence still proves", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "picked", func() {
			mustWrite(t, filepath.Join(root, "picked.txt"), []byte("picked\n"), 0o644)
			commitAll(t, root, "picked")
		})
		gitRun(t, root, "commit", "--allow-empty", "-qm", "diverge")
		gitRun(t, root, "cherry-pick", "picked")
		landed, byContent := verdict(t, root, "picked")
		requireTest(t, landed && byContent, "cherry-picked branch verdict = %t/%t, want landed by patch", landed, byContent)
	})
	t.Run("merge-carrying branch is still kept even with its content landed", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "side", func() {
			mustWrite(t, filepath.Join(root, "side.txt"), []byte("side\n"), 0o644)
			commitAll(t, root, "side")
		})
		onNewBranch(t, root, "merged", func() {
			mustWrite(t, filepath.Join(root, "merged.txt"), []byte("merged\n"), 0o644)
			commitAll(t, root, "merged")
			gitRun(t, root, "merge", "-q", "--no-ff", "-m", "merge side", "side")
		})
		// Land every byte of the merge-carrying branch on main, as one composed commit.
		mustWrite(t, filepath.Join(root, "side.txt"), []byte("side\n"), 0o644)
		mustWrite(t, filepath.Join(root, "merged.txt"), []byte("merged\n"), 0o644)
		commitAll(t, root, "compose both")
		landed, _ := verdict(t, root, "merged")
		requireTest(t, !landed, "merge-carrying branch was proven landed; merge-only content must stay kept")
	})
}

// [RL2]
func TestLandedInDefaultRefusesUnlandedContent(t *testing.T) {
	t.Run("strict superset of what landed", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "feature", func() {
			mustWrite(t, filepath.Join(root, "one.txt"), []byte("one\n"), 0o644)
			commitAll(t, root, "one")
			mustWrite(t, filepath.Join(root, "two.txt"), []byte("two\n"), 0o644)
			commitAll(t, root, "two")
		})
		mustWrite(t, filepath.Join(root, "one.txt"), []byte("one\n"), 0o644)
		commitAll(t, root, "land only one")
		landed, _ := verdict(t, root, "feature")
		requireTest(t, !landed, "branch with an unlanded commit was proven landed")
	})
	t.Run("same path added on both sides with different content", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "feature", func() {
			mustWrite(t, filepath.Join(root, "shared.txt"), []byte("feature\n"), 0o644)
			commitAll(t, root, "feature shared")
		})
		mustWrite(t, filepath.Join(root, "shared.txt"), []byte("mainline\n"), 0o644)
		commitAll(t, root, "mainline shared")
		landed, _ := verdict(t, root, "feature")
		requireTest(t, !landed, "diverging same-path addition was proven landed")
	})
	t.Run("binary divergence", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "feature", func() {
			mustWrite(t, filepath.Join(root, "bin.dat"), []byte{0x00, 0x01, 0x02}, 0o644)
			commitAll(t, root, "feature binary")
		})
		mustWrite(t, filepath.Join(root, "bin.dat"), []byte{0x00, 0xff, 0x02}, 0o644)
		commitAll(t, root, "mainline binary")
		landed, _ := verdict(t, root, "feature")
		requireTest(t, !landed, "diverging binary content was proven landed")
	})
	t.Run("mode-only divergence", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "feature", func() {
			if err := os.Chmod(filepath.Join(root, "tracked.txt"), 0o755); err != nil {
				t.Fatal(err)
			}
			commitAll(t, root, "make executable")
		})
		gitRun(t, root, "commit", "--allow-empty", "-qm", "diverge")
		landed, _ := verdict(t, root, "feature")
		requireTest(t, !landed, "unlanded mode change was proven landed")
	})
	t.Run("symlink versus regular file", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "feature", func() {
			if err := os.Remove(filepath.Join(root, "README.md")); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("tracked.txt", filepath.Join(root, "README.md")); err != nil {
				t.Fatal(err)
			}
			commitAll(t, root, "readme becomes symlink")
		})
		gitRun(t, root, "commit", "--allow-empty", "-qm", "diverge")
		landed, _ := verdict(t, root, "feature")
		requireTest(t, !landed, "unlanded symlink conversion was proven landed")
	})
	t.Run("symlink added where a same-content file landed", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "feature", func() {
			if err := os.Symlink("tracked.txt", filepath.Join(root, "link")); err != nil {
				t.Fatal(err)
			}
			commitAll(t, root, "symlink")
			mustWrite(t, filepath.Join(root, "other.txt"), []byte("other\n"), 0o644)
			commitAll(t, root, "other")
		})
		mustWrite(t, filepath.Join(root, "link"), []byte("tracked.txt"), 0o644)
		mustWrite(t, filepath.Join(root, "other.txt"), []byte("other\n"), 0o644)
		commitAll(t, root, "compose with a regular file")
		landed, _ := verdict(t, root, "feature")
		requireTest(t, !landed, "a symlink was proven landed by a same-content regular file")
	})
	// The squash shape keeps this case away from git cherry, whose patch-id is
	// whitespace-blind and would prove a lone CRLF commit landed by its LF counterpart;
	// the content proof itself must stay byte-exact.
	t.Run("line-ending divergence inside a squash", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "feature", func() {
			mustWrite(t, filepath.Join(root, "eol.txt"), []byte("one\r\ntwo\r\n"), 0o644)
			commitAll(t, root, "crlf")
			mustWrite(t, filepath.Join(root, "other.txt"), []byte("other\n"), 0o644)
			commitAll(t, root, "other")
		})
		mustWrite(t, filepath.Join(root, "eol.txt"), []byte("one\ntwo\n"), 0o644)
		mustWrite(t, filepath.Join(root, "other.txt"), []byte("other\n"), 0o644)
		commitAll(t, root, "compose with lf")
		landed, _ := verdict(t, root, "feature")
		requireTest(t, !landed, "CRLF content was proven landed by its LF counterpart")
	})
}

// [RL3]
func TestLandedInDefaultReportsAmbiguityAsNotLanded(t *testing.T) {
	t.Run("no merge base", func(t *testing.T) {
		root := newWorktreeRepo(t)
		gitRun(t, root, "switch", "-q", "--orphan", "lone")
		mustWrite(t, filepath.Join(root, "lone.txt"), []byte("lone\n"), 0o644)
		commitAll(t, root, "unrelated history")
		gitRun(t, root, "switch", "-q", "main")
		landed, _ := verdict(t, root, "lone")
		requireTest(t, !landed, "branch with no merge base was proven landed")
	})
	t.Run("apply that does not cleanly reverse", func(t *testing.T) {
		root := newWorktreeRepo(t)
		onNewBranch(t, root, "feature", func() {
			mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("feature\n"), 0o644)
			commitAll(t, root, "feature edit")
		})
		mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("mainline\n"), 0o644)
		commitAll(t, root, "mainline edit")
		landed, _ := verdict(t, root, "feature")
		requireTest(t, !landed, "conflicting edit was proven landed")
	})
	t.Run("submodule pointer is refused even when present", func(t *testing.T) {
		root := newWorktreeRepo(t)
		gitlink := gitOutput(t, root, "rev-parse", "HEAD")
		addGitlink := func() {
			mustMkdirAll(t, filepath.Join(root, "vendor", "dep"), 0o755)
			gitRun(t, root, "update-index", "--add", "--cacheinfo", "160000,"+gitlink+",vendor/dep")
		}
		onNewBranch(t, root, "feature", func() {
			addGitlink()
			gitRun(t, root, "commit", "-qm", "gitlink")
			mustWrite(t, filepath.Join(root, "extra.txt"), []byte("extra\n"), 0o644)
			commitAll(t, root, "extra")
		})
		addGitlink()
		mustWrite(t, filepath.Join(root, "extra.txt"), []byte("extra\n"), 0o644)
		gitRun(t, root, "add", "extra.txt")
		gitRun(t, root, "commit", "-qm", "compose gitlink and extra")
		landed, _ := verdict(t, root, "feature")
		requireTest(t, !landed, "diff carrying a submodule pointer was proven landed")
	})
}

// [RL4]
func TestLandedInDefaultIsReadOnlyAndDeterministic(t *testing.T) {
	root := newWorktreeRepo(t)
	onNewBranch(t, root, "feature", func() {
		mustWrite(t, filepath.Join(root, "feature.txt"), []byte("feature\n"), 0o644)
		commitAll(t, root, "feature")
		mustWrite(t, filepath.Join(root, "more.txt"), []byte("more\n"), 0o644)
		commitAll(t, root, "more")
	})
	squashLand(t, root, "feature")

	status := gitOutput(t, root, "status", "--porcelain")
	index, err := os.ReadFile(filepath.Join(root, ".git", "index"))
	mustNoError(t, err)

	first, _ := verdict(t, root, "feature")
	second, _ := verdict(t, root, "feature")
	requireTest(t, first && second && first == second, "verdict was not stable: %t then %t", first, second)

	indexAfter, err := os.ReadFile(filepath.Join(root, ".git", "index"))
	mustNoError(t, err)
	requireTest(t, string(index) == string(indexAfter), "the proof rewrote the real index")
	requireTest(t, gitOutput(t, root, "status", "--porcelain") == status, "the proof changed working-tree status")

	// The proof reads committed trees only, so an uncommitted edit neither changes the
	// verdict nor is touched by it.
	mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("dirty\n"), 0o644)
	dirty, _ := verdict(t, root, "feature")
	requireTest(t, dirty, "an uncommitted edit changed the verdict")
	content, err := os.ReadFile(filepath.Join(root, "tracked.txt"))
	mustNoError(t, err)
	requireTest(t, string(content) == "dirty\n", "the proof rewrote a working-tree file: %q", content)
}
