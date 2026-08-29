package worktree

import (
	"bytes"
	"testing"
)

// TestShowPrintsTheBlobAtTheRevision is S1: the verb prints the committed bytes of one
// tracked path at exit 0, so an agent reads a revision without the worktree path.
func TestShowPrintsTheBlobAtTheRevision(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "show-blob")
	commitInWorktree(t, creation.Path, "tracked.txt", "one\ntwo\n", "add tracked")
	var stdout, stderr bytes.Buffer
	if code := ShowCommand(root, home, []string{creation.Assignment.Label, "HEAD:tracked.txt"}, &stdout, &stderr); code != 0 {
		t.Fatalf("show exited %d, want 0: %s", code, stderr.String())
	}
	if stdout.String() != "one\ntwo\n" {
		t.Fatalf("show printed %q, want the committed bytes", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("show wrote %q to stderr, want nothing", stderr.String())
	}
}

// TestShowPassesNULBytesThrough is S2: the verb writes bytes and not lines, so a binary
// blob arrives as Git stores it.
func TestShowPassesNULBytesThrough(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "show-nul")
	blob := "a\x00b\n"
	commitInWorktree(t, creation.Path, "binary.bin", blob, "add binary")
	var stdout, stderr bytes.Buffer
	if code := ShowCommand(root, home, []string{creation.Assignment.Label, "HEAD:binary.bin"}, &stdout, &stderr); code != 0 {
		t.Fatalf("show exited %d, want 0: %s", code, stderr.String())
	}
	if stdout.String() != blob {
		t.Fatalf("show printed %q, want %q", stdout.String(), blob)
	}
}

// TestShowPassesGitsOwnFailureThrough is S3: a missing object returns Git's exit code
// and Git's own stderr, so a bad revision names itself rather than a Bench sentence.
func TestShowPassesGitsOwnFailureThrough(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "show-missing")
	commitInWorktree(t, creation.Path, "tracked.txt", "one\n", "add tracked")
	direct := descendant(t, "git", "-C", creation.Path, "cat-file", "blob", "HEAD:no-such-file")
	var directErr bytes.Buffer
	direct.Stderr = &directErr
	if err := direct.Run(); err == nil {
		t.Fatal("a direct cat-file of a missing object succeeded")
	}
	var stdout, stderr bytes.Buffer
	if code := ShowCommand(root, home, []string{creation.Assignment.Label, "HEAD:no-such-file"}, &stdout, &stderr); code != 128 {
		t.Fatalf("show exited %d, want 128: %s", code, stderr.String())
	}
	if stderr.String() != directErr.String() {
		t.Fatalf("show stderr = %q, want Git's own %q", stderr.String(), directErr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("show printed %q on a missing object, want nothing", stdout.String())
	}
}

// TestShowRefusesAnOperandThatIsNotARevision is S4 and S5: an operand without a `:`, and
// an operand that starts with `-`, each return the grammar line at exit 2. The target is
// unresolvable, so a run that reached the resolver would print the target refusal
// instead; the grammar line proves no Git ran.
func TestShowRefusesAnOperandThatIsNotARevision(t *testing.T) {
	t.Parallel()
	root, _, home := newOwnedAssignment(t, "show-operand")
	want := "usage: bench worktree show <target> <rev>:<path>\n"
	for name, operand := range map[string]string{
		"no colon":    "tracked.txt",
		"dash option": "--output=/tmp/x:tracked.txt",
	} {
		var stdout, stderr bytes.Buffer
		if code := ShowCommand(root, home, []string{"no-such-label", operand}, &stdout, &stderr); code != 2 {
			t.Fatalf("%s operand %q exited %d, want 2: %s", name, operand, code, stderr.String())
		}
		if stderr.String() != want {
			t.Fatalf("%s operand %q printed %q, want %q", name, operand, stderr.String(), want)
		}
		if stdout.String() != "" {
			t.Fatalf("%s operand %q printed %q on stdout, want nothing", name, operand, stdout.String())
		}
	}
}
