package responsebound

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/poolkey"
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

// retainedPath answers the store levels below home of the assignment scope of the worktree
// at root, the retained response path in that scope, and false for a root outside any
// assignment worktree.
func retainedPath(home, root string) ([]string, string, bool) {
	scope, ok := AssignmentScope(home, root)
	if !ok {
		return nil, "", false
	}
	levels := []string{storeDirName, poolkey.Key(root), scope}
	return levels, filepath.Join(slices.Concat([]string{home}, levels, []string{retainedName})...), true
}

// Retain replaces the retained response of the assignment worktree at root with the complete
// output of the finished response o and the call that call answers. A response with no
// complete copy clears the retained response, and so does a failed call; call runs only when
// the write follows. A root outside any assignment worktree keeps nothing.
func (o *Owner) Retain(home, root string, call func() (Retained, error)) error {
	levels, path, ok := retainedPath(home, root)
	if !ok {
		return nil
	}
	if err := privateDirs(home, levels...); err != nil {
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
	_, path, ok := retainedPath(home, root)
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
