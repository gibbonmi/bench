package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// TestReadIfExistsPassesAFIFO covers a graded tree that holds a named pipe under a source
// name. The pipe sorts ahead of a stale reference, so the sweep reports that reference
// only when the read passes the pipe without a block in open(2).
func TestReadIfExistsPassesAFIFO(t *testing.T) {
	docs := filepath.Join(t.TempDir(), "docs")
	if err := os.MkdirAll(docs, 0o755); err != nil {
		t.Fatal(err)
	}
	pipe := filepath.Join(docs, "a-pipe.md")
	if err := syscall.Mkfifo(pipe, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docs, "release-runbook.md"), []byte("dist/"+retiredReproducibilityRecord+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	window := bounds.TestDeadline(0)
	done := make(chan []string, 1)
	go func() { done <- checkRetiredReproducibilityRecord(filepath.Dir(docs)) }()
	select {
	case diags := <-done:
		if len(diags) != 1 || !strings.HasSuffix(diags[0], "docs/release-runbook.md") {
			t.Fatalf("sweep past a FIFO: want one diagnostic naming docs/release-runbook.md, got %v", diags)
		}
	case <-time.After(window):
		// A read-write open of a FIFO never blocks, so it releases the blocked reader.
		if release, err := os.OpenFile(pipe, os.O_RDWR, 0); err == nil {
			release.Close()
		}
		t.Fatal(bounds.TestTimeoutVerdict("the retired-record sweep to pass a FIFO", window))
	}
}
