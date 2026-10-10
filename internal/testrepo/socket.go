package testrepo

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// MakeSocket leaves a Unix socket special file at path, whatever the length of path.
// A socket address has a platform length limit near 104 bytes, so MakeSocket binds in a
// short private directory under the temporary directory and renames the socket file to
// path. No listener stays open. The rename requires path to be on the filesystem of the
// temporary directory.
func MakeSocket(path string) (err error) {
	dir, err := os.MkdirTemp("", "s")
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, os.RemoveAll(dir)) }()
	bound := filepath.Join(dir, "s")
	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: bound, Net: "unix"})
	if err != nil {
		return err
	}
	listener.SetUnlinkOnClose(false)
	if err := listener.Close(); err != nil {
		return err
	}
	return os.Rename(bound, path)
}

// PlantSocket makes a socket special file at path through MakeSocket. A MakeSocket error
// alone cannot tell a host without Unix sockets from a defect in the helper or the fixture.
// So after a MakeSocket error, PlantSocket binds a socket of its own at a short path. When
// that bind fails, the host cannot bind a Unix socket at all, and PlantSocket returns the
// bind error for the caller's no-socket route. When that bind succeeds, PlantSocket fails t.
func PlantSocket(t testing.TB, path string) error {
	t.Helper()
	err := MakeSocket(path)
	if err == nil {
		return nil
	}
	if bindErr := bindShortSocket(); bindErr != nil {
		return bindErr
	}
	t.Fatalf("MakeSocket(%s) failed on a host that binds Unix sockets: %v", path, err)
	return err
}

// RequireSocket makes a socket special file at path for t. It skips t through the
// capability route only when the host cannot bind a Unix socket at all.
func RequireSocket(t testing.TB, path string) {
	t.Helper()
	if err := PlantSocket(t, path); err != nil {
		skipWithoutSockets(t, err)
	}
}

// skipWithoutSockets skips t through the capability route for a host where a short-path
// bind failed with err.
func skipWithoutSockets(t testing.TB, err error) {
	t.Helper()
	capability.Capability(t, capability.Fifo, fmt.Sprintf("unix sockets unavailable: %v", err))
}

// bindShortSocket binds and closes a listener at a short path. It shares no step with
// MakeSocket, so a defect in MakeSocket cannot make the host look like a host without
// Unix sockets.
func bindShortSocket() (err error) {
	dir, err := os.MkdirTemp("", "p")
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, os.RemoveAll(dir)) }()
	listener, err := net.Listen("unix", filepath.Join(dir, "p"))
	if err != nil {
		return err
	}
	return listener.Close()
}
