package repository_test

import (
	"testing"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gittest"
)

// TestDefaultRevisionAnswersTheCommitIdentity pins that the default revision is the commit
// of the default-branch tip. An annotated tag named like the branch does not put the tag
// object in place of the commit.
func TestDefaultRevisionAnswersTheCommitIdentity(t *testing.T) {
	t.Parallel()
	root := gittest.RepoOnBranch(t, "main")
	gittest.Commit(t, root, "c0", "--allow-empty")
	commit := gittest.Output(t, root, "rev-parse", "HEAD")
	gittest.MustRun(t, root, "tag", "-a", "main", "-m", "a tag named like the branch")

	branch, revision, err := commitrepo.DefaultRevision(commitrepo.Store{Root: root})
	if err != nil || branch != "main" || revision != commit {
		t.Fatalf("default revision = (%q, %q, %v), want (main, %s, nil)", branch, revision, err, commit)
	}
}
