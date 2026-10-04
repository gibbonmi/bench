// Package commitmenttest owns repository fixtures for commitment tests.
package commitmenttest

import (
	"os"
	"path/filepath"
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
	data, err := commitment.Bytes(policy)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".bench", "commitment.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// Commit records the fixture's current tree on main.
func Commit(t testing.TB, root, message string) {
	t.Helper()
	gittest.Output(t, root, "add", "-A")
	gittest.Output(t, root, "commit", "-m", message)
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
