package treetarget

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The argv line and the refusal pair here are authored apart from Run: the spec fixes the
// family verb and its leaf as two child tokens and names the leaf in each refusal.

// ER29: a family leaf call starts its child with the verb and the leaf as separate argv
// tokens, before the leaf's own arguments, and a target that names no active assignment
// refuses with the leaf in its command before any child starts.
func TestRunStartsFamilyLeafChild(t *testing.T) {
	r := newTargetRepo(t, t.TempDir())
	r.create(t, "tree-target-alpha", "alpha")
	dir := t.TempDir()
	script, marker := filepath.Join(dir, "child"), filepath.Join(dir, "marker")
	writeScript(t, script, "#!/bin/sh\nprintf '[%s]' \"$@\" > '"+marker+"'\n")
	run := func(args ...string) (string, string, int) {
		var out, errOut bytes.Buffer
		code := Run(Call{Name: "worktree", Leaf: "read", Flag: "--in", Args: args, Root: r.root, Home: r.home, Running: script,
			Stdin: strings.NewReader(""), Stdout: &out, Stderr: &errOut})
		return out.String(), errOut.String(), code
	}
	if stdout, stderr, code := run("beta", "notes.md"); code != 1 || stdout != "" || stderr != "bench worktree read --in: target is unassigned\nnext=bench worktree list\n" {
		t.Fatalf("worktree read --in beta = (%d, %q, %q), want the unassigned refusal of bench worktree read --in", code, stdout, stderr)
	}
	requireNoChild(t, marker)
	if stdout, stderr, code := run("alpha", "notes.md", "--lines", "1:2"); code != 0 || stdout != "" || stderr != "" {
		t.Fatalf("worktree read --in alpha = (%d, %q, %q), want the child's exit 0", code, stdout, stderr)
	}
	argv, err := os.ReadFile(marker)
	if err != nil || string(argv) != "[worktree][read][notes.md][--lines][1:2]" {
		t.Fatalf("child argv = %q (%v), want [worktree][read][notes.md][--lines][1:2]", argv, err)
	}
}
