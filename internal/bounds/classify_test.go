package bounds

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
)

// TestClassifyStates grades every state the classifier can reach in one table, so an
// implementation that answers only absent and complete fails the rows it skipped.
// StateUnsupportedSchema is absent by design: no parser runs here, so the only test
// this package could write for it would assert a constant.
func TestClassifyStates(t *testing.T) {
	for _, tt := range []struct {
		name    string
		fixture func(*testing.T, string) string
		limit   int64
		state   FileState
		stream  ReadStatus
		data    string
		reason  bool
	}{
		{
			name:    "absent",
			fixture: func(_ *testing.T, dir string) string { return filepath.Join(dir, "missing.md") },
			limit:   64,
			state:   StateAbsent,
		},
		{
			name:    "empty",
			fixture: func(t *testing.T, dir string) string { return writeFixture(t, dir, "empty.md", "") },
			limit:   64,
			state:   StateEmpty,
			stream:  ReadComplete,
		},
		{
			name:    "parsed",
			fixture: func(t *testing.T, dir string) string { return writeFixture(t, dir, "record.md", "# journal\n") },
			limit:   64,
			state:   StateParsed,
			stream:  ReadComplete,
			data:    "# journal\n",
		},
		{
			name:    "malformed",
			fixture: func(t *testing.T, dir string) string { return writeFixture(t, dir, "binary.md", "head\xff\xfetail") },
			limit:   64,
			state:   StateMalformed,
			stream:  ReadComplete,
			data:    "head\xff\xfetail",
			reason:  true,
		},
		{
			name:    "oversized",
			fixture: func(t *testing.T, dir string) string { return writeFixture(t, dir, "huge.md", "123456") },
			limit:   5,
			state:   StateUnreadable,
			stream:  ReadOversized,
			reason:  true,
		},
		{
			// A socket rather than a FIFO: this table has no deadline of its own, and a
			// type-blind implementation opening a writerless FIFO would hang the whole
			// table instead of failing one row. TestClassifyFIFOWithoutOpen owns that case.
			name: "wrong type",
			fixture: func(t *testing.T, dir string) string {
				return requireSocket(t, filepath.Join(dir, "control.sock"))
			},
			limit:  64,
			state:  StateWrongType,
			reason: true,
		},
		{
			// The one wrong-type fixture that asks nothing of the host: every filesystem
			// makes directories. The socket row above is the richer fixture but sits
			// behind a capability guard, so without this row a host with no unix sockets
			// asserts wrong-type for a file nowhere at all.
			name: "wrong type (directory)",
			fixture: func(t *testing.T, dir string) string {
				return makeDir(t, dir, "learnings.md")
			},
			limit:  64,
			state:  StateWrongType,
			reason: true,
		},
		{
			name: "unreadable",
			fixture: func(t *testing.T, dir string) string {
				path := writeFixture(t, dir, "denied.md", "secret\n")
				requireUnreadable(t, path)
				return path
			},
			limit:  64,
			state:  StateUnreadable,
			reason: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Classify(tt.fixture(t, t.TempDir()), tt.limit)
			if got.State != tt.state {
				t.Fatalf("Classify state = %q, want %q (reason=%q)", got.State, tt.state, got.Reason)
			}
			if got.Stream != tt.stream {
				t.Errorf("Classify stream = %q, want %q", got.Stream, tt.stream)
			}
			if string(got.Data) != tt.data {
				t.Errorf("Classify data = %q, want %q", got.Data, tt.data)
			}
			if tt.reason && got.Reason == "" {
				t.Errorf("Classify reason is empty, want the underlying diagnostic preserved")
			}
			if !tt.reason && got.Reason != "" {
				t.Errorf("Classify reason = %q, want empty for a non-diagnostic state", got.Reason)
			}
		})
	}
}

// TestClassifyAbsentVsEmpty pins the distinction the whole fail-closed posture rests
// on: absence is the only authoritative empty state, so a present-but-empty record
// must not collapse onto it, nor it onto the present one.
func TestClassifyAbsentVsEmpty(t *testing.T) {
	dir := t.TempDir()
	absent := Classify(filepath.Join(dir, "missing.md"), 64)
	empty := Classify(writeFixture(t, dir, "present.md", ""), 64)
	if absent.State != StateAbsent {
		t.Errorf("missing path state = %q, want %q", absent.State, StateAbsent)
	}
	if empty.State != StateEmpty {
		t.Errorf("present empty file state = %q, want %q", empty.State, StateEmpty)
	}
	if absent.State == empty.State {
		t.Fatal("absent and present-but-empty collapsed onto one state")
	}
}

// TestClassifyDanglingSymlink is the trap the classifier exists to close: os.ReadFile
// reports a broken link as ENOENT, so a read-first implementation reports authoritative
// absence for a link whose target someone deleted.
func TestClassifyDanglingSymlink(t *testing.T) {
	path := filepath.Join(t.TempDir(), "journal.md")
	requireSymlink(t, "target-that-does-not-exist.md", path)
	got := Classify(path, 64)
	if got.State != StateUnreadable {
		t.Fatalf("dangling symlink state = %q, want %q", got.State, StateUnreadable)
	}
	if got.Reason == "" {
		t.Error("dangling symlink reason is empty, want the underlying diagnostic preserved")
	}
}

// TestClassifySymlinkFollowed is story 2's companion: a control record behind a live
// link is read, not rejected, so a blanket "any symlink is unreadable" rule fails here.
func TestClassifySymlinkFollowed(t *testing.T) {
	dir := t.TempDir()
	target := writeFixture(t, dir, "target.md", "# linked\n")
	link := filepath.Join(dir, "journal.md")
	requireSymlink(t, target, link)
	got := Classify(link, 64)
	if got.State != StateParsed || string(got.Data) != "# linked\n" {
		t.Fatalf("symlink to a regular file = %q/%q, want %q and the target's bytes", got.State, got.Data, StateParsed)
	}
}

// TestClassifyFIFOWithoutOpen drives a FIFO with no writer under a real deadline. An
// implementation that opens before checking the type blocks in open(2) forever, so it
// fails by expiring the deadline rather than by returning the wrong state.
func TestClassifyFIFOWithoutOpen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pipe.md")
	requireFifo(t, path)
	done := make(chan Classified, 1)
	go func() { done <- Classify(path, 64) }()
	select {
	case got := <-done:
		if got.State != StateWrongType {
			t.Fatalf("FIFO state = %q, want %q", got.State, StateWrongType)
		}
	case <-time.After(TestDeadline(0)):
		t.Fatal("Classify blocked on a FIFO with no writer, so it opened the path before checking its type")
	}
}

// TestClassifyNonRegular covers the non-regular kinds that would read cleanly: /dev/null
// yields EOF immediately and a socket accepts an open, so a type-blind implementation
// reports either as an empty control record.
func TestClassifyNonRegular(t *testing.T) {
	device := "/dev/null"
	if info, err := os.Lstat(device); err != nil || info.Mode()&os.ModeDevice == 0 {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("no device node at %s", device))
	}
	socket := requireSocket(t, filepath.Join(t.TempDir(), "control.sock"))
	for _, path := range []string{device, socket} {
		got := Classify(path, 64)
		if got.State != StateWrongType {
			t.Errorf("Classify(%s) state = %q, want %q", path, got.State, StateWrongType)
		}
	}
}

// TestClassifyPermissionDenied asserts the reason alongside the state, which is what
// forbids the implementation that maps every failed open onto a bare unreadable and
// throws the underlying error away.
func TestClassifyPermissionDenied(t *testing.T) {
	path := writeFixture(t, t.TempDir(), "denied.md", "secret\n")
	requireUnreadable(t, path)
	got := Classify(path, 64)
	if got.State != StateUnreadable {
		t.Fatalf("mode 0o000 state = %q, want %q", got.State, StateUnreadable)
	}
	if got.Reason == "" {
		t.Error("mode 0o000 reason is empty, want the underlying error preserved")
	}
}

// TestClassifyDir is the directory arm of the same contract. An unreadable directory and
// an empty one both yield zero entries, so only the state can tell them apart.
func TestClassifyDir(t *testing.T) {
	root := t.TempDir()
	absent := ClassifyDir(filepath.Join(root, "missing"))
	if absent.State != StateAbsent {
		t.Errorf("missing directory state = %q, want %q", absent.State, StateAbsent)
	}
	empty := ClassifyDir(makeDir(t, root, "empty"))
	if empty.State != StateEmpty || len(empty.Entries) != 0 {
		t.Errorf("empty directory = %q/%d entries, want %q and none", empty.State, len(empty.Entries), StateEmpty)
	}
	populated := makeDir(t, root, "specs")
	writeFixture(t, populated, "ft86.md", "# spec\n")
	if got := ClassifyDir(populated); got.State != StateParsed || len(got.Entries) != 1 {
		t.Errorf("populated directory = %q/%d entries, want %q and one", got.State, len(got.Entries), StateParsed)
	}
	denied := makeDir(t, root, "denied")
	writeFixture(t, denied, "ft86.md", "# spec\n")
	requireUnreadableDir(t, denied)
	got := ClassifyDir(denied)
	if got.State != StateUnreadable {
		t.Fatalf("unreadable directory state = %q, want %q", got.State, StateUnreadable)
	}
	if got.Reason == "" {
		t.Error("unreadable directory reason is empty, want the underlying error preserved")
	}
	if file := ClassifyDir(writeFixture(t, root, "ROADMAP.md", "# roadmap\n")); file.State != StateWrongType {
		t.Errorf("regular file read as a directory = %q, want %q", file.State, StateWrongType)
	}
}

func writeFixture(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
	return path
}

func makeDir(t *testing.T, root, name string) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("make fixture directory %s: %v", path, err)
	}
	return path
}

func requireSymlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}
}

func requireFifo(t *testing.T, path string) {
	t.Helper()
	if err := syscall.Mkfifo(path, 0o644); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
	}
}

// requireUnreadable strips a fixture's permissions and then proves the strip bit, since
// root ignores the mode entirely and would otherwise read the file straight through the
// assertion. The restore is registered before the chmod so it runs ahead of TempDir's
// own removal, which cannot descend into a directory it cannot enter.
func requireUnreadable(t *testing.T, path string) {
	t.Helper()
	t.Cleanup(func() { _ = os.Chmod(path, 0o600) })
	if err := os.Chmod(path, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip permissions: %v", err))
	}
	f, err := os.Open(path)
	if err == nil {
		f.Close()
		capability.Capability(t, capability.Privilege, "mode 0o000 is still readable by this user")
	}
}

func requireSocket(t *testing.T, path string) string {
	t.Helper()
	listener, err := net.Listen("unix", path)
	if err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("unix sockets unavailable on this filesystem: %v", err))
	}
	t.Cleanup(func() { listener.Close() })
	return path
}

func requireUnreadableDir(t *testing.T, path string) {
	t.Helper()
	t.Cleanup(func() { _ = os.Chmod(path, 0o700) })
	if err := os.Chmod(path, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip directory permissions: %v", err))
	}
	if _, err := os.ReadDir(path); err == nil {
		capability.Capability(t, capability.Privilege, "mode 0o000 directory is still readable by this user")
	}
}

// TestClassifyNoFollowStates grades the no-follow form's complete disposition in one
// table. It is separate from TestClassifyStates because the two forms disagree on
// exactly one input — a live symlink — and a shared table could not assert both
// contracts. Every row is a real filesystem object: a synthetic mode seam would let an
// implementation pass the table while still opening the path it refuses.
func TestClassifyNoFollowStates(t *testing.T) {
	for _, tt := range []struct {
		name    string
		fixture func(*testing.T, string) string
		state   FileState
		stream  ReadStatus
		data    string
		reason  bool
	}{
		{
			name:    "absent",
			fixture: func(_ *testing.T, dir string) string { return filepath.Join(dir, "SKILL.md") },
			state:   StateAbsent,
		},
		{
			name:    "empty",
			fixture: func(t *testing.T, dir string) string { return writeFixture(t, dir, "SKILL.md", "") },
			state:   StateEmpty,
			stream:  ReadComplete,
		},
		{
			name:    "regular file",
			fixture: func(t *testing.T, dir string) string { return writeFixture(t, dir, "SKILL.md", "---\nname: x\n---\n") },
			state:   StateParsed,
			stream:  ReadComplete,
			data:    "---\nname: x\n---\n",
		},
		{
			name: "exactly the control-record limit",
			fixture: func(t *testing.T, dir string) string {
				return writeFixture(t, dir, "SKILL.md", strings.Repeat("a", int(ControlRecordLimit)))
			},
			state:  StateParsed,
			stream: ReadComplete,
			data:   strings.Repeat("a", int(ControlRecordLimit)),
		},
		{
			name: "one byte over the control-record limit",
			fixture: func(t *testing.T, dir string) string {
				return writeFixture(t, dir, "SKILL.md", strings.Repeat("a", int(ControlRecordLimit)+1))
			},
			state:  StateUnreadable,
			stream: ReadOversized,
			reason: true,
		},
		{
			name:    "invalid UTF-8",
			fixture: func(t *testing.T, dir string) string { return writeFixture(t, dir, "SKILL.md", "head\xff\xfetail") },
			state:   StateMalformed,
			stream:  ReadComplete,
			data:    "head\xff\xfetail",
			reason:  true,
		},
		{
			name: "live symlink",
			fixture: func(t *testing.T, dir string) string {
				target := writeFixture(t, dir, "target.md", "# linked\n")
				link := filepath.Join(dir, "SKILL.md")
				requireSymlink(t, target, link)
				return link
			},
			state:  StateWrongType,
			reason: true,
		},
		{
			name: "dangling symlink",
			fixture: func(t *testing.T, dir string) string {
				link := filepath.Join(dir, "SKILL.md")
				requireSymlink(t, filepath.Join(dir, "gone.md"), link)
				return link
			},
			state:  StateWrongType,
			reason: true,
		},
		{
			name: "socket",
			fixture: func(t *testing.T, dir string) string {
				return requireSocket(t, filepath.Join(dir, "control.sock"))
			},
			state:  StateWrongType,
			reason: true,
		},
		{
			name:    "directory",
			fixture: func(t *testing.T, dir string) string { return makeDir(t, dir, "SKILL.md") },
			state:   StateWrongType,
			reason:  true,
		},
		{
			name: "device",
			fixture: func(t *testing.T, _ string) string {
				const device = "/dev/null"
				if info, err := os.Lstat(device); err != nil || info.Mode()&os.ModeDevice == 0 {
					capability.Capability(t, capability.Fifo, fmt.Sprintf("no device node at %s", device))
				}
				return device
			},
			state:  StateWrongType,
			reason: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyNoFollow(tt.fixture(t, t.TempDir()))
			if got.State != tt.state {
				t.Fatalf("ClassifyNoFollow state = %q, want %q (reason=%q)", got.State, tt.state, got.Reason)
			}
			if got.Stream != tt.stream {
				t.Errorf("ClassifyNoFollow stream = %q, want %q", got.Stream, tt.stream)
			}
			if string(got.Data) != tt.data {
				t.Errorf("ClassifyNoFollow data length = %d, want %d", len(got.Data), len(tt.data))
			}
			if tt.reason && got.Reason == "" {
				t.Errorf("ClassifyNoFollow reason is empty, want the underlying diagnostic preserved")
			}
			if !tt.reason && got.Reason != "" {
				t.Errorf("ClassifyNoFollow reason = %q, want empty for a non-diagnostic state", got.Reason)
			}
		})
	}
}

// TestClassifyNoFollowFIFOWithoutOpen is the row the table above cannot hold: a FIFO
// with no writer blocks in open(2) forever, so a type-blind implementation hangs the
// whole table instead of failing one row. It fails here by expiring the deadline.
func TestClassifyNoFollowFIFOWithoutOpen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "SKILL.md")
	requireFifo(t, path)
	done := make(chan Classified, 1)
	go func() { done <- ClassifyNoFollow(path) }()
	select {
	case got := <-done:
		if got.State != StateWrongType {
			t.Fatalf("FIFO state = %q, want %q", got.State, StateWrongType)
		}
	case <-time.After(TestDeadline(0)):
		t.Fatal("ClassifyNoFollow blocked on a FIFO with no writer, so it opened the path before checking its type")
	}
}

// TestClassifyFormsDisagreeOnlyOnSymlinks pins the one input the two forms answer
// differently, so a later edit cannot quietly collapse them into one policy: Classify
// still reads a control record behind a live link, and the no-follow form refuses it.
func TestClassifyFormsDisagreeOnlyOnSymlinks(t *testing.T) {
	dir := t.TempDir()
	target := writeFixture(t, dir, "target.md", "# linked\n")
	link := filepath.Join(dir, "SKILL.md")
	requireSymlink(t, target, link)
	if got := Classify(link, ControlRecordLimit); got.State != StateParsed {
		t.Errorf("Classify(live symlink) = %q, want %q: the follow contract other callers depend on changed", got.State, StateParsed)
	}
	if got := ClassifyNoFollow(link); got.State != StateWrongType {
		t.Errorf("ClassifyNoFollow(live symlink) = %q, want %q", got.State, StateWrongType)
	}
}
