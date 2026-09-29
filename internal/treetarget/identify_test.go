package treetarget

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
)

// The cell values here are authored apart from Identify: the spec fixes `primary`,
// `unassigned`, `none`, `unknown`, and a bare bool for the dirty cell. A fallback that
// prints another word, or an empty cell, reds these rows.

// gitIn runs git in root and answers its trimmed output.
func gitIn(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// committedRepo answers a primary checkout whose one commit tracks tracked.txt.
func committedRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, root, "add", "tracked.txt")
	gitIn(t, root, "commit", "-qm", "base")
	return root
}

// linkedWorktree adds a linked worktree of root in a pool directory whose name no label
// shares, and answers its path.
func linkedWorktree(t *testing.T, root string) string {
	t.Helper()
	linked := filepath.Join(t.TempDir(), "pool-dir-7f3a")
	gitIn(t, root, "worktree", "add", "-q", "-b", "linked", linked)
	return linked
}

// putAssignment records one assignment labeled label in state that owns the tree at path.
func putAssignment(t *testing.T, root, path, label string, state intent.AssignmentState) {
	t.Helper()
	id, owner := strings.Repeat("a", 32), strings.Repeat("b", 32)
	err := intent.PutAssignment(root, intent.Assignment{
		Schema: intent.AssignmentRecordSchema, ID: id, OwnerID: owner, Request: intent.RequestDigest("tree-target-" + label),
		Label: label, Start: gitIn(t, root, "rev-parse", "HEAD"), Branch: intent.AssignmentBranchRef(owner, id),
		Worktree: path, State: state,
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TT12, TT23: an active assignment names its tree by label, and the row holds no part of
// the tree's path.
func TestIdentifyNamesActiveLabel(t *testing.T) {
	root := committedRepo(t)
	linked := linkedWorktree(t, root)
	putAssignment(t, root, linked, "alpha", intent.StateActive)
	identity := Identify(linked)
	row, err := Row(identity)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(row, linked) || strings.Contains(row, filepath.Base(linked)) {
		t.Fatalf("row %q holds the worktree path %q", row, linked)
	}
	if identity.Target != "alpha" {
		t.Fatalf("target = %q, want the active label alpha", identity.Target)
	}
}

// The ledger accepts a label with a control byte that TOON refuses. The row still prints,
// with the byte escaped, so no raw control byte reaches the response.
func TestRowEscapesHostileLabel(t *testing.T) {
	root := committedRepo(t)
	linked := linkedWorktree(t, root)
	putAssignment(t, root, linked, "al\x07pha", intent.StateActive)
	row, err := Row(Identify(linked))
	if want := `  "al\\u0007pha",`; err != nil || !strings.Contains(row, "\n"+want) || strings.ContainsRune(row, '\x07') {
		t.Fatalf("hostile label row = (%q, %v), want the row with target cell %q", row, err, want)
	}
}

// TT13, TT14: a linked worktree that no active assignment owns is unassigned, both with
// no ledger and with a released assignment.
func TestIdentifyUnownedWorktree(t *testing.T) {
	t.Run("no assignment", func(t *testing.T) {
		root := committedRepo(t)
		if target := Identify(linkedWorktree(t, root)).Target; target != "unassigned" {
			t.Fatalf("target = %q, want unassigned", target)
		}
	})
	t.Run("released assignment", func(t *testing.T) {
		root := committedRepo(t)
		linked := linkedWorktree(t, root)
		putAssignment(t, root, linked, "alpha", intent.StateComplete)
		if target := Identify(linked).Target; target != "unassigned" {
			t.Fatalf("target = %q, want unassigned for a released assignment", target)
		}
	})
}

// TT16: an unborn HEAD renders the head cell `none`, and the dirty cell prints bare.
func TestIdentifyUnbornHead(t *testing.T) {
	row, err := Row(Identify(gittest.Repo(t)))
	if want := "tree[1]{target,head,dirty}:\n  primary,none,false\n"; err != nil || row != want {
		t.Fatalf("unborn row = (%q, %v), want %q", row, err, want)
	}
}

// TT17, TT18, TT19: an untracked file and a modified tracked file render dirty `true`, and
// a corrupt index renders dirty `unknown`.
func TestIdentifyDirtyStates(t *testing.T) {
	for _, row := range []struct {
		name  string
		file  string
		body  string
		dirty any
	}{
		{"untracked file", "untracked.txt", "new\n", true},
		{"modified tracked file", "tracked.txt", "changed\n", true},
		{"corrupt index", filepath.Join(".git", "index"), "garbage", "unknown"},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := committedRepo(t)
			if err := os.WriteFile(filepath.Join(root, row.file), []byte(row.body), 0o644); err != nil {
				t.Fatal(err)
			}
			if dirty := Identify(root).Dirty; dirty != row.dirty {
				t.Fatalf("dirty = %#v, want %#v", dirty, row.dirty)
			}
		})
	}
}
