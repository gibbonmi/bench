package prospectiveartifact

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
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

	var replaceErr *durablefile.Error
	if !errors.As(err, &replaceErr) || !replaceErr.Published || replaceErr.Stage != "directory-sync" || !errors.Is(err, cause) {
		t.Fatalf("err = %v, want a published directory-sync *durablefile.Error that wraps %v", err, cause)
	}
	path := filepath.Join(root, RecordName)
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		t.Fatalf("owner record mode = %v, want a regular 0600 file", info.Mode())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"schema":1,"owner_pid":4242,"common_dir":"/repo/.git"}` + "\n"; string(data) != want {
		t.Fatalf("owner record = %q, want %q", data, want)
	}
}
