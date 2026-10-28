package gate

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// BG117: a detached grade of the checkout writes its state during the guard window of a
// foreground gate, and the guard reports no change. Another Bench admin file still shows.
func TestCheckoutGuardIgnoresBackgroundGradeState(t *testing.T) {
	t.Parallel()
	for i, name := range backgroundstate.GradePaths("", "") {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root := gittest.RepoOnBranch(t, "main")
			admin, common := guardedDirs(t, root)
			ctx := withCheckoutGuard(context.Background(), root)
			path := backgroundstate.GradePaths(admin, common)[i]
			claim, err := backgroundstate.ClaimFor("key", root)
			if err != nil {
				t.Fatal(err)
			}
			// The claim directory holds the claims; each other path is one state file.
			if path == filepath.Dir(backgroundstate.ClaimPath(common, claim.Key)) {
				err = backgroundstate.WriteClaim(durablefile.Native(), common, claim)
			} else {
				err = os.WriteFile(path, []byte("state\n"), 0o600)
			}
			if err != nil {
				t.Fatal(err)
			}
			if rows, _, failed := checkoutReport(ctx)(); failed || len(rows) != 0 {
				t.Fatalf("checkout guard reported background grade state %s: %q", path, rows)
			}
		})
	}
	t.Run("bench-other", func(t *testing.T) {
		t.Parallel()
		root := gittest.RepoOnBranch(t, "main")
		admin, _ := guardedDirs(t, root)
		ctx := withCheckoutGuard(context.Background(), root)
		if err := os.WriteFile(filepath.Join(admin, "bench-other"), []byte("state\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		rows, _, failed := checkoutReport(ctx)()
		if want := `checkout guard: changed path ".git/bench-other"`; !failed || len(rows) != 1 || rows[0] != want {
			t.Fatalf("checkout guard rows = %q, want [%s]", rows, want)
		}
	})
}

// guardedDirs answers the admin directory and the common directory of the checkout root.
func guardedDirs(t *testing.T, root string) (string, string) {
	t.Helper()
	admin, err := benchgit.AdminDir(root)
	if err != nil {
		t.Fatal(err)
	}
	common, err := benchgit.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	return admin, common
}
