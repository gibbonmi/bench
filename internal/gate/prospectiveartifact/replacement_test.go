package prospectiveartifact

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
)

// D18: the owner record writer returns the injected directory-sync failure of its
// replacement after the rename, and the record holds the encoded record at mode 0600.
func TestOwnerReplacement(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	cause := errors.New("injected owner directory-sync failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"directory-sync": cause}}

	err := publish(rec.Replacer(), root, Record{Schema: 1, OwnerPID: 4242, CommonDir: "/repo/.git"})

	durablefiletest.RequireFailure(t, err, "directory-sync", true, cause)
	durablefiletest.RequireFile(t, filepath.Join(root, RecordName), 0o600, `{"schema":1,"owner_pid":4242,"common_dir":"/repo/.git"}`+"\n")
}

// D36: a rename failure of the owner record replacement returns an unpublished failure
// that wraps its cause and publishes no record.
func TestOwnerPreRenameFailure(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	cause := errors.New("injected owner rename failure")
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": cause}}

	err := publish(rec.Replacer(), root, Record{Schema: 1, OwnerPID: 4242, CommonDir: "/repo/.git"})

	durablefiletest.RequireFailure(t, err, "rename", false, cause)
	if _, err := os.Lstat(filepath.Join(root, RecordName)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("owner record stat = %v, want no published record", err)
	}
}
