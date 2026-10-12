// Package durablefiletest supplies the fault-injection harness of a durable replacement to
// the owner tests and to the caller tests.
package durablefiletest

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// Recorder supplies the native operations to a replacement and records in Trace each
// stage that the replacement reaches, with "remove" for the removal of a temporary file.
// A stage in Fail returns its cause instead of its effect, but a close always releases the
// real file. Short makes each write keep its last byte back and report no error. Temps
// holds each temporary file name, Regular reports that the renamed file was a regular
// file, and Opened is the directory that the replacement opened.
type Recorder struct {
	Fail    map[string]error
	Short   bool
	Trace   []string
	Temps   []string
	Regular bool
	Opened  string
}

// step records stage and returns its injected cause.
func (rec *Recorder) step(stage string) error {
	rec.Trace = append(rec.Trace, stage)
	return rec.Fail[stage]
}

// Replacer returns the native replacer with each operation routed through rec.
func (rec *Recorder) Replacer() durablefile.Replacer {
	native := durablefile.Native()
	return durablefile.Replacer{
		CreateTemp: func(dir, pattern string) (durablefile.File, error) {
			if err := rec.step("create"); err != nil {
				return nil, err
			}
			file, err := native.CreateTemp(dir, pattern)
			if err != nil {
				return nil, err
			}
			rec.Temps = append(rec.Temps, file.Name())
			return &recordedFile{File: file, rec: rec}, nil
		},
		Rename: func(oldpath, newpath string) error {
			if err := rec.step("rename"); err != nil {
				return err
			}
			if info, err := os.Lstat(oldpath); err == nil {
				rec.Regular = info.Mode().IsRegular()
			}
			return native.Rename(oldpath, newpath)
		},
		OpenDirectory: func(name string) (durablefile.Directory, error) {
			if err := rec.step("directory-open"); err != nil {
				return nil, err
			}
			dir, err := native.OpenDirectory(name)
			if err != nil {
				return nil, err
			}
			rec.Opened = name
			return &recordedDirectory{Directory: dir, rec: rec}, nil
		},
		Remove: func(name string) error {
			if err := rec.step("remove"); err != nil {
				return err
			}
			return native.Remove(name)
		},
	}
}

// Siblings returns the names in the directory of path other than the base of path.
func Siblings(t testing.TB, path string) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		if entry.Name() != filepath.Base(path) {
			names = append(names, entry.Name())
		}
	}
	return names
}

type recordedFile struct {
	durablefile.File
	rec *Recorder
}

func (f *recordedFile) Chmod(mode fs.FileMode) error {
	if err := f.rec.step("mode"); err != nil {
		return err
	}
	return f.File.Chmod(mode)
}

func (f *recordedFile) Write(p []byte) (int, error) {
	if err := f.rec.step("write"); err != nil {
		return 0, err
	}
	if f.rec.Short && len(p) > 0 {
		return f.File.Write(p[:len(p)-1])
	}
	return f.File.Write(p)
}

func (f *recordedFile) Sync() error {
	if err := f.rec.step("file-sync"); err != nil {
		return err
	}
	return f.File.Sync()
}

func (f *recordedFile) Close() error {
	err := f.rec.step("file-close")
	if closed := f.File.Close(); err == nil {
		err = closed
	}
	return err
}

type recordedDirectory struct {
	durablefile.Directory
	rec *Recorder
}

func (d *recordedDirectory) Sync() error {
	if err := d.rec.step("directory-sync"); err != nil {
		return err
	}
	return d.Directory.Sync()
}

func (d *recordedDirectory) Close() error {
	err := d.rec.step("directory-close")
	if closed := d.Directory.Close(); err == nil {
		err = closed
	}
	return err
}
