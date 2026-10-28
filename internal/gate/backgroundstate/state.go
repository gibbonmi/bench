// Package backgroundstate owns the records of the background grade: the claim of a running
// grade, the background record of a worktree, and the machine queue with its slot. It imports
// no internal/gate code, so a package that internal/gate imports can read and stop a grade.
//
// An absent record reads as the empty value. A present record that does not read is an
// error, and the caller retains its grade and starts none. Every record is written through
// the durable replacement, so a reader sees the old bytes or the new bytes.
package backgroundstate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/jsonfile"
)

// recordMode is the mode of each record, and dirMode is the mode of each record directory.
const (
	recordMode fs.FileMode = 0o600
	dirMode    fs.FileMode = 0o700
)

// readRecord decodes the record at path into value and validates it with valid. An absent
// record leaves value empty and answers no error.
func readRecord(path string, value any, valid func() error) error {
	read := bounds.Classify(path, bounds.ControlRecordLimit)
	if read.State == bounds.StateAbsent {
		return nil
	}
	if read.State.Failed() {
		return fmt.Errorf("%s is %s: %s", path, read.State, read.Reason)
	}
	if err := jsonfile.Decode(read.Data, value); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if err := valid(); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

// writeRecord installs value with a final newline as the private record at path through
// replacer. The parent directory of path must exist.
func writeRecord(replacer durablefile.Replacer, path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return replacer.Replace(path, append(data, '\n'), recordMode)
}

// ensureDir creates dir with dirMode when it is absent, and it refuses a dir that is not a
// private directory. A fresh directory is made durable through its parent.
func ensureDir(dir string) error {
	created := true
	if err := os.Mkdir(dir, dirMode); errors.Is(err, fs.ErrExist) {
		created = false
	} else if err != nil {
		return err
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode().Perm() != dirMode {
		return fmt.Errorf("%s is not a private directory", dir)
	}
	if !created {
		return nil
	}
	parent, err := os.Open(filepath.Dir(dir))
	if err != nil {
		return err
	}
	return errors.Join(parent.Sync(), parent.Close())
}

// pathName refuses a name that is not one path element, so a key never leaves its directory.
func pathName(field, name string) error {
	if name == "" || name == "." || name == ".." || strings.ContainsRune(name, filepath.Separator) {
		return fmt.Errorf("%s %q is not one path name", field, name)
	}
	return nil
}

// required refuses an empty value of each named field.
func required(fields ...string) error {
	for i := 0; i < len(fields); i += 2 {
		if fields[i+1] == "" {
			return fmt.Errorf("the %s field is empty", fields[i])
		}
	}
	return nil
}
