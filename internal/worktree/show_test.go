package worktree

import (
	"bytes"
	"testing"
)

// TestShowPrintsTheBlobAtTheRevision is S1: the verb prints the committed bytes of one
// tracked path at exit 0, so an agent reads a revision without the worktree path.
func TestShowPrintsTheBlobAtTheRevision(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "show-blob")
	commitInWorktree(t, f.creation.Path, "tracked.txt", "one\ntwo\n", "add tracked")
	r := runVerb(t, verbShow, f.call(f.creation.Assignment.Label, "HEAD:tracked.txt"))
	if r.exit != 0 {
		t.Fatalf("show exited %d, want 0: %s", r.exit, r.stderr)
	}
	if r.stdout != "one\ntwo\n" {
		t.Fatalf("show printed %q, want the committed bytes", r.stdout)
	}
	if r.stderr != "" {
		t.Fatalf("show wrote %q to stderr, want nothing", r.stderr)
	}
}

// TestShowPassesNULBytesThrough is S2: the verb writes bytes and not lines, so a binary
// blob arrives as Git stores it.
func TestShowPassesNULBytesThrough(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "show-nul")
	blob := "a\x00b\n"
	commitInWorktree(t, f.creation.Path, "binary.bin", blob, "add binary")
	r := runVerb(t, verbShow, f.call(f.creation.Assignment.Label, "HEAD:binary.bin"))
	if r.exit != 0 {
		t.Fatalf("show exited %d, want 0: %s", r.exit, r.stderr)
	}
	if r.stdout != blob {
		t.Fatalf("show printed %q, want %q", r.stdout, blob)
	}
}

// TestShowPassesGitsOwnFailureThrough is S3: a missing object returns Git's exit code
// and Git's own stderr, so a bad revision names itself rather than a Bench sentence.
func TestShowPassesGitsOwnFailureThrough(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "show-missing")
	commitInWorktree(t, f.creation.Path, "tracked.txt", "one\n", "add tracked")
	direct := descendant(t, "git", "-C", f.creation.Path, "cat-file", "blob", "HEAD:no-such-file")
	var directErr bytes.Buffer
	direct.Stderr = &directErr
	if err := direct.Run(); err == nil {
		t.Fatal("a direct cat-file of a missing object succeeded")
	}
	r := runVerb(t, verbShow, f.call(f.creation.Assignment.Label, "HEAD:no-such-file"))
	if r.exit != 128 {
		t.Fatalf("show exited %d, want 128: %s", r.exit, r.stderr)
	}
	if r.stderr != directErr.String() {
		t.Fatalf("show stderr = %q, want Git's own %q", r.stderr, directErr.String())
	}
	if r.stdout != "" {
		t.Fatalf("show printed %q on a missing object, want nothing", r.stdout)
	}
}

// TestShowRefusesAnOperandThatIsNotARevision is S4 and S5: an operand without a `:`, and
// an operand that starts with `-`, each return the grammar line at exit 2. The target is
// unresolvable, so a run that reached the resolver would print the target refusal
// instead; the grammar line proves no Git ran.
func TestShowRefusesAnOperandThatIsNotARevision(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "show-operand")
	want := "usage: bench worktree show <target> <rev>:<path>\n"
	for name, operand := range map[string]string{
		"no colon":    "tracked.txt",
		"dash option": "--output=/tmp/x:tracked.txt",
	} {
		r := runVerb(t, verbShow, f.call("no-such-label", operand))
		if r.exit != 2 {
			t.Fatalf("%s operand %q exited %d, want 2: %s", name, operand, r.exit, r.stderr)
		}
		if r.stderr != want {
			t.Fatalf("%s operand %q printed %q, want %q", name, operand, r.stderr, want)
		}
		if r.stdout != "" {
			t.Fatalf("%s operand %q printed %q on stdout, want nothing", name, operand, r.stdout)
		}
	}
}

// TestShowHelpPrintsTheGrammarLine is S4's help half: the sole help spelling is a request
// and not an operand, so it answers with the grammar line at exit 0.
func TestShowHelpPrintsTheGrammarLine(t *testing.T) {
	t.Parallel()
	r := runVerb(t, verbShow, repoHome{t.TempDir(), t.TempDir()}.call("--help"))
	if r.exit != 0 {
		t.Fatalf("help exited %d, want 0: %s", r.exit, r.stderr)
	}
	if want := "usage: bench worktree show <target> <rev>:<path>\n"; r.stderr != want {
		t.Fatalf("help printed %q, want %q", r.stderr, want)
	}
	if r.stdout != "" {
		t.Fatalf("help printed %q on stdout, want nothing", r.stdout)
	}
}

// TestShowRefusesAControlByteInTheOperand pins the operand's line-safety check. A newline
// would break the one-line refusal shape, so the operand refuses at exit 2 with the
// grammar line. The target is unresolvable, so the grammar line proves no Git ran.
func TestShowRefusesAControlByteInTheOperand(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "show-control")
	r := runVerb(t, verbShow, f.call("no-such-label", "HEAD:a\nb"))
	if r.exit != 2 {
		t.Fatalf("control-byte operand exited %d, want 2: %s", r.exit, r.stderr)
	}
	if want := "usage: bench worktree show <target> <rev>:<path>\n"; r.stderr != want {
		t.Fatalf("control-byte operand printed %q, want %q", r.stderr, want)
	}
	if r.stdout != "" {
		t.Fatalf("control-byte operand printed %q on stdout, want nothing", r.stdout)
	}
}
