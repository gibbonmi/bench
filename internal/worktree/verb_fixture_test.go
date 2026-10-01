package worktree

// This file declares each named fixture value that a worktree test builder returns. A
// builder stays in its own test file. A value method can build a verb call, but it runs
// no verb, so the verb runner stays the one way to run a verb.

// restoredAssignment is an owned assignment after one reset preserved its dirty checkout
// under ref.
type restoredAssignment struct {
	root     string
	creation Creation
	home     string
	ref      string
}

// call builds a verb call at the fixture's root and home.
func (f restoredAssignment) call(args ...string) verbCall {
	return verbCall{root: f.root, home: f.home, args: args}
}
