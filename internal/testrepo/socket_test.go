package testrepo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// TestMakeSocketAtLongPath uses a path longer than every platform's socket address
// limit, so a bind at the path itself fails and only the rename can put the socket there.
func TestMakeSocketAtLongPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), strings.Repeat("d", 60), strings.Repeat("s", 60))
	if err := os.Mkdir(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := MakeSocket(path); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("unix sockets unavailable: %v", err))
	}
	info, err := os.Lstat(path)
	if err != nil || info.Mode().Type() != os.ModeSocket {
		t.Fatalf("Lstat(%d-byte path) = %v, %v; want a socket", len(path), info, err)
	}
}
