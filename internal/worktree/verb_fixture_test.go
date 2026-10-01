package worktree

// This file declares each named fixture value that a worktree test builder returns. A
// builder stays in its own test file. A value method can build a verb call, but it runs
// no verb, so the verb runner stays the one way to run a verb.

// ownedAssignment is a repository, its one owned registration, and the private home
// that the registration lives under.
type ownedAssignment struct {
	root     string
	creation Creation
	home     string
}

// call builds a verb call at the fixture's root and home.
func (f ownedAssignment) call(args ...string) verbCall {
	return verbCall{root: f.root, home: f.home, args: args}
}

// restoredAssignment is an owned assignment after one reset preserved its dirty checkout
// under ref.
type restoredAssignment struct {
	ownedAssignment
	ref string
}
