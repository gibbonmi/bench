package status

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

func TestStagedSpecCountUsesFactsStatusReader(t *testing.T) {
	root := t.TempDir()
	write := func(path, body string) {
		t.Helper()
		full := filepath.Join(root, "specs", path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("staged/spec.md", "Status: staged\n")
	write("implemented/spec.md", "Status: implemented\n")
	write("fenced/spec.md", "```md\nStatus: staged\n```\n")
	if err := os.MkdirAll(filepath.Join(root, "specs", "missing"), 0o755); err != nil {
		t.Fatal(err)
	}
	fifo := filepath.Join(root, "specs", "fifo", "spec.md")
	if err := os.MkdirAll(filepath.Dir(fifo), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(fifo, 0o644); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	if got, _ := stagedSpecCount(root); got != 1 {
		t.Fatalf("stagedSpecCount = %d, want 1", got)
	}

	fifoRoot := t.TempDir()
	if err := syscall.Mkfifo(filepath.Join(fifoRoot, "specs"), 0o644); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	done := make(chan int, 1)
	go func() {
		count, _ := stagedSpecCount(fifoRoot)
		done <- count
	}()
	window := bounds.TestDeadline(0)
	select {
	case got := <-done:
		if got != 0 {
			t.Fatalf("FIFO specs count = %d, want 0", got)
		}
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("stagedSpecCount to return for a FIFO specs path", window))
	}
}
