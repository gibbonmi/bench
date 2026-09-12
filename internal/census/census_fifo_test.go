package census

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

// TestCountsRefusesAFifoWithoutBlocking proves a refused file type reads as zero and
// never holds the board open on a reader that has no writer.
func TestCountsRefusesAFifoWithoutBlocking(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	dir := Dir(home, root)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(dir, knownID), 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	done := make(chan map[string]int, 1)
	go func() {
		counts, _ := Counts(home, root)
		done <- counts
	}()
	window := bounds.TestDeadline(0)
	select {
	case counts := <-done:
		if len(counts) != 0 {
			t.Fatalf("Counts on a FIFO record = %v, want none", counts)
		}
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("Counts to return for a FIFO record file", window))
	}
}

// TestHeadBreakdownRefusesAFifoWithoutBlocking proves the breakdown has the same
// file-type posture as Counts: a refused file type renders no text and never holds
// the landing open on a reader that has no writer.
func TestHeadBreakdownRefusesAFifoWithoutBlocking(t *testing.T) {
	t.Parallel()
	home, root, _ := fixtureHome(t)
	dir := Dir(home, root)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(dir, knownID), 0o600); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	done := make(chan string, 1)
	go func() {
		done <- HeadBreakdown(home, root, knownID)
	}()
	window := bounds.TestDeadline(0)
	select {
	case got := <-done:
		if got != "" {
			t.Fatalf("HeadBreakdown on a FIFO record = %q, want no text", got)
		}
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("HeadBreakdown to return for a FIFO record file", window))
	}
}
