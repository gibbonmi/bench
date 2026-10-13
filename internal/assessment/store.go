package assessment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/poolkey"
)

// Store keeps assessment records outside the disposable assignment pool. A nil Files
// replaces each record through the operating-system operations.
type Store struct {
	Home, Root string
	Files      *FileOps
}

// FileOps holds the operations of one record replacement, so a test injects a failure
// into the shared replacement sequence.
type FileOps struct {
	CreateTemp    func(dir, pattern string) (durablefile.File, error)
	Rename        func(oldpath, newpath string) error
	OpenDirectory func(name string) (durablefile.Directory, error)
	Remove        func(name string) error
}

var safeID = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,127}$`)

func (s Store) Dir() string { return filepath.Join(s.Home, "assessment", poolkey.Key(s.Root)) }
func (s Store) path(id string) (string, error) {
	if !ValidID(id) {
		return "", fmt.Errorf("unsafe run ID")
	}
	return filepath.Join(s.Dir(), id+".json"), nil
}

func readJSON(path string, target any) error {
	return decodeFile(path, target, jsonfile.DecodeDocument)
}

func decodeFile(path string, target any, decode func([]byte, any) error) error {
	if err := bounds.RefuseLinks(path); err != nil {
		return err
	}
	read := bounds.ClassifyNoFollow(path)
	if read.State == bounds.StateAbsent {
		return os.ErrNotExist
	}
	if read.State != bounds.StateParsed {
		return fmt.Errorf("%s: %s %s", path, read.State, read.Reason)
	}
	return decode(read.Data, target)
}

// ValidID reports whether value uses the bounded portable identity grammar.
func ValidID(value string) bool { return safeID.MatchString(value) }

func (s Store) Read(id string) (Run, error) {
	var r Run
	path, err := s.path(id)
	if err != nil {
		return r, err
	}
	if err = decodeFile(path, &r, jsonfile.Decode); err != nil {
		return r, err
	}
	if r.RunID != id || r.RepoKey != poolkey.Key(s.Root) {
		return r, fmt.Errorf("record identity mismatch")
	}
	return r, Validate(r)
}

func (s Store) Record(r Run) error {
	if err := Validate(r); err != nil {
		return err
	}
	path, err := s.path(r.RunID)
	if err != nil {
		return err
	}
	if r.RepoKey != poolkey.Key(s.Root) {
		return fmt.Errorf("foreign repository")
	}
	if err = bounds.RefuseLinks(path); err != nil {
		return err
	}
	if err = os.MkdirAll(s.Dir(), 0700); err != nil {
		return err
	}
	lock := path + ".lock"
	if err = os.Mkdir(lock, 0700); err != nil {
		return fmt.Errorf("record locked: %w", err)
	}
	defer os.Remove(lock)
	old, readErr := s.Read(r.RunID)
	if readErr == nil {
		if err = compatible(old, r); err != nil {
			return err
		}
		if reflect.DeepEqual(old, r) {
			return nil
		}
	} else if !os.IsNotExist(readErr) {
		return readErr
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if int64(len(data)) > bounds.ControlRecordLimit {
		return fmt.Errorf("record oversized")
	}
	return replacer(s.Files).Replace(path, data, 0600)
}

func replacer(ops *FileOps) durablefile.Replacer {
	if ops == nil {
		return durablefile.Native()
	}
	return durablefile.Replacer(*ops)
}
