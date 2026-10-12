package reviewrecord

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// failingFile is a native temporary file whose Sync fails with cause.
type failingFile struct {
	durablefile.File
	cause error
}

func (f failingFile) Sync() error { return f.cause }

// failingDirectory is a native parent directory whose Sync fails with cause.
type failingDirectory struct {
	durablefile.Directory
	cause error
}

func (d failingDirectory) Sync() error { return d.cause }

// failingSync returns the native replacement operation with the file sync, or the
// directory sync when directory is true, failing with cause.
func failingSync(cause error, directory bool) replacement {
	replacer := durablefile.Native()
	create, open := replacer.CreateTemp, replacer.OpenDirectory
	if directory {
		replacer.OpenDirectory = func(name string) (durablefile.Directory, error) {
			dir, err := open(name)
			if err != nil {
				return nil, err
			}
			return failingDirectory{dir, cause}, nil
		}
	} else {
		replacer.CreateTemp = func(dir, pattern string) (durablefile.File, error) {
			file, err := create(dir, pattern)
			if err != nil {
				return nil, err
			}
			return failingFile{file, cause}, nil
		}
	}
	return replacer.Replace
}

// recordRoot returns a root whose name holds spaces, quotes, glob characters, and Unicode,
// and the path of the review record of renderSpec under it. When old is not nil, the record
// holds old at mode 0600, so a replacement that keeps the old mode shows.
func recordRoot(t *testing.T, old []byte) (string, string) {
	t.Helper()
	root := filepath.Join(t.TempDir(), `root "quoted" [glob]*? ünïcödé`)
	relative, err := RecordPath(renderSpec)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if old != nil {
		if err := os.WriteFile(path, old, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root, path
}

// rendered is the document that Render makes of record over document.
func rendered(t *testing.T, document []byte, record Record) []byte {
	t.Helper()
	data, err := Render(document, record)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// heldBytes returns the bytes of the record at path, and it fails the test when the
// directory of path holds any other file.
func heldBytes(t *testing.T, path string) []byte {
	t.Helper()
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Name() != filepath.Base(path) {
			t.Fatalf("the record directory holds %q beside the record", entry.Name())
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestReviewReplacement(t *testing.T) {
	newRecord := renderRecord("new")
	update := func(record *Record) error { *record = newRecord; return nil }
	old := rendered(t, nil, renderRecord("old"))

	t.Run("directory sync failure after the rename", func(t *testing.T) {
		root, path := recordRoot(t, old)
		cause := errors.New("injected directory-sync failure")
		err := writeWith(failingSync(cause, true), root, renderSpec, absent(""), update)
		var replaceErr *durablefile.Error
		if !errors.As(err, &replaceErr) || replaceErr.Stage != "directory-sync" || !replaceErr.Published || !errors.Is(err, cause) {
			t.Fatalf("write with a failed directory sync = %v, want the published directory-sync cause", err)
		}
		if !strings.Contains(err.Error(), "new bytes may already be visible") {
			t.Fatalf("write with a failed directory sync = %q, want the uncertainty phrase", err)
		}
		if got, want := heldBytes(t, path), rendered(t, old, newRecord); !bytes.Equal(got, want) {
			t.Fatalf("record = %q, want the rendered new record %q", got, want)
		}
		if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0o644 {
			t.Fatalf("record mode = %v (%v), want 0644", info, err)
		}
		if loaded, err := Read(root, renderSpec); err != nil || loaded.Completion.Verification[0].NativeRef.Excerpt != "new" {
			t.Fatalf("a reader of the record gets %v, %v; want the new record", loaded, err)
		}
	})

	t.Run("file sync failure before the rename", func(t *testing.T) {
		cause := errors.New("injected file-sync failure")
		for _, existing := range [][]byte{old, nil} {
			root, path := recordRoot(t, existing)
			create := func(string) (Record, error) { return newRecord, nil }
			err := writeWith(failingSync(cause, false), root, renderSpec, create, update)
			var replaceErr *durablefile.Error
			if !errors.As(err, &replaceErr) || replaceErr.Stage != "file-sync" || replaceErr.Published || !errors.Is(err, cause) {
				t.Fatalf("write with a failed file sync = %v, want the unpublished file-sync cause", err)
			}
			if existing == nil {
				if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
					t.Fatalf("a failed first write left a record: %v", err)
				}
				continue
			}
			if got := heldBytes(t, path); !bytes.Equal(got, old) {
				t.Fatalf("record = %q, want the old record", got)
			}
		}
	})

	t.Run("refusal before the replacement", func(t *testing.T) {
		refused := errors.New("injected change refusal")
		for _, change := range []func(*Record) error{
			func(*Record) error { return refused },
			func(record *Record) error { *record = newRecord; record.Spec = "specs/other/spec.md"; return nil },
		} {
			root, path := recordRoot(t, old)
			calls := 0
			replaceFile := func(string, []byte, fs.FileMode) error { calls++; return nil }
			if err := writeWith(replaceFile, root, renderSpec, absent(""), change); err == nil || calls != 0 {
				t.Fatalf("write = %v after %d replacements, want a refusal before any replacement", err, calls)
			}
			if got := heldBytes(t, path); !bytes.Equal(got, old) {
				t.Fatalf("record = %q, want the old record", got)
			}
		}
	})
}
