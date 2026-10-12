package durablefile_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
)

// destination returns a destination path in a fresh directory whose names hold spaces,
// embedded quotes, glob characters, and Unicode. When old is not nil, the destination
// holds old.
func destination(t *testing.T, old []byte) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), `parent "quoted" [glob]*? ünïcödé`)
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, `record 'one' *.md`)
	if old != nil {
		if err := os.WriteFile(path, old, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return path
}

// readBack returns the bytes at path.
func readBack(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestReplaceBytes(t *testing.T) {
	path := destination(t, []byte("old bytes that are longer than the new payload\n"))
	payload := []byte("new\x00payload\n\twith a NUL, a tab, and ü\n")
	if err := durablefile.Replace(path, payload, 0o644); err != nil {
		t.Fatalf("Replace: %v", err)
	}
	if got := readBack(t, path); !bytes.Equal(got, payload) {
		t.Fatalf("destination = %q, want %q", got, payload)
	}
}

func TestReplaceEmpty(t *testing.T) {
	path := destination(t, nil)
	if err := durablefile.Replace(path, []byte{}, 0o644); err != nil {
		t.Fatalf("Replace: %v", err)
	}
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("an empty payload left no destination: %v", err)
	}
	if !info.Mode().IsRegular() || info.Size() != 0 {
		t.Fatalf("destination mode %v size %d, want an empty regular file", info.Mode(), info.Size())
	}
}

func TestReplaceMode(t *testing.T) {
	for _, mode := range []fs.FileMode{0o755, 0o640} {
		path := destination(t, []byte("old"))
		if err := durablefile.Replace(path, []byte("#!/bin/sh\n"), mode); err != nil {
			t.Fatalf("Replace at %v: %v", mode, err)
		}
		info, err := os.Lstat(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != mode {
			t.Fatalf("destination permission bits = %v, want %v", got, mode)
		}
	}
}

func TestReplaceVisibility(t *testing.T) {
	payloads := [][]byte{bytes.Repeat([]byte("a"), 1<<20), bytes.Repeat([]byte("b"), 1<<20)}
	path := destination(t, payloads[0])
	stop, partial, reads := make(chan struct{}), make(chan []byte, 1), make(chan int, 1)
	go func() {
		count := 0
		defer func() { reads <- count }()
		for {
			select {
			case <-stop:
				return
			default:
			}
			data, err := os.ReadFile(path)
			count++
			if err != nil || !bytes.Equal(data, payloads[0]) && !bytes.Equal(data, payloads[1]) {
				partial <- data
				return
			}
		}
	}()
	for i := 1; i <= 32; i++ {
		if err := durablefile.Replace(path, payloads[i%2], 0o644); err != nil {
			t.Fatalf("Replace %d: %v", i, err)
		}
	}
	close(stop)
	count := <-reads
	select {
	case data := <-partial:
		t.Fatalf("a concurrent reader observed %d bytes that are neither complete payload", len(data))
	default:
	}
	if count == 0 {
		t.Fatal("the concurrent reader made no read")
	}
}

func TestAdjacentTemporary(t *testing.T) {
	path := destination(t, []byte("old"))
	rec := &recorder{}
	if err := rec.replacer().Replace(path, []byte("new"), 0o644); err != nil {
		t.Fatalf("Replace: %v", err)
	}
	if len(rec.temps) != 1 || filepath.Dir(rec.temps[0]) != filepath.Dir(path) {
		t.Fatalf("temporary files %q, want one file in %q", rec.temps, filepath.Dir(path))
	}
	if !rec.regular {
		t.Fatalf("temporary file %q was not a regular file at the rename", rec.temps[0])
	}
}

func TestReplaceOrder(t *testing.T) {
	path := destination(t, []byte("old"))
	rec := &recorder{}
	if err := rec.replacer().Replace(path, []byte("new"), 0o644); err != nil {
		t.Fatalf("Replace: %v", err)
	}
	want := beforeRename
	if len(rec.trace) < len(want) || !slices.Equal(rec.trace[:len(want)], want) {
		t.Fatalf("trace = %q, want the prefix %q", rec.trace, want)
	}
}

func TestReplaceDirectorySync(t *testing.T) {
	path := destination(t, []byte("old"))
	rec := &recorder{}
	if err := rec.replacer().Replace(path, []byte("new"), 0o644); err != nil {
		t.Fatalf("Replace: %v", err)
	}
	rename := slices.Index(rec.trace, "rename")
	want := afterRename
	if rename < 0 || !slices.Equal(rec.trace[rename+1:], want) {
		t.Fatalf("trace = %q, want %q after the rename", rec.trace, want)
	}
	if rec.opened != filepath.Dir(path) {
		t.Fatalf("opened directory %q, want the parent %q", rec.opened, filepath.Dir(path))
	}
}

func TestReplaceShortWrite(t *testing.T) {
	old := []byte("old")
	path := destination(t, old)
	rec := &recorder{short: true}
	err := rec.replacer().Replace(path, []byte("new payload"), 0o644)
	if !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("Replace after a short write = %v, want io.ErrShortWrite", err)
	}
	if slices.Contains(rec.trace, "rename") {
		t.Fatalf("trace = %q renames after a short write", rec.trace)
	}
	if got := readBack(t, path); !bytes.Equal(got, old) {
		t.Fatalf("destination = %q, want the old %q", got, old)
	}
}

// failAt replaces path with data and fails stage with a fresh cause. It returns the
// replacement error as a *durablefile.Error and the cause.
func failAt(t *testing.T, path string, data []byte, stage string) (*durablefile.Error, error) {
	t.Helper()
	cause := fmt.Errorf("injected %s failure", stage)
	rec := &recorder{fail: map[string]error{stage: cause}}
	err := rec.replacer().Replace(path, data, 0o644)
	var replaceErr *durablefile.Error
	if !errors.As(err, &replaceErr) {
		t.Fatalf("Replace with a %s failure = %v, want a *durablefile.Error", stage, err)
	}
	if replaceErr.Stage != stage {
		t.Fatalf("Stage = %q, want %q", replaceErr.Stage, stage)
	}
	return replaceErr, cause
}

func TestReplaceBeforeRenameFailures(t *testing.T) {
	old := []byte("old bytes")
	for _, stage := range beforeRename {
		path := destination(t, old)
		if replaceErr, _ := failAt(t, path, []byte("new"), stage); replaceErr.Published {
			t.Fatalf("a %s failure reports a published replacement", stage)
		}
		if got := readBack(t, path); !bytes.Equal(got, old) {
			t.Fatalf("after a %s failure the destination = %q, want the old %q", stage, got, old)
		}
	}
}

func TestReplaceAbsentFailures(t *testing.T) {
	for _, stage := range beforeRename {
		path := destination(t, nil)
		failAt(t, path, []byte("new"), stage)
		if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("after a %s failure the absent destination is present: %v", stage, err)
		}
	}
}

func TestReplaceAfterRenameFailures(t *testing.T) {
	for _, stage := range afterRename {
		replaceErr, _ := failAt(t, destination(t, []byte("old")), []byte("new"), stage)
		if !replaceErr.Published || !strings.Contains(replaceErr.Error(), "new bytes may already be visible") {
			t.Fatalf("a %s failure = %q with Published %v, want the published uncertainty", stage, replaceErr, replaceErr.Published)
		}
	}
}

func TestReplaceAfterRenameState(t *testing.T) {
	path, payload := destination(t, []byte("old")), []byte("new payload")
	failAt(t, path, payload, "directory-sync")
	if got := readBack(t, path); !bytes.Equal(got, payload) {
		t.Fatalf("after a directory-sync failure the destination = %q, want the new %q", got, payload)
	}
}

func TestReplaceErrorCause(t *testing.T) {
	for _, stage := range slices.Concat(beforeRename, afterRename) {
		replaceErr, cause := failAt(t, destination(t, []byte("old")), []byte("new"), stage)
		if !errors.Is(replaceErr, cause) {
			t.Fatalf("a %s failure = %v, which does not hold the injected cause", stage, replaceErr)
		}
	}
}

func TestReplaceTemporaryCleanup(t *testing.T) {
	for _, stage := range beforeRename[1:] {
		path := destination(t, []byte("old"))
		failAt(t, path, []byte("new"), stage)
		if left := siblings(t, path); len(left) != 0 {
			t.Fatalf("after a %s failure the directory holds %q beside the destination", stage, left)
		}
	}
}

func TestReplaceCleanupFailure(t *testing.T) {
	for _, stages := range [][2]string{{"write", "remove"}, {"write", "file-close"}, {"directory-sync", "directory-close"}} {
		path := destination(t, []byte("old"))
		primary, cleanup := fmt.Errorf("injected %s failure", stages[0]), fmt.Errorf("injected %s failure", stages[1])
		rec := &recorder{fail: map[string]error{stages[0]: primary, stages[1]: cleanup}}
		err := rec.replacer().Replace(path, []byte("new"), 0o644)
		if !errors.Is(err, primary) || !errors.Is(err, cleanup) {
			t.Fatalf("Replace = %v, want both the %s cause and the %s cleanup cause", err, stages[0], stages[1])
		}
	}
}
