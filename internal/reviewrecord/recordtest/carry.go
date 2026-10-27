package recordtest

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
)

// Carrying prepares a delegated one-chunk fixture whose chunk plan is owed, and records the
// chunk at its own commit with a passing prior entry for each requirement, bound to that
// commit. A later CommitDelta then gives the delta that a carry decision grades.
func Carrying(t testing.TB, owed []rr.Requirement) *Fixture {
	t.Helper()
	f := Prepare(t, gittest.RepoOnBranch(t, "main"), 1, Spec, Body, Owing(owed), Delegate)
	f.AddChunk()
	return f
}

// Owing makes owed the verification plan of the first chunk. It precedes Delegate in the option
// list, so Delegate gives each requirement its ticket.
func Owing(owed []rr.Requirement) Option {
	return func(plan *rr.Plan) {
		plan.Chunks[0].Verification = append([]rr.Requirement{}, owed...)
	}
}

// CommitDelta appends one line to each of paths, which creates an absent file, commits the
// change, and moves the last recorded chunk to that commit through Refreeze. Each entry keeps
// the commit that it graded.
func (f *Fixture) CommitDelta(paths ...string) {
	f.T.Helper()
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join(f.Root, filepath.FromSlash(path)))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			f.T.Fatal(err)
		}
		f.Write(path, string(data)+"// carry delta\n")
	}
	f.Commit("carry delta")
	f.Refreeze()
}

// Legacy makes each entry of requirement id in the last recorded chunk a hand-built legacy
// entry: it has no source commit, and its digest is the chunk source digest.
func (f *Fixture) Legacy(id string) {
	f.T.Helper()
	chunk := &f.Record.Chunks[len(f.Record.Chunks)-1]
	for i := range chunk.Verification {
		if chunk.Verification[i].Requirement == id {
			chunk.Verification[i].SourceCommit, chunk.Verification[i].SourceDigest = "", chunk.SourceDigest
		}
	}
}

// SavedChunk saves the fixture record and answers its last chunk entry as the record reader
// reads it from the saved file.
func (f *Fixture) SavedChunk() rr.Chunk {
	f.T.Helper()
	f.Save()
	record, err := rr.Read(f.Root, f.Record.Spec)
	if err != nil {
		f.T.Fatal(err)
	}
	return record.Chunks[len(record.Chunks)-1]
}
