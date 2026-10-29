package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// kitRootRepo answers a committed primary checkout and one linked worktree of it. The
// linked worktree sits under a parent directory whose name holds a space, so each path
// reaches git as one argument.
func kitRootRepo(t *testing.T) (primary, linked string) {
	t.Helper()
	primary = gittest.RepoOnBranch(t, "main")
	gittest.Commit(t, primary, "seed", "--allow-empty")
	linked = filepath.Join(t.TempDir(), "with space", "linked")
	gittest.MustRun(t, primary, "worktree", "add", "-q", "--detach", linked)
	return primary, linked
}

func requireKitForRoot(t *testing.T, kit, root, want string) {
	t.Helper()
	if got := KitForRoot(kit, root); got != want {
		t.Fatalf("KitForRoot(%q, %q) = %q, want %q", kit, root, got, want)
	}
}

// TestKitForRootNamesALinkedWorktreeOfThePrimaryKit pins the match. A form that keeps the
// kit answers the primary checkout.
func TestKitForRootNamesALinkedWorktreeOfThePrimaryKit(t *testing.T) {
	t.Parallel()
	primary, linked := kitRootRepo(t)
	requireKitForRoot(t, primary, linked, linked)
}

// TestKitForRootNamesALinkedWorktreeOfASiblingKit pins the common-directory match. A form
// that matches only the primary checkout keeps the sibling kit.
func TestKitForRootNamesALinkedWorktreeOfASiblingKit(t *testing.T) {
	t.Parallel()
	primary, linked := kitRootRepo(t)
	sibling := filepath.Join(t.TempDir(), "sibling")
	gittest.MustRun(t, primary, "worktree", "add", "-q", "--detach", sibling)
	requireKitForRoot(t, sibling, linked, linked)
}

// TestKitForRootKeepsAnAdoptedSubdirectoryKit pins the top-level condition. A form without
// it answers the worktree. The kit directory exists, because git fails on a missing one
// and that failure hides the condition.
func TestKitForRootKeepsAnAdoptedSubdirectoryKit(t *testing.T) {
	t.Parallel()
	primary, linked := kitRootRepo(t)
	adopted := filepath.Join(primary, ".bench")
	if err := os.Mkdir(adopted, 0o755); err != nil {
		t.Fatal(err)
	}
	requireKitForRoot(t, adopted, linked, adopted)
}

// TestKitForRootKeepsAKitFromAnotherRepository pins the common-directory condition. A form
// that skips it answers the root.
func TestKitForRootKeepsAKitFromAnotherRepository(t *testing.T) {
	t.Parallel()
	_, linked := kitRootRepo(t)
	other := gittest.RepoOnBranch(t, "main")
	requireKitForRoot(t, other, linked, other)
}

// TestKitForRootKeepsASymlinkedKitFromAnotherRepository pins the kit spelling. A form that
// answers the canonical kit loses the link.
func TestKitForRootKeepsASymlinkedKitFromAnotherRepository(t *testing.T) {
	t.Parallel()
	_, linked := kitRootRepo(t)
	link := filepath.Join(t.TempDir(), "other-link")
	if err := os.Symlink(gittest.RepoOnBranch(t, "main"), link); err != nil {
		t.Fatal(err)
	}
	requireKitForRoot(t, link, linked, link)
}

// TestKitForRootKeepsTheKitOutsideARepository pins the git-failure fallback. A form that
// treats a git failure as a match answers the root.
func TestKitForRootKeepsTheKitOutsideARepository(t *testing.T) {
	t.Parallel()
	primary, _ := kitRootRepo(t)
	requireKitForRoot(t, primary, t.TempDir(), primary)
}

// TestKitForRootKeepsAnEmptyKit pins the empty-kit guard. A form without it resolves the
// empty kit to the working directory, the primary checkout here, and answers the root.
// t.Chdir refuses a parallel test.
func TestKitForRootKeepsAnEmptyKit(t *testing.T) {
	primary, linked := kitRootRepo(t)
	t.Chdir(primary)
	requireKitForRoot(t, "", linked, "")
}

// TestKitForRootMatchesASymlinkedKitSpelling pins the canonical compare. A raw string
// comparison of the top level keeps the kit.
func TestKitForRootMatchesASymlinkedKitSpelling(t *testing.T) {
	t.Parallel()
	primary, linked := kitRootRepo(t)
	link := filepath.Join(t.TempDir(), "kit-link")
	if err := os.Symlink(primary, link); err != nil {
		t.Fatal(err)
	}
	requireKitForRoot(t, link, linked, linked)
}

// TestKitForRootKeepsASymlinkedRootSpelling pins the root spelling. A form that answers the
// canonical root loses the link.
func TestKitForRootKeepsASymlinkedRootSpelling(t *testing.T) {
	t.Parallel()
	primary, linked := kitRootRepo(t)
	link := filepath.Join(t.TempDir(), "root-link")
	if err := os.Symlink(linked, link); err != nil {
		t.Fatal(err)
	}
	requireKitForRoot(t, primary, link, link)
}

// TestKitForRootKeepsAnAbsentKit pins the absent-kit fallback. A form that treats an
// unreadable kit as a match answers the root.
func TestKitForRootKeepsAnAbsentKit(t *testing.T) {
	t.Parallel()
	_, linked := kitRootRepo(t)
	absent := filepath.Join(t.TempDir(), "absent")
	requireKitForRoot(t, absent, linked, absent)
}
