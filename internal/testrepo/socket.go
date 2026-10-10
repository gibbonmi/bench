package testrepo

import (
	"errors"
	"net"
	"os"
	"path/filepath"
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
