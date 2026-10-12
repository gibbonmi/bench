package durablefile_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// beforeRename is each stage, in order, that can fail before the rename publishes the new
// bytes.
var beforeRename = []string{"create", "mode", "write", "file-sync", "file-close", "rename"}

// afterRename is each stage, in order, that can fail after the rename publishes the new
// bytes.
var afterRename = []string{"directory-open", "directory-sync", "directory-close"}

// recorder supplies the native operations to a replacement and records in trace each
// stage that the replacement reaches, with "remove" for the removal of a temporary file.
// A stage in fail returns its cause instead of its effect, but a close always releases the
// real file. short makes each write keep its last byte back and report no error. temps
// holds each temporary file name, regular reports that the renamed file was a regular
// file, and opened is the directory that the replacement opened.
type recorder struct {
	fail    map[string]error
	short   bool
	trace   []string
	temps   []string
	regular bool
	opened  string
}

// step records stage and returns its injected cause.
func (rec *recorder) step(stage string) error {
	rec.trace = append(rec.trace, stage)
	return rec.fail[stage]
}

// replacer returns the native replacer with each operation routed through rec.
func (rec *recorder) replacer() durablefile.Replacer {
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
			rec.temps = append(rec.temps, file.Name())
			return &recordedFile{File: file, rec: rec}, nil
		},
		Rename: func(oldpath, newpath string) error {
			if err := rec.step("rename"); err != nil {
				return err
			}
			if info, err := os.Lstat(oldpath); err == nil {
				rec.regular = info.Mode().IsRegular()
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
			rec.opened = name
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

type recordedFile struct {
	durablefile.File
	rec *recorder
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
	if f.rec.short && len(p) > 0 {
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
	rec *recorder
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

// siblings returns the names in the directory of path other than the base of path.
func siblings(t testing.TB, path string) []string {
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
