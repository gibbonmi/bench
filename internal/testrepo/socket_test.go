package testrepo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMakeSocketAtLongPath uses a path longer than every platform's socket address
// limit, so a bind at the path itself fails and only the rename can put the socket there.
func TestMakeSocketAtLongPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), strings.Repeat("d", 60), strings.Repeat("s", 60))
	if err := os.Mkdir(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	RequireSocket(t, path)
	info, err := os.Lstat(path)
	if err != nil || info.Mode().Type() != os.ModeSocket {
		t.Fatalf("Lstat(%d-byte path) = %v, %v; want a socket", len(path), info, err)
	}
}

// fatalRecorder keeps a Fatalf message instead of ending the test, so a test can see
// that PlantSocket fails its caller.
type fatalRecorder struct {
	testing.TB
	fatal string
}

func (r *fatalRecorder) Helper() {}

func (r *fatalRecorder) Fatalf(format string, args ...any) { r.fatal = fmt.Sprintf(format, args...) }

// TestPlantSocketFailsOnMakeSocketError gives MakeSocket a path in a missing directory.
// The rename fails on a host that binds Unix sockets, and PlantSocket must fail its
// caller instead of reporting a host without sockets.
func TestPlantSocketFailsOnMakeSocketError(t *testing.T) {
	if err := bindShortSocket(); err != nil {
		skipWithoutSockets(t, err)
	}
	recorder := &fatalRecorder{TB: t}
	err := PlantSocket(recorder, filepath.Join(t.TempDir(), "missing", "s"))
	if err == nil || recorder.fatal == "" {
		t.Fatalf("PlantSocket(path in a missing directory) = %v with Fatalf message %q; want an error and a Fatalf call", err, recorder.fatal)
	}
}
