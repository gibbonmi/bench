package bounds

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// failingReader fails every read, so a caller reaches ReadFailed without a filesystem.
type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("read denied") }

// TestClassifyBytesGradesABoundedRead pins the byte form's whole disposition: an oversized
// or failed read is unreadable, invalid UTF-8 is malformed, no bytes is empty, and bytes
// are parsed. The file forms grade their own bytes through this function, so a caller that
// holds bytes and no path — a Git index blob — cannot answer a state the file forms do not.
func TestClassifyBytesGradesABoundedRead(t *testing.T) {
	const limit = 8
	for _, tc := range []struct {
		name  string
		body  string
		want  FileState
		data  string
		fails bool
	}{
		{name: "oversized", body: strings.Repeat("a", limit+1), want: StateUnreadable},
		{name: "invalid UTF-8", body: "\xff", want: StateMalformed, data: "\xff"},
		{name: "empty", want: StateEmpty},
		{name: "parsed", body: "ok", want: StateParsed, data: "ok"},
		{name: "failed read", want: StateUnreadable, fails: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			read := Read(strings.NewReader(tc.body), limit)
			if tc.fails {
				read = Read(failingReader{}, limit)
			}
			c := ClassifyBytes(read)
			if c.State != tc.want {
				t.Fatalf("ClassifyBytes state = %q, want %q", c.State, tc.want)
			}
			if c.State.Failed() && c.Reason == "" {
				t.Fatalf("state %q carries no reason", c.State)
			}
			if !c.State.Failed() && string(c.Data) != tc.data {
				t.Fatalf("ClassifyBytes data = %q, want %q", c.Data, tc.data)
			}
			if c.Stream != read.Status {
				t.Fatalf("Stream = %q, want the read's own status %q", c.Stream, read.Status)
			}
		})
	}
}

// TestClassifyBytesAndClassifyNoFollowAgree pins the split: the same bytes under the same
// limit grade the same whether they arrive as a file or as a read result. A byte form that
// grew its own oversized or UTF-8 rule reds here.
func TestClassifyBytesAndClassifyNoFollowAgree(t *testing.T) {
	for _, body := range []string{"", "ok", "\xff"} {
		path := filepath.Join(t.TempDir(), "subject")
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		file := ClassifyNoFollow(path)
		bytesForm := ClassifyBytes(Read(strings.NewReader(body), ControlRecordLimit))
		if file.State != bytesForm.State || string(file.Data) != string(bytesForm.Data) {
			t.Fatalf("body %q: file form = %+v, byte form = %+v", body, file, bytesForm)
		}
	}
}
