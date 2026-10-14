package responsebound

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// retainedName is the file of the retained response in an assignment scope. It carries no
// spill suffix, so the primary prune never reads it, and Drop removes it with the scope.
const retainedName = "retained.json"

// Retained is the last complete response of one assignment worktree and the call that
// printed it. Verb is the verb word as the operator typed it, and Args are the arguments
// after it. Head and Dirty are the worktree's commit and dirty paths after the verb returned.
type Retained struct {
	Verb   string   `json:"verb"`
	Args   []string `json:"args"`
	Exit   int      `json:"exit"`
	Head   string   `json:"head"`
	Dirty  []string `json:"dirty"`
	Output []byte   `json:"output"`
}

// retainedPath answers the assignment scope directory of the worktree at root and the
// retained response path in it. It answers false where assignmentDir answers false.
func retainedPath(home, root string) (scopeDir, string, bool, error) {
	dir, ok, err := assignmentDir(home, root)
	if err != nil || !ok {
		return dir, "", false, err
	}
	path, err := dir.file(retainedName)
	return dir, path, err == nil, err
}

// Retain replaces the retained response of the assignment worktree at root with the
// complete output of the finished response o. The call function answers the other fields
// of the retained response, and Retain runs it only when the write follows. A response with
// no complete copy clears the retained response. A call function that returns an error
// clears it too. A root outside any assignment worktree keeps nothing.
func (o *Owner) Retain(home, root string, call func() (Retained, error)) error {
	dir, path, ok, err := retainedPath(home, root)
	if err != nil {
		return fmt.Errorf("retain response: %w", err)
	}
	if !ok {
		return nil
	}
	if err := dir.create(); err != nil {
		return fmt.Errorf("retain response: %w", err)
	}
	output, ok := o.completeOutput()
	if !ok {
		return clearRetained(path)
	}
	retained, err := call()
	if err != nil {
		return errors.Join(fmt.Errorf("retain response: %w", err), clearRetained(path))
	}
	retained.Output = output
	data, err := json.Marshal(retained)
	if err != nil {
		return fmt.Errorf("retain response: %w", err)
	}
	return durablefile.Replace(path, data, privateFile)
}

func clearRetained(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("clear retained response: %w", err)
	}
	return nil
}

// ReadRetained answers the retained response of the assignment worktree at root. An absent
// response, and a root outside any assignment worktree, answer an error that matches
// fs.ErrNotExist.
func ReadRetained(home, root string) (Retained, error) {
	_, path, ok, err := retainedPath(home, root)
	if err != nil {
		return Retained{}, fmt.Errorf("read retained response: %w", err)
	}
	if !ok {
		return Retained{}, fmt.Errorf("read retained response: no assignment worktree: %w", fs.ErrNotExist)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Retained{}, fmt.Errorf("read retained response: %w", err)
	}
	var retained Retained
	if err := json.Unmarshal(data, &retained); err != nil {
		return Retained{}, fmt.Errorf("read retained response: %w", err)
	}
	return retained, nil
}
