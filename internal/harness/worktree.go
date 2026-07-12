// Package harness owns deterministic adapters from harness event envelopes to
// Bench's harness-neutral lifecycle. It parses events and derives opaque request
// IDs; ownership, assignment, locking, recovery, and cleanup remain in worktree.
package harness

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/worktree"
)

type claudeCreateEvent struct {
	SessionID string `json:"session_id"`
	Cwd       string `json:"cwd"`
	Name      string `json:"name"`
}

type claudeRemoveEvent struct {
	SessionID    string `json:"session_id"`
	WorktreePath string `json:"worktree_path"`
}

// WorktreeCommand implements the plumbing-only `bench worktree-hook create|remove`
// surface. The shell hook passes stdin through unchanged; all event validation and
// request derivation happen here before the shared lifecycle is called.
func WorktreeCommand(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) != 1 || (args[0] != "create" && args[0] != "remove") {
		fmt.Fprintln(stderr, "usage: bench worktree-hook create|remove")
		return 2
	}
	switch args[0] {
	case "create":
		return claudeCreate(stdin, stdout, stderr)
	default:
		return claudeRemove(stdin, stderr)
	}
}

func claudeCreate(stdin io.Reader, stdout, stderr io.Writer) int {
	var event claudeCreateEvent
	if err := decodeEvent(stdin, &event); err != nil {
		fmt.Fprintf(stderr, "bench worktree-hook create: invalid event JSON: %v\n", err)
		return 1
	}
	if event.SessionID == "" || event.Cwd == "" || strings.TrimSpace(event.Name) == "" {
		fmt.Fprintln(stderr, "bench worktree-hook create: event requires session_id, cwd, and name")
		return 1
	}
	root, err := git.RootAt(event.Cwd)
	if err != nil {
		fmt.Fprintf(stderr, "bench worktree-hook create: cwd is not in a Git repository: %v\n", err)
		return 1
	}
	creation, err := worktree.Create(root, claudeRequestID(event.SessionID), event.Name, nil)
	if err != nil {
		fmt.Fprintf(stderr, "bench worktree-hook create: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, creation.Path)
	return 0
}

func claudeRemove(stdin io.Reader, stderr io.Writer) int {
	var event claudeRemoveEvent
	if err := decodeEvent(stdin, &event); err != nil {
		fmt.Fprintf(stderr, "bench worktree-hook remove: invalid event JSON: %v\n", err)
		return 1
	}
	if event.SessionID == "" || event.WorktreePath == "" {
		fmt.Fprintln(stderr, "bench worktree-hook remove: event requires session_id and worktree_path")
		return 1
	}
	repositoryHint := os.Getenv("CLAUDE_PROJECT_DIR")
	if repositoryHint == "" {
		repositoryHint = event.WorktreePath
	}
	root, err := git.RootAt(repositoryHint)
	if err != nil {
		fmt.Fprintf(stderr, "bench worktree-hook remove: repository unavailable: %v\n", err)
		return 1
	}
	registrations, err := git.Worktrees(root)
	if err != nil || len(registrations) == 0 {
		fmt.Fprintf(stderr, "bench worktree-hook remove: worktree registration unavailable: %v\n", err)
		return 1
	}
	return worktree.ReleaseCommand(registrations[0].Path, []string{"--request", claudeRequestID(event.SessionID), event.WorktreePath}, io.Discard, stderr)
}

func decodeEvent(stdin io.Reader, dst any) error {
	const maxEventBytes = 1 << 20
	data, err := io.ReadAll(io.LimitReader(stdin, maxEventBytes+1))
	if err != nil {
		return err
	}
	if len(data) > maxEventBytes {
		return fmt.Errorf("event exceeds %d-byte limit", maxEventBytes)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return fmt.Errorf("empty input")
	}
	return json.Unmarshal(data, dst)
}

func claudeRequestID(sessionID string) string {
	sum := sha256.Sum256([]byte("claude-worktree/v1\x00" + sessionID))
	return "claude-worktree/v1:" + hex.EncodeToString(sum[:])
}
