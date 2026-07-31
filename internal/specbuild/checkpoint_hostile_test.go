package specbuild

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestCheckpointRejectsHostileReceiptsBeforeMutationOrBlocking(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*testing.T, checkpointFixture) string
	}{
		{"empty", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "empty.json")
			write(t, path, "")
			return path
		}},
		{"oversized", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "oversized.json")
			write(t, path, strings.Repeat("x", int(bounds.ControlRecordLimit)+1))
			return path
		}},
		{"malformed", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "malformed.json")
			write(t, path, "{\n")
			return path
		}},
		{"unreadable", func(t *testing.T, fixture checkpointFixture) string {
			path := writeCheckpointReceipt(t, fixture.receipt, "\n")
			if err := os.Chmod(path, 0); err != nil {
				t.Fatal(err)
			}
			return path
		}},
		{"fifo", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "receipt.fifo")
			if err := syscall.Mkfifo(path, 0o600); err != nil {
				t.Fatal(err)
			}
			return path
		}},
		{"device", func(_ *testing.T, _ checkpointFixture) string { return "/dev/null" }},
		{"socket", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "receipt.sock")
			listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = listener.Close() })
			return path
		}},
		{"regular symlink", func(t *testing.T, fixture checkpointFixture) string {
			target := writeCheckpointReceipt(t, fixture.receipt, "\n")
			path := filepath.Join(t.TempDir(), "receipt-link.json")
			if err := os.Symlink(target, path); err != nil {
				t.Fatal(err)
			}
			return path
		}},
		{"dangling symlink", func(t *testing.T, _ checkpointFixture) string {
			path := filepath.Join(t.TempDir(), "dangling-link.json")
			if err := os.Symlink("missing.json", path); err != nil {
				t.Fatal(err)
			}
			return path
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newCheckpointFixture(t)
			before := checkpointSnapshotFor(t, fixture)
			path := tc.setup(t, fixture)
			done := make(chan error, 1)
			go func() {
				_, err := fixture.service.Checkpoint(t.Context(), "build demo", fixture.assigned.ID, path)
				done <- err
			}()
			select {
			case err := <-done:
				if err == nil {
					t.Fatal("Checkpoint unexpectedly accepted hostile receipt")
				}
			case <-time.After(time.Second):
				t.Fatal("Checkpoint blocked on hostile receipt")
			}
			if after := checkpointSnapshotFor(t, fixture); after != before {
				t.Fatalf("hostile receipt mutated state: before=%#v after=%#v", before, after)
			}
		})
	}
}
