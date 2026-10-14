package responsebound

import (
	"crypto/sha256"
	"encoding/hex"
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

// scopedPath answers the assignment scope directory of the worktree at root and the path of
// the file name in it. It answers false where assignmentDir answers false.
func scopedPath(home, root, name string) (scopeDir, string, bool, error) {
	dir, ok, err := assignmentDir(home, root)
	if err != nil || !ok {
		return dir, "", false, err
	}
	path, err := dir.file(name)
	return dir, path, err == nil, err
}

// Retain replaces the retained response of the assignment worktree at root with the
// complete output of the finished response o. The call function answers the other fields
// of the retained response, and Retain runs it only when the write follows. A response with
// no complete copy clears the retained response. A call function that returns an error
// clears it too. A root outside any assignment worktree keeps nothing.
func (o *Owner) Retain(home, root string, call func() (Retained, error)) error {
	dir, path, ok, err := scopedPath(home, root, retainedName)
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
	var retained Retained
	if err := readScoped(home, root, retainedName, &retained); err != nil {
		return Retained{}, fmt.Errorf("read retained response: %w", err)
	}
	return retained, nil
}

// Staged is one staged excerpt: a copy of the retained response, and the row of a mutant-set
// probe answer that its import reads, or 0 for any other answer.
type Staged struct {
	Retained
	Mutant int `json:"mutant,omitempty"`
}

// stageName is the file name of the staged excerpt of one requirement: the hex SHA-256 digest
// of the spec slug, the chunk ID, and the requirement ID, joined by NUL bytes. No plan string
// forms a path component, so no ID can name a path outside the scope.
func stageName(slug, chunk, requirement string) string {
	sum := sha256.Sum256([]byte(slug + "\x00" + chunk + "\x00" + requirement))
	return hex.EncodeToString(sum[:])
}

// Stage replaces the staged excerpt of one requirement in the assignment scope of the
// worktree at root, so Drop removes it with the scope. A root outside any assignment
// worktree refuses.
func Stage(home, root, slug, chunk, requirement string, staged Staged) error {
	dir, path, ok, err := scopedPath(home, root, stageName(slug, chunk, requirement))
	if err == nil && !ok {
		err = errNoAssignment
	}
	if err == nil {
		err = dir.create()
	}
	if err != nil {
		return fmt.Errorf("stage excerpt: %w", err)
	}
	data, err := json.Marshal(staged)
	if err != nil {
		return fmt.Errorf("stage excerpt: %w", err)
	}
	return durablefile.Replace(path, data, privateFile)
}

// ReadStaged answers the staged excerpt of one requirement in the assignment worktree at
// root. An absent stage, and a root outside any assignment worktree, answer an error that
// matches fs.ErrNotExist.
func ReadStaged(home, root, slug, chunk, requirement string) (Staged, error) {
	var staged Staged
	if err := readScoped(home, root, stageName(slug, chunk, requirement), &staged); err != nil {
		return Staged{}, fmt.Errorf("read staged excerpt: %w", err)
	}
	return staged, nil
}

// Unstage removes the staged excerpt of one requirement in the assignment worktree at root.
// Each store level must be a real directory. An absent stage, and a root outside any
// assignment worktree, remove nothing.
func Unstage(home, root, slug, chunk, requirement string) error {
	dir, path, ok, err := scopedPath(home, root, stageName(slug, chunk, requirement))
	if err == nil && ok {
		err = dir.realLevels("remove", nil)
	}
	if err == nil && ok {
		err = os.Remove(path)
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove staged excerpt: %w", err)
	}
	return nil
}

// errNoAssignment matches fs.ErrNotExist, because a root outside any assignment worktree
// holds no retained response and no stage.
var errNoAssignment = fmt.Errorf("no assignment worktree: %w", fs.ErrNotExist)

// readScoped decodes the file name in the assignment scope of the worktree at root into v.
// It never follows a symlink: each store level must be a real directory, and the file must
// be a regular file.
func readScoped(home, root, name string, v any) error {
	dir, path, ok, err := scopedPath(home, root, name)
	if err != nil {
		return err
	}
	if !ok {
		return errNoAssignment
	}
	if err := dir.realLevels("read", nil); err != nil {
		return err
	}
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return &fs.PathError{Op: "read", Path: path, Err: errNotOwnFile}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

var errNotOwnFile = errors.New("not a regular file")
