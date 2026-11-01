// The rebuild remedy rows grade the freshness owner against the pool guard it must pass.
// They live here, beside the guard, because internal/freshness is at its file budget.
package benchguard_test

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/shellcommand"
)

// TestRebuildRemedyKind is RB01, RB02, RB05, RB06, and RB10. A root under the pool of
// home takes the build verb, and every other root keeps the cd form.
func TestRebuildRemedyKind(t *testing.T) {
	home := t.TempDir()
	pool := home + "/worktrees/k/x"
	for _, tc := range []struct {
		row, home, root, action string
		worktree                bool
	}{
		{"RB01 a pool root", home, pool, "bench worktree build '" + pool + "'", true},
		{"RB02 a root outside the pool", home, t.TempDir(), "", false},
		{"RB05 a sibling of the pool", home, home + "/worktrees-old/x", "", false},
		{"RB06 an empty home", "", "/worktrees/k/x", "", false},
		// The empty home's pool is the relative path worktrees, which this root is under.
		{"RB06 an empty home with a relative root", "", "worktrees/k/x", "", false},
	} {
		want := tc.action
		if want == "" {
			want = freshness.CheckoutRebuildAction(tc.root)
		}
		if got := freshness.RebuildIn(tc.home, tc.root); got.Action != want || got.Worktree != tc.worktree {
			t.Errorf("%s: RebuildIn(%q, %q) = %+v, want Action %q and Worktree %t", tc.row, tc.home, tc.root, got, want, tc.worktree)
		}
	}
	// RB10 pins the kept cd form itself, so its expectation is written out by hand.
	want := "cd '" + pool + "' && bash scripts/go-build.sh '" + pool + "' '" + pool + "/dist/bench'"
	if got := freshness.CheckoutRebuildAction(pool); got != want {
		t.Errorf("RB10: CheckoutRebuildAction(%q) = %q, want the cd form %q", pool, got, want)
	}
}

// TestRebuildRemedyPassesTheHooks is RB03 and RB04. The build remedy passes the pool guard
// and the follow-on rule, and the cd form for the same root is the one the guard refuses.
func TestRebuildRemedyPassesTheHooks(t *testing.T) {
	home := t.TempDir()
	root := home + "/worktrees/k/x"
	pools := poolkey.Pools(home)
	action := freshness.RebuildIn(home, root).Action
	if target := benchguard.PoolReference(action, pools); target != "" {
		t.Errorf("RB03: PoolReference(%q) = %q, want no pool reference", action, target)
	}
	if verdict := benchguard.Classify(action, benchguard.DefaultResolver()); verdict.Blocked {
		t.Errorf("RB03: Classify(%q) blocks: %s", action, verdict.Message())
	}
	if checkout := freshness.CheckoutRebuildAction(root); benchguard.PoolReference(checkout, pools) == "" {
		t.Errorf("RB04: PoolReference(%q) is empty, want the guard to refuse the cd form", checkout)
	}
}

// TestRebuildRemedyHostileRoot is RB07. A space, a quote, and a glob character in the
// root stay inside the one target word.
func TestRebuildRemedyHostileRoot(t *testing.T) {
	home := t.TempDir()
	root := home + "/worktrees/k/it's a *"
	action := freshness.RebuildIn(home, root).Action
	stream := shellcommand.Parse(action)
	if len(stream.Commands) != 1 {
		t.Fatalf("RB07: %q parses to %d commands, want 1", action, len(stream.Commands))
	}
	span := stream.Commands[0]
	words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
	if want := []string{"bench", "worktree", "build", root}; !slices.Equal(words, want) {
		t.Errorf("RB07: %q parses to %q, want %q", action, words, want)
	}
}

// TestRebuildActionReadsTheProcessHome is RB08. The ambient form grades the pool of the
// process Bench home, which is the pool the hook grades.
func TestRebuildActionReadsTheProcessHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	root := home + "/worktrees/k/x"
	if got, want := freshness.RebuildAction(root), "bench worktree build '"+root+"'"; got != want {
		t.Errorf("RB08: RebuildAction(%q) = %q, want %q", root, got, want)
	}
}

// TestRefusalCarriesThePoolRemedy is RB09. Each printer of a refusal for a pool root
// names the build remedy once and no cd.
func TestRefusalCarriesThePoolRemedy(t *testing.T) {
	home := t.TempDir()
	t.Setenv(benchhome.Env, home)
	root := home + "/worktrees/k/x"
	refusal := &freshness.Refusal{Executable: filepath.Join(root, "dist", "bench"), RepairRoot: root, Cause: errors.New("fixture")}
	message := refusal.Error()
	if got := strings.Count(message, "; rebuild with bench worktree build '"+root+"'"); got != 1 {
		t.Errorf("RB09: Error() = %q holds the build remedy %d times, want 1", message, got)
	}
	if strings.Contains(message, "cd ") {
		t.Errorf("RB09: Error() = %q, want no cd", message)
	}
}
