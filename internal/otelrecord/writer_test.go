package otelrecord

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/poolkey"
)

// TestWriterWritesTheKeyedPathUnderTheHome covers OT8: the record lands at
// otel/<repository key>/traces.jsonl below the resolved home.
func TestWriterWritesTheKeyedPathUnderTheHome(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv(benchhome.Env, home)

	if err := NewWriter(benchhome.Dir(), root).Append([]byte(`{"resourceSpans":[]}`)); err != nil {
		t.Fatalf("append: %v", err)
	}

	want := filepath.Join(home, "otel", poolkey.Key(root), "traces.jsonl")
	text, err := os.ReadFile(want)
	if err != nil {
		t.Fatalf("read %s: %v", want, err)
	}
	if string(text) != "{\"resourceSpans\":[]}\n" {
		t.Fatalf("record content = %q", text)
	}
}

// TestTwoWritersLeaveOnlyIntactLines covers OT9: two writers with independent file
// openers append many lines each, and every line parses.
func TestTwoWritersLeaveOnlyIntactLines(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	const perWriter = 200

	var group sync.WaitGroup
	for writer := 0; writer < 2; writer++ {
		group.Add(1)
		go func(writer int) {
			defer group.Done()
			appender := NewWriter(home, root)
			for line := 0; line < perWriter; line++ {
				payload := fmt.Sprintf(`{"resourceSpans":[{"writer":%d,"line":%d,"pad":%q}]}`,
					writer, line, strings.Repeat("x", 4096))
				if err := appender.Append([]byte(payload)); err != nil {
					t.Errorf("append: %v", err)
					return
				}
			}
		}(writer)
	}
	group.Wait()

	text, err := os.ReadFile(Path(home, root))
	if err != nil {
		t.Fatalf("read record: %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(string(text), "\n"), "\n")
	if len(lines) != 2*perWriter {
		t.Fatalf("line count = %d, want %d", len(lines), 2*perWriter)
	}
	for index, line := range lines {
		var decoded map[string]any
		if err := json.Unmarshal([]byte(line), &decoded); err != nil {
			t.Fatalf("line %d does not parse: %v", index, err)
		}
		if _, ok := decoded["resourceSpans"]; !ok {
			t.Fatalf("line %d has no resourceSpans key", index)
		}
	}
}

// TestWriterRefusesASymlinkedDirectory covers OT10: a redirected record directory
// must not write outside the Bench home.
func TestWriterRefusesASymlinkedDirectory(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	outside := t.TempDir()

	dir := Dir(home, root)
	if err := os.MkdirAll(filepath.Dir(dir), 0o700); err != nil {
		t.Fatalf("create parent: %v", err)
	}
	if err := os.Symlink(outside, dir); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	err := NewWriter(home, root).Append([]byte(`{"resourceSpans":[]}`))
	if err == nil {
		t.Fatal("append accepted a symlinked record directory")
	}
	if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("error = %v, want a symlink refusal", err)
	}
	if entries, readErr := os.ReadDir(outside); readErr != nil || len(entries) != 0 {
		t.Fatalf("the write followed the symlink: entries=%v err=%v", entries, readErr)
	}
}

// TestWriterReportsAnUnwritableDirectory covers OT11: the writer swallows nothing.
func TestWriterReportsAnUnwritableDirectory(t *testing.T) {
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root ignores the directory mode")
	}
	home := t.TempDir()
	root := t.TempDir()

	dir := Dir(home, root)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("create record directory: %v", err)
	}
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatalf("make the record directory unwritable: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	if err := NewWriter(home, root).Append([]byte(`{"resourceSpans":[]}`)); err == nil {
		t.Fatal("append reported no error on an unwritable directory")
	}
}

// TestWriterRefusesANonRegularRecordFile covers OT29: a FIFO at the record path would
// block the first append, so each recorded verb would hang. The row opens the FIFO for
// reading in the background only after the append answers, so a passing run leaves no
// blocked descriptor and a regression fails the deadline rather than the assertion.
func TestWriterRefusesANonRegularRecordFile(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()

	dir := Dir(home, root)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("create record directory: %v", err)
	}
	record := filepath.Join(dir, "traces.jsonl")
	if err := syscall.Mkfifo(record, 0o600); err != nil {
		t.Fatalf("create fifo: %v", err)
	}

	answered := make(chan error, 1)
	go func() { answered <- NewWriter(home, root).Append([]byte(`{"resourceSpans":[]}`)) }()
	select {
	case err := <-answered:
		if err == nil {
			t.Fatal("append accepted a FIFO at the record path")
		}
		if !strings.Contains(err.Error(), "not a regular file") {
			t.Fatalf("error = %v, want a non-regular-file refusal", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("append blocked on the FIFO at the record path")
	}
}

// TestWriterRefusesASymlinkedParent covers OT10 at the parent level: os.MkdirAll follows
// a link above the leaf, so a redirected <home>/otel writes the record outside the home.
func TestWriterRefusesASymlinkedParent(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	outside := t.TempDir()

	if err := os.Symlink(outside, filepath.Join(home, "otel")); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	err := NewWriter(home, root).Append([]byte(`{"resourceSpans":[]}`))
	if err == nil {
		t.Fatal("append accepted a symlinked record parent")
	}
	if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("error = %v, want a symlink refusal", err)
	}
	if entries, readErr := os.ReadDir(outside); readErr != nil || len(entries) != 0 {
		t.Fatalf("the write followed the symlink: entries=%v err=%v", entries, readErr)
	}
}
