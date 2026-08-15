package greenmarker

import (
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := benchgit.Output(append([]string{"-C", root}, args...)...)
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(out)
}

func repo(t *testing.T) (root, tip string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	git(t, root, "commit", "-q", "--allow-empty", "-m", "one")
	return root, git(t, root, "rev-parse", "HEAD")
}

func TestReadAbsentMarkerIsNotPresent(t *testing.T) {
	root, _ := repo(t)
	if commit, present, err := Read(root, "main"); commit != "" || present || err != nil {
		t.Fatalf("Read = (%q, %v, %v), want absent", commit, present, err)
	}
}

func TestReadPeelsAnnotatedTagToCommit(t *testing.T) {
	root, tip := repo(t)
	git(t, root, "tag", "-a", "-m", "green", "green-tag", tip)
	git(t, root, "update-ref", Ref("main"), git(t, root, "rev-parse", "green-tag"))
	commit, present, err := Read(root, "main")
	if !present || err != nil || commit != tip {
		t.Fatalf("Read = (%q, %v, %v), want peeled tip %s", commit, present, err, tip)
	}
}

func TestReadDanglingSymbolicMarkerIsPresentWithError(t *testing.T) {
	root, _ := repo(t)
	git(t, root, "symbolic-ref", Ref("main"), "refs/heads/nowhere")
	if _, present, err := Read(root, "main"); !present || err == nil || !strings.Contains(err.Error(), "dangling symbolic ref") {
		t.Fatalf("Read = (present %v, err %v), want present dangling error", present, err)
	}
}

func TestAdvanceCreatesThenRefusesStaleExpectation(t *testing.T) {
	root, tip := repo(t)
	if err := Advance(root, "main", tip, ""); err != nil {
		t.Fatalf("Advance from absent: %v", err)
	}
	git(t, root, "commit", "-q", "--allow-empty", "-m", "two")
	next := git(t, root, "rev-parse", "HEAD")
	if err := Advance(root, "main", next, ""); err == nil {
		t.Fatal("Advance with absent expectation over an existing marker succeeded, want compare-and-swap refusal")
	}
	if err := Advance(root, "main", next, tip); err != nil {
		t.Fatalf("Advance with matching expectation: %v", err)
	}
	if err := Advance(root, "main", next, tip); err != nil {
		t.Fatalf("Advance already at destination: %v, want tolerated", err)
	}
	if got, _, _ := Read(root, "main"); got != next {
		t.Fatalf("marker = %s, want %s", got, next)
	}
}
