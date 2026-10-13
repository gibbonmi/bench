package reviewrecord_test

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// failing returns the native replacement operation with stage failing with cause.
func failing(stage string, cause error) func(string, []byte, fs.FileMode) error {
	recorder := &durablefiletest.Recorder{Fail: map[string]error{stage: cause}}
	return recorder.Replacer().Replace
}

// recordRoot returns a root whose name holds spaces, quotes, glob characters, and Unicode,
// and the path of the review record of RenderSpec under it. When old is not nil, the record
// holds old at mode 0600, so a replacement that keeps the old mode shows.
func recordRoot(t *testing.T, old []byte) (string, string) {
	t.Helper()
	root := filepath.Join(t.TempDir(), `root "quoted" [glob]*? ünïcödé`)
	relative, err := rr.RecordPath(rr.RenderSpec)
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
func rendered(t *testing.T, document []byte, record rr.Record) []byte {
	t.Helper()
	data, err := rr.Render(document, record)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// heldBytes returns the bytes of the record at path.
func heldBytes(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// requirePublished fails the test unless err is the published directory-sync failure that
// holds cause and states that the new bytes may already be visible.
func requirePublished(t *testing.T, err, cause error) {
	t.Helper()
	var replaceErr *durablefile.Error
	if !errors.As(err, &replaceErr) || replaceErr.Stage != "directory-sync" || !replaceErr.Published || !errors.Is(err, cause) {
		t.Fatalf("write with a failed directory sync = %v, want the published directory-sync cause", err)
	}
	if !strings.Contains(err.Error(), "new bytes may already be visible") {
		t.Fatalf("write with a failed directory sync = %q, want the uncertainty phrase", err)
	}
}

// requireMode fails the test unless the file at path has the permission bits 0644.
func requireMode(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o644 {
		t.Fatalf("%s permission bits = %v, want 0644", path, info.Mode().Perm())
	}
}

func TestReviewReplacement(t *testing.T) {
	newRecord := rr.RenderRecord("new")
	update := func(record *rr.Record) error { *record = newRecord; return nil }
	old := rendered(t, nil, rr.RenderRecord("old"))

	t.Run("record write directory sync failure after the rename", func(t *testing.T) {
		root, path := recordRoot(t, old)
		cause := errors.New("injected directory-sync failure")
		rr.BindReplacement(t, failing("directory-sync", cause))
		requirePublished(t, rr.Write(root, rr.RenderSpec, rr.Absent(""), update), cause)
		if got, want := heldBytes(t, path), rendered(t, old, newRecord); !bytes.Equal(got, want) {
			t.Fatalf("record = %q, want the rendered new record %q", got, want)
		}
		requireMode(t, path)
		if loaded, err := rr.Read(root, rr.RenderSpec); err != nil || loaded.Completion.Verification[0].NativeRef.Excerpt != "new" {
			t.Fatalf("a reader of the record gets %v, %v; want the new record", loaded, err)
		}
	})

	t.Run("assignment write directory sync failure after the rename", func(t *testing.T) {
		f, _ := recordtest.NewLinked(t, 1, recordtest.Delegate, recordtest.Undispatched)
		path := filepath.Join(f.Root, filepath.FromSlash(recordtest.Spec))
		if err := os.Chmod(path, 0o600); err != nil {
			t.Fatal(err)
		}
		cause := errors.New("injected directory-sync failure")
		replaceFile, handed := failing("directory-sync", cause), []byte(nil)
		rr.BindReplacement(t, func(path string, data []byte, mode fs.FileMode) error {
			handed = data
			return replaceFile(path, data, mode)
		})
		source, execution := f.Tip(), *f.Plan.Execution
		_, err := rr.RecordAssignment(f.Root, recordtest.Spec, rr.AssignmentCall{Source: source, Ticket: "1.md", Assignment: recordtest.Assign("1.md"),
			RunID: execution.RunID, OrchestratorSession: execution.OrchestratorSession, AuthorLimit: execution.AuthorLimit, ReviewMode: execution.ReviewMode})
		requirePublished(t, err, cause)
		got, err := os.ReadFile(path)
		if err != nil || handed == nil || !bytes.Equal(got, handed) {
			t.Fatalf("spec = %q (%v), want the replaced plan document %q", got, err, handed)
		}
		requireMode(t, path)
		f.Commit("dispatch 1.md")
		f.Reload()
		want := recordtest.Assign("1.md")
		want.Source = source
		if history := f.Plan.Execution.Assignments["1.md"]; !reflect.DeepEqual(history, []rr.Assignment{want}) {
			t.Fatalf("1.md history = %+v, want %+v", history, want)
		}
	})

	t.Run("file sync failure before the rename", func(t *testing.T) {
		cause := errors.New("injected file-sync failure")
		rr.BindReplacement(t, failing("file-sync", cause))
		for _, existing := range [][]byte{old, nil} {
			root, path := recordRoot(t, existing)
			create := func(string) (rr.Record, error) { return newRecord, nil }
			err := rr.Write(root, rr.RenderSpec, create, update)
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
		refusals := map[string]func(*rr.Record) error{
			"injected change refusal":      func(*rr.Record) error { return errors.New("injected change refusal") },
			"invalid record spec identity": func(record *rr.Record) error { *record = newRecord; record.Spec = "specs/other/spec.md"; return nil },
		}
		calls := 0
		rr.BindReplacement(t, func(string, []byte, fs.FileMode) error { calls++; return nil })
		for refusal, change := range refusals {
			root, path := recordRoot(t, old)
			calls = 0
			if err := rr.Write(root, rr.RenderSpec, rr.Absent(""), change); err == nil || !strings.Contains(err.Error(), refusal) || calls != 0 {
				t.Fatalf("write = %v after %d replacements, want the refusal %q before any replacement", err, calls, refusal)
			}
			if got := heldBytes(t, path); !bytes.Equal(got, old) {
				t.Fatalf("record = %q, want the old record", got)
			}
		}
	})
}
