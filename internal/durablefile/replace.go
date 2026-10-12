// Package durablefile replaces one regular file so that a reader sees either the old bytes
// or the new bytes, and a successful return includes synchronization of the file and of
// its parent directory.
package durablefile

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// File is the temporary file operations that a replacement uses.
type File interface {
	Write(p []byte) (int, error)
	Chmod(mode fs.FileMode) error
	Sync() error
	Close() error
	Name() string
}

// Directory is the parent directory operations that a replacement uses after the rename.
type Directory interface {
	Sync() error
	Close() error
}

// Replacer holds the operations of one replacement. Native supplies the operating-system
// operations, and a test supplies its own operations to the same replacement sequence.
type Replacer struct {
	CreateTemp    func(dir, pattern string) (File, error)
	Rename        func(oldpath, newpath string) error
	OpenDirectory func(name string) (Directory, error)
	Remove        func(name string) error
}

// Error is a failed replacement of Path. Stage names the step that failed, and Err holds
// its cause with any cleanup failure joined. Published reports that the rename completed,
// so the new bytes may already be visible; nothing rolls them back.
type Error struct {
	Path      string
	Stage     string
	Published bool
	Err       error
}

// Error names the path, the stage, and the cause, and a published failure adds the
// uncertainty of the visible bytes.
func (e *Error) Error() string {
	message := fmt.Sprintf("replace %s: %s: %v", e.Path, e.Stage, e.Err)
	if e.Published {
		message += "; the rename completed, so new bytes may already be visible"
	}
	return message
}

// Unwrap returns the cause of the failure.
func (e *Error) Unwrap() error { return e.Err }

// Native returns the replacer that uses the operating-system operations.
func Native() Replacer {
	return Replacer{
		CreateTemp: func(dir, pattern string) (File, error) {
			file, err := os.CreateTemp(dir, pattern)
			if err != nil {
				return nil, err
			}
			return file, nil
		},
		Rename: os.Rename,
		OpenDirectory: func(name string) (Directory, error) {
			dir, err := os.Open(name)
			if err != nil {
				return nil, err
			}
			return dir, nil
		},
		Remove: os.Remove,
	}
}

// Replace replaces the file at path with data at mode through the native operations. The
// parent directory of path must exist.
func Replace(path string, data []byte, mode fs.FileMode) error {
	return Native().Replace(path, data, mode)
}

// Replace replaces the file at path with data at mode through the operations of r. It
// creates a temporary regular file beside path, applies mode, writes data, synchronizes
// and closes the file, renames it over path, and then synchronizes and closes the parent
// directory. A failure before the rename leaves path as it was and removes the temporary
// file.
func (r Replacer) Replace(path string, data []byte, mode fs.FileMode) error {
	dir := filepath.Dir(path)
	temp, err := r.CreateTemp(dir, "."+filepath.Base(path)+".*")
	if err != nil {
		return &Error{Path: path, Stage: "create", Err: err}
	}
	stage, err := "mode", temp.Chmod(mode)
	if err == nil {
		stage, err = "write", write(temp, data)
	}
	if err == nil {
		stage, err = "file-sync", temp.Sync()
	}
	if err != nil {
		err = errors.Join(err, temp.Close())
	} else {
		stage, err = "file-close", temp.Close()
	}
	if err == nil {
		stage, err = "rename", r.Rename(temp.Name(), path)
	}
	if err != nil {
		return &Error{Path: path, Stage: stage, Err: errors.Join(err, r.Remove(temp.Name()))}
	}
	directory, err := r.OpenDirectory(dir)
	if err != nil {
		return &Error{Path: path, Stage: "directory-open", Published: true, Err: err}
	}
	stage, err = "directory-sync", directory.Sync()
	if err != nil {
		err = errors.Join(err, directory.Close())
	} else {
		stage, err = "directory-close", directory.Close()
	}
	if err != nil {
		return &Error{Path: path, Stage: stage, Published: true, Err: err}
	}
	return nil
}

// write writes all of data to file. A short write without an error is io.ErrShortWrite.
func write(file File, data []byte) error {
	n, err := file.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	return err
}
