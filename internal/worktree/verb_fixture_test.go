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

// reclaimPoolFixture is a repository, a private home, and the pool parent under that
// home. The pool parent holds one key for each repository, and the builder creates it.
type reclaimPoolFixture struct {
	root string
	home string
	pool string
}

// repoPoolFixture is a repository, a private home, and the pool root of that one
// repository under that home. The builder computes the pool root and does not create it.
type repoPoolFixture struct {
	root     string
	home     string
	poolRoot string
}

// restoredAssignment is an owned assignment after one reset preserved its dirty checkout
// under ref.
type restoredAssignment struct {
	ownedAssignment
	ref string
}
