// Package commitmenttest owns repository fixtures for commitment tests.
package commitmenttest

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
)

// Repo creates one main-branch repository with policy committed.
func Repo(t testing.TB, policy commitment.Policy) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	WritePolicy(t, root, policy)
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte("# Roadmap\n\n## Recommended sequence\n\n1. old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	Commit(t, root, "initial")
	return root
}

// WritePolicy replaces the tracked commitment policy in root.
func WritePolicy(t testing.TB, root string, policy commitment.Policy) {
	t.Helper()
	if err := writePolicy(root, policy); err != nil {
		t.Fatal(err)
	}
}

// Commit records the fixture's current tree on main.
func Commit(t testing.TB, root, message string) {
	t.Helper()
	if err := commit(root, message); err != nil {
		t.Fatal(err)
	}
}

// CommitPolicy replaces the tracked commitment policy in root and commits that file alone.
// It returns its failure instead of failing a test, so a goroutine can call it.
func CommitPolicy(root string, policy commitment.Policy, message string) error {
	if err := writePolicy(root, policy); err != nil {
		return err
	}
	return commit(root, message, commitment.PolicyPath)
}

func writePolicy(root string, policy commitment.Policy) error {
	data, err := commitment.Bytes(policy)
	if err != nil {
		return err
	}
	path := filepath.Join(root, filepath.FromSlash(commitment.PolicyPath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// commit stages paths, or the whole tree when no path is named, and commits them.
func commit(root, message string, paths ...string) error {
	for _, args := range [][]string{append([]string{"add", "-A", "--"}, paths...), append([]string{"commit", "-m", message, "--"}, paths...)} {
		if _, err := gittest.Run(root, args...); err != nil {
			return err
		}
	}
	return nil
}

// Write writes one regular fixture file below root.
func Write(t testing.TB, root, path, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Planning creates an owned linked checkout for approval command tests.
func Planning(t testing.TB, root string) string {
	t.Helper()
	return Assignment(t, root, "planning")
}

// Assignment creates one owned linked checkout for an independent caller request.
func Assignment(t testing.TB, root, request string) string {
	t.Helper()
	id := intent.RequestDigest(request)[:32]
	owner := strings.Repeat("b", 32)
	branch := intent.AssignmentBranchRef(owner, id)
	target := filepath.Join(t.TempDir(), "planning")
	start := gittest.Output(t, root, "rev-parse", "HEAD")
	gittest.Output(t, root, "worktree", "add", "-b", strings.TrimPrefix(branch, "refs/heads/"), target, "HEAD")
	if err := intent.PutAssignment(root, intent.Assignment{Schema: intent.AssignmentRecordSchema, ID: id, OwnerID: owner, Request: intent.RequestDigest(request), Label: request, Start: start, Branch: branch, Worktree: target, State: intent.StateActive}); err != nil {
		t.Fatal(err)
	}
	return target
}

// SeedProtected writes an active milestone whose outcomes A then B own roadmap rows FT1
// and FT2, with the matching board and recommended sequence. The caller commits it.
func SeedProtected(t testing.TB, root string) {
	t.Helper()
	outcomes := []commitment.Outcome{}
	for i, id := range []string{"A", "B"} {
		row := "FT" + strconv.Itoa(i+1)
		path := "roadmap/" + row + ".md"
		body := "**" + row + " — " + id + "**\n\nKeep the " + id + " obligation.\n"
		Write(t, root, path, body)
		outcomes = append(outcomes, commitment.Outcome{ID: id, Criteria: []commitment.Criterion{{ID: id + "-done", Text: "The obligation is satisfied."}}, Sources: []commitment.SourceBinding{{ID: row, Path: path, Identity: commitment.Identity([]byte(body))}}})
	}
	WritePolicy(t, root, commitment.Policy{Version: 1, ActiveMilestone: "M", Milestones: []commitment.Milestone{{ID: "M", Outcomes: outcomes}}})
	Write(t, root, "ROADMAP.md", "# Roadmap\n\n## Parked\n\n**FT1 — A**\n\n**FT2 — B**\n\n## Recommended sequence\n\n1. A\n2. B\n")
}
