package recordtest

import (
	benchgit "github.com/gibbonmi/bench/internal/git"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
)

// mainline is the primary checkout of a NewLinked fixture, which holds the default branch.
func (f *Fixture) mainline() *Fixture {
	f.T.Helper()
	if f.primary == "" {
		f.T.Fatal("a default-branch change needs a NewLinked fixture")
	}
	return &Fixture{T: f.T, Root: f.primary}
}

// Mainline commits files onto the default branch and returns the new default-branch tip.
// The fixture checkout keeps its branch.
func (f *Fixture) Mainline(files map[string]string) string {
	f.T.Helper()
	main := f.mainline()
	for path, data := range files {
		main.Write(path, data)
	}
	main.Commit("default branch change")
	return main.Tip()
}

// SyncMainline fast-forwards the default branch to the fixture HEAD, so a later
// default-branch change can conflict with a later fixture change.
func (f *Fixture) SyncMainline() {
	f.T.Helper()
	f.mainline().Git("merge", "--ff-only", f.Tip())
}

// HandFold merges incoming into the fixture checkout with git merge, lets resolve edit the
// checkout, and commits the fold. A conflict keeps Git's markers for resolve to replace. A
// nil resolve commits the merge as Git wrote it. The answer is the fold commit.
func (f *Fixture) HandFold(incoming string, resolve func()) string {
	f.T.Helper()
	// A conflict exits nonzero, so the merge state below is the test of a started merge.
	_, _ = benchgit.Raw("-C", f.Root, "merge", "--no-ff", "--no-commit", incoming)
	f.Git("rev-parse", "--verify", "--quiet", "MERGE_HEAD")
	if resolve != nil {
		resolve()
	}
	f.Commit("fold " + incoming)
	return f.Tip()
}

// RecordFold records the fold at HEAD through the record owner with delta in place of the
// derived delta, and reloads the record. The record file must hold f.Record first.
func (f *Fixture) RecordFold(delta rr.FoldDelta) rr.Fold {
	f.T.Helper()
	derive := func(string, string) (rr.FoldDelta, error) { return delta, nil }
	fold, err := rr.RecordFold(f.Root, f.Record.Spec, f.Tip(), derive)
	if err != nil {
		f.T.Fatal(err)
	}
	if f.Record, err = rr.Read(f.Root, f.Record.Spec); err != nil {
		f.T.Fatal(err)
	}
	return fold
}
