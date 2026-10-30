package responsebound

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

// TestExactFileStates grades the file owner's states: a present empty file pages explicitly,
// and a linked component, a special or unreadable file, and a missing file refuse before any
// content without opening a linked target.
func TestExactFileStates(t *testing.T) {
	dir := t.TempDir()
	write := func(path, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(dir, "empty.txt"), "")
	write(filepath.Join(dir, "oversized.txt"), strings.Repeat("x", int(bounds.ControlRecordLimit)+1))
	write(filepath.Join(dir, "real", "file.txt"), "linked\n")
	write(filepath.Join(dir, "locked.txt"), "locked\n")
	if err := os.Chmod(filepath.Join(dir, "locked.txt"), 0); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(dir, "fifo"), 0o600); err != nil {
		t.Fatal(err)
	}
	for link, target := range map[string]string{"linkdir": "real", "leaflink": "fifo"} {
		if err := os.Symlink(filepath.Join(dir, target), filepath.Join(dir, link)); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range []struct {
		name   string
		path   string
		stream bool
		state  bounds.FileState
	}{
		{"present empty", "empty.txt", false, ""},
		{"present empty artifact", "empty.txt", true, ""},
		{"linked parent", "linkdir/file.txt", false, bounds.StateWrongType},
		{"linked leaf to a FIFO", "leaflink", true, bounds.StateWrongType},
		{"FIFO", "fifo", true, bounds.StateWrongType},
		{"directory", "real", false, bounds.StateWrongType},
		{"unreadable", "locked.txt", false, bounds.StateUnreadable},
		{"missing", "missing.txt", false, bounds.StateAbsent},
		{"repository file past the control-record limit", "oversized.txt", false, bounds.StateUnreadable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "unreadable" && os.Geteuid() == 0 {
				capability.Capability(t, capability.Privilege, "root reads a file without permission bits")
			}
			source, err := OpenSource(filepath.Join(dir, tc.path), tc.stream)
			if tc.state != "" {
				var refusal *SourceError
				if source != nil || !errors.As(err, &refusal) || refusal.State != tc.state {
					t.Fatalf("OpenSource = %v, %v; want the %s refusal", source, err, tc.state)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer source.Close()
			request := pageRequest(readFrame, nil)
			request.Payload, request.Size, request.Digest = source, source.Size, source.Digest
			page, err := PageRead(request)
			if err != nil {
				t.Fatal(err)
			}
			decoded := decodePage(t, page.Output)
			if cell(decoded.row, "present") != "true" || cell(decoded.row, "next") != "0" || cell(decoded.row, "complete") != "true" || len(decoded.data) != 0 || decoded.argv != nil {
				t.Fatalf("empty page = %q, want present, terminal, empty content at position zero", page.Output)
			}
		})
	}
}

// TestExactSourceReplaced grades the change check after the open: a file that a rename
// replaces or a write rewrites refuses with the changed-file reason.
func TestExactSourceReplaced(t *testing.T) {
	for name, change := range map[string]func(t *testing.T, path string){
		"renamed over": func(t *testing.T, path string) {
			other := path + ".new"
			if err := os.WriteFile(other, []byte("other bytes\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.Rename(other, path); err != nil {
				t.Fatal(err)
			}
		},
		"rewritten": func(t *testing.T, path string) {
			if err := os.WriteFile(path, []byte("original bytes and more\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "source.txt")
			if err := os.WriteFile(path, []byte("original bytes\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			source, err := OpenSource(path, false)
			if err != nil {
				t.Fatal(err)
			}
			defer source.Close()
			if err := source.Unchanged(); err != nil {
				t.Fatalf("an unchanged source refused: %v", err)
			}
			change(t, path)
			var refusal *SourceError
			if err := source.Unchanged(); !errors.As(err, &refusal) || refusal.State != bounds.StateUnreadable || refusal.Reason != changedReason {
				t.Fatalf("Unchanged = %v, want the changed-file refusal", err)
			}
		})
	}
}

// TestExactLineRanges grades the inclusive one-based physical-line selection with its
// original terminators, and its refusal of a range beyond the last line.
func TestExactLineRanges(t *testing.T) {
	data := []byte("first\r\nsecond\nthird")
	for _, tc := range []struct {
		value, want string
		refused     bool
	}{
		{value: "1:1", want: "first\r\n"},
		{value: "2:3", want: "second\nthird"},
		{value: "3:3", want: "third"},
		{value: "3:4", refused: true},
	} {
		first, last, err := ParseLineRange(tc.value)
		if err != nil {
			t.Fatal(err)
		}
		got, err := SelectLines(data, first, last)
		if tc.refused != (err != nil) || string(got) != tc.want {
			t.Errorf("SelectLines(%s) = %q, %v; want %q refused=%t", tc.value, got, err, tc.want, tc.refused)
		}
	}
	if _, err := SelectLines(nil, 1, 1); err == nil {
		t.Error("an empty source selected a line")
	}
	for _, value := range []string{"0:1", "2:1", "1", "1:", "+1:2", "a:b", "1:99999999999999999999"} {
		if _, _, err := ParseLineRange(value); err == nil {
			t.Errorf("ParseLineRange(%q) accepted a malformed range", value)
		}
	}
}
