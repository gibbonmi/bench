package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/usage"
)

// TestCleanSetDiscardModifiers holds the explicit-set path to the discard modifiers one
// call carries. Each case reads its modifier back out of the rendered apply command, then
// applies the set and checks the durable effect only that modifier authorizes.
func TestCleanSetDiscardModifiers(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		modifier string
		fixture  func(*testing.T, string, string) Creation
		effect   func(*testing.T, string, Creation)
	}{
		{
			name:     "discard branch",
			modifier: "--discard-branch",
			fixture: func(t *testing.T, root, home string) Creation {
				creation := mustCreate(t, root, home, "set-discard-branch", "discard branch")
				landAssignment(t, root, creation, "branch.txt")
				return creation
			},
			effect: func(t *testing.T, root string, creation Creation) {
				if git.OK("-C", root, "show-ref", "--verify", "--quiet", creation.Assignment.Branch) {
					t.Fatalf("set apply left branch %q", creation.Assignment.Branch)
				}
			},
		},
		{
			name:     "discard ignored",
			modifier: "--discard-ignored",
			fixture: func(t *testing.T, root, home string) Creation {
				mustWrite(t, filepath.Join(root, ".gitignore"), []byte("ignored.txt\n"), 0o644)
				gitRun(t, root, "add", ".gitignore")
				gitRun(t, root, "commit", "-qm", "ignore set residue")
				creation := mustCreate(t, root, home, "set-discard-ignored", "discard ignored")
				landAssignment(t, root, creation, "landed.txt")
				mustWrite(t, filepath.Join(creation.Path, "ignored.txt"), []byte("residue\n"), 0o644)
				bare := runVerb(t, verbClean, repoHome{root, home}.call("--target", creation.Assignment.ID))
				if bare.exit != 0 || !strings.Contains(bare.stdout, creation.Path+",retain,") ||
					!strings.Contains(bare.stdout, "ignored residuals require --discard-ignored") {
					t.Fatalf("bare set plan = (%d, %q), want the ignored residue retained", bare.exit, bare.stdout)
				}
				return creation
			},
			effect: func(t *testing.T, _ string, creation Creation) {
				if _, err := os.Lstat(filepath.Join(creation.Path, "ignored.txt")); !os.IsNotExist(err) {
					t.Fatalf("set apply left the ignored residue in %s: %v", creation.Path, err)
				}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := newWorktreeRepo(t)
			f := repoHome{root, filepath.Join(root, ".bench-home")}
			creation := tc.fixture(t, f.root, f.home)
			plan := runVerb(t, verbClean, f.call(tc.modifier, "--target", creation.Assignment.ID))
			if plan.exit != 0 || plan.stderr != "" {
				t.Fatalf("%s plan = (%d, %q, %q), want one applicable plan", tc.name, plan.exit, plan.stdout, plan.stderr)
			}
			if !strings.Contains(plan.stdout, "bench worktree clean "+tc.modifier+" --target ") {
				t.Fatalf("%s plan = %q, want the modifier in the rendered apply command", tc.name, plan.stdout)
			}
			applied := runVerb(t, verbClean, f.call(tc.modifier, "--target", creation.Assignment.ID, "--apply", plan.mustFingerprint(t)))
			if applied.exit != 0 || applied.stderr != "" || strings.Count(applied.stdout, ",removed,") != 1 {
				t.Fatalf("%s apply = (%d, %q, %q), want one removal", tc.name, applied.exit, applied.stdout, applied.stderr)
			}
			if _, err := os.Lstat(creation.Path); !os.IsNotExist(err) {
				t.Fatalf("%s apply left %s: %v", tc.name, creation.Path, err)
			}
			tc.effect(t, f.root, creation)
		})
	}
}

func TestCleanSetGrammar(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	first := f.first
	before, err := intent.Assignments(f.root)
	mustNoError(t, err)
	for _, args := range [][]string{
		// A call that names no selection mode is as unscoped as one that names two, so the
		// same refusal answers both before the command reaches a planner.
		{},
		{"--discard-branch"},
		{"--discard-ignored", "--full"},
		{"--apply", strings.Repeat("a", 64)},
		{"--target", first.Assignment.ID, "--landed"},
		{"--landed", "--target", first.Assignment.ID},
		{"--discard-branch", "--unclaimed", "--target", first.Assignment.ID},
		{"--target", first.Assignment.ID, "--discard-branch", "--unclaimed"},
		{"--target", first.Assignment.ID, first.Path},
		{first.Path, "--target", first.Assignment.ID},
		{"--target"},
		{"--target", first.Assignment.ID, "--target"},
		{"--target", ""},
		{"--target", first.Assignment.ID, "--apply-current"},
		{"--target", first.Assignment.ID, "--apply", "short"},
	} {
		result := runVerb(t, verbClean, f.call(args...))
		if result.exit != 2 || result.stderr != "" || !strings.Contains(result.stdout, "invalid invocation; run "+usage.WorktreeClean) {
			t.Fatalf("args=%q exit=%d stdout=%q stderr=%q, want usage refusal", args, result.exit, result.stdout, result.stderr)
		}
		if rows := result.mustRows(t, cleanupTable); len(rows) != 1 {
			t.Fatalf("args=%q rows=%#v, want the usage row alone and no plan", args, rows)
		}
	}
	after, err := intent.Assignments(f.root)
	mustNoError(t, err)
	if len(after) != len(before) {
		t.Fatalf("grammar refusals changed the ledger: %#v -> %#v", before, after)
	}
	for _, creation := range []Creation{f.first, f.second} {
		if _, statErr := os.Stat(creation.Path); statErr != nil {
			t.Fatalf("grammar refusal removed %s: %v", creation.Path, statErr)
		}
	}
	if runVerb(t, verbClean, f.call("--target", first.Assignment.ID, "--discard-ignored", "--full")).exit == 2 {
		t.Fatal("the explicit set refused its own modifiers")
	}
}

func TestCleanSetHostileOperand(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, target string
	}{
		{name: "space", target: "hostile target name"},
		{name: "newline", target: "hostile\ntarget"},
		{name: "leading dash", target: "--landed"},
		{name: "command substitution", target: "$(touch pwned)"},
		{name: "shell separator", target: "; rm -rf /"},
		{name: "control byte", target: "hostile\x1btarget"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := landedSetFixture(t)
			plan := runVerb(t, verbClean, f.call("--target", f.first.Assignment.ID, "--target", tc.target))
			if plan.exit != 1 || plan.stderr != "" {
				t.Fatalf("hostile plan exit=%d stdout=%q stderr=%q, want a reported failure", plan.exit, plan.stdout, plan.stderr)
			}
			rows := plan.mustRows(t, cleanupTable)
			if len(rows) != 2 {
				t.Fatalf("hostile plan rows = %#v, want every selection outcome", rows)
			}
			plan.mustNoFingerprint(t)
			for _, row := range textRows(t, rows) {
				if row["fingerprint"] != "none" {
					t.Fatalf("hostile plan = %q, want no applicable fingerprint", plan.stdout)
				}
			}
			if strings.Contains(plan.stdout, "bench worktree clean --target") {
				t.Fatalf("hostile plan = %q, want no replayable action", plan.stdout)
			}
			if strings.ContainsRune(plan.stdout, '\x1b') || strings.Contains(plan.stdout, tc.target) && strings.ContainsAny(tc.target, "\n\x1b") {
				t.Fatalf("hostile plan = %q, want the operand rendered as data", plan.stdout)
			}
			for _, creation := range []Creation{f.first, f.second} {
				if _, err := os.Stat(creation.Path); err != nil {
					t.Fatalf("hostile plan removed %s: %v", creation.Path, err)
				}
			}
			if _, err := os.Lstat(filepath.Join(f.root, "pwned")); !os.IsNotExist(err) {
				t.Fatalf("hostile operand reached a shell: %v", err)
			}
		})
	}
}
