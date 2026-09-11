package assessment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/poolkey"
)

// Store keeps assessment records outside the disposable assignment pool.
type Store struct {
	Home, Root string
	Files      FileOps
}

type FileOps struct {
	WriteFile func(string, []byte, os.FileMode) error
	Rename    func(string, string) error
}

var safeID = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,127}$`)

func (s Store) Dir() string { return filepath.Join(s.Home, "assessment", poolkey.Key(s.Root)) }
func (s Store) path(id string) (string, error) {
	if !safeID.MatchString(id) {
		return "", fmt.Errorf("unsafe run ID")
	}
	return filepath.Join(s.Dir(), id+".json"), nil
}

func noLinks(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	for {
		info, err := os.Lstat(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink refused: %s", path)
		}
		parent := filepath.Dir(path)
		if parent == path {
			return nil
		}
		path = parent
	}
}

func readJSON(path string, target any) error {
	if err := noLinks(path); err != nil {
		return err
	}
	read := bounds.ClassifyNoFollow(path)
	if read.State == bounds.StateAbsent {
		return os.ErrNotExist
	}
	if read.State != bounds.StateParsed {
		return fmt.Errorf("%s: %s %s", path, read.State, read.Reason)
	}
	dec := json.NewDecoder(bytes.NewReader(read.Data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return err
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("trailing JSON input")
	}
	return nil
}

func (s Store) Read(id string) (Run, error) {
	var r Run
	path, err := s.path(id)
	if err != nil {
		return r, err
	}
	if err = readJSON(path, &r); err != nil {
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
	if err = noLinks(path); err != nil {
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
	return replace(s.Files, s.Dir(), path, data)
}
func replace(ops FileOps, dir, path string, data []byte) error {
	write := ops.WriteFile
	if write == nil {
		write = os.WriteFile
	}
	rename := ops.Rename
	if rename == nil {
		rename = os.Rename
	}
	temp, err := os.CreateTemp(dir, ".assessment-*")
	if err != nil {
		return err
	}
	name := temp.Name()
	defer os.Remove(name)
	if err = temp.Close(); err != nil {
		return err
	}
	if err = write(name, data, 0600); err != nil {
		return err
	}
	return rename(name, path)
}
