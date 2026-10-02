package worktree

// This file declares each named fixture value that a worktree test builder returns. A
// builder stays in its own test file. A value method can build a verb call, but it runs
// no verb, so the verb runner stays the one way to run a verb. The text row reader takes
// the rows that mustRows already decoded, so it parses no rendered text.

import "testing"

// repoHome is a repository and the private home that its registrations live under.
type repoHome struct {
	root string
	home string
}

// call builds a verb call at the fixture's root and home.
func (f repoHome) call(args ...string) verbCall {
	return verbCall{root: f.root, home: f.home, args: args}
}

// callWith builds a verb call at the fixture's root and home that runs the verb's joins
// form with j.
func (f repoHome) callWith(j joins, args ...string) verbCall {
	call := f.call(args...)
	call.joins = &j
	return call
}

// ownedAssignment is a repository, its one owned registration, and the private home
// that the registration lives under.
type ownedAssignment struct {
	repoHome
	creation Creation
}

// reauthorizeSet is an owned assignment and the reviewed base and tip that a reauthorize
// names for it.
type reauthorizeSet struct {
	ownedAssignment
	base, tip string
}

// mergeSet is a repository, one owned registration for each label, the seam set that the
// merge verb runs with, and the file that the seam set's lane appends to.
type mergeSet struct {
	repoHome
	joins   joins
	tally   string
	created []Creation
}

// merge builds a merge verb call that runs the merge verb's joins form with the set's seams.
func (f mergeSet) merge(args ...string) verbCall {
	return f.callWith(f.joins, args...)
}

// reclaimPoolFixture is a repository, a private home, and the pool parent under that
// home. The pool parent holds one key for each repository, and the builder creates it.
type reclaimPoolFixture struct {
	repoHome
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

// landedSet is a repository with two landed clean members and one landed member whose
// tracked file is dirty.
type landedSet struct {
	repoHome
	first, second, dirty Creation
}

// removableSet is a repository whose landed clean members each plan a removal, with each
// member's landed file by assignment identity.
type removableSet struct {
	repoHome
	creations []Creation
	files     map[string]string
}

// retainedMemberSet is a repository with one removable member and one member whose ignored
// residue the plan retains.
type retainedMemberSet struct {
	repoHome
	removable, retained Creation
}

// refusedRelease is an owned assignment whose release refused, with the stderr of that
// refusal.
type refusedRelease struct {
	ownedAssignment
	stderr string
}

// landingFixture is an owned assignment that a landing publishes: the landing base, the
// reviewed source tip, and the file that the fixture's gate scripts append to on each run.
type landingFixture struct {
	ownedAssignment
	base, tip, tally string
}

// foldedLanding is a landing fixture whose source folded a destination advance. Its base is
// the advanced destination tip, and fold is the fold commit that a review reads as its base.
type foldedLanding struct {
	landingFixture
	fold string
}

// foldedSibling is a sibling assignment that a landing source folded, and the source tip
// after that fold.
type foldedSibling struct {
	sibling Creation
	tip     string
}

// cleanupTable is the table block that the clean verb renders its rows in.
const cleanupTable = "worktree_cleanup"

// textRow is one decoded table row: each field and its text cell.
type textRow map[string]string

// textRows reads decoded table rows as text rows. It fails t on a row that is not an object
// and on a cell that is not text.
func textRows(t testing.TB, rows []any) []textRow {
	t.Helper()
	read := make([]textRow, 0, len(rows))
	for i, row := range rows {
		fields, ok := row.(map[string]any)
		if !ok {
			t.Fatalf("row %d decoded as %T, want an object", i, row)
			return nil
		}
		cells := make(textRow, len(fields))
		for field, value := range fields {
			text, ok := value.(string)
			if !ok {
				t.Fatalf("row %d field %q = %#v, want text", i, field, value)
				return nil
			}
			cells[field] = text
		}
		read = append(read, cells)
	}
	return read
}

// textRowsBy reads decoded table rows as text rows keyed by the cell of field.
func textRowsBy(t testing.TB, rows []any, field string) map[string]textRow {
	t.Helper()
	keyed := map[string]textRow{}
	for _, row := range textRows(t, rows) {
		keyed[row[field]] = row
	}
	return keyed
}
