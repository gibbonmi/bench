// Package writeguard classifies a harness file-tool write against the primary
// checkout. Within Bench, main receives writes only through landings, and the Bash
// guards plus `bench commit` already hold that boundary. A file tool reaches the same
// tree without a shell, so this guard closes the remaining route.
package writeguard

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/usage"
)

// PathFromEnvelope returns the absolute path a PreToolUse file-tool envelope writes.
// A relative file_path resolves against the envelope's own cwd, because the hook
// process's ambient working directory is not the directory the tool call ran in.
//
// It returns an error when the envelope carries no readable path. The caller owns the
// fail posture: a file edit the guard cannot read stays open, so nothing is refused on
// a path nobody read.
func PathFromEnvelope(data []byte) (string, error) {
	var envelope struct {
		Cwd       string                     `json:"cwd"`
		ToolInput map[string]json.RawMessage `json:"tool_input"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return "", err
	}
	raw, ok := envelope.ToolInput["file_path"]
	if !ok {
		return "", errors.New("file_path field missing")
	}
	var path string
	if err := json.Unmarshal(raw, &path); err != nil {
		return "", err
	}
	if path == "" {
		return "", errors.New("file_path field empty")
	}
	if strings.IndexByte(path, 0) >= 0 {
		return "", errors.New("file_path field has control byte")
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	if envelope.Cwd == "" {
		return "", errors.New("relative file_path with no cwd")
	}
	return filepath.Join(envelope.Cwd, path), nil
}

// Checker holds the four repository facts the verdict rests on. Each one is injected,
// so the rule is graded without a git process and the verb keeps the one Go owner of
// each fact.
type Checker struct {
	// RootAt returns the top-level directory of the repository holding dir.
	RootAt func(dir string) (string, error)
	// IsPrimary reports whether root is the repository's primary checkout.
	IsPrimary func(root string) (bool, error)
	// IsTracked reports whether git tracks path in root.
	IsTracked func(root, path string) bool
	// IsIgnored reports whether root's ignore rules cover path.
	IsIgnored func(root, path string) bool
}

// Verdict is the guard's decision for one write. Path is the absolute path the
// envelope named, so the refusal can name the file the reader must move.
type Verdict struct {
	Blocked bool
	Path    string
}

// Classify decides one write. It refuses a tracked path under the primary checkout,
// and it allows everything else: a path in a Bench worktree, a path outside every
// repository, and a path git ignores.
//
// The ignore test runs before the tracked test, and it allows. An ignored path is a
// local working note, and git.LocalNoteRoot sends such a note to the primary checkout
// on purpose, so the guard must not refuse the write Bench itself directs there.
//
// Any fact the checker cannot answer allows the write. The guard denies only on a fact
// it read.
func Classify(path string, c Checker) Verdict {
	root, err := c.RootAt(filepath.Dir(path))
	if err != nil || root == "" {
		return Verdict{Path: path}
	}
	primary, err := c.IsPrimary(root)
	if err != nil || !primary {
		return Verdict{Path: path}
	}
	if c.IsIgnored(root, path) || !c.IsTracked(root, path) {
		return Verdict{Path: path}
	}
	return Verdict{Blocked: true, Path: path}
}

// Message returns the refusal line: the fixed prefix, the path, and the one refusal
// every Bench write verb prints from the primary checkout, so the reader gets the same
// repair here as from `bench commit`.
func (v Verdict) Message() string {
	return "BLOCKED: " + v.Path + " is tracked in the primary checkout. " + usage.PrimaryCheckoutRefusal()
}
