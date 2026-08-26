package census

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/poolkey"
)

// TestDirSitsBesideThePool proves the records live under <home>/census and never
// under the worktree pool, which `bench worktree reclaim` enumerates.
// (Coverage row EC26.)
func TestDirSitsBesideThePool(t *testing.T) {
	t.Parallel()
	home := filepath.Join(string(filepath.Separator), "home-a", ".bench")
	root := filepath.Join(string(filepath.Separator), "repos", "example")
	got := Dir(home, root)
	if want := filepath.Join(home, "census", poolkey.Key(root)); got != want {
		t.Fatalf("Dir = %q, want %q", got, want)
	}
	if !strings.HasPrefix(got, filepath.Join(home, "census")+string(filepath.Separator)) {
		t.Fatalf("Dir %q is not under the census home", got)
	}
	if strings.Contains(got, filepath.Join(home, "worktrees")) {
		t.Fatalf("Dir %q is inside the worktree pool", got)
	}
}

// TestDirDependsOnTheInjectedHome proves the home arrives explicitly, with no
// environment read below the effect boundary.
func TestDirDependsOnTheInjectedHome(t *testing.T) {
	t.Parallel()
	root := filepath.Join(string(filepath.Separator), "repos", "example")
	a := Dir(filepath.Join(string(filepath.Separator), "home-a"), root)
	b := Dir(filepath.Join(string(filepath.Separator), "home-b"), root)
	if filepath.Base(a) != filepath.Base(b) {
		t.Fatalf("the census key must depend only on the root: %s vs %s", a, b)
	}
	if a == b {
		t.Fatalf("Dir ignored the injected home: %s", a)
	}
}
